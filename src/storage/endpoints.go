package storage

import (
	"encoding/json"
	"net/http"
	"time"
)

type SecretRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SecretDecryptRequest struct {
	Secret string `json:"secret"`
}

type SecretDecryptResponse struct {
	Password string `json:"password"`
}

type VaultSecretsResponse struct {
	Secrets []VaultSecretResponse `json:"secrets"`
}

type VaultSecretResponse struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Value     string `json:"value"`
	ExpiresIn int64  `json:"expires_in"`
}

func (v *Vault) EndpointEncrypt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		v.EndpointEncryptSecret(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (v *Vault) EndpointDecrypt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		v.EndpointDecryptSecret(w, r)
	case http.MethodDelete:
		v.EndpointDeleteSecret(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (v *Vault) Endpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		v.EndpointGetSecrets(w, r)
	case http.MethodDelete:
		v.EndpointClearVault(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (v *Vault) EndpointDecryptSecret(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	id := r.URL.Query().Get("id")
	var secretRequest SecretDecryptRequest
	err := json.NewDecoder(r.Body).Decode(&secretRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decryptedSecret, err := v.GetPasswordFromSecret(w, id, secretRequest.Secret, v.ServerSecret)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Secret not found"}`))
		return
	}
	// Delete secret from vault
	v.DeleteSecret(id)
	var secretDecryptResponse SecretDecryptResponse
	secretDecryptResponse.Password = decryptedSecret
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(secretDecryptResponse)
}

func (v *Vault) EndpointGetSecrets(w http.ResponseWriter, r *http.Request) {
	vault := v.GetSecrets()
	var vaultSecrets []VaultSecretResponse
	for _, secret := range vault {
		diffNowToExp := time.Until(secret.ExpiresAt).Seconds()
		vaultSecrets = append(vaultSecrets, VaultSecretResponse{
			Id:        secret.Id.String(),
			Username:  secret.Username,
			Value:     string(secret.Value),
			ExpiresIn: int64(diffNowToExp),
		})
	}
	vaultSecretsResponse := VaultSecretsResponse{
		Secrets: vaultSecrets,
	}
	if len(vaultSecrets) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"secrets": []}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vaultSecretsResponse)
}

func (v *Vault) EndpointEncryptSecret(w http.ResponseWriter, r *http.Request) {
	var secretRequest SecretRequest
	err := json.NewDecoder(r.Body).Decode(&secretRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	clientSecret := v.NewSecret(secretRequest.Username, secretRequest.Password, v.ServerSecret)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(clientSecret)
}

func (v *Vault) EndpointClearVault(w http.ResponseWriter, r *http.Request) {
	v.Secrets = []Secret{}
	w.WriteHeader(http.StatusOK)
}

func (v *Vault) EndpointDeleteSecret(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	v.DeleteSecret(id)
	w.WriteHeader(http.StatusOK)
}

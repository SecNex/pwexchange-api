package storage

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/secnex/pwexchange/utils"
	"github.com/secnex/pwexchange/utils/secure"
)

type Secret struct {
	Id        *uuid.UUID
	Salt      []byte
	Value     string
	Username  string
	ExpiresAt time.Time
}

type SecretResponse struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Secret    string `json:"secret"`
	ExpiresAt string `json:"expires_at"`
}

func NewResponse(secret Secret, clientSecret string) SecretResponse {
	return SecretResponse{
		Id:        secret.Id.String(),
		Username:  secret.Username,
		Secret:    clientSecret,
		ExpiresAt: secret.ExpiresAt.Format(time.RFC3339),
	}
}

func GenerateExpirationTime(seconds int) time.Time {
	return time.Now().Add(time.Duration(seconds) * time.Second)
}

func (v *Vault) AddSecret(secret Secret) {
	v.Secrets = append(v.Secrets, secret)
}

func (v *Vault) GetSecrets() []Secret {
	return v.Secrets
}

func (v *Vault) GetSecret(id string) *Secret {
	for _, secret := range v.Secrets {
		if secret.Id.String() == id {
			return &secret
		}
	}
	return nil
}

func (v *Vault) DeleteSecret(id string) {
	var newSecrets []Secret
	for _, secret := range v.Secrets {
		if secret.Id.String() != id {
			newSecrets = append(newSecrets, secret)
		}
	}
	v.Secrets = newSecrets
}
func (v *Vault) NewSecret(username string, password string, serverSecret []byte) SecretResponse {
	id := uuid.New()
	clientSecret := utils.NewRandom(32).String()
	clientSecretBytes := []byte(clientSecret)

	salt := utils.NewSalt()

	encryptionKey := secure.DeriveKey(serverSecret, clientSecretBytes, salt)

	encryptedPassword, err := secure.Encrypt([]byte(password), encryptionKey)
	if err != nil {
		return SecretResponse{}
	}

	secret := Secret{
		Id:        &id,
		Salt:      salt,
		Value:     encryptedPassword,
		Username:  username,
		ExpiresAt: GenerateExpirationTime(60 * 60 * 24),
	}

	v.AddSecret(secret)

	return NewResponse(secret, clientSecret)
}

func (v *Vault) GetPasswordFromSecret(w http.ResponseWriter, id string, clientSecret string, serverSecret []byte) (string, error) {
	secret := v.GetSecret(id)
	if secret == nil {
		return "", fmt.Errorf("secret not found")
	}

	decryptionClientSecretBytes := []byte(clientSecret)
	decryptionKey := secure.DeriveKey(serverSecret, decryptionClientSecretBytes, secret.Salt)

	decryptedPassword, err := secure.Decrypt(secret.Value, decryptionKey)
	if err != nil {
		return "", fmt.Errorf("decryption failed")
	}

	return decryptedPassword, nil
}

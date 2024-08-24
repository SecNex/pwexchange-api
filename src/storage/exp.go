package storage

import "time"

// Go Routine to handle expiration of secrets
func (v *Vault) ExpirationRoutine() {
	for {
		for i, secret := range v.Secrets {
			if secret.ExpiresAt.Before(time.Now()) {
				v.Secrets = append(v.Secrets[:i], v.Secrets[i+1:]...)
			}
		}
		time.Sleep(time.Second)
	}
}

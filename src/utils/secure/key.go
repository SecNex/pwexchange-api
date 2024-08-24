package secure

import "golang.org/x/crypto/argon2"

func DeriveKey(serverSecret, clientSecret, salt []byte) []byte {
	// Kombiniere Server- und Client-Geheimnis
	combinedSecret := append(serverSecret, clientSecret...)

	// Argon2id key derivation
	return argon2.IDKey(combinedSecret, salt, 3, 64*1024, 4, 32)
}

package storage

type Vault struct {
	Secrets      []Secret
	ServerSecret []byte
}

func NewVault(serverSecret []byte) *Vault {
	return &Vault{
		Secrets:      []Secret{},
		ServerSecret: serverSecret,
	}
}

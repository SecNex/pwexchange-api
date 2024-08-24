package utils

import "math/rand"

type Random struct {
	Length int
	Value  string
}

func NewRandom(length int) *Random {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charactersLength := len(characters)
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = characters[rand.Intn(charactersLength)]
	}
	return &Random{Length: length, Value: string(randomString)}
}

func (r Random) String() string {
	return r.Value
}

func (r Random) Bytes() []byte {
	return []byte(r.Value)
}

func NewSalt() []byte {
	return NewRandom(16).Bytes()
}

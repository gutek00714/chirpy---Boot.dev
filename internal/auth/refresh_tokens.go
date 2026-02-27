package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		panic(err)
	}
	return hex.EncodeToString(token)
}

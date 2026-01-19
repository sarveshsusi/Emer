package services

import (

	"crypto/sha256"
	"encoding/hex"
)

func HashToken(t string) string {
	h := sha256.Sum256([]byte(t))
	return hex.EncodeToString(h[:])
}
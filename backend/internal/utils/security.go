// internal/utils/security.go

package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateAPIKeyPtr gera uma chave segura e retorna como ponteiro
func GenerateAPIKeyPtr() *string {
	bytes := make([]byte, 32) // 32 bytes = 64 caracteres hexadecimais
	rand.Read(bytes)
	key := hex.EncodeToString(bytes)
	return &key
}

package tokens

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}

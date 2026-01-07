package tokenhash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

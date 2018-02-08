package helper

import (
	"crypto/sha256"
	"fmt"
)

func CreateID(text string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
}

package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// FileHash for determining hash of file
func FileHash(filePath string) string {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err == nil {
			return hex.EncodeToString(hash.Sum(nil))
		}
	}

	return ""
}

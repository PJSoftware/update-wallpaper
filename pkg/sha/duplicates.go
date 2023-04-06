package sha

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// FileHash for determining hash of file
func FileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
		
	str := hex.EncodeToString(hash.Sum(nil))
	return str, nil
}

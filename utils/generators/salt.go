package generators

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSalt(length int) (string, error) {

	numBytes := length / 4 * 3

	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	salt := base64.URLEncoding.EncodeToString(randomBytes)

	return salt[:length], nil
}

package encryptors

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, salt string) (string, error) {

	passwordWithSalt := append([]byte(password), []byte(salt)...)

	hashedPassword, err := bcrypt.GenerateFromPassword(
		passwordWithSalt,
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashPassword(hashedPassword string, password string, salt string) error {

	passwordWithSalt := append([]byte(password), []byte(salt)...)

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), passwordWithSalt)
}

package crypto_utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GetHash(input string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

package database

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte(""), err
	}
	return hashedPass, nil
}

func ValidatePassword(password string, storedPass []byte) error {
	err := bcrypt.CompareHashAndPassword(storedPass, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func GenerateRandom() (string, error) {
	length := 32
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

package lib

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/pbkdf2"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func EncryptPassword(password string, salt string) (string, error) {
	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		return "", errors.New("invalid salt format")
	}

	dk := pbkdf2.Key([]byte(password), saltBytes, 10000, 32, sha256.New)
    
	return hex.EncodeToString(dk), nil
}
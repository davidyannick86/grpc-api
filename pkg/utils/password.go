package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrorHandler(errors.New("password cannot be empty"), "Password hashing failed")
	}

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", ErrorHandler(err, "Failed to generate salt")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}

func VerifyPassword(password, encodedHash string) error {
	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		return ErrorHandler(errors.New("invalid hash format"), "Password verification failed")
	}

	saltBase64 := parts[0]
	hashPassword64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return ErrorHandler(err, "Failed to decode salt")
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(hashPassword64)
	if err != nil {
		return ErrorHandler(err, "Failed to decode hash")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, uint32(len(hashedPassword)))

	if len(hash) != len(hashedPassword) {
		return ErrorHandler(errors.New("password does not match"), "Password verification failed")
	}

	if subtle.ConstantTimeCompare(hash, hashedPassword) == 1 {
		return nil
	}

	return ErrorHandler(errors.New("password does not match"), "Password verification failed")
}

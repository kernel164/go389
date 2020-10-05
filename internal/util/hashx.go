package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func Sha256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func CheckSha256(data string, hash string) (bool, error) {
	return Sha256(data) == hash, nil
}

func Md5(data string) string {
	h := md5.New()
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func CheckMd5(data string, hash string) (bool, error) {
	return Md5(data) == hash, nil
}

func CheckBcrypt(data string, hash string) (bool, error) {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(hash), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(hash))
	if err != nil {
		return false, err
	}
	return true, nil
}

func Bcrypt(data string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

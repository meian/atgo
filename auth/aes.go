package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/pkg/errors"
)

var (
	key = "fedcba9876543210"
	iv  = []byte(key)
)

func encrypt(data string) (string, error) {
	d := []byte(data)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "failed to create new cipher")
	}
	cipherText := make([]byte, len(d))
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText, d)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decrypt(encryptedBase64 string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "failed to create new cipher")
	}
	plainText := make([]byte, len(cipherText))
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

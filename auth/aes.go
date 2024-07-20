package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/pkg/errors"
)

var (
	key = []byte("fedcba9876543210")
	iv  = key
)

func encrypt(data string) string {
	d := []byte(data)
	block, _ := aes.NewCipher(key)
	cipherText := make([]byte, len(d))
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText, d)
	return base64.StdEncoding.EncodeToString(cipherText)
}

func decrypt(encryptedBase64 string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64")
	}
	block, _ := aes.NewCipher(key)
	plainText := make([]byte, len(cipherText))
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

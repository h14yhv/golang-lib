package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

const (
	ErrCipherTextTooShort = "cypher text too short"
)

func GCMEncrypt(plain, key []byte) (encoded []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// Success
	return gcm.Seal(nonce, nonce, plain, nil), nil
}

func GCMDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New(ErrCipherTextTooShort)
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	// Success
	return gcm.Open(nil, nonce, cipherText, nil)
}

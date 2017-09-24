package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Encrypt ...
func Encrypt(key AESKey, msg []byte) ([]byte, error) {
	// Generate Cipher block
	cipherBlock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	// Generate GCM
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, msg, nil), nil
}

// Decrypt ...
func Decrypt(key AESKey, msg []byte) ([]byte, error) {
	// Generate Cipher block
	cipherBlock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	// Generate GCM
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		panic(err.Error())
	}

	// Check for nonce existence in ciphertext
	if len(msg) < gcm.NonceSize() {
		return nil, errors.New("invalid msg")
	}

	// Obtain plaintext msg
	plaintext, err := gcm.Open(nil,
		msg[:gcm.NonceSize()], msg[gcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

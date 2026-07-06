package github

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"git.shi.foo/config"
)

func encryptionKey() [32]byte {
	return sha256.Sum256([]byte(config.GitHub.EncryptionKey))
}

func Encrypt(plaintext string) (string, error) {
	key := encryptionKey()
	block, blockError := aes.NewCipher(key[:])
	if blockError != nil {
		return "", blockError
	}

	gcm, gcmError := cipher.NewGCM(block)
	if gcmError != nil {
		return "", gcmError
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, randError := io.ReadFull(rand.Reader, nonce); randError != nil {
		return "", randError
	}

	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func Decrypt(encoded string) (string, error) {
	key := encryptionKey()
	block, blockError := aes.NewCipher(key[:])
	if blockError != nil {
		return "", blockError
	}

	gcm, gcmError := cipher.NewGCM(block)
	if gcmError != nil {
		return "", gcmError
	}

	sealed, decodeError := base64.StdEncoding.DecodeString(encoded)
	if decodeError != nil {
		return "", decodeError
	}

	nonceSize := gcm.NonceSize()
	if len(sealed) < nonceSize {
		return "", errors.New(CiphertextTooShort)
	}

	nonce, ciphertext := sealed[:nonceSize], sealed[nonceSize:]
	plaintext, openError := gcm.Open(nil, nonce, ciphertext, nil)
	if openError != nil {
		return "", openError
	}

	return string(plaintext), nil
}

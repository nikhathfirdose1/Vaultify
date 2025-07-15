package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

const KeySize = 32   // 32 bytes for AES-256
const NonceSize = 12 // Standard GCM nonce size

var encrypter cipher.AEAD

// LoadOrCreateKey loads the master key from disk or creates one if not present
func LoadOrCreateKey(path string) error {
	key, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new random key
			key = make([]byte, KeySize)
			_, err := rand.Read(key)
			if err != nil {
				return fmt.Errorf("failed to generate key: %w", err)
			}
			if err := os.WriteFile(path, key, 0600); err != nil {
				return fmt.Errorf("failed to write key: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read key: %w", err)
		}
	}
	if len(key) != KeySize {
		return errors.New("invalid key length: must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher block: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to wrap cipher in GCM: %w", err)
	}

	encrypter = aead
	return nil
}

// Encrypt encrypts plaintext with AES-GCM and returns nonce+ciphertext+tag
func Encrypt(plaintext []byte) ([]byte, error) {
	if encrypter == nil {
		return nil, errors.New("encrypter not initialized")
	}
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	ciphertext := encrypter.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

// Decrypt reverses the encryption
func Decrypt(blob []byte) ([]byte, error) {
	if encrypter == nil {
		return nil, errors.New("encrypter not initialized")
	}
	if len(blob) < NonceSize {
		return nil, errors.New("invalid blob length")
	}
	nonce := blob[:NonceSize]
	ciphertext := blob[NonceSize:]
	return encrypter.Open(nil, nonce, ciphertext, nil)
}

var whitelist []string

func SetWhitelist(tokens []string) {
	whitelist = tokens
}

func WhitelistTokens() []string {
	return whitelist
}

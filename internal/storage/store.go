package storage

import (
	"errors"
	"sync"
	"time"
)

type StoredSecret struct {
	Blob      []byte
	ExpiresAt time.Time
}

type SecretStore struct {
	mu      sync.RWMutex
	secrets map[string]StoredSecret
}

var store *SecretStore

func InitStore() {
	store = &SecretStore{
		secrets: make(map[string]StoredSecret),
	}
	go store.startEvictionLoop()
}

// StoreSecret adds a secret with TTL in seconds
func StoreSecret(name string, blob []byte, ttlSeconds int) {
	expires := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	store.mu.Lock()
	defer store.mu.Unlock()
	store.secrets[name] = StoredSecret{
		Blob:      blob,
		ExpiresAt: expires,
	}
}

// FetchSecret returns the blob if it hasn't expired
func FetchSecret(name string) ([]byte, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	secret, ok := store.secrets[name]
	if !ok {
		return nil, errors.New("secret not found")
	}
	if time.Now().After(secret.ExpiresAt) {
		return nil, errors.New("secret expired")
	}
	return secret.Blob, nil
}

// startEvictionLoop runs every minute to clean expired secrets
func (s *SecretStore) startEvictionLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for name, secret := range s.secrets {
			if now.After(secret.ExpiresAt) {
				delete(s.secrets, name)
			}
		}
		s.mu.Unlock()
	}
}


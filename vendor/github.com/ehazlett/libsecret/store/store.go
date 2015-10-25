package store

import (
	"errors"
	"time"
)

type Backend string

const (
	VAULT Backend = "vault"
)

var (
	// ErrBackendNotSupported is thrown when the backend store is not supported
	ErrBackendNotSupported = errors.New("Backend store not supported yet, please choose one of")
)

type SecretStore interface {
	// Get returns a secret from the store
	Get(path string) (*Secret, error)
	// Put stores a secret in the store
	Put(path string, value interface{}) error
	// Delete removes a secret from the store
	Delete(path string) error
	// Revoke is used to revoke access for the path
	Revoke(path string) error
	// Renew is used to renew a lease for the secret
	Renew(path string, duration time.Duration) error
}

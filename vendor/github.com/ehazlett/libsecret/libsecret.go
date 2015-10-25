package libsecret

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ehazlett/libsecret/store"
)

// Initialize creates a new Store object, initializing the client
type Initialize func(addr string, options *store.Config) (store.SecretStore, error)

var (
	// Backend initializers
	initializers = make(map[store.Backend]Initialize)
)

func supportedBackends() string {
	keys := make([]string, 0, len(initializers))
	for k := range initializers {
		keys = append(keys, string(k))
	}

	sort.Strings(keys)

	return strings.Join(keys, ", ")
}

// NewSecretStore creates a an instance of store
func NewSecretStore(backend store.Backend, addr string, options *store.Config) (store.SecretStore, error) {
	if init, exists := initializers[backend]; exists {
		return init(addr, options)
	}

	return nil, fmt.Errorf("%s %s", store.ErrBackendNotSupported.Error(), supportedBackends())
}

// AddSecretStore adds a new store backend to libkv
func AddSecretStore(store store.Backend, init Initialize) {
	initializers[store] = init
}

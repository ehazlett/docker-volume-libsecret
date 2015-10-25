package vault

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/libsecret"
	"github.com/ehazlett/libsecret/store"
	"github.com/hashicorp/vault/api"
)

var (
	ErrAddressNotSpecified = errors.New("address not specified")
	ErrNotImplemented      = errors.New("not yet implemented")
	ErrSecretDoesNotExist  = errors.New("secret does not exist")
)

type Vault struct {
	client *api.Client
}

// Register registers with libsecret
func Register() {
	libsecret.AddSecretStore(store.VAULT, NewVault)
}

func NewVault(addr string, config *store.Config) (store.SecretStore, error) {
	defaultCfg := api.DefaultConfig()
	defaultCfg.Address = addr

	if addr == "" {
		return nil, ErrAddressNotSpecified
	}

	c, err := api.NewClient(defaultCfg)
	if err != nil {
		return nil, err
	}

	if t, ok := config.StoreOpts["token"]; ok {
		log.Debugf("using token: %s", t)
		c.SetToken(t.(string))
	}

	return &Vault{
		client: c,
	}, nil
}

func (v *Vault) Get(path string) (*store.Secret, error) {
	s, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, ErrSecretDoesNotExist
	}

	secret := &store.Secret{
		Path:  path,
		Value: s.Data["value"],
	}

	return secret, nil
}

func (v *Vault) Put(path string, value interface{}) error {
	data := map[string]interface{}{
		"value": value,
	}

	if _, err := v.client.Logical().Write(path, data); err != nil {
		return err
	}

	return nil
}

func (v *Vault) Delete(path string) error {
	if _, err := v.client.Logical().Delete(path); err != nil {
		return err
	}

	return nil
}

func (v *Vault) Revoke(path string) error {
	return ErrNotImplemented
}

func (v *Vault) Renew(path string, duration time.Duration) error {
	return ErrNotImplemented
}

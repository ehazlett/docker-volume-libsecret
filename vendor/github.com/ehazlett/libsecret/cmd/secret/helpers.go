package main

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/ehazlett/libsecret"
	"github.com/ehazlett/libsecret/store"
)

func backend(b string) store.Backend {
	switch b {
	case "vault":
		return store.VAULT
	}

	return ""
}

func getStore(c *cli.Context) (store.SecretStore, error) {
	b := backend(c.GlobalString("backend"))
	addr := c.GlobalString("addr")
	opts := getStoreOpts(c)

	storeConfig := &store.Config{
		StoreOpts: opts,
	}

	return libsecret.NewSecretStore(b, addr, storeConfig)
}

func getStoreOpts(c *cli.Context) map[string]interface{} {
	opts := c.GlobalStringSlice("store-opt")
	data := map[string]interface{}{}

	for _, o := range opts {
		parts := strings.Split(o, "=")
		if len(parts) > 1 {
			data[parts[0]] = parts[1]
		} else {
			data[parts[0]] = ""
		}
	}

	return data
}

package main

import (
	"os"
	"os/signal"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/calavera/dkvolume"
	"github.com/codegangsta/cli"
	"github.com/ehazlett/docker-volume-libsecret/version"
	"github.com/ehazlett/libsecret/store"

	// supported backends
	"github.com/ehazlett/libsecret/store/vault"
)

func init() {
	log.SetLevel(log.DebugLevel)

	// register vault backend
	vault.Register()
}

func getStoreOpts(c *cli.Context) map[string]interface{} {
	opts := c.StringSlice("store-opt")
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

func getBackend(b string) store.Backend {
	switch b {
	case "vault":
		return store.VAULT
	}

	return ""
}

func main() {
	app := cli.NewApp()
	app.Version = version.FullVersion()
	app.Name = "docker-volume-libsecret"
	app.Usage = "docker volume driver plugin for libsecret"
	app.Author = "@ehazlett"
	app.Email = "github.com/ehazlett/docker-volume-libsecret"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "root, r",
			Usage: "root path for volumes",
			Value: "/var/lib/docker/volumes/libsecret",
		},
		cli.StringFlag{
			Name:  "addr",
			Usage: "secret backend addr",
			Value: "",
		},
		cli.StringFlag{
			Name:  "backend",
			Usage: "secret backend",
			Value: "",
		},
		cli.StringSliceFlag{
			Name:  "store-opt",
			Usage: "options to pass to libsecret store (key=val)",
			Value: &cli.StringSlice{},
		},
	}

	app.Action = func(c *cli.Context) {
		log.Infof("%s v%s", app.Name, app.Version)

		rootPath := c.String("root")
		backendAddr := c.String("addr")
		backendName := c.String("backend")
		opts := getStoreOpts(c)

		if backendAddr == "" {
			log.Fatal("you must specify a backend address")
		}

		if backendName == "" {
			log.Fatal("you must specify a backend type")
		}

		backend := getBackend(backendName)

		log.Debugf("initializing secret driver: backend=%s addr=%s", backend, backendAddr)
		d, err := NewSecretDriver(rootPath, backend, backendAddr, opts)
		if err != nil {
			log.Fatal(err)
		}

		h := dkvolume.NewHandler(d)
		h.ServeUnix("root", "libsecret")

		cs := make(chan os.Signal, 1)
		signal.Notify(cs, os.Interrupt)
		go func() {
			for _ = range cs {
				// cleanup
			}
		}()

	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

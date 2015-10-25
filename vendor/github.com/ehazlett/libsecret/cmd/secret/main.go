package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ehazlett/libsecret/version"

	"github.com/ehazlett/libsecret/store/vault"
	"github.com/ehazlett/simplelog"
)

func init() {
	// simple log formatter
	f := &simplelog.SimpleFormatter{}
	log.SetFormatter(f)

	// register vault store
	vault.Register()
}

func main() {
	app := cli.NewApp()
	app.Name = "secret"
	app.Usage = "libsecret cli"
	app.Version = version.FullVersion()
	app.Author = "@ehazlett"
	app.Email = "github.com/ehazlett/libsecret"
	app.Before = func(c *cli.Context) error {
		// enable debug
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("debug enabled")
		}

		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "backend, b",
			Usage:  "secret store backend",
			EnvVar: "SECRET_BACKEND",
		},
		cli.StringFlag{
			Name:   "addr, a",
			Usage:  "address to backend store",
			EnvVar: "SECRET_ADDRESS",
		},
		cli.StringSliceFlag{
			Name:  "store-opt, o",
			Usage: "secret store option (key=val)",
		},
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug",
		},
	}
	app.Commands = []cli.Command{
		getCommand,
		setCommand,
		rmCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

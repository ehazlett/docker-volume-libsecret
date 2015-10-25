package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var setCommand = cli.Command{
	Name:   "set",
	Action: setAction,
	Before: func(c *cli.Context) error {
		if c.String("path") == "" {
			return fmt.Errorf("you must specify a path")
		}

		if c.String("value") == "" {
			return fmt.Errorf("you must specify a value")
		}

		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Usage: "path to credential",
		},
		cli.StringFlag{
			Name:  "value, v",
			Usage: "credential value",
		},
	},
}

func setAction(c *cli.Context) {
	s, err := getStore(c)
	if err != nil {
		log.Fatal(err)
	}

	path := c.String("path")
	val := c.String("value")

	if err := s.Put(path, val); err != nil {
		log.Fatal(err)
	}
}

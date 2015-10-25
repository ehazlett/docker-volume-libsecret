package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var rmCommand = cli.Command{
	Name:   "rm",
	Action: rmAction,
	Before: func(c *cli.Context) error {
		if c.String("path") == "" {
			return fmt.Errorf("you must specify a path")
		}

		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Usage: "path to credential",
		},
	},
}

func rmAction(c *cli.Context) {
	s, err := getStore(c)
	if err != nil {
		log.Fatal(err)
	}

	path := c.String("path")
	if err := s.Delete(path); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var getCommand = cli.Command{
	Name:   "get",
	Action: getAction,
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

func getAction(c *cli.Context) {
	s, err := getStore(c)
	if err != nil {
		log.Fatal(err)
	}

	path := c.String("path")
	secret, err := s.Get(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(secret.Value)
}

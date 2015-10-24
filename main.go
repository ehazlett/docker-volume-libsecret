package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
	"github.com/calavera/dkvolume"
)

var (
	rootPath string
)

func init() {
	log.SetLevel(log.DebugLevel)

	flag.StringVar(&rootPath, "root-path", "/var/lib/docker/volumes/libsecret", "root path for volumes")
}

func main() {
	flag.Parse()

	d, err := NewSecretDriver(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	h := dkvolume.NewHandler(d)
	h.ServeUnix("root", "libsecret")
}

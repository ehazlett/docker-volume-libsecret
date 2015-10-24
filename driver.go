package main

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/calavera/dkvolume"
)

type SecretDriver struct {
	root string
}

func NewSecretDriver(root string) (*SecretDriver, error) {
	return &SecretDriver{
		root: root,
	}, nil
}

func (d *SecretDriver) resolvePath(name string) string {
	return filepath.Join(d.root, name)
}

func (d *SecretDriver) Create(r dkvolume.Request) dkvolume.Response {
	log.Debugf("create: %v", r)

	p := d.resolvePath(r.Name)

	errStr := ""
	if err := os.MkdirAll(p, 0755); err != nil {
		errStr = err.Error()
	}

	return dkvolume.Response{Err: errStr}
}

func (d *SecretDriver) Remove(r dkvolume.Request) dkvolume.Response {
	log.Debugf("remove: %v", r)

	p := d.resolvePath(r.Name)

	errStr := ""
	if err := os.RemoveAll(p); err != nil {
		errStr = err.Error()
	}

	return dkvolume.Response{Err: errStr}
}

func (d *SecretDriver) Path(r dkvolume.Request) dkvolume.Response {
	log.Debugf("path: %v", r)

	p := d.resolvePath(r.Name)
	return dkvolume.Response{Mountpoint: p}
}

func (d *SecretDriver) Mount(r dkvolume.Request) dkvolume.Response {
	log.Debugf("mount: %v", r)

	return dkvolume.Response{Mountpoint: filepath.Join(d.root, r.Name)}
}

func (d *SecretDriver) Unmount(r dkvolume.Request) dkvolume.Response {
	log.Debugf("unmount: %v", r)

	return dkvolume.Response{}
}

package main

import (
	"os"
	"path/filepath"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/calavera/dkvolume"
	"github.com/ehazlett/libsecret/store"
)

type SecretDriver struct {
	root         string
	fs           map[string]*FS
	storeAddr    string
	storeBackend store.Backend
	storeOpts    map[string]interface{}
}

func NewSecretDriver(root string, backend store.Backend, addr string, opts map[string]interface{}) (*SecretDriver, error) {
	log.Debugf("opts: %v", opts)

	log.Debugf("backend: %s", backend)

	return &SecretDriver{
		root:         root,
		storeAddr:    addr,
		storeBackend: backend,
		storeOpts:    opts,
		fs:           map[string]*FS{},
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

	p := d.resolvePath(r.Name)

	errStr := ""

	fs, err := NewFS(p, d.storeBackend, d.storeAddr, d.storeOpts)
	if err != nil {
		errStr = err.Error()
	}

	if err := fs.Mount(r.Name); err != nil {
		errStr = err.Error()
	}

	d.fs[r.Name] = fs

	return dkvolume.Response{
		Mountpoint: filepath.Join(d.root, r.Name),
		Err:        errStr,
	}
}

func (d *SecretDriver) Unmount(r dkvolume.Request) dkvolume.Response {
	log.Debugf("unmount: %v", r)

	p := d.resolvePath(r.Name)
	if err := syscall.Unmount(p, 0); err != nil {
		log.Fatal(err)
	}

	return dkvolume.Response{}
}

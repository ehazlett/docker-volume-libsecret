package main

import (
	"os"
	"path"
	"strings"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/libsecret"
	"github.com/ehazlett/libsecret/store"
	"golang.org/x/net/context"
)

type FS struct {
	mountpoint string
	volumeName string
	conn       *fuse.Conn
	errChan    chan (error)
	server     *fs.Server
	store      store.SecretStore
	files      map[string]*File
	tick       *time.Ticker
}

func NewFS(mountpoint string, storeBackend store.Backend, storeAddr string, storeOpts map[string]interface{}) (*FS, error) {
	c := make(chan error)
	go func() {
		err := <-c
		log.Errorf("fs: %s", err.Error())
	}()

	storeConfig := &store.Config{
		StoreOpts: storeOpts,
	}

	secretStore, err := libsecret.NewSecretStore(storeBackend, storeAddr, storeConfig)
	if err != nil {
		return nil, err
	}

	return &FS{
		mountpoint: mountpoint,
		errChan:    c,
		store:      secretStore,
		files:      map[string]*File{},
	}, nil
}

func (f *FS) Mount(volumeName string) error {
	log.Debugf("setting up fuse: volume=%s", volumeName)
	c, err := fuse.Mount(
		f.mountpoint,
		fuse.FSName("libsecret"),
		fuse.Subtype("libsecretfs"),
		fuse.LocalVolume(),
		fuse.VolumeName(volumeName),
	)
	if err != nil {
		return err
	}

	srv := fs.New(c, nil)

	f.server = srv
	f.volumeName = volumeName
	f.conn = c

	go func() {
		err = f.server.Serve(f)
		if err != nil {
			f.errChan <- err
		}
	}()

	// check if the mount process has an error to report
	log.Debug("waiting for mount")
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}

	return nil
}

func (f *FS) Root() (fs.Node, error) {
	return &Dir{f, f.server, f.volumeName}, nil
}

type Dir struct {
	fs       *FS
	fuse     *fs.Server
	basePath string
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	a.Valid = time.Second * 1
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	secretPath := strings.TrimPrefix(path.Join(d.basePath, name), "/")
	log.Debugf("looking up secret: path=%s", secretPath)

	s, err := d.fs.store.Get(secretPath)
	if err != nil {
		log.Warn(err)
		return &Dir{d.fs, d.fuse, path.Join(d.basePath, name)}, nil
	}

	f := &File{
		d.fs,
		d.fuse,
		name,
		path.Join(d.basePath, name),
		s.Value.(string),
	}

	return f, nil
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return nil, nil
}

type File struct {
	fs      *FS
	fuse    *fs.Server
	name    string
	path    string
	content string
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Debugf("path: %s", f.path)

	f.fs.files[f.path] = f

	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(f.content))
	a.Valid = time.Second * 1

	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(f.content), nil
}

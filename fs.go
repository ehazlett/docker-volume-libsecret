package main

import (
	"os"
	"path"
	"strings"

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
	store      store.SecretStore
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
	}, nil
}

func (f *FS) Mount(volumeName string) error {
	log.Debugf("setting up fuse: volume=%s containerPath=%s", volumeName)
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

	f.volumeName = volumeName
	f.conn = c

	go func() {
		err = fs.Serve(c, f)
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
	return &Dir{f, f.volumeName}, nil
}

type Dir struct {
	fs       *FS
	basePath string
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	secretPath := strings.TrimPrefix(path.Join(d.basePath, name), "/")
	log.Debugf("looking up secret: path=%s", secretPath)

	s, err := d.fs.store.Get(secretPath)
	if err != nil {
		return &Dir{d.fs, path.Join(d.basePath, name)}, nil
	}

	return &File{
		d.fs,
		name,
		path.Join(d.basePath, name),
		s.Value.(string),
	}, nil
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return nil, nil
}

type File struct {
	fs      *FS
	name    string
	path    string
	content string
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Debugf("path: %s", f.path)

	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(f.content))
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(f.content), nil
}

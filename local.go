package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const defFileMode os.FileMode = os.ModePerm

var (
	ErrEmptyURI = errors.New("empty uri")
)

// Local is local storage in local filesystem
type Local string

func NewLocal(path string) (*Local, error) {
	// TODO(zouying): add check path is valid and safe
	// ...

	// make sure dir exsists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	l := Local(path)

	return &l, nil
}

// Put data and return uri for object
// uri could be multi-level directory, split by '/'
func (l Local) Put(ctx context.Context, uri string, r io.Reader) error {
	if err := validURI(uri); err != nil {
		return err
	}

	path := l.fullpath(uri)

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	if _, err := os.Stat(path); os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, defFileMode)
}

// Get object by uri
func (l Local) Get(ctx context.Context, uri string) (rc io.ReadCloser, err error) {
	if err := validURI(uri); err != nil {
		return nil, err
	}

	path := l.fullpath(uri)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	return os.Open(path)
}

// Delete object by uri
func (l Local) Delete(ctx context.Context, uri string) error {
	if err := validURI(uri); err != nil {
		return err
	}

	path := l.fullpath(uri)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(path)
}

// Close storage
func (l Local) Close() error {
	return nil
}

func (l Local) fullpath(uri string) string {
	return filepath.Join(string(l), uri)
}

func validURI(uri string) error {
	if uri == "" {
		return ErrEmptyURI
	}

	if strings.HasPrefix(uri, "/") || strings.HasSuffix(uri, "/") {
		return fmt.Errorf("uri invalid: %s", uri)
	}

	return nil
}

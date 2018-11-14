package storage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const defFileMode os.FileMode = os.ModePerm

var (
	ErrEmptyUri = errors.New("empty uri")
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
func (l Local) Put(ctx context.Context, uri string, r io.Reader) error {
	if uri == "" {
		return ErrEmptyUri
	}

	path := l.fullpath(uri)

	if _, err := os.Stat(path); os.IsExist(err) {
		return err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, defFileMode)
}

// Get object by uri
func (l Local) Get(ctx context.Context, uri string) (rc io.ReadCloser, err error) {
	if uri == "" {
		return nil, ErrEmptyUri
	}

	path := l.fullpath(uri)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	return os.Open(path)
}

// Delete object by uri
func (l Local) Delete(ctx context.Context, uri string) error {
	if uri == "" {
		return ErrEmptyUri
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

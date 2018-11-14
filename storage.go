package storage

import (
	"context"
	"io"
)

// Storage for object
type Storage interface {
	// Open(ctx context.Context, path string) (io.ReadCloser, error)

	// Put data and return uri for object
	Put(ctx context.Context, r io.Reader) (uri string, err error)

	// Get object by uri
	Get(ctx context.Context, uri string) (rc io.ReadCloser, err error)

	// Delete object by uri
	Delete(ctx context.Context, uri string) error

	// Close storage
	Close() error
}

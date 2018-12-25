package storage

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocal(t *testing.T) {
	ts := []struct {
		name string
		uri  string
		data []byte
	}{
		{
			name: "normal data",
			uri:  "testfile1",
			data: []byte("hello local storage"),
		},
	}

	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			dir := filepath.Join(os.TempDir(), "local-test")
			l, err := NewLocal(dir)
			assert.NoError(t, err)
			defer func() {
				os.RemoveAll(dir)
			}()

			buf := bytes.NewBuffer(tc.data)
			err = l.Put(context.Background(), tc.uri, buf)
			assert.NoError(t, err)

			rc, err := l.Get(context.Background(), tc.uri)
			assert.NoError(t, err)
			got, err := ioutil.ReadAll(rc)
			assert.NoError(t, err)
			rc.Close()
			assert.Equal(t, tc.data, got)

		})
	}
}

func TestPutDuplicateFile(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "local-test")
	l, err := NewLocal(dir)
	assert.NoError(t, err)
	defer func() {
		os.RemoveAll(dir)
	}()

	for i := 0; i < 2; i++ {
		err = l.Put(context.Background(), "test1.txt", bytes.NewBuffer([]byte("file data")))
		assert.NoError(t, err)
	}
}

func TestDeleteNotExistsFile(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "local-test")
	l, err := NewLocal(dir)
	assert.NoError(t, err)
	defer func() {
		os.RemoveAll(dir)
	}()

	err = l.Delete(context.Background(), "test-not-exists.txt")
	assert.NoError(t, err)
}

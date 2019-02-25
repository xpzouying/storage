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
		{
			name: "multiple directory",
			uri:  "test_dir/testfile2",
			data: []byte("this is test file2"),
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

func TestPutErrorURI(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "local-test")
	l, err := NewLocal(dir)
	assert.NoError(t, err)
	defer func() {
		os.RemoveAll(dir)
	}()

	ts := []struct {
		name string
		uri  string
	}{
		{name: "begin_with_backlash", uri: "/test_dir"},
		{name: "end_with_backlash", uri: "test_dir/"},
		{name: "empty_uri", uri: ""},
	}

	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			err = l.Put(context.Background(), tc.uri, bytes.NewBuffer([]byte("only dir")))
			assert.Error(t, err)
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

func TestDeleteDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "local-test")
	l, err := NewLocal(dir)
	assert.NoError(t, err)
	defer func() {
		os.RemoveAll(dir)
	}()

	// add test file
	dirName := "test-dir"
	dirFullPath := filepath.Join(dir, dirName)
	err = os.Mkdir(dirFullPath, os.ModePerm)
	assert.NoError(t, err)
	testFilePath := filepath.Join(dirFullPath, "test1")
	_, err = os.OpenFile(testFilePath, os.O_CREATE|os.O_RDONLY, os.ModePerm)
	assert.NoError(t, err)

	err = l.Delete(context.Background(), dirName)
	assert.NoError(t, err)

	if _, err = os.Stat(dirFullPath); err != nil {
		if os.IsExist(err) {
			t.Errorf("dir should be delete: %s", dirFullPath)
		}
	}
}

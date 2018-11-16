package storage

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
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
			dir := "/tmp/oss"
			l, err := NewLocal(dir)
			assert.NoError(t, err)

			buf := bytes.NewBuffer(tc.data)
			err = l.Put(context.Background(), tc.uri, buf)
			assert.NoError(t, err)

			rc, err := l.Get(context.Background(), tc.uri)
			assert.NoError(t, err)
			got, err := ioutil.ReadAll(rc)
			assert.NoError(t, err)
			rc.Close()
			assert.Equal(t, tc.data, got)

			cleanup := func() {
				os.RemoveAll(dir)
			}
			cleanup()
		})
	}
}

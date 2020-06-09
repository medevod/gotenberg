package fscleaner

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestRemoveAll(t *testing.T) {
	filenames := []string{
		"foo.txt",
		"foo.jpg",
		"bar.txt",
	}

	logger := log.NewNopLogger()
	fs := afero.NewMemMapFs()
	fscleaner := New(logger, fs)

	err := fs.Mkdir("foo", os.FileMode(0755))
	require.NoError(t, err)

	err = fs.Mkdir("bar", os.FileMode(0755))
	require.NoError(t, err)

	for _, filename := range filenames {
		_, err = fs.Create(fmt.Sprintf("foo/%s", filename))
		require.NoError(t, err)
	}

	fscleaner.RemoveAll("foo")
	time.Sleep(100 * time.Millisecond)

	// Directory foo and its children should not exist anymore.
	_, err = fs.Stat("foo")
	assert.Error(t, err)

	fs = afero.NewReadOnlyFs(fs)
	fscleaner = New(logger, fs)
	fscleaner.RemoveAll("bar")

	// Directory bar should still exist.
	_, err = fs.Stat("bar")
	assert.NoError(t, err)
}

package fs

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

func TestRemoveIfOlderThan(t *testing.T) {
	fs := New(
		log.NewNopLogger(),
		afero.NewMemMapFs(),
	)
	olderThan := 10 * time.Minute

	// It cannot remove a non existing file.
	err := fs.RemoveIfOlderThan("foo.txt", olderThan)
	assert.Error(t, err)

	// It does not removes the file as
	// it is not older as given duration.
	_, err = fs.Create("foo.txt")
	require.NoError(t, err)

	err = fs.RemoveIfOlderThan("foo.txt", olderThan)
	assert.NoError(t, err)

	// It removes the file.
	err = fs.Chtimes(
		"foo.txt",
		time.Now(),
		time.Now().Add(60*time.Minute),
	)
	require.NoError(t, err)

	err = fs.RemoveIfOlderThan("foo.txt", olderThan)
	assert.NoError(t, err)
}

func TestNonBlockingRemoveAll(t *testing.T) {
	fs := New(
		log.NewNopLogger(),
		afero.NewMemMapFs(),
	)

	err := fs.Mkdir("foo", os.FileMode(0700))
	require.NoError(t, err)

	err = fs.Mkdir("bar", os.FileMode(0700))
	require.NoError(t, err)

	for _, filename := range []string{
		"foo.txt",
		"foo.jpg",
		"bar.txt",
	} {
		_, err = fs.Create(fmt.Sprintf("foo/%s", filename))
		require.NoError(t, err)
	}

	// Directory foo (and its children) should not exist anymore.
	fs.NonBlockingRemoveAll("foo")
	time.Sleep(100 * time.Millisecond)

	exists, err := afero.DirExists(fs, "foo")
	require.NoError(t, err)
	assert.False(t, exists)

	// Directory bar should still exist.
	fs = New(
		log.NewNopLogger(),
		afero.NewReadOnlyFs(fs),
	)
	fs.NonBlockingRemoveAll("bar")

	exists, err = afero.DirExists(fs, "bar")
	require.NoError(t, err)
	assert.True(t, exists)
}

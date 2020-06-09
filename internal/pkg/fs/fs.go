package fs

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/afero"
	"time"
)

type Fs struct {
	logger log.Logger

	afero.Fs
}

// New returns a filesystem.
func New(logger log.Logger, fs afero.Fs) *Fs {
	return &Fs{logger, fs}
}

// RemoveIfOlderThan removes a file identified by name
// if it is older than given duration.
func (fs *Fs) RemoveIfOlderThan(
	name string,
	olderThan time.Duration,
) error {
	f, err := fs.Stat(name)
	if err != nil {
		return err
	}

	t := time.Now().Add(olderThan)
	if f.ModTime().Before(t) {
		return nil
	}

	return fs.Remove(name)
}

// NonBlockingRemoveAll works like RemoveAll,
// but in non blocking fashion.
//
// If an error happens, it will be logged.
func (fs *Fs) NonBlockingRemoveAll(path string) {
	go func() {
		if err := fs.RemoveAll(path); err != nil {
			level.Error(fs.logger).Log(
				"msg", "failed to remove all",
				"path", path,
				"err", err.Error(),
			)
		}
	}()
}

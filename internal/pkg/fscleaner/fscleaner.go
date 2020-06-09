package fscleaner

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/afero"
)

type FsCleaner struct {
	logger log.Logger
	fs     afero.Fs
}

func New(logger log.Logger, fs afero.Fs) *FsCleaner {
	return &FsCleaner{logger, fs}
}

func (c *FsCleaner) RemoveAll(path string) {
	go func() {
		if err := c.fs.RemoveAll(path); err != nil {
			level.Error(c.logger).Log(
				"msg", "failed to remove all",
				"path", path,
				"err", err.Error(),
				"trace", "system",
			)
		}
	}()
}

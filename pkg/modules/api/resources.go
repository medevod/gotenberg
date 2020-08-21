package api

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Resources struct {
	dirPath string
	values  map[string][]string
	files   map[string]string
}

func NewResources(ctx *fiber.Ctx) (*Resources, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}

	// TODO tmp from env var.
	res := &Resources{
		dirPath: fmt.Sprintf("/%s/%s", "tmp", uuid.New()),
		values:  form.Value,
		files:   make(map[string]string),
	}

	err = os.MkdirAll(res.dirPath, 0755)
	if err != nil {
		return nil, err
	}

	for _, files := range form.File {
		for _, fh := range files {
			// Avoid directory traversal and normalize filename.
			// See https://github.com/thecodingmachine/gotenberg/issues/104.
			// See https://github.com/thecodingmachine/gotenberg/issues/228.
			t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
			filename, _, err := transform.String(t, strings.ToLower(filepath.Base(fh.Filename)))
			if err != nil {
				return nil, err
			}

			path := fmt.Sprintf("%s/%s", res.dirPath, filename)

			err = ctx.SaveFile(fh, path)
			if err != nil {
				return nil, err
			}

			res.files[filename] = path
		}
	}

	return res, nil
}

func (res *Resources) Close() error {
	// TODO custom errors?
	if res.dirPath == "" {
		return errors.New("resources are already closed")
	}

	_, err := os.Stat(res.dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' does not exist", res.dirPath)
	}

	err = os.RemoveAll(res.dirPath)
	if err != nil {
		return err
	}

	res.dirPath = ""
	res.values = nil
	res.files = nil

	return nil
}

func (res *Resources) WorkingDir() string {
	return res.dirPath
}

func (res *Resources) Paths(extensions ...string) ([]string, error) {
	// TODO check if closed.

	var paths []string

	for filename, path := range res.files {
		for _, ext := range extensions {
			if filepath.Ext(filename) == ext {
				paths = append(paths, path)
			}
		}
	}

	if len(paths) == 0 {
		// TODO custom error?
		return nil, fmt.Errorf("no files found for extensions '%v'", extensions)
	}

	return paths, nil
}

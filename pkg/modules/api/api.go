package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber"
	"github.com/gotenberg/gotenberg/v7/pkg/core"
	flag "github.com/spf13/pflag"
)

func init() {
	core.RegisterModule(API{})
}

type API struct {
	Port           int
	BodyLimit      int
	ReadTimeout    float64
	ProcessTimeout float64
	WriteTimeout   float64
	RootPath       string

	srv *fiber.App
}

func (API) Descriptor() core.ModuleDescriptor {
	return core.ModuleDescriptor{
		ID: "API",
		FlagSet: func() *flag.FlagSet {
			fs := flag.NewFlagSet("api", flag.ExitOnError)
			fs.Int("api-port", 3000, "Set the port on which the API should listen")
			fs.Int("api-body-limit", 16, "Set the maximum allowed size in MB for a request body")
			fs.Float64("api-read-timeout", 10, "Set the maximum duration in seconds allowed to read the full request, including body")
			fs.Float64("api-process-timeout", 10, "Set the maximum duration in seconds allowed to process a request")
			fs.Float64("api-write-timeout", 10, "Set the maximum duration in seconds before timing out writes of the response")
			fs.String("api-root-path", "/", "Set the root path of the API (useful for service discovery via URL paths)")

			return fs
		}(),
		New: func() core.Module { return new(API) },
	}
}

func (a *API) Provision(flags core.ParsedFlags) error {
	a.Port = flags.MustInt("api-port")
	a.BodyLimit = flags.MustInt("api-body-limit")
	a.ReadTimeout = flags.MustFloat("api-read-timeout")
	a.ProcessTimeout = flags.MustFloat("api-process-timeout")
	a.WriteTimeout = flags.MustFloat("api-write-timeout")
	a.RootPath = flags.MustString("api-root-path")

	// As the root path must begin and end with a slash,
	// we correct its value if necessary.
	if !strings.HasPrefix(a.RootPath, "/") {
		a.RootPath = fmt.Sprintf("/%s", a.RootPath)
	}

	if !strings.HasSuffix(a.RootPath, "/") {
		a.RootPath = fmt.Sprintf("%s/", a.RootPath)
	}

	return nil
}

func (a *API) Validate() error {
	var err core.ErrorArray

	// TODO: migrate common validations to a dedicated package.
	if a.Port < 1 || a.Port > 65535 {
		err = append(err,
			errors.New("port must be more than 1 and less than 65535"),
		)
	}

	if a.BodyLimit <= 0 {
		err = append(err,
			errors.New("body limit must be more than 0"),
		)
	}

	if a.ReadTimeout <= 0 {
		err = append(err,
			errors.New("read timeout must be more than 0"),
		)
	}

	if a.ProcessTimeout <= 0 {
		err = append(err,
			errors.New("process timeout must be more than 0"),
		)
	}

	if a.WriteTimeout <= 0 {
		err = append(err,
			errors.New("write timeout must be more than 0"),
		)
	}

	if len(err) > 0 {
		return err
	}

	return nil
}

func (a *API) Start() error {
	// TODO: try prefork setting.
	a.srv = fiber.New(&fiber.Settings{
		DisableStartupMessage: false, // TODO: remove.
		//DisableStartupMessage: true,
		BodyLimit:    a.BodyLimit * 1024 * 1024, // TODO: bytes?
		ReadTimeout:  time.Duration(a.ReadTimeout*1000) * time.Millisecond,
		WriteTimeout: time.Duration(a.WriteTimeout*1000) * time.Millisecond,
	})

	a.srv.Get(a.withRootPath("foo"), func(c *fiber.Ctx) {
		c.Send("Bar")
	})

	return a.srv.Listen(a.Port)
}

func (a *API) withRootPath(path string) string {
	return fmt.Sprintf("%s%s", a.RootPath, path)
}

func (a *API) Stop() error {
	return a.srv.Shutdown()
}

// Interface guards.
var (
	_ core.Module      = (*API)(nil)
	_ core.Provisioner = (*API)(nil)
	_ core.Validator   = (*API)(nil)
	_ core.App         = (*API)(nil)
)

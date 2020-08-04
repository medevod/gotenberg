package ghostscript

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gotenberg/gotenberg/v7/pkg/core"
	"github.com/gotenberg/gotenberg/v7/pkg/modules/api"
)

func init() {
	core.MustRegisterModule(Ghostscript{})
}

type Ghostscript struct {
	binPath string
}

func (Ghostscript) Descriptor() core.ModuleDescriptor {
	return core.ModuleDescriptor{
		ID:      "Ghostscript",
		FlagSet: nil,
		New:     func() core.Module { return new(Ghostscript) },
	}
}

func (gs *Ghostscript) Provision(_ *core.Context) error {
	path, ok := os.LookupEnv("GS_BIN_PATH")
	if !ok {
		return fmt.Errorf("'%s' environment variable is not set", "GS_BIN_PATH")
	}

	gs.binPath = path

	return nil
}

func (gs *Ghostscript) Validate() error {
	info, err := os.Stat(gs.binPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("'%s': %w", gs.binPath, err)
	}

	if info.IsDir() {
		return fmt.Errorf("'%s' is a directory", gs.binPath)
	}

	return nil
}

func (gs *Ghostscript) Routes() []*api.Route {
	return []*api.Route{
		{
			Method: fiber.MethodGet,
			Path:   "api/merge",
			Handlers: []fiber.Handler{
				func(ctx *fiber.Ctx) {
					ctx.Send("Foo")
				},
			},
		},
	}
}

// Interface guards.
var (
	_ core.Module      = (*Ghostscript)(nil)
	_ core.Provisioner = (*Ghostscript)(nil)
	_ core.Validator   = (*Ghostscript)(nil)
	_ api.Router       = (*Ghostscript)(nil)
)

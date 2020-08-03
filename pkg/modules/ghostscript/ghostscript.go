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

func (Ghostscript) ValidateOnRegistration() error {
	// TODO: dedicated package for env var bin path?
	path, ok := os.LookupEnv("GS_BIN_PATH")
	if !ok {
		return fmt.Errorf("'%s' environment variable is not set", "GS_BIN_PATH")
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("'%s': %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("'%s' is a directory", path)
	}

	return nil
}

func (gs *Ghostscript) Provision(_ core.ParsedFlags) error {
	gs.binPath = os.Getenv("GS_BIN_PATH")

	return nil
}

func (gs *Ghostscript) Inject() error {
	mergeRoute := &api.Route{
		Method: fiber.MethodPost,
		Path:   "api/merge",
		Handlers: []fiber.Handler{
			func(ctx *fiber.Ctx) {
				ctx.Send("Foo")
			},
		},
	}

	return api.RegisterRoute(mergeRoute)
}

// Interface guards.
var (
	_ core.Module                = (*Ghostscript)(nil)
	_ core.RegistrationValidator = (*Ghostscript)(nil)
	_ core.Provisioner           = (*Ghostscript)(nil)
	_ core.Dependency            = (*Ghostscript)(nil)
)

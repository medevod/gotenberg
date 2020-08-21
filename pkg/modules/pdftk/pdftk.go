package pdftk

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sort"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/thecodingmachine/gotenberg/v7/pkg/core"
	"github.com/thecodingmachine/gotenberg/v7/pkg/modules/api"
	"github.com/thecodingmachine/gotenberg/v7/pkg/sys"
)

func init() {
	core.MustRegisterModule(PDFtk{})
}

type PDFtk struct {
	binPath string

	logger *zap.Logger
}

func (PDFtk) Descriptor() core.ModuleDescriptor {
	return core.ModuleDescriptor{
		ID:      "PDFtk",
		FlagSet: nil,
		New:     func() core.Module { return new(PDFtk) },
	}
}

func (tk *PDFtk) Provision(ctx *core.Context) error {
	// TODO move to a common package?
	path, ok := os.LookupEnv("PDFTK_BIN_PATH")
	if !ok {
		return fmt.Errorf("'%s' environment variable is not set", "PDFTK_BIN_PATH")
	}

	tk.binPath = path
	tk.logger = ctx.NewLogger(tk)

	return nil
}

func (tk *PDFtk) Validate() error {
	// TODO move to a common package?
	info, err := os.Stat(tk.binPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("'%s': %w", tk.binPath, err)
	}

	if info.IsDir() {
		return fmt.Errorf("'%s' is a directory", tk.binPath)
	}

	return nil
}

func (tk *PDFtk) Routes() []*api.Route {
	return []*api.Route{
		{
			Method: fiber.MethodPost,
			Path:   "api/merge",
			Handlers: []fiber.Handler{
				func(ctx *fiber.Ctx) {
					res, err := api.NewResources(ctx)
					if err != nil {
						ctx.Next(err)
					}

					pdfs, err := res.Paths(".pdf")
					if err != nil {
						ctx.Next(err)
					}

					destPath := fmt.Sprintf("%s/%s.pdf", res.WorkingDir(), uuid.New())

					err = tk.Merge(ctx.Context(), destPath, pdfs)
					if err != nil {
						ctx.Next(err)
					}

					// TODO close resources.
					// TODO remove sent file.
					// TODO compress option?
					err = ctx.SendFile(destPath)
					if err != nil {
						ctx.Next(err)
					}
				},
			},
		},
	}
}

func (tk *PDFtk) Merge(ctx context.Context, destPath string, paths []string) error {
	if len(paths) == 1 {
		err := os.Rename(paths[0], destPath)
		if err != nil {
			return err
		}

		return nil
	}

	// See https://github.com/thecodingmachine/gotenberg/issues/139.
	sort.Strings(paths)

	var args []string
	args = append(args, paths...)
	args = append(args, "cat", "output", destPath)

	err := sys.Exec(ctx, tk.binPath, args...)
	if err != nil {
		return err
	}

	return nil
}

// Interface guards.
var (
	_ core.Module      = (*PDFtk)(nil)
	_ core.Provisioner = (*PDFtk)(nil)
	_ core.Validator   = (*PDFtk)(nil)
	_ api.Router       = (*PDFtk)(nil)
)

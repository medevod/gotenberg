package logging

import (
	"errors"
	"strings"

	"github.com/gotenberg/gotenberg/v7/pkg/core"
	flag "github.com/spf13/pflag"
)

func init() {
	core.MustRegisterModule(Logging{})
}

type Logging struct {
	level  string
	format string
}

func (Logging) Descriptor() core.ModuleDescriptor {
	return core.ModuleDescriptor{
		ID: "Logging",
		FlagSet: func() *flag.FlagSet {
			fs := flag.NewFlagSet("logging", flag.ExitOnError)
			fs.String("log-level", "error", "Set the log level - either error, warn, info or debug")
			fs.String("log-format", "auto", "Set log format - either auto, text or json")

			return fs
		}(),
		New: func() core.Module { return new(Logging) },
	}
}

func (log *Logging) Provision(flags core.ParsedFlags) error {
	log.level = strings.ToLower(flags.MustString("log-level"))
	log.format = strings.ToLower(flags.MustString("log-format"))

	return nil
}

func (log *Logging) Validate() error {
	var errs core.ErrorArray

	switch log.level {
	case "error", "warn", "info", "debug":
		break
	default:
		errs = append(errs, errors.New("log level must be either error, warn, info or debug"))
	}

	switch log.format {
	case "auto", "text", "json":
		break
	default:
		errs = append(errs, errors.New("log format must be either auto, text or json"))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (log *Logging) ModifyCore() error {
	return core.ModifiyLogger(log.level, log.format)
}

// Interface guards.
var (
	_ core.Module       = (*Logging)(nil)
	_ core.Provisioner  = (*Logging)(nil)
	_ core.Validator    = (*Logging)(nil)
	_ core.CoreModifier = (*Logging)(nil)
)

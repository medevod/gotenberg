package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gotenberg/gotenberg/v7/pkg/core"
	flag "github.com/spf13/pflag"
)

func Run() {
	// Creates the root FlagSet.
	fs := flag.NewFlagSet("gotenberg", flag.ExitOnError)

	// Adds the modules flags to the root FlagSet.
	mods := core.GetModules()
	for _, mod := range mods {
		fs.AddFlagSet(mod.FlagSet)
	}

	// Parses the flags...
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// ...and creates a wrapper around those.
	parsedFlags := core.ParsedFlags{FlagSet: fs}

	// Initializes module's instances.
	var (
		instances []core.Module
		apps      []core.App
	)
	for _, mod := range mods {
		instance := mod.New()
		instances = append(instances, instance)

		if p, ok := instance.(core.Provisioner); ok {
			err := p.Provision(parsedFlags)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// TODO: map all validation errors?
		if v, ok := instance.(core.Validator); ok {
			err := v.Validate()
			if err != nil {
				fmt.Printf("%s module validation failed: %s\n", mod.ID, err)
				os.Exit(1)
			}
		}

		if a, ok := instance.(core.App); ok {
			apps = append(apps, a)
		}
	}

	// TODO: improve startup (and shutdown) of apps.

	for _, app := range apps {
		go func(a core.App) {
			err := a.Start()
			if err != nil {
				fmt.Println(err)
			}
		}(app)
	}

	quit := make(chan os.Signal, 1)

	// we'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(quit, os.Interrupt)

	// block until we receive our signal.
	<-quit

	for _, app := range apps {
		go func(a core.App) {
			err := a.Stop()
			if err != nil {
				fmt.Println(err)
			}
		}(app)
	}

	os.Exit(0)
}

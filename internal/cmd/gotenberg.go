package cmd

import (
	"context"
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
	descriptors := core.GetModuleDescriptors()
	for _, desc := range descriptors {
		fs.AddFlagSet(desc.FlagSet)
	}

	// Parses the flags...
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		// TODO: use dedicated exit codes.
		os.Exit(1)
	}

	// ...and creates a wrapper around those.
	parsedFlags := &core.ParsedFlags{FlagSet: fs}

	ctx := core.NewContext(context.Background(), parsedFlags)

	if err := ctx.CoreModifiers(); err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	apps, err := ctx.Apps()
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	// TODO: improve startup (and shutdown) of apps.

	for _, app := range apps {
		go func(a core.App) {
			err := a.Start()
			if err != nil {
				// TODO print module ID.
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
				// TODO print module ID.
				fmt.Println(err)
			}
		}(app)
	}

	os.Exit(0)
}

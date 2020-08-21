package gotenberg

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	flag "github.com/spf13/pflag"
	"github.com/thecodingmachine/gotenberg/v7/pkg/core"
)

var version = "snapshot"

func Run() {
	fmt.Printf("Gotenberg version %s\n", version)

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
	parsedFlags := core.ParsedFlags{FlagSet: fs}

	ctx := core.NewContext(context.Background(), descriptors, parsedFlags)

	modifiers, err := ctx.Modules(new(core.CoreModifier))
	if err != nil {
		fmt.Printf("[FATAL] %s\n", err)
		os.Exit(1)
	}

	for _, modifier := range modifiers {
		if err := modifier.(core.CoreModifier).ModifyCore(); err != nil {
			fmt.Printf("[FATAL] %s\n", err)
			os.Exit(1)
		}
	}

	apps, err := ctx.Modules(new(core.App))
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
		}(app.(core.App))
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
		}(app.(core.App))
	}

	os.Exit(0)
}

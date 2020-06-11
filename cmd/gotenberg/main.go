package main

import (
	"fmt"
	"os"

	"github.com/thecodingmachine/gotenberg/internal/app/api"
	"github.com/thecodingmachine/gotenberg/internal/pkg/cli"
	"github.com/thecodingmachine/gotenberg/internal/pkg/codes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("[ERROR] expected a command")
		os.Exit(codes.ExitCodeFailedStartup)
	}

	var cmd *cli.Command

	switch os.Args[1] {
	case "api":
		cmd = api.NewCommand()
	default:
		fmt.Printf("[ERROR] '%s' is not a recognized command\n", os.Args[1])
		os.Exit(codes.ExitCodeFailedStartup)
	}

	fs := cmd.Flags
	if fs == nil {
		fmt.Printf("[FATAL] nil flags\n")
		os.Exit(codes.ExitCodeFailedStartup)
	}

	err := fs.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		os.Exit(codes.ExitCodeFailedStartup)
	}

	exitCode, err := cmd.Func(cli.Flags{FlagSet: fs})
	if err != nil {
		fmt.Println(err)
		os.Exit(exitCode)
	}

	os.Exit(exitCode)
}

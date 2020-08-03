package main

import (
	"github.com/gotenberg/gotenberg/v7/internal/cmd"

	_ "github.com/gotenberg/gotenberg/v7/pkg/modules/api"
	_ "github.com/gotenberg/gotenberg/v7/pkg/modules/ghostscript"
	_ "github.com/gotenberg/gotenberg/v7/pkg/modules/logging"
)

func main() {
	cmd.Run()
}

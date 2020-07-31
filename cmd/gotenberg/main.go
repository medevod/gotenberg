package main

import (
	"github.com/gotenberg/gotenberg/v7/internal/cmd"

	_ "github.com/gotenberg/gotenberg/v7/pkg/modules/api"
)

func main() {
	cmd.Run()
}

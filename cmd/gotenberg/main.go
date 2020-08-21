package main

import (
	"github.com/thecodingmachine/gotenberg/v7/internal/app/gotenberg"

	_ "github.com/thecodingmachine/gotenberg/v7/pkg/modules/api"
	_ "github.com/thecodingmachine/gotenberg/v7/pkg/modules/logging"
	_ "github.com/thecodingmachine/gotenberg/v7/pkg/modules/pdftk"
)

func main() {
	gotenberg.Run()
}

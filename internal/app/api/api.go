package api

import (
	"fmt"

	flag "github.com/spf13/pflag"
	"github.com/thecodingmachine/gotenberg/internal/pkg/cli"
	"github.com/thecodingmachine/gotenberg/internal/pkg/codes"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "api",
		Func: cmdAPI,
		Flags: func() *flag.FlagSet {
			fs := flag.NewFlagSet("api", flag.ExitOnError)
			fs.Bool("enable-html", false, "enable HTML conversion")
			fs.Bool("enable-url", false, "enable URL conversion")
			fs.Bool("enable-markdown", false, "enable Markdown conversion")
			fs.Bool("enable-office", false, "enable Office conversion")
			fs.Bool("enable-merge", false, "enable merge")

			return fs
		}(),
	}
}

func cmdAPI(fs cli.Flags) (int, error) {
	enableHTMLFlag := fs.Bool("enable-html")
	enableURLFlag := fs.Bool("enable-url")
	enableMarkdownFlag := fs.Bool("enable-markdown")
	enableOfficeFlag := fs.Bool("enable-office")
	enableMergeFlag := fs.Bool("enable-merge")

	fmt.Println("html ", enableHTMLFlag)
	fmt.Println("url ", enableURLFlag)
	fmt.Println("md ", enableMarkdownFlag)
	fmt.Println("office ", enableOfficeFlag)
	fmt.Println("merge ", enableMergeFlag)

	return codes.ExitCodeSuccess, nil
}

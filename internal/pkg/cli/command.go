package cli

import flag "github.com/spf13/pflag"

type Command struct {
	Name string

	Func CommandFunc

	Flags *flag.FlagSet
}

type CommandFunc func(Flags) (int, error)

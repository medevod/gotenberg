package cli

import (
	"strconv"

	flag "github.com/spf13/pflag"
)

type Flags struct {
	*flag.FlagSet
}

func (f Flags) String(name string) string {
	return f.FlagSet.Lookup(name).Value.String()
}

func (f Flags) Bool(name string) bool {
	val, _ := strconv.ParseBool(f.String(name))

	return val
}

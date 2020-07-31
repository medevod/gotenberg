package core

import flag "github.com/spf13/pflag"

// ParsedFlags wraps a flag.FlagSet so that retrieving
// the typed values is easier.
type ParsedFlags struct {
	*flag.FlagSet
}

// MustString returns the string value of a
// flag given by name. It panics if an error occurs.
func (f ParsedFlags) MustString(name string) string {
	val, err := f.GetString(name)
	if err != nil {
		panic(err)
	}

	return val
}

// MustBool returns the boolean value of a
// flag given by name. It panics if an error occurs.
func (f ParsedFlags) MustBool(name string) bool {
	val, err := f.GetBool(name)
	if err != nil {
		panic(err)
	}

	return val
}

// MustInt returns the int value of a
// flag given by name. It panics if an error occurs.
func (f ParsedFlags) MustInt(name string) int {
	val, err := f.GetInt(name)
	if err != nil {
		panic(err)
	}

	return val
}

// MustFloat returns the float value of a
// flag given by name. It panics if an error occurs.
func (f ParsedFlags) MustFloat(name string) float64 {
	val, err := f.GetFloat64(name)
	if err != nil {
		panic(err)
	}

	return val
}

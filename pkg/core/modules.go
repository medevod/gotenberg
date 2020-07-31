package core

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"sort"
	"sync"
)

type Module interface {
	Descriptor() ModuleDescriptor
}

type ModuleID string

type ModuleDescriptor struct {
	// ID is the unique name
	// of the module.
	// Required.
	ID ModuleID

	// FlagSet is the definition of the flags
	// of the module.
	FlagSet *flag.FlagSet

	// New returns a new and empty instance of
	// the module's type.
	// Required.
	New func() Module
}

type Provisioner interface {
	Provision(ParsedFlags) error
}

type Validator interface {
	Validate() error
}

type App interface {
	Start() error
	Stop() error
}

func RegisterModule(m Module) {
	desc := m.Descriptor()

	if desc.ID == "" {
		panic("module with an empty ID cannot be registered")
	}

	if desc.New == nil {
		panic("module New function cannot be nil")
	}

	if val := desc.New(); val == nil {
		panic("module New function cannot return a nil instance")
	}

	modulesMu.Lock()
	defer modulesMu.Unlock()

	if _, ok := modules[string(desc.ID)]; ok {
		panic(fmt.Sprintf("module '%s' is already registered", desc.ID))
	}

	modules[string(desc.ID)] = desc
}

func GetModules() []ModuleDescriptor {
	modulesMu.RLock()
	defer modulesMu.RUnlock()

	var mods []ModuleDescriptor
	for _, m := range modules {
		mods = append(mods, m)
	}

	sort.Slice(mods, func(i, j int) bool {
		return mods[i].ID < mods[j].ID
	})

	return mods
}

var (
	modules   = make(map[string]ModuleDescriptor)
	modulesMu sync.RWMutex
)

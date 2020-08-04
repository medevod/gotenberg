package core

import (
	"fmt"
	"sort"
	"sync"

	flag "github.com/spf13/pflag"
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
	// Optional.
	FlagSet *flag.FlagSet

	// New returns a new and empty instance of
	// the module's type.
	// Required.
	New func() Module
}

type Provisioner interface {
	Provision(*Context) error
}

type Validator interface {
	Validate() error
}

type CoreModifier interface {
	ModifyCore() error
}

type App interface {
	Start() error
	Stop() error
}

func MustRegisterModule(mod Module) {
	desc := mod.Descriptor()

	if desc.ID == "" {
		panic("module with an empty ID cannot be registered")
	}

	if desc.New == nil {
		panic("module New function cannot be nil")
	}

	if val := desc.New(); val == nil {
		panic("module New function cannot return a nil instance")
	}

	descriptorsMu.Lock()
	defer descriptorsMu.Unlock()

	if _, ok := descriptors[string(desc.ID)]; ok {
		panic(fmt.Sprintf("module '%s' is already registered", desc.ID))
	}

	descriptors[string(desc.ID)] = desc
}

func GetModuleDescriptors() []ModuleDescriptor {
	descriptorsMu.RLock()
	defer descriptorsMu.RUnlock()

	var mods []ModuleDescriptor
	for _, desc := range descriptors {
		mods = append(mods, desc)
	}

	sort.Slice(mods, func(i, j int) bool {
		return mods[i].ID < mods[j].ID
	})

	return mods
}

var (
	descriptors   = make(map[string]ModuleDescriptor)
	descriptorsMu sync.RWMutex
)

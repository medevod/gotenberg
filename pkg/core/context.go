package core

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"reflect"
)

type Context struct {
	context.Context
	descriptors     []ModuleDescriptor
	flags           ParsedFlags
	moduleInstances map[ModuleID]interface{}
}

func NewContext(
	ctx context.Context,
	descriptors []ModuleDescriptor,
	flags ParsedFlags,
) *Context {
	// TODO cancel?
	newCtx := &Context{
		Context:         ctx,
		descriptors:     descriptors,
		flags:           flags,
		moduleInstances: make(map[ModuleID]interface{}),
	}

	return newCtx
}

func (ctx *Context) ParsedFlags() ParsedFlags {
	return ctx.flags
}

func (ctx *Context) Modules(kind interface{}) ([]interface{}, error) {
	realKind := reflect.TypeOf(kind).Elem()

	var mods []interface{}
	for _, desc := range ctx.descriptors {
		newInstance := desc.New()

		if ok := reflect.TypeOf(newInstance).Implements(realKind); ok {
			// The module implements the requested interface.
			// We check if it has already been initialized.
			instance, ok := ctx.moduleInstances[desc.ID]

			if ok {
				mods = append(mods, instance)
			} else {
				err := ctx.loadModule(desc.ID, newInstance)
				if err != nil {
					return nil, err
				}

				mods = append(mods, newInstance)
			}
		}
	}

	return mods, nil
}

func (ctx *Context) loadModule(ID ModuleID, instance interface{}) error {
	if prov, ok := instance.(Provisioner); ok {
		// The instance can be provisioned.
		err := prov.Provision(ctx)
		if err != nil {
			return fmt.Errorf("cannot provision module '%s': %w", ID, err)
		}
	}

	if validator, ok := instance.(Validator); ok {
		// The instance can be validated.
		err := validator.Validate()
		if err != nil {
			return fmt.Errorf("validation failed for module '%s': %w", ID, err)
		}
	}

	ctx.moduleInstances[ID] = instance

	return nil
}

func (ctx *Context) NewLogger(mod Module) *zap.Logger {
	modID := string(mod.Descriptor().ID)

	return Log().With(
		zap.String("module", modID),
	)
}

package core

import (
	"context"
	"fmt"
	"reflect"
)

type Context struct {
	context.Context
	flags           *ParsedFlags
	moduleInstances map[string]interface{}
}

func NewContext(ctx context.Context, flags *ParsedFlags) *Context {
	// TODO cancel?
	newCtx := &Context{
		Context:         ctx,
		flags:           flags,
		moduleInstances: make(map[string]interface{}),
	}

	return newCtx
}

func (ctx *Context) ParsedFlags() *ParsedFlags {
	return ctx.flags
}

func (ctx *Context) Apps() ([]App, error) {
	descriptorsMu.RLock()
	defer descriptorsMu.RUnlock()

	var apps []App
	for _, desc := range descriptors {
		mod := desc.New()

		if app, ok := mod.(App); ok {
			err := ctx.loadModule(string(desc.ID), app)
			if err != nil {
				return nil, err
			}

			apps = append(apps, app)
		}
	}

	return apps, nil
}

func (ctx *Context) CoreModifiers() error {
	descriptorsMu.RLock()
	defer descriptorsMu.RUnlock()

	for _, desc := range descriptors {
		mod := desc.New()

		if modifier, ok := mod.(CoreModifier); ok {
			if err := ctx.loadModule(string(desc.ID), modifier); err != nil {
				return err
			}

			if err := modifier.ModifyCore(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ctx *Context) Modules(kind interface{}) ([]interface{}, error) {
	realKind := reflect.TypeOf(kind).Elem()

	descriptorsMu.RLock()
	defer descriptorsMu.RUnlock()

	var mods []interface{}
	for _, desc := range descriptors {
		newInstance := desc.New()

		if ok := reflect.TypeOf(newInstance).Implements(realKind); ok {
			// The module implements the requested interface.
			// We check if it has already been initialized.
			instance, ok := ctx.moduleInstances[string(desc.ID)]

			if ok {
				mods = append(mods, instance)
			} else {
				err := ctx.loadModule(string(desc.ID), newInstance)
				if err != nil {
					return nil, err
				}

				mods = append(mods, newInstance)
			}
		}
	}

	return mods, nil
}

func (ctx *Context) loadModule(ID string, instance interface{}) error {
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

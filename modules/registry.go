package modules

import "barista.run/bar"

// Registry registers bar modules. It can be used to easily register modules
// until an error is encountered and pass them to `barista.Run`. Modules from
// registry appear in the bar in the same order as they were added.
//
//   registry := modules.NewRegistry()
//
//   registry.
//       Add(someModule, someOtherModule).
//       Addf(func() (bar.Module, error) {
//           return someModuleFromFactory, nil
//       })
//
//   if err := registry.Err(); err != nil {
//     panic(err)
//   }
//
//   panic(barista.Run(registry.Modules()...))
type Registry struct {
	modules []bar.Module
	err     error
}

// NewRegistry creates a new *Registry for bar modules.
func NewRegistry() *Registry {
	return &Registry{
		modules: make([]bar.Module, 0),
	}
}

// Add adds modules to the registry. Modules that are nil are ignored. If a
// factory func passed to `Addf` previously returned an error, adding modules
// here is a no-op.
func (r *Registry) Add(modules ...bar.Module) *Registry {
	if r.err != nil {
		return r
	}

	for _, module := range modules {
		if module != nil {
			r.modules = append(r.modules, module)
		}
	}
	return r
}

// Addf adds a module using a factory func. If the factory returns a nil
// module, it is ignored. Errors returned by the factory will cause the
// registry to not accept any more modules via `Add` or `Addf`.
func (r *Registry) Addf(factory func() (bar.Module, error)) *Registry {
	if r.err != nil {
		return r
	}

	var module bar.Module

	module, r.err = factory()

	return r.Add(module)
}

// Err returns the first error returned by a factory func or nil if there was
// none.
func (r *Registry) Err() error {
	return r.err
}

// Modules returns the bar modules. This can be used to pass modules directly
// to `barista.Run`.
func (r *Registry) Modules() []bar.Module {
	return r.modules
}

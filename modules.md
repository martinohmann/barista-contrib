---
title: Module Helpers
---

This package contains helpers which provide functionality around module management.

## Module registry

The module registry can be used to assemble the modules for a bar using a
fluent interface. Modules can be added directly using `Add`, or via `Addf`
using a factory func if the creation of the module is more complex. `nil`
modules are ignored and not added to the bar. The first error returned by any
of the factories passed to `Addf` will cause subsequent calls to `Add` or
`Addf` to be no-op.

### Usage example

```go
registry := modules.NewRegistry()

err := registry.
    Addf(func() (bar.Module, error) {
        ifaces, err := net.Interfaces()
        if err != nil {
            return nil, err
        }

        mods := make([]bar.Module, len(ifaces))

        for i, iface := range ifaces {
            mods[i] = netspeed.New(iface.Name).Output(func(s netspeed.Speeds) bar.Output {
                return outputs.Textf("%s, tx:%s, rx:%s",
                    iface.Name, format.IByterate(s.Tx), format.IByterate(s.Rx))
            })
        }

        group, _ := switching.Group(mods...)

        return group, nil
    }).
    Add(
        ipify.New(),
        sysinfo.New(),
        diskspace.New("/"),
        xset.New().Output(func(info dpms.Info) bar.Output {
            out := outputs.Text("dpms")

            if info.Enabled {
                return out.Color(colors.Scheme("enabled"))
            }

            return out.Color(colors.Scheme("disabled"))
        }),
    ).
    Err()

if err != nil {
    log.Fatal(err)
}

barista.Run(registry.Modules()...)
```

---
title: Updates
---

Show available updates: `updates.New(someProvider)`.

Updates supports displaying currently available updates using pluggable
providers, with the ability to add custom providers fairly easily. Provider is
just

```go
type Provider interface {
	Updates() (Info, error)
}
```

The following update providers are available in barista-contrib:

* [pacman](https://godoc.org/github.com/martinohmann/barista-contrib/modules/updates/pacman):
  Checks for pacman updates using `checkupdates` from
* [pacman-contrib](https://www.archlinux.org/packages/community/x86_64/pacman-contrib/).

* [yay](https://godoc.org/github.com/martinohmann/barista-contrib/modules/updates/yay):
  Checks for Arch Linux updates using [`yay`](https://github.com/Jguer/yay).

## Configuration

* `Output(func(updates.Info) bar.Output)`: Sets the output format.

* `Every(time.Duration)`: Sets the interval to wait before checking for updates again. Defaults to 1 hour.

* `Refresh()`: Triggers an update check and refresh the module.

## Examples

<div class="module-example-out">5 updates</div>
Show updates if available and display package details using
[beeep](https://github.com/gen2brain/beeep) on left-click:

```go
updates.New(yay.New()).Output(func(info updates.Info) bar.Output {
  if info.Updates == 0 {
    return nil
  }

  return outputs.Textf("%d updates", info.Updates).
    OnClick(click.Left(func() {
      beeep.Notify("Available Pacman Updates", info.PackageDetails.String(), "")
    }))
})
```

## Data: `type Info struct`

### Fields

* `Updates int`: Number of available updates.
* `PackageDetails PackageDetails`: Contains details about package updates if the provider supports it.

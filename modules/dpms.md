---
title: DPMS
---

Show the Display Power Management Signaling (DPMS) state: `ip.New(someProvider)`.

DPMS supports displaying the DPMS using pluggable providers, with the ability to add custom providers fairly easily. Provider is just

```go
type Provider interface {
	// Get retrieves the current DPMS status, returning true if it is enabled.
	Get() (bool, error)

	// Set enables or disables DPMS.
	Set(enabled bool) error
}
```

The following DPMS providers are available in barista-contrib:

* [xset](https://godoc.org/github.com/martinohmann/barista-contrib/modules/dpms/xset): Query and update the DPMS state using `xset`.

## Configuration

* `Output(func(dpms.Info) bar.Output)`: Sets the output format.

  If a segment does not have a click handler, the module will set a default click handler, which:
  - Toggles DPMS state on left click

* `Every(time.Duration)`: Sets the interval to wait before checking the DPMS state again. Defaults to 1 minute.

* `Refresh()`: Checks the DPMS state and refreshes the module.

## Examples

<div class="module-example-out">dpms: on</div>
Show DPMS state an react on click events:

```go
func ifLeft(dofn func()) func(bar.Event) {
	return func(e bar.Event) {
		if e.Button == bar.ButtonLeft {
			dofn()
		}
	}
}

xset.New().Output(func(info dpms.Info) bar.Output {
    if info.Enabled {
        return outputs.Text("dpms: on").OnClick(ifLeft(info.Disable))
    }

    return outputs.Text("dpms: off").OnClick(ifLeft(info.Enable))
})
```

Using `.OnClick(nil)` prevents the default click handler of the dpms module
from being added to part of the output.

## Data: `type Info struct`

### Fields

* `Enabled bool`: `true` if DPMS is enabled.

### Methods

* `String() string`: Returns `dpms enabled` if DPMS is enabled, `dpms disabled` otherwise.

### Controller Methods

In addition to the data methods listed above, dpms' `Info` type also provides controller
methods to interact with the DPMS state:

* `Enable()`: Enables DPMS.
* `Disable()`: Disables DPMS.
* `Toggle()`: Enables DPMS if it's off, disables it otherwise.

---
title: Keyboard
---

Show the currently active keyboard layout: `keyboard.New(someProvider, layouts...)`.

Keyboard supports displaying the keyboard layout using pluggable providers,
with the ability to add custom providers fairly easily. Provider is just

```go
type Provider interface {
	// GetLayout retrieves the name of the currently active keyboard layout.
	GetLayout() (string, error)

	// SetLayout sets a new keyboard layout.
	SetLayout(layout string) error
}
```

The following keyboard providers are available in barista-contrib:

* [xkbmap](https://godoc.org/github.com/martinohmann/barista-contrib/modules/keyboard/xkbmap): Retrieve and set the keyboard layout `xkbmap`.

## Configuration

* `Output(func(keyboard.Layout) bar.Output)`: Sets the output format.

  If a segment does not have a click handler, the module will set a default click handler, which:
  - Switched to next keyboard layout on left click or scrollup.
  - Switched to previous keyboard layout on right click or scrolldown.

* `Every(time.Duration)`: Sets the interval to wait before checking the keyboard layout again. Defaults to 10 seconds.

* `Refresh()`: Retrieve the keyboard layout and refreshes the module.

## Examples

<div class="module-example-out">us</div>
Show the current keyboard layout and add custom click handlers:

```go
xkbmap.New("us", "de").Output(func(layout keyboard.Layout) bar.Output {
    return outputs.Textf("%s", strings.ToUpper(layout.Name)).OnClick(func(e bar.Event) {
		switch e.Button {
		case bar.ButtonLeft:
			l.Next()
		case bar.ButtonRight:
			l.Previous()
		}
	})
})
```

Using `.OnClick(nil)` prevents the default click handler of the keyboard module
from being added to part of the output.

## Data: `type Layout struct`

### Fields

* `Name string`: Name of the currently active keyboard layout.

### Methods

* `String() string`: Returns the name of the current keyboard layout.

### Controller Methods

In addition to the data methods listed above, keyboard's `Layout` type also provides controller
methods to control the active keyboard layout:

* `Next()`: Switches to the next keyboard layout from the list.
* `Previous()`: Switches to the previous keyboard from the list.
* `SetLayout(string)`: Set the keyboard layout. The passed layout name has to
  be one of the layouts passed to `keyboard.New()` when initializing the
  module, otherwise this is a no-op.
* `GetLayouts() []string`: Gets the names of the keyboard layouts the
  controller knows about. These are the names that were passed to
  `keyboard.New()` on module initialization.

package xkbmap

import (
	"github.com/martinohmann/barista-contrib/internal/xkbmap"
	"github.com/martinohmann/barista-contrib/modules/keyboard"
)

// New creates a new *keyboard.Module using xkbmap as provider for keyboard
// layouts.
func New(layouts ...string) *keyboard.Module {
	return keyboard.New(&provider{}, layouts...)
}

type provider struct{}

// SetLayout implements keyboard.Provider.
func (p *provider) SetLayout(layout string) error {
	return xkbmap.SetLayout(layout)
}

// GetLayout implements keyboard.Provider.
func (p *provider) GetLayout() (string, error) {
	info, err := xkbmap.Query()
	if err != nil {
		return "", err
	}

	return info.Layout, nil
}

// Package yay contains an updates.Provider that uses `yay` to check for Arch
// Linux package updates.
package yay

import (
	"github.com/martinohmann/barista-contrib/internal/exec"
	"github.com/martinohmann/barista-contrib/modules/updates"
	"github.com/martinohmann/barista-contrib/modules/updates/pacman"
)

// Option is a func that can be passed to New to configure the yay update
// provider.
type Option func(p *Provider)

// AUROnly option makes yay only check for updates for AUR packages.
func AUROnly(p *Provider) {
	p.aurOnly = true
}

// New creates a new *Provider and configures it with the provided options.
func New(options ...Option) *Provider {
	p := &Provider{}

	for _, option := range options {
		option(p)
	}

	return p
}

type Provider struct {
	aurOnly bool
}

// Updates implements updates.Provider.
func (p *Provider) Updates() (updates.Info, error) {
	args := []string{"-Qu"}
	if p.aurOnly {
		args = []string{"-Qua"}
	}

	out, err := exec.CommandOutput("yay", args...)
	if err != nil {
		return updates.Info{}, err
	}

	details, err := pacman.ParsePackageDetails(out)
	if err != nil {
		return updates.Info{}, err
	}

	info := updates.Info{
		Updates:        len(details),
		PackageDetails: details,
	}

	return info, nil
}

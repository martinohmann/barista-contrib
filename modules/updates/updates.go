package updates

import (
	"fmt"
	"strings"
	"time"

	"barista.run/bar"
	"barista.run/base/notifier"
	"barista.run/base/value"
	"barista.run/outputs"
	"barista.run/timing"
)

// Provider provides the count of currently available updates for the bar.
type Provider interface {
	Updates() (Info, error)
}

// ProviderFunc is a func that satisfies the Provider interface.
type ProviderFunc func() (Info, error)

// Updates implements Provider.
func (f ProviderFunc) Updates() (Info, error) {
	return f()
}

// Info contains information about available updates.
type Info struct {
	// Updates is the number of available updates.
	Updates int
	// PackageDetails are optional details for the packages that updates are
	// available for.
	PackageDetails PackageDetails
}

// PackageDetails contains details about package updates.
type PackageDetails []PackageDetail

// String implements fmt.Stringer.
func (d PackageDetails) String() string {
	var sb strings.Builder

	for i, detail := range d {
		sb.WriteString(detail.String())
		if i < len(d)-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

// PackageDetail contains information about a single package update.
type PackageDetail struct {
	// PackageName is the name of the package.
	PackageName string
	// CurrentVersion is the currently installed package version.
	CurrentVersion string
	// TargetVersion is the version of the available package update.
	TargetVersion string
}

// String implements fmt.Stringer.
func (d PackageDetail) String() string {
	return fmt.Sprintf("%s %s -> %s", d.PackageName, d.CurrentVersion, d.TargetVersion)
}

// Module is a module for displaying currently available updates in the bar.
type Module struct {
	outputFunc value.Value // of func(Info) bar.Output
	provider   Provider
	notifyCh   <-chan struct{}
	notifyFn   func()
	scheduler  *timing.Scheduler
}

// New creates a new *Module with the given update count provider. By default,
// the module will refresh the update counts every hour. The refresh interval
// can be configured using `Every`.
func New(provider Provider) *Module {
	m := &Module{
		provider:  provider,
		scheduler: timing.NewScheduler(),
	}

	m.notifyFn, m.notifyCh = notifier.New()
	m.outputFunc.Set(func(info Info) bar.Output {
		if info.Updates == 1 {
			return outputs.Text("1 update")
		}
		return outputs.Textf("%d updates", info.Updates)
	})

	m.Every(time.Hour)

	return m
}

// Stream implements bar.Module.
func (m *Module) Stream(s bar.Sink) {
	info, err := m.provider.Updates()
	outputFunc := m.outputFunc.Get().(func(Info) bar.Output)
	for {
		if !s.Error(err) {
			s.Output(outputFunc(info))
		}

		select {
		case <-m.outputFunc.Next():
			outputFunc = m.outputFunc.Get().(func(Info) bar.Output)
		case <-m.notifyCh:
			info, err = m.provider.Updates()
		case <-m.scheduler.C:
			info, err = m.provider.Updates()
		}
	}
}

// Output updates the output format func.
func (m *Module) Output(format func(Info) bar.Output) *Module {
	m.outputFunc.Set(format)
	return m
}

// Every configures the refresh interval for the module. Passing a zero
// interval will disable refreshing.
func (m *Module) Every(interval time.Duration) *Module {
	if interval == 0 {
		m.scheduler.Stop()
	} else {
		m.scheduler.Every(interval)
	}
	return m
}

// Refresh forces a refresh of the module output.
func (m *Module) Refresh() {
	m.notifyFn()
}

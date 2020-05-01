package dpms

import (
	"time"

	"barista.run/bar"
	"barista.run/base/notifier"
	"barista.run/base/value"
	l "barista.run/logging"
	"barista.run/outputs"
	"barista.run/timing"
)

// Provider provides means the get and set the DPMS status.
type Provider interface {
	// Get retrieves the current DPMS status, returning true if it is enabled.
	Get() (bool, error)

	// Set enables or disables DPMS.
	Set(enabled bool) error
}

// Info contains the current DPMS status. It also exposes controller methods to
// change the DPMS status.
type Info struct {
	Enabled bool

	provider Provider
	update   func()
}

// String implements fmt.Stringer.
func (i Info) String() string {
	if i.Enabled {
		return "dpms enabled"
	}

	return "dpms disabled"
}

// Enable enables DPMS.
func (i Info) Enable() {
	i.setEnabled(true)
}

// Disable disables DPMS.
func (i Info) Disable() {
	i.setEnabled(false)
}

// Toggle enables DPMS if it is disabled and vice versa.
func (i Info) Toggle() {
	enabled, err := i.provider.Get()
	if err != nil {
		l.Log("Error obtaining DPMS status: %v", err)
		return
	}

	i.setEnabled(!enabled)
}

func (i Info) setEnabled(enabled bool) {
	if err := i.provider.Set(enabled); err != nil {
		l.Log("Error updating DPMS status: %v", err)
		return
	}

	i.update()
}

// Module is a module for displaying and interacting with the current DPMS
// status.
type Module struct {
	provider   Provider
	outputFunc value.Value // of func(Info) bar.Output
	notifyCh   <-chan struct{}
	notifyFn   func()
	scheduler  *timing.Scheduler
}

// New creates a new *Module which uses given provider to query and update the
// DPMS status. By default, the module will refresh every minute. The refresh
// interval can be configured using `Every`.
func New(provider Provider) *Module {
	m := &Module{
		provider:  provider,
		scheduler: timing.NewScheduler(),
	}

	m.notifyFn, m.notifyCh = notifier.New()
	m.outputFunc.Set(func(info Info) bar.Output {
		return outputs.Text(info.String())
	})

	m.Every(1 * time.Minute)

	return m
}

func defaultClickHandler(i Info) func(bar.Event) {
	return func(e bar.Event) {
		if e.Button == bar.ButtonLeft {
			i.Toggle()
		}
	}
}

// Stream implements bar.Module.
func (m *Module) Stream(s bar.Sink) {
	enabled, err := m.provider.Get()
	outputFunc := m.outputFunc.Get().(func(Info) bar.Output)
	for {
		if !s.Error(err) {
			info := Info{
				Enabled:  enabled,
				update:   m.notifyFn,
				provider: m.provider,
			}

			s.Output(outputs.Group(outputFunc(info)).OnClick(defaultClickHandler(info)))
		}

		select {
		case <-m.outputFunc.Next():
			outputFunc = m.outputFunc.Get().(func(Info) bar.Output)
		case <-m.notifyCh:
			enabled, err = m.provider.Get()
		case <-m.scheduler.C:
			enabled, err = m.provider.Get()
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

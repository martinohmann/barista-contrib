package ip

import (
	"net"
	"time"

	"barista.run/bar"
	"barista.run/base/notifier"
	"barista.run/base/value"
	"barista.run/outputs"
	"barista.run/timing"
)

// Provider provides the current public ip of the client.
type Provider interface {
	// GetIP retrieves the current public client IP. Must return nil for both
	// return values if there is no internet connection.
	GetIP() (net.IP, error)
}

// ProviderFunc is a func that satisfies the Provider interface.
type ProviderFunc func() (net.IP, error)

// GetIP implements Provider.
func (f ProviderFunc) GetIP() (net.IP, error) {
	return f()
}

// Info contains the client's public IP address or nil of not connected.
type Info struct {
	net.IP
}

// Connected returns true when the client is connected to the internet, that is
// the IP address is not nil.
func (i Info) Connected() bool {
	return i.IP != nil
}

// Module is a module for displaying the client's current public IP address in
// the bar.
type Module struct {
	provider   Provider
	outputFunc value.Value // of func(Info) bar.Output
	notifyCh   <-chan struct{}
	notifyFn   func()
	scheduler  *timing.Scheduler
}

// New creates a new *Module with the given provider for looking up the ip
// address. By default, the module will refresh the IP address every 10
// minutes. The refresh interval can be configured using `Every`. Clicking on
// the bar output will also update the module output if not overridden.
func New(provider Provider) *Module {
	m := &Module{
		provider:  provider,
		scheduler: timing.NewScheduler(),
	}

	m.notifyFn, m.notifyCh = notifier.New()
	m.outputFunc.Set(func(info Info) bar.Output {
		if info.Connected() {
			return outputs.Text(info.String())
		}
		return outputs.Text("offline")
	})

	m.Every(10 * time.Minute)

	return m
}

func defaultClickHandler(m *Module) func(bar.Event) {
	return func(e bar.Event) {
		if e.Button == bar.ButtonLeft {
			m.Refresh()
		}
	}
}

// Stream implements bar.Module.
func (m *Module) Stream(s bar.Sink) {
	ip, err := m.provider.GetIP()
	outputFunc := m.outputFunc.Get().(func(Info) bar.Output)
	for {
		if !s.Error(err) {
			info := Info{
				IP: ip,
			}

			s.Output(outputs.Group(outputFunc(info)).OnClick(defaultClickHandler(m)))
		}

		select {
		case <-m.outputFunc.Next():
			outputFunc = m.outputFunc.Get().(func(Info) bar.Output)
		case <-m.notifyCh:
			ip, err = m.provider.GetIP()
		case <-m.scheduler.C:
			ip, err = m.provider.GetIP()
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

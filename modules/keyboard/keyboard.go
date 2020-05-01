package keyboard

import (
	"sync"
	"time"

	"barista.run/bar"
	"barista.run/base/notifier"
	"barista.run/base/value"
	l "barista.run/logging"
	"barista.run/outputs"
	"barista.run/timing"
	"golang.org/x/time/rate"
)

// Provider provides the current keyboard layout and is also able to change it.
type Provider interface {
	// GetLayout retrieves the name of the currently active keyboard layout.
	GetLayout() (string, error)

	// SetLayout sets a new keyboard layout.
	SetLayout(layout string) error
}

// Controller can switch between keyboard layouts.
type Controller interface {
	// Next switches to the next layout in the layout list. This will wrap
	// around if the last layout is reached.
	Next()

	// Next switches to the previous layout in the layout list. This will wrap
	// around if the first layout is reached.
	Previous()

	// SetLayout sets a new layout. Layout strings that were not configured on
	// the module are ignored.
	SetLayout(layout string)

	// GetLayouts returns all layouts that are configured on the keyboard module
	// instance that the controller belongs to.
	GetLayouts() []string
}

// Layout contains the name of the currently set keyboard layout. It also
// exposes a Controller to switch between available layouts.
type Layout struct {
	Controller

	// Name is the name of the keyboard layout, e.g "us".
	Name string
}

// String implements fmt.Stringer.
func (l Layout) String() string {
	return l.Name
}

type controller struct {
	sync.Mutex
	layoutMap map[string]int
	layouts   []string
	current   int
	provider  Provider
	update    func()
}

func newController(provider Provider, layouts []string, updateFn func()) *controller {
	c := &controller{
		layouts:   layouts,
		layoutMap: make(map[string]int),
		provider:  provider,
		update:    updateFn,
	}

	for i, layout := range layouts {
		c.layoutMap[layout] = i
	}

	currentLayout, _ := provider.GetLayout()

	// Set the current layout as active, add it to the list of layouts if not
	// present yet.
	i, ok := c.layoutMap[currentLayout]
	if ok {
		c.current = i
	} else {
		c.current = len(c.layouts)
		c.layoutMap[currentLayout] = c.current
		c.layouts = append(c.layouts, currentLayout)
	}

	return c
}

func (c *controller) GetLayouts() []string {
	c.Lock()
	defer c.Unlock()
	return c.layouts
}

func (c *controller) Next() {
	c.Lock()
	defer c.Unlock()

	c.setLayout(c.current + 1)
}

func (c *controller) Previous() {
	c.Lock()
	defer c.Unlock()

	c.setLayout(c.current - 1)
}

func (c *controller) SetLayout(layout string) {
	c.Lock()
	defer c.Unlock()

	index, ok := c.layoutMap[layout]
	if !ok {
		return
	}

	c.setLayout(index)
}

func (c *controller) setLayout(index int) {
	count := len(c.layouts)

	// handle wrap around on either side
	index = (index + count) % count

	layout := c.layouts[index]

	if err := c.provider.SetLayout(layout); err != nil {
		l.Log("Error setting keyboard layout: %v", err)
		return
	}

	c.current = index

	c.update()
}

// Module is a module for displaying and interacting with the keyboard layout
// that is configured by the user.
type Module struct {
	controller Controller
	provider   Provider
	outputFunc value.Value // of func(Info) bar.Output
	notifyCh   <-chan struct{}
	notifyFn   func()
	scheduler  *timing.Scheduler
}

// New creates a new *Module with given keyboard provider. By default, the
// lists of layouts is cycled through whenever the keyboard layout display in
// the bar is clicked or scrolled. By default, the module will refresh every 10
// seconds. The refresh interval can be configured using `Every`.
func New(provider Provider, layouts ...string) *Module {
	m := &Module{
		provider:  provider,
		scheduler: timing.NewScheduler(),
	}

	m.notifyFn, m.notifyCh = notifier.New()
	m.controller = newController(provider, layouts, m.notifyFn)
	m.outputFunc.Set(func(layout Layout) bar.Output {
		return outputs.Text(layout.String())
	})

	m.Every(10 * time.Second)

	return m
}

// RateLimiter throttles layout updates to once every ~20ms to avoid unexpected
// behaviour.
var RateLimiter = rate.NewLimiter(rate.Every(20*time.Millisecond), 1)

func defaultClickHandler(l Layout) func(bar.Event) {
	return func(e bar.Event) {
		if !RateLimiter.Allow() {
			return
		}

		switch {
		case e.Button == bar.ButtonLeft || e.Button == bar.ScrollUp:
			l.Next()
		case e.Button == bar.ButtonRight || e.Button == bar.ScrollDown:
			l.Previous()
		}
	}
}

// Stream implements bar.Module.
func (m *Module) Stream(s bar.Sink) {
	layout, err := m.provider.GetLayout()
	outputFunc := m.outputFunc.Get().(func(Layout) bar.Output)
	for {
		if !s.Error(err) {
			l := Layout{
				Name:       layout,
				Controller: m.controller,
			}

			s.Output(outputs.Group(outputFunc(l)).OnClick(defaultClickHandler(l)))
		}

		select {
		case <-m.outputFunc.Next():
			outputFunc = m.outputFunc.Get().(func(Layout) bar.Output)
		case <-m.notifyCh:
			layout, err = m.provider.GetLayout()
		case <-m.scheduler.C:
			layout, err = m.provider.GetLayout()
		}
	}
}

// Output updates the output format func.
func (m *Module) Output(format func(Layout) bar.Output) *Module {
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

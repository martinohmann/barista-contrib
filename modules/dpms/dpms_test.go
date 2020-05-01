package dpms

import (
	"errors"
	"sync"
	"testing"

	"barista.run/bar"
	"barista.run/outputs"
	testBar "barista.run/testing/bar"
)

type testProvider struct {
	sync.Mutex
	err     error
	enabled bool
}

func (p *testProvider) Get() (bool, error) {
	p.Lock()
	defer p.Unlock()
	if p.err != nil {
		return false, p.err
	}

	return p.enabled, nil
}

func (p *testProvider) Set(enabled bool) error {
	p.Lock()
	defer p.Unlock()
	if p.err != nil {
		return p.err
	}

	p.enabled = enabled
	return nil
}

func (p *testProvider) setError(err error) {
	p.Lock()
	defer p.Unlock()
	p.err = err
}

func TestModule(t *testing.T) {
	testBar.New(t)

	testProvider := &testProvider{
		enabled: true,
	}

	m := New(testProvider)
	testBar.Run(m)

	out := testBar.NextOutput("on start")
	out.AssertText([]string{"dpms enabled"})
	_ = testProvider.Set(false)
	m.Refresh()
	out = testBar.NextOutput("dpms disabled")
	out.AssertText([]string{"dpms disabled"})
	_ = testProvider.Set(true)
	testBar.Tick()
	out = testBar.NextOutput("reenabled")
	out.AssertText([]string{"dpms enabled"})

	out.At(0).Click(bar.Event{Button: bar.ButtonLeft})
	out = testBar.NextOutput("disabled via default click handler")
	out.AssertText([]string{"dpms disabled"})

	m.Output(func(info Info) bar.Output {
		return outputs.Textf("dpms: %v", info.Enabled).
			OnClick(func(e bar.Event) {
				switch e.Button {
				case bar.ButtonLeft:
					info.Enable()
				case bar.ButtonRight:
					info.Disable()
				case bar.ScrollUp:
					info.Toggle()
				}
			})
	})

	out = testBar.NextOutput("on output format change")

	out.At(0).Click(bar.Event{Button: bar.ButtonLeft})
	out = testBar.NextOutput("enable")
	out.AssertText([]string{"dpms: true"})

	out.At(0).Click(bar.Event{Button: bar.ButtonRight})
	out = testBar.NextOutput("disable")
	out.AssertText([]string{"dpms: false"})

	out.At(0).Click(bar.Event{Button: bar.ScrollUp})
	out = testBar.NextOutput("toggle")
	out.AssertText([]string{"dpms: true"})

	testProvider.setError(errors.New("whoops"))
	testBar.Tick()
	out = testBar.NextOutput("error")
	out.AssertError()

	testProvider.setError(nil)
	testBar.Tick()
	out = testBar.NextOutput("error")
	out.AssertText([]string{"dpms: true"})
}

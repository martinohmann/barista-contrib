package ip

import (
	"errors"
	"net"
	"sync"
	"testing"

	"barista.run/bar"
	"barista.run/outputs"
	testBar "barista.run/testing/bar"
)

type testProvider struct {
	sync.Mutex
	err error
	ip  net.IP
}

func (p *testProvider) GetIP() (net.IP, error) {
	p.Lock()
	defer p.Unlock()
	if p.err != nil {
		return nil, p.err
	}

	return p.ip, nil
}

func (p *testProvider) setIP(ip net.IP) {
	p.Lock()
	defer p.Unlock()
	p.ip = ip
}

func (p *testProvider) setError(err error) {
	p.Lock()
	defer p.Unlock()
	p.err = err
}

func TestModule(t *testing.T) {
	testBar.New(t)

	testProvider := &testProvider{
		ip: net.ParseIP("127.0.0.1"),
	}

	m := New(testProvider)
	testBar.Run(m)

	out := testBar.NextOutput("on start")
	out.AssertText([]string{"127.0.0.1"})
	testProvider.setIP(net.ParseIP("1.1.1.1"))
	m.Refresh()
	out = testBar.NextOutput("ip changed")
	out.AssertText([]string{"1.1.1.1"})
	testProvider.setIP(nil)
	testBar.Tick()
	out = testBar.NextOutput("disconnected")
	out.AssertText([]string{"offline"})

	testProvider.setError(errors.New("whoops"))

	testBar.Tick()
	out = testBar.NextOutput("next layout")
	out.AssertError()

	testProvider.setError(nil)

	m.Output(func(info Info) bar.Output {
		return outputs.Textf("ip: %s", info.IP)
	})

	out = testBar.NextOutput("outputFunc changed")
	out.AssertText([]string{"Error"})
	testBar.Tick()

	out = testBar.NextOutput("next interval")
	out.AssertText([]string{"ip: <nil>"})
	testBar.Tick()

	testProvider.setIP(net.ParseIP("10.10.10.10"))
	out = testBar.NextOutput("ip changed")
	out.AssertText([]string{"ip: 10.10.10.10"})

	testProvider.setIP(net.ParseIP("20.20.20.20"))

	out.At(0).Click(bar.Event{Button: bar.ButtonLeft})
	out = testBar.NextOutput("click")
	out.AssertText([]string{"ip: 20.20.20.20"})
}

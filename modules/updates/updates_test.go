package updates

import (
	"testing"

	"barista.run/bar"
	"barista.run/outputs"
	testBar "barista.run/testing/bar"
)

func TestModule(t *testing.T) {
	testBar.New(t)

	p := ProviderFunc(func() func() (Info, error) {
		var i int
		return func() (Info, error) {
			i++
			return Info{Updates: i}, nil
		}
	}())

	m := New(p)
	testBar.Run(m)

	testBar.LatestOutput().AssertText([]string{"1 update"})
	testBar.Tick()
	testBar.LatestOutput().AssertText([]string{"2 updates"})

	m.Output(func(info Info) bar.Output {
		return outputs.Textf("%d", info.Updates)
	})

	testBar.LatestOutput().AssertText([]string{"2"})
	m.Refresh()
	testBar.LatestOutput().AssertText([]string{"3"})
}

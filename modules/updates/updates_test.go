package updates

import (
	"testing"

	"barista.run/bar"
	"barista.run/outputs"
	testBar "barista.run/testing/bar"
	"github.com/stretchr/testify/assert"
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

func TestPackageDetailsString(t *testing.T) {
	details := PackageDetails{
		{PackageName: "foo", CurrentVersion: "v1.0", TargetVersion: "v1.1"},
		{PackageName: "bar", CurrentVersion: "v0.1", TargetVersion: "v1.0"},
	}

	expected := "foo v1.0 -> v1.1\nbar v0.1 -> v1.0"

	assert.Equal(t, expected, details.String())
}

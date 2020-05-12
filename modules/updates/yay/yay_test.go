package yay

import (
	"errors"
	"testing"

	"github.com/martinohmann/barista-contrib/internal/exec"
	"github.com/martinohmann/barista-contrib/modules/updates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	restore := exec.FakeCommandOutput(func(cmd exec.Cmd) ([]byte, error) {
		if cmd.Matches("yay", "-Qu") {
			return []byte("foo 1.0 -> 1.2\nbar 2.0 -> 2.1\n"), nil
		}

		return nil, errors.New("invalid command")
	})
	defer restore()

	p := New()

	info, err := p.Updates()
	require.NoError(t, err)

	expected := updates.Info{
		Updates: 2,
		PackageDetails: updates.PackageDetails{
			{PackageName: "foo", CurrentVersion: "1.0", TargetVersion: "1.2"},
			{PackageName: "bar", CurrentVersion: "2.0", TargetVersion: "2.1"},
		},
	}

	assert.Equal(t, expected, info)
}

func TestProvider_AUROnly(t *testing.T) {
	restore := exec.FakeCommandOutput(func(cmd exec.Cmd) ([]byte, error) {
		if cmd.Matches("yay", "-Qua") {
			return []byte("foo 1.0 -> 1.2\nbar 2.0 -> 2.1\n"), nil
		}

		return nil, errors.New("invalid command")
	})
	defer restore()

	p := New(AUROnly)

	info, err := p.Updates()
	require.NoError(t, err)

	expected := updates.Info{
		Updates: 2,
		PackageDetails: updates.PackageDetails{
			{PackageName: "foo", CurrentVersion: "1.0", TargetVersion: "1.2"},
			{PackageName: "bar", CurrentVersion: "2.0", TargetVersion: "2.1"},
		},
	}

	assert.Equal(t, expected, info)
}

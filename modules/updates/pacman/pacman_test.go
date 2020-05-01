package pacman

import (
	"testing"

	"github.com/martinohmann/barista-contrib/modules/updates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePackageDetails(t *testing.T) {
	raw := []byte(`
xf86-video-intel 1:2.99.917+904+gf2853658-1 -> 1:2.99.917+906+g846b53da-1
xmlsec 1.2.29-1 -> 1.2.30-1
xorgproto 2019.2-2 -> 2020.1-1
zip 3.0-8 -> 3.0-9

	`)

	expected := updates.PackageDetails{
		{PackageName: "xf86-video-intel", CurrentVersion: "1:2.99.917+904+gf2853658-1", TargetVersion: "1:2.99.917+906+g846b53da-1"},
		{PackageName: "xmlsec", CurrentVersion: "1.2.29-1", TargetVersion: "1.2.30-1"},
		{PackageName: "xorgproto", CurrentVersion: "2019.2-2", TargetVersion: "2020.1-1"},
		{PackageName: "zip", CurrentVersion: "3.0-8", TargetVersion: "3.0-9"},
	}

	details, err := parsePackageDetails(raw)
	require.NoError(t, err)
	assert.Equal(t, expected, details)
}

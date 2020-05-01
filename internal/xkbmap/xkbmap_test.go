package xkbmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQueryOutput(t *testing.T) {
	raw := []byte(`rules:      evdev
model:      pc105
layout:     us
`)

	expected := Info{
		Rules:  "evdev",
		Model:  "pc105",
		Layout: "us",
	}

	assert.Equal(t, expected, parseQueryOutput(raw))
}

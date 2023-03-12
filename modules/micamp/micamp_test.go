package micamp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratePercentageBar(t *testing.T) {
	require.Equal(t, "0%   .......... 🎙", generatePercentageBar(0.00001))
	require.Equal(t, "1%   .......... 🎙", generatePercentageBar(0.01))
	require.Equal(t, "12%  :......... 🎙", generatePercentageBar(0.12))
	require.Equal(t, "22%  ::........ 🎙", generatePercentageBar(0.22))
	require.Equal(t, "32%  :::....... 🎙", generatePercentageBar(0.32))
	require.Equal(t, "42%  ::::...... 🎙", generatePercentageBar(0.42))
	require.Equal(t, "52%  :::::..... 🎙", generatePercentageBar(0.52))
	require.Equal(t, "62%  ::::::.... 🎙", generatePercentageBar(0.62))
	require.Equal(t, "72%  :::::::... 🎙", generatePercentageBar(0.72))
	require.Equal(t, "82%  ::::::::.. 🎙", generatePercentageBar(0.82))
	require.Equal(t, "92%  :::::::::. 🎙", generatePercentageBar(0.92))
	require.Equal(t, "99%  :::::::::. 🎙", generatePercentageBar(0.99))
	require.Equal(t, "100% :::::::::: 🎙", generatePercentageBar(1))
	require.Equal(t, "ERR 200% (amp=2.000) 🎙", generatePercentageBar(2))
}

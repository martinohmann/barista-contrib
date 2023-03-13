/*
 * Copyright (c) 2023, Simon Gottschlag <simon@gottschlag.se>
 *
 * SPDX-License-Identifier: MIT
 */

package micamp

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWavSampler(t *testing.T) {
	wavSampler := newWavSampler()
	wavSampler.Write(testGenerateFakeAmplitudes(t, 8))
	require.Equal(t, float64(4), wavSampler.amplitude())

	wavSampler = newWavSampler()
	wavSampler.Write(testGenerateFakeAmplitudes(t, 16))
	require.Equal(t, float64(8), wavSampler.amplitude())

	wavSampler = newWavSampler()
	wavSampler.Write(testGenerateFakeAmplitudes(t, 32))
	require.Equal(t, float64(16), wavSampler.amplitude())

	wavSampler = newWavSampler()
	wavSampler.Write(testGenerateFakeAmplitudes(t, 64))
	require.Equal(t, float64(32), wavSampler.amplitude())

	wavSampler = newWavSampler()
	wavSampler.Write(testGenerateFakeAmplitudes(t, 8))
	wavSampler.Write(testGenerateFakeAmplitudes(t, 16))
	require.Equal(t, float64(6), wavSampler.amplitude())
}

func TestWavSamplerTimeout(t *testing.T) {
	wavSampler := newWavSampler()
	wavSampler.Write(testGenerateFakeAmplitudes(t, 8))
	require.Equal(t, float64(4), wavSampler.amplitude())
	wavSampler.lastUpdate = time.Now().Add(-1 * time.Second)
	require.Equal(t, float64(4), wavSampler.amplitude())
	wavSampler.lastUpdate = time.Now().Add(-3 * time.Second)
	require.True(t, math.IsNaN(wavSampler.amplitude()))
}

func testGenerateFakeAmplitudes(t *testing.T, c int) []float32 {
	t.Helper()

	result := []float32{}
	for i := 1; i < c; i++ {
		result = append(result, float32(i))
	}
	return result
}

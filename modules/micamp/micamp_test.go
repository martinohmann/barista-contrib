package micamp

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	testbar "barista.run/testing/bar"
	"barista.run/timing"
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

type testProvider struct {
	t          *testing.T
	closeCount int
}

func (p *testProvider) close() {
	p.t.Helper()

	p.closeCount++
}

type testSampler struct {
	t                *testing.T
	currentAmplitude float64
	amplitudeCount   int
}

func (s *testSampler) Write(p []float32) (int, error) {
	s.t.Helper()

	return 0, nil
}

func (s *testSampler) amplitude() float64 {
	s.t.Helper()

	s.amplitudeCount++

	return s.currentAmplitude
}

func TestModule(t *testing.T) {
	testbar.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &testProvider{
		t: t,
	}
	s := &testSampler{
		t:                t,
		currentAmplitude: 0,
	}
	m := &Module{
		ctx:       ctx,
		scheduler: timing.NewScheduler().Every(5 * time.Millisecond),
		newMicProviderFn: func() (provider, error) {
			return p, nil
		},
		wavSampler: s,
	}
	m.outputFunc.Set(defaultOutputFunc)

	testbar.Run(m)

	testbar.Tick()
	out := testbar.NextOutput("initial state")
	out.AssertText([]string{"0%   .......... 🎙"})

	s.currentAmplitude = 0.001
	testbar.Tick()
	out = testbar.NextOutput("on change, 0.001")
	out.AssertText([]string{"0%   .......... 🎙"})

	s.currentAmplitude = 0.5
	testbar.Tick()
	out = testbar.NextOutput("on change, 0.500")
	out.AssertText([]string{"50%  :::::..... 🎙"})

	s.currentAmplitude = 0
	testbar.Tick()
	out = testbar.NextOutput("on change, error")
	out.AssertText([]string{"0%   .......... 🎙"})

	s.currentAmplitude = math.NaN()
	testbar.Tick()
	out = testbar.NextOutput("on change, error")
	out.AssertText([]string{"NaN .......... 🎙"})

	m.newMicProviderFn = func() (provider, error) {
		return nil, fmt.Errorf("ze-failure")
	}
	testbar.Tick()
	out = testbar.NextOutput("on change, provider nil")
	require.Nil(t, m.micProvider)
	out.AssertText([]string{"NaN .......... 🎙"})

	testbar.Tick()
	out = testbar.NextOutput("on change, provider nil")
	require.Nil(t, m.micProvider)
	out.AssertText([]string{"NaN .......... 🎙"})

	m.newMicProviderFn = func() (provider, error) {
		return p, nil
	}
	s.currentAmplitude = 0.01
	testbar.Tick()
	out = testbar.NextOutput("on change, back to normal")
	out.AssertText([]string{"1%   .......... 🎙"})

	s.currentAmplitude = 1
	testbar.Tick()
	out = testbar.NextOutput("on change, 1")
	out.AssertText([]string{"100% :::::::::: 🎙"})

	s.currentAmplitude = 2
	testbar.Tick()
	out = testbar.NextOutput("on change, error")
	out.AssertText([]string{"ERR 200% (amp=2.000) 🎙"})

	s.currentAmplitude = 0
	testbar.Tick()
	out = testbar.NextOutput("on change, error")
	out.AssertText([]string{"0%   .......... 🎙"})
}

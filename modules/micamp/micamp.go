/*
 * Copyright (c) 2023, Simon Gottschlag <simon@gottschlag.se>
 *
 * SPDX-License-Identifier: MIT
 */

package micamp

import (
	"context"
	"fmt"
	"math"
	"time"

	"barista.run/bar"
	"barista.run/base/value"
	"barista.run/colors"
	"barista.run/outputs"
	"barista.run/timing"
)

type provider interface {
	close()
}

type sampler interface {
	Write(p []float32) (int, error)
	amplitude() float64
}

type module struct {
	outputFunc       value.Value
	ctx              context.Context
	scheduler        *timing.Scheduler
	micProvider      provider
	newMicProviderFn func() (provider, error)
	wavSampler       sampler
}

func generatePercentageBar(amp float64) string {
	switch percentage := int(amp * 100); {
	case percentage >= 0 && percentage < 10:
		return fmt.Sprintf("%d%%   .......... 🎙", percentage)
	case percentage >= 10 && percentage < 20:
		return fmt.Sprintf("%d%%  :......... 🎙", percentage)
	case percentage >= 20 && percentage < 30:
		return fmt.Sprintf("%d%%  ::........ 🎙", percentage)
	case percentage >= 30 && percentage < 40:
		return fmt.Sprintf("%d%%  :::....... 🎙", percentage)
	case percentage >= 40 && percentage < 50:
		return fmt.Sprintf("%d%%  ::::...... 🎙", percentage)
	case percentage >= 50 && percentage < 60:
		return fmt.Sprintf("%d%%  :::::..... 🎙", percentage)
	case percentage >= 60 && percentage < 70:
		return fmt.Sprintf("%d%%  ::::::.... 🎙", percentage)
	case percentage >= 70 && percentage < 80:
		return fmt.Sprintf("%d%%  :::::::... 🎙", percentage)
	case percentage >= 80 && percentage < 90:
		return fmt.Sprintf("%d%%  ::::::::.. 🎙", percentage)
	case percentage >= 90 && percentage < 100:
		return fmt.Sprintf("%d%%  :::::::::. 🎙", percentage)
	case percentage == 100:
		return fmt.Sprintf("%d%% :::::::::: 🎙", percentage)
	default:
		return fmt.Sprintf("ERR %d%% (amp=%0.3f) 🎙", percentage, amp)
	}
}

var defaultOutputFunc = func(amp float64) bar.Output {
	if math.IsNaN(amp) {
		return outputs.Text("NaN .......... 🎙").Color(colors.Hex("#f00"))
	}

	if amp == float64(0) {
		return outputs.Text("0%   .......... 🎙").Color(colors.Hex("#ff0"))
	}

	return outputs.Text(generatePercentageBar(amp))
}

// New creates the microphone amplitude (micamp) module for barista.
// It is used to give a visual indication that audio is passing
// through the microphone.
//
// Default output will be:
//
//   - When the mic is muted or doesn't receive any audio:
//     NaN .......... 🎙
//
//   - When the mic is receiving audio:
//     0%   .......... 🎙
//     50%  :::::..... 🎙
//     100% :::::::::: 🎙
//
//   - When the amplitude isn't between 0-100:
//     ERR 200% (amp=2.000) 🎙
//
// Parameters:
//
//   - ctx:
//     context which when done will stop the stream and gracefully
//     shut down the pulse audio client.
//
//   - micSourceNamePrefix:
//     The prefix of the microphone name as seen by the pulse audio
//     description (which can be found using `pactl list sources`).
//     If it's empty (`""`) then the pulse audio default source will
//     be used.
func New(ctx context.Context, micSourceNamePrefix string) *module {
	wavSampler := newWavSampler()
	m := &module{
		ctx:       ctx,
		scheduler: timing.NewScheduler().Every(1 * time.Second),
		newMicProviderFn: func() (provider, error) {
			return newPulseProvider(micSourceNamePrefix, wavSampler)
		},
		wavSampler: wavSampler,
	}

	m.outputFunc.Set(defaultOutputFunc)

	return m
}

func (m *module) Stream(s bar.Sink) {
	defer m.close()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-m.scheduler.C:
			m.process(s)
		}
	}
}

func (m *module) close() {
	if m.micProvider == nil {
		return
	}

	m.micProvider.close()
}

func (m *module) process(s bar.Sink) {
	if !m.isProviderReady(false) {
		m.output(s, math.NaN())
		return
	}

	amp := m.wavSampler.amplitude()
	if amp == float64(0) || math.IsNaN(amp) {
		m.isProviderReady(true)
		m.output(s, amp)
		return
	}

	m.output(s, amp)

	return
}

func (m *module) output(s bar.Sink, amp float64) {
	format := m.outputFunc.Get().(func(float64) bar.Output)
	s.Output(format(amp))
}

func (m *module) isProviderReady(force bool) bool {
	if m.micProvider != nil && !force {
		return true
	}

	provider, err := m.newMicProviderFn()
	if err != nil {
		m.close()
		m.micProvider = nil
		return false
	}

	m.close()
	m.micProvider = provider

	return true
}

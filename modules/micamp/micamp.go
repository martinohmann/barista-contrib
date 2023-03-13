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

type Module struct {
	outputFunc       value.Value
	ctx              context.Context
	scheduler        *timing.Scheduler
	micProvider      provider
	newMicProviderFn func() (provider, error)
	wavSampler       sampler
}

func generatePercentageBar(amp float64) string {
	// 0%   .......... 🎙
	// 2%   .......... 🎙
	// 12%  :......... 🎙
	// 22%  ::........ 🎙
	// ...
	// 92%  :::::::::. 🎙
	// 100% :::::::::: 🎙

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

func New(ctx context.Context, micSourceNamePrefix string) *Module {
	wavSampler := newWavSampler()
	m := &Module{
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

func (m *Module) Stream(s bar.Sink) {
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

func (m *Module) close() {
	if m.micProvider == nil {
		return
	}

	m.micProvider.close()
}

func (m *Module) process(s bar.Sink) {
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

func (m *Module) output(s bar.Sink, amp float64) {
	format := m.outputFunc.Get().(func(float64) bar.Output)
	s.Output(format(amp))
}

func (m *Module) isProviderReady(force bool) bool {
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

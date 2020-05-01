package cpufreq

import (
	"time"

	"barista.run/bar"
	"barista.run/base/notifier"
	"barista.run/base/value"
	"barista.run/outputs"
	"barista.run/timing"
	"github.com/martinlindhe/unit"
	"github.com/prometheus/procfs/sysfs"
)

type Provider interface {
	GetCPUFrequency() (Info, error)
}

type Info struct {
	Stats []sysfs.SystemCPUCpufreqStats
}

func (i Info) NumCPUs() int {
	return len(i.Stats)
}

func (i Info) Freq(cpu int) unit.Frequency {
	if cpu < 0 || cpu >= i.NumCPUs() {
		return 0
	}

	freq := i.Stats[cpu].ScalingCurrentFrequency
	if freq == nil {
		return 0
	}

	return unit.Frequency(float64(*freq) * 1000)
}

func (i Info) AverageFreq() unit.Frequency {
	var count int
	var sum uint64
	for _, stat := range i.Stats {
		if stat.ScalingCurrentFrequency == nil {
			continue
		}

		count++
		sum += *stat.ScalingCurrentFrequency
	}

	return unit.Frequency(float64(sum) / float64(count) * 1000)
}

type Module struct {
	provider   Provider
	outputFunc value.Value // of func(Info) bar.Output
	notifyCh   <-chan struct{}
	notifyFn   func()
	scheduler  *timing.Scheduler
}

func New(provider Provider) *Module {
	m := &Module{
		provider:  provider,
		scheduler: timing.NewScheduler(),
	}

	m.notifyFn, m.notifyCh = notifier.New()
	m.outputFunc.Set(func(info Info) bar.Output {
		return outputs.Textf("%.2fGHz", info.AverageFreq().Gigahertz())
	})

	m.Every(10 * time.Second)

	return m
}

func (m *Module) Stream(s bar.Sink) {
	info, err := m.provider.GetCPUFrequency()
	outputFunc := m.outputFunc.Get().(func(Info) bar.Output)
	for {
		if s.Error(err) {
			continue
		}

		s.Output(outputFunc(info))

		select {
		case <-m.outputFunc.Next():
			outputFunc = m.outputFunc.Get().(func(Info) bar.Output)
		case <-m.notifyCh:
			info, err = m.provider.GetCPUFrequency()
		case <-m.scheduler.C:
			info, err = m.provider.GetCPUFrequency()
		}
	}
}

func (m *Module) Output(format func(Info) bar.Output) *Module {
	m.outputFunc.Set(format)
	return m
}

func (m *Module) Every(interval time.Duration) *Module {
	if interval == 0 {
		m.scheduler.Stop()
	} else {
		m.scheduler.Every(interval)
	}
	return m
}

func (m *Module) Refresh() {
	m.notifyFn()
}

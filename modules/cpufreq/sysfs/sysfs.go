package sysfs

import (
	"github.com/martinohmann/barista-contrib/modules/cpufreq"
	"github.com/prometheus/procfs/sysfs"
)

// New creates a new *cpufreq.Module using sysfs as CPU frequency provider.
func New(fs sysfs.FS) *cpufreq.Module {
	return cpufreq.New(&provider{
		fs: fs,
	})
}

type provider struct {
	fs sysfs.FS
}

// Set implements cpufreq.Provider.
func (p *provider) GetCPUFrequency() (cpufreq.Info, error) {
	stats, err := p.fs.SystemCpufreq()
	if err != nil {
		return cpufreq.Info{}, err
	}

	return cpufreq.Info{Stats: stats}, nil
}

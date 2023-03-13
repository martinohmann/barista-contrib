package micamp

import (
	"github.com/jfreymuth/pulse"
)

type sampler interface {
	Write(p []float32) (int, error)
	amplitude() float64
}

type pulseProvider struct {
	pulseClient *pulse.Client
	pulseStream *pulse.RecordStream
	wavSampler  sampler
}

func newPulseProvider(micSourceName string, wavSampler sampler) (*pulseProvider, error) {
	pulseClient, err := pulse.NewClient(pulse.ClientApplicationName("barista-micamp"))
	if err != nil {
		return nil, err
	}

	pulseStream, err := pulseClient.NewRecord(pulse.Float32Writer(wavSampler.Write))
	if err != nil {
		return nil, err
	}

	pulseStream.Start()

	return &pulseProvider{
		pulseClient: pulseClient,
		pulseStream: pulseStream,
		wavSampler:  wavSampler,
	}, nil
}

func (p *pulseProvider) close() {
	p.pulseStream.Stop()
	p.pulseStream.Close()
	p.pulseClient.Close()
}

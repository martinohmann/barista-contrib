/*
 * Copyright (c) 2023, Simon Gottschlag <simon@gottschlag.se>
 *
 * SPDX-License-Identifier: MIT
 */

package micamp

import (
	"fmt"
	"strings"

	"github.com/jfreymuth/pulse"
)

type pulseProvider struct {
	pulseClient *pulse.Client
	pulseStream *pulse.RecordStream
	wavSampler  sampler
}

func newPulseProvider(micSourceNamePrefix string, wavSampler sampler) (*pulseProvider, error) {
	pulseClient, err := pulse.NewClient(pulse.ClientApplicationName("barista-micamp"))
	if err != nil {
		return nil, err
	}

	pulseSource, err := getPulseSource(pulseClient, micSourceNamePrefix)
	if err != nil {
		return nil, err
	}

	pulseStream, err := pulseClient.NewRecord(pulse.Float32Writer(wavSampler.Write), pulse.RecordSource(pulseSource))
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

func getPulseSource(pulseClient *pulse.Client, micSourceNamePrefix string) (*pulse.Source, error) {
	if micSourceNamePrefix == "" {
		return pulseClient.DefaultSource()
	}

	sources, err := pulseClient.ListSources()
	if err != nil {
		return nil, err
	}

	for _, source := range sources {
		if strings.HasPrefix(source.Name(), micSourceNamePrefix) {
			return source, nil
		}
	}

	return nil, fmt.Errorf("unable to find any source with the prefix %q", micSourceNamePrefix)
}

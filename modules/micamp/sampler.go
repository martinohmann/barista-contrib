/*
 * Copyright (c) 2023, Simon Gottschlag <simon@gottschlag.se>
 *
 * SPDX-License-Identifier: MIT
 */

package micamp

import (
	"container/ring"
	"math"
	"sync"
	"time"
)

type wavSampler struct {
	mu         sync.Mutex
	r          *ring.Ring
	lastUpdate time.Time
}

func newWavSampler() *wavSampler {
	return &wavSampler{
		r:          ring.New(32),
		lastUpdate: time.Now(),
	}
}

func (s *wavSampler) Write(p []float32) (int, error) {
	averageAmplitude := calculateAverage32(p)
	s.mu.Lock()
	s.r.Value = averageAmplitude
	s.r = s.r.Next()
	s.lastUpdate = time.Now()
	s.mu.Unlock()

	return len(p), nil
}
func (s *wavSampler) amplitude() float64 {
	// Note: If we haven't received any updates using Write() for
	//       2 seconds or more, return Not a Number (NaN).
	if time.Since(s.lastUpdate) > 2*time.Second {
		return math.NaN()
	}

	buf := []float64{}
	s.mu.Lock()
	s.r.Do(func(a interface{}) {
		if a == nil {
			return
		}

		v, ok := a.(float64)
		if ok {
			buf = append(buf, v)
		}
	})
	s.mu.Unlock()

	return calculateAverage64(buf)
}

func calculateAverage32(p []float32) float64 {
	sum := float64(0)
	for _, f := range p {
		sum += float64(math.Abs(float64(f)))
	}
	return sum / float64(len(p))
}

func calculateAverage64(p []float64) float64 {
	sum := float64(0)
	for _, f := range p {
		sum += float64(math.Abs(f))
	}
	return sum / float64(len(p))
}

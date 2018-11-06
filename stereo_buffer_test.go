package main

import (
	"math/rand"
	"reflect"
	"testing"
)

type stream struct {
	samples [][2]float64
	index   int
}

func (s *stream) Stream(samples [][2]float64) (int, bool) {
	samplesToCopy := len(samples)
	remainingLength := len(s.samples[s.index:])
	if remainingLength == 0 {
		return 0, false
	}
	if remainingLength >= samplesToCopy {
		copy(samples, s.samples[s.index:s.index+samplesToCopy])
		s.index += samplesToCopy
		return samplesToCopy, true
	}
	copy(samples, s.samples[s.index:])
	s.index = 0
	return remainingLength, true
}

func (s *stream) Err() error {
	return nil
}

func Test_stereoBuffer_StreamLeft_StreamRight(t *testing.T) {
	tests := []struct {
		name                 string
		samples              [][2]float64
		leftReads            int
		rightReads           int
		expectedLeftSamples  [][2]float64
		expectedRightSamples [][2]float64
	}{
		{
			name: "parallel reads",
			samples: [][2]float64{
				[2]float64{2, 3},
				[2]float64{1, 6},
				[2]float64{22, 3},
				[2]float64{2, 33},
				[2]float64{16, 7},
				[2]float64{4, 30},
			},
			leftReads:  6,
			rightReads: 6,
			expectedLeftSamples: [][2]float64{
				[2]float64{2, 2},
				[2]float64{1, 1},
				[2]float64{22, 22},
				[2]float64{2, 2},
				[2]float64{16, 16},
				[2]float64{4, 4},
			},
			expectedRightSamples: [][2]float64{
				[2]float64{3, 3},
				[2]float64{6, 6},
				[2]float64{3, 3},
				[2]float64{33, 33},
				[2]float64{7, 7},
				[2]float64{30, 30},
			},
		},
		{
			name: "partial left reads",
			samples: [][2]float64{
				[2]float64{2, 3},
				[2]float64{1, 6},
				[2]float64{22, 3},
				[2]float64{2, 33},
				[2]float64{16, 7},
				[2]float64{4, 30},
			},
			leftReads:  3,
			rightReads: 6,
			expectedLeftSamples: [][2]float64{
				[2]float64{2, 2},
				[2]float64{1, 1},
				[2]float64{22, 22},
			},
			expectedRightSamples: [][2]float64{
				[2]float64{3, 3},
				[2]float64{6, 6},
				[2]float64{3, 3},
				[2]float64{33, 33},
				[2]float64{7, 7},
				[2]float64{30, 30},
			},
		},
		{
			name: "partial right reads",
			samples: [][2]float64{
				[2]float64{2, 3},
				[2]float64{1, 6},
				[2]float64{22, 3},
				[2]float64{2, 33},
				[2]float64{16, 7},
				[2]float64{4, 30},
			},
			leftReads:  6,
			rightReads: 2,
			expectedLeftSamples: [][2]float64{
				[2]float64{2, 2},
				[2]float64{1, 1},
				[2]float64{22, 22},
				[2]float64{2, 2},
				[2]float64{16, 16},
				[2]float64{4, 4},
			},
			expectedRightSamples: [][2]float64{
				[2]float64{3, 3},
				[2]float64{6, 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamer := &stream{
				samples: tt.samples,
			}

			buffer := &stereoBuffer{
				streamer: streamer,
			}

			readFn := func(sampleFn func(s [][2]float64) (int, bool)) [2]float64 {
				samples := make([][2]float64, 1)
				sampleFn(samples)
				return samples[0]
			}

			var rightSamples [][2]float64
			var leftSamples [][2]float64

			var reads []func()
			for i := 0; i < tt.rightReads; i++ {
				reads = append(reads, func() {
					rightSamples = append(rightSamples, readFn(buffer.StreamRight))
				})
			}
			for i := 0; i < tt.leftReads; i++ {
				reads = append(reads, func() {
					leftSamples = append(leftSamples, readFn(buffer.StreamLeft))
				})
			}

			rand.Shuffle(len(reads), func(i, j int) {
				reads[i], reads[j] = reads[j], reads[i]
			})
			for _, r := range reads {
				r()
			}

			if !reflect.DeepEqual(tt.expectedLeftSamples, leftSamples) {
				t.Errorf("Expected left samples to equal\n%v\n\tgot\n%v", tt.expectedLeftSamples, leftSamples)
			}

			if !reflect.DeepEqual(tt.expectedRightSamples, rightSamples) {
				t.Errorf("Expected right samples to equal\n%v\n\tgot\n%v", tt.expectedRightSamples, rightSamples)
			}
		})
	}
}

func Test_stereoBuffer_StreamLeft(t *testing.T) {

}

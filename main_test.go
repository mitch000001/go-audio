package main

import (
	"reflect"
	"testing"
)

func Test_splitToMono(t *testing.T) {
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

			left, right := splitToMono(streamer)

			{
				leftSamples := make([][2]float64, tt.leftReads)

				n, ok := left.Stream(leftSamples)
				if !ok {
					t.Errorf("left stream failed")
				}
				if n != tt.leftReads {
					t.Errorf("Expected left reads to equal %d, got %d", tt.leftReads, n)
				}

				if !reflect.DeepEqual(tt.expectedLeftSamples, leftSamples) {
					t.Errorf("Expected left samples to equal\n%v\n\tgot\n%v", tt.expectedLeftSamples, leftSamples)
				}
			}

			{
				rightSamples := make([][2]float64, tt.rightReads)

				n, ok := right.Stream(rightSamples)
				if !ok {
					t.Errorf("right stream failed")
				}
				if n != tt.rightReads {
					t.Errorf("Expected right reads to equal %d, got %d", tt.rightReads, n)
				}

				if !reflect.DeepEqual(tt.expectedRightSamples, rightSamples) {
					t.Errorf("Expected right samples to equal\n%v\n\tgot\n%v", tt.expectedRightSamples, rightSamples)
				}
			}
		})
	}
}

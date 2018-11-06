package main

import "github.com/faiface/beep"

type stereoBuffer struct {
	streamer beep.Streamer
	index    int
	left     [][2]float64
	right    [][2]float64
}

func (s *stereoBuffer) stream(samples [][2]float64) (int, bool) {
	n, ok := s.streamer.Stream(samples)
	for i := 0; i < n; i++ {
		s.left = append(s.left, [2]float64{
			samples[i][0],
			samples[i][0],
		})
		s.right = append(s.right, [2]float64{
			samples[i][1],
			samples[i][1],
		})
	}
	return n, ok
}

func (s *stereoBuffer) StreamLeft(leftSamples [][2]float64) (int, bool) {
	sampleLength := len(leftSamples)
	bufferLength := len(s.left)
	if bufferLength == 0 {
		samples := make([][2]float64, sampleLength)
		n, ok := s.stream(samples)
		copy(leftSamples, s.left)
		s.left = make([][2]float64, 0)
		return n, ok
	}
	if bufferLength < sampleLength {
		diff := sampleLength - bufferLength
		samples := make([][2]float64, diff)
		s.stream(samples)
		copy(leftSamples, s.left)
		s.left = make([][2]float64, 0)
		return sampleLength, true
	}

	copy(leftSamples, s.left[:sampleLength])
	s.left = s.left[sampleLength:]
	return sampleLength, true
}

func (s *stereoBuffer) StreamRight(rightSamples [][2]float64) (int, bool) {
	sampleLength := len(rightSamples)
	bufferLength := len(s.right)
	if bufferLength == 0 {
		samples := make([][2]float64, sampleLength)
		n, ok := s.stream(samples)
		copy(rightSamples, s.right)
		s.right = make([][2]float64, 0)
		return n, ok
	}
	if bufferLength < sampleLength {
		diff := sampleLength - bufferLength
		samples := make([][2]float64, diff)
		s.stream(samples)
		copy(rightSamples, s.right)
		s.right = make([][2]float64, 0)
		return sampleLength, true
	}

	copy(rightSamples, s.right[:sampleLength])
	s.right = s.right[sampleLength:]
	return sampleLength, true
}

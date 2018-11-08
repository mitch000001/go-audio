package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/faiface/beep"
)

var noise = beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		samples[i][0] = rand.Float64()*2 - 1
		samples[i][1] = rand.Float64()*2 - 1
	}
	return len(samples), true
})

func sineWave(sr beep.SampleRate, freq float64) beep.Streamer {
	t := 0.0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			y := math.Sin(math.Pi * freq * t)
			samples[i][0] = y
			samples[i][1] = y
			t += sr.D(1).Seconds()
		}
		return sampleLength, true
	})
}

func lfo(sr beep.SampleRate, rate int, amount, cutoffFreq float64) beep.Streamer {
	lfoFreq := 1 / float64(rate)
	lfoOscillator := sineWave(sr, lfoFreq)
	sineFreq := cutoffFreq - (amount / 2)
	sineOscillator := sineWave(sr, sineFreq)
	x := 0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		lfoSamples := make([][2]float64, sampleLength)
		lfoOscillator.Stream(lfoSamples)
		sineOscillator.Stream(samples)
		for i := range samples {
			samples[i][0] = math.Cos(samples[i][0] + lfoSamples[i][0])
			samples[i][1] = math.Cos(samples[i][1] + lfoSamples[i][1])
			x += i + 1
			fmt.Println("x:", x, "samples:", samples[i])
		}
		return sampleLength, true
	})
}

// LFO is a lfo
func LFO(sr beep.SampleRate, rate int, amount, freq float64) beep.Streamer {
	Ac := 1.0
	Am := amount
	Kf := 1.0
	fc := freq
	fm := 1 / float64(rate)
	deltaF := Kf * Am

	t := 0.0
	x := 0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			y := Ac * math.Cos(2*math.Pi*fc*t+(Am*deltaF)/fm*math.Sin(2*math.Pi*fm*t))
			samples[i][0] = y
			samples[i][1] = y
			t += sr.D(1).Seconds()
			x += i + 1
			fmt.Println("x:", x, "samples:", samples[i])
		}
		return sampleLength, true
	})
}

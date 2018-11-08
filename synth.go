package main

import (
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

func sawtooth(sr beep.SampleRate, freq float64) beep.Streamer {
	a := 1.0
	t := 0.0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			samples[i][0] = 2 * (t/a - math.Floor(1/2+t/a))
			samples[i][1] = 2 * (t/a - math.Floor(1/2+t/a))
			t += sr.D(1).Seconds()
		}
		return sampleLength, true
	})
}

func lfo(sr beep.SampleRate, rate int, amount, freq float64) beep.Streamer {
	Ac := 1.0
	Am := 1.0
	Kf := amount
	fc := freq
	fm := 1 / float64(rate)
	deltaF := Kf * Am

	t := 0.0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			y := Ac * math.Cos(2*math.Pi*fc*t+(Am*deltaF)/fm*math.Sin(2*math.Pi*fm*t))
			samples[i][0] = y
			samples[i][1] = y
			t += sr.D(1).Seconds()
		}
		return sampleLength, true
	})
}

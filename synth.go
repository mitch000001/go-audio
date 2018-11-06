package main

import (
	"math"
	"math/rand"
	"time"

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

func lfo(sr beep.SampleRate, rate int, amount, cutoff float64) beep.Streamer {
	min := cutoff - amount
	freq := cutoff - amount/2
	rising := true
	if rate <= 0 {
		rate = 1
	}
	if rate > 20 {
		rate = 20
	}

	changeAfter := sr.N(500 / time.Duration(rate) * time.Millisecond)
	if changeAfter == 0 {
		changeAfter = 10
	}
	readSamples := 0
	sine := beep.Take(changeAfter, sineWave(sr, freq))
	return beep.StreamerFunc(func(samples [][2]float64) (int, bool) {
		sampleLength := len(samples)
		n, _ := sine.Stream(samples)
		readSamples += n
		if readSamples%changeAfter == 0 {
			if rising {
				freq += 20
			} else {
				freq -= 20
			}
			sine = beep.Take(changeAfter, sineWave(sr, freq))
			readSamples = 0
			if freq == cutoff {
				rising = false
			} else if freq == min {
				rising = true
			}
		}
		return sampleLength, true
	})
}

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

func wave(sr beep.SampleRate, waveForm func(t float64) float64) beep.Streamer {
	t := 0.0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			y := waveForm(t)
			samples[i][0] = y
			samples[i][1] = y
			t += sr.D(1).Seconds()
		}
		return sampleLength, true
	})
}

func sineWave(sr beep.SampleRate, freq float64) beep.Streamer {
	t := 0.0
	sineFn := sine(1, freq, 0)
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			y := sineFn(t)
			samples[i][0] = y
			samples[i][1] = y
			t += sr.D(1).Seconds()
		}
		return sampleLength, true
	})
}

func sine(amplitude float64, frequency float64, phaseShift float64) func(t float64) float64 {
	return func(t float64) float64 {
		return amplitude * math.Sin(
			2*math.Pi*frequency*t+phaseShift,
		)
	}
}

func sawtoothWave(sr beep.SampleRate, freq float64) beep.Streamer {
	t := 0.0
	sawFn := sawtooth(1, freq, 0)
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			samples[i][0] = sawFn(t)
			samples[i][1] = sawFn(t)
			t += sr.D(1).Seconds()
		}
		return sampleLength, true
	})
}

func sawtooth(amplitude float64, frequency float64, phaseShift float64) func(t float64) float64 {
	return func(t float64) float64 {
		return -(amplitude / 2 / math.Pi) * math.Atan(
			1/math.Tan(
				(t*math.Pi)/(1/frequency),
			),
		)
	}
}

func triangleWave(sr beep.SampleRate, freq float64) beep.Streamer {
	t := 0.0
	triangleFn := triangle(1, freq, 0)
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			samples[i][0] = triangleFn(t)
			samples[i][1] = triangleFn(t)
			t += sr.D(1).Seconds()
		}
		return sampleLength, true
	})
}

func triangle(amplitude float64, frequency float64, phaseShift float64) func(t float64) float64 {
	return func(t float64) float64 {
		return (2 * amplitude / math.Pi) * math.Asin(
			math.Sin(
				((2*math.Pi)/(1/frequency))*t,
			),
		)
	}
}

func squareWave(sr beep.SampleRate, freq float64) beep.Streamer {
	t := 0.0
	squareFn := square(1, freq, 0)
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		sampleLength := len(samples)
		for i := range samples {
			samples[i][0] = squareFn(t)
			samples[i][1] = squareFn(t)
			t += sr.D(1).Seconds()
		}
		return sampleLength, true
	})
}

func square(amplitude float64, frequency float64, phaseShift float64) func(t float64) float64 {
	return func(t float64) float64 {
		return amplitude / 4 * math.Copysign(1,
			math.Sin(
				2*math.Pi*frequency*t+phaseShift,
			),
		)
	}
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

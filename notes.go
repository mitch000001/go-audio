package main

import (
	"fmt"
	"math"

	"github.com/faiface/beep"
)

func pitch(semiTones int, s beep.Streamer) beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (int, bool) {
		sampleLength := len(samples)
		origSamples := make([][2]float64, sampleLength)
		n, ok := s.Stream(origSamples)
		for i := range samples[:n] {
			samples[i][0] = samples[i][0] * pitchPerSemitone(semiTones, 0)
			samples[i][1] = samples[i][1] * pitchPerSemitone(semiTones, 0)
		}
		return n, ok
	})
}

func pitchPerSemitone(semiTones int, freq float64) float64 {
	pow := float64(semiTones) / 12.0
	fmt.Println("Pow:", pow)
	return math.Pow(2, pow) * freq
}

var notes = map[string]float64{
	"C8":      4186.01,
	"B7":      3951.07,
	"A#7/Bb7": 3729.31,
	"A7":      3520.00,
	"G#7/Ab7": 3322.44,
	"G7":      3135.96,
	"F#7/Gb7": 2959.96,
	"F7":      2793.83,
	"E7":      2637.02,
	"D#7/Eb7": 2489.02,
	"D7":      2349.32,
	"C#7/Db7": 2217.46,
	"C7":      2093.00,
	"B6":      1975.53,
	"A#6/Bb6": 1864.66,
	"A6":      1760.00,
	"G#6/Ab6": 1661.22,
	"G6":      1567.98,
	"F#6/Gb6": 1479.98,
	"F6":      1396.91,
	"E6":      1318.51,
	"D#6/Eb6": 1244.51,
	"D6":      1174.66,
	"C#6/Db6": 1108.73,
	"C6":      1046.50,
	"B5":      987.767,
	"A#5/Bb5": 932.328,
	"A5":      880.000,
	"G#5/Ab5": 830.609,
	"G5":      783.991,
	"F#5/Gb5": 739.989,
	"F5":      698.456,
	"E5":      659.255,
	"D#5/Eb5": 622.254,
	"D5":      587.330,
	"C#5/Db5": 554.365,
	"C5":      523.251,
	"B4":      493.883,
	"A#4/Bb4": 466.164,
	"A4":      440.000,
	"G#4/Ab4": 415.305,
	"G4":      391.995,
	"F#4/Gb4": 369.994,
	"F4":      349.228,
	"E4":      329.628,
	"D#4/Eb4": 311.127,
	"D4":      293.665,
	"C#4/Db4": 277.183,
	"C4[3]":   261.626,
	"B3":      246.942,
	"A#3/Bb3": 233.082,
	"A3":      220.000,
	"G#3/Ab3": 207.652,
	"G3":      195.998,
	"F#3/Gb3": 184.997,
	"F3":      174.614,
	"E3":      164.814,
	"D#3/Eb3": 155.563,
	"D3":      146.832,
	"C#3/Db3": 138.591,
	"C3":      130.813,
	"B2":      123.471,
	"A#2/Bb2": 116.541,
	"A2":      110.000,
	"G#2/Ab2": 103.826,
	"G2":      97.9989,
	"F#2/Gb2": 92.4986,
	"F2":      87.3071,
	"E2":      82.4069,
	"D#2/Eb2": 77.7817,
	"D2":      73.4162,
	"C#2/Db2": 69.2957,
	"C2":      65.4064,
	"B1":      61.7354,
	"A#1/Bb1": 58.2705,
	"A1":      55.0000,
	"G#1/Ab1": 51.9131,
	"G1":      48.9994,
	"F#1/Gb1": 46.2493,
	"F1":      43.6535,
	"E1":      41.2034,
	"D#1/Eb1": 38.8909,
	"D1":      36.7081,
	"C#1/Db1": 34.6478,
	"C1":      32.7032,
	"B0":      30.8677,
	"A#0/Bb0": 29.1352,
	"A0":      27.5000,
	"G#0/Ab0": 25.9565,
	"G0":      24.4997,
	"F#0/Gb0": 23.1247,
	"F0":      21.8268,
	"E0":      20.6017,
	"D#0/Eb0": 19.4454,
	"D0":      18.3540,
	"C#0/Db0": 17.3239,
	"C0":      16.3516,
}

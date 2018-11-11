package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"

	"github.com/faiface/beep/wav"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	speakerSampleRate := beep.SampleRate(48000)
	fmt.Println("Speaker SampleRate:", speakerSampleRate)
	speaker.Init(speakerSampleRate, speakerSampleRate.N(time.Second/10))

	hasAdditionalArg := len(os.Args) > 1
	var hasFileArg bool
	if hasAdditionalArg {
		info, _ := os.Stat(os.Args[1])
		if info != nil && !info.IsDir() {
			hasFileArg = true
		}
	}
	hasAudioArg := hasAdditionalArg && hasFileArg
	audioStreamer := beep.Callback(func(){})
	if hasAudioArg {
		f, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		buffer, format, err := wav.Decode(f)
		if err != nil {
			panic(err)
		}
		left, right := splitToMono(buffer)
		audioStreamer = beep.Seq(
			beep.Callback(func() {
				fmt.Println("splitted stereo mix")
				fmt.Println("Format:", format)
			}),
			beep.Resample(4, format.SampleRate, speakerSampleRate, beep.Mix(
				left,
				right,
			)),
		)
	}
	done := make(chan struct{})
	octavePitch := pitchPerSemitone(-12, 440)
	speaker.Play(
		beep.Seq(
			beep.Callback(func() { fmt.Println("sinewave") }),
			beep.Take(
				speakerSampleRate.N(2*time.Second),
				sineWave(speakerSampleRate, 440),
			),
			beep.Callback(func() { fmt.Println("sinewave octave pitch") }),
			beep.Take(
				speakerSampleRate.N(2*time.Second),
				sineWave(speakerSampleRate, octavePitch),
			),
			beep.Callback(func() { fmt.Println("sawtooth 2 seconds") }),
			beep.Take(
				speakerSampleRate.N(2*time.Second),
				sawtoothWave(speakerSampleRate, 440),
			),
			beep.Callback(func() { fmt.Println("triangle 2 seconds") }),
			beep.Take(
				speakerSampleRate.N(2*time.Second),
				triangleWave(speakerSampleRate, 440),
			),
			beep.Callback(func() { fmt.Println("square 2 seconds") }),
			beep.Take(
				speakerSampleRate.N(2*time.Second),
				squareWave(speakerSampleRate, 440),
			),
			beep.Callback(func() { fmt.Println("pulse 2 seconds") }),
			beep.Take(
				speakerSampleRate.N(2*time.Second),
				pulseWave(speakerSampleRate, 440, 0.3),
			),
			beep.Callback(func() { fmt.Println("LFO 2 seconds") }),
			beep.Take(
				speakerSampleRate.N(2*time.Second),
				lfo(speakerSampleRate, 2, 175, 440),
			),
			beep.Callback(func() { fmt.Println("noise") }),
			beep.Take(
				speakerSampleRate.N(1*time.Second),
				&effects.Volume{
					Streamer: noise,
					Base:     2,
					Volume:   -5,
				},
			),
			audioStreamer,
			beep.Callback(func() {
				<-time.After(4 * time.Second)
				close(done)
			}),
		),
	)

	// pixelgl.Run(run)

	var cl = make(chan os.Signal, 1)
	signal.Notify(cl, os.Interrupt)

	select {
	case <-done:
	case sig := <-cl:
		fmt.Println("\ncaught signal ", sig)
	}
}

func streamErrorPrinter(s beep.Streamer) beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		return beep.Seq(
			s,
			beep.Callback(func() {
				if s.Err() != nil {
					fmt.Printf("Stream error: %v\n", s.Err())
				}
			}),
		).Stream(samples)
	})
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Audio rocks!",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	sr := beep.SampleRate(44100)
	s := streamDrawer(
		beep.Take(
			sr.N(time.Second),
			sineWave(sr, 600),
		),
		imd,
	)
	for {
		tmp := make([][2]float64, 512)
		_, ok := s.Stream(tmp)
		if !ok {
			break
		}
	}
	imd.SetMatrix(pixel.IM.Moved(pixel.V(300, 0)))

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)
		win.Update()
	}
}

func splitToMono(s beep.Streamer) (left beep.Streamer, right beep.Streamer) {
	buffer := &stereoBuffer{
		streamer: s,
	}
	left = beep.StreamerFunc(func(leftSamples [][2]float64) (n int, ok bool) {
		return buffer.StreamLeft(leftSamples)
	})
	right = beep.StreamerFunc(func(rightSamples [][2]float64) (n int, ok bool) {
		return buffer.StreamRight(rightSamples)
	})
	return left, right
}

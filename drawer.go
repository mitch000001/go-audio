package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

func streamDrawer(s beep.Streamer, drawer *imdraw.IMDraw) beep.Streamer {
	x := 0.0
	return beep.StreamerFunc(func(samples [][2]float64) (int, bool) {
		n, ok := s.Stream(samples)
		for _, sample := range samples {
			x++
			drawer.Color = colornames.Limegreen
			drawer.Push(pixel.V(x, 0), pixel.V(x, sample[0]*100))
			drawer.Line(3)
		}
		return n, ok
	})
}

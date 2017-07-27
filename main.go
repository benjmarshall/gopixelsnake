package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	// Setup Window Configuration
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	// Create the window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	imd := imdraw.New(nil)

	imd.Color = colornames.Blueviolet
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(pixel.V(100, 100), pixel.V(700, 100))
	imd.EndShape = imdraw.SharpEndShape
	imd.Push(pixel.V(100, 500), pixel.V(700, 500))
	imd.Line(30)

	imd.Color = colornames.Limegreen
	imd.Push(pixel.V(500, 500))
	imd.Circle(300, 50)
	imd.Color = colornames.Navy
	imd.Push(pixel.V(200, 500), pixel.V(800, 500))
	imd.Ellipse(pixel.V(120, 80), 0)

	imd.Color = colornames.Red
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(pixel.V(500, 350))
	imd.CircleArc(150, -math.Pi, 0, 30)

	// Keep going till the window is closed
	for !win.Closed() {

		// Clear the frame
		win.Clear(colornames.Aliceblue)

		imd.Draw(win)

		// Update the window
		win.Update()
	}
}

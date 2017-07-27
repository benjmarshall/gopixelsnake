package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
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
		VSync:  false,
	}

	// Create the window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	fps := time.Tick(time.Second / 120)

	face, err := loadTTF("intuitive.ttf", 80)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(win.Bounds().Center(), atlas)

	txt.Color = colornames.Lightgray

	// Keep going till the window is closed
	for !win.Closed() {

		txt.WriteString(win.Typed())
		if win.JustPressed(pixelgl.KeyEnter) || win.Repeated(pixelgl.KeyEnter) {
			txt.WriteRune('\n')
		}
		if win.JustPressed(pixelgl.KeyTab) || win.Repeated(pixelgl.KeyTab) {
			txt.WriteRune('\t')
		}

		// Clear the frame
		win.Clear(colornames.Darkcyan)

		txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

		// Update the window
		win.Update()

		// Synchonise the framerate
		<-fps
	}
}

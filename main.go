package main

import (
	"github.com/benjmarshall/gosnake/snake"
	"github.com/benjmarshall/gosnake/types"
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

	// Setup Game Configuration
	gameCFG := types.NewGameCFG(700, 700, 5, cfg)

	// Create Game Area Shape
	imd := imdraw.New(nil)
	imd.Color = colornames.White
	imd.Push(gameCFG.GetGameAreaAsVecs())
	imd.Rectangle(2)

	// Initialize a new snake
	snake := snake.NewSnake(gameCFG)

	// Create a snake shape
	imd.Color = colornames.Purple
	imd.Push(snake.GetHeadPos(), snake.GetTailPos())
	imd.Line(gameCFG.GetGridSize())

	// Keep going till the window is closed
	for !win.Closed() {

		// Clear the frame
		win.Clear(colornames.Darkcyan)

		// Draw the Game Area
		imd.Draw(win)

		// Update the window
		win.Update()
	}
}

package main

import (
	"fmt"
	"log"
	"time"

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
	gameCFG := types.NewGameCFG(700, 700, 10, cfg)

	// Initialize a new snake
	s := snake.NewSnake(gameCFG)
	log.Println(s)

	// Create the Game Background Shape
	imdArea := imdraw.New(nil)

	// Create the Game Contents Shape
	imdGame := imdraw.New(nil)

	// Create some variables
	var (
		dir        snake.Direction
		frames     = 0
		second     = time.Tick(time.Second)
		keyPressed = false
	)

	// Draw the initial frames	// Clear the frame
	win.Clear(colornames.Darkcyan)
	drawGameBackground(win, imdArea, &gameCFG)
	drawSnake(win, imdGame, &gameCFG, &s)
	win.Update()

	// Start the snake timer
	snakeTicker := time.NewTicker(time.Second * time.Duration(s.GetSpeed()))

	// Keep going till the window is closed
	for !win.Closed() {

		// Clear the screen
		win.Clear(colornames.Darkcyan)

		// Catch user input
		if keyPressed == false {
			if win.JustPressed(pixelgl.KeyUp) {
				dir = snake.UP
				keyPressed = true
			} else if win.JustPressed(pixelgl.KeyDown) {
				dir = snake.DOWN
				keyPressed = true
			} else if win.JustPressed(pixelgl.KeyLeft) {
				dir = snake.LEFT
				keyPressed = true
			} else if win.JustPressed(pixelgl.KeyRight) {
				dir = snake.RIGHT
				keyPressed = true
			}
		}

		// Update the snake
		select {
		case <-snakeTicker.C:
			// Update the snake
			s.Update(false, dir)
			// Reset the user inputs
			dir = snake.NOCHANGE
			keyPressed = false
		default:
		}

		// Draw the sframe
		drawGameBackground(win, imdArea, &gameCFG)
		drawSnake(win, imdGame, &gameCFG, &s)
		win.Update()
		frames++

		// Update FPS
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}

	}
}

func drawGameBackground(win *pixelgl.Window, imd *imdraw.IMDraw, gameCFG *types.GameCFGType) {
	imd.Clear()
	imd.Color = colornames.White
	imd.Push(gameCFG.GetGameAreaAsVecs())
	imd.Rectangle(2)
	imd.Draw(win)
}

func drawSnake(win *pixelgl.Window, imd *imdraw.IMDraw, gameCFG *types.GameCFGType, s *snake.Type) {
	imd.Clear()
	imd.Color = colornames.Purple
	positions := []pixel.Vec{s.GetHeadPos()}
	positions = append(positions, s.GetPositionPoints()...)
	positions = append(positions, s.GetTailPos())
	imd.Push(positions...)
	imd.Line(gameCFG.GetGridSize())
	imd.Draw(win)
}

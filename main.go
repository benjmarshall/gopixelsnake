package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/benjmarshall/gosnake/snake"
	"github.com/benjmarshall/gosnake/types"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
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

	// Creat a text Atlas
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// Create the window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Setup Game Configuration
	gameCFG := types.NewGameCFG(700, 700, 2, 10, cfg)

	// Initialize a new snake
	s := snake.NewSnake(gameCFG)

	// Generate a berry
	berry := generateRandomBerry(&gameCFG)

	// Create the Game Background Shape
	imdArea := imdraw.New(nil)

	// Create the Game Contents Shape
	imdGame := imdraw.New(nil)

	// Create a berry Contents Shape
	imdBerry := imdraw.New(nil)

	// Create some variables
	var (
		dir         snake.Direction
		frames      = 0
		second      = time.Tick(time.Second)
		keyPressed  = false
		gameRunning = true
		gameOver    = false
		eaten       = false
	)

	// Draw the initial frames	// Clear the frame
	win.Clear(colornames.Darkcyan)
	drawGameBackground(win, imdArea, &gameCFG)
	drawSnakeRect(win, imdGame, &gameCFG, &s)
	drawBerry(win, imdBerry, &gameCFG, berry)
	win.Update()

	// Start the snake timer
	snakeTicker := time.NewTicker(time.Second * time.Duration(s.GetSpeed()))

	// Keep going till the window is closed
	for !win.Closed() {

		// Clear the screen
		win.Clear(colornames.Darkcyan)

		// Do game logic only if the game is actually running!
		if gameRunning {

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
				s.Update(eaten, dir)
				// Reset the user inputs
				dir = snake.NOCHANGE
				keyPressed = false
				// Debug
				// log.Println(s.GetHeadPos())
				// log.Println(s.GetPositionPoints())
				// log.Println(s.GetTailPos())
				// Check the snake is still in bounds
				if !s.CheckSnakeOK(&gameCFG) {
					gameOver = true
					gameRunning = false
				}
				// Check if the snake has eaten
				eaten = s.CheckIfSnakeHasEaten(&gameCFG, berry)
				if eaten {
					berry = generateRandomBerry(&gameCFG)
				}
			default:
			}

		}

		// Always draw the game
		drawGameBackground(win, imdArea, &gameCFG)
		drawSnakeRect(win, imdGame, &gameCFG, &s)
		drawBerry(win, imdBerry, &gameCFG, berry)

		// Check if the game is over
		if gameOver {
			drawGameOver(win, basicAtlas, &gameCFG)
			if win.JustPressed(pixelgl.KeyEnter) {
				win.SetClosed(true)
			}
		}

		// Always update the window
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
	min, max := gameCFG.GetGameAreaAsVecs()
	min = gameCFG.GetWindowMatrix().Project(min)
	max = gameCFG.GetWindowMatrix().Project(max)
	vec := pixel.V(gameCFG.GetBorderWeight()/2, gameCFG.GetBorderWeight()/2)
	min = min.Sub(vec)
	max = max.Add(vec)
	imd.Push(min, max)
	imd.Rectangle(2)
	imd.Draw(win)
}

func drawSnakeRect(win *pixelgl.Window, imd *imdraw.IMDraw, gameCFG *types.GameCFGType, s *snake.Type) {
	imd.Clear()
	imd.Color = colornames.Purple
	positions := []pixel.Vec{s.GetHeadPos()}
	positions = append(positions, s.GetPositionPoints()...)
	positions = append(positions, s.GetTailPos())
	for _, pos := range positions {
		m := gameCFG.GetWindowMatrix()
		vec := pixel.V(gameCFG.GetGridSize()/2, gameCFG.GetGridSize()/2)
		min := m.Project(pos).Sub(vec)
		max := m.Project(pos).Add(vec)
		imd.Push(min, max)
	}
	imd.Rectangle(0)
	imd.Draw(win)
}

func drawBerry(win *pixelgl.Window, imd *imdraw.IMDraw, gameCFG *types.GameCFGType, berry pixel.Vec) {
	berry = gameCFG.GetWindowMatrix().Project(berry)
	imd.Clear()
	imd.Color = colornames.Orangered
	imd.Push(berry)
	imd.Circle(gameCFG.GetGridSize()/2, 0)
	imd.Draw(win)
}

func drawGameOver(win *pixelgl.Window, atlas *text.Atlas, gameCFG *types.GameCFGType) {
	textOrig := gameCFG.GetWindowMatrix().Project(gameCFG.GetGameAreaAsRec().Center())
	gameoverMessage := text.New(textOrig, atlas)
	lines := []string{
		"Game Over!",
		"Press Enter to exit...",
	}
	gameoverMessage.Color = colornames.Black
	for _, line := range lines {
		gameoverMessage.Dot.X -= gameoverMessage.BoundsOf(line).W() / 2
		fmt.Fprintln(gameoverMessage, line)
	}
	gameoverMessage.Orig.Add(pixel.V(0, gameoverMessage.BoundsOf(lines[0]).H()))
	gameoverMessage.Draw(win, pixel.IM.Scaled(gameoverMessage.Orig, 4))
}

func generateRandomBerry(gameCFG *types.GameCFGType) pixel.Vec {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x, y := gameCFG.GetGameAreaDims()
	berryX := float64(r.Intn(int(x/gameCFG.GetGridSize()) - 1))
	berryY := float64(r.Intn(int(y/gameCFG.GetGridSize()) - 1))
	berry := pixel.V(berryX, berryY)
	log.Printf("Berry: %v", berry)
	return gameCFG.GetGridMatrix().Project(berry)
}

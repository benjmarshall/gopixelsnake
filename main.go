package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/benjmarshall/gosnake/gametext"
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
	gameCFG := types.NewGameCFG(700, 700, 2, 10, cfg)

	// Setup text structure
	textStruct := gametext.NewGameText(win, gameCFG)

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
		frames         = 0
		second         = time.Tick(time.Second)
		gameRunning    = false
		gameOver       = false
		eaten          = false
		inputKeyBuffer = []snake.Direction{}
		dir            snake.Direction
		score          = 0
	)

	// Draw the initial frames	// Clear the frame
	win.Clear(colornames.Darkcyan)
	drawGameBackground(win, imdArea, &gameCFG)
	drawSnakeRect(win, imdGame, &gameCFG, &s)
	drawBerry(win, imdBerry, &gameCFG, berry)
	textStruct.DrawTitleText(win)
	textStruct.DrawScoreText(win, score)
	textStruct.DrawControlsText(win)
	win.Update()

	// Keep going till the window is closed
	for !win.Closed() {

		// Clear the screen
		win.Clear(colornames.Darkcyan)

		if !gameRunning && !gameOver {
			// Game is not running so wait for user to do something!
			// Catch user input
			startDir := snake.NOCHANGE
			if win.JustPressed(pixelgl.KeyUp) {
				startDir = s.StartOfGame(snake.UP)
				gameRunning = true
			} else if win.JustPressed(pixelgl.KeyDown) {
				startDir = s.StartOfGame(snake.DOWN)
				gameRunning = true
			} else if win.JustPressed(pixelgl.KeyLeft) {
				startDir = s.StartOfGame(snake.LEFT)
				gameRunning = true
			} else if win.JustPressed(pixelgl.KeyRight) {
				startDir = s.StartOfGame(snake.RIGHT)
				gameRunning = true
			} else if win.JustPressed(pixelgl.KeyX) {
				win.SetClosed(true)
			}
			if startDir != snake.NOCHANGE {
				inputKeyBuffer = append(inputKeyBuffer, startDir)
			}
		}

		// Do game logic only if the game is actually running!
		if gameRunning {

			// Catch user input
			if win.JustPressed(pixelgl.KeyUp) {
				inputKeyBuffer = append(inputKeyBuffer, snake.UP)
			} else if win.JustPressed(pixelgl.KeyDown) {
				inputKeyBuffer = append(inputKeyBuffer, snake.DOWN)
			} else if win.JustPressed(pixelgl.KeyLeft) {
				inputKeyBuffer = append(inputKeyBuffer, snake.LEFT)
			} else if win.JustPressed(pixelgl.KeyRight) {
				inputKeyBuffer = append(inputKeyBuffer, snake.RIGHT)
			}

			// Update the snake
			select {
			case <-s.GetTicker():
				// Update the snake
				if len(inputKeyBuffer) == 0 {
					dir = snake.NOCHANGE
				} else {
					dir = inputKeyBuffer[0]
					inputKeyBuffer = inputKeyBuffer[1:len(inputKeyBuffer)]
				}
				s.Update(eaten, dir)
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
					s.IncreaseSpeed()
				}
				// Update the score
				score += int((s.GetSpeed() * 10))
				if eaten {
					score += int((1000 * s.GetSpeed()))
				}
			default:
			}

		} else if gameOver {
			// Game has ended wait for user to continue
			if win.JustPressed(pixelgl.KeyEnter) {
				// reset the board
				gameOver = false
				score = 0
				berry = generateRandomBerry(&gameCFG)
				s = snake.NewSnake(gameCFG)
			}
		}

		// Always draw the game
		drawGameBackground(win, imdArea, &gameCFG)
		drawSnakeRect(win, imdGame, &gameCFG, &s)
		drawBerry(win, imdBerry, &gameCFG, berry)
		textStruct.DrawTitleText(win)
		textStruct.DrawScoreText(win, score)
		textStruct.DrawControlsText(win)
		if !gameRunning && !gameOver {
			// Show the start game message
			textStruct.DrawStartGameText(win)
		} else if gameOver {
			textStruct.DrawGameOverText(win)
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

func generateRandomBerry(gameCFG *types.GameCFGType) pixel.Vec {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x, y := gameCFG.GetGameAreaDims()
	berryX := float64(r.Intn(int(x/gameCFG.GetGridSize()) - 1))
	berryY := float64(r.Intn(int(y/gameCFG.GetGridSize()) - 1))
	berry := pixel.V(berryX, berryY)
	log.Printf("Berry: %v", berry)
	return gameCFG.GetGridMatrix().Project(berry)
}

package main

import (
	"fmt"
	"time"

	"github.com/benjmarshall/gopixelsnake/drawing"
	"github.com/benjmarshall/gopixelsnake/game"
	"github.com/benjmarshall/gopixelsnake/gametext"
	"github.com/benjmarshall/gopixelsnake/scores"
	"github.com/benjmarshall/gopixelsnake/snake"
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
		Title:     "Pixel Rocks!",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: false,
		VSync:     true,
	}

	// Create the window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Setup Game Configuration
	gameCFG := game.NewGameConfig(700, 700, 2, 10, cfg)

	// Setup text structure
	textStruct := gametext.NewGameText(win, gameCFG)

	// Setup a scores structure
	scoresTable := scores.NewScores("high_scores.csv", 10)

	// Initialize a new snake
	s := snake.NewSnake(gameCFG)

	// Generate a berry
	berry := game.GenerateRandomBerry(&gameCFG)

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
		showScores     = false
		scoreName      string
		highScore      = false
	)

	// Draw the initial frames
	win.Clear(colornames.Darkcyan)
	drawing.DrawGameBackground(win, imdArea, &gameCFG)
	drawing.DrawSnakeRect(win, imdGame, &gameCFG, &s)
	drawing.DrawBerry(win, imdBerry, &gameCFG, berry)
	textStruct.DrawTitleText(win)
	textStruct.DrawScoreText(win, score)
	textStruct.DrawControlsText(win)
	win.Update()

	// Keep going till the window is closed
	for !win.Closed() {

		// Clear the screen
		win.Clear(colornames.Darkcyan)

		if !gameRunning && !gameOver && !showScores {
			// Game is not running so wait for user to do something!
			if win.JustPressed(pixelgl.KeyUp) {
				s.StartOfGame(snake.UP)
				gameRunning = true
			} else if win.JustPressed(pixelgl.KeyDown) {
				s.StartOfGame(snake.DOWN)
				gameRunning = true
			} else if win.JustPressed(pixelgl.KeyLeft) {
				s.StartOfGame(snake.LEFT)
				gameRunning = true
			} else if win.JustPressed(pixelgl.KeyRight) {
				s.StartOfGame(snake.RIGHT)
				gameRunning = true
			} else if win.JustPressed(pixelgl.KeyX) {
				win.SetClosed(true)
			}
			if win.JustPressed(pixelgl.KeyS) {
				showScores = true
			}
		} else if !gameRunning && !gameOver && showScores {
			if win.JustPressed(pixelgl.KeyS) {
				showScores = false
			} else if win.JustPressed(pixelgl.KeyX) {
				win.SetClosed(true)
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
				// Check the snake is still in bounds
				if !s.CheckSnakeOK(&gameCFG) {
					gameOver = true
					gameRunning = false
					if score >= scoresTable.GetBottomScore() {
						highScore = true
					}
					break
				}
				// Check if the snake has eaten
				eaten = s.CheckIfSnakeHasEaten(&gameCFG, berry)
				if eaten {
					berry = game.GenerateRandomBerry(&gameCFG)
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
			// Game has ended, wait for user to continue
			if win.JustPressed(pixelgl.KeyEnter) {
				// Submit score and reset for a new game
				if highScore {
					scoresTable.AddScore(score, scoreName)
				}
				// reset the board
				scoreName = ""
				gameOver = false
				highScore = false
				score = 0
				berry = game.GenerateRandomBerry(&gameCFG)
				s = snake.NewSnake(gameCFG)
			} else if win.JustPressed(pixelgl.KeyBackspace) {
				// Add support for deleting charaters from score name
				scoreName = scoreName[0 : len(scoreName)-1]
			} else if len(scoreName) < 3 {
				// Capture input for score name (up to 3 chars)
				scoreName = scoreName + win.Typed()
			}
		}

		// Always draw the game
		drawing.DrawGameBackground(win, imdArea, &gameCFG)
		if !showScores {
			// Hide game elements if high scores are being diplayed
			drawing.DrawSnakeRect(win, imdGame, &gameCFG, &s)
			drawing.DrawBerry(win, imdBerry, &gameCFG, berry)
		}
		textStruct.DrawTitleText(win)
		textStruct.DrawScoreText(win, score)
		textStruct.DrawControlsText(win)
		if !gameRunning && !gameOver && !showScores {
			// Show the start game message
			textStruct.DrawStartGameText(win)
		} else if gameOver {
			textStruct.DrawGameOverText(win, &gameCFG, scoreName, highScore)
		} else if showScores {
			textStruct.DrawScoresListText(win, &gameCFG, &scoresTable)
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

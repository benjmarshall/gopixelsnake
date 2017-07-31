package gametext

import (
	"fmt"
	"strconv"

	"github.com/benjmarshall/gosnake/types"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// Type holds the various text object for the game
type Type struct {
	title    snaketext
	score    snaketext
	controls snaketext
	gameover snaketext
	atlas    *text.Atlas
}

// snaketext is a wrapper around pixel.text which also holds scale information for drawing
type snaketext struct {
	text      *text.Text
	drawScale pixel.Matrix
}

// NewGameText generates a new game text structure used to control all text display for the game
func NewGameText(win *pixelgl.Window, gameCFG types.GameCFGType) Type {
	t := new(Type)
	// Create a text Atlas
	t.atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// Create Game title
	textColumnWidth := win.Bounds().W() - gameCFG.GetWindowMatrix().Project(gameCFG.GetGameAreaAsRec().Max).X
	textOrigX := win.Bounds().W() - (textColumnWidth / 2)
	textOrigY := win.Bounds().H() * 0.9
	textOrig := pixel.V(textOrigX, textOrigY)
	t.title.text = text.New(textOrig, t.atlas)
	t.title.text.Color = colornames.Black
	lines := []string{
		"Go Pixel",
		"Snake",
	}
	for _, line := range lines {
		t.title.text.Dot.X -= t.title.text.BoundsOf(line).W() / 2
		fmt.Fprintln(t.title.text, line)
	}
	t.title.drawScale = pixel.IM.Scaled(t.title.text.Orig, 4)

	// Create Game Over Text
	textOrig = gameCFG.GetWindowMatrix().Project(gameCFG.GetGameAreaAsRec().Center())
	t.gameover.text = text.New(textOrig, t.atlas)
	lines = []string{
		"Game Over!",
		"Press Enter to exit...",
	}
	t.gameover.text.Color = colornames.Black
	for _, line := range lines {
		t.gameover.text.Dot.X -= t.gameover.text.BoundsOf(line).W() / 2
		fmt.Fprintln(t.gameover.text, line)
	}
	t.gameover.text.Orig.Add(pixel.V(0, t.gameover.text.BoundsOf(lines[0]).H()))
	t.gameover.drawScale = pixel.IM.Scaled(t.gameover.text.Orig, 4)

	// Create Score Text
	textOrigY = win.Bounds().H() * 0.8
	textOrig = pixel.V(textOrigX, textOrigY)
	t.score.text = text.New(textOrig, t.atlas)
	t.score.text.Color = colornames.Black
	scoreText := "0"
	t.score.text.Dot.X = t.score.text.Orig.X - t.score.text.BoundsOf(scoreText).W()/2
	fmt.Fprintln(t.score.text, scoreText)
	t.score.drawScale = pixel.IM.Scaled(t.score.text.Orig, 6)

	// Create Controls Text
	textOrigY = win.Bounds().H() * 0.6
	textOrig = pixel.V(textOrigX, textOrigY)
	t.controls.text = text.New(textOrig, t.atlas)
	lines = []string{
		"Start Game",
		"Enter\n",
		"Control Snake",
		"Arrows\n",
		"View Scores",
		"S\n",
		"Exit",
		"X",
	}
	t.controls.text.Color = colornames.Black
	for _, line := range lines {
		t.controls.text.Dot.X -= t.controls.text.BoundsOf(line).W() / 2
		fmt.Fprintln(t.controls.text, line)
	}
	t.controls.text.LineHeight = 1.5
	t.controls.drawScale = pixel.IM.Scaled(t.controls.text.Orig, 3)

	return *t

}

// DrawTitleText draws the title text on the window provided
func (t *Type) DrawTitleText(win *pixelgl.Window) {
	t.title.text.Draw(win, t.title.drawScale)
}

// DrawGameOverText draws the game over text on the provided window
func (t *Type) DrawGameOverText(win *pixelgl.Window) {
	t.gameover.text.Draw(win, t.gameover.drawScale)
}

// DrawControlsText draws the controls text on the provided window
func (t *Type) DrawControlsText(win *pixelgl.Window) {
	t.controls.text.Draw(win, t.controls.drawScale)
}

// DrawScoreText draws the text on the provide window
func (t *Type) DrawScoreText(win *pixelgl.Window, score int) {
	scoreText := strconv.Itoa(score)
	t.score.text.Clear()
	t.score.text.Dot.X = t.score.text.Orig.X - t.score.text.BoundsOf(scoreText).W()/2
	fmt.Fprintf(t.score.text, scoreText)
	t.score.text.Draw(win, t.score.drawScale)
}

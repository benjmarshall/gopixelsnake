package drawing

import (
	"github.com/benjmarshall/gosnake/game"
	"github.com/benjmarshall/gosnake/snake"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// DrawGameBackground draws the game area border
func DrawGameBackground(win *pixelgl.Window, imd *imdraw.IMDraw, gameCFG *game.Config) {
	imd.Clear()
	imd.Color = colornames.White
	min, max := gameCFG.GetGameAreaAsVecs()
	min = gameCFG.GetWindowMatrix().Project(min)
	max = gameCFG.GetWindowMatrix().Project(max)
	vec := pixel.V(gameCFG.GetBorderWeight()/2, gameCFG.GetBorderWeight()/2)
	min = min.Sub(vec)
	max = max.Add(vec)
	imd.Push(min, max)
	imd.Rectangle(gameCFG.GetBorderWeight())
	imd.Draw(win)
}

// DrawSnakeRect draws the snake shape using rectangles
func DrawSnakeRect(win *pixelgl.Window, imd *imdraw.IMDraw, gameCFG *game.Config, s *snake.Type) {
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

// DrawBerry draws the berry shape
func DrawBerry(win *pixelgl.Window, imd *imdraw.IMDraw, gameCFG *game.Config, berry pixel.Vec) {
	berry = gameCFG.GetWindowMatrix().Project(berry)
	imd.Clear()
	imd.Color = colornames.Orangered
	imd.Push(berry)
	imd.Circle(gameCFG.GetGridSize()/2, 0)
	imd.Draw(win)
}

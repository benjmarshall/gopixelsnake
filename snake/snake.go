package snake

import (
	"math/rand"
	"time"

	"github.com/benjmarshall/gosnake/types"
	"github.com/faiface/pixel"
)

// Type is a struct which represents a snake in the game
type Type struct {
	headPos          pixel.Vec
	tailPos          pixel.Vec
	length           float64
	speed            float64
	currentDirection direction
	pointsList       []pixel.Vec
	gameCFG          *types.GameCFGType
}

type direction struct {
	val pixel.Vec
}

var (
	up    = direction{pixel.V(0, 1)}
	down  = direction{pixel.V(0, -1)}
	left  = direction{pixel.V(-1, 0)}
	right = direction{pixel.V(1, 0)}
)

// NewSnake returns an initialised snake
func NewSnake(gameCFG types.GameCFGType) Type {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	snake := new(Type)
	snake.gameCFG = &gameCFG
	snake.length = 3
	snake.speed = 1.0
	x, y := gameCFG.GetGameAreaDims()
	m := gameCFG.GetGridMatrix()
	snake.headPos = m.Project(pixel.V(float64(r.Intn(int(x))), float64(r.Intn(int(y)))))
	switch i := r.Intn(3); {
	case i == 0:
		snake.currentDirection = up
	case i == 1:
		snake.currentDirection = down
	case i == 2:
		snake.currentDirection = left
	case i == 3:
		snake.currentDirection = right
	default:
		snake.currentDirection = up
	}
	snake.tailPos = snake.headPos.Add(snake.dir().Scaled(snake.length))
	return *snake
}

func (s *Type) dir() pixel.Vec {
	return s.currentDirection.val
}

// GetHeadPos returns the position of the head of the snake in the game area coordinate plane
func (s *Type) GetHeadPos() pixel.Vec {
	return s.gameCFG.GetGridMatrix().Unproject(s.headPos)
}

// GetTailPos returns the position of the tail of the snake in the game area coordinate plane
func (s *Type) GetTailPos() pixel.Vec {
	return s.gameCFG.GetGridMatrix().Unproject(s.tailPos)
}

package snake

import (
	"log"
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
	currentDirection Direction
	pointsList       []pixel.Vec
	gameCFG          *types.GameCFGType
}

// Direction is used to define the direction the snake is heading
type Direction struct {
	val pixel.Vec
}

var (
	// UP is the Direction defining travel towards the top of the game area.
	UP = Direction{pixel.V(0, 1)}
	// DOWN is the Direction defining travel towards the bottom of the game area.
	DOWN = Direction{pixel.V(0, -1)}
	// LEFT is the Direction defining travel towards the left of the game area.
	LEFT = Direction{pixel.V(-1, 0)}
	// RIGHT is the Direction defining travel towards the right of the game area.
	RIGHT = Direction{pixel.V(1, 0)}
	// NOCHANGE is a blank Direction, it can be used to not alter the current heading.
	NOCHANGE = Direction{pixel.V(0, 0)}
)

// NewSnake returns an initialised snake
func NewSnake(gameCFG types.GameCFGType) Type {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	snake := new(Type)
	snake.gameCFG = &gameCFG
	snake.length = 5
	snake.speed = 1
	x, y := gameCFG.GetGameAreaDims()
	snake.headPos = pixel.V(float64(r.Intn(int(x/gameCFG.GetGridSize())-10))+5, float64(r.Intn(int(y/gameCFG.GetGridSize())-10))+5)
	switch i := r.Intn(3); {
	case i == 0:
		snake.currentDirection = UP
	case i == 1:
		snake.currentDirection = DOWN
	case i == 2:
		snake.currentDirection = LEFT
	case i == 3:
		snake.currentDirection = RIGHT
	default:
		snake.currentDirection = UP
	}
	snake.tailPos = snake.headPos.Sub(snake.dir().Scaled(snake.length - 1))
	return *snake
}

func (s *Type) dir() pixel.Vec {
	return s.currentDirection.val
}

// GetHeadPos returns the position of the head of the snake in the game area coordinate plane
func (s *Type) GetHeadPos() pixel.Vec {
	return s.gameCFG.GetGridMatrix().Project(s.headPos)
}

// GetTailPos returns the position of the tail of the snake in the game area coordinate plane
func (s *Type) GetTailPos() pixel.Vec {
	return s.gameCFG.GetGridMatrix().Project(s.tailPos)
}

// GetPositionPoints returns the list of the snakes previous turn positions in the game area coordinate plane
func (s *Type) GetPositionPoints() []pixel.Vec {
	positions := []pixel.Vec{}
	for _, pos := range s.pointsList {
		positions = append(positions, s.gameCFG.GetGridMatrix().Project(pos))
	}
	return positions
}

// GetSpeed returns the snake speed multiplier
func (s *Type) GetSpeed() float64 {
	return s.speed
}

// Update is used to UPdate the status of snake position and speed.
func (s *Type) Update(eaten bool, dir Direction) {
	// If the snake has eaten let's up the speed
	if eaten {
		s.speed *= 1.1
	}

	if dir != NOCHANGE {
		//log.Println("Changing direction")
		// Update the direction
		s.currentDirection = dir
		// Push the current head position into the points stack
		s.pointsList = append([]pixel.Vec{s.headPos}, s.pointsList...)
	}
	if len(s.pointsList) > 0 {
		if s.tailPos == s.pointsList[len(s.pointsList)-1] {
			//log.Println("Checking tail pos")
			// If the tail is on our last point the remove it from the current stack
			if len(s.pointsList) <= 1 {
				//log.Println("only 1")
				s.pointsList = []pixel.Vec{}
			} else {
				//log.Println("mod list")
				s.pointsList = s.pointsList[0 : len(s.pointsList)-1]
			}
		}
	}

	// Update the head position
	s.headPos = s.headPos.Add(s.currentDirection.val)
	// Update the tail position
	if len(s.pointsList) == 0 {
		s.tailPos = s.tailPos.Add(s.currentDirection.val)
	} else {
		//log.Println("Moving tail towards point")
		vec := s.tailPos.To(s.pointsList[len(s.pointsList)-1]).Unit()
		s.tailPos = s.tailPos.Add(vec)
		log.Println(vec)
	}
	// log.Println("Snake after update:")
	// log.Println(s.headPos)
	// log.Println(s.pointsList)
	// log.Println(s.tailPos)
}

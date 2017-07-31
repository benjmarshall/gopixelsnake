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
	ticker           time.Ticker
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
	snake.speed = 2
	x, y := gameCFG.GetGameAreaDims()
	snakeStartingMargin := 10
	startingDimX := int(x/gameCFG.GetGridSize()) - snakeStartingMargin
	startingDimY := int(y/gameCFG.GetGridSize()) - snakeStartingMargin
	startX := r.Intn(startingDimX) + (snakeStartingMargin / 2)
	startY := r.Intn(startingDimY) + (snakeStartingMargin / 2)
	snake.headPos = pixel.V(float64(startX), float64(startY))
	snake.ticker = *time.NewTicker(time.Second / time.Duration(snake.speed))
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

// GetTicker returns the snake speed ticker
func (s *Type) GetTicker() time.Ticker {
	return s.ticker
}

// IncreaseSpeed increase the speed of the snake
func (s *Type) IncreaseSpeed() {
	s.speed++
	s.ticker.Stop()
	s.ticker = *time.NewTicker(time.Second / time.Duration(s.speed))
}

// Update is used to Update the status of snake position and speed.
func (s *Type) Update(eaten bool, dir Direction) {
	// If the snake has eaten let's  the length
	if eaten {
		s.length++
	}

	if dir != NOCHANGE {
		//log.Println("Changing direction")
		// Ignore a request to change to the opposite direction
		if !(dir == UP && s.currentDirection == DOWN) &&
			!(dir == DOWN && s.currentDirection == UP) &&
			!(dir == LEFT && s.currentDirection == RIGHT) &&
			!(dir == RIGHT && s.currentDirection == LEFT) {
			// Update the direction
			s.currentDirection = dir
			// Push the current head position into the points stack
			s.pointsList = append([]pixel.Vec{s.headPos}, s.pointsList...)
		}
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
	// Update the tail position (if we have eaten a berry, leave the tail where it is)
	if !eaten {
		if len(s.pointsList) == 0 {
			s.tailPos = s.tailPos.Add(s.currentDirection.val)
		} else {
			//log.Println("Moving tail towards point")
			vec := s.tailPos.To(s.pointsList[len(s.pointsList)-1]).Unit()
			s.tailPos = s.tailPos.Add(vec)
			log.Println(vec)
		}
	}
	// log.Println("Snake after update:")
	// log.Println(s.headPos)
	// log.Println(s.pointsList)
	// log.Println(s.tailPos)
}

// CheckSnakeOK is used to check the snake hasn't exicted the game area and has not hit itself
func (s *Type) CheckSnakeOK(gameCFG *types.GameCFGType) bool {

	// Check snake is inside the game boundary
	if !gameCFG.GetGameAreaAsRec().Contains(s.GetHeadPos()) {
		log.Println("Game Over")
		log.Printf("Snake Head: %v", s.GetHeadPos())
		log.Printf("Game Area: %v", gameCFG.GetGameAreaAsRec())
		return false
	}

	// Check snake hasn't hit itself
	// First collect all the turn positions of the snake (minus the head)
	positions := []pixel.Vec{}
	positions = append(positions, s.pointsList...)
	positions = append(positions, s.tailPos)
	// log.Println("Checking snake hit")
	// log.Printf("Head Pos: %v", s.headPos)
	// log.Printf("Tail Pos: %v", s.tailPos)
	// Loop over the positions
	for i := 0; i < len(positions)-1; i++ {
		// Get the length of the line from one turn position to the next
		l := positions[i].To(positions[i+1]).Len() + 1
		// Now add all of the points along the line to a new slice
		subPositions := []pixel.Vec{positions[i]}
		if l > 2 {
			for j := 1.0; j < l-2+1; j++ {
				// We use linear interpolation here to get the points inbetween the end vectors
				interpPoint := j / (l - 1)
				interpVal := pixel.Lerp(positions[i], positions[i+1], interpPoint)
				subPositions = append(subPositions, interpVal)
				// log.Printf("Adding subpixel %f, interpolated using interp point %v, from length of %v", interpVal, interpPoint, l)
			}
		}
		subPositions = append(subPositions, positions[i+1])
		// Now loop over the slice of points along the line and see if we have a colision.
		// log.Printf("Checking pixels: %v", subPositions)
		for _, subPos := range subPositions {
			if s.headPos == subPos {
				return false
			}
		}
	}

	return true
}

// CheckIfSnakeHasEaten is used to check the snake has easten the berry
func (s *Type) CheckIfSnakeHasEaten(gameCFG *types.GameCFGType, berry pixel.Vec) bool {
	berryTransformed := gameCFG.GetGridMatrix().Unproject(berry)
	if s.headPos == berryTransformed {
		return true
	}
	return false
}

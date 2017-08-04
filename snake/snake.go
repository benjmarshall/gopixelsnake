package snake

import (
	"math/rand"
	"time"

	"github.com/benjmarshall/gopixelsnake/game"
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
	gameCFG          *game.Config
	ticker           time.Ticker
	tickerChannel    chan time.Time
	startChannel     chan time.Time
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
func NewSnake(gameCFG game.Config) Type {
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
	snake.tailPos = snake.headPos.Sub(snake.currentDirection.val.Scaled(snake.length - 1))
	// Debug
	// log.Println("__Snake Config__")
	// log.Printf("Direction: %v", snake.currentDirection)
	return *snake
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
func (s *Type) GetTicker() <-chan time.Time {
	return s.tickerChannel
}

// IncreaseSpeed increase the speed of the snake
func (s *Type) IncreaseSpeed() {
	s.speed++
	// Shut down the old ticker and channel multiplex
	close(s.startChannel)
	s.ticker.Stop()
	// Start up the new ones
	s.startChannel = make(chan time.Time)
	s.ticker = *time.NewTicker(time.Second / time.Duration(s.speed))
	go tickerMultiplex(s.tickerChannel, s.ticker.C, s.startChannel)
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
			// If the tail is on our last point the remove it from the current stack
			if len(s.pointsList) <= 1 {
				s.pointsList = []pixel.Vec{}
			} else {
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
			vec := s.tailPos.To(s.pointsList[len(s.pointsList)-1]).Unit()
			s.tailPos = s.tailPos.Add(vec)
		}
	}
}

// CheckSnakeOK is used to check the snake hasn't exicted the game area and has not hit itself
func (s *Type) CheckSnakeOK(gameCFG *game.Config) bool {

	// Check snake is inside the game boundary
	if !gameCFG.GetGameAreaAsRec().Contains(s.GetHeadPos()) {
		return false
	}

	// Check snake hasn't hit itself
	// First collect all the turn positions of the snake (minus the head)
	positions := []pixel.Vec{}
	positions = append(positions, s.pointsList...)
	positions = append(positions, s.tailPos)
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
			}
		}
		subPositions = append(subPositions, positions[i+1])
		// Now loop over the slice of points along the line and see if we have a colision.
		for _, subPos := range subPositions {
			if s.headPos == subPos {
				return false
			}
		}
	}

	return true
}

// CheckIfSnakeHasEaten is used to check the snake has easten the berry
func (s *Type) CheckIfSnakeHasEaten(gameCFG *game.Config, berry pixel.Vec) bool {
	berryTransformed := gameCFG.GetGridMatrix().Unproject(berry)
	if s.headPos == berryTransformed {
		return true
	}
	return false
}

// StartOfGame is used to allow the starting of the game with the arrow keys to choose
// the initial direction of the snake.
func (s *Type) StartOfGame(dir Direction) {
	// Start the snake ticker now so it is synced with the users key press
	// We use a goroutine running a channel multiplex here so we can
	// send an immediate trigger to start the first frame
	s.tickerChannel = make(chan time.Time)
	s.startChannel = make(chan time.Time)
	s.ticker = *time.NewTicker(time.Second / time.Duration(s.speed))
	go tickerMultiplex(s.tickerChannel, s.ticker.C, s.startChannel)
	s.startChannel <- time.Now()

	if (dir == UP && s.currentDirection == DOWN) ||
		(dir == DOWN && s.currentDirection == UP) ||
		(dir == LEFT && s.currentDirection == RIGHT) ||
		(dir == RIGHT && s.currentDirection == LEFT) {
		// User has started in opposite direction, switch head and tail.
		tempPos := s.headPos
		s.headPos = s.tailPos
		s.tailPos = tempPos
		s.currentDirection = dir
	}
}

// tickerMultiplex is used to allow us to send an immdidiate pulse into the 'ticker' channel upon creation
func tickerMultiplex(out chan<- time.Time, tickerIn <-chan time.Time, startIn <-chan time.Time) {
	for true {
		select {
		case t := <-tickerIn:
			out <- t
		case t, ok := <-startIn:
			if !ok {
				return
			}
			out <- t
		}
	}
}

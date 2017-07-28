package main

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type gameCFGType struct {
	gameAreaOrig   pixel.Vec
	gameAreaDims   gameAreaDimsType
	gameArea       pixel.Rect
	gameGridSize   float64
	gameGridMatrix pixel.Matrix
}

type gameAreaDimsType struct {
	x float64
	y float64
}

func newGameCFG(Dims gameAreaDimsType, gridSize float64, winCFG pixelgl.WindowConfig) gameCFGType {
	gameCFG := new(gameCFGType)
	if math.Mod(Dims.x, gridSize) != 0 || math.Mod(Dims.y, gridSize) != 0 {
		panic(errors.New("game Area must be a multiple of the grid size"))
	}
	gameAreaMargin := (winCFG.Bounds.H() - Dims.y) / 2
	gameCFG.gameAreaDims = Dims
	gameCFG.gameArea = pixel.R(gameAreaMargin, gameAreaMargin, gameAreaMargin+Dims.x, gameAreaMargin+Dims.y)
	gameCFG.gameGridSize = gridSize
	gameCFG.gameGridMatrix = pixel.IM.Scaled(pixel.ZV, 1/gridSize)
	return *gameCFG
}

type snakeType struct {
	headPos          pixel.Vec
	tailPos          pixel.Vec
	length           float64
	speed            float64
	currentDirection direction
	pointsList       []pixel.Vec
	gameCFG          *gameCFGType
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

func newSnake(gameCFG gameCFGType) snakeType {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	snake := new(snakeType)
	snake.gameCFG = &gameCFG
	snake.length = 3
	snake.speed = 1.0
	snake.headPos = gameCFG.gameGridMatrix.Project(pixel.V(float64(r.Intn(int(gameCFG.gameAreaDims.x))), float64(r.Intn(int(gameCFG.gameAreaDims.x)))))
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

func (s *snakeType) dir() pixel.Vec {
	return s.currentDirection.val
}

func (s *snakeType) getHeadPos() pixel.Vec {
	return s.gameCFG.gameGridMatrix.Unproject(s.headPos)
}

func (s *snakeType) getTailPos() pixel.Vec {
	return s.gameCFG.gameGridMatrix.Unproject(s.tailPos)
}

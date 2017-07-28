package types

import (
	"errors"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// GameCFGType is a struct used to define the configuration of the game
type GameCFGType struct {
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

// NewGameCFG returns and initialised Game Configuration Struct
func NewGameCFG(xSize float64, ySize float64, gridSize float64, winCFG pixelgl.WindowConfig) GameCFGType {
	gameCFG := new(GameCFGType)
	if math.Mod(xSize, gridSize) != 0 || math.Mod(ySize, gridSize) != 0 {
		panic(errors.New("game Area must be a multiple of the grid size"))
	}
	gameAreaMargin := (winCFG.Bounds.H() - ySize) / 2
	gameCFG.gameAreaDims = gameAreaDimsType{x: xSize, y: ySize}
	gameCFG.gameArea = pixel.R(gameAreaMargin, gameAreaMargin, gameAreaMargin+xSize, gameAreaMargin+ySize)
	gameCFG.gameGridSize = gridSize
	gameCFG.gameGridMatrix = pixel.IM.Scaled(pixel.ZV, 1/gridSize)
	return *gameCFG
}

// GetGridMatrix returns the matrix for the game grid
func (cfg *GameCFGType) GetGridMatrix() pixel.Matrix {
	return cfg.gameGridMatrix
}

// GetGridSize returns the pixel size of the game grid
func (cfg *GameCFGType) GetGridSize() float64 {
	return cfg.gameGridSize
}

// GetGameAreaDims returns the dimensions of the game area
func (cfg *GameCFGType) GetGameAreaDims() (x float64, y float64) {
	return cfg.gameAreaDims.x, cfg.gameAreaDims.y
}

// GetGameAreaAsVecs returns the vectors representing the game area
func (cfg *GameCFGType) GetGameAreaAsVecs() (min pixel.Vec, max pixel.Vec) {
	return cfg.gameArea.Min, cfg.gameArea.Max
}

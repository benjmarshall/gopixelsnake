package types

import (
	"errors"
	"log"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// GameCFGType is a struct used to define the configuration of the game
type GameCFGType struct {
	gameAreaDims            gameAreaDimsType
	gameArea                pixel.Rect
	gameAreaBorderThickness float64
	gameGridSize            float64
	gameGridMatrix          pixel.Matrix
}

type gameAreaDimsType struct {
	x float64
	y float64
}

// NewGameCFG returns and initialised Game Configuration Struct
func NewGameCFG(xSize float64, ySize float64, borderWeight float64, gridSize float64, winCFG pixelgl.WindowConfig) GameCFGType {
	gameCFG := new(GameCFGType)
	if math.Mod(xSize, gridSize) != 0 || math.Mod(ySize, gridSize) != 0 {
		panic(errors.New("game Area must be a multiple of the grid size"))
	}
	gameAreaMargin := (winCFG.Bounds.H() - ySize) / 2
	gameCFG.gameAreaDims = gameAreaDimsType{x: xSize, y: ySize}
	gameCFG.gameArea = pixel.R(gameAreaMargin, gameAreaMargin, gameAreaMargin+xSize, gameAreaMargin+ySize)
	gameCFG.gameAreaBorderThickness = borderWeight
	gameCFG.gameGridSize = gridSize
	gameCFG.gameGridMatrix = pixel.IM.Scaled(pixel.ZV, gridSize).Moved(pixel.V(gridSize/2, gridSize/2))
	// Debug
	log.Println("__Game Config__")
	log.Printf("Game Area Margin: %v", gameAreaMargin)
	log.Printf("Area Dims: %v", gameCFG.gameAreaDims)
	log.Printf("Area Rectangle: %v", gameCFG.gameArea)
	log.Printf("Border Thichness: %v", gameCFG.gameAreaBorderThickness)
	log.Printf("Grid Size: %v", gameCFG.gameGridSize)
	log.Printf("Grid Matrix: %v", gameCFG.gameGridMatrix)
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

// GetGameAreaAsRec returns the rectangle representing the game area
func (cfg *GameCFGType) GetGameAreaAsRec() pixel.Rect {
	return cfg.gameArea
}

// GetBorderWeight returns the weight in pixels of the border of the game area
func (cfg *GameCFGType) GetBorderWeight() float64 {
	return cfg.gameAreaBorderThickness
}

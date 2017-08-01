package game

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Config is a struct used to define the configuration of the game
type Config struct {
	gameAreaDims            gameAreaDimsType
	gameArea                pixel.Rect
	gameAreaBorderThickness float64
	gameGridSize            float64
	gameGridMatrix          pixel.Matrix
	gameWindowMatrix        pixel.Matrix
}

type gameAreaDimsType struct {
	x float64
	y float64
}

// NewGameConfig returns and initialised Game Configuration Struct
func NewGameConfig(xSize float64, ySize float64, borderWeight float64, gridSize float64, winCFG pixelgl.WindowConfig) Config {
	gameCFG := new(Config)
	if math.Mod(xSize, gridSize) != 0 || math.Mod(ySize, gridSize) != 0 {
		panic(errors.New("game Area must be a multiple of the grid size"))
	}
	gameAreaMargin := (winCFG.Bounds.H() - ySize) / 2
	gameCFG.gameAreaDims = gameAreaDimsType{x: xSize, y: ySize}
	gameCFG.gameArea = pixel.R(0, 0, xSize, ySize)
	gameCFG.gameAreaBorderThickness = borderWeight
	gameCFG.gameGridSize = gridSize
	gameCFG.gameGridMatrix = pixel.IM.Scaled(pixel.ZV, gridSize).Moved(pixel.V(gridSize/2, gridSize/2))
	gameCFG.gameWindowMatrix = pixel.IM.Moved(pixel.V(gameAreaMargin, gameAreaMargin))
	// Debug
	log.Println("__Game Config__")
	log.Printf("Game Area Margin: %v", gameAreaMargin)
	log.Printf("Area Dims: %v", gameCFG.gameAreaDims)
	log.Printf("Area Rectangle: %v", gameCFG.gameArea)
	log.Printf("Border Thichness: %v", gameCFG.gameAreaBorderThickness)
	log.Printf("Grid Size: %v", gameCFG.gameGridSize)
	log.Printf("Grid Matrix: %v", gameCFG.gameGridMatrix)
	log.Printf("Window Matrix: %v", gameCFG.gameWindowMatrix)
	return *gameCFG
}

// GetGridMatrix returns the matrix for the game grid which is used to translate the snake coordinates onto the game area.
func (cfg *Config) GetGridMatrix() pixel.Matrix {
	return cfg.gameGridMatrix
}

// GetWindowMatrix returns the matrix for the game area which is used to translate the game area coordinates onto the window.
func (cfg *Config) GetWindowMatrix() pixel.Matrix {
	return cfg.gameWindowMatrix
}

// GetGridSize returns the pixel size of the game grid
func (cfg *Config) GetGridSize() float64 {
	return cfg.gameGridSize
}

// GetGameAreaDims returns the dimensions of the game area
func (cfg *Config) GetGameAreaDims() (x float64, y float64) {
	return cfg.gameAreaDims.x, cfg.gameAreaDims.y
}

// GetGameAreaAsVecs returns the vectors representing the game area in it's native coordinates
func (cfg *Config) GetGameAreaAsVecs() (min pixel.Vec, max pixel.Vec) {
	return cfg.gameArea.Min, cfg.gameArea.Max
}

// GetGameAreaAsRec returns the rectangle representing the game area in it's native coordinates
func (cfg *Config) GetGameAreaAsRec() pixel.Rect {
	return cfg.gameArea
}

// GetBorderWeight returns the weight in pixels of the border of the game area
func (cfg *Config) GetBorderWeight() float64 {
	return cfg.gameAreaBorderThickness
}

// GenerateRandomBerry generates a new berry in a random location
func GenerateRandomBerry(gameCFG *Config) pixel.Vec {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x, y := gameCFG.GetGameAreaDims()
	berryX := float64(r.Intn(int(x/gameCFG.GetGridSize()) - 1))
	berryY := float64(r.Intn(int(y/gameCFG.GetGridSize()) - 1))
	berry := pixel.V(berryX, berryY)
	return gameCFG.GetGridMatrix().Project(berry)
}

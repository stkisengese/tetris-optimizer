package solver

import (
	"fmt"
	"github.com/stkisengese/tetris-optimizer/internal/grid"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
	"math"
	"sort"
)

// Result represents the solution result
type Result struct {
	Grid    *grid.Grid
	Success bool
	Size    int
}

// CalculateMinSquareSize calculates the theoretical minimum square size
// needed to fit all tetrominoes
func CalculateMinSquareSize(tetrominoes []*tetromino.Tetromino) int {
	// Handle edge cases
	if len(tetrominoes) == 0 {
		return 0
	}

	// Count total blocks across all tetrominoes
	totalBlocks := len(tetrominoes) * 4 // Each tetromino has exactly 4 blocks

	// Calculate theoretical minimum: ceil(sqrt(total_blocks))
	minSize := int(math.Ceil(math.Sqrt(float64(totalBlocks))))

	// Safety bounds - ensure minimum size is at least 1
	if minSize < 1 {
		minSize = 1
	}

	// For a single piece, minimum size should be at least the piece's dimensions
	if len(tetrominoes) == 1 {
		t := tetrominoes[0]
		if t.Width > minSize {
			minSize = t.Width
		}
		if t.Height > minSize {
			minSize = t.Height
		}
	}

	return minSize
}

// SolveTetris solves the tetris puzzle using backtracking
func SolveTetris(tetrominoes []*tetromino.Tetromino, gridSize int) (*Result, error) {
	if len(tetrominoes) == 0 {
		return &Result{Success: false, Size: gridSize}, nil
	}

	// Create grid
	g, err := grid.NewGrid(gridSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create grid: %v", err)
	}

	// Sort tetrominoes by size (largest first) for better pruning
	sortedTetrominoes := make([]*tetromino.Tetromino, len(tetrominoes))
	copy(sortedTetrominoes, tetrominoes)
	sort.Slice(sortedTetrominoes, func(i, j int) bool {
		// Sort by area (width * height), then by width, then by height
		areaI := sortedTetrominoes[i].Width * sortedTetrominoes[i].Height
		areaJ := sortedTetrominoes[j].Width * sortedTetrominoes[j].Height

		if areaI != areaJ {
			return areaI > areaJ
		}
		if sortedTetrominoes[i].Width != sortedTetrominoes[j].Width {
			return sortedTetrominoes[i].Width > sortedTetrominoes[j].Width
		}
		return sortedTetrominoes[i].Height > sortedTetrominoes[j].Height
	})

	// Start backtracking
	success := backtrack(g, sortedTetrominoes, 0)

	return &Result{
		Grid:    g,
		Success: success,
		Size:    gridSize,
	}, nil
}

// backtrack implements the recursive backtracking algorithm
func backtrack(g *grid.Grid, tetrominoes []*tetromino.Tetromino, index int) bool {
	// Base case: all tetrominoes placed
	if index >= len(tetrominoes) {
		return true
	}

	current := tetrominoes[index]

	// Try all rotations of the current tetromino
	rotations := current.GenerateRotations()

	for _, rotation := range rotations {
		// Try all positions on the grid
		for y := 0; y <= g.Size-rotation.Height; y++ {
			for x := 0; x <= g.Size-rotation.Width; x++ {
				// Check if we can place the tetromino at this position
				if g.CanPlaceTetromino(rotation, x, y) {
					// Early pruning: check if remaining pieces can fit
					if !canFitRemaining(g, tetrominoes, index, rotation, x, y) {
						continue
					}

					// Place the tetromino
					err := g.PlaceTetromino(rotation, x, y)
					if err != nil {
						continue
					}

					// Recursively try to place the next tetromino
					if backtrack(g, tetrominoes, index+1) {
						return true
					}

					// Backtrack: remove the tetromino
					g.RemoveTetromino(rotation)
				}
			}
		}
	}

	return false
}

// canFitRemaining checks if the remaining tetrominoes can theoretically fit
// in the remaining empty space (pruning optimization)
func canFitRemaining(g *grid.Grid, tetrominoes []*tetromino.Tetromino, currentIndex int, currentPiece *tetromino.Tetromino, x, y int) bool {
	// Calculate remaining pieces
	remainingPieces := len(tetrominoes) - currentIndex - 1
	if remainingPieces <= 0 {
		return true
	}

	// Calculate empty cells after placing current piece
	emptyCells := g.CountEmpty() - 4 // 4 blocks will be occupied by current piece

	// Each remaining piece needs exactly 4 cells
	requiredCells := remainingPieces * 4

	// Simple check: do we have enough empty cells?
	if emptyCells < requiredCells {
		return false
	}

	return true
}

// SolveOptimal finds the optimal solution by trying increasing grid sizes
func SolveOptimal(tetrominoes []*tetromino.Tetromino) (*Result, error) {
	if len(tetrominoes) == 0 {
		return &Result{Success: false, Size: 0}, nil
	}

	// Calculate minimum possible size
	minSize := CalculateMinSquareSize(tetrominoes)

	// Try increasing sizes until we find a solution
	for size := minSize; size <= minSize+4; size++ { // Reasonable upper bound
		result, err := SolveTetris(tetrominoes, size)
		if err != nil {
			return nil, err
		}

		if result.Success {
			return result, nil
		}
	}

	// If no solution found in reasonable range, return the last attempt
	return SolveTetris(tetrominoes, minSize+4)
}

// ValidateSolution checks if a solution is valid
func ValidateSolution(g *grid.Grid, tetrominoes []*tetromino.Tetromino) error {
	// Check that all tetrominoes are placed correctly
	placedPieces := make(map[rune]int)

	for y := 0; y < g.Size; y++ {
		for x := 0; x < g.Size; x++ {
			cell, err := g.GetCell(x, y)
			if err != nil {
				return err
			}

			if cell != '.' {
				placedPieces[cell]++
			}
		}
	}

	// Verify each tetromino is placed exactly once with 4 blocks
	for _, t := range tetrominoes {
		count, exists := placedPieces[t.ID]
		if !exists {
			return fmt.Errorf("tetromino %c not found in solution", t.ID)
		}
		if count != 4 {
			return fmt.Errorf("tetromino %c has %d blocks instead of 4", t.ID, count)
		}
	}

	// Check for extra pieces
	expectedPieces := len(tetrominoes)
	if len(placedPieces) != expectedPieces {
		return fmt.Errorf("solution has %d pieces instead of %d", len(placedPieces), expectedPieces)
	}

	return nil
}

// GetSolutionStats returns statistics about the solution
func GetSolutionStats(result *Result) map[string]interface{} {
	if result == nil || result.Grid == nil {
		return map[string]interface{}{
			"success": false,
		}
	}

	stats := map[string]interface{}{
		"success":     result.Success,
		"grid_size":   result.Size,
		"total_cells": result.Size * result.Size,
		"empty_cells": result.Grid.CountEmpty(),
	}

	if result.Success {
		stats["filled_cells"] = stats["total_cells"].(int) - stats["empty_cells"].(int)
		stats["utilization"] = float64(stats["filled_cells"].(int)) / float64(stats["total_cells"].(int))
	}

	return stats
}

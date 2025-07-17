package solver

import (
	"fmt"
	"math"
	"github.com/stkisengese/tetris-optimizer/internal/grid"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
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
	return int(math.Ceil(math.Sqrt(float64(totalBlocks))))
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

	success := backtrack(g, tetrominoes, 0)

	return &Result{
		Grid:    g,
		Success: success,
		Size:    gridSize,
	}, nil
}

// backtrack implements simple recursive backtracking
func backtrack(g *grid.Grid, tetrominoes []*tetromino.Tetromino, index int) bool {
	// Base case: all tetrominoes placed
	if index >= len(tetrominoes) {
		return true
	}

	current := tetrominoes[index]
	
	// Try all possible rotations
	rotations := current.GenerateRotations()
	for _, rotation := range rotations {
		// Try all possible positions
		for y := 0; y <= g.Size-rotation.Height; y++ {
			for x := 0; x <= g.Size-rotation.Width; x++ {
				if g.CanPlaceTetromino(rotation, x, y) {
					// Place the piece/tetromino
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

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


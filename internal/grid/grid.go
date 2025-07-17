package grid

import (
	"fmt"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
	"strings"
)

// Grid represents the solution board
type Grid struct {
	// Size is the dimension of the square grid (size x size)
	Size int

	// Cells contains the grid data, where each cell contains:
	// - '.' for empty
	// - Letter (A-Z) for tetromino pieces
	Cells [][]rune
}

// NewGrid creates a new empty grid of the specified size
func NewGrid(size int) (*Grid, error) {
	if size <= 0 {
		return nil, fmt.Errorf("grid size must be positive, got %d", size)
	}

	cells := make([][]rune, size)
	for i := range cells {
		cells[i] = make([]rune, size)
		for j := range cells[i] {
			cells[i][j] = '.'
		}
	}

	return &Grid{
		Size:  size,
		Cells: cells,
	}, nil
}

// IsEmpty checks if a cell is empty
func (g *Grid) IsEmpty(x, y int) bool {
	if !g.IsValidPosition(x, y) {
		return false
	}
	return g.Cells[y][x] == '.'
}

// IsValidPosition checks if the coordinates are within grid bounds
func (g *Grid) IsValidPosition(x, y int) bool {
	return x >= 0 && x < g.Size && y >= 0 && y < g.Size
}

// CanPlaceTetromino checks if a tetromino can be placed at the given position
func (g *Grid) CanPlaceTetromino(t *tetromino.Tetromino, x, y int) bool {
	for _, point := range t.Points {
		newX := x + point.X
		newY := y + point.Y

		// Check bounds
		if !g.IsValidPosition(newX, newY) {
			return false
		}

		// Check if cell is empty
		if !g.IsEmpty(newX, newY) {
			return false
		}
	}

	return true
}

// PlaceTetromino places a tetromino on the grid at the given position
func (g *Grid) PlaceTetromino(t *tetromino.Tetromino, x, y int) error {
	if !g.CanPlaceTetromino(t, x, y) {
		return fmt.Errorf("cannot place tetromino %c at position (%d, %d)", t.ID, x, y)
	}

	for _, point := range t.Points {
		newX := x + point.X
		newY := y + point.Y
		g.Cells[newY][newX] = t.ID
	}

	// Update tetromino position
	t.SetPosition(x, y)

	return nil
}

// RemoveTetromino removes a tetromino from the grid
func (g *Grid) RemoveTetromino(t *tetromino.Tetromino) {
	absolutePoints := t.GetAbsolutePoints()
	for _, point := range absolutePoints {
		if g.IsValidPosition(point.X, point.Y) {
			g.Cells[point.Y][point.X] = '.'
		}
	}
}

// CountEmpty returns the number of empty cells in the grid
func (g *Grid) CountEmpty() int {
	count := 0
	for _, row := range g.Cells {
		for _, cell := range row {
			if cell == '.' {
				count++
			}
		}
	}
	return count
}

// String returns a string representation of the grid
func (g *Grid) String() string {
	var builder strings.Builder

	for _, row := range g.Cells {
		builder.WriteString(string(row))
		builder.WriteString("\n")
	}

	return builder.String()
}

// Clear resets all cells to empty
func (g *Grid) Clear() {
	for i := range g.Cells {
		for j := range g.Cells[i] {
			g.Cells[i][j] = '.'
		}
	}
}

// GetCell returns the content of a cell at the given position
func (g *Grid) GetCell(x, y int) (rune, error) {
	if !g.IsValidPosition(x, y) {
		return '.', fmt.Errorf("position (%d, %d) is out of bounds", x, y)
	}
	return g.Cells[y][x], nil
}

// SetCell sets the content of a cell at the given position
func (g *Grid) SetCell(x, y int, value rune) error {
	if !g.IsValidPosition(x, y) {
		return fmt.Errorf("position (%d, %d) is out of bounds", x, y)
	}
	g.Cells[y][x] = value
	return nil
}

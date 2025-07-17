package tests

import (
	"github.com/stkisengese/tetris-optimizer/internal/grid"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
	"testing"
)

func TestNewGrid(t *testing.T) {
	g, err := grid.NewGrid(4)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if g.Size != 4 {
		t.Errorf("Expected size 4, got %d", g.Size)
	}

	// Check that all cells are empty
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if !g.IsEmpty(i, j) {
				t.Errorf("Cell (%d, %d) should be empty", i, j)
			}
		}
	}

	// Test invalid size
	_, err = grid.NewGrid(-1)
	if err == nil {
		t.Error("Expected error for negative size")
	}
}

func TestGridBounds(t *testing.T) {
	g, _ := grid.NewGrid(3)

	// Test valid positions
	if !g.IsValidPosition(0, 0) {
		t.Error("(0,0) should be valid")
	}

	if !g.IsValidPosition(2, 2) {
		t.Error("(2,2) should be valid")
	}

	// Test invalid positions
	if g.IsValidPosition(-1, 0) {
		t.Error("(-1,0) should be invalid")
	}

	if g.IsValidPosition(3, 2) {
		t.Error("(3,2) should be invalid")
	}
}

func TestGridTetrominoPlacement(t *testing.T) {
	g, _ := grid.NewGrid(4)

	// Create a simple tetromino
	tetrominoGrid := []string{
		"##..",
		"##..",
		"....",
		"....",
	}

	tetro, err := tetromino.NewTetromino('O', tetrominoGrid)
	if err != nil {
		t.Fatalf("Expected no error creating tetromino, got %v", err)
	}

	// Test valid placement
	if !g.CanPlaceTetromino(tetro, 0, 0) {
		t.Error("Should be able to place tetromino at (0,0)")
	}

	err = g.PlaceTetromino(tetro, 0, 0)
	if err != nil {
		t.Errorf("Expected no error placing tetromino, got %v", err)
	}

	// Check that cells are occupied
	cell, _ := g.GetCell(0, 0)
	if cell != 'O' {
		t.Errorf("Expected 'O' at (0,0), got %c", cell)
	}

	// Test invalid placement (overlapping)
	if g.CanPlaceTetromino(tetro, 0, 0) {
		t.Error("Should not be able to place tetromino on occupied cells")
	}
}

func TestGridCopy(t *testing.T) {
	g, _ := grid.NewGrid(2)
	g.SetCell(0, 0, 'A')

	copy := g.Copy()

	// Check that copy has the same content
	cell, _ := copy.GetCell(0, 0)
	if cell != 'A' {
		t.Error("Copy should have same content as original")
	}

	// Modify copy and ensure original is unchanged
	copy.SetCell(0, 0, 'B')

	originalCell, _ := g.GetCell(0, 0)
	if originalCell != 'A' {
		t.Error("Original should not be affected by copy modification")
	}
}

func TestGridUtilities(t *testing.T) {
	g, _ := grid.NewGrid(2)

	// Test empty count
	if g.CountEmpty() != 4 {
		t.Errorf("Expected 4 empty cells, got %d", g.CountEmpty())
	}

	// Test IsFull
	if g.IsFull() {
		t.Error("Grid should not be full")
	}

	// Fill grid
	g.SetCell(0, 0, 'A')
	g.SetCell(0, 1, 'B')
	g.SetCell(1, 0, 'C')
	g.SetCell(1, 1, 'D')

	if !g.IsFull() {
		t.Error("Grid should be full")
	}

	if g.CountEmpty() != 0 {
		t.Errorf("Expected 0 empty cells, got %d", g.CountEmpty())
	}
}

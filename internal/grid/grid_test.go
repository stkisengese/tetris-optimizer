package grid_test

import (
	"testing"

	"github.com/stkisengese/tetris-optimizer/internal/grid"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
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

	// Test invalid placement (overlapping)
	if g.CanPlaceTetromino(tetro, 0, 0) {
		t.Error("Should not be able to place tetromino on occupied cells")
	}
}

// Add these test functions to your existing grid_test.go file

func TestRemoveTetromino(t *testing.T) {
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

	// Place the tetromino first
	err = g.PlaceTetromino(tetro, 0, 0)
	if err != nil {
		t.Fatalf("Expected no error placing tetromino, got %v", err)
	}

	// Verify it's placed (cells should not be empty)
	if g.IsEmpty(0, 0) || g.IsEmpty(0, 1) || g.IsEmpty(1, 0) || g.IsEmpty(1, 1) {
		t.Error("Tetromino should be placed on the grid")
	}

	// Remove the tetromino
	g.RemoveTetromino(tetro)

	// Verify it's removed (cells should be empty again)
	if !g.IsEmpty(0, 0) || !g.IsEmpty(0, 1) || !g.IsEmpty(1, 0) || !g.IsEmpty(1, 1) {
		t.Error("Tetromino should be removed from the grid")
	}

	// Test removing tetromino at different position
	err = g.PlaceTetromino(tetro, 2, 2)
	if err != nil {
		t.Fatalf("Expected no error placing tetromino at (2,2), got %v", err)
	}

	g.RemoveTetromino(tetro)

	// Verify removal at new position
	if !g.IsEmpty(2, 2) || !g.IsEmpty(2, 3) || !g.IsEmpty(3, 2) || !g.IsEmpty(3, 3) {
		t.Error("Tetromino should be removed from position (2,2)")
	}
}

func TestRemoveTetrominoOutOfBounds(t *testing.T) {
	g, _ := grid.NewGrid(3)

	// Create a tetromino that would extend beyond grid bounds
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

	// Position tetromino so it would extend beyond bounds
	// This should not crash the program due to the IsValidPosition check
	tetro.SetPosition(2, 2) // This would make some points go beyond 3x3 grid

	// This should not panic or cause errors
	g.RemoveTetromino(tetro)

	// Grid should remain unchanged since no valid positions were affected
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !g.IsEmpty(i, j) {
				t.Errorf("Grid cell (%d, %d) should remain empty", i, j)
			}
		}
	}
}

func TestGridString(t *testing.T) {
	g, _ := grid.NewGrid(3)

	// Test empty grid string representation
	expected := "...\n...\n...\n"
	result := g.String()

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}

	// Test grid with placed tetromino
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

	err = g.PlaceTetromino(tetro, 0, 0)
	if err != nil {
		t.Fatalf("Expected no error placing tetromino, got %v", err)
	}

	expectedWithTetromino := "OO.\nOO.\n...\n"
	resultWithTetromino := g.String()

	if resultWithTetromino != expectedWithTetromino {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedWithTetromino, resultWithTetromino)
	}
}

func TestGridStringLargerGrid(t *testing.T) {
	g, _ := grid.NewGrid(5)

	// Test larger empty grid
	expected := ".....\n.....\n.....\n.....\n.....\n"
	result := g.String()

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}

	// Place multiple tetrominoes and test string representation
	tetrominoGrid1 := []string{
		"##..",
		"##..",
		"....",
		"....",
	}

	tetrominoGrid2 := []string{
		"....",
		"..##",
		"..##",
		"....",
	}

	tetro1, _ := tetromino.NewTetromino('A', tetrominoGrid1)
	tetro2, _ := tetromino.NewTetromino('B', tetrominoGrid2)

	g.PlaceTetromino(tetro1, 0, 0)
	g.PlaceTetromino(tetro2, 2, 2)

	expectedMultiple := "AA...\nAA...\n..BB.\n..BB.\n.....\n"
	resultMultiple := g.String()

	if resultMultiple != expectedMultiple {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedMultiple, resultMultiple)
	}
}

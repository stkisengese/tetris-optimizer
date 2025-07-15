package tests

import (
	"testing"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
)

func TestPoint(t *testing.T) {
	p1 := tetromino.Point{X: 1, Y: 2}
	p2 := tetromino.Point{X: 3, Y: 4}
	
	// Test String
	if p1.String() != "(1,2)" {
		t.Errorf("Expected (1,2), got %s", p1.String())
	}
	
	// Test Equals
	if p1.Equals(p2) {
		t.Error("Points should not be equal")
	}
	
	if !p1.Equals(tetromino.Point{X: 1, Y: 2}) {
		t.Error("Points should be equal")
	}
	
	// Test Add
	result := p1.Add(p2)
	expected := tetromino.Point{X: 4, Y: 6}
	if !result.Equals(expected) {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestNewTetromino(t *testing.T) {
	// Test valid L-shaped tetromino
	grid := []string{
		"#...",
		"#...",
		"##..",
		"....",
	}
	
	tetro, err := tetromino.NewTetromino('L', grid)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if tetro.ID != 'L' {
		t.Errorf("Expected ID 'L', got %c", tetro.ID)
	}
	
	if len(tetro.Points) != 4 {
		t.Errorf("Expected 4 points, got %d", len(tetro.Points))
	}
	
	if tetro.Width != 2 || tetro.Height != 3 {
		t.Errorf("Expected dimensions 2x3, got %dx%d", tetro.Width, tetro.Height)
	}
	
	// Test invalid grid (wrong number of blocks)
	invalidGrid := []string{
		"#...",
		"#...",
		"....",
		"....",
	}
	
	_, err = tetromino.NewTetromino('I', invalidGrid)
	if err == nil {
		t.Error("Expected error for invalid tetromino")
	}
}

func TestTetrominoRotation(t *testing.T) {
	// Test L-shaped tetromino rotation
	grid := []string{
		"#...",
		"#...",
		"##..",
		"....",
	}
	
	tetro, err := tetromino.NewTetromino('L', grid)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	originalWidth := tetro.Width
	originalHeight := tetro.Height
	
	// Rotate 90 degrees
	tetro.Rotate90()
	
	// Check dimensions swapped
	if tetro.Width != originalHeight || tetro.Height != originalWidth {
		t.Errorf("Expected dimensions to swap after rotation")
	}
	
	// Should still have 4 points
	if len(tetro.Points) != 4 {
		t.Errorf("Expected 4 points after rotation, got %d", len(tetro.Points))
	}
}

func TestTetrominoCopy(t *testing.T) {
	grid := []string{
		"####",
		"....",
		"....",
		"....",
	}
	
	original, err := tetromino.NewTetromino('I', grid)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	copy := original.Copy()
	
	// Check that copy is equal to original
	if !original.Equals(copy) {
		t.Error("Copy should be equal to original")
	}
	
	// Modify copy and ensure original is unchanged
	copy.SetPosition(5, 5)
	
	if original.Position.X != 0 || original.Position.Y != 0 {
		t.Error("Original should not be affected by copy modification")
	}
}
package tetromino_test

import (
	"testing"

	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
)

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

func TestGenerateRotations(t *testing.T) {
	// Test I-piece (should have 2 unique rotations)
	grid := []string{
		"....",
		"####",
		"....",
		"....",
	}

	tetro, err := tetromino.NewTetromino('I', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}

	rotations := tetro.GenerateRotations()

	// I-piece should have 2 unique rotations (horizontal and vertical)
	if len(rotations) != 2 {
		t.Errorf("Expected 2 unique rotations for I-piece, got %d", len(rotations))
	}

	// Test O-piece (should have 1 unique rotation)
	grid = []string{
		"....",
		".##.",
		".##.",
		"....",
	}

	tetro, err = tetromino.NewTetromino('O', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}

	rotations = tetro.GenerateRotations()

	// O-piece should have 1 unique rotation (it's symmetric)
	if len(rotations) != 1 {
		t.Errorf("Expected 1 unique rotation for O-piece, got %d", len(rotations))
	}

	// Test L-piece (should have 4 unique rotations)
	grid = []string{
		"#...",
		"#...",
		"##..",
		"....",
	}

	tetro, err = tetromino.NewTetromino('L', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}

	rotations = tetro.GenerateRotations()

	// L-piece should have 4 unique rotations
	if len(rotations) != 4 {
		t.Errorf("Expected 4 unique rotations for L-piece, got %d", len(rotations))
	}
}

func TestShapeKey(t *testing.T) {
	// Create two identical tetrominoes
	grid := []string{
		"#...",
		"#...",
		"##..",
		"....",
	}

	tetro1, err := tetromino.NewTetromino('L', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}

	tetro2, err := tetromino.NewTetromino('L', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}

	// They should have the same shape key
	key1 := tetro1.ShapeKey()
	key2 := tetro2.ShapeKey()

	if key1 != key2 {
		t.Errorf("Expected identical tetrominoes to have same shape key, got %s vs %s", key1, key2)
	}

	// Rotate one and they should have different shape keys
	tetro2.Rotate90()
	key2 = tetro2.ShapeKey()

	if key1 == key2 {
		t.Error("Expected rotated tetrominoes to have different shape keys")
	}
}

func TestCloneAlias(t *testing.T) {
	grid := []string{
		"#...",
		"#...",
		"##..",
		"....",
	}

	tetro, err := tetromino.NewTetromino('L', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}

	// Test that Clone() is an alias for Copy()
	cloned := tetro.Clone()

	// Test that it's a deep copy
	cloned.SetPosition(5, 5)
	if tetro.Position.X != 0 || tetro.Position.Y != 0 {
		t.Error("Expected original tetromino to be unchanged after cloning")
	}
}

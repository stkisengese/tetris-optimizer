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


func TestGetBlocks(t *testing.T) {
	// Test L-piece
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
	
	// Test at origin
	blocks := tetro.GetBlocks()
	expected := []tetromino.Point{{0, 0}, {0, 1}, {0, 2}, {1, 2}}
	
	if len(blocks) != 4 {
		t.Errorf("Expected 4 blocks, got %d", len(blocks))
	}
	
	for i, block := range blocks {
		if !block.Equals(expected[i]) {
			t.Errorf("Block %d: expected %v, got %v", i, expected[i], block)
		}
	}
	
	// Test after moving position
	tetro.SetPosition(5, 3)
	blocks = tetro.GetBlocks()
	expectedMoved := []tetromino.Point{{5, 3}, {5, 4}, {5, 5}, {6, 5}}
	
	for i, block := range blocks {
		if !block.Equals(expectedMoved[i]) {
			t.Errorf("Moved block %d: expected %v, got %v", i, expectedMoved[i], block)
		}
	}
}

func TestGetBounds(t *testing.T) {
	// Test T-piece
	grid := []string{
		"....",
		".###",
		"..#.",
		"....",
	}
	
	tetro, err := tetromino.NewTetromino('T', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}
	
	// Test at origin
	minX, minY, maxX, maxY := tetro.GetBounds()
	if minX != 0 || minY != 0 || maxX != 2 || maxY != 1 {
		t.Errorf("Expected bounds (0,0,2,1), got (%d,%d,%d,%d)", minX, minY, maxX, maxY)
	}
	
	// Test after moving position
	tetro.SetPosition(10, 5)
	minX, minY, maxX, maxY = tetro.GetBounds()
	if minX != 10 || minY != 5 || maxX != 12 || maxY != 6 {
		t.Errorf("Expected bounds (10,5,12,6), got (%d,%d,%d,%d)", minX, minY, maxX, maxY)
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

func TestNormalize(t *testing.T) {
	// Create a tetromino and move it away from origin
	grid := []string{
		"....",
		".##.",
		".##.",
		"....",
	}
	
	tetro, err := tetromino.NewTetromino('O', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}
	
	// Modify points to simulate a moved tetromino
	tetro.Points = []tetromino.Point{{5, 3}, {6, 3}, {5, 4}, {6, 4}}
	
	// Normalize should move back to origin
	tetro.Normalize()
	
	expected := []tetromino.Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	for i, point := range tetro.Points {
		if !point.Equals(expected[i]) {
			t.Errorf("Point %d: expected %v, got %v", i, expected[i], point)
		}
	}
}

func TestTranslate(t *testing.T) {
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
	
	// Initial position should be (0,0)
	if tetro.Position.X != 0 || tetro.Position.Y != 0 {
		t.Errorf("Expected initial position (0,0), got (%d,%d)", tetro.Position.X, tetro.Position.Y)
	}
	
	// Translate by (3, 5)
	tetro.Translate(3, 5)
	
	if tetro.Position.X != 3 || tetro.Position.Y != 5 {
		t.Errorf("Expected position (3,5), got (%d,%d)", tetro.Position.X, tetro.Position.Y)
	}
	
	// Translate again by (-1, 2)
	tetro.Translate(-1, 2)
	
	if tetro.Position.X != 2 || tetro.Position.Y != 7 {
		t.Errorf("Expected position (2,7), got (%d,%d)", tetro.Position.X, tetro.Position.Y)
	}
}

func TestIsEquivalent(t *testing.T) {
	// Create two L-pieces
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
	
	// They should be equivalent
	if !tetro1.IsEquivalent(tetro2) {
		t.Error("Expected identical tetrominoes to be equivalent")
	}
	
	// Rotate one and they should still be equivalent
	tetro2.Rotate90()
	if !tetro1.IsEquivalent(tetro2) {
		t.Error("Expected rotated tetrominoes to be equivalent")
	}
	
	// Test with different IDs
	tetro3, err := tetromino.NewTetromino('T', grid)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}
	
	if tetro1.IsEquivalent(tetro3) {
		t.Error("Expected tetrominoes with different IDs to not be equivalent")
	}
	
	// Test with different shapes
	gridI := []string{
		"....",
		"####",
		"....",
		"....",
	}
	
	tetroI, err := tetromino.NewTetromino('I', gridI)
	if err != nil {
		t.Fatalf("Failed to create tetromino: %v", err)
	}
	
	if tetro1.IsEquivalent(tetroI) {
		t.Error("Expected tetrominoes with different shapes to not be equivalent")
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
	copied := tetro.Copy()
	
	if !cloned.Equals(copied) {
		t.Error("Expected Clone() to work the same as Copy()")
	}
	
	// Test that it's a deep copy
	cloned.SetPosition(5, 5)
	if tetro.Position.X != 0 || tetro.Position.Y != 0 {
		t.Error("Expected original tetromino to be unchanged after cloning")
	}
}

func TestRotationPreservesShape(t *testing.T) {
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
	
	original := tetro.Copy()
	
	// Rotate 4 times should bring it back to original
	for i := 0; i < 4; i++ {
		tetro.Rotate90()
	}
	
	// Should have same shape (though points might be in different order)
	if !tetro.IsEquivalent(original) {
		t.Error("Expected tetromino to be equivalent to original after 4 rotations")
	}
}

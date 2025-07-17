package tests

import (
	"github.com/stkisengese/tetris-optimizer/internal/solver"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
	"testing"
)

func TestCalculateMinSquareSize(t *testing.T) {
	tests := []struct {
		name        string
		tetrominoes []*tetromino.Tetromino
		expected    int
	}{
		{
			name:        "empty input",
			tetrominoes: []*tetromino.Tetromino{},
			expected:    0,
		},
		{
			name:        "single piece",
			tetrominoes: createTestTetrominoes(1),
			expected:    4, // ceil(sqrt(4)) = 2
		},
		{
			name:        "two pieces",
			tetrominoes: createTestTetrominoes(2),
			expected:    3, // ceil(sqrt(8)) = 3
		},
		{
			name:        "four pieces",
			tetrominoes: createTestTetrominoes(4),
			expected:    4, // ceil(sqrt(16)) = 4
		},
		{
			name:        "five pieces",
			tetrominoes: createTestTetrominoes(5),
			expected:    5, // ceil(sqrt(20)) = 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := solver.CalculateMinSquareSize(tt.tetrominoes)
			if result != tt.expected {
				t.Errorf("CalculateMinSquareSize() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestGetSolutionStats(t *testing.T) {
	// Test with nil result
	stats := solver.GetSolutionStats(nil)
	if stats["success"] != false {
		t.Errorf("Expected success=false for nil result")
	}

	// Test with valid result
	tetrominoes := createLPiece()
	result, err := solver.SolveTetris(tetrominoes, 3)
	if err != nil {
		t.Fatalf("Failed to create test solution: %v", err)
	}

	stats = solver.GetSolutionStats(result)
	if stats["grid_size"] != 3 {
		t.Errorf("Expected grid_size=3, got %v", stats["grid_size"])
	}
	if stats["total_cells"] != 9 {
		t.Errorf("Expected total_cells=9, got %v", stats["total_cells"])
	}
}

// Helper functions for creating test tetrominoes

func createTestTetrominoes(count int) []*tetromino.Tetromino {
	tetrominoes := make([]*tetromino.Tetromino, count)
	for i := 0; i < count; i++ {
		// Create a simple I-piece for testing
		grid := []string{
			"#...",
			"#...",
			"#...",
			"#...",
		}
		tetro, _ := tetromino.NewTetromino(rune('A'+i), grid)
		tetrominoes[i] = tetro
	}
	return tetrominoes
}

func createLPiece() []*tetromino.Tetromino {
	grid := []string{
		"#...",
		"#...",
		"##..",
		"....",
	}
	tetro, _ := tetromino.NewTetromino('A', grid)
	return []*tetromino.Tetromino{tetro}
}

func createIPiece() []*tetromino.Tetromino {
	grid := []string{
		"#...",
		"#...",
		"#...",
		"#...",
	}
	tetro, _ := tetromino.NewTetromino('A', grid)
	return []*tetromino.Tetromino{tetro}
}

func createSquarePiece() []*tetromino.Tetromino {
	grid := []string{
		"##..",
		"##..",
		"....",
		"....",
	}
	tetro, _ := tetromino.NewTetromino('A', grid)
	return []*tetromino.Tetromino{tetro}
}

// Benchmark tests
func BenchmarkCalculateMinSquareSize(b *testing.B) {
	tetrominoes := createTestTetrominoes(10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		solver.CalculateMinSquareSize(tetrominoes)
	}
}

func BenchmarkSolveTetris(b *testing.B) {
	tetrominoes := createTestTetrominoes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		solver.SolveTetris(tetrominoes, 3)
	}
}

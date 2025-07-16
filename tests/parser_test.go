package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"github.com/stkisengese/tetris-optimizer/internal/parser"
)

// Helper function to create temporary test files
func createTempFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	return tmpFile
}

func TestValidTetrominoFile(t *testing.T) {
	content := `#...
#...
##..
....

....
....
....
####

.###
...#
....
....`

	tmpFile := createTempFile(t, content)
	
	tetrominoes, err := parser.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error for valid file, got: %v", err)
	}
	
	if len(tetrominoes) != 3 {
		t.Errorf("Expected 3 tetrominoes, got %d", len(tetrominoes))
	}
	
	// Check IDs
	expectedIDs := []rune{'A', 'B', 'C'}
	for i, tetro := range tetrominoes {
		if tetro.ID != expectedIDs[i] {
			t.Errorf("Expected tetromino %d to have ID %c, got %c", i, expectedIDs[i], tetro.ID)
		}
	}
}

func TestInvalidFileFormat(t *testing.T) {
	testCases := []struct {
		name    string
		content string
	}{
		{
			name: "Wrong number of blocks",
			content: `##..
##..
....
....`,
		},
		{
			name: "Invalid character",
			content: `#...
#...
##..
.x..`,
		},
		{
			name: "Line too long",
			content: `#....
#...
##..
....`,
		},
		{
			name: "Incomplete grid",
			content: `#...
#...
##..`,
		},
		{
			name: "Disconnected blocks",
			content: `#...
....
....
...#`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := createTempFile(t, tc.content)
			
			_, err := parser.ReadFile(tmpFile)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tc.name)
			}
		})
	}
}

func TestEmptyFile(t *testing.T) {
	tmpFile := createTempFile(t, "")
	
	_, err := parser.ReadFile(tmpFile)
	if err == nil {
		t.Error("Expected error for empty file, got nil")
	}
}

func TestWhitespaceOnlyFile(t *testing.T) {
	tmpFile := createTempFile(t, "   \n\t\n  \n")
	
	_, err := parser.ReadFile(tmpFile)
	if err == nil {
		t.Error("Expected error for whitespace-only file, got nil")
	}
}

func TestFileWithTrailingWhitespace(t *testing.T) {
	content := `#...  
#...	
##..   
....`

	tmpFile := createTempFile(t, content)
	
	tetrominoes, err := parser.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error for file with trailing whitespace, got: %v", err)
	}
	
	if len(tetrominoes) != 1 {
		t.Errorf("Expected 1 tetromino, got %d", len(tetrominoes))
	}
}

func TestFileWithShortLines(t *testing.T) {
	content := `#
#
##
.`

	tmpFile := createTempFile(t, content)
	
	tetrominoes, err := parser.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error for file with short lines, got: %v", err)
	}
	
	if len(tetrominoes) != 1 {
		t.Errorf("Expected 1 tetromino, got %d", len(tetrominoes))
	}
}

func TestMultipleTetrominoesWithEmptyLines(t *testing.T) {
	content := `#...
#...
##..
....

....
....
....
####


.###
...#
....
....`

	tmpFile := createTempFile(t, content)
	
	tetrominoes, err := parser.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error for valid file with multiple empty lines, got: %v", err)
	}
	
	if len(tetrominoes) != 3 {
		t.Errorf("Expected 3 tetrominoes, got %d", len(tetrominoes))
	}
}

func TestValidateFile(t *testing.T) {
	validContent := `#...
#...
##..
....`

	tmpFile := createTempFile(t, validContent)
	
	err := parser.ValidateFile(tmpFile)
	if err != nil {
		t.Errorf("Expected no error for valid file, got: %v", err)
	}
	
	// Test invalid file
	invalidContent := `#...
#...
##..`

	tmpFile2 := createTempFile(t, invalidContent)
	
	err = parser.ValidateFile(tmpFile2)
	if err == nil {
		t.Error("Expected error for invalid file, got nil")
	}
}

func TestParseError(t *testing.T) {
	err := parser.NewParseError("test error", 5, "test.txt")
	
	expected := "parse error at line 5: test error"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
	
	// Test error without line number
	err2 := parser.NewParseError("test error", 0, "test.txt")
	expected2 := "parse error: test error"
	if err2.Error() != expected2 {
		t.Errorf("Expected error message %q, got %q", expected2, err2.Error())
	}
}

func TestUtilityFunctions(t *testing.T) {
	grid := []string{
		"#...",
		"#...",
		"##..",
		"....",
	}
	
	// Test CountBlocks
	blocks := parser.CountBlocks(grid)
	if blocks != 4 {
		t.Errorf("Expected 4 blocks, got %d", blocks)
	}
	
	// Test NormalizeGrid
	shortGrid := []string{"#", "#", "##", "."}
	normalized, err := parser.NormalizeGrid(shortGrid)
	if err != nil {
		t.Errorf("Expected no error normalizing grid, got: %v", err)
	}
	
	for i, line := range normalized {
		if len(line) != 4 {
			t.Errorf("Expected line %d to be 4 characters, got %d", i, len(line))
		}
	}
	
	// Test GetTetrominoStats
	blocks, minX, maxX, minY, maxY := parser.GetTetrominoStats(grid)
	if blocks != 4 {
		t.Errorf("Expected 4 blocks in stats, got %d", blocks)
	}
	if minX != 0 || maxX != 1 || minY != 0 || maxY != 2 {
		t.Errorf("Expected bounds (0,1,0,2), got (%d,%d,%d,%d)", minX, maxX, minY, maxY)
	}
}

func TestNonExistentFile(t *testing.T) {
	_, err := parser.ReadFile("non_existent_file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	
	if !strings.Contains(err.Error(), "cannot open file") {
		t.Errorf("Expected 'cannot open file' error, got: %v", err)
	}
}

func TestComplexValidTetromino(t *testing.T) {
	// Test all standard tetromino shapes
	tetrominoes := []string{
		// I-piece
		`....
####
....
....`,
		// O-piece
		`....
.##.
.##.
....`,
		// T-piece
		`....
.###
..#.
....`,
		// S-piece
		`....
..##
.##.
....`,
		// Z-piece
		`....
.##.
..##
....`,
		// J-piece
		`....
.#..
.###
....`,
		// L-piece
		`....
...#
.###
....`,
	}
	
	for i, tetrominoStr := range tetrominoes {
		t.Run(fmt.Sprintf("Tetromino_%d", i), func(t *testing.T) {
			tmpFile := createTempFile(t, tetrominoStr)
			
			pieces, err := parser.ReadFile(tmpFile)
			if err != nil {
				t.Errorf("Expected no error for valid tetromino %d, got: %v", i, err)
			}
			
			if len(pieces) != 1 {
				t.Errorf("Expected 1 piece for tetromino %d, got %d", i, len(pieces))
			}
		})
	}
}
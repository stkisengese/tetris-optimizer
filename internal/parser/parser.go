package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
)

// ParseError represents errors that occur during parsing
type ParseError struct {
	Message string
	Line    int
	File    string
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("parse error at line %d: %s", e.Line, e.Message)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

// NewParseError creates a new parse error
func NewParseError(message string, line int, file string) *ParseError {
	return &ParseError{
		Message: message,
		Line:    line,
		File:    file,
	}
}

// ReadFile reads and parses tetromino definitions from a file
func ReadFile(filename string) ([]*tetromino.Tetromino, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("cannot open file: %v", err), 0, filename)
	}
	defer file.Close()

	return ParseTetrominoes(file, filename)
}

// ParseTetrominoes parses tetrominoes from a reader
func ParseTetrominoes(file *os.File, filename string) ([]*tetromino.Tetromino, error) {
	scanner := bufio.NewScanner(file)
	var tetrominoes []*tetromino.Tetromino
	var currentGrid []string
	var currentID rune = 'A'
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Handle empty lines as separators
		if strings.TrimSpace(line) == "" {
			if len(currentGrid) > 0 {
				// Process current tetromino
				tetro, err := processTetromino(currentGrid, currentID, lineNumber, filename)
				if err != nil {
					return nil, err
				}
				tetrominoes = append(tetrominoes, tetro)
				currentGrid = []string{}
				currentID++
			}
			continue
		}

		// Add line to current grid
		currentGrid = append(currentGrid, line)

		// Check if we have a complete 4x4 grid
		if len(currentGrid) == 4 {
			tetro, err := processTetromino(currentGrid, currentID, lineNumber, filename)
			if err != nil {
				return nil, err
			}
			tetrominoes = append(tetrominoes, tetro)
			currentGrid = []string{}
			currentID++
		}
	}

	// Handle any remaining grid at end of file
	if len(currentGrid) > 0 {
		tetro, err := processTetromino(currentGrid, currentID, lineNumber, filename)
		if err != nil {
			return nil, err
		}
		tetrominoes = append(tetrominoes, tetro)
	}

	if err := scanner.Err(); err != nil {
		return nil, NewParseError(fmt.Sprintf("error reading file: %v", err), lineNumber, filename)
	}

	// Validate that we have at least one tetromino
	if len(tetrominoes) == 0 {
		return nil, NewParseError("no valid tetrominoes found in file", 0, filename)
	}

	return tetrominoes, nil
}

// processTetromino validates and creates a tetromino from a 4x4 grid
func processTetromino(grid []string, id rune, lineNumber int, filename string) (*tetromino.Tetromino, error) {
	// Validate grid size
	if len(grid) != 4 {
		return nil, NewParseError(fmt.Sprintf("tetromino must be 4x4 grid, got %d lines", len(grid)), lineNumber, filename)
	}

	// Validate and normalize each line
	normalizedGrid := make([]string, 4)
	for i, line := range grid {
		// Remove trailing whitespace but preserve length validation
		trimmed := strings.TrimRight(line, " \t")
		
		// Check for invalid characters
		for j, char := range trimmed {
			if char != '#' && char != '.' {
				return nil, NewParseError(fmt.Sprintf("invalid character '%c' at position %d in line %d", char, j, i+1), lineNumber, filename)
			}
		}
		
		// Ensure line is exactly 4 characters (pad with dots if needed)
		if len(trimmed) > 4 {
			return nil, NewParseError(fmt.Sprintf("line %d too long: expected 4 characters, got %d", i+1, len(trimmed)), lineNumber, filename)
		}
		
		// Pad with dots if line is shorter than 4 characters
		for len(trimmed) < 4 {
			trimmed += "."
		}
		
		normalizedGrid[i] = trimmed
	}

	// Validate tetromino has exactly 4 blocks
	blockCount := 0
	for _, line := range normalizedGrid {
		blockCount += strings.Count(line, "#")
	}
	
	if blockCount != 4 {
		return nil, NewParseError(fmt.Sprintf("tetromino must have exactly 4 blocks, got %d", blockCount), lineNumber, filename)
	}

	// Validate that blocks are connected (each block must be adjacent to at least one other)
	if !isConnected(normalizedGrid) {
		return nil, NewParseError("tetromino blocks must be connected", lineNumber, filename)
	}

	// Create tetromino
	tetro, err := tetromino.NewTetromino(id, normalizedGrid)
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("failed to create tetromino: %v", err), lineNumber, filename)
	}

	return tetro, nil
}

// isConnected checks if all blocks in a tetromino are connected
func isConnected(grid []string) bool {
	// Find all block positions
	var blocks []tetromino.Point
	for y, line := range grid {
		for x, char := range line {
			if char == '#' {
				blocks = append(blocks, tetromino.Point{X: x, Y: y})
			}
		}
	}

	if len(blocks) == 0 {
		return false
	}

	// Use flood fill to check connectivity
	visited := make(map[tetromino.Point]bool)
	queue := []tetromino.Point{blocks[0]}
	visited[blocks[0]] = true
	connectedCount := 1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Check all 4 adjacent positions
		directions := []tetromino.Point{
			{X: 0, Y: -1}, // Up
			{X: 0, Y: 1},  // Down
			{X: -1, Y: 0}, // Left
			{X: 1, Y: 0},  // Right
		}

		for _, dir := range directions {
			next := current.Add(dir)
			
			// Check if next position is a block and not visited
			if !visited[next] && isBlockAt(grid, next.X, next.Y) {
				visited[next] = true
				queue = append(queue, next)
				connectedCount++
			}
		}
	}

	// All blocks should be connected
	return connectedCount == len(blocks)
}

// isBlockAt checks if there's a block at the given position
func isBlockAt(grid []string, x, y int) bool {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[y]) {
		return false
	}
	return grid[y][x] == '#'
}

// ValidateFile performs quick validation of file format without full parsing
func ValidateFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return NewParseError(fmt.Sprintf("cannot open file: %v", err), 0, filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	gridLines := 0
	hasContent := false

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			if gridLines > 0 && gridLines != 4 {
				return NewParseError(fmt.Sprintf("incomplete tetromino grid (expected 4 lines, got %d)", gridLines), lineNumber, filename)
			}
			gridLines = 0
			continue
		}

		hasContent = true
		gridLines++

		// Check line length and characters
		trimmed := strings.TrimRight(line, " \t")
		if len(trimmed) > 4 {
			return NewParseError(fmt.Sprintf("line too long: expected max 4 characters, got %d", len(trimmed)), lineNumber, filename)
		}

		// Check for invalid characters
		for i, char := range trimmed {
			if char != '#' && char != '.' {
				return NewParseError(fmt.Sprintf("invalid character '%c' at position %d", char, i), lineNumber, filename)
			}
		}

		// Reset grid line count after complete tetromino
		if gridLines == 4 {
			gridLines = 0
		}
	}

	// Check for incomplete tetromino at end
	if gridLines > 0 && gridLines != 4 {
		return NewParseError(fmt.Sprintf("incomplete tetromino grid at end of file (expected 4 lines, got %d)", gridLines), lineNumber, filename)
	}

	if !hasContent {
		return NewParseError("file is empty or contains only whitespace", 0, filename)
	}

	return scanner.Err()
}
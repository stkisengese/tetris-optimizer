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

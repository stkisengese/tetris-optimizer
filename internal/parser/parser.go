package parser

import (
	"bufio"
	"fmt"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
	"os"
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

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(currentGrid) != 0 {
				tetro, err := processTetromino(currentGrid, currentID, filename)
				if err != nil {
					return nil, err
				}
				tetrominoes = append(tetrominoes, tetro)
				currentGrid = []string{}
				currentID++
			}
			continue
		}

		currentGrid = append(currentGrid, line)
	}

	// Process last tetromino if file doesn't end with a newline
	if len(currentGrid) != 0 {
		tetro, err := processTetromino(currentGrid, currentID, filename)
		if err != nil {
			return nil, err
		}
		tetrominoes = append(tetrominoes, tetro)
	}

	if err := scanner.Err(); err != nil {
		return nil, NewParseError(fmt.Sprintf("error reading file: %v", err), 0, filename)
	}

	if len(tetrominoes) == 0 {
		return nil, NewParseError("no valid tetrominoes found in file", 0, filename)
	}

	return tetrominoes, nil
}

// validateAndCreateTetromino validates a tetromino and creates it
func processTetromino(lines []string, id rune, filename string) (*tetromino.Tetromino, error) {
	if len(lines) != 4 {
		return nil, NewParseError(fmt.Sprintf("tetromino must be 4x4 grid, got %d lines", len(lines)), 0, filename)
	}

	var count int
	var grid [4][4]byte
	startX, startY := -1, -1

	for y, line := range lines {
		if len(line) != 4 {
			return nil, NewParseError(fmt.Sprintf("line %d must be exactly 4 characters, got %d", y+1, len(line)), 0, filename)
		}

		for x := 0; x < 4; x++ {
			ch := line[x]
			if ch != '#' && ch != '.' {
				return nil, NewParseError(fmt.Sprintf("invalid character '%c' at position %d in line %d", ch, x, y+1), 0, filename)
			}
			grid[y][x] = ch
			if ch == '#' {
				count++
				if startX == -1 {
					startX, startY = x, y
				}
			}
		}
	}

	if count != 4 {
		return nil, NewParseError(fmt.Sprintf("tetromino must have exactly 4 blocks, got %d", count), 0, filename)
	}

	if !isConnected(grid, startX, startY) {
		return nil, NewParseError("tetromino blocks must be connected", 0, filename)
	}

	// Create tetromino
	tetro, err := tetromino.NewTetromino(id, lines)
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("failed to create tetromino: %v", err), 0, filename)
	}

	return tetro, nil
}

// isConnected checks if all blocks are connected using DFS
func isConnected(grid [4][4]byte, startX, startY int) bool {
	visited := [4][4]bool{}
	stack := [][2]int{{startY, startX}}
	visited[startY][startX] = true
	count := 1

	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for len(stack) > 0 {
		y, x := stack[len(stack)-1][0], stack[len(stack)-1][1]
		stack = stack[:len(stack)-1]

		for _, d := range dirs {
			ny, nx := y+d[0], x+d[1]
			if ny >= 0 && ny < 4 && nx >= 0 && nx < 4 &&
				grid[ny][nx] == '#' && !visited[ny][nx] {
				visited[ny][nx] = true
				stack = append(stack, [2]int{ny, nx})
				count++
			}
		}
	}
	return count == 4
}

// ValidateFile performs quick validation of file format without full parsing
func ValidateFile(filename string) error {
	_, err := ReadFile(filename)
	return err
}

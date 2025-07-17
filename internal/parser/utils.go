package parser

import (
	"fmt"
	"strings"
)

// FormatGridForDisplay formats a 4x4 grid for display purposes
func FormatGridForDisplay(grid []string) string {
	var builder strings.Builder
	for _, line := range grid {
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	return builder.String()
}

// IsValidTetrominoChar checks if a character is valid for tetromino representation
func IsValidTetrominoChar(char rune) bool {
	return char == '#' || char == '.'
}

// CountBlocks counts the number of blocks (#) in a grid
func CountBlocks(grid []string) int {
	count := 0
	for _, line := range grid {
		count += strings.Count(line, "#")
	}
	return count
}

// NormalizeGrid ensures all lines in the grid are exactly 4 characters
func NormalizeGrid(grid []string) ([]string, error) {
	if len(grid) != 4 {
		return nil, fmt.Errorf("grid must have exactly 4 lines, got %d", len(grid))
	}

	normalized := make([]string, 4)
	for i, line := range grid {
		trimmed := strings.TrimRight(line, " \t")

		if len(trimmed) > 4 {
			return nil, fmt.Errorf("line %d too long: expected max 4 characters, got %d", i+1, len(trimmed))
		}

		// Pad with dots if needed
		for len(trimmed) < 4 {
			trimmed += "."
		}

		normalized[i] = trimmed
	}

	return normalized, nil
}

// GetTetrominoStats returns statistics about a tetromino grid
func GetTetrominoStats(grid []string) (blocks int, minX, maxX, minY, maxY int) {
	blocks = 0
	minX, minY = 4, 4
	maxX, maxY = -1, -1

	for y, line := range grid {
		for x, char := range line {
			if char == '#' {
				blocks++
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	return blocks, minX, maxX, minY, maxY
}

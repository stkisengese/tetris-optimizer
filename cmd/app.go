package main

import (
	"fmt"
	"io"

	"github.com/stkisengese/tetris-optimizer/internal/parser"
	"github.com/stkisengese/tetris-optimizer/internal/solver"
)

// AppResult represents the result of running the application
type AppResult struct {
	Output   string
	ExitCode int
	Error    error
}

// RunApp contains the main application logic, extracted for testing
func RunApp(args []string, writer io.Writer) AppResult {
	if len(args) < 2 {
		fmt.Fprintln(writer, "Usage: go run . <input_file>")
		return AppResult{ExitCode: 1}
	}

	filename := args[1]

	// Parse tetrominoes from file
	tetrominoes, err := parser.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(writer, "ERROR")
		return AppResult{ExitCode: 1, Error: err}
	}

	// Solve the tetris puzzle
	result, err := solver.SolveOptimal(tetrominoes)
	if err != nil {
		fmt.Fprintln(writer, "ERROR")
		return AppResult{ExitCode: 1, Error: err}
	}

	// Check if solution was found
	if !result.Success {
		fmt.Fprintln(writer, "ERROR")
		return AppResult{ExitCode: 1}
	}

	// Print the solution
	output := result.Grid.String()
	fmt.Fprint(writer, output)
	return AppResult{Output: output, ExitCode: 0}
}

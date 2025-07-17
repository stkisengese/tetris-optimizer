package main

import (
	"fmt"
	"github.com/stkisengese/tetris-optimizer/internal/parser"
	"github.com/stkisengese/tetris-optimizer/internal/solver"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <input_file>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Parse tetrominoes from file
	tetrominoes, err := parser.ReadFile(filename)
	if err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}

	// Solve the tetris puzzle
	result, err := solver.SolveOptimal(tetrominoes)
	if err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}

	// Check if solution was found
	if !result.Success {
		fmt.Println("ERROR")
		os.Exit(1)
	}

	// Validate the solution
	if err := solver.ValidateSolution(result.Grid, tetrominoes); err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}

	// Print the solution
	fmt.Print(result.Grid.String())
}

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <input_file>")
		os.Exit(1)
	}
	
	filename := os.Args[1]
	
	fmt.Printf("Input file: %s\n", filename)
	fmt.Println("Full implementation coming in next issues!")
}
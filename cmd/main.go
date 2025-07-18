package main

import (
	"os"
)

func main() {
	result := RunApp(os.Args, os.Stdout)
	if result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}
}

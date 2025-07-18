package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunAppInvalidArgs(t *testing.T) {
	var buf bytes.Buffer

	// Test with no arguments
	result := RunApp([]string{"program"}, &buf)

	if result.ExitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", result.ExitCode)
	}

	if !strings.Contains(buf.String(), "Usage:") {
		t.Errorf("Expected usage message, got: %s", buf.String())
	}
}

func TestRunAppFileNotFound(t *testing.T) {
	var buf bytes.Buffer

	// Test with non-existent file
	result := RunApp([]string{"program", "nonexistent.txt"}, &buf)

	if result.ExitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", result.ExitCode)
	}

	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("Expected ERROR message, got: %s", buf.String())
	}
}

func TestRunAppWithValidFile(t *testing.T) {
	// Create a temporary test file
	content := `....
##..
.#..
.#..

....
####
....
....
`

	tmpFile, err := os.CreateTemp("", "test_tetris_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	var buf bytes.Buffer
	result := RunApp([]string{"program", tmpFile.Name()}, &buf)

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Output: %s", result.ExitCode, buf.String())
	}

	// Check that some output was produced
	if len(result.Output) == 0 && result.ExitCode == 0 {
		t.Error("Expected some output for successful solve")
	}
}

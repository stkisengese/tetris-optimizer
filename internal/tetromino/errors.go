package tetromino

import "fmt"

// TetrominoError represents errors specific to tetromino operations
type TetrominoError struct {
	Operation string
	Message   string
}

func (e *TetrominoError) Error() string {
	return fmt.Sprintf("tetromino %s error: %s", e.Operation, e.Message)
}

// NewTetrominoError creates a new tetromino error
func NewTetrominoError(operation, message string) *TetrominoError {
	return &TetrominoError{
		Operation: operation,
		Message:   message,
	}
}
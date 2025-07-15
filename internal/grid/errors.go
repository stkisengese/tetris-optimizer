package grid

import (
	"fmt"
	"github.com/stkisengese/tetris-optimizer/internal/tetromino"
)

// GridError represents errors specific to grid operations
type GridError struct {
	Operation string
	Message   string
}

func (e *GridError) Error() string {
	return fmt.Sprintf("grid %s error: %s", e.Operation, e.Message)
}

// NewGridError creates a new grid error
func NewGridError(operation, message string) *GridError {
	return &GridError{
		Operation: operation,
		Message:   message,
	}
}

// Common error types
var (
	ErrOutOfBounds    = NewGridError("bounds", "position is out of bounds")
	ErrInvalidSize    = NewGridError("size", "invalid grid size")
	ErrCellOccupied   = NewGridError("placement", "cell is already occupied")
	ErrInvalidFormat  = tetromino.NewTetrominoError("parsing", "invalid tetromino format")
	ErrInvalidBlocks  = tetromino.NewTetrominoError("validation", "tetromino must have exactly 4 blocks")
)
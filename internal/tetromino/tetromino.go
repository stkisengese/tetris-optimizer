package tetromino

import (
	"fmt"
	"sort"
	"strings"
)

// Point represents a coordinate in 2D space
type Point struct {
	X int
	Y int
}

// String returns string representation of Point
func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}


// Add returns a new point that is the sum of this point and another
func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

// Tetromino represents a tetris piece with its shape and position
type Tetromino struct {
	// ID is the identifier for this tetromino (A, B, C, etc.)
	ID rune

	// Points contains the relative coordinates of the tetromino blocks
	// All coordinates are relative to the top-left corner (0,0)
	Points []Point

	// Width and Height represent the bounding box dimensions
	Width  int
	Height int

	// Position represents the current position on the grid
	Position Point
}

// NewTetromino creates a new tetromino from a 4x4 grid representation
func NewTetromino(id rune, grid []string) (*Tetromino, error) {
	if len(grid) != 4 {
		return nil, fmt.Errorf("tetromino grid must be 4x4, got %d rows", len(grid))
	}

	var points []Point
	minX, minY := 4, 4
	maxX, maxY := -1, -1

	// Parse the grid and find all '#' positions
	for y, row := range grid {
		if len(row) != 4 {
			return nil, fmt.Errorf("tetromino row %d must be 4 characters, got %d", y, len(row))
		}

		for x, char := range row {
			if char == '#' {
				points = append(points, Point{X: x, Y: y})
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

	if len(points) != 4 {
		return nil, fmt.Errorf("tetromino must have exactly 4 blocks, got %d", len(points))
	}

	// Normalize coordinates to start from (0,0)
	normalizedPoints := make([]Point, len(points))
	for i, p := range points {
		normalizedPoints[i] = Point{X: p.X - minX, Y: p.Y - minY}
	}

	return &Tetromino{
		ID:       id,
		Points:   normalizedPoints,
		Width:    maxX - minX + 1,
		Height:   maxY - minY + 1,
		Position: Point{X: 0, Y: 0},
	}, nil
}

// Copy creates a deep copy of the tetromino
func (t *Tetromino) Clone() *Tetromino {
	points := make([]Point, len(t.Points))
	copy(points, t.Points)

	return &Tetromino{
		ID:       t.ID,
		Points:   points,
		Width:    t.Width,
		Height:   t.Height,
		Position: t.Position,
	}
}

// SetPosition updates the tetromino's position on the grid
func (t *Tetromino) SetPosition(x, y int) {
	t.Position = Point{X: x, Y: y}
}

// GetAbsolutePoints returns the absolute positions of all blocks
func (t *Tetromino) GetAbsolutePoints() []Point {
	result := make([]Point, len(t.Points))
	for i, p := range t.Points {
		result[i] = p.Add(t.Position)
	}
	return result
}


// Rotate90 rotates the tetromino 90 degrees clockwise
func (t *Tetromino) Rotate90() {
	newPoints := make([]Point, len(t.Points))

	for i, p := range t.Points {
		// Rotation formula: (x, y) -> (y, -x)
		// But we need to adjust for the new bounding box
		newPoints[i] = Point{X: p.Y, Y: -p.X}
	}

	// Normalize the rotated points
	t.Points = t.normalizePoints(newPoints)

	// Swap width and height
	t.Width, t.Height = t.Height, t.Width
}

// GenerateRotations generates all unique rotations of the tetromino
func (t *Tetromino) GenerateRotations() []*Tetromino {
	rotations := make([]*Tetromino, 0, 4)
	current := t.Clone()

	seen := make(map[string]bool)

	for i := 0; i < 4; i++ {
		key := current.ShapeKey()
		if !seen[key] {
			rotations = append(rotations, current.Clone())
			seen[key] = true
		}
		current.Rotate90()
	}

	return rotations
}

// normalizePoints adjusts points so the minimum x and y are 0
func (t *Tetromino) normalizePoints(points []Point) []Point {
	if len(points) == 0 {
		return points
	}

	minX, minY := points[0].X, points[0].Y
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
	}

	normalized := make([]Point, len(points))
	for i, p := range points {
		normalized[i] = Point{X: p.X - minX, Y: p.Y - minY}
	}

	return normalized
}

// shapeKey generates a unique string key for the tetromino shape
func (t *Tetromino) ShapeKey() string {
	// Sort points to ensure consistent ordering
	points := make([]Point, len(t.Points))
	copy(points, t.Points)

	sort.Slice(points, func(i, j int) bool {
		if points[i].Y == points[j].Y {
			return points[i].X < points[j].X
		}
		return points[i].Y < points[j].Y
	})

	var builder strings.Builder
	for i, p := range points {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("%d:%d", p.X, p.Y))
	}

	return builder.String()
}

// String returns a string representation of the tetromino
func (t *Tetromino) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Tetromino %c (%dx%d) at %s:\n",
		t.ID, t.Width, t.Height, t.Position))

	// Create a visual representation
	grid := make([][]rune, t.Height)
	for i := range grid {
		grid[i] = make([]rune, t.Width)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}

	for _, p := range t.Points {
		if p.Y < t.Height && p.X < t.Width {
			grid[p.Y][p.X] = t.ID
		}
	}

	for _, row := range grid {
		builder.WriteString(string(row))
		builder.WriteString("\n")
	}

	return builder.String()
}

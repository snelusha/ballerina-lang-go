package text

import "fmt"

// LinePosition represents a line number and a character offset from the start of the line.
type LinePosition interface {
	Line() int
	Offset() int
	String() string
	LinePositionLookupKey() LinePositionLookupKey
}

// LinePositionLookupKey represents the comparable fields of LinePosition for equality/hashing.
type LinePositionLookupKey struct {
	Line   int
	Offset int
}

// linePositionImpl is the concrete implementation of LinePosition.
type linePositionImpl struct {
	line   int
	offset int
}

// LinePositionFromLineAndOffset constructs a LinePosition with the given line and offset.
func LinePositionFromLineAndOffset(line, offset int) LinePosition {
	return &linePositionImpl{
		line:   line,
		offset: offset,
	}
}

// Line returns the line number.
func (lp linePositionImpl) Line() int {
	return lp.line
}

// Offset returns the character offset from the start of the line.
func (lp linePositionImpl) Offset() int {
	return lp.offset
}

// String returns a string representation of the line position.
func (lp linePositionImpl) String() string {
	return fmt.Sprintf("%d:%d", lp.line, lp.offset)
}

// LinePositionLookupKey returns the lookup key for equality comparisons.
func (lp linePositionImpl) LinePositionLookupKey() LinePositionLookupKey {
	return LinePositionLookupKey{
		Line:   lp.line,
		Offset: lp.offset,
	}
}

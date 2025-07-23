package text

import "fmt"

// LineRange represents a pair of LinePosition.
type LineRange interface {
	FileName() string
	StartLine() LinePosition
	EndLine() LinePosition
	String() string
	LineRangeLookupKey() LineRangeLookupKey
}

// LineRangeLookupKey represents the comparable fields of LineRange for equality/hashing.
type LineRangeLookupKey struct {
	StartLine LinePositionLookupKey
	EndLine   LinePositionLookupKey
}

// lineRangeImpl is the concrete implementation of LineRange.
type lineRangeImpl struct {
	fileName  string
	startLine LinePosition
	endLine   LinePosition
}

// LineRangeFromFileNameAndLinePositions constructs a LineRange with the given file name and line positions.
func LineRangeFromFileNameAndLinePositions(fileName string, startLine, endLine LinePosition) LineRange {
	return &lineRangeImpl{
		fileName:  fileName,
		startLine: startLine,
		endLine:   endLine,
	}
}

// FileName returns the file name.
func (lr lineRangeImpl) FileName() string {
	return lr.fileName
}

// StartLine returns the starting line position.
func (lr lineRangeImpl) StartLine() LinePosition {
	return lr.startLine
}

// EndLine returns the ending line position.
func (lr lineRangeImpl) EndLine() LinePosition {
	return lr.endLine
}

// String returns a string representation of the line range.
func (lr lineRangeImpl) String() string {
	return fmt.Sprintf("(%s,%s)", lr.startLine.String(), lr.endLine.String())
}

// LineRangeLookupKey returns the lookup key for equality comparisons.
func (lr lineRangeImpl) LineRangeLookupKey() LineRangeLookupKey {
	return LineRangeLookupKey{
		StartLine: lr.startLine.LinePositionLookupKey(),
		EndLine:   lr.endLine.LinePositionLookupKey(),
	}
}

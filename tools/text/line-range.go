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

type lineRangeImpl struct {
	fileName  string
	startLine LinePosition
	endLine   LinePosition
}

func LineRangeFromLinePositions(fileName string, startLine, endLine LinePosition) LineRange {
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

func (lr lineRangeImpl) StartLine() LinePosition {
	return lr.startLine
}

func (lr lineRangeImpl) EndLine() LinePosition {
	return lr.endLine
}

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

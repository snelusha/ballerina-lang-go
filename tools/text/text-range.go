package text

import "fmt"

// TextRange describes a contiguous sequence of unicode code points in the TextDocument.
type TextRange interface {
	StartOffset() int
	EndOffset() int
	Length() int
	Contains(position int) bool
	IntersectionExists(textRange TextRange) bool
	String() string
	TextRangeLookupKey() TextRangeLookupKey
}

// TextRangeLookupKey represents the comparable fields of TextRange for equality/hashing.
type TextRangeLookupKey struct {
	StartOffset int
	EndOffset   int
}

type textRangeImpl struct {
	startOffset int
	endOffset   int
	length      int
}

func TextRangeFromStartOffsetAndLength(startOffset, length int) TextRange {
	return &textRangeImpl{
		startOffset: startOffset,
		length:      length,
		endOffset:   startOffset + length,
	}
}

func (tr textRangeImpl) StartOffset() int {
	return tr.startOffset
}

func (tr textRangeImpl) EndOffset() int {
	return tr.endOffset
}

func (tr textRangeImpl) Length() int {
	return tr.length
}

func (tr textRangeImpl) Contains(position int) bool {
	return tr.startOffset <= position && position < tr.endOffset
}

// IntersectionExists tests whether there exists an intersection of this range and the given range.
// The ranges R1(S1, E1) and R2(S2, E2) intersects if S1 is greater than or equal to E2 and
// S2 is less than or equal to E1.
func (tr textRangeImpl) IntersectionExists(textRange TextRange) bool {
	return tr.startOffset <= textRange.EndOffset() && textRange.StartOffset() <= tr.endOffset
}

func (tr textRangeImpl) String() string {
	return fmt.Sprintf("(%d,%d)", tr.startOffset, tr.endOffset)
}

// TextRangeLookupKey returns the lookup key for equality comparisons.
func (tr textRangeImpl) TextRangeLookupKey() TextRangeLookupKey {
	return TextRangeLookupKey{
		StartOffset: tr.startOffset,
		EndOffset:   tr.endOffset,
	}
}

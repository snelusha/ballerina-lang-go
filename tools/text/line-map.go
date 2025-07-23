package text

import (
	"ballerina-lang-go/tools/utils"
	"fmt"
)

// LineMap represents a collection of text lines in the TextDocument.
type LineMap interface {
	TextLine(line int) TextLine
	LinePositionFrom(position int) LinePosition
	TextPositionFrom(linePosition LinePosition) int
	TextLines() []string
}

// lineMapImpl is the concrete implementation of LineMap.
type lineMapImpl struct {
	textLines []TextLine
	length    int
}

// NewLineMap constructs a LineMap with the given text lines.
func NewLineMap(textLines []TextLine) LineMap {
	return &lineMapImpl{
		textLines: textLines,
		length:    len(textLines),
	}
}

// TextLine returns the text line at the given line number.
func (lm lineMapImpl) TextLine(line int) (TextLine, error) {
	if err := lm.lineRangeCheck(line); err != nil {
		return nil, err
	}
	return lm.textLines[line], nil
}

// LinePositionFrom converts a text position to a line position.
func (lm lineMapImpl) LinePositionFrom(position int) (LinePosition, error) {
	if err := lm.positionRangeCheck(position); err != nil {
		return nil, err
	}
	textLine := lm.findLineFrom(position)
	return FromLineAndOffset(textLine.LineNo(), position-textLine.StartOffset()), nil
}

// TextPositionFrom converts a line position to a text position.
func (lm lineMapImpl) TextPositionFrom(linePosition LinePosition) (int, error) {
	if err := lm.lineRangeCheck(linePosition.Line()); err != nil {
		return 0, err
	}
	textLine := lm.textLines[linePosition.Line()]
	if textLine.Length() < linePosition.Offset() {
		return 0, utils.NewIllegalArgumentError(fmt.Sprintf("Cannot find a line with the character offset '%d'", linePosition.Offset()))
	}
	return textLine.StartOffset() + linePosition.Offset(), nil
}

// TextLines returns an immutable list of text content from all lines.
func (lm lineMapImpl) TextLines() []string {
	lines := make([]string, len(lm.textLines))
	for i, textLine := range lm.textLines {
		lines[i] = textLine.Text()
	}
	return lines
}

// positionRangeCheck validates that the position is within bounds.
func (lm lineMapImpl) positionRangeCheck(position int) error {
	if position < 0 || position > lm.textLines[lm.length-1].EndOffset() {
		return utils.NewIndexOutOfBoundsError(fmt.Sprintf("Index: '%d', Size: '%d'", position, lm.textLines[lm.length-1].EndOffset()))
	}
	return nil
}

// lineRangeCheck validates that the line number is within bounds.
func (lm lineMapImpl) lineRangeCheck(lineNo int) error {
	if lineNo < 0 || lineNo > lm.length {
		return utils.NewIndexOutOfBoundsError(fmt.Sprintf("Line number: '%d', Size: '%d'", lineNo, lm.length))
	}
	return nil
}

// findLineFrom returns the TextLine to which the given position belongs.
// Performs a binary search to find the matching text line.
func (lm lineMapImpl) findLineFrom(position int) TextLine {
	// Check boundary conditions
	if position == 0 {
		return lm.textLines[0]
	} else if position == lm.textLines[lm.length-1].EndOffset() {
		return lm.textLines[lm.length-1]
	}

	var foundTextLine TextLine
	left := 0
	right := lm.length - 1

	for left <= right {
		// Using bit shift to handle overflow when sum of left and right is greater than max int
		middle := (left + right) >> 1
		startOffset := lm.textLines[middle].StartOffset()
		endOffset := lm.textLines[middle].EndOffsetWithNewLines()

		if startOffset <= position && position < endOffset {
			foundTextLine = lm.textLines[middle]
			break
		} else if endOffset <= position {
			left = middle + 1
		} else {
			right = middle - 1
		}
	}

	return foundTextLine
}

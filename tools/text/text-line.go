package text

// TextLine represents a single line in the TextDocument.
type TextLine interface {
	LineNo() int
	Text() string
	StartOffset() int
	EndOffset() int
	EndOffsetWithNewLines() int
	Length() int
	LengthWithNewLineChars() int
}

// textLineImpl is the concrete implementation of TextLine.
type textLineImpl struct {
	lineNo               int
	text                 string
	startOffset          int
	endOffset            int
	lengthOfNewLineChars int
}

// NewTextLine constructs a TextLine with the given parameters.
func NewTextLine(lineNo int, text string, startOffset, endOffset, lengthOfNewLineChars int) TextLine {
	return &textLineImpl{
		lineNo:               lineNo,
		text:                 text,
		startOffset:          startOffset,
		endOffset:            endOffset,
		lengthOfNewLineChars: lengthOfNewLineChars,
	}
}

// LineNo returns the line number.
func (tl textLineImpl) LineNo() int {
	return tl.lineNo
}

// Text returns the text content of the line.
func (tl textLineImpl) Text() string {
	return tl.text
}

// StartOffset returns the start offset of the line within the TextDocument.
func (tl textLineImpl) StartOffset() int {
	return tl.startOffset
}

// EndOffset returns the end offset of the line within the TextDocument.
func (tl textLineImpl) EndOffset() int {
	return tl.endOffset
}

// EndOffsetWithNewLines returns the end offset including new line characters.
func (tl textLineImpl) EndOffsetWithNewLines() int {
	return tl.endOffset + tl.lengthOfNewLineChars
}

// Length returns the length of the line without new line characters.
func (tl textLineImpl) Length() int {
	return tl.endOffset - tl.startOffset
}

// LengthWithNewLineChars returns the length of the line including new line characters.
func (tl textLineImpl) LengthWithNewLineChars() int {
	return tl.endOffset - tl.startOffset + tl.lengthOfNewLineChars
}

package bir

type Location interface {
	LineRange() LineRange
	TextRange() TextRange
}

type LineRange interface {
	StartLine() int
	EndLine() int
}

type TextRange interface {
	StartOffset() int
	EndOffset() int
	Length() int
}

type lineRangeImpl struct {
	startLine int
	endLine   int
}

func NewLineRange(startLine, endLine int) LineRange {
	return &lineRangeImpl{
		startLine: startLine,
		endLine:   endLine,
	}
}

func (l *lineRangeImpl) StartLine() int {
	return l.startLine
}

func (l *lineRangeImpl) EndLine() int {
	return l.endLine
}

type textRangeImpl struct {
	startOffset int
	endOffset   int
}

func NewTextRange(startOffset, endOffset int) TextRange {
	return &textRangeImpl{
		startOffset: startOffset,
		endOffset:   endOffset,
	}
}

func (t *textRangeImpl) StartOffset() int {
	return t.startOffset
}

func (t *textRangeImpl) EndOffset() int {
	return t.endOffset
}

func (t *textRangeImpl) Length() int {
	return t.endOffset - t.startOffset
}

type locationImpl struct {
	lineRange LineRange
	textRange TextRange
}

func NewLocation(lineRange LineRange, textRange TextRange) Location {
	return &locationImpl{
		lineRange: lineRange,
		textRange: textRange,
	}
}

func (l *locationImpl) LineRange() LineRange {
	return l.lineRange
}

func (l *locationImpl) TextRange() TextRange {
	return l.textRange
}

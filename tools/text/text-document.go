package text

// TextDocument is an abstract representation of a Ballerina source file (.bal).
type TextDocument interface {
	Apply(textDocumentChange TextDocumentChange) TextDocument
	ToCharArray() []rune
	Line(line int) (TextLine, error)
	LinePositionFromTextPosition(textPosition int) (LinePosition, error)
	TextPositionFromLinePosition(linePosition LinePosition) (int, error)
	TextLines() []string
	Lines() LineMap
	PopulateTextLineMap() LineMap
}

type textDocumentBase struct {
	lineMap LineMap
}

func (td textDocumentBase) Line(line int) (TextLine, error) {
	return td.lineMap.TextLine(line)
}

func (td textDocumentBase) LinePositionFromTextPosition(textPosition int) (LinePosition, error) {
	return td.lineMap.LinePositionFromPosition(textPosition)
}

func (td textDocumentBase) TextPositionFromLinePosition(linePosition LinePosition) (int, error) {
	return td.lineMap.TextPositionFromLinePosition(linePosition)
}

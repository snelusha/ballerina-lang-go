package text

// TextDocument is an abstract representation of a Ballerina source file (.bal).
type TextDocument interface {
	Apply(textDocumentChange TextDocumentChange) TextDocument
	ToCharArray() []rune
	Line(line int) (TextLine, error)
	LinePositionFrom(textPosition int) (LinePosition, error)
	TextPositionFrom(linePosition LinePosition) (int, error)
	TextLines() []string
	Lines() LineMap
	PopulateTextLineMap() LineMap
}

// textDocumentBase holds shared state and default implementations.
type textDocumentBase struct {
	lineMap LineMap
}

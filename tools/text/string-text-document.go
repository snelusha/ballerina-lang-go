package text

import "strings"

// StringTextDocument represents a TextDocument created with a string.
type StringTextDocument interface {
	TextDocument
	String() string
}

// stringTextDocumentImpl is the concrete implementation of StringTextDocument.
type stringTextDocumentImpl struct {
	textDocumentBase
	text        string
	textLineMap LineMap
}

// NewStringTextDocument constructs a StringTextDocument with the given text.
func NewStringTextDocument(text string) StringTextDocument {
	return &stringTextDocumentImpl{
		textDocumentBase: textDocumentBase{},
		text:             text,
	}
}

// Apply implements the abstract method for applying text document changes.
func (std *stringTextDocumentImpl) Apply(textDocumentChange TextDocumentChange) TextDocument {
	startOffset := 0
	var sb strings.Builder
	textEditCount := textDocumentChange.GetTextEditCount()

	for i := range textEditCount {
		textEdit := textDocumentChange.GetTextEdit(i)
		textRange := textEdit.Range()
		sb.WriteString(std.text[startOffset:textRange.StartOffset()])
		sb.WriteString(textEdit.Text())
		startOffset = textRange.EndOffset()
	}
	sb.WriteString(std.text[startOffset:])

	return NewStringTextDocument(sb.String())
}

// PopulateTextLineMap implements the abstract method for populating the text line map.
func (std *stringTextDocumentImpl) PopulateTextLineMap() LineMap {
	if std.textLineMap != nil {
		return std.textLineMap
	}
	std.textLineMap = NewLineMap(std.calculateTextLines())
	return std.textLineMap
}

// ToCharArray implements the abstract method for converting to character array.
func (std *stringTextDocumentImpl) ToCharArray() []rune {
	return []rune(std.text)
}

// String returns the text content as a string.
func (std *stringTextDocumentImpl) String() string {
	return std.text
}

// Line returns the text line at the given line number.
func (std *stringTextDocumentImpl) Line(line int) (TextLine, error) {
	return std.Lines().TextLine(line)
}

// LinePositionFrom converts a text position to a line position.
func (std *stringTextDocumentImpl) LinePositionFrom(textPosition int) (LinePosition, error) {
	return std.Lines().LinePositionFrom(textPosition)
}

// TextPositionFrom converts a line position to a text position.
func (std *stringTextDocumentImpl) TextPositionFrom(linePosition LinePosition) (int, error) {
	return std.Lines().TextPositionFrom(linePosition)
}

// TextLines returns the text content of all lines.
func (std *stringTextDocumentImpl) TextLines() []string {
	if std.lineMap != nil {
		return std.lineMap.TextLines()
	}
	std.lineMap = std.PopulateTextLineMap()
	return std.lineMap.TextLines()
}

// Lines returns the line map, populating it if necessary.
func (std *stringTextDocumentImpl) Lines() LineMap {
	if std.lineMap != nil {
		return std.lineMap
	}
	std.lineMap = std.PopulateTextLineMap()
	return std.lineMap
}

// calculateTextLines parses the text and creates TextLine objects for each line.
func (std *stringTextDocumentImpl) calculateTextLines() []TextLine {
	startOffset := 0
	var textLines []TextLine
	var lineBuilder strings.Builder
	index := 0
	line := 0
	textLength := len(std.text)
	var lengthOfNewLineChars int

	for index < textLength {
		c := rune(std.text[index])
		if c == '\r' || c == '\n' {
			nextCharIndex := index + 1
			if c == '\r' && textLength != nextCharIndex && rune(std.text[nextCharIndex]) == '\n' {
				lengthOfNewLineChars = 2
			} else {
				lengthOfNewLineChars = 1
			}

			strLine := lineBuilder.String()
			endOffset := startOffset + len(strLine)
			textLines = append(textLines, NewTextLine(line, strLine, startOffset, endOffset, lengthOfNewLineChars))
			line++
			startOffset = endOffset + lengthOfNewLineChars
			lineBuilder.Reset()
			index += lengthOfNewLineChars
		} else {
			lineBuilder.WriteRune(c)
			index++
		}
	}

	strLine := lineBuilder.String()
	textLines = append(textLines, NewTextLine(line, strLine, startOffset, startOffset+len(strLine), 0))

	return textLines
}

package text

import "strings"

// StringTextDocument represents a TextDocument created with a string.
type StringTextDocument interface {
	TextDocument
	String() string
}

type stringTextDocumentImpl struct {
	textDocumentBase
	text        string
	textLineMap LineMap
}

func NewStringTextDocument(text string) StringTextDocument {
	return &stringTextDocumentImpl{
		textDocumentBase: textDocumentBase{},
		text:             text,
	}
}

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

func (std *stringTextDocumentImpl) PopulateTextLineMap() LineMap {
	if std.textLineMap != nil {
		return std.textLineMap
	}
	std.textLineMap = NewLineMap(std.calculateTextLines())
	return std.textLineMap
}

func (std *stringTextDocumentImpl) ToCharArray() []rune {
	return []rune(std.text)
}

func (std *stringTextDocumentImpl) String() string {
	return std.text
}

func (std *stringTextDocumentImpl) TextLines() []string {
	if std.textDocumentBase.lineMap != nil {
		return std.textDocumentBase.lineMap.TextLines()
	}
	std.textDocumentBase.lineMap = std.PopulateTextLineMap()
	return std.textDocumentBase.lineMap.TextLines()
}

func (std *stringTextDocumentImpl) Lines() LineMap {
	if std.textDocumentBase.lineMap != nil {
		return std.textDocumentBase.lineMap
	}
	std.textDocumentBase.lineMap = std.PopulateTextLineMap()
	return std.textDocumentBase.lineMap
}

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

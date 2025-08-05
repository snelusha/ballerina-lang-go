package text

import "fmt"

// TextEdit represents a text edit on a TextDocument.
type TextEdit interface {
	Range() TextRange
	Text() string
	String() string
}

type textEditImpl struct {
	textRange TextRange
	text      string
}

func TextEditFromTextRangeAndText(textRange TextRange, text string) TextEdit {
	return &textEditImpl{
		textRange: textRange,
		text:      text,
	}
}

func (te textEditImpl) Range() TextRange {
	return te.textRange
}

func (te textEditImpl) Text() string {
	return te.text
}

func (te textEditImpl) String() string {
	return fmt.Sprintf("%s%s", te.textRange.String(), te.text)
}

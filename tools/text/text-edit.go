package text

import "fmt"

// TextEdit represents a text edit on a TextDocument.
type TextEdit interface {
	Range() TextRange
	Text() string
	String() string
}

// textEditImpl is the concrete implementation of TextEdit.
type textEditImpl struct {
	textRange TextRange
	text      string
}

// TextEditFromTextRangeAndText constructs a TextEdit with the given range and text.
func TextEditFromTextRangeAndText(textRange TextRange, text string) TextEdit {
	return &textEditImpl{
		textRange: textRange,
		text:      text,
	}
}

// Range returns the text range of the edit.
func (te textEditImpl) Range() TextRange {
	return te.textRange
}

// Text returns the replacement text.
func (te textEditImpl) Text() string {
	return te.text
}

// String returns a string representation of the text edit.
func (te textEditImpl) String() string {
	return fmt.Sprintf("%s%s", te.textRange.String(), te.text)
}

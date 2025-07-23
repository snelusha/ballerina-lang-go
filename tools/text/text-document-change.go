package text

import "strings"

// TextDocumentChange represents textual changes on a single TextDocument.
type TextDocumentChange interface {
	GetTextEditCount() int
	GetTextEdit(index int) TextEdit
	String() string
}

// textDocumentChangeImpl is the concrete implementation of TextDocumentChange.
type textDocumentChangeImpl struct {
	textEdits []TextEdit
}

// TextDocumentChangeFromTextEdits constructs a TextDocumentChange with the given text edits.
func TextDocumentChangeFromTextEdits(textEdits []TextEdit) TextDocumentChange {
	// Create a copy of the slice to ensure immutability
	editsCopy := make([]TextEdit, len(textEdits))
	copy(editsCopy, textEdits)

	return &textDocumentChangeImpl{
		textEdits: editsCopy,
	}
}

// GetTextEditCount returns the number of text edits.
func (tdc textDocumentChangeImpl) GetTextEditCount() int {
	return len(tdc.textEdits)
}

// GetTextEdit returns the text edit at the given index.
func (tdc textDocumentChangeImpl) GetTextEdit(index int) TextEdit {
	return tdc.textEdits[index]
}

// String returns a string representation of the text document change.
func (tdc textDocumentChangeImpl) String() string {
	if len(tdc.textEdits) == 0 {
		return ""
	}

	var editStrings []string
	for _, textEdit := range tdc.textEdits {
		editStrings = append(editStrings, textEdit.String())
	}

	return strings.Join(editStrings, ",")
}

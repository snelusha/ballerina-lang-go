package text

import "strings"

// TextDocumentChange represents textual changes on a single TextDocument.
type TextDocumentChange interface {
	GetTextEditCount() int
	GetTextEdit(index int) TextEdit
	String() string
}

type textDocumentChangeImpl struct {
	textEdits []TextEdit
}

func TextDocumentChangeFromTextEdits(textEdits []TextEdit) TextDocumentChange {
	// Create a copy of the slice to ensure immutability
	editsCopy := make([]TextEdit, len(textEdits))
	copy(editsCopy, textEdits)

	return &textDocumentChangeImpl{
		textEdits: editsCopy,
	}
}

func (tdc textDocumentChangeImpl) GetTextEditCount() int {
	return len(tdc.textEdits)
}

func (tdc textDocumentChangeImpl) GetTextEdit(index int) TextEdit {
	return tdc.textEdits[index]
}

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

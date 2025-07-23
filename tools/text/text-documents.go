package text

// TextDocumentFromText creates a TextDocument from the given text string.
func TextDocumentFromText(text string) TextDocument {
	return NewStringTextDocument(text)
}

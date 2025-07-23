package text

import "unicode"

// CharReader is a character reader utility used by the Ballerina lexer.
type CharReader interface {
	Reset(offset int)
	Peek() rune
	PeekN(k int) rune
	Advance()
	AdvanceN(k int)
	Mark()
	GetMarkedChars() string
	IsEOF() bool
}

// charReaderImpl is the concrete implementation of CharReader.
type charReaderImpl struct {
	charBuffer       []rune
	offset           int
	charBufferLength int
	lexemeStartPos   int
}

// newCharReader constructs a CharReader with the given character buffer.
func newCharReader(charBuffer []rune) CharReader {
	return &charReaderImpl{
		charBuffer:       charBuffer,
		offset:           0,
		charBufferLength: len(charBuffer),
		lexemeStartPos:   0,
	}
}

// CharReaderFromTextDocument creates a CharReader from the given TextDocument.
func CharReaderFromTextDocument(textDocument TextDocument) CharReader {
	return newCharReader(textDocument.ToCharArray())
}

// CharReaderFromText creates a CharReader from the given text string.
func CharReaderFromText(text string) CharReader {
	charBuffer := []rune(text)
	return newCharReader(charBuffer)
}

// Reset sets the offset to the given position.
func (cr *charReaderImpl) Reset(offset int) {
	cr.offset = offset
}

// Peek returns the character at the current offset without advancing.
func (cr charReaderImpl) Peek() rune {
	if cr.offset < cr.charBufferLength {
		return cr.charBuffer[cr.offset]
	} else {
		// TODO Revisit this branch
		return unicode.MaxRune
	}
}

// PeekN returns the character at offset + k without advancing.
func (cr charReaderImpl) PeekN(k int) rune {
	n := cr.offset + k
	if n < cr.charBufferLength {
		return cr.charBuffer[n]
	} else {
		// TODO Revisit this branch
		return unicode.MaxRune
	}
}

// Advance moves the offset forward by one position.
func (cr *charReaderImpl) Advance() {
	cr.offset++
}

// AdvanceN moves the offset forward by k positions.
func (cr *charReaderImpl) AdvanceN(k int) {
	cr.offset += k
}

// Mark sets the lexeme start position to the current offset.
func (cr *charReaderImpl) Mark() {
	cr.lexemeStartPos = cr.offset
}

// GetMarkedChars returns the string consisting of the marked characters.
func (cr charReaderImpl) GetMarkedChars() string {
	return string(cr.charBuffer[cr.lexemeStartPos:cr.offset])
}

// IsEOF returns true if the reader has reached the end of the buffer.
func (cr charReaderImpl) IsEOF() bool {
	return cr.offset >= cr.charBufferLength
}

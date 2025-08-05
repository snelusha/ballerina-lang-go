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

func CharReaderFromTextDocument(textDocument TextDocument) CharReader {
	return newCharReader(textDocument.ToCharArray())
}

func CharReaderFromText(text string) CharReader {
	charBuffer := []rune(text)
	return newCharReader(charBuffer)
}

func (cr *charReaderImpl) Reset(offset int) {
	cr.offset = offset
}

func (cr charReaderImpl) Peek() rune {
	if cr.offset < cr.charBufferLength {
		return cr.charBuffer[cr.offset]
	} else {
		// TODO Revisit this branch
		return unicode.MaxRune
	}
}

func (cr charReaderImpl) PeekN(k int) rune {
	n := cr.offset + k
	if n < cr.charBufferLength {
		return cr.charBuffer[n]
	} else {
		// TODO Revisit this branch
		return unicode.MaxRune
	}
}

func (cr *charReaderImpl) Advance() {
	cr.offset++
}

func (cr *charReaderImpl) AdvanceN(k int) {
	cr.offset += k
}

func (cr *charReaderImpl) Mark() {
	cr.lexemeStartPos = cr.offset
}

func (cr charReaderImpl) GetMarkedChars() string {
	return string(cr.charBuffer[cr.lexemeStartPos:cr.offset])
}

func (cr charReaderImpl) IsEOF() bool {
	return cr.offset >= cr.charBufferLength
}

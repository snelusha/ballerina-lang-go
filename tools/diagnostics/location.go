package diagnostics

import (
	"ballerina-lang-go/tools/text"
)

// Location represents the location in TextDocument.
// It is a combination of source file path, start and end line numbers, and start and end column numbers.
type Location interface {
	LineRange() text.LineRange
	TextRange() text.TextRange
}

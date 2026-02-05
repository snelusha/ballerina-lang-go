package parser

import (
	"ballerina-lang-go/parser/common"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/tools/text"
	"fmt"
	"testing"
)

func TestReproLexer(t *testing.T) {
	// The input string mirrors the failing case
	// # + path - file path. Example: `C:\users\OddThinking\Documents\My Source\Widget\foo.src`
	input := "# + path - file path. Example: `C:\\users\\OddThinking\\Documents\\My Source\\Widget\\foo.src`"

	leadingTrivia := make([]tree.STNode, 0)
	diagnostics := make([]tree.STNodeDiagnostic, 0)
	lexer := NewDocumentationLexer(text.CharReaderFromText(input), leadingTrivia, diagnostics, nil)

	for {
		token := lexer.NextToken()
		kind := token.Kind()

		fmt.Printf("Kind: %v, Token: %+v\n", kind, token)

		if kind == common.EOF_TOKEN {
			break
		}
	}
}

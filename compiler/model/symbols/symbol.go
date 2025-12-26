package symbols

import (
	"ballerina-lang-go/compiler/common"
	"ballerina-lang-go/compiler/model/elements"
	"ballerina-lang-go/tools/diagnostics"
)

type Symbol interface {
	GetName() common.Name
	GetOriginalName() common.Name
	GetKind() SymbolKind
	GetType() any
	GetFlags() map[elements.Flag]any
	GetEnclosingSymbol() Symbol
	GetEnclosedSymbols() []any
	GetPosition() diagnostics.Location
	GetOrigin() SymbolOrigin
}

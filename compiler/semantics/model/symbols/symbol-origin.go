package symbols

type SymbolOrigin int

const (
	SymbolOriginVirtual SymbolOrigin = iota
	SymbolOriginSource
	SymbolOriginCompiled
)

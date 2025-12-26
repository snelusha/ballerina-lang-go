package symbols

type SymbolOrigin byte

const (
	SymbolOriginBuiltin        SymbolOrigin = 1
	SymbolOriginSource         SymbolOrigin = 2
	SymbolOriginCompiledSource SymbolOrigin = 3
	SymbolOriginVirtual        SymbolOrigin = 4
)

func (s SymbolOrigin) ToBIROrigin() SymbolOrigin {
	switch s {
	case SymbolOriginBuiltin:
		return SymbolOriginBuiltin
	case SymbolOriginSource:
		return SymbolOriginCompiledSource
	case SymbolOriginCompiledSource:
		return SymbolOriginCompiledSource
	case SymbolOriginVirtual:
		return SymbolOriginVirtual
	default:
		panic("Invalid symbol origin")
	}
}

func (s SymbolOrigin) Value() byte {
	return byte(s)
}

func ToOrigin(value byte) SymbolOrigin {
	switch value {
	case 1:
		return SymbolOriginBuiltin
	case 2:
		return SymbolOriginSource
	case 3:
		return SymbolOriginCompiledSource
	case 4:
		return SymbolOriginVirtual
	default:
		panic("Invalid symbol origin value")
	}
}

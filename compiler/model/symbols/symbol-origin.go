package symbols

type SymbolOrigin byte

const (
	BUILTIN SymbolOrigin = iota + 1
	SOURCE
	COMPILED_SOURCE
	VIRTUAL
)

func (so SymbolOrigin) ToBIROrigin() SymbolOrigin {
	switch so {
	case BUILTIN:
		return BUILTIN
	case SOURCE:
		return SOURCE
	case COMPILED_SOURCE:
		return COMPILED_SOURCE
	case VIRTUAL:
		return VIRTUAL
	default:
		panic("invalid symbol origin")
	}
}

func (so SymbolOrigin) Value() byte {
	return byte(so)
}

func ToOrigin(value byte) SymbolOrigin {
	switch SymbolOrigin(value) {
	case BUILTIN:
		return BUILTIN
	case SOURCE:
		return SOURCE
	case COMPILED_SOURCE:
		return COMPILED_SOURCE
	case VIRTUAL:
		return VIRTUAL
	default:
		panic("invalid symbol origin value")
	}
}

package symbols

type VariableSymbol interface {
	Symbol
	GetConstValue() any
}

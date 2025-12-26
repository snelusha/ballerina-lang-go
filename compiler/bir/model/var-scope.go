package model

type VarScope byte

const (
	VarScopeFunction VarScope = 1
	VarScopeGlobal   VarScope = 2
)

func (v VarScope) GetValue() byte {
	return byte(v)
}

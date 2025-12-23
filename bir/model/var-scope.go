package model

type VarScope byte

const (
	VarScopeFunction VarScope = 1 + iota
	VarScopeGlobal
)

func (v VarScope) Value() byte {
	return byte(v)
}

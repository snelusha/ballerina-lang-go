package model

type VarScope byte

const (
	VarScopeFunction VarScope = 1 + iota
	VarScopeGlobal
)

func (s VarScope) Value() byte {
	return byte(s)
}

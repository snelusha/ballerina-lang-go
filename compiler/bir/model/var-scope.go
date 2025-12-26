package model

type VarScope byte

const (
	VarScopeFunction VarScope = iota + 1
	VarScopeGlobal
)

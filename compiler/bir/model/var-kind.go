package model

type VarKind byte

const (
	VarKindLocal VarKind = 1 + iota
	VarKindArg
	VarKindTemp
	VarKindReturn
	VarKindGlobal
	VarKindSelf
	VarKindConstant
	VarKindSynthetic
)

func (k VarKind) Value() byte {
	return byte(k)
}

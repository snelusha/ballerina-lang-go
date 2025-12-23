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

func (v VarKind) Value() byte {
	return byte(v)
}

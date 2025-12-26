package model

type VarKind byte

const (
	VarKindLocal VarKind = iota + 1
	VarKindArg
	VarKindTemp
	VarKindReturn
	VarKindGlobal
	VarKindSelf
	VarKindConstant
	VarKindSynthetic
)

func (v VarKind) GetValue() byte {
	return byte(v)
}

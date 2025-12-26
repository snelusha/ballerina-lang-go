package model

type VarKind byte

const (
	VarKindLocal     VarKind = 1
	VarKindArg       VarKind = 2
	VarKindTemp      VarKind = 3
	VarKindReturn    VarKind = 4
	VarKindGlobal    VarKind = 5
	VarKindSelf      VarKind = 6
	VarKindConstant  VarKind = 7
	VarKindSynthetic VarKind = 8
)

func (v VarKind) GetValue() byte {
	return byte(v)
}

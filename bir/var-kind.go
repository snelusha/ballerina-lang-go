package bir

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

func (v VarKind) GetValue() byte {
	return byte(v)
}

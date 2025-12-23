package bir

type VarScope byte

const (
	VarScopeFunction VarScope = 1 + iota
	VarScopeGlobal
)

func (v VarScope) GetValue() byte {
	return byte(v)
}

package types

type InvokableType interface {
	Type
	GetParameterTypes() []Type
	GetReturnType() Type
}

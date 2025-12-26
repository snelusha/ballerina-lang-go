package types

type Type interface {
	GetKind() TypeKind
}

type TypeKind int

const (
	TypeKindInt TypeKind = iota
	TypeKindByte
	TypeKindFloat
	TypeKindDecimal
	TypeKindString
	TypeKindBoolean
	TypeKindTypedesc
	TypeKindNil
	TypeKindNever
	TypeKindError
	TypeKindReadonly
	TypeKindParameterized
	TypeKindFunction
	TypeKindOther
)

func (t TypeKind) TypeName() string {
	switch t {
	case TypeKindInt:
		return "int"
	case TypeKindByte:
		return "byte"
	case TypeKindFloat:
		return "float"
	case TypeKindDecimal:
		return "decimal"
	case TypeKindString:
		return "string"
	case TypeKindBoolean:
		return "boolean"
	case TypeKindTypedesc:
		return "typedesc"
	case TypeKindNil:
		return "nil"
	case TypeKindNever:
		return "never"
	case TypeKindError:
		return "error"
	case TypeKindReadonly:
		return "readonly"
	case TypeKindParameterized:
		return "parameterized"
	case TypeKindFunction:
		return "function"
	default:
		return "other"
	}
}

type ValueType interface {
	Type
}

type InvokableType interface {
	Type
	GetParameterTypes() []Type
	GetReturnType() Type
}

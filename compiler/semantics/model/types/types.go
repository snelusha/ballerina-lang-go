package types

import "ballerina-lang-go/compiler/util"

// Type represents a type in Ballerina (base interface)
type Type interface {
	GetKind() TypeKind
}

// TypeKind represents the kind of a type
type TypeKind int

type NamedNode interface {
	GetName() util.Name
}

type BType interface {
	NamedNode
}

type BInvokableType interface {
	BType
}

type bTypeImpl struct {
	name util.Name
}

func (t *bTypeImpl) GetName() util.Name {
	return t.name
}

type bInvokableTypeImpl struct {
	bTypeImpl
	paramTypes []BType
	returnType BType
	restType   BType
}

func NewBInvokableType(paramTypes []BType, restType BType, returnType BType) BInvokableType {
	return &bInvokableTypeImpl{
		paramTypes: paramTypes,
		restType:   restType,
		returnType: returnType,
	}
}

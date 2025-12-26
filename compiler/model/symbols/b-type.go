package symbols

import (
	"ballerina-lang-go/compiler/common"
)

type BType interface {
	GetTag() int
	GetTSymbol() BTypeSymbol
	GetName() common.Name
	GetFlags() int64
	SetFlags(flags int64)
	AddFlags(flags int64)
	GetReturnType() BType
	IsNullable() bool
	GetQualifiedTypeName() string
	String() string
}

type bTypeImpl struct {
	tag     int
	tsymbol BTypeSymbol
	name    common.Name
	flags   int64
}

func NewBType(tag int, tsymbol BTypeSymbol) BType {
	return &bTypeImpl{
		tag:     tag,
		tsymbol: tsymbol,
		name:    common.NewName(""),
		flags:   0,
	}
}

func NewBTypeWithName(tag int, tsymbol BTypeSymbol, name common.Name, flags int64) BType {
	return &bTypeImpl{
		tag:     tag,
		tsymbol: tsymbol,
		name:    name,
		flags:   flags,
	}
}

func (b *bTypeImpl) GetTag() int {
	return b.tag
}

func (b *bTypeImpl) GetTSymbol() BTypeSymbol {
	return b.tsymbol
}

func (b *bTypeImpl) GetName() common.Name {
	return b.name
}

func (b *bTypeImpl) GetFlags() int64 {
	return b.flags
}

func (b *bTypeImpl) SetFlags(flags int64) {
	b.flags = flags
}

func (b *bTypeImpl) AddFlags(flags int64) {
	b.flags |= flags
}

func (b *bTypeImpl) GetReturnType() BType {
	return nil
}

func (b *bTypeImpl) IsNullable() bool {
	return false
}

func (b *bTypeImpl) GetQualifiedTypeName() string {
	if b.tsymbol == nil {
		return ""
	}
	pkgID := b.tsymbol.GetPkgID()
	name := b.tsymbol.GetName()
	if pkgID != nil && name != nil {
		return pkgID.String() + ":" + name.GetValue()
	}
	return ""
}

func (b *bTypeImpl) String() string {
	return "type"
}

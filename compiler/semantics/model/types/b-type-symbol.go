package types

import (
	"fmt"

	"ballerina-lang-go/compiler/common"
	"ballerina-lang-go/compiler/model/elements"
	"ballerina-lang-go/tools/diagnostics"
)

type BTypeSymbol interface {
	BSymbol
	IsTypeParamResolved() bool
	SetTypeParamResolved(resolved bool)
	GetTypeParamTSymbol() BTypeSymbol
	SetTypeParamTSymbol(symbol BTypeSymbol)
	GetAnnotations() BVarSymbol
	SetAnnotations(annotations BVarSymbol)
}

type bTypeSymbolImpl struct {
	*bSymbolImpl
	isTypeParamResolved bool
	typeParamTSymbol    BTypeSymbol
	annotations         BVarSymbol
}

func NewBTypeSymbol(symTag, flags int64, name common.Name, pkgID elements.PackageID, bType BType, owner BSymbol, pos diagnostics.Location, origin SymbolOrigin) BTypeSymbol {
	return NewBTypeSymbolWithOriginalName(symTag, flags, name, name, pkgID, bType, owner, pos, origin)
}

func NewBTypeSymbolWithOriginalName(symTag, flags int64, name, originalName common.Name, pkgID elements.PackageID, bType BType, owner BSymbol, pos diagnostics.Location, origin SymbolOrigin) BTypeSymbol {
	return &bTypeSymbolImpl{
		bSymbolImpl: NewBSymbolWithOriginalName(symTag, flags, name, originalName, pkgID, bType, owner, pos, origin).(*bSymbolImpl),
	}
}

func (b *bTypeSymbolImpl) IsTypeParamResolved() bool {
	return b.isTypeParamResolved
}

func (b *bTypeSymbolImpl) SetTypeParamResolved(resolved bool) {
	b.isTypeParamResolved = resolved
}

func (b *bTypeSymbolImpl) GetTypeParamTSymbol() BTypeSymbol {
	return b.typeParamTSymbol
}

func (b *bTypeSymbolImpl) SetTypeParamTSymbol(symbol BTypeSymbol) {
	b.typeParamTSymbol = symbol
}

func (b *bTypeSymbolImpl) GetAnnotations() BVarSymbol {
	return b.annotations
}

func (b *bTypeSymbolImpl) SetAnnotations(annotations BVarSymbol) {
	b.annotations = annotations
}

func (b *bTypeSymbolImpl) String() string {
	pkgID := b.GetPkgID()
	name := b.GetName()

	if pkgID == nil || pkgID.GetName().GetValue() == "." {
		return name.GetValue()
	}

	return fmt.Sprintf("%s:%s", pkgID.String(), name.GetValue())
}

type BVarSymbol interface {
	BSymbol
}

type bVarSymbolImpl struct {
	*bSymbolImpl
}

func NewBVarSymbol(tag, flags int64, name common.Name, pkgID elements.PackageID, bType BType, owner BSymbol, pos diagnostics.Location, origin SymbolOrigin) BVarSymbol {
	return &bVarSymbolImpl{
		bSymbolImpl: NewBSymbol(tag, flags, name, pkgID, bType, owner, pos, origin).(*bSymbolImpl),
	}
}

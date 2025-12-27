package symbols

import (
	"ballerina-lang-go/compiler/semantics/model/elements"
	"ballerina-lang-go/compiler/semantics/model/types"
	"ballerina-lang-go/compiler/util"
	"ballerina-lang-go/tools/diagnostics"
)

type BTypeSymbol struct {
	*BSymbol
	IsTypeParamResolved bool
	TypeParamTSymbol    *BTypeSymbol
	Annotations         *BVarSymbol
}

func NewBTypeSymbol(symTag, flags int64, name util.Name, pkgID elements.PackageID, typ types.BType, owner *BSymbol, pos diagnostics.Location, origin SymbolOrigin) *BTypeSymbol {
	return NewBTypeSymbolWithOriginalName(symTag, flags, name, name, pkgID, typ, owner, pos, origin)
}

func NewBTypeSymbolWithOriginalName(symTag, flags int64, name, originalName util.Name, pkgID elements.PackageID, typ types.BType, owner *BSymbol, pos diagnostics.Location, origin SymbolOrigin) *BTypeSymbol {
	return &BTypeSymbol{
		BSymbol: NewBSymbolWithOriginalName(symTag, flags, name, originalName, pkgID, typ, owner, pos, origin),
	}
}

func (t *BTypeSymbol) String() string {
	return t.Name.String()
}

type BVarSymbol struct {
	*BSymbol
}

func NewBVarSymbol(tag, flags int64, name util.Name, pkgID elements.PackageID, typ types.BType, owner *BSymbol, location diagnostics.Location, origin SymbolOrigin) *BVarSymbol {
	return &BVarSymbol{
		BSymbol: NewBSymbol(tag, flags, name, pkgID, typ, owner, location, origin),
	}
}

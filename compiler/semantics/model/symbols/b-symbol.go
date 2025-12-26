package symbols

import (
	"ballerina-lang-go/compiler/common"
	"ballerina-lang-go/compiler/model/elements"
	"ballerina-lang-go/compiler/model/symbols"
	semanticTypes "ballerina-lang-go/compiler/semantics/model/types"
	"ballerina-lang-go/compiler/util"
	"ballerina-lang-go/tools/diagnostics"
)

type BSymbol interface {
	symbols.Symbol
	GetTag() uint64
	GetName() common.Name
	GetOriginalName() common.Name
	GetPkgID() elements.PackageID
	GetKind() symbols.SymbolKind
	GetOwner() BSymbol
	IsTainted() bool
	IsClosure() bool
	GetMarkdownDocumentation() elements.MarkdownDocAttachment
	GetOrigin() symbols.SymbolOrigin
	String() string
}

type bSymbolImpl struct {
	tag                   uint64
	flags                 uint64
	name                  common.Name
	originalName          common.Name
	pkgID                 elements.PackageID
	kind                  symbols.SymbolKind
	bType                 semanticTypes.BType
	owner                 BSymbol
	tainted               bool
	closure               bool
	markdownDocumentation elements.MarkdownDocAttachment
	pos                   diagnostics.Location
	origin                symbols.SymbolOrigin
}

func NewBSymbol(tag, flags uint64, name common.Name, pkgID elements.PackageID, bType semanticTypes.BType, owner BSymbol, location diagnostics.Location, origin symbols.SymbolOrigin) BSymbol {
	return NewBSymbolWithOriginalName(tag, flags, name, name, pkgID, bType, owner, location, origin)
}

func NewBSymbolWithOriginalName(tag, flags uint64, name, originalName common.Name, pkgID elements.PackageID, bType semanticTypes.BType, owner BSymbol, location diagnostics.Location, origin symbols.SymbolOrigin) BSymbol {
	return &bSymbolImpl{
		tag:          tag,
		flags:        flags,
		name:         name,
		originalName: originalName,
		pkgID:        pkgID,
		bType:        bType,
		owner:        owner,
		pos:          location,
		origin:       origin,
	}
}

func (b *bSymbolImpl) GetTag() uint64 {
	return b.tag
}

func (b *bSymbolImpl) GetFlags() map[elements.Flag]any {
	return util.UnMask(b.flags)
}

func (b *bSymbolImpl) GetName() common.Name {
	return b.name
}

func (b *bSymbolImpl) GetOriginalName() common.Name {
	if b.originalName != nil && b.originalName.GetValue() != "" {
		return b.originalName
	}
	return b.name
}

func (b *bSymbolImpl) GetPkgID() elements.PackageID {
	return b.pkgID
}

func (b *bSymbolImpl) GetKind() symbols.SymbolKind {
	return symbols.OTHER
}

func (b *bSymbolImpl) GetType() any {
	return b.bType
}

func (b *bSymbolImpl) GetOwner() BSymbol {
	return b.owner
}

func (b *bSymbolImpl) IsTainted() bool {
	return b.tainted
}

func (b *bSymbolImpl) IsClosure() bool {
	return b.closure
}

func (b *bSymbolImpl) GetMarkdownDocumentation() elements.MarkdownDocAttachment {
	return b.markdownDocumentation
}

func (b *bSymbolImpl) GetPosition() diagnostics.Location {
	return b.pos
}

func (b *bSymbolImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *bSymbolImpl) GetEnclosingSymbol() symbols.Symbol {
	return b.owner
}

func (b *bSymbolImpl) GetEnclosedSymbols() []any {
	return make([]any, 0)
}

func (b *bSymbolImpl) String() string {
	return b.name.GetValue()
}

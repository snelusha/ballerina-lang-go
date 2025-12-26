package symbols

import (
	"ballerina-lang-go/compiler/common"
	"ballerina-lang-go/compiler/model/elements"
	"ballerina-lang-go/diagnostics"
)

type BSymbol interface {
	GetTag() int64
	GetFlags() int64
	SetFlags(flags int64)
	GetName() common.Name
	GetOriginalName() common.Name
	GetPkgID() elements.PackageID
	GetKind() SymbolKind
	GetType() BType
	SetType(bType BType)
	GetOwner() BSymbol
	SetOwner(owner BSymbol)
	IsTainted() bool
	SetTainted(tainted bool)
	IsClosure() bool
	SetClosure(closure bool)
	GetMarkdownDocumentation() elements.MarkdownDocAttachment
	SetMarkdownDocumentation(doc elements.MarkdownDocAttachment)
	GetPos() diagnostics.Location
	GetOrigin() SymbolOrigin
	GetScope() Scope
	SetScope(scope Scope)
	GetEnclosedSymbols() []BSymbol
	String() string
}

type bSymbolImpl struct {
	tag                   int64
	flags                 int64
	name                  common.Name
	originalName          common.Name
	pkgID                 elements.PackageID
	kind                  SymbolKind
	bType                 BType
	owner                 BSymbol
	tainted               bool
	closure               bool
	markdownDocumentation elements.MarkdownDocAttachment
	pos                   diagnostics.Location
	origin                SymbolOrigin
	scope                 Scope
}

func NewBSymbol(tag, flags int64, name common.Name, pkgID elements.PackageID, bType BType, owner BSymbol, location diagnostics.Location, origin SymbolOrigin) BSymbol {
	return NewBSymbolWithOriginalName(tag, flags, name, name, pkgID, bType, owner, location, origin)
}

func NewBSymbolWithOriginalName(tag, flags int64, name, originalName common.Name, pkgID elements.PackageID, bType BType, owner BSymbol, location diagnostics.Location, origin SymbolOrigin) BSymbol {
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

func (b *bSymbolImpl) GetTag() int64 {
	return b.tag
}

func (b *bSymbolImpl) GetFlags() int64 {
	return b.flags
}

func (b *bSymbolImpl) SetFlags(flags int64) {
	b.flags = flags
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

func (b *bSymbolImpl) GetKind() SymbolKind {
	return SymbolKindOther
}

func (b *bSymbolImpl) GetType() BType {
	return b.bType
}

func (b *bSymbolImpl) SetType(bType BType) {
	b.bType = bType
}

func (b *bSymbolImpl) GetOwner() BSymbol {
	return b.owner
}

func (b *bSymbolImpl) SetOwner(owner BSymbol) {
	b.owner = owner
}

func (b *bSymbolImpl) IsTainted() bool {
	return b.tainted
}

func (b *bSymbolImpl) SetTainted(tainted bool) {
	b.tainted = tainted
}

func (b *bSymbolImpl) IsClosure() bool {
	return b.closure
}

func (b *bSymbolImpl) SetClosure(closure bool) {
	b.closure = closure
}

func (b *bSymbolImpl) GetMarkdownDocumentation() elements.MarkdownDocAttachment {
	return b.markdownDocumentation
}

func (b *bSymbolImpl) SetMarkdownDocumentation(doc elements.MarkdownDocAttachment) {
	b.markdownDocumentation = doc
}

func (b *bSymbolImpl) GetPos() diagnostics.Location {
	return b.pos
}

func (b *bSymbolImpl) GetOrigin() SymbolOrigin {
	return b.origin
}

func (b *bSymbolImpl) GetScope() Scope {
	return b.scope
}

func (b *bSymbolImpl) SetScope(scope Scope) {
	b.scope = scope
}

func (b *bSymbolImpl) GetEnclosedSymbols() []BSymbol {
	return make([]BSymbol, 0)
}

func (b *bSymbolImpl) String() string {
	return b.name.GetValue()
}

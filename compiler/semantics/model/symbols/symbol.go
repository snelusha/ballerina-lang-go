package symbols

import (
	"ballerina-lang-go/compiler/semantics/model/elements"
	"ballerina-lang-go/compiler/semantics/model/types"
	"ballerina-lang-go/compiler/util"
	"ballerina-lang-go/tools/diagnostics"
)

type Symbol interface {
	GetName() util.Name
	GetOriginalName() util.Name
	GetKind() SymbolKind
	GetType() types.BType
	GetFlags() map[elements.Flag]struct{}
	GetEnclosingSymbol() Symbol
	GetEnclosedSymbols() []Symbol
	GetPosition() diagnostics.Location
	GetOrigin() SymbolOrigin
}

type BSymbol struct {
	Tag                   int64
	Flags                 int64
	Name                  util.Name
	OriginalName          util.Name
	PkgID                 elements.PackageID
	Kind                  SymbolKind
	Type                  types.BType
	Owner                 *BSymbol
	Tainted               bool
	Closure               bool
	MarkdownDocumentation *elements.MarkdownDocAttachment
	Pos                   diagnostics.Location
	Origin                SymbolOrigin
	Scope                 *Scope
}

func NewBSymbol(tag, flags int64, name util.Name, pkgID elements.PackageID, typ types.BType, owner *BSymbol, location diagnostics.Location, origin SymbolOrigin) *BSymbol {
	return &BSymbol{
		Tag:          tag,
		Flags:        flags,
		Name:         name,
		OriginalName: name,
		PkgID:        pkgID,
		Type:         typ,
		Owner:        owner,
		Pos:          location,
		Origin:       origin,
	}
}

func NewBSymbolWithOriginalName(tag, flags int64, name, originalName util.Name, pkgID elements.PackageID, typ types.BType, owner *BSymbol, location diagnostics.Location, origin SymbolOrigin) *BSymbol {
	return &BSymbol{
		Tag:          tag,
		Flags:        flags,
		Name:         name,
		OriginalName: originalName,
		PkgID:        pkgID,
		Type:         typ,
		Owner:        owner,
		Pos:          location,
		Origin:       origin,
	}
}

func (s *BSymbol) GetMarkdownDocAttachment() *elements.MarkdownDocAttachment {
	return s.MarkdownDocumentation
}

func (s *BSymbol) GetName() util.Name {
	return s.Name
}

func (s *BSymbol) GetOriginalName() util.Name {
	if s.OriginalName.Value != "" {
		return s.OriginalName
	}
	return s.Name
}

func (s *BSymbol) GetKind() SymbolKind {
	return SymbolKindOther
}

func (s *BSymbol) GetType() types.BType {
	return s.Type
}

func (s *BSymbol) GetFlags() map[elements.Flag]struct{} {
	return elements.UnmaskFlags(s.Flags)
}

func (s *BSymbol) GetEnclosingSymbol() Symbol {
	if s.Owner != nil {
		return s.Owner
	}
	return nil
}

func (s *BSymbol) GetEnclosedSymbols() []Symbol {
	return make([]Symbol, 0)
}

func (s *BSymbol) GetPosition() diagnostics.Location {
	return s.Pos
}

func (s *BSymbol) GetOrigin() SymbolOrigin {
	return s.Origin
}

func (s *BSymbol) String() string {
	return s.Name.String()
}

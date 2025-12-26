package model

import "ballerina-lang-go/compiler/common"

type ScopeEntry interface {
	GetSymbol() BSymbol
	GetNext() ScopeEntry
}

type scopeEntryImpl struct {
	symbol BSymbol
	next   ScopeEntry
}

func NewScopeEntry(symbol BSymbol, next ScopeEntry) ScopeEntry {
	return &scopeEntryImpl{
		symbol: symbol,
		next:   next,
	}
}

func (se *scopeEntryImpl) GetSymbol() BSymbol {
	return se.symbol
}

func (se *scopeEntryImpl) GetNext() ScopeEntry {
	return se.next
}

type Scope interface {
	GetOwner() BSymbol
	GetEntries() map[string]ScopeEntry
	Define(name common.Name, symbol BSymbol)
	Lookup(name common.Name) ScopeEntry
}

type scopeImpl struct {
	owner   BSymbol
	entries map[string]ScopeEntry
}

var NotFoundEntry = NewScopeEntry(nil, nil)

func NewScope(owner BSymbol) Scope {
	return &scopeImpl{
		owner:   owner,
		entries: make(map[string]ScopeEntry),
	}
}

func (s *scopeImpl) GetOwner() BSymbol {
	return s.owner
}

func (s *scopeImpl) GetEntries() map[string]ScopeEntry {
	return s.entries
}

func (s *scopeImpl) Define(name common.Name, symbol BSymbol) {
	entry := NewScopeEntry(symbol, s.entries[name.GetValue()])
	s.entries[name.GetValue()] = entry
}

func (s *scopeImpl) Lookup(name common.Name) ScopeEntry {
	if entry, exists := s.entries[name.GetValue()]; exists {
		return entry
	}
	return NotFoundEntry
}

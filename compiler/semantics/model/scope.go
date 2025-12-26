package model

import (
	"ballerina-lang-go/compiler/common"
	"ballerina-lang-go/compiler/model/symbols"
)

type ScopeEntry interface {
	GetSymbol() symbols.BSymbol
	GetNext() ScopeEntry
}

type scopeEntryImpl struct {
	symbol symbols.BSymbol
	next   ScopeEntry
}

func NewScopeEntry(symbol symbols.BSymbol, next ScopeEntry) ScopeEntry {
	return &scopeEntryImpl{
		symbol: symbol,
		next:   next,
	}
}

func (s *scopeEntryImpl) GetSymbol() symbols.BSymbol {
	return s.symbol
}

func (s *scopeEntryImpl) GetNext() ScopeEntry {
	return s.next
}

type Scope interface {
	GetOwner() symbols.BSymbol
	GetEntries() map[string]ScopeEntry
	Define(name common.Name, symbol symbols.BSymbol)
	Lookup(name common.Name) ScopeEntry
}

type scopeImpl struct {
	owner   symbols.BSymbol
	entries map[string]ScopeEntry
}

var NotFoundEntry = NewScopeEntry(nil, nil)

func NewScope(owner symbols.BSymbol) Scope {
	return &scopeImpl{
		owner:   owner,
		entries: make(map[string]ScopeEntry, 10),
	}
}

func (s *scopeImpl) GetOwner() symbols.BSymbol {
	return s.owner
}

func (s *scopeImpl) GetEntries() map[string]ScopeEntry {
	return s.entries
}

func (s *scopeImpl) Define(name common.Name, symbol symbols.BSymbol) {
	current := s.entries[name.GetValue()]
	if current == nil {
		current = NotFoundEntry
	}

	newEntry := NewScopeEntry(symbol, current)
	s.entries[name.GetValue()] = newEntry
}

func (s *scopeImpl) Lookup(name common.Name) ScopeEntry {
	entry := s.entries[name.GetValue()]
	if entry == nil {
		return NotFoundEntry
	}
	return entry
}

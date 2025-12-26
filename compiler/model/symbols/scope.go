package symbols

import (
	"ballerina-lang-go/compiler/common"
)

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

func (s *scopeEntryImpl) GetSymbol() BSymbol {
	return s.symbol
}

func (s *scopeEntryImpl) GetNext() ScopeEntry {
	return s.next
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
		entries: make(map[string]ScopeEntry, 10),
	}
}

func (s *scopeImpl) GetOwner() BSymbol {
	return s.owner
}

func (s *scopeImpl) GetEntries() map[string]ScopeEntry {
	return s.entries
}

func (s *scopeImpl) Define(name common.Name, symbol BSymbol) {
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

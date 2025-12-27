package symbols

import (
	"ballerina-lang-go/compiler/util"
)

type Scope struct {
	Owner   *BSymbol
	Entries map[util.Name]*ScopeEntry
}

type ScopeEntry struct {
	Symbol *BSymbol
	Next   *ScopeEntry
}

var NotFoundEntry = &ScopeEntry{Symbol: nil, Next: nil}

func NewScope(owner *BSymbol) *Scope {
	return &Scope{
		Owner:   owner,
		Entries: make(map[util.Name]*ScopeEntry, 10),
	}
}

func (s *Scope) Define(name util.Name, symbol *BSymbol) {
	current := s.Entries[name]
	if current == nil {
		current = NotFoundEntry
	}

	newEntry := &ScopeEntry{
		Symbol: symbol,
		Next:   current,
	}
	s.Entries[name] = newEntry
}

func (s *Scope) Lookup(name util.Name) *ScopeEntry {
	entry := s.Entries[name]
	if entry == nil {
		return NotFoundEntry
	}
	return entry
}

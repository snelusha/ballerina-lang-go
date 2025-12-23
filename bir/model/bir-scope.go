package model

type BirScope struct {
	ID     int
	Parent *BirScope
}

func NewBirScope(id int, parent *BirScope) *BirScope {
	return &BirScope{
		ID:     id,
		Parent: parent,
	}
}

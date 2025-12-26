package model

type BirScope interface {
	GetID() int
	GetParent() BirScope
}

type birScopeImpl struct {
	id     int
	parent BirScope
}

func NewBirScope(id int, parent BirScope) BirScope {
	return &birScopeImpl{
		id:     id,
		parent: parent,
	}
}

func (b *birScopeImpl) GetID() int {
	return b.id
}

func (b *birScopeImpl) GetParent() BirScope {
	return b.parent
}

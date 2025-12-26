package common

type Name interface {
	GetValue() string
}

type nameImpl struct {
	value string
}

func NewName(value string) Name {
	return &nameImpl{value: value}
}

func (n *nameImpl) GetValue() string {
	return n.value
}

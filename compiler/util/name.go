package util

type Name struct {
	Value string
}

func NewName(value string) *Name {
	return &Name{Value: value}
}

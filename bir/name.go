package bir

type Name struct {
	Value string
}

func NewName(value string) *Name {
	return &Name{Value: value}
}

func (n *Name) GetValue() string {
	return n.Value
}

func (n *Name) String() string {
	return n.Value
}

func (n *Name) Equals(other *Name) bool {
	if n == other {
		return true
	}
	if other == nil {
		return false
	}
	return n.Value == other.Value
}

package bir

type BIROperand struct {
	Pos         Location
	VariableDcl *BIRVariableDcl
}

func NewBIROperand(variableDcl *BIRVariableDcl) *BIROperand {
	return &BIROperand{
		Pos:         nil,
		VariableDcl: variableDcl,
	}
}

func (o *BIROperand) Accept(visitor BIRVisitor) {
	visitor.VisitOperand(o)
}

func (o *BIROperand) Equals(other *BIROperand) bool {
	if o == other {
		return true
	}
	if other == nil {
		return false
	}
	return o.VariableDcl.Equals(other.VariableDcl)
}

func (o *BIROperand) String() string {
	return o.VariableDcl.String()
}

package model

type BIROperand interface {
	BIRNode
	GetVariableDcl() BIRVariableDcl
}

type birOperandImpl struct {
	variableDcl BIRVariableDcl
}

func NewBIROperand(variableDcl BIRVariableDcl) BIROperand {
	return &birOperandImpl{
		variableDcl: variableDcl,
	}
}

func (b *birOperandImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIROperand(b)
}

func (b *birOperandImpl) GetVariableDcl() BIRVariableDcl {
	return b.variableDcl
}

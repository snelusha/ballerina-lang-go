package model

type BIROperand interface {
	BIRNode
	GetVariableDcl() BIRVariableDcl
	String() string
}

type birOperandImpl struct {
	*birNodeImpl
	variableDcl BIRVariableDcl
}

func NewBIROperand(variableDcl BIRVariableDcl) BIROperand {
	return &birOperandImpl{
		birNodeImpl: newBIRNode(nil),
		variableDcl: variableDcl,
	}
}

func (b *birOperandImpl) GetVariableDcl() BIRVariableDcl {
	return b.variableDcl
}

func (b *birOperandImpl) String() string {
	return b.variableDcl.String()
}

func (b *birOperandImpl) Accept(visitor BIRVisitor) {
	visitor.VisitOperand(b)
}

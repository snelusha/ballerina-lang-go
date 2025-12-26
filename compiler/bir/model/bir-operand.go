package model

type BIROperand interface {
	BIRNode
	GetVariableDcl() BIRVariableDcl
}

type birOperandImpl struct {
	*birNodeImpl
	variableDcl BIRVariableDcl
}

func NewBIROperand(variableDcl BIRVariableDcl) BIROperand {
	return &birOperandImpl{
		birNodeImpl: NewBIRNode(nil).(*birNodeImpl),
		variableDcl: variableDcl,
	}
}

func (b *birOperandImpl) GetVariableDcl() BIRVariableDcl {
	return b.variableDcl
}

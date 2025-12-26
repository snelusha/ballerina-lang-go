package model

type BIRAssignInstruction interface {
	BIRInstruction
	GetLhsOperand() BIROperand
}

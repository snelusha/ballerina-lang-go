package model

import "ballerina-lang-go/tools/diagnostics"

type BIRAbstractInstruction interface {
	BIRInstruction
	GetKind() InstructionKind
	GetRhsOperands() []BIROperand
	SetRhsOperands(operands []BIROperand)
}

type birAbstractInstructionImpl struct {
	pos   diagnostics.Location
	kind  InstructionKind
	lhsOp BIROperand
	scope *BirScope
}

func newBIRAbstractInstruction(pos diagnostics.Location, kind InstructionKind) *birAbstractInstructionImpl {
	return &birAbstractInstructionImpl{
		pos:  pos,
		kind: kind,
	}
}

func (b *birAbstractInstructionImpl) GetKind() InstructionKind {
	return b.kind
}

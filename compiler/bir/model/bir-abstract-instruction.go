package model

import (
	"ballerina-lang-go/diagnostics"
)

type BIRAbstractInstruction interface {
	BIRNode
	BIRInstruction
	GetLhsOp() BIROperand
	SetLhsOp(op BIROperand)
	GetScope() BirScope
	SetScope(scope BirScope)
	GetRhsOperands() []BIROperand
	SetRhsOperands(operands []BIROperand)
}

type birAbstractInstructionImpl struct {
	*birNodeImpl
	kind  InstructionKind
	lhsOp BIROperand
	scope BirScope
}

func newBIRAbstractInstruction(pos diagnostics.Location, kind InstructionKind) *birAbstractInstructionImpl {
	return &birAbstractInstructionImpl{
		birNodeImpl: newBIRNode(pos),
		kind:        kind,
	}
}

func (b *birAbstractInstructionImpl) GetKind() InstructionKind {
	return b.kind
}

func (b *birAbstractInstructionImpl) GetLhsOp() BIROperand {
	return b.lhsOp
}

func (b *birAbstractInstructionImpl) SetLhsOp(op BIROperand) {
	b.lhsOp = op
}

func (b *birAbstractInstructionImpl) GetScope() BirScope {
	return b.scope
}

func (b *birAbstractInstructionImpl) SetScope(scope BirScope) {
	b.scope = scope
}

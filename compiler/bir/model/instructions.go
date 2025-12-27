package model

import (
	"ballerina-lang-go/tools/diagnostics"
)

type BIRInstruction interface {
	GetKind() InstructionKind
}

type BIRAbstractInstruction interface {
	BIRNode
	BIRInstruction
	GetLhsOp() *BIROperand
	GetRhsOperands() []*BIROperand
	SetRhsOperands(operands []*BIROperand)
	GetScope() *BirScope
	SetScope(scope *BirScope)
}

type birAbstractInstructionBase struct {
	birNodeBase
	kind  InstructionKind
	lhsOp *BIROperand
	scope *BirScope
}

func (b *birAbstractInstructionBase) GetKind() InstructionKind {
	return b.kind
}

func (b *birAbstractInstructionBase) GetLhsOp() *BIROperand {
	return b.lhsOp
}

func (b *birAbstractInstructionBase) GetScope() *BirScope {
	return b.scope
}

func (b *birAbstractInstructionBase) SetScope(scope *BirScope) {
	b.scope = scope
}

type BIRAssignInstruction interface {
	BIRInstruction
	GetLhsOperand() *BIROperand
}

type BIRNonTerminator interface {
	BIRAbstractInstruction
}

type birNonTerminatorBase struct {
	birAbstractInstructionBase
}

func (b *birNonTerminatorBase) GetRhsOperands() []*BIROperand {
	return []*BIROperand{}
}

func (b *birNonTerminatorBase) SetRhsOperands(operands []*BIROperand) {
}

type BIRTerminator interface {
	BIRAbstractInstruction
	GetNextBasicBlocks() []*BIRBasicBlock
	GetThenBB() *BIRBasicBlock
	SetThenBB(bb *BIRBasicBlock)
}

type birTerminatorBase struct {
	birAbstractInstructionBase
	thenBB *BIRBasicBlock
}

func (b *birTerminatorBase) GetThenBB() *BIRBasicBlock {
	return b.thenBB
}

func (b *birTerminatorBase) SetThenBB(bb *BIRBasicBlock) {
	b.thenBB = bb
}

func (b *birTerminatorBase) GetRhsOperands() []*BIROperand {
	return []*BIROperand{}
}

func (b *birTerminatorBase) SetRhsOperands(operands []*BIROperand) {
}

type BIROperandImpl struct {
	birNodeBase
	VariableDcl *BIRVariableDcl
}

func NewBIROperand(variableDcl *BIRVariableDcl) *BIROperandImpl {
	return &BIROperandImpl{
		birNodeBase: birNodeBase{pos: nil},
		VariableDcl: variableDcl,
	}
}

func (o *BIROperandImpl) String() string {
	if o.VariableDcl != nil {
		return o.VariableDcl.String()
	}
	return ""
}

type BIRMove struct {
	birNonTerminatorBase
	RhsOp *BIROperand
}

func NewBIRMove(pos diagnostics.Location, fromOperand, toOperand *BIROperand) *BIRMove {
	if toOperand != nil && toOperand.VariableDcl != nil {
		toOperand.VariableDcl.Initialized = true
	}
	return &BIRMove{
		birNonTerminatorBase: birNonTerminatorBase{
			birAbstractInstructionBase: birAbstractInstructionBase{
				birNodeBase: birNodeBase{pos: pos},
				kind:        InstructionKindMove,
				lhsOp:       toOperand,
			},
		},
		RhsOp: fromOperand,
	}
}

func (m *BIRMove) GetLhsOperand() *BIROperand {
	return m.lhsOp
}

func (m *BIRMove) GetRhsOperands() []*BIROperand {
	return []*BIROperand{m.RhsOp}
}

func (m *BIRMove) SetRhsOperands(operands []*BIROperand) {
	if len(operands) > 0 {
		m.RhsOp = operands[0]
	}
}

type BIRBinaryOp struct {
	birNonTerminatorBase
	RhsOp1 *BIROperand
	RhsOp2 *BIROperand
}

func NewBIRBinaryOp(pos diagnostics.Location, kind InstructionKind, lhsOp, rhsOp1, rhsOp2 *BIROperand) *BIRBinaryOp {
	return &BIRBinaryOp{
		birNonTerminatorBase: birNonTerminatorBase{
			birAbstractInstructionBase: birAbstractInstructionBase{
				birNodeBase: birNodeBase{pos: pos},
				kind:        kind,
				lhsOp:       lhsOp,
			},
		},
		RhsOp1: rhsOp1,
		RhsOp2: rhsOp2,
	}
}

func (b *BIRBinaryOp) GetLhsOperand() *BIROperand {
	return b.lhsOp
}

func (b *BIRBinaryOp) GetRhsOperands() []*BIROperand {
	return []*BIROperand{b.RhsOp1, b.RhsOp2}
}

func (b *BIRBinaryOp) SetRhsOperands(operands []*BIROperand) {
	if len(operands) >= 2 {
		b.RhsOp1 = operands[0]
		b.RhsOp2 = operands[1]
	}
}

type BIRUnaryOp struct {
	birNonTerminatorBase
	RhsOp *BIROperand
}

func NewBIRUnaryOp(pos diagnostics.Location, kind InstructionKind, lhsOp, rhsOp *BIROperand) *BIRUnaryOp {
	return &BIRUnaryOp{
		birNonTerminatorBase: birNonTerminatorBase{
			birAbstractInstructionBase: birAbstractInstructionBase{
				birNodeBase: birNodeBase{pos: pos},
				kind:        kind,
				lhsOp:       lhsOp,
			},
		},
		RhsOp: rhsOp,
	}
}

func (u *BIRUnaryOp) GetLhsOperand() *BIROperand {
	return u.lhsOp
}

func (u *BIRUnaryOp) GetRhsOperands() []*BIROperand {
	return []*BIROperand{u.RhsOp}
}

func (u *BIRUnaryOp) SetRhsOperands(operands []*BIROperand) {
	if len(operands) > 0 {
		u.RhsOp = operands[0]
	}
}

type BIRGoto struct {
	birTerminatorBase
	TargetBB *BIRBasicBlock
}

func NewBIRGoto(pos diagnostics.Location, targetBB *BIRBasicBlock) *BIRGoto {
	return &BIRGoto{
		birTerminatorBase: birTerminatorBase{
			birAbstractInstructionBase: birAbstractInstructionBase{
				birNodeBase: birNodeBase{pos: pos},
				kind:        InstructionKindGoto,
			},
		},
		TargetBB: targetBB,
	}
}

func NewBIRGotoWithScope(pos diagnostics.Location, targetBB *BIRBasicBlock, scope *BirScope) *BIRGoto {
	g := NewBIRGoto(pos, targetBB)
	g.scope = scope
	return g
}

func (g *BIRGoto) GetNextBasicBlocks() []*BIRBasicBlock {
	return []*BIRBasicBlock{g.TargetBB}
}

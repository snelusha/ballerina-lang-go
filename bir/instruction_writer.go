// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package bir

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"ballerina-lang-go/tools/diagnostics"
)

// BIRInstructionWriter writes BIR instructions to buffers
type BIRInstructionWriter struct {
	buf               *bytes.Buffer
	scopeBuf          *bytes.Buffer
	cp                *ConstantPool
	instructionOffset int
	completedScopes   map[*BIRScope]bool
	scopeCount        int
}

// NewBIRInstructionWriter creates a new instruction writer
func NewBIRInstructionWriter(buf *bytes.Buffer, scopeBuf *bytes.Buffer, cp *ConstantPool) *BIRInstructionWriter {
	return &BIRInstructionWriter{
		buf:               buf,
		scopeBuf:          scopeBuf,
		cp:                cp,
		instructionOffset: 0,
		completedScopes:   make(map[*BIRScope]bool),
		scopeCount:        0,
	}
}

// GetScopeCount returns the number of scopes written
func (iw *BIRInstructionWriter) GetScopeCount() int {
	return iw.scopeCount
}

// WriteBBs writes basic blocks
func (iw *BIRInstructionWriter) WriteBBs(bbList []BIRBasicBlock) error {
	if err := binary.Write(iw.buf, binary.BigEndian, int32(len(bbList))); err != nil {
		return err
	}

	for _, bb := range bbList {
		if err := iw.writeBB(bb); err != nil {
			return err
		}
	}
	return nil
}

func (iw *BIRInstructionWriter) writeBB(bb BIRBasicBlock) error {
	// Write basic block name
	bbNameCPIndex := AddStringCPEntry(bb.Id.Value(), iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(bbNameCPIndex)); err != nil {
		return err
	}

	// Write instruction count (including terminator)
	instructionCount := len(bb.Instructions) + 1
	if err := binary.Write(iw.buf, binary.BigEndian, int32(instructionCount)); err != nil {
		return err
	}

	// Write non-terminator instructions
	for _, instruction := range bb.Instructions {
		if err := iw.writeInstruction(instruction); err != nil {
			return err
		}
		iw.instructionOffset++
		iw.writeScopes(instruction)
	}

	// Write terminator
	if bb.Terminator == nil {
		return fmt.Errorf("basic block without terminator: %s", bb.Id.Value())
	}

	if err := iw.writeTerminator(bb.Terminator); err != nil {
		return err
	}
	iw.writeScope(bb.Terminator)

	return nil
}

func (iw *BIRInstructionWriter) writeInstruction(instruction BIRNonTerminator) error {
	// Write position - get from instruction base
	pos := iw.getInstructionPos(instruction)
	if err := WritePosition(pos, iw.buf, iw.cp); err != nil {
		return err
	}

	// Write instruction kind
	kind := instruction.GetKind()
	if err := binary.Write(iw.buf, binary.BigEndian, byte(kind)); err != nil {
		return err
	}

	// Write instruction-specific data
	return iw.writeInstructionData(instruction, kind)
}

func (iw *BIRInstructionWriter) writeInstructionData(instruction BIRNonTerminator, _ InstructionKind) error {
	// This is a simplified implementation
	// Full implementation would handle all instruction types
	switch inst := instruction.(type) {
	case *Move:
		return iw.writeMove(inst)
	case *ConstantLoad:
		return iw.writeConstantLoad(inst)
	case *BinaryOp:
		return iw.writeBinaryOp(inst)
	case *UnaryOp:
		return iw.writeUnaryOp(inst)
	default:
		// For other instruction types, write minimal data
		return nil
	}
}

func (iw *BIRInstructionWriter) writeMove(move *Move) error {
	if err := iw.writeOperand(move.RhsOp); err != nil {
		return err
	}
	return iw.writeOperand(move.LhsOp)
}

func (iw *BIRInstructionWriter) writeConstantLoad(cl *ConstantLoad) error {
	if err := WriteType(iw.cp, iw.buf, cl.Type); err != nil {
		return fmt.Errorf("writing constant load type: %w", err)
	}
	if err := iw.writeOperand(cl.LhsOp); err != nil {
		return fmt.Errorf("writing constant load lhs operand: %w", err)
	}

	// Write constant value based on type
	if err := WriteConstValueWithType(iw.cp, iw.buf, cl.Value, cl.Type); err != nil {
		return fmt.Errorf("writing constant load value: %w", err)
	}
	return nil
}

func (iw *BIRInstructionWriter) writeBinaryOp(binOp *BinaryOp) error {
	if err := iw.writeOperand(&binOp.RhsOp1); err != nil {
		return err
	}
	if err := iw.writeOperand(&binOp.RhsOp2); err != nil {
		return err
	}
	return iw.writeOperand(binOp.LhsOp)
}

func (iw *BIRInstructionWriter) writeUnaryOp(unaryOp *UnaryOp) error {
	if err := iw.writeOperand(unaryOp.RhsOp); err != nil {
		return err
	}
	return iw.writeOperand(unaryOp.LhsOp)
}

func (iw *BIRInstructionWriter) writeTerminator(terminator BIRTerminator) error {
	// Write position - get from terminator base
	pos := iw.getTerminatorPos(terminator)
	if err := WritePosition(pos, iw.buf, iw.cp); err != nil {
		return err
	}

	// Write terminator kind
	kind := terminator.GetKind()
	if err := binary.Write(iw.buf, binary.BigEndian, byte(kind)); err != nil {
		return err
	}

	// Write terminator-specific data
	switch term := terminator.(type) {
	case *Goto:
		return iw.writeGoto(term)
	case *Return:
		return nil // Return has no additional data
	case *Branch:
		return iw.writeBranch(term)
	case *Call:
		return iw.writeCall(term)
	default:
		return fmt.Errorf("unsupported terminator type: %T", terminator)
	}
}

func (iw *BIRInstructionWriter) writeGoto(gotoInst *Goto) error {
	if gotoInst.ThenBB == nil {
		return fmt.Errorf("goto terminator has nil ThenBB")
	}
	bbNameCPIndex := AddStringCPEntry(gotoInst.ThenBB.Id.Value(), iw.cp)
	return binary.Write(iw.buf, binary.BigEndian, int32(bbNameCPIndex))
}

func (iw *BIRInstructionWriter) writeBranch(branch *Branch) error {
	if err := iw.writeOperand(branch.Op); err != nil {
		return err
	}

	if branch.TrueBB == nil {
		return fmt.Errorf("branch terminator has nil TrueBB")
	}
	trueBBNameCPIndex := AddStringCPEntry(branch.TrueBB.Id.Value(), iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(trueBBNameCPIndex)); err != nil {
		return err
	}

	if branch.FalseBB == nil {
		return fmt.Errorf("branch terminator has nil FalseBB")
	}
	falseBBNameCPIndex := AddStringCPEntry(branch.FalseBB.Id.Value(), iw.cp)
	return binary.Write(iw.buf, binary.BigEndian, int32(falseBBNameCPIndex))
}

func (iw *BIRInstructionWriter) writeCall(call *Call) error {
	// Write isVirtual flag
	if err := binary.Write(iw.buf, binary.BigEndian, call.IsVirtual); err != nil {
		return err
	}

	// Write package index
	pkgIndex := AddPkgCPEntry(&call.CalleePkg, iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(pkgIndex)); err != nil {
		return err
	}

	// Write call name
	nameCPIndex := AddStringCPEntry(call.Name.Value(), iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	// Write arguments
	if err := binary.Write(iw.buf, binary.BigEndian, int32(len(call.Args))); err != nil {
		return err
	}
	for _, arg := range call.Args {
		if err := iw.writeOperand(&arg); err != nil {
			return err
		}
	}

	// Write LHS operand
	hasLhsOp := call.LhsOp != nil
	if err := binary.Write(iw.buf, binary.BigEndian, hasLhsOp); err != nil {
		return err
	}
	if hasLhsOp {
		if err := iw.writeOperand(call.LhsOp); err != nil {
			return err
		}
	}

	// Write then BB
	if call.ThenBB == nil {
		return fmt.Errorf("call terminator has nil ThenBB")
	}
	thenBBNameCPIndex := AddStringCPEntry(call.ThenBB.Id.Value(), iw.cp)
	return binary.Write(iw.buf, binary.BigEndian, int32(thenBBNameCPIndex))
}

func (iw *BIRInstructionWriter) writeOperand(operand *BIROperand) error {
	if operand == nil {
		return fmt.Errorf("operand cannot be nil")
	}

	if operand.VariableDcl == nil {
		return fmt.Errorf("operand variable declaration cannot be nil")
	}

	// Check if variable is ignored
	if operand.VariableDcl.IgnoreVariable {
		if err := binary.Write(iw.buf, binary.BigEndian, true); err != nil {
			return err
		}
		return WriteType(iw.cp, iw.buf, operand.VariableDcl.Type)
	}

	// Variable is not ignored
	if err := binary.Write(iw.buf, binary.BigEndian, false); err != nil {
		return err
	}

	// Write variable kind
	if err := binary.Write(iw.buf, binary.BigEndian, byte(operand.VariableDcl.Kind)); err != nil {
		return err
	}

	// Write variable scope
	if err := binary.Write(iw.buf, binary.BigEndian, byte(operand.VariableDcl.Scope)); err != nil {
		return err
	}

	// Write variable name
	varNameCPIndex := AddStringCPEntry(operand.VariableDcl.Name.Value(), iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(varNameCPIndex)); err != nil {
		return err
	}

	// For global/constant variables, write package index and type
	if operand.VariableDcl.Kind == VAR_KIND_GLOBAL || operand.VariableDcl.Kind == VAR_KIND_CONSTANT {
		// Get package ID from global variable
		// This requires access to the global variable's package ID
		// For now, write -1 as placeholder
		pkgIndex := int32(-1)
		if err := binary.Write(iw.buf, binary.BigEndian, pkgIndex); err != nil {
			return err
		}
		return WriteType(iw.cp, iw.buf, operand.VariableDcl.Type)
	}

	return nil
}

func (iw *BIRInstructionWriter) writeScopes(instruction BIRNonTerminator) {
	iw.instructionOffset++
	scope := iw.getInstructionScope(instruction)
	if scope == nil {
		return
	}
	iw.writeScopeFromInstruction(scope)
}

func (iw *BIRInstructionWriter) writeScope(terminator BIRTerminator) {
	if terminator.GetKind() == INSTRUCTION_KIND_RETURN {
		return
	}
	scope := iw.getTerminatorScope(terminator)
	if scope == nil {
		return
	}
	iw.writeScopeFromInstruction(scope)
}

func (iw *BIRInstructionWriter) getInstructionScope(instruction BIRNonTerminator) *BIRScope {
	switch inst := instruction.(type) {
	case *Move:
		return inst.BIRInstructionBase.Scope
	case *ConstantLoad:
		return inst.BIRInstructionBase.Scope
	case *BinaryOp:
		return inst.BIRInstructionBase.Scope
	case *UnaryOp:
		return inst.BIRInstructionBase.Scope
	default:
		return nil
	}
}

func (iw *BIRInstructionWriter) getTerminatorScope(terminator BIRTerminator) *BIRScope {
	switch term := terminator.(type) {
	case *Goto:
		return term.BIRTerminatorBase.BIRInstructionBase.Scope
	case *Return:
		return term.BIRTerminatorBase.BIRInstructionBase.Scope
	case *Branch:
		return term.BIRTerminatorBase.BIRInstructionBase.Scope
	case *Call:
		return term.BIRTerminatorBase.BIRInstructionBase.Scope
	default:
		return nil
	}
}

func (iw *BIRInstructionWriter) writeScopeFromInstruction(scope *BIRScope) {
	if scope == nil {
		return
	}

	if iw.completedScopes[scope] {
		return
	}

	iw.completedScopes[scope] = true
	iw.scopeCount++

	// Write scope ID
	if err := binary.Write(iw.scopeBuf, binary.BigEndian, int32(scope.Id)); err != nil {
		return
	}

	// Write instruction offset
	if err := binary.Write(iw.scopeBuf, binary.BigEndian, int32(iw.instructionOffset)); err != nil {
		return
	}

	// Write parent scope
	if scope.Parent != nil {
		if err := binary.Write(iw.scopeBuf, binary.BigEndian, true); err != nil {
			return
		}
		if err := binary.Write(iw.scopeBuf, binary.BigEndian, int32(scope.Parent.Id)); err != nil {
			return
		}
		iw.writeScopeFromInstruction(scope.Parent)
	} else {
		binary.Write(iw.scopeBuf, binary.BigEndian, false)
	}
}

// WriteErrorTable writes error table entries
func (iw *BIRInstructionWriter) WriteErrorTable(errorEntries []any) error {
	// TODO: Implement error table writing
	// For now, write count as 0
	return binary.Write(iw.buf, binary.BigEndian, int32(0))
}

// Helper methods to get position from instructions
func (iw *BIRInstructionWriter) getInstructionPos(instruction BIRNonTerminator) diagnostics.Location {
	// Try to get position from instruction base
	switch inst := instruction.(type) {
	case *Move:
		return inst.BIRInstructionBase.BIRNodeBase.Pos
	case *ConstantLoad:
		return inst.BIRInstructionBase.BIRNodeBase.Pos
	case *BinaryOp:
		return inst.BIRInstructionBase.BIRNodeBase.Pos
	case *UnaryOp:
		return inst.BIRInstructionBase.BIRNodeBase.Pos
	default:
		// Try to get from embedded BIRInstructionBase
		if base, ok := instruction.(interface{ GetBIRInstructionBase() *BIRInstructionBase }); ok {
			return base.GetBIRInstructionBase().BIRNodeBase.Pos
		}
		return nil
	}
}

func (iw *BIRInstructionWriter) getTerminatorPos(terminator BIRTerminator) diagnostics.Location {
	// Try to get position from terminator base
	switch term := terminator.(type) {
	case *Goto:
		return term.BIRTerminatorBase.BIRInstructionBase.BIRNodeBase.Pos
	case *Return:
		return term.BIRTerminatorBase.BIRInstructionBase.BIRNodeBase.Pos
	case *Branch:
		return term.BIRTerminatorBase.BIRInstructionBase.BIRNodeBase.Pos
	case *Call:
		return term.BIRTerminatorBase.BIRInstructionBase.BIRNodeBase.Pos
	default:
		return nil
	}
}

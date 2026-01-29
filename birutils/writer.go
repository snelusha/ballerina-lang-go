package birutils

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	"ballerina-lang-go/model"
)

const (
	BIR_MAGIC   = "\xba\x10\xc0\xde"
	BIR_VERSION = 75
)

type BIRWriter struct {
	cp  *ConstantPool
	pkg *bir.BIRPackage
}

func NewBIRWriter(pkg *bir.BIRPackage) *BIRWriter {
	return &BIRWriter{
		cp:  NewConstantPool(),
		pkg: pkg,
	}
}

func (bw *BIRWriter) Serialize() ([]byte, error) {
	birbuf := &bytes.Buffer{}

	// Write the package details
	if err := bw.writeInt32(birbuf, bw.cp.AddPackageCPEntry(bw.pkg.PackageID)); err != nil {
		return nil, err
	}

	// Write import module declarations
	if err := bw.writeImportModuleDecls(birbuf); err != nil {
		return nil, err
	}

	// Write constants
	if err := bw.writeConstants(birbuf); err != nil {
		return nil, err
	}

	// Write global vars
	if err := bw.writeGlobalVars(birbuf); err != nil {
		return nil, err
	}

	// Write functions
	if err := bw.writeFunctions(birbuf); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	_, err := buf.Write([]byte(BIR_MAGIC))
	if err != nil {
		return nil, err
	}

	if err := bw.writeInt32(buf, int32(BIR_VERSION)); err != nil {
		return nil, err
	}

	cpBytes, err := bw.cp.Serialize()
	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(cpBytes); err != nil {
		return nil, err
	}

	if _, err := buf.Write(birbuf.Bytes()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (bw *BIRWriter) writeImportModuleDecls(buf *bytes.Buffer) error {
	if err := bw.writeInt32(buf, int32(len(bw.pkg.ImportModules))); err != nil {
		return err
	}
	for _, imp := range bw.pkg.ImportModules {
		OrgName := imp.PackageID.OrgName.Value()
		PkgName := imp.PackageID.PkgName.Value()
		ModuleName := imp.PackageID.Name.Value()
		Version := imp.PackageID.Version.Value()

		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&OrgName)); err != nil {
			return err
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&PkgName)); err != nil {
			return err
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&ModuleName)); err != nil {
			return err
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&Version)); err != nil {
			return err
		}
	}

	return nil
}

func (bw *BIRWriter) writeConstants(buf *bytes.Buffer) error {
	if err := bw.writeInt32(buf, int32(len(bw.pkg.Constants))); err != nil {
		return err
	}

	for _, c := range bw.pkg.Constants {
		if err := bw.writeConstant(buf, &c); err != nil {
			return err
		}
	}

	return nil
}

func (bw *BIRWriter) writeConstant(buf *bytes.Buffer, constant *bir.BIRConstant) error {
	name := constant.Name.Value()
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&name)); err != nil {
		return err
	}

	if err := bw.writeInt64(buf, constant.Flags); err != nil {
		return err
	}

	if err := bw.writeUInt8(buf, uint8(constant.Origin)); err != nil {
		return err
	}

	constantType, err := bw.castToBType(constant.Type)
	if err != nil {
		return err
	}
	err = bw.writeType(buf, constantType)
	if err != nil {
		return err
	}

	birbuf := &bytes.Buffer{}
	constValueType, err := bw.castToBType(constant.ConstValue.Type)
	if err != nil {
		return err
	}
	err = bw.writeType(birbuf, constValueType)
	if err != nil {
		return err
	}
	err = bw.writeConstValue(birbuf, &constant.ConstValue)
	if err != nil {
		return err
	}

	if err := bw.writeInt64(buf, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err = buf.Write(birbuf.Bytes())
	return err
}

func (bw *BIRWriter) writeType(buf *bytes.Buffer, t ast.BType) error {
	idx := bw.cp.AddShapeCPEntry(t)
	return bw.writeInt32(buf, idx)
}

func (bw *BIRWriter) writeConstValue(buf *bytes.Buffer, cv *bir.ConstValue) error {
	bType, err := bw.castToBType(cv.Type)
	if err != nil {
		return err
	}
	switch model.TypeTags(bType.BTypeGetTag()) {
	case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
		if err := bw.writeInt32(buf, bw.cp.AddIntegerCPEntry(cv.Value.(int64))); err != nil {
			return err
		}
	case model.TypeTags_BYTE:
		if err := bw.writeInt32(buf, bw.cp.AddByteCPEntry(cv.Value.(byte))); err != nil {
			return err
		}
	case model.TypeTags_FLOAT:
		if err := bw.writeInt32(buf, bw.cp.AddFloatCPEntry(cv.Value.(float64))); err != nil {
			return err
		}
	case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(cv.Value.(*string))); err != nil {
			return err
		}
	case model.TypeTags_BOOLEAN:
		if err := bw.writeInt32(buf, bw.cp.AddBooleanCPEntry(cv.Value.(bool))); err != nil {
			return err
		}
	}
	return nil
}

func (bw *BIRWriter) writeGlobalVars(buf *bytes.Buffer) error {
	if err := bw.writeInt32(buf, int32(len(bw.pkg.GlobalVars))); err != nil {
		return err
	}

	for _, gv := range bw.pkg.GlobalVars {
		if err := bw.writeUInt8(buf, uint8(gv.Kind)); err != nil {
			return err
		}

		name := gv.Name.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&name)); err != nil {
			return err
		}

		if err := bw.writeInt64(buf, gv.Flags); err != nil {
			return err
		}

		if err := bw.writeUInt8(buf, uint8(gv.Origin)); err != nil {
			return err
		}

		gvType, err := bw.castToBType(gv.Type)
		if err != nil {
			return err
		}
		bw.writeType(buf, gvType)
	}

	return nil
}

func (bw *BIRWriter) writeFunctions(buf *bytes.Buffer) error {
	if err := bw.writeInt32(buf, int32(len(bw.pkg.Functions))); err != nil {
		return err
	}

	for _, fn := range bw.pkg.Functions {
		if err := bw.writeFunction(buf, &fn); err != nil {
			return err
		}
	}

	return nil
}

func (bw *BIRWriter) writeFunction(buf *bytes.Buffer, fn *bir.BIRFunction) error {
	name := fn.Name.Value()
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&name)); err != nil {
		return err
	}

	originalName := fn.OriginalName.Value()
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&originalName)); err != nil {
		return err
	}

	if err := bw.writeInt64(buf, fn.Flags); err != nil {
		return err
	}

	if err := bw.writeUInt8(buf, uint8(fn.Origin)); err != nil {
		return err
	}

	if err := bw.writeInt32(buf, int32(len(fn.RequiredParams))); err != nil {
		return err
	}
	for _, requiredParam := range fn.RequiredParams {
		name := requiredParam.Name.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&name)); err != nil {
			return err
		}
		if err := bw.writeInt64(buf, requiredParam.Flags); err != nil {
			return err
		}
	}

	birbuf := &bytes.Buffer{}

	if err := bw.writeInt32(birbuf, int32(fn.ArgsCount)); err != nil {
		return err
	}
	if err := bw.writeBool(birbuf, fn.ReturnVariable != nil); err != nil {
		return err
	}
	if fn.ReturnVariable != nil {
		if err := bw.writeUInt8(birbuf, uint8(fn.ReturnVariable.Kind)); err != nil {
			return err
		}

		retVarType, err := bw.castToBType(fn.ReturnVariable.Type)
		if err != nil {
			return err
		}
		err = bw.writeType(birbuf, retVarType)
		if err != nil {
			return err
		}
		returnVarName := fn.ReturnVariable.Name.Value()
		if err := bw.writeInt32(birbuf, bw.cp.AddStringCPEntry(&returnVarName)); err != nil {
			return err
		}
	}

	if err := bw.writeInt32(birbuf, int32(len(fn.LocalVars))); err != nil {
		return err
	}
	for _, localVar := range fn.LocalVars {
		if err := bw.writeLocalVar(birbuf, &localVar); err != nil {
			return err
		}
	}

	if err := bw.writeInt32(birbuf, int32(len(fn.BasicBlocks))); err != nil {
		return err
	}

	for _, bb := range fn.BasicBlocks {
		if err := bw.writeBasicBlock(birbuf, &bb); err != nil {
			return err
		}
	}

	if err := bw.writeInt64(buf, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err := buf.Write(birbuf.Bytes())

	return err
}

func (bw *BIRWriter) writeLocalVar(buf *bytes.Buffer, localVar *bir.BIRVariableDcl) error {
	if err := bw.writeUInt8(buf, uint8(localVar.Kind)); err != nil {
		return err
	}

	localVarType, err := bw.castToBType(localVar.Type)
	if err != nil {
		return err
	}
	if err := bw.writeType(buf, localVarType); err != nil {
		return err
	}
	localVarName := localVar.Name.Value()
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&localVarName)); err != nil {
		return err
	}

	if localVar.Kind == bir.VAR_KIND_ARG {
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&localVar.MetaVarName)); err != nil {
			return err
		}
	}

	if localVar.Kind == bir.VAR_KIND_LOCAL {
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&localVar.MetaVarName)); err != nil {
			return err
		}

		endBBId := ""
		if localVar.EndBB != nil {
			endBBId = localVar.EndBB.Id.Value()
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&endBBId)); err != nil {
			return err
		}

		startBBId := ""
		if localVar.StartBB != nil {
			startBBId = localVar.StartBB.Id.Value()
		}
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&startBBId)); err != nil {
			return err
		}

		if err := bw.writeInt32(buf, int32(localVar.InsOffset)); err != nil {
			return err
		}
	}
	return nil
}

func (bw *BIRWriter) writeBasicBlock(buf *bytes.Buffer, bb *bir.BIRBasicBlock) error {
	id := bb.Id.Value()
	if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&id)); err != nil {
		return err
	}
	// TODO: Adding the terminator instruction as well! Why?
	if err := bw.writeInt32(buf, int32(len(bb.Instructions))); err != nil {
		return err
	}

	for _, instr := range bb.Instructions {
		if err := bw.writeUInt8(buf, uint8(instr.GetKind())); err != nil {
			return err
		}
		if err := bw.writeInstruction(buf, instr); err != nil {
			return err
		}
	}

	if bb.Terminator == nil {
		panic(fmt.Sprintf("Basic block without a terminator %s", bb.Id.Value()))
	}

	if err := bw.writeUInt8(buf, uint8(bb.Terminator.GetKind())); err != nil {
		return err
	}

	return bw.writeTerminator(buf, bb.Terminator)
}

func (bw *BIRWriter) writeTerminator(buf *bytes.Buffer, term bir.BIRTerminator) error {
	switch term := term.(type) {
	case *bir.Goto:
		fmt.Println("Writing GOTO")
		id := term.ThenBB.Id.Value()
		fmt.Println("	" + id)
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&id)); err != nil {
			return err
		}
	case *bir.Branch:
		fmt.Println("Writing BRANCH")
		if err := bw.writeOperand(buf, term.Op); err != nil {
			return err
		}

		trueId := term.TrueBB.Id.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&trueId)); err != nil {
			return err
		}
		falseId := term.FalseBB.Id.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&falseId)); err != nil {
			return err
		}
	case *bir.Call:
		fmt.Println("Writing CALL")
		if err := bw.writeBool(buf, term.IsVirtual); err != nil {
			return err
		}
		pkgIdx := bw.cp.AddPackageCPEntry(term.CalleePkg)
		if err := bw.writeInt32(buf, pkgIdx); err != nil {
			return err
		}
		callName := term.Name.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&callName)); err != nil {
			return err
		}

		if err := bw.writeInt32(buf, int32(len(term.Args))); err != nil {
			return err
		}

		for _, arg := range term.Args {
			if err := bw.writeOperand(buf, &arg); err != nil {
				return err
			}
		}

		if term.LhsOp != nil {
			if err := bw.writeUInt8(buf, uint8(1)); err != nil {
				return err
			}
			if err := bw.writeOperand(buf, term.LhsOp); err != nil {
				return err
			}
		} else {
			if err := bw.writeUInt8(buf, uint8(0)); err != nil {
				return err
			}
		}

		thenBBId := term.ThenBB.Id.Value()
		if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&thenBBId)); err != nil {
			return err
		}

	case *bir.Return:
		fmt.Println("Writing RETURN")
	default:
		return fmt.Errorf("unsupported terminator type: %T", term)
	}
	return nil
}

func (bw *BIRWriter) writeInstruction(buf *bytes.Buffer, instr bir.BIRInstruction) error {
	switch instr := instr.(type) {
	case *bir.Move:
		if err := bw.writeOperand(buf, instr.RhsOp); err != nil {
			return err
		}
		return bw.writeOperand(buf, instr.LhsOp)
	case *bir.BinaryOp:
		if err := bw.writeOperand(buf, &instr.RhsOp1); err != nil {
			return err
		}
		if err := bw.writeOperand(buf, &instr.RhsOp2); err != nil {
			return err
		}
		return bw.writeOperand(buf, instr.LhsOp)
	case *bir.UnaryOp:
		if err := bw.writeOperand(buf, instr.RhsOp); err != nil {
			return err
		}
		return bw.writeOperand(buf, instr.LhsOp)
	case *bir.ConstantLoad:
		instrTypeCast, err := bw.castToBType(instr.Type)
		if err != nil {
			return err
		}
		if err := bw.writeType(buf, instrTypeCast); err != nil {
			return err
		}
		if err := bw.writeOperand(buf, instr.LhsOp); err != nil {
			return err
		}

		switch model.TypeTags(instrTypeCast.BTypeGetTag()) {
		case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
			if err := bw.writeInt32(buf, bw.cp.AddIntegerCPEntry(instr.Value.(int64))); err != nil {
				return err
			}
		case model.TypeTags_BYTE:
			if err := bw.writeInt32(buf, bw.cp.AddByteCPEntry(instr.Value.(byte))); err != nil {
				return err
			}
		case model.TypeTags_FLOAT:
			if err := bw.writeInt32(buf, bw.cp.AddFloatCPEntry(instr.Value.(float64))); err != nil {
				return err
			}
		case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
			val := instr.Value.(string)
			if err := bw.writeInt32(buf, bw.cp.AddStringCPEntry(&val)); err != nil {
				return err
			}
		case model.TypeTags_BOOLEAN:
			if err := bw.writeInt32(buf, bw.cp.AddBooleanCPEntry(instr.Value.(bool))); err != nil {
				return err
			}
		}

	case *bir.FieldAccess:
		panic("FieldAccess not implemented")
	case *bir.NewArray:
		panic("NewArray not implemented")
	}
	return nil
}

func (bw *BIRWriter) writeOperand(buf *bytes.Buffer, op *bir.BIROperand) error {
	if op.VariableDcl.IgnoreVariable {
		if err := bw.writeBool(buf, true); err != nil {
			return err
		}

		opType, err := bw.castToBType(op.VariableDcl.Type)
		if err != nil {
			return err
		}
		return bw.writeType(buf, opType)
	}

	if err := bw.writeBool(buf, false); err != nil {
		return err
	}
	if err := bw.writeUInt8(buf, uint8(op.VariableDcl.Kind)); err != nil {
		return err
	}

	if err := bw.writeUInt8(buf, uint8(op.VariableDcl.Scope)); err != nil {
		return err
	}

	varName := op.VariableDcl.Name.Value()
	nameIdx := bw.cp.AddStringCPEntry(&varName)

	return bw.writeInt32(buf, nameIdx)
}

func (bw *BIRWriter) writeInt8(buf *bytes.Buffer, val int8) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeUInt8(buf *bytes.Buffer, val uint8) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeInt32(buf *bytes.Buffer, val int32) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeInt64(buf *bytes.Buffer, val int64) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeFloat64(buf *bytes.Buffer, val float64) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) writeBool(buf *bytes.Buffer, val bool) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (bw *BIRWriter) castToBType(t any) (ast.BType, error) {
	bType, ok := t.(ast.BType)
	if !ok {
		return nil, fmt.Errorf("expected ast.BType, got %T", t)
	}
	return bType, nil
}

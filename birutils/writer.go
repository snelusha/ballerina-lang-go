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

func (w *BIRWriter) Serialize() ([]byte, error) {
	birbuf := &bytes.Buffer{}

	// Write the package details
	if err := w.writeInt32(birbuf, w.cp.AddPackageCPEntry(w.pkg.PackageID)); err != nil {
		return nil, err
	}

	// Write import module declarations
	if err := w.writeImportModuleDecls(birbuf); err != nil {
		return nil, err
	}

	// Write constants
	if err := w.writeConstants(birbuf); err != nil {
		return nil, err
	}

	// Write global vars
	if err := w.writeGlobalVars(birbuf); err != nil {
		return nil, err
	}

	// Write functions
	if err := w.writeFunctions(birbuf); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	_, err := buf.Write([]byte(BIR_MAGIC))
	if err != nil {
		return nil, err
	}

	if err := w.writeInt32(buf, int32(BIR_VERSION)); err != nil {
		return nil, err
	}

	cpBytes, err := w.cp.Serialize()
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

func (w *BIRWriter) writeImportModuleDecls(buf *bytes.Buffer) error {
	if err := w.writeInt32(buf, int32(len(w.pkg.ImportModules))); err != nil {
		return err
	}
	for _, imp := range w.pkg.ImportModules {
		OrgName := imp.PackageID.OrgName.Value()
		PkgName := imp.PackageID.PkgName.Value()
		ModuleName := imp.PackageID.Name.Value()
		Version := imp.PackageID.Version.Value()

		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&OrgName)); err != nil {
			return err
		}
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&PkgName)); err != nil {
			return err
		}
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&ModuleName)); err != nil {
			return err
		}
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&Version)); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRWriter) writeConstants(buf *bytes.Buffer) error {
	if err := w.writeInt32(buf, int32(len(w.pkg.Constants))); err != nil {
		return err
	}

	for _, c := range w.pkg.Constants {
		if err := w.writeConstant(buf, &c); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRWriter) writeConstant(buf *bytes.Buffer, constant *bir.BIRConstant) error {
	name := constant.Name.Value()
	if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&name)); err != nil {
		return err
	}

	if err := w.writeInt64(buf, constant.Flags); err != nil {
		return err
	}

	if err := w.writeUInt8(buf, uint8(constant.Origin)); err != nil {
		return err
	}

	err := w.writeType(buf, constant.Type.(ast.BType))
	if err != nil {
		return err
	}

	birbuf := &bytes.Buffer{}
	err = w.writeType(birbuf, constant.ConstValue.Type.(ast.BType))
	if err != nil {
		return err
	}
	err = w.writeConstValue(birbuf, &constant.ConstValue)
	if err != nil {
		return err
	}

	if err := w.writeInt64(buf, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err = buf.Write(birbuf.Bytes())
	return err
}

func (w *BIRWriter) writeType(buf *bytes.Buffer, t ast.BType) error {
	idx := w.cp.AddShapeCPEntry(t)
	return w.writeInt32(buf, idx)
}

func (w *BIRWriter) writeConstValue(buf *bytes.Buffer, cv *bir.ConstValue) error {
	bType, ok := cv.Type.(ast.BType)
	if !ok {
		return fmt.Errorf("unsupported const value type: %T", cv.Type)
	}
	switch model.TypeTags(bType.BTypeGetTag()) {
	case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
		if err := w.writeInt32(buf, w.cp.AddIntegerCPEntry(cv.Value.(int64))); err != nil {
			return err
		}
	case model.TypeTags_BYTE:
		if err := w.writeInt32(buf, w.cp.AddByteCPEntry(cv.Value.(byte))); err != nil {
			return err
		}
	case model.TypeTags_FLOAT:
		if err := w.writeInt32(buf, w.cp.AddFloatCPEntry(cv.Value.(float64))); err != nil {
			return err
		}
	case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(cv.Value.(*string))); err != nil {
			return err
		}
	case model.TypeTags_BOOLEAN:
		if err := w.writeInt32(buf, w.cp.AddBooleanCPEntry(cv.Value.(bool))); err != nil {
			return err
		}
	}
	return nil
}

func (w *BIRWriter) writeGlobalVars(buf *bytes.Buffer) error {
	if err := w.writeInt32(buf, int32(len(w.pkg.GlobalVars))); err != nil {
		return err
	}

	for _, gv := range w.pkg.GlobalVars {
		if err := w.writeUInt8(buf, uint8(gv.Kind)); err != nil {
			return err
		}

		name := gv.Name.Value()
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&name)); err != nil {
			return err
		}

		if err := w.writeInt64(buf, gv.Flags); err != nil {
			return err
		}

		if err := w.writeUInt8(buf, uint8(gv.Origin)); err != nil {
			return err
		}

		w.writeType(buf, gv.Type.(ast.BType))
	}

	return nil
}

func (w *BIRWriter) writeFunctions(buf *bytes.Buffer) error {
	if err := w.writeInt32(buf, int32(len(w.pkg.Functions))); err != nil {
		return err
	}

	for _, fn := range w.pkg.Functions {
		if err := w.writeFunction(buf, &fn); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRWriter) writeFunction(buf *bytes.Buffer, fn *bir.BIRFunction) error {
	name := fn.Name.Value()
	if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&name)); err != nil {
		return err
	}

	originalName := fn.OriginalName.Value()
	if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&originalName)); err != nil {
		return err
	}

	if err := w.writeInt64(buf, fn.Flags); err != nil {
		return err
	}

	if err := w.writeUInt8(buf, uint8(fn.Origin)); err != nil {
		return err
	}

	if err := w.writeInt32(buf, int32(len(fn.RequiredParams))); err != nil {
		return err
	}
	for _, requiredParam := range fn.RequiredParams {
		name := requiredParam.Name.Value()
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&name)); err != nil {
			return err
		}
		if err := w.writeInt64(buf, requiredParam.Flags); err != nil {
			return err
		}
	}

	birbuf := &bytes.Buffer{}

	if err := w.writeInt32(birbuf, int32(fn.ArgsCount)); err != nil {
		return err
	}
	if err := w.writeBool(birbuf, fn.ReturnVariable != nil); err != nil {
		return err
	}
	if fn.ReturnVariable != nil {
		if err := w.writeUInt8(birbuf, uint8(fn.ReturnVariable.Kind)); err != nil {
			return err
		}
		err := w.writeType(birbuf, fn.ReturnVariable.Type.(ast.BType))
		if err != nil {
			return err
		}
		returnVarName := fn.ReturnVariable.Name.Value()
		if err := w.writeInt32(birbuf, w.cp.AddStringCPEntry(&returnVarName)); err != nil {
			return err
		}
	}

	if err := w.writeInt32(birbuf, int32(len(fn.LocalVars))); err != nil {
		return err
	}
	for _, localVar := range fn.LocalVars {
		if err := w.writeLocalVar(birbuf, &localVar); err != nil {
			return err
		}
	}

	if err := w.writeInt32(birbuf, int32(len(fn.BasicBlocks))); err != nil {
		return err
	}

	for _, bb := range fn.BasicBlocks {
		if err := w.writeBasicBlock(birbuf, &bb); err != nil {
			return err
		}
	}

	if err := w.writeInt64(buf, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err := buf.Write(birbuf.Bytes())

	return err
}

func (w *BIRWriter) writeLocalVar(buf *bytes.Buffer, localVar *bir.BIRVariableDcl) error {
	if err := w.writeUInt8(buf, uint8(localVar.Kind)); err != nil {
		return err
	}
	if err := w.writeType(buf, localVar.Type.(ast.BType)); err != nil {
		return err
	}
	localVarName := localVar.Name.Value()
	if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&localVarName)); err != nil {
		return err
	}

	if localVar.Kind == bir.VAR_KIND_ARG {
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&localVar.MetaVarName)); err != nil {
			return err
		}
	}

	if localVar.Kind == bir.VAR_KIND_LOCAL {
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&localVar.MetaVarName)); err != nil {
			return err
		}

		endBBId := ""
		if localVar.EndBB != nil {
			endBBId = localVar.EndBB.Id.Value()
		}
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&endBBId)); err != nil {
			return err
		}

		startBBId := ""
		if localVar.StartBB != nil {
			startBBId = localVar.StartBB.Id.Value()
		}
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&startBBId)); err != nil {
			return err
		}

		if err := w.writeInt32(buf, int32(localVar.InsOffset)); err != nil {
			return err
		}
	}
	return nil
}

func (w *BIRWriter) writeBasicBlock(buf *bytes.Buffer, bb *bir.BIRBasicBlock) error {
	id := bb.Id.Value()
	if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&id)); err != nil {
		return err
	}
	// TODO: Adding the terminator instruction as well! Why?
	if err := w.writeInt32(buf, int32(len(bb.Instructions))); err != nil {
		return err
	}

	for _, instr := range bb.Instructions {
		if err := w.writeUInt8(buf, uint8(instr.GetKind())); err != nil {
			return err
		}
		if err := w.writeInstruction(buf, instr); err != nil {
			return err
		}
	}

	if bb.Terminator == nil {
		panic(fmt.Sprintf("Basic block without a terminator %s", bb.Id.Value()))
	}

	if err := w.writeUInt8(buf, uint8(bb.Terminator.GetKind())); err != nil {
		return err
	}

	return w.writeTerminator(buf, bb.Terminator)
}

func (w *BIRWriter) writeTerminator(buf *bytes.Buffer, term bir.BIRTerminator) error {
	switch term := term.(type) {
	case *bir.Goto:
		fmt.Println("Writing GOTO")
		id := term.ThenBB.Id.Value()
		fmt.Println("	" + id)
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&id)); err != nil {
			return err
		}
	case *bir.Branch:
		fmt.Println("Writing BRANCH")
		if err := w.writeOperand(buf, term.Op); err != nil {
			return err
		}

		trueId := term.TrueBB.Id.Value()
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&trueId)); err != nil {
			return err
		}
		falseId := term.FalseBB.Id.Value()
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&falseId)); err != nil {
			return err
		}
	case *bir.Call:
		fmt.Println("Writing CALL")
		if err := w.writeBool(buf, term.IsVirtual); err != nil {
			return err
		}
		pkgIdx := w.cp.AddPackageCPEntry(term.CalleePkg)
		if err := w.writeInt32(buf, pkgIdx); err != nil {
			return err
		}
		callName := term.Name.Value()
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&callName)); err != nil {
			return err
		}

		if err := w.writeInt32(buf, int32(len(term.Args))); err != nil {
			return err
		}

		for _, arg := range term.Args {
			if err := w.writeOperand(buf, &arg); err != nil {
				return err
			}
		}

		if term.LhsOp != nil {
			if err := w.writeUInt8(buf, uint8(1)); err != nil {
				return err
			}
			if err := w.writeOperand(buf, term.LhsOp); err != nil {
				return err
			}
		} else {
			if err := w.writeUInt8(buf, uint8(0)); err != nil {
				return err
			}
		}

		thenBBId := term.ThenBB.Id.Value()
		if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&thenBBId)); err != nil {
			return err
		}

	case *bir.Return:
		fmt.Println("Writing RETURN")
	default:
		return fmt.Errorf("unsupported terminator type: %T", term)
	}
	return nil
}

func (w *BIRWriter) writeInstruction(buf *bytes.Buffer, instr bir.BIRInstruction) error {
	switch instr := instr.(type) {
	case *bir.Move:
		if err := w.writeOperand(buf, instr.RhsOp); err != nil {
			return err
		}
		return w.writeOperand(buf, instr.LhsOp)
	case *bir.BinaryOp:
		if err := w.writeOperand(buf, &instr.RhsOp1); err != nil {
			return err
		}
		if err := w.writeOperand(buf, &instr.RhsOp2); err != nil {
			return err
		}
		return w.writeOperand(buf, instr.LhsOp)
	case *bir.UnaryOp:
		if err := w.writeOperand(buf, instr.RhsOp); err != nil {
			return err
		}
		return w.writeOperand(buf, instr.LhsOp)
	case *bir.ConstantLoad:
		if err := w.writeType(buf, instr.Type.(ast.BType)); err != nil {
			return err
		}
		if err := w.writeOperand(buf, instr.LhsOp); err != nil {
			return err
		}

		instrType, ok := instr.Type.(ast.BType)
		if !ok {
			return fmt.Errorf("unsupported constant load type: %T", instr.Type)
		}

		switch model.TypeTags(instrType.BTypeGetTag()) {
		case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
			if err := w.writeInt32(buf, w.cp.AddIntegerCPEntry(instr.Value.(int64))); err != nil {
				return err
			}
		case model.TypeTags_BYTE:
			if err := w.writeInt32(buf, w.cp.AddByteCPEntry(instr.Value.(byte))); err != nil {
				return err
			}
		case model.TypeTags_FLOAT:
			if err := w.writeInt32(buf, w.cp.AddFloatCPEntry(instr.Value.(float64))); err != nil {
				return err
			}
		case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
			val := instr.Value.(string)
			if err := w.writeInt32(buf, w.cp.AddStringCPEntry(&val)); err != nil {
				return err
			}
		case model.TypeTags_BOOLEAN:
			if err := w.writeInt32(buf, w.cp.AddBooleanCPEntry(instr.Value.(bool))); err != nil {
				return err
			}
		}

	case *bir.FieldAccess:
		// fmt.Println("Writing FieldAccess instruction")
	case *bir.NewArray:
		// fmt.Println("Writing NewArray instruction")
	}
	return nil
}

func (w *BIRWriter) writeOperand(buf *bytes.Buffer, op *bir.BIROperand) error {
	if op.VariableDcl.IgnoreVariable {
		if err := w.writeBool(buf, true); err != nil {
			return err
		}
		return w.writeType(buf, op.VariableDcl.Type.(ast.BType))
	}

	if err := w.writeBool(buf, false); err != nil {
		return err
	}
	if err := w.writeUInt8(buf, uint8(op.VariableDcl.Kind)); err != nil {
		return err
	}

	if err := w.writeUInt8(buf, uint8(op.VariableDcl.Scope)); err != nil {
		return err
	}

	varName := op.VariableDcl.Name.Value()
	nameIdx := w.cp.AddStringCPEntry(&varName)

	return w.writeInt32(buf, nameIdx)
}

func (w *BIRWriter) writeInt8(buf *bytes.Buffer, val int8) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (w *BIRWriter) writeUInt8(buf *bytes.Buffer, val uint8) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (w *BIRWriter) writeInt32(buf *bytes.Buffer, val int32) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (w *BIRWriter) writeInt64(buf *bytes.Buffer, val int64) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (w *BIRWriter) writeFloat64(buf *bytes.Buffer, val float64) error {
	return binary.Write(buf, binary.BigEndian, val)
}

func (w *BIRWriter) writeBool(buf *bytes.Buffer, val bool) error {
	return binary.Write(buf, binary.BigEndian, val)
}

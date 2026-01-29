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
	if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddPackageCPEntry(w.pkg.PackageID)); err != nil {
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

	if err := binary.Write(buf, binary.BigEndian, int32(BIR_VERSION)); err != nil {
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
	if err := binary.Write(buf, binary.BigEndian, int32(len(w.pkg.ImportModules))); err != nil {
		return err
	}
	for _, imp := range w.pkg.ImportModules {
		OrgName := imp.PackageID.OrgName.Value()
		PkgName := imp.PackageID.PkgName.Value()
		ModuleName := imp.PackageID.Name.Value()
		Version := imp.PackageID.Version.Value()

		if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&OrgName)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&PkgName)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&ModuleName)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&Version)); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRWriter) writeConstants(buf *bytes.Buffer) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(w.pkg.Constants))); err != nil {
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
	if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&name)); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, constant.Flags); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, constant.Origin); err != nil {
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

	if err := binary.Write(buf, binary.BigEndian, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err = buf.Write(birbuf.Bytes())
	return err
}

func (w *BIRWriter) writeType(buf *bytes.Buffer, t ast.BType) error {
	idx := w.cp.AddShapeCPEntry(t)
	return binary.Write(buf, binary.BigEndian, idx)
}

func (w *BIRWriter) writeConstValue(buf *bytes.Buffer, cv *bir.ConstValue) error {
	bType, ok := cv.Type.(ast.BType)
	if !ok {
		return fmt.Errorf("unsupported const value type: %T", cv.Type)
	}
	switch model.TypeTags(bType.BTypeGetTag()) {
	case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddIntegerCPEntry(cv.Value.(int64))); err != nil {
			return err
		}
	case model.TypeTags_BYTE:
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddByteCPEntry(cv.Value.(byte))); err != nil {
			return err
		}
	case model.TypeTags_FLOAT:
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddFloatCPEntry(cv.Value.(float64))); err != nil {
			return err
		}
	case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(cv.Value.(*string))); err != nil {
			return err
		}
	case model.TypeTags_BOOLEAN:
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddBooleanCPEntry(cv.Value.(bool))); err != nil {
			return err
		}
	}
	return nil
}

func (w *BIRWriter) writeGlobalVars(buf *bytes.Buffer) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(w.pkg.GlobalVars))); err != nil {
		return err
	}

	for _, gv := range w.pkg.GlobalVars {
		if err := binary.Write(buf, binary.BigEndian, gv.Kind); err != nil {
			return err
		}

		name := gv.Name.Value()
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&name)); err != nil {
			return err
		}

		if err := binary.Write(buf, binary.BigEndian, gv.Flags); err != nil {
			return err
		}

		if err := binary.Write(buf, binary.BigEndian, gv.Origin); err != nil {
			return err
		}

		w.writeType(buf, gv.Type.(ast.BType))
	}

	return nil
}

func (w *BIRWriter) writeFunctions(buf *bytes.Buffer) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(w.pkg.Functions))); err != nil {
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
	if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&name)); err != nil {
		return err
	}

	originalName := fn.OriginalName.Value()
	if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&originalName)); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, fn.Flags); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, fn.Origin); err != nil {
		return err
	}

	// TODO: Type

	if err := binary.Write(buf, binary.BigEndian, int32(len(fn.RequiredParams))); err != nil {
		return err
	}
	for _, requiredParam := range fn.RequiredParams {
		name := requiredParam.Name.Value()
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&name)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, requiredParam.Flags); err != nil {
			return err
		}
	}

	birbuf := &bytes.Buffer{}

	if err := binary.Write(birbuf, binary.BigEndian, int32(fn.ArgsCount)); err != nil {
		return err
	}
	if err := binary.Write(birbuf, binary.BigEndian, fn.ReturnVariable != nil); err != nil {
		return err
	}
	if fn.ReturnVariable != nil {
		if err := binary.Write(birbuf, binary.BigEndian, fn.ReturnVariable.Kind); err != nil {
			return err
		}
		err := w.writeType(birbuf, fn.ReturnVariable.Type.(ast.BType))
		if err != nil {
			return err
		}
		returnVarName := fn.ReturnVariable.Name.Value()
		if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&returnVarName)); err != nil {
			return err
		}
	}

	if err := binary.Write(birbuf, binary.BigEndian, int32(len(fn.LocalVars))); err != nil {
		return err
	}
	for _, localVar := range fn.LocalVars {
		if err := binary.Write(birbuf, binary.BigEndian, localVar.Kind); err != nil {
			return err
		}
		if err := w.writeType(birbuf, localVar.Type.(ast.BType)); err != nil {
			return err
		}
		localVarName := localVar.Name.Value()
		if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&localVarName)); err != nil {
			return err
		}

		if localVar.Kind == bir.VAR_KIND_ARG {
			if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&localVar.MetaVarName)); err != nil {
				return err
			}
		}

		if localVar.Kind == bir.VAR_KIND_LOCAL {
			if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&localVar.MetaVarName)); err != nil {
				return err
			}

			if localVar.EndBB != nil {
				endBBId := localVar.EndBB.Id.Value()
				if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&endBBId)); err != nil {
					return err
				}
			} else {
				endBBId := ""
				if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&endBBId)); err != nil {
					return err
				}
			}

			if localVar.StartBB != nil {
				startBBId := localVar.StartBB.Id.Value()
				if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&startBBId)); err != nil {
					return err
				}
			} else {
				startBBId := ""
				if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&startBBId)); err != nil {
					return err
				}
			}

			if err := binary.Write(birbuf, binary.BigEndian, int32(localVar.InsOffset)); err != nil {
				return err
			}
		}

	}

	// basic blocks
	if err := binary.Write(birbuf, binary.BigEndian, int32(len(fn.BasicBlocks))); err != nil {
		return err
	}

	for _, bb := range fn.BasicBlocks {
		id := bb.Id.Value()
		if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&id)); err != nil {
			return err
		}
		// Adding the terminator instruction as well
		if err := binary.Write(birbuf, binary.BigEndian, int32(len(bb.Instructions))); err != nil {
			return err
		}

		for _, instr := range bb.Instructions {
			if err := binary.Write(birbuf, binary.BigEndian, instr.GetKind()); err != nil {
				return err
			}
			if err := w.writeInstruction(birbuf, instr); err != nil {
				return err
			}
		}

		if bb.Terminator == nil {
			panic(fmt.Sprintf("Basic block without a terminator %s", bb.Id.Value()))
		}

		if err := binary.Write(birbuf, binary.BigEndian, bb.Terminator.GetKind()); err != nil {
			return err
		}
		// FIXME: Write terminator

		switch term := bb.Terminator.(type) {
		case *bir.Goto:
			fmt.Println("Writing GOTO")
			id := term.ThenBB.Id.Value()
			fmt.Println("	" + id)
			if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&id)); err != nil {
				return err
			}
		case *bir.Branch:
			fmt.Println("Writing BRANCH")
			if err := w.writeOperand(birbuf, term.Op); err != nil {
				return err
			}

			trueId := term.TrueBB.Id.Value()
			if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&trueId)); err != nil {
				return err
			}
			falseId := term.FalseBB.Id.Value()
			if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&falseId)); err != nil {
				return err
			}
		case *bir.Call:
			fmt.Println("Writing CALL")
			if err := binary.Write(birbuf, binary.BigEndian, term.IsVirtual); err != nil {
				return err
			}
			pkgIdx := w.cp.AddPackageCPEntry(term.CalleePkg)
			if err := binary.Write(birbuf, binary.BigEndian, pkgIdx); err != nil {
				return err
			}
			callName := term.Name.Value()
			if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&callName)); err != nil {
				return err
			}

			if err := binary.Write(birbuf, binary.BigEndian, int32(len(term.Args))); err != nil {
				return err
			}

			for _, arg := range term.Args {
				if err := w.writeOperand(birbuf, &arg); err != nil {
					return err
				}
			}

			if term.LhsOp != nil {
				if err := binary.Write(birbuf, binary.BigEndian, uint8(1)); err != nil {
					return err
				}
				if err := w.writeOperand(birbuf, term.LhsOp); err != nil {
					return err
				}
			} else {
				if err := binary.Write(birbuf, binary.BigEndian, uint8(0)); err != nil {
					return err
				}
			}

			thenBBId := term.ThenBB.Id.Value()
			if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&thenBBId)); err != nil {
				return err
			}

		case *bir.Return:
			fmt.Println("Writing RETURN")
		default:
			fmt.Printf("Done with %T", bb.Terminator)
			panic("yeah!")
		}
	}

	if err := binary.Write(buf, binary.BigEndian, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err := buf.Write(birbuf.Bytes())

	return err
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
			if err := binary.Write(buf, binary.BigEndian, w.cp.AddIntegerCPEntry(instr.Value.(int64))); err != nil {
				return err
			}
		case model.TypeTags_BYTE:
			if err := binary.Write(buf, binary.BigEndian, w.cp.AddByteCPEntry(instr.Value.(byte))); err != nil {
				return err
			}
		case model.TypeTags_FLOAT:
			if err := binary.Write(buf, binary.BigEndian, w.cp.AddFloatCPEntry(instr.Value.(float64))); err != nil {
				return err
			}
		case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
			val := instr.Value.(string)
			if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&val)); err != nil {
				return err
			}
		case model.TypeTags_BOOLEAN:
			if err := binary.Write(buf, binary.BigEndian, w.cp.AddBooleanCPEntry(instr.Value.(bool))); err != nil {
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
		if err := binary.Write(buf, binary.BigEndian, true); err != nil {
			return err
		}
		return w.writeType(buf, op.VariableDcl.Type.(ast.BType))
	}

	if err := binary.Write(buf, binary.BigEndian, false); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, op.VariableDcl.Kind); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, op.VariableDcl.Scope); err != nil {
		return err
	}

	varName := op.VariableDcl.Name.Value()
	nameIdx := w.cp.AddStringCPEntry(&varName)

	return binary.Write(buf, binary.BigEndian, nameIdx)
}

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

	fmt.Printf("[WRITER] REQUIRED PARAMS COUNT: %d IN FUNCTION: %s\n", len(fn.RequiredParams), name)
	if err := binary.Write(buf, binary.BigEndian, int32(len(fn.RequiredParams))); err != nil {
		return err
	}
	for _, requiredParam := range fn.RequiredParams {
		name := requiredParam.Name.Value()
		fmt.Printf("[WRITER] REQUIRED PARAM NAME: %s\n", name)
		if err := binary.Write(buf, binary.BigEndian, w.cp.AddStringCPEntry(&name)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, requiredParam.Flags); err != nil {
			return err
		}
	}
	fmt.Println()

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

	// fmt.Printf("[WRITER] LOCAL VARS COUNT: %d\n", len(fn.LocalVars))
	//
	// if err := binary.Write(birbuf, binary.BigEndian, int32(len(fn.LocalVars))); err != nil {
	// 	return err
	// }
	// for _, localVar := range fn.LocalVars {
	// 	_ = localVar // To avoid unused variable error

	// if err := binary.Write(birbuf, binary.BigEndian, localVar.Kind); err != nil {
	// 	return err
	// }
	// if err := w.writeType(birbuf, localVar.Type.(ast.BType)); err != nil {
	// 	return err
	// }
	// localVarName := localVar.Name.Value()
	// if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&localVarName)); err != nil {
	// 	return err
	// }

	// if localVar.Kind == bir.VAR_KIND_ARG {
	// 	fmt.Printf("[WRITER] KIND ARG\n")
	// 	if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&localVar.MetaVarName)); err != nil {
	// 		return err
	// 	}
	// }

	// if localVar.Kind == bir.VAR_KIND_LOCAL {
	//
	// 	fmt.Printf("[WRITER] KIND LOCAL\n")
	// 	if err := binary.Write(birbuf, binary.BigEndian, localVar.MetaVarName); err != nil {
	// 		return err
	// 	}
	//
	// 	if localVar.EndBB != nil {
	// 		fmt.Printf("[WRITER] END BB NOT NIL\n")
	// 		endBBId := localVar.EndBB.Id.Value()
	// 		if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&endBBId)); err != nil {
	// 			return err
	// 		}
	// 	}
	//
	// 	if localVar.StartBB != nil {
	// 		fmt.Printf("[WRITER] START BB NOT NIL\n")
	// 		startBBId := localVar.StartBB.Id.Value()
	// 		if err := binary.Write(birbuf, binary.BigEndian, w.cp.AddStringCPEntry(&startBBId)); err != nil {
	// 			return err
	// 		}
	// 	}
	//
	// 	fmt.Printf("[WRITER] INS OFFSET: %d\n", localVar.InsOffset)
	//
	// 	if err := binary.Write(birbuf, binary.BigEndian, int32(localVar.InsOffset)); err != nil {
	// 		return err
	// 	}
	// }
	// }

	if err := binary.Write(buf, binary.BigEndian, int64(birbuf.Len())); err != nil {
		return err
	}
	_, err := buf.Write(birbuf.Bytes())

	return err
}

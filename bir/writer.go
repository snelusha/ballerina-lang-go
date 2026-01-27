package bir

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"ballerina-lang-go/model"
)

const (
	BIR_MAGIC   = "\xba\x10\xc0\xde"
	BIR_VERSION = 75
)

type BIRWriter struct {
	cp  *ConstantPool
	pkg *BIRPackage
}

func NewBIRWriter(pkg *BIRPackage) *BIRWriter {
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

func (w *BIRWriter) writeConstant(buf *bytes.Buffer, constant *BIRConstant) error {
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

	err := w.writeType(buf, &constant.Type)
	if err != nil {
		return err
	}

	birbuf := &bytes.Buffer{}
	err = w.writeType(birbuf, &constant.ConstValue.Type)
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
	fmt.Printf("Constant Value Length: %d\n", birbuf.Len())
	_, err = buf.Write(birbuf.Bytes())
	return err
}

func (w *BIRWriter) writeType(buf *bytes.Buffer, t *minimalBType) error {
	idx := w.cp.AddShapeCPEntry(t)
	fmt.Printf("[WRITER] idx %d\n", idx)
	return binary.Write(buf, binary.BigEndian, idx)
}

func (w *BIRWriter) writeConstValue(buf *bytes.Buffer, cv *ConstValue) error {
	fmt.Printf("Writing Const Value of type tag: %d\n", cv.Type.tag)
	switch model.TypeTags(cv.Type.tag) {
	case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
		fmt.Printf("Writing INT Const Value: %d\n", cv.Value.(int64))
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

		// fmt.Printf("Global Var Type: %+v\n", gv.Type)

		w.writeType(buf, &gv.Type)
	}

	return nil
}

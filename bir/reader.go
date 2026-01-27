package bir

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"ballerina-lang-go/model"
)

type BIRReader struct {
	r  *bytes.Reader
	cp []any
}

func NewBIRReader(data []byte) *BIRReader {
	return &BIRReader{
		r: bytes.NewReader(data),
	}
}

func (br *BIRReader) LoadBIRPackage() (*BIRPackage, error) {
	return br.readPackage()
}

func (br *BIRReader) readPackage() (*BIRPackage, error) {
	magic := make([]byte, 4)
	_, err := br.r.Read(magic)
	if err != nil {
		return nil, err
	}

	if string(magic) != BIR_MAGIC {
		return nil, fmt.Errorf("invalid BIR magic: %x", magic)
	}

	var version int32
	if err := binary.Read(br.r, binary.BigEndian, &version); err != nil {
		return nil, err
	}

	if version != BIR_VERSION {
		return nil, fmt.Errorf("unsupported BIR version: %d", version)
	}

	if err := br.readConstantPool(); err != nil {
		return nil, fmt.Errorf("reading constant pool: %w", err)
	}

	var pkgCPIndex int32
	if err := binary.Read(br.r, binary.BigEndian, &pkgCPIndex); err != nil {
		return nil, fmt.Errorf("reading package CP index: %w", err)
	}

	pkgID := br.getPackageFromCP(int(pkgCPIndex))

	imports, err := br.readImports()
	if err != nil {
		return nil, fmt.Errorf("reading imports: %w", err)
	}

	constants, err := br.readConstants()
	if err != nil {
		return nil, fmt.Errorf("reading constants: %w", err)
	}

	globalVars, err := br.readGlobalVars()
	if err != nil {
		return nil, fmt.Errorf("reading global vars: %w", err)
	}

	return &BIRPackage{
		PackageID:     pkgID,
		ImportModules: imports,
		Constants:     constants,
		GlobalVars:    globalVars,
	}, nil
}

func (br *BIRReader) readConstantPool() error {
	var cpSize int32
	if err := binary.Read(br.r, binary.BigEndian, &cpSize); err != nil {
		return err
	}

	br.cp = make([]any, cpSize)

	for i := 0; i < int(cpSize); i++ {
		var tag int8
		if err := binary.Read(br.r, binary.BigEndian, &tag); err != nil {
			return fmt.Errorf("reading CP entry %d tag: %w", i, err)
		}

		switch tag {
		case 0: // NULL/placeholder entry
			br.cp[i] = nil
		case 1: // INTEGER
			var value int64
			if err := binary.Read(br.r, binary.BigEndian, &value); err != nil {
				return err
			}
			br.cp[i] = value

		case 2: // FLOAT
			var value float64
			if err := binary.Read(br.r, binary.BigEndian, &value); err != nil {
				return err
			}
			br.cp[i] = value

		case 3: // BOOLEAN
			var b byte
			if err := binary.Read(br.r, binary.BigEndian, &b); err != nil {
				return err
			}
			br.cp[i] = b != 0

		case 4: // STRING
			var length int32
			if err := binary.Read(br.r, binary.BigEndian, &length); err != nil {
				return err
			}
			if length < 0 {
				br.cp[i] = (*string)(nil)
			} else {
				strBytes := make([]byte, length)
				if _, err := br.r.Read(strBytes); err != nil {
					return err
				}
				str := string(strBytes)
				br.cp[i] = str
			}

		case 5: // PACKAGE
			var orgIdx, pkgNameIdx, moduleNameIdx, versionIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &orgIdx); err != nil {
				return err
			}
			if err := binary.Read(br.r, binary.BigEndian, &pkgNameIdx); err != nil {
				return err
			}
			if err := binary.Read(br.r, binary.BigEndian, &moduleNameIdx); err != nil {
				return err
			}
			if err := binary.Read(br.r, binary.BigEndian, &versionIdx); err != nil {
				return err
			}

			org := model.Name(br.getStringFromCP(int(orgIdx)))
			pkgName := model.Name(br.getStringFromCP(int(pkgNameIdx)))
			moduleName := model.Name(br.getStringFromCP(int(moduleNameIdx)))
			version := model.Name(br.getStringFromCP(int(versionIdx)))

			br.cp[i] = &model.PackageID{
				OrgName: &org,
				PkgName: &pkgName,
				Name:    &moduleName,
				Version: &version,
			}

		case 6: // BYTE
			var value int32
			if err := binary.Read(br.r, binary.BigEndian, &value); err != nil {
				return err
			}
			br.cp[i] = value

		case 7: // SHAPE (type)
			var shapeLen int32
			if err := binary.Read(br.r, binary.BigEndian, &shapeLen); err != nil {
				return err
			}

			// tag
			var tag uint8
			if err := binary.Read(br.r, binary.BigEndian, &tag); err != nil {
				return err
			}

			// name
			var nameIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &nameIdx); err != nil {
				return err
			}
			name := model.Name(br.getStringFromCP(int(nameIdx)))

			// flags
			var flags int64
			if err := binary.Read(br.r, binary.BigEndian, &flags); err != nil {
				return err
			}

			br.cp[i] = minimalBType{
				tag:   int(tag),
				name:  name,
				flags: flags,
			}
		default:
			return fmt.Errorf("unknown CP tag: %d at entry %d", tag, i)
		}
	}

	return nil
}

func (br *BIRReader) getStringFromCP(index int) string {
	if index < 0 || index >= len(br.cp) {
		return ""
	}

	if str, ok := br.cp[index].(string); ok {
		return str
	}
	return ""
}

func (r *BIRReader) getPackageFromCP(index int) *model.PackageID {
	if index < 0 || index >= len(r.cp) {
		return nil
	}
	if pkg, ok := r.cp[index].(*model.PackageID); ok {
		return pkg
	}
	return nil
}

func (r *BIRReader) getTypeFromCP(index int) *minimalBType {
	if index < 0 || index >= len(r.cp) {
		return nil
	}
	if t, ok := r.cp[index].(minimalBType); ok {
		return &t
	}
	return nil
}

func (r *BIRReader) getIntegerFromCP(index int) int64 {
	if index < 0 || index >= len(r.cp) {
		return 0
	}
	if val, ok := r.cp[index].(int64); ok {
		return val
	}
	return 0
}

func (br *BIRReader) readImports() ([]BIRImportModule, error) {
	var count int32
	if err := binary.Read(br.r, binary.BigEndian, &count); err != nil {
		return nil, err
	}

	imports := make([]BIRImportModule, count)
	for i := 0; i < int(count); i++ {
		var orgIdx, pkgNameIdx, moduleNameIdx, versionIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &orgIdx); err != nil {
			return nil, err
		}
		if err := binary.Read(br.r, binary.BigEndian, &pkgNameIdx); err != nil {
			return nil, err
		}
		if err := binary.Read(br.r, binary.BigEndian, &moduleNameIdx); err != nil {
			return nil, err
		}
		if err := binary.Read(br.r, binary.BigEndian, &versionIdx); err != nil {
			return nil, err
		}

		org := model.Name(br.getStringFromCP(int(orgIdx)))
		pkgName := model.Name(br.getStringFromCP(int(pkgNameIdx)))
		moduleName := model.Name(br.getStringFromCP(int(moduleNameIdx)))
		version := model.Name(br.getStringFromCP(int(versionIdx)))

		imports[i] = BIRImportModule{
			PackageID: &model.PackageID{
				OrgName: &org,
				PkgName: &pkgName,
				Name:    &moduleName,
				Version: &version,
			},
		}
	}

	return imports, nil
}

func (br *BIRReader) readConstants() ([]BIRConstant, error) {
	var count int32
	if err := binary.Read(br.r, binary.BigEndian, &count); err != nil {
		return nil, err
	}

	constants := make([]BIRConstant, count)
	for i := 0; i < int(count); i++ {
		var nameIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &nameIdx); err != nil {
			return nil, err
		}

		name := model.Name(br.getStringFromCP(int(nameIdx)))

		var flags int64
		if err := binary.Read(br.r, binary.BigEndian, &flags); err != nil {
			return nil, err
		}

		var origin uint8
		if err := binary.Read(br.r, binary.BigEndian, &origin); err != nil {
			return nil, err
		}

		constant := BIRConstant{
			Name:   name,
			Flags:  flags,
			Origin: model.SymbolOrigin(origin),
		}

		var typeIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &typeIdx); err != nil {
			return nil, err
		}

		t := br.getTypeFromCP(int(typeIdx))
		constant.Type = *t

		var length int64
		if err := binary.Read(br.r, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		var cTypeIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &cTypeIdx); err != nil {
			return nil, err
		}

		cv := br.getTypeFromCP(int(cTypeIdx))
		switch model.TypeTags(cv.tag) {
		case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
			var valueIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &valueIdx); err != nil {
				return nil, err
			}
			constant.ConstValue = ConstValue{
				Type:  *cv,
				Value: br.getIntegerFromCP(int(valueIdx)),
			}
		}

		constants[i] = constant
	}

	return constants, nil
}

func (br *BIRReader) readGlobalVars() ([]BIRGlobalVariableDcl, error) {
	var count int32
	if err := binary.Read(br.r, binary.BigEndian, &count); err != nil {
		return nil, err
	}

	variables := make([]BIRGlobalVariableDcl, count)
	for i := 0; i < int(count); i++ {
		var kind uint8
		if err := binary.Read(br.r, binary.BigEndian, &kind); err != nil {
			return nil, err
		}

		var nameIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &nameIdx); err != nil {
			return nil, err
		}

		name := model.Name(br.getStringFromCP(int(nameIdx)))

		var flags int64
		if err := binary.Read(br.r, binary.BigEndian, &flags); err != nil {
			return nil, err
		}

		var origin uint8
		if err := binary.Read(br.r, binary.BigEndian, &origin); err != nil {
			return nil, err
		}

		var typeIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &typeIdx); err != nil {
			return nil, err
		}

		fmt.Printf("Global Var Type Index: %d\n", typeIdx)

		t := br.getTypeFromCP(int(typeIdx))

		variables[i] = BIRGlobalVariableDcl{
			BIRVariableDcl: BIRVariableDcl{
				Kind: VarKind(kind),
				Name: name,
				Type: *t,
			},
			Flags:  flags,
			Origin: model.SymbolOrigin(origin),
		}
	}

	return variables, nil
}

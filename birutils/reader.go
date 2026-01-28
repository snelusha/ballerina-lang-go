package birutils

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
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

func (br *BIRReader) LoadBIRPackage() (*bir.BIRPackage, error) {
	return br.readPackage()
}

func (br *BIRReader) readPackage() (*bir.BIRPackage, error) {
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

	functions, err := br.readFunctions()

	return &bir.BIRPackage{
		PackageID:     pkgID,
		ImportModules: imports,
		Constants:     constants,
		GlobalVars:    globalVars,
		Functions:     functions,
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

			t := ast.NewBType(model.TypeTags(tag), nil, name, uint64(flags))

			// FIXME: Revisit this
			br.cp[i] = *t.(*ast.BTypeImpl)
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

func (r *BIRReader) getTypeFromCP(index int) ast.BType {
	if index < 0 || index >= len(r.cp) {
		return nil
	}
	if t, ok := r.cp[index].(ast.BTypeImpl); ok {
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

func (br *BIRReader) readImports() ([]bir.BIRImportModule, error) {
	var count int32
	if err := binary.Read(br.r, binary.BigEndian, &count); err != nil {
		return nil, err
	}

	imports := make([]bir.BIRImportModule, count)
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

		imports[i] = bir.BIRImportModule{
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

func (br *BIRReader) readConstants() ([]bir.BIRConstant, error) {
	var count int32
	if err := binary.Read(br.r, binary.BigEndian, &count); err != nil {
		return nil, err
	}

	constants := make([]bir.BIRConstant, count)
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

		constant := bir.BIRConstant{
			Name:   name,
			Flags:  flags,
			Origin: model.SymbolOrigin(origin),
		}

		var typeIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &typeIdx); err != nil {
			return nil, err
		}

		t := br.getTypeFromCP(int(typeIdx))
		constant.Type = t

		var length int64
		if err := binary.Read(br.r, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		var cTypeIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &cTypeIdx); err != nil {
			return nil, err
		}

		cv := br.getTypeFromCP(int(cTypeIdx))
		switch model.TypeTags(cv.BTypeGetTag()) {
		case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
			var valueIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &valueIdx); err != nil {
				return nil, err
			}
			constant.ConstValue = bir.ConstValue{
				Type:  cv,
				Value: br.getIntegerFromCP(int(valueIdx)),
			}
		}

		constants[i] = constant
	}

	return constants, nil
}

func (br *BIRReader) readGlobalVars() ([]bir.BIRGlobalVariableDcl, error) {
	var count int32
	if err := binary.Read(br.r, binary.BigEndian, &count); err != nil {
		return nil, err
	}

	variables := make([]bir.BIRGlobalVariableDcl, count)
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

		t := br.getTypeFromCP(int(typeIdx))

		variables[i] = bir.BIRGlobalVariableDcl{
			BIRVariableDcl: bir.BIRVariableDcl{
				Kind: bir.VarKind(kind),
				Name: name,
				Type: t,
			},
			Flags:  flags,
			Origin: model.SymbolOrigin(origin),
		}
	}

	return variables, nil
}

func (br *BIRReader) readFunctions() ([]bir.BIRFunction, error) {
	var count int32
	if err := binary.Read(br.r, binary.BigEndian, &count); err != nil {
		return nil, err
	}

	functions := make([]bir.BIRFunction, count)
	for i := 0; i < int(count); i++ {
		var nameIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &nameIdx); err != nil {
			return nil, err
		}

		name := model.Name(br.getStringFromCP(int(nameIdx)))

		var originalNameIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &originalNameIdx); err != nil {
			return nil, err
		}

		originalName := model.Name(br.getStringFromCP(int(originalNameIdx)))

		var flag int64
		if err := binary.Read(br.r, binary.BigEndian, &flag); err != nil {
			return nil, err
		}

		var origin uint8
		if err := binary.Read(br.r, binary.BigEndian, &origin); err != nil {
			return nil, err
		}

		function := bir.BIRFunction{
			Name:         name,
			OriginalName: originalName,
			Flags:        flag,
			Origin:       model.SymbolOrigin(origin),
		}

		// var typeIdx int32
		// if err := binary.Read(br.r, binary.BigEndian, &typeIdx); err != nil {
		// 	return nil, err
		// }

		// t := br.getTypeFromCP(int(typeIdx))
		// fmt.Printf("Function %s type: %+v\n", name, t)

		var requiredParamsCount int32
		if err := binary.Read(br.r, binary.BigEndian, &requiredParamsCount); err != nil {
			return nil, err
		}
		fmt.Printf("[READER] REQUIRED PARAMS COUNT: %d IN FUNCTION %s\n", requiredParamsCount, name)

		requiredParams := make([]bir.BIRParameter, requiredParamsCount)
		for j := 0; j < int(requiredParamsCount); j++ {
			var paramNameIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &paramNameIdx); err != nil {
				return nil, err
			}
			paramName := model.Name(br.getStringFromCP(int(paramNameIdx)))
			fmt.Printf("[READER] PARAM NAME: %s\n", paramName)

			var paramFlags int64
			if err := binary.Read(br.r, binary.BigEndian, &paramFlags); err != nil {
				return nil, err
			}

			requiredParams[j] = bir.BIRParameter{
				Name:  paramName,
				Flags: paramFlags,
			}
		}
		function.RequiredParams = requiredParams

		var length int64
		if err := binary.Read(br.r, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		var argsCount int32
		if err := binary.Read(br.r, binary.BigEndian, &argsCount); err != nil {
			return nil, err
		}
		function.ArgsCount = int(argsCount)

		var hasReturnVar bool
		if err := binary.Read(br.r, binary.BigEndian, &hasReturnVar); err != nil {
			return nil, err
		}

		if hasReturnVar {
			var returnVarKind uint8
			if err := binary.Read(br.r, binary.BigEndian, &returnVarKind); err != nil {
				return nil, err
			}

			var returnVarTypeIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &returnVarTypeIdx); err != nil {
				return nil, err
			}
			returnVarType := br.getTypeFromCP(int(returnVarTypeIdx))

			var returnVarNameIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &returnVarNameIdx); err != nil {
				return nil, err
			}
			returnVarName := model.Name(br.getStringFromCP(int(returnVarNameIdx)))

			function.ReturnVariable = &bir.BIRVariableDcl{
				Kind: bir.VarKind(returnVarKind),
				Name: returnVarName,
				Type: returnVarType,
			}
		}

		// var localVarCount int32
		// if err := binary.Read(br.r, binary.BigEndian, &localVarCount); err != nil {
		// 	return nil, err
		// }

		// localVars := make([]bir.BIRVariableDcl, localVarCount)
		// for range localVars {
		// 	var localVarKind uint8
		// 	if err := binary.Read(br.r, binary.BigEndian, &localVarKind); err != nil {
		// 		return nil, err
		// 	}
		//
		// 	var localVarTypeIdx int32
		// 	if err := binary.Read(br.r, binary.BigEndian, &localVarTypeIdx); err != nil {
		// 		return nil, err
		// 	}
		// 	localVarType := br.getTypeFromCP(int(localVarTypeIdx))
		//
		// 	var localVarNameIdx int32
		// 	if err := binary.Read(br.r, binary.BigEndian, &localVarNameIdx); err != nil {
		// 		return nil, err
		// 	}
		// 	localVarName := model.Name(br.getStringFromCP(int(localVarNameIdx)))
		//
		// 	localVar := bir.BIRVariableDcl{
		// 		Kind: bir.VarKind(localVarKind),
		// 		Name: localVarName,
		// 		Type: localVarType,
		// 	}

		// fmt.Printf("[READER] LOCAL VAR: %+v\n", localVar)

		// if localVarKind == uint8(bir.VAR_KIND_ARG) {
		// 	fmt.Printf("[READER] KIND ARG\n")
		//
		// 	var metaVarNameIdx int32
		// 	if err := binary.Read(br.r, binary.BigEndian, &metaVarNameIdx); err != nil {
		// 		return nil, err
		// 	}
		// 	metaVarName := br.getStringFromCP(int(metaVarNameIdx))
		// 	localVar.MetaVarName = metaVarName
		// }

		// localVars = append(localVars, localVar)
		// }
		// function.LocalVars = localVars

		functions[i] = function
	}

	return functions, nil
}

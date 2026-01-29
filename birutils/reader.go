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

func (r *BIRReader) getByteFromCP(index int) uint8 {
	if index < 0 || index >= len(r.cp) {
		return 0
	}
	if val, ok := r.cp[index].(uint8); ok {
		return val
	}
	return 0
}

func (r *BIRReader) getFloatFromCP(index int) float64 {
	if index < 0 || index >= len(r.cp) {
		return 0
	}
	if val, ok := r.cp[index].(float64); ok {
		return val
	}
	return 0
}

func (r *BIRReader) getBooleanFromCP(index int) bool {
	if index < 0 || index >= len(r.cp) {
		return false
	}
	if val, ok := r.cp[index].(bool); ok {
		return val
	}
	return false
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

		requiredParams := make([]bir.BIRParameter, requiredParamsCount)
		for j := 0; j < int(requiredParamsCount); j++ {
			var paramNameIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &paramNameIdx); err != nil {
				return nil, err
			}
			paramName := model.Name(br.getStringFromCP(int(paramNameIdx)))

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

		var localVarCount int32
		if err := binary.Read(br.r, binary.BigEndian, &localVarCount); err != nil {
			return nil, err
		}

		localVars := make([]bir.BIRVariableDcl, localVarCount)

		for j := 0; j < int(localVarCount); j++ {
			var localVarKind uint8
			if err := binary.Read(br.r, binary.BigEndian, &localVarKind); err != nil {
				return nil, err
			}

			var localVarTypeIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &localVarTypeIdx); err != nil {
				return nil, err
			}
			localVarType := br.getTypeFromCP(int(localVarTypeIdx))

			var localVarNameIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &localVarNameIdx); err != nil {
				return nil, err
			}
			localVarName := model.Name(br.getStringFromCP(int(localVarNameIdx)))

			localVar := bir.BIRVariableDcl{
				Kind: bir.VarKind(localVarKind),
				Name: localVarName,
				Type: localVarType,
			}

			if localVarKind == uint8(bir.VAR_KIND_ARG) {
				var metaVarNameIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &metaVarNameIdx); err != nil {
					return nil, err
				}
				metaVarName := br.getStringFromCP(int(metaVarNameIdx))
				localVar.MetaVarName = metaVarName
			} else if localVarKind == uint8(bir.VAR_KIND_LOCAL) {
				var metaVarNameIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &metaVarNameIdx); err != nil {
					return nil, err
				}
				metaVarName := br.getStringFromCP(int(metaVarNameIdx))
				localVar.MetaVarName = metaVarName

				var endBBIdIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &endBBIdIdx); err != nil {
					return nil, err
				}
				endBBId := model.Name(br.getStringFromCP(int(endBBIdIdx)))
				localVar.EndBB = &bir.BIRBasicBlock{
					Id: endBBId,
				}

				var startBBIdIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &startBBIdIdx); err != nil {
					return nil, err
				}
				startBBId := model.Name(br.getStringFromCP(int(startBBIdIdx)))
				localVar.StartBB = &bir.BIRBasicBlock{
					Id: startBBId,
				}

				var insOffset int32
				if err := binary.Read(br.r, binary.BigEndian, &insOffset); err != nil {
					return nil, err
				}
				localVar.InsOffset = int(insOffset)
			}

			localVars[j] = localVar
		}
		function.LocalVars = localVars

		var basicBlockCount int32
		if err := binary.Read(br.r, binary.BigEndian, &basicBlockCount); err != nil {
			return nil, err
		}

		basicBlocks := make([]bir.BIRBasicBlock, basicBlockCount)
		for j := 0; j < int(basicBlockCount); j++ {
			var idIdx int32
			if err := binary.Read(br.r, binary.BigEndian, &idIdx); err != nil {
				return nil, err
			}
			id := model.Name(br.getStringFromCP(int(idIdx)))

			block := bir.BIRBasicBlock{
				Id: id,
			}

			var instructionCount int32
			if err := binary.Read(br.r, binary.BigEndian, &instructionCount); err != nil {
				return nil, err
			}

			instructions := make([]bir.BIRInstruction, instructionCount)
			for k := 0; k < int(instructionCount); k++ {
				var insKind uint8
				if err := binary.Read(br.r, binary.BigEndian, &insKind); err != nil {
					return nil, err
				}
				instructionKind := bir.InstructionKind(insKind)

				switch instructionKind {
				case bir.INSTRUCTION_KIND_MOVE:
					rhsOp, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					lhsOp, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					instructions[k] = &bir.Move{
						BIRInstructionBase: bir.BIRInstructionBase{
							LhsOp: lhsOp,
						},
						RhsOp: rhsOp,
					}
				case bir.INSTRUCTION_KIND_ADD, bir.INSTRUCTION_KIND_SUB, bir.INSTRUCTION_KIND_MUL, bir.INSTRUCTION_KIND_DIV, bir.INSTRUCTION_KIND_MOD, bir.INSTRUCTION_KIND_EQUAL, bir.INSTRUCTION_KIND_NOT_EQUAL, bir.INSTRUCTION_KIND_GREATER_THAN, bir.INSTRUCTION_KIND_GREATER_EQUAL, bir.INSTRUCTION_KIND_LESS_THAN, bir.INSTRUCTION_KIND_LESS_EQUAL, bir.INSTRUCTION_KIND_AND, bir.INSTRUCTION_KIND_OR, bir.INSTRUCTION_KIND_REF_EQUAL, bir.INSTRUCTION_KIND_REF_NOT_EQUAL, bir.INSTRUCTION_KIND_CLOSED_RANGE, bir.INSTRUCTION_KIND_HALF_OPEN_RANGE, bir.INSTRUCTION_KIND_ANNOT_ACCESS, bir.INSTRUCTION_KIND_BITWISE_AND, bir.INSTRUCTION_KIND_BITWISE_OR, bir.INSTRUCTION_KIND_BITWISE_XOR, bir.INSTRUCTION_KIND_BITWISE_LEFT_SHIFT, bir.INSTRUCTION_KIND_BITWISE_RIGHT_SHIFT, bir.INSTRUCTION_KIND_BITWISE_UNSIGNED_RIGHT_SHIFT:
					rhsOp1, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					rhsOp2, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					lhsOp, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					instructions[k] = &bir.BinaryOp{
						BIRInstructionBase: bir.BIRInstructionBase{
							LhsOp: lhsOp,
						},
						Kind:   instructionKind,
						RhsOp1: *rhsOp1,
						RhsOp2: *rhsOp2,
					}
				case bir.INSTRUCTION_KIND_TYPEOF, bir.INSTRUCTION_KIND_NOT, bir.INSTRUCTION_KIND_NEGATE:
					rhsOp, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					lhsOp, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					instructions[k] = &bir.UnaryOp{
						BIRInstructionBase: bir.BIRInstructionBase{
							LhsOp: lhsOp,
						},
						Kind:  instructionKind,
						RhsOp: rhsOp,
					}
				case bir.INSTRUCTION_KIND_CONST_LOAD:
					var constLoadTypeIdx int32
					if err := binary.Read(br.r, binary.BigEndian, &constLoadTypeIdx); err != nil {
						return nil, err
					}
					constLoadType := br.getTypeFromCP(int(constLoadTypeIdx))

					lhsOp, err := br.readOperand()
					if err != nil {
						return nil, err
					}

					var value any

					switch constLoadType.BTypeGetTag() {
					case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
						var valueIdx int32
						if err := binary.Read(br.r, binary.BigEndian, &valueIdx); err != nil {
							return nil, err
						}
						value = br.getIntegerFromCP(int(valueIdx))
					case model.TypeTags_BYTE:
						var valueIdx int32
						if err := binary.Read(br.r, binary.BigEndian, &valueIdx); err != nil {
							return nil, err
						}
						value = br.getByteFromCP(int(valueIdx))
					case model.TypeTags_FLOAT:
						var valueIdx int32
						if err := binary.Read(br.r, binary.BigEndian, &valueIdx); err != nil {
							return nil, err
						}
						value = br.getFloatFromCP(int(valueIdx))
					case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
						var valueIdx int32
						if err := binary.Read(br.r, binary.BigEndian, &valueIdx); err != nil {
							return nil, err
						}
						value = br.getStringFromCP(int(valueIdx))
					case model.TypeTags_BOOLEAN:
						var valueIdx int32
						if err := binary.Read(br.r, binary.BigEndian, &valueIdx); err != nil {
							return nil, err
						}
						value = br.getBooleanFromCP(int(valueIdx))
					}

					instructions[k] = &bir.ConstantLoad{
						BIRInstructionBase: bir.BIRInstructionBase{
							LhsOp: lhsOp,
						},
						Type:  constLoadType,
						Value: value,
					}
				}
			}
			block.Instructions = instructions

			var terminatorKind uint8
			if err := binary.Read(br.r, binary.BigEndian, &terminatorKind); err != nil {
				return nil, err
			}

			termInstructionKind := bir.InstructionKind(terminatorKind)

			switch termInstructionKind {
			case bir.INSTRUCTION_KIND_RETURN:
				fmt.Println("Reading RETURN")
				block.Terminator = &bir.Return{}
			case bir.INSTRUCTION_KIND_GOTO:
				fmt.Println("Reading GOTO")
				var idIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &idIdx); err != nil {
					return nil, err
				}
				id := br.getStringFromCP(int(idIdx))
				block.Terminator = &bir.Goto{
					BIRTerminatorBase: bir.BIRTerminatorBase{
						ThenBB: &bir.BIRBasicBlock{
							Id: model.Name(id),
						},
					},
				}
			case bir.INSTRUCTION_KIND_BRANCH:
				fmt.Println("Reading BRANCH")

				op, err := br.readOperand()
				if err != nil {
					return nil, err
				}

				var trueBBIdIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &trueBBIdIdx); err != nil {
					return nil, err
				}
				trueBBId := br.getStringFromCP(int(trueBBIdIdx))

				var falseBBIdIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &falseBBIdIdx); err != nil {
					return nil, err
				}
				falseBBId := br.getStringFromCP(int(falseBBIdIdx))

				block.Terminator = &bir.Branch{
					Op: op,
					TrueBB: &bir.BIRBasicBlock{
						Id: model.Name(trueBBId),
					},
					FalseBB: &bir.BIRBasicBlock{
						Id: model.Name(falseBBId),
					},
				}
			case bir.INSTRUCTION_KIND_CALL:
				var isVirtual bool
				if err := binary.Read(br.r, binary.BigEndian, &isVirtual); err != nil {
					return nil, err
				}

				var pkgIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &pkgIdx); err != nil {
					return nil, err
				}
				pkg := br.getPackageFromCP(int(pkgIdx))

				var nameIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &nameIdx); err != nil {
					return nil, err
				}
				name := br.getStringFromCP(int(nameIdx))

				var argsCount int32
				if err := binary.Read(br.r, binary.BigEndian, &argsCount); err != nil {
					return nil, err
				}

				args := make([]bir.BIROperand, argsCount)
				for k := 0; k < int(argsCount); k++ {
					arg, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					args[k] = *arg
				}

				var lshOpExists bool
				if err := binary.Read(br.r, binary.BigEndian, &lshOpExists); err != nil {
					return nil, err
				}

				var lhsOp *bir.BIROperand
				if lshOpExists {
					op, err := br.readOperand()
					if err != nil {
						return nil, err
					}
					lhsOp = op
				}

				var thenBBIdIdx int32
				if err := binary.Read(br.r, binary.BigEndian, &thenBBIdIdx); err != nil {
					return nil, err
				}
				thenBBId := br.getStringFromCP(int(thenBBIdIdx))

				block.Terminator = &bir.Call{
					Kind:      termInstructionKind,
					IsVirtual: isVirtual,
					CalleePkg: pkg,
					Name:      model.Name(name),
					Args:      args,
					BIRTerminatorBase: bir.BIRTerminatorBase{
						ThenBB: &bir.BIRBasicBlock{
							Id: model.Name(thenBBId),
						},
						BIRInstructionBase: bir.BIRInstructionBase{
							LhsOp: lhsOp,
						},
					},
				}
			}

			basicBlocks[j] = block
		}
		function.BasicBlocks = basicBlocks

		functions[i] = function
	}

	return functions, nil
}

func (br *BIRReader) readOperand() (*bir.BIROperand, error) {
	var ignoreVariable bool
	if err := binary.Read(br.r, binary.BigEndian, &ignoreVariable); err != nil {
		return nil, err
	}

	if ignoreVariable {
		var varTypeIdx int32
		if err := binary.Read(br.r, binary.BigEndian, &varTypeIdx); err != nil {
			return nil, err
		}
		varType := br.getTypeFromCP(int(varTypeIdx))
		return &bir.BIROperand{
			VariableDcl: &bir.BIRVariableDcl{
				Type: varType,
			},
		}, nil
	}

	var varKind uint8
	if err := binary.Read(br.r, binary.BigEndian, &varKind); err != nil {
		return nil, err
	}

	var scope uint8
	if err := binary.Read(br.r, binary.BigEndian, &scope); err != nil {
		return nil, err
	}

	var nameIdx int32
	if err := binary.Read(br.r, binary.BigEndian, &nameIdx); err != nil {
		return nil, err
	}
	name := model.Name(br.getStringFromCP(int(nameIdx)))

	return &bir.BIROperand{
		VariableDcl: &bir.BIRVariableDcl{
			Kind:  bir.VarKind(varKind),
			Scope: bir.VarScope(scope),
			Name:  name,
		},
	}, nil
}

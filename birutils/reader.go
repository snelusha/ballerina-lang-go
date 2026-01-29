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

	version, err := br.readInt32()
	if err != nil {
		return nil, err
	}

	if version != BIR_VERSION {
		return nil, fmt.Errorf("unsupported BIR version: %d", version)
	}

	if err := br.readConstantPool(); err != nil {
		return nil, fmt.Errorf("reading constant pool: %w", err)
	}

	pkgCPIndex, err := br.readInt32()
	if err != nil {
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

func (br *BIRReader) readInt8() (int8, error) {
	var v int8
	if err := binary.Read(br.r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (br *BIRReader) readUInt8() (uint8, error) {
	var v uint8
	if err := binary.Read(br.r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (br *BIRReader) readInt32() (int32, error) {
	var v int32
	if err := binary.Read(br.r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (br *BIRReader) readInt64() (int64, error) {
	var v int64
	if err := binary.Read(br.r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (br *BIRReader) readFloat64() (float64, error) {
	var v float64
	if err := binary.Read(br.r, binary.BigEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (br *BIRReader) readBool() (bool, error) {
	var v bool
	if err := binary.Read(br.r, binary.BigEndian, &v); err != nil {
		return false, err
	}
	return v, nil
}

func (br *BIRReader) readConstantPool() error {
	cpSize, err := br.readInt32()
	if err != nil {
		return err
	}

	br.cp = make([]any, cpSize)

	for i := 0; i < int(cpSize); i++ {
		tag, err := br.readInt8()
		if err != nil {
			return fmt.Errorf("reading CP entry %d tag: %w", i, err)
		}

		if err := br.readConstantPoolEntry(tag, i); err != nil {
			return fmt.Errorf("reading CP entry %d (tag %d): %w", i, tag, err)
		}
	}

	return nil
}

func (br *BIRReader) readConstantPoolEntry(tag int8, i int) error {
	switch tag {
	case 0: // NULL/placeholder entry
		br.cp[i] = nil
	case 1: // INTEGER
		value, err := br.readInt64()
		if err != nil {
			return err
		}
		br.cp[i] = value

	case 2: // FLOAT
		value, err := br.readFloat64()
		if err != nil {
			return err
		}
		br.cp[i] = value

	case 3: // BOOLEAN
		b, err := br.readUInt8()
		if err != nil {
			return err
		}
		br.cp[i] = b != 0

	case 4: // STRING
		length, err := br.readInt32()
		if err != nil {
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
		orgIdx, err := br.readInt32()
		if err != nil {
			return err
		}
		pkgNameIdx, err := br.readInt32()
		if err != nil {
			return err
		}
		moduleNameIdx, err := br.readInt32()
		if err != nil {
			return err
		}
		versionIdx, err := br.readInt32()
		if err != nil {
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
		value, err := br.readInt32()
		if err != nil {
			return err
		}
		br.cp[i] = value

	case 7: // SHAPE (type)
		_, err := br.readInt32() // shapeLen
		if err != nil {
			return err
		}

		// tag
		tag, err := br.readUInt8()
		if err != nil {
			return err
		}

		// name
		nameIdx, err := br.readInt32()
		if err != nil {
			return err
		}
		name := model.Name(br.getStringFromCP(int(nameIdx)))

		// flags
		flags, err := br.readInt64()
		if err != nil {
			return err
		}

		t := ast.NewBType(model.TypeTags(tag), nil, name, uint64(flags))

		// FIXME: Revisit this
		br.cp[i] = *t.(*ast.BTypeImpl)
	default:
		return fmt.Errorf("unknown CP tag: %d", tag)
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
	count, err := br.readInt32()
	if err != nil {
		return nil, err
	}

	imports := make([]bir.BIRImportModule, count)
	for i := 0; i < int(count); i++ {
		orgIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		pkgNameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		moduleNameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		versionIdx, err := br.readInt32()
		if err != nil {
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
	count, err := br.readInt32()
	if err != nil {
		return nil, err
	}

	constants := make([]bir.BIRConstant, count)
	for i := 0; i < int(count); i++ {
		nameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}

		name := model.Name(br.getStringFromCP(int(nameIdx)))

		flags, err := br.readInt64()
		if err != nil {
			return nil, err
		}

		origin, err := br.readUInt8()
		if err != nil {
			return nil, err
		}

		constant := bir.BIRConstant{
			Name:   name,
			Flags:  flags,
			Origin: model.SymbolOrigin(origin),
		}

		typeIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}

		t := br.getTypeFromCP(int(typeIdx))
		constant.Type = t

		_, err = br.readInt64() // length
		if err != nil {
			return nil, err
		}

		cTypeIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}

		cv := br.getTypeFromCP(int(cTypeIdx))
		switch model.TypeTags(cv.BTypeGetTag()) {
		case model.TypeTags_INT, model.TypeTags_SIGNED32_INT, model.TypeTags_SIGNED16_INT, model.TypeTags_SIGNED8_INT, model.TypeTags_UNSIGNED32_INT, model.TypeTags_UNSIGNED16_INT, model.TypeTags_UNSIGNED8_INT:
			valueIdx, err := br.readInt32()
			if err != nil {
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
	count, err := br.readInt32()
	if err != nil {
		return nil, err
	}

	variables := make([]bir.BIRGlobalVariableDcl, count)
	for i := 0; i < int(count); i++ {
		kind, err := br.readUInt8()
		if err != nil {
			return nil, err
		}

		nameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}

		name := model.Name(br.getStringFromCP(int(nameIdx)))

		flags, err := br.readInt64()
		if err != nil {
			return nil, err
		}

		origin, err := br.readUInt8()
		if err != nil {
			return nil, err
		}

		typeIdx, err := br.readInt32()
		if err != nil {
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
	count, err := br.readInt32()
	if err != nil {
		return nil, err
	}

	functions := make([]bir.BIRFunction, count)
	for i := 0; i < int(count); i++ {
		fn, err := br.readFunction()
		if err != nil {
			return nil, err
		}
		functions[i] = *fn
	}

	return functions, nil
}

func (br *BIRReader) readFunction() (*bir.BIRFunction, error) {
	nameIdx, err := br.readInt32()
	if err != nil {
		return nil, err
	}
	name := model.Name(br.getStringFromCP(int(nameIdx)))

	originalNameIdx, err := br.readInt32()
	if err != nil {
		return nil, err
	}
	originalName := model.Name(br.getStringFromCP(int(originalNameIdx)))

	flag, err := br.readInt64()
	if err != nil {
		return nil, err
	}

	origin, err := br.readUInt8()
	if err != nil {
		return nil, err
	}

	requiredParamsCount, err := br.readInt32()
	if err != nil {
		return nil, err
	}

	requiredParams := make([]bir.BIRParameter, requiredParamsCount)
	for j := 0; j < int(requiredParamsCount); j++ {
		paramNameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		paramName := model.Name(br.getStringFromCP(int(paramNameIdx)))

		paramFlags, err := br.readInt64()
		if err != nil {
			return nil, err
		}

		requiredParams[j] = bir.BIRParameter{
			Name:  paramName,
			Flags: paramFlags,
		}
	}

	_, err = br.readInt64() // length, unused?
	if err != nil {
		return nil, err
	}

	argsCount, err := br.readInt32()
	if err != nil {
		return nil, err
	}

	hasReturnVar, err := br.readBool()
	if err != nil {
		return nil, err
	}

	var returnVar *bir.BIRVariableDcl
	if hasReturnVar {
		returnVarKind, err := br.readUInt8()
		if err != nil {
			return nil, err
		}
		returnVarTypeIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		returnVarType := br.getTypeFromCP(int(returnVarTypeIdx))

		returnVarNameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		returnVarName := model.Name(br.getStringFromCP(int(returnVarNameIdx)))

		returnVar = &bir.BIRVariableDcl{
			Kind: bir.VarKind(returnVarKind),
			Name: returnVarName,
			Type: returnVarType,
		}
	}

	localVarCount, err := br.readInt32()
	if err != nil {
		return nil, err
	}
	localVars := make([]bir.BIRVariableDcl, localVarCount)
	for j := 0; j < int(localVarCount); j++ {
		localVar, err := br.readLocalVar()
		if err != nil {
			return nil, err
		}
		localVars[j] = *localVar
	}

	basicBlockCount, err := br.readInt32()
	if err != nil {
		return nil, err
	}
	basicBlocks := make([]bir.BIRBasicBlock, basicBlockCount)
	for j := 0; j < int(basicBlockCount); j++ {
		block, err := br.readBasicBlock()
		if err != nil {
			return nil, err
		}
		basicBlocks[j] = *block
	}

	return &bir.BIRFunction{
		Name:           name,
		OriginalName:   originalName,
		Flags:          flag,
		Origin:         model.SymbolOrigin(origin),
		RequiredParams: requiredParams,
		ArgsCount:      int(argsCount),
		ReturnVariable: returnVar,
		LocalVars:      localVars,
		BasicBlocks:    basicBlocks,
	}, nil
}

func (br *BIRReader) readLocalVar() (*bir.BIRVariableDcl, error) {
	kind, err := br.readUInt8()
	if err != nil {
		return nil, err
	}
	typeIdx, err := br.readInt32()
	if err != nil {
		return nil, err
	}
	t := br.getTypeFromCP(int(typeIdx))

	nameIdx, err := br.readInt32()
	if err != nil {
		return nil, err
	}
	name := model.Name(br.getStringFromCP(int(nameIdx)))

	localVar := &bir.BIRVariableDcl{
		Kind: bir.VarKind(kind),
		Name: name,
		Type: t,
	}

	if kind == uint8(bir.VAR_KIND_ARG) {
		metaVarNameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		localVar.MetaVarName = br.getStringFromCP(int(metaVarNameIdx))
	} else if kind == uint8(bir.VAR_KIND_LOCAL) {
		metaVarNameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		localVar.MetaVarName = br.getStringFromCP(int(metaVarNameIdx))

		endBBIdIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		endBBId := model.Name(br.getStringFromCP(int(endBBIdIdx)))
		localVar.EndBB = &bir.BIRBasicBlock{Id: endBBId}

		startBBIdIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		startBBId := model.Name(br.getStringFromCP(int(startBBIdIdx)))
		localVar.StartBB = &bir.BIRBasicBlock{Id: startBBId}

		insOffset, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		localVar.InsOffset = int(insOffset)
	}
	return localVar, nil
}

func (br *BIRReader) readBasicBlock() (*bir.BIRBasicBlock, error) {
	idIdx, err := br.readInt32()
	if err != nil {
		return nil, err
	}
	id := model.Name(br.getStringFromCP(int(idIdx)))

	instructionCount, err := br.readInt32()
	if err != nil {
		return nil, err
	}

	instructions := make([]bir.BIRInstruction, instructionCount)
	for k := 0; k < int(instructionCount); k++ {
		ins, err := br.readInstruction()
		if err != nil {
			return nil, err
		}
		instructions[k] = ins
	}

	term, err := br.readTerminator()
	if err != nil {
		return nil, err
	}

	return &bir.BIRBasicBlock{
		Id:           id,
		Instructions: instructions,
		Terminator:   term,
	}, nil
}

func (br *BIRReader) readInstruction() (bir.BIRInstruction, error) {
	insKind, err := br.readUInt8()
	if err != nil {
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
		return &bir.Move{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			RhsOp: rhsOp,
		}, nil
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
		return &bir.BinaryOp{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			Kind:   instructionKind,
			RhsOp1: *rhsOp1,
			RhsOp2: *rhsOp2,
		}, nil
	case bir.INSTRUCTION_KIND_TYPEOF, bir.INSTRUCTION_KIND_NOT, bir.INSTRUCTION_KIND_NEGATE:
		rhsOp, err := br.readOperand()
		if err != nil {
			return nil, err
		}
		lhsOp, err := br.readOperand()
		if err != nil {
			return nil, err
		}
		return &bir.UnaryOp{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			Kind:  instructionKind,
			RhsOp: rhsOp,
		}, nil
	case bir.INSTRUCTION_KIND_CONST_LOAD:
		constLoadTypeIdx, err := br.readInt32()
		if err != nil {
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
			valueIdx, err := br.readInt32()
			if err != nil {
				return nil, err
			}
			value = br.getIntegerFromCP(int(valueIdx))
		case model.TypeTags_BYTE:
			valueIdx, err := br.readInt32()
			if err != nil {
				return nil, err
			}
			value = br.getByteFromCP(int(valueIdx))
		case model.TypeTags_FLOAT:
			valueIdx, err := br.readInt32()
			if err != nil {
				return nil, err
			}
			value = br.getFloatFromCP(int(valueIdx))
		case model.TypeTags_STRING, model.TypeTags_CHAR_STRING, model.TypeTags_DECIMAL:
			valueIdx, err := br.readInt32()
			if err != nil {
				return nil, err
			}
			value = br.getStringFromCP(int(valueIdx))
		case model.TypeTags_BOOLEAN:
			valueIdx, err := br.readInt32()
			if err != nil {
				return nil, err
			}
			value = br.getBooleanFromCP(int(valueIdx))
		}

		return &bir.ConstantLoad{
			BIRInstructionBase: bir.BIRInstructionBase{
				LhsOp: lhsOp,
			},
			Type:  constLoadType,
			Value: value,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported instruction kind: %d", instructionKind)
	}
}

func (br *BIRReader) readTerminator() (bir.BIRTerminator, error) {
	terminatorKind, err := br.readUInt8()
	if err != nil {
		return nil, err
	}
	termInstructionKind := bir.InstructionKind(terminatorKind)

	switch termInstructionKind {
	case bir.INSTRUCTION_KIND_RETURN:
		return &bir.Return{}, nil
	case bir.INSTRUCTION_KIND_GOTO:
		idIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		id := br.getStringFromCP(int(idIdx))
		return &bir.Goto{
			BIRTerminatorBase: bir.BIRTerminatorBase{
				ThenBB: &bir.BIRBasicBlock{
					Id: model.Name(id),
				},
			},
		}, nil
	case bir.INSTRUCTION_KIND_BRANCH:
		op, err := br.readOperand()
		if err != nil {
			return nil, err
		}

		trueBBIdIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		trueBBId := br.getStringFromCP(int(trueBBIdIdx))

		falseBBIdIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		falseBBId := br.getStringFromCP(int(falseBBIdIdx))

		return &bir.Branch{
			Op: op,
			TrueBB: &bir.BIRBasicBlock{
				Id: model.Name(trueBBId),
			},
			FalseBB: &bir.BIRBasicBlock{
				Id: model.Name(falseBBId),
			},
		}, nil
	case bir.INSTRUCTION_KIND_CALL:
		isVirtual, err := br.readBool()
		if err != nil {
			return nil, err
		}

		pkgIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		pkg := br.getPackageFromCP(int(pkgIdx))

		nameIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		name := br.getStringFromCP(int(nameIdx))

		argsCount, err := br.readInt32()
		if err != nil {
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

		lshOpExists, err := br.readBool()
		if err != nil {
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

		thenBBIdIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		thenBBId := br.getStringFromCP(int(thenBBIdIdx))

		return &bir.Call{
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
		}, nil
	default:
		return nil, fmt.Errorf("unsupported terminator kind: %d", termInstructionKind)
	}
}

func (br *BIRReader) readOperand() (*bir.BIROperand, error) {
	ignoreVariable, err := br.readBool()
	if err != nil {
		return nil, err
	}

	if ignoreVariable {
		varTypeIdx, err := br.readInt32()
		if err != nil {
			return nil, err
		}
		varType := br.getTypeFromCP(int(varTypeIdx))
		return &bir.BIROperand{
			VariableDcl: &bir.BIRVariableDcl{
				Type: varType,
			},
		}, nil
	}

	varKind, err := br.readUInt8()
	if err != nil {
		return nil, err
	}

	scope, err := br.readUInt8()
	if err != nil {
		return nil, err
	}

	nameIdx, err := br.readInt32()
	if err != nil {
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

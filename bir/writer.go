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

	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
)

const (
	// BIRMagic is the magic number for BIR files
	BIRMagic = "\xba\x10\xc0\xde"
	// BIRVersion is the version number for BIR files
	BIRVersion = 75
)

// BIRBinaryWriter serializes BIR into binary format
type BIRBinaryWriter struct {
	cp         *ConstantPool
	birPackage *BIRPackage
}

// NewBIRBinaryWriter creates a new BIR binary writer
func NewBIRBinaryWriter(birPackage *BIRPackage, typeEnv any) *BIRBinaryWriter {
	return &BIRBinaryWriter{
		birPackage: birPackage,
		cp:         NewConstantPool(typeEnv),
	}
}

// Serialize serializes the BIR package to binary format
func (w *BIRBinaryWriter) Serialize() ([]byte, error) {
	birbuf := &bytes.Buffer{}

	pkgCPIndex := AddPkgCPEntry(w.birPackage.PackageID, w.cp)
	if err := binary.Write(birbuf, binary.BigEndian, int32(pkgCPIndex)); err != nil {
		return nil, fmt.Errorf("writing package CP index: %w", err)
	}

	if err := w.writeImportModuleDecls(birbuf, w.birPackage.ImportModules); err != nil {
		return nil, err
	}
	if err := w.writeConstants(birbuf, w.birPackage.Constants); err != nil {
		return nil, err
	}
	if err := w.writeTypeDefs(birbuf, w.birPackage.TypeDefs); err != nil {
		return nil, err
	}
	if err := w.writeGlobalVars(birbuf, w.birPackage.GlobalVars); err != nil {
		return nil, err
	}
	if err := w.writeTypeDefBodies(birbuf, w.birPackage.TypeDefs); err != nil {
		return nil, err
	}
	if err := w.writeFunctions(birbuf, w.birPackage.Functions); err != nil {
		return nil, err
	}
	if err := w.writeAnnotations(birbuf, []any{}); err != nil {
		return nil, err
	}
	if err := w.writeServiceDeclarations(birbuf, []any{}); err != nil {
		return nil, err
	}

	w.cp.PrintEntries()

	cpBytes, err := w.cp.Serialize()
	if err != nil {
		return nil, fmt.Errorf("serializing constant pool: %w", err)
	}

	result := &bytes.Buffer{}
	if _, err := result.Write([]byte(BIRMagic)); err != nil {
		return nil, err
	}
	if err := binary.Write(result, binary.BigEndian, int32(BIRVersion)); err != nil {
		return nil, err
	}
	if _, err := result.Write(cpBytes); err != nil {
		return nil, err
	}
	if _, err := result.Write(birbuf.Bytes()); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (w *BIRBinaryWriter) writeImportModuleDecls(buf *bytes.Buffer, importModules []BIRImportModule) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(importModules))); err != nil {
		return err
	}

	for _, impMod := range importModules {
		if impMod.PackageID == nil {
			continue
		}

		orgCPIndex := AddStringCPEntry(getPackageIDOrgName(impMod.PackageID), w.cp)
		if err := binary.Write(buf, binary.BigEndian, int32(orgCPIndex)); err != nil {
			return err
		}

		pkgNameCPIndex := AddStringCPEntry(getPackageIDPkgName(impMod.PackageID), w.cp)
		if err := binary.Write(buf, binary.BigEndian, int32(pkgNameCPIndex)); err != nil {
			return err
		}

		nameCPIndex := AddStringCPEntry(getPackageIDName(impMod.PackageID), w.cp)
		if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
			return err
		}

		versionCPIndex := AddStringCPEntry(getPackageIDVersion(impMod.PackageID), w.cp)
		if err := binary.Write(buf, binary.BigEndian, int32(versionCPIndex)); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeConstants(buf *bytes.Buffer, constants []BIRConstant) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(constants))); err != nil {
		return err
	}

	for i, constant := range constants {
		if err := w.writeConstant(buf, constant); err != nil {
			return fmt.Errorf("writing constant %d (%s): %w", i, constant.Name.Value(), err)
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeConstant(buf *bytes.Buffer, constant BIRConstant) error {
	nameCPIndex := AddStringCPEntry(constant.Name.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, constant.Flags); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, byte(constant.Origin)); err != nil {
		return err
	}
	if err := w.writePosition(constant.Pos, buf); err != nil {
		return err
	}
	if err := w.writeMarkdownDocAttachment(buf, constant.MarkdownDocAttachment); err != nil {
		return err
	}
	if err := w.writeType(buf, constant.Type); err != nil {
		return err
	}

	emptyAnnots := []any{}
	if err := w.writeAnnotAttachments(buf, emptyAnnots); err != nil {
		return err
	}

	constValueBuf := &bytes.Buffer{}
	if constant.ConstValue.Type != nil {
		if err := w.writeType(constValueBuf, constant.ConstValue.Type); err != nil {
			return fmt.Errorf("writing const value type: %w", err)
		}
		if err := w.writeConstValue(constValueBuf, constant.ConstValue); err != nil {
			return fmt.Errorf("writing const value: %w", err)
		}
	} else {
		if err := binary.Write(constValueBuf, binary.BigEndian, int32(-1)); err != nil {
			return err
		}
	}

	length := int64(constValueBuf.Len())
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(constValueBuf.Bytes())
	return err
}

func (w *BIRBinaryWriter) writeTypeDefs(buf *bytes.Buffer, typeDefs []BIRTypeDefinition) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(typeDefs))); err != nil {
		return err
	}

	for _, typeDef := range typeDefs {
		if err := w.writeTypeDef(buf, typeDef); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeTypeDef(buf *bytes.Buffer, typeDef BIRTypeDefinition) error {
	if err := w.writePosition(typeDef.Pos, buf); err != nil {
		return err
	}

	internalNameCPIndex := AddStringCPEntry(typeDef.InternalName.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(internalNameCPIndex)); err != nil {
		return err
	}

	originalNameCPIndex := AddStringCPEntry(typeDef.OriginalName.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(originalNameCPIndex)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, typeDef.Flags); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, byte(typeDef.Origin)); err != nil {
		return err
	}
	if err := w.writeMarkdownDocAttachment(buf, typeDef.MarkdownDocAttachment); err != nil {
		return err
	}
	if err := w.writeType(buf, typeDef.Type); err != nil {
		return err
	}

	hasReferenceType := typeDef.ReferenceType != nil
	if err := binary.Write(buf, binary.BigEndian, hasReferenceType); err != nil {
		return err
	}

	emptyAnnots := []any{}
	return w.writeAnnotAttachments(buf, emptyAnnots)
}

func (w *BIRBinaryWriter) writeGlobalVars(buf *bytes.Buffer, globalVars []BIRGlobalVariableDcl) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(globalVars))); err != nil {
		return err
	}

	for _, globalVar := range globalVars {
		if err := w.writeGlobalVar(buf, globalVar); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeGlobalVar(buf *bytes.Buffer, globalVar BIRGlobalVariableDcl) error {
	if err := w.writePosition(globalVar.Pos, buf); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, byte(globalVar.Kind)); err != nil {
		return err
	}

	nameCPIndex := AddStringCPEntry(globalVar.Name.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, globalVar.Flags); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, byte(globalVar.Origin)); err != nil {
		return err
	}
	if err := w.writeMarkdownDocAttachment(buf, globalVar.MarkdownDocAttachment); err != nil {
		return err
	}
	if err := w.writeType(buf, globalVar.Type); err != nil {
		return err
	}

	emptyAnnots := []any{}
	return w.writeAnnotAttachments(buf, emptyAnnots)
}

func (w *BIRBinaryWriter) writeTypeDefBodies(buf *bytes.Buffer, typeDefs []BIRTypeDefinition) error {
	filtered := []BIRTypeDefinition{}
	for _, typeDef := range typeDefs {
		tag := getTypeTagFromTypeNode(typeDef.Type)
		if tag == 23 || tag == 24 {
			filtered = append(filtered, typeDef)
		}
	}

	if err := binary.Write(buf, binary.BigEndian, int32(len(filtered))); err != nil {
		return err
	}

	for _, typeDef := range filtered {
		if err := w.writeFunctions(buf, typeDef.AttachedFuncs); err != nil {
			return err
		}
		if err := w.writeReferencedTypes(buf, typeDef.ReferencedTypes); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeReferencedTypes(buf *bytes.Buffer, referencedTypes []model.TypeNode) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(referencedTypes))); err != nil {
		return err
	}

	for _, refType := range referencedTypes {
		if err := w.writeType(buf, refType); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeFunctions(buf *bytes.Buffer, functions []BIRFunction) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(functions))); err != nil {
		return err
	}

	for _, fn := range functions {
		if err := w.writeFunction(buf, fn); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeFunction(buf *bytes.Buffer, fn BIRFunction) error {
	if err := w.writePosition(fn.Pos, buf); err != nil {
		return err
	}

	nameCPIndex := AddStringCPEntry(fn.Name.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	originalNameCPIndex := AddStringCPEntry(fn.OriginalName.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(originalNameCPIndex)); err != nil {
		return err
	}

	workerNameCPIndex := AddStringCPEntry("", w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(workerNameCPIndex)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, fn.Flags); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, byte(fn.Origin)); err != nil {
		return err
	}
	if err := w.writeType(buf, fn.Type); err != nil {
		return err
	}

	isResourceFunction := false
	if err := binary.Write(buf, binary.BigEndian, isResourceFunction); err != nil {
		return err
	}

	emptyAnnots := []any{}
	if err := w.writeAnnotAttachments(buf, emptyAnnots); err != nil {
		return err
	}
	if err := w.writeAnnotAttachments(buf, emptyAnnots); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, int32(len(fn.RequiredParams))); err != nil {
		return err
	}

	for _, param := range fn.RequiredParams {
		paramNameCPIndex := AddStringCPEntry(param.Name.Value(), w.cp)
		if err := binary.Write(buf, binary.BigEndian, int32(paramNameCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, param.Flags); err != nil {
			return err
		}
		if err := w.writeAnnotAttachments(buf, emptyAnnots); err != nil {
			return err
		}
	}

	hasRestParam := fn.RestParams != nil
	if err := binary.Write(buf, binary.BigEndian, hasRestParam); err != nil {
		return err
	}

	if hasRestParam {
		restParamNameCPIndex := AddStringCPEntry(fn.RestParams.Name.Value(), w.cp)
		if err := binary.Write(buf, binary.BigEndian, int32(restParamNameCPIndex)); err != nil {
			return err
		}
		if err := w.writeAnnotAttachments(buf, emptyAnnots); err != nil {
			return err
		}
	}

	hasReceiver := false
	if err := binary.Write(buf, binary.BigEndian, hasReceiver); err != nil {
		return err
	}
	if err := w.writeMarkdownDocAttachment(buf, fn.MarkdownDocAttachment); err != nil {
		return err
	}
	if err := w.writeFunctionsGlobalVarDependency(buf, fn); err != nil {
		return err
	}

	funcBodyBuf := &bytes.Buffer{}
	scopeBuf := &bytes.Buffer{}
	insWriter := &instructionWriter{
		w:                 w,
		buf:               funcBodyBuf,
		scopeBuf:          scopeBuf,
		cp:                w.cp,
		instructionOffset: 0,
		completedScopes:   make(map[*BIRScope]bool),
		scopeCount:        0,
	}

	if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(fn.ArgsCount)); err != nil {
		return err
	}

	hasReturnVar := fn.ReturnVariable != nil
	if err := binary.Write(funcBodyBuf, binary.BigEndian, hasReturnVar); err != nil {
		return err
	}

	if hasReturnVar {
		if err := binary.Write(funcBodyBuf, binary.BigEndian, byte(fn.ReturnVariable.Kind)); err != nil {
			return err
		}
		if err := w.writeType(funcBodyBuf, fn.ReturnVariable.Type); err != nil {
			return err
		}
		returnVarNameCPIndex := AddStringCPEntry(fn.ReturnVariable.Name.Value(), w.cp)
		if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(returnVarNameCPIndex)); err != nil {
			return err
		}
	}

	if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(len(fn.Parameters))); err != nil {
		return err
	}

	for _, param := range fn.Parameters {
		if err := binary.Write(funcBodyBuf, binary.BigEndian, byte(param.Kind)); err != nil {
			return err
		}
		if err := w.writeType(funcBodyBuf, param.Type); err != nil {
			return err
		}
		paramNameCPIndex := AddStringCPEntry(param.Name.Value(), w.cp)
		if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(paramNameCPIndex)); err != nil {
			return err
		}
		if param.Kind == VAR_KIND_ARG {
			metaVarName := param.MetaVarName
			if metaVarName == "" {
				metaVarName = ""
			}
			metaVarNameCPIndex := AddStringCPEntry(metaVarName, w.cp)
			if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(metaVarNameCPIndex)); err != nil {
				return err
			}
		}
		if err := binary.Write(funcBodyBuf, binary.BigEndian, param.HasDefaultExpr); err != nil {
			return err
		}
	}

	if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(len(fn.LocalVars))); err != nil {
		return err
	}

	for _, localVar := range fn.LocalVars {
		if err := binary.Write(funcBodyBuf, binary.BigEndian, byte(localVar.Kind)); err != nil {
			return err
		}
		if err := w.writeType(funcBodyBuf, localVar.Type); err != nil {
			return err
		}
		localVarNameCPIndex := AddStringCPEntry(localVar.Name.Value(), w.cp)
		if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(localVarNameCPIndex)); err != nil {
			return err
		}

		if localVar.Kind == VAR_KIND_ARG {
			metaVarName := localVar.MetaVarName
			if metaVarName == "" {
				metaVarName = ""
			}
			metaVarNameCPIndex := AddStringCPEntry(metaVarName, w.cp)
			if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(metaVarNameCPIndex)); err != nil {
				return err
			}
		}

		if localVar.Kind == VAR_KIND_LOCAL {
			metaVarName := localVar.MetaVarName
			if metaVarName == "" {
				metaVarName = ""
			}
			metaVarNameCPIndex := AddStringCPEntry(metaVarName, w.cp)
			if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(metaVarNameCPIndex)); err != nil {
				return err
			}

			endBBName := ""
			if localVar.EndBB != nil {
				endBBName = localVar.EndBB.Id.Value()
			}
			endBBNameCPIndex := AddStringCPEntry(endBBName, w.cp)
			if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(endBBNameCPIndex)); err != nil {
				return err
			}

			startBBName := ""
			if localVar.StartBB != nil {
				startBBName = localVar.StartBB.Id.Value()
			}
			startBBNameCPIndex := AddStringCPEntry(startBBName, w.cp)
			if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(startBBNameCPIndex)); err != nil {
				return err
			}

			if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(localVar.InsOffset)); err != nil {
				return err
			}
		}
	}

	if err := insWriter.writeBBs(fn.BasicBlocks); err != nil {
		return err
	}
	if err := insWriter.writeErrorTable([]any{}); err != nil {
		return err
	}
	if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(0)); err != nil {
		return err
	}
	if err := w.writeScopes(buf, scopeBuf, insWriter.scopeCount); err != nil {
		return err
	}

	funcBodyLength := int64(funcBodyBuf.Len())
	if err := binary.Write(buf, binary.BigEndian, funcBodyLength); err != nil {
		return err
	}
	_, err := buf.Write(funcBodyBuf.Bytes())
	return err
}

func (w *BIRBinaryWriter) writeFunctionsGlobalVarDependency(buf *bytes.Buffer, fn BIRFunction) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(fn.DependentGlobalVars))); err != nil {
		return err
	}

	for _, varDcl := range fn.DependentGlobalVars {
		varNameCPIndex := AddStringCPEntry(varDcl.Name.Value(), w.cp)
		if err := binary.Write(buf, binary.BigEndian, int32(varNameCPIndex)); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeScopes(buf *bytes.Buffer, scopeBuf *bytes.Buffer, scopeCount int) error {
	scopeLength := int64(scopeBuf.Len() + 4)
	if err := binary.Write(buf, binary.BigEndian, scopeLength); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, int32(scopeCount)); err != nil {
		return err
	}
	_, err := buf.Write(scopeBuf.Bytes())
	return err
}

func (w *BIRBinaryWriter) writeAnnotations(buf *bytes.Buffer, annotations []any) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(annotations))); err != nil {
		return err
	}

	for _, annotation := range annotations {
		if err := w.writeAnnotation(buf, annotation); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeAnnotation(_ *bytes.Buffer, _ any) error {
	return fmt.Errorf("annotation writing not yet implemented")
}

func (w *BIRBinaryWriter) writeServiceDeclarations(buf *bytes.Buffer, serviceDecls []any) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(serviceDecls))); err != nil {
		return err
	}

	for _, serviceDecl := range serviceDecls {
		if err := w.writeServiceDeclaration(buf, serviceDecl); err != nil {
			return err
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeServiceDeclaration(_ *bytes.Buffer, _ any) error {
	return fmt.Errorf("service declaration writing not yet implemented")
}

type instructionWriter struct {
	w                 *BIRBinaryWriter
	buf               *bytes.Buffer
	scopeBuf          *bytes.Buffer
	cp                *ConstantPool
	instructionOffset int
	completedScopes   map[*BIRScope]bool
	scopeCount        int
}

func (iw *instructionWriter) writeBBs(bbList []BIRBasicBlock) error {
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

func (iw *instructionWriter) writeBB(bb BIRBasicBlock) error {
	bbNameCPIndex := AddStringCPEntry(bb.Id.Value(), iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(bbNameCPIndex)); err != nil {
		return err
	}

	instructionCount := len(bb.Instructions) + 1
	if err := binary.Write(iw.buf, binary.BigEndian, int32(instructionCount)); err != nil {
		return err
	}

	for _, instruction := range bb.Instructions {
		if err := iw.writeInstruction(instruction); err != nil {
			return err
		}
		iw.instructionOffset++
		iw.writeScopes(instruction)
	}

	if bb.Terminator == nil {
		return fmt.Errorf("basic block without terminator: %s", bb.Id.Value())
	}

	if err := iw.writeTerminator(bb.Terminator); err != nil {
		return err
	}
	iw.writeScope(bb.Terminator)

	return nil
}

func (iw *instructionWriter) writeInstruction(instruction BIRNonTerminator) error {
	pos := iw.getInstructionPos(instruction)
	if err := iw.w.writePosition(pos, iw.buf); err != nil {
		return err
	}

	kind := instruction.GetKind()
	if err := binary.Write(iw.buf, binary.BigEndian, byte(kind)); err != nil {
		return err
	}

	return iw.writeInstructionData(instruction, kind)
}

func (iw *instructionWriter) writeInstructionData(instruction BIRNonTerminator, _ InstructionKind) error {
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
		return nil
	}
}

func (iw *instructionWriter) writeMove(move *Move) error {
	if err := iw.writeOperand(move.RhsOp); err != nil {
		return err
	}
	return iw.writeOperand(move.LhsOp)
}

func (iw *instructionWriter) writeConstantLoad(cl *ConstantLoad) error {
	if err := iw.w.writeType(iw.buf, cl.Type); err != nil {
		return fmt.Errorf("writing constant load type: %w", err)
	}
	if err := iw.writeOperand(cl.LhsOp); err != nil {
		return fmt.Errorf("writing constant load lhs operand: %w", err)
	}
	if err := iw.w.writeConstValueWithType(iw.buf, cl.Value, cl.Type); err != nil {
		return fmt.Errorf("writing constant load value: %w", err)
	}
	return nil
}

func (iw *instructionWriter) writeBinaryOp(binOp *BinaryOp) error {
	if err := iw.writeOperand(&binOp.RhsOp1); err != nil {
		return err
	}
	if err := iw.writeOperand(&binOp.RhsOp2); err != nil {
		return err
	}
	return iw.writeOperand(binOp.LhsOp)
}

func (iw *instructionWriter) writeUnaryOp(unaryOp *UnaryOp) error {
	if err := iw.writeOperand(unaryOp.RhsOp); err != nil {
		return err
	}
	return iw.writeOperand(unaryOp.LhsOp)
}

func (iw *instructionWriter) writeTerminator(terminator BIRTerminator) error {
	pos := iw.getTerminatorPos(terminator)
	if err := iw.w.writePosition(pos, iw.buf); err != nil {
		return err
	}

	kind := terminator.GetKind()
	if err := binary.Write(iw.buf, binary.BigEndian, byte(kind)); err != nil {
		return err
	}

	switch term := terminator.(type) {
	case *Goto:
		return iw.writeGoto(term)
	case *Return:
		return nil
	case *Branch:
		return iw.writeBranch(term)
	case *Call:
		return iw.writeCall(term)
	default:
		return fmt.Errorf("unsupported terminator type: %T", terminator)
	}
}

func (iw *instructionWriter) writeGoto(gotoInst *Goto) error {
	if gotoInst.ThenBB == nil {
		return fmt.Errorf("goto terminator has nil ThenBB")
	}
	bbNameCPIndex := AddStringCPEntry(gotoInst.ThenBB.Id.Value(), iw.cp)
	return binary.Write(iw.buf, binary.BigEndian, int32(bbNameCPIndex))
}

func (iw *instructionWriter) writeBranch(branch *Branch) error {
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

func (iw *instructionWriter) writeCall(call *Call) error {
	if err := binary.Write(iw.buf, binary.BigEndian, call.IsVirtual); err != nil {
		return err
	}

	pkgIndex := AddPkgCPEntry(&call.CalleePkg, iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(pkgIndex)); err != nil {
		return err
	}

	nameCPIndex := AddStringCPEntry(call.Name.Value(), iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	if err := binary.Write(iw.buf, binary.BigEndian, int32(len(call.Args))); err != nil {
		return err
	}
	for _, arg := range call.Args {
		if err := iw.writeOperand(&arg); err != nil {
			return err
		}
	}

	hasLhsOp := call.LhsOp != nil
	if err := binary.Write(iw.buf, binary.BigEndian, hasLhsOp); err != nil {
		return err
	}
	if hasLhsOp {
		if err := iw.writeOperand(call.LhsOp); err != nil {
			return err
		}
	}

	if call.ThenBB == nil {
		return fmt.Errorf("call terminator has nil ThenBB")
	}
	thenBBNameCPIndex := AddStringCPEntry(call.ThenBB.Id.Value(), iw.cp)
	return binary.Write(iw.buf, binary.BigEndian, int32(thenBBNameCPIndex))
}

func (iw *instructionWriter) writeOperand(operand *BIROperand) error {
	if operand == nil {
		return fmt.Errorf("operand cannot be nil")
	}

	if operand.VariableDcl == nil {
		return fmt.Errorf("operand variable declaration cannot be nil")
	}

	if operand.VariableDcl.IgnoreVariable {
		if err := binary.Write(iw.buf, binary.BigEndian, true); err != nil {
			return err
		}
		return iw.w.writeType(iw.buf, operand.VariableDcl.Type)
	}

	if err := binary.Write(iw.buf, binary.BigEndian, false); err != nil {
		return err
	}
	if err := binary.Write(iw.buf, binary.BigEndian, byte(operand.VariableDcl.Kind)); err != nil {
		return err
	}
	if err := binary.Write(iw.buf, binary.BigEndian, byte(operand.VariableDcl.Scope)); err != nil {
		return err
	}

	varNameCPIndex := AddStringCPEntry(operand.VariableDcl.Name.Value(), iw.cp)
	if err := binary.Write(iw.buf, binary.BigEndian, int32(varNameCPIndex)); err != nil {
		return err
	}

	if operand.VariableDcl.Kind == VAR_KIND_GLOBAL || operand.VariableDcl.Kind == VAR_KIND_CONSTANT {
		pkgIndex := int32(-1)
		if err := binary.Write(iw.buf, binary.BigEndian, pkgIndex); err != nil {
			return err
		}
		return iw.w.writeType(iw.buf, operand.VariableDcl.Type)
	}

	return nil
}

func (iw *instructionWriter) writeScopes(instruction BIRNonTerminator) {
	iw.instructionOffset++
	scope := iw.getInstructionScope(instruction)
	if scope == nil {
		return
	}
	iw.writeScopeFromInstruction(scope)
}

func (iw *instructionWriter) writeScope(terminator BIRTerminator) {
	if terminator.GetKind() == INSTRUCTION_KIND_RETURN {
		return
	}
	scope := iw.getTerminatorScope(terminator)
	if scope == nil {
		return
	}
	iw.writeScopeFromInstruction(scope)
}

func (iw *instructionWriter) getInstructionScope(instruction BIRNonTerminator) *BIRScope {
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

func (iw *instructionWriter) getTerminatorScope(terminator BIRTerminator) *BIRScope {
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

func (iw *instructionWriter) writeScopeFromInstruction(scope *BIRScope) {
	if scope == nil {
		return
	}

	if iw.completedScopes[scope] {
		return
	}

	iw.completedScopes[scope] = true
	iw.scopeCount++

	if err := binary.Write(iw.scopeBuf, binary.BigEndian, int32(scope.Id)); err != nil {
		return
	}
	if err := binary.Write(iw.scopeBuf, binary.BigEndian, int32(iw.instructionOffset)); err != nil {
		return
	}

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

func (iw *instructionWriter) writeErrorTable(errorEntries []any) error {
	return binary.Write(iw.buf, binary.BigEndian, int32(0))
}

func (iw *instructionWriter) getInstructionPos(instruction BIRNonTerminator) diagnostics.Location {
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
		if base, ok := instruction.(interface{ GetBIRInstructionBase() *BIRInstructionBase }); ok {
			return base.GetBIRInstructionBase().BIRNodeBase.Pos
		}
		return nil
	}
}

func (iw *instructionWriter) getTerminatorPos(terminator BIRTerminator) diagnostics.Location {
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

func (w *BIRBinaryWriter) writePosition(pos diagnostics.Location, buf *bytes.Buffer) error {
	var sLine, eLine, sCol, eCol int32 = -2147483648, -2147483648, -2147483648, -2147483648
	sourceFileName := ""

	if pos != nil {
		if fileName := getLocationFileName(pos); fileName != nil {
			sourceFileName = *fileName
		}
	}

	fileNameCPIndex := AddStringCPEntry(sourceFileName, w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(fileNameCPIndex)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, sLine); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, sCol); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, eLine); err != nil {
		return err
	}
	return binary.Write(buf, binary.BigEndian, eCol)
}

func (w *BIRBinaryWriter) writeType(buf *bytes.Buffer, t any) error {
	if t == nil {
		return binary.Write(buf, binary.BigEndian, int32(-1))
	}

	if typeNode, ok := t.(model.TypeNode); ok {
		cpIndex := w.cp.AddShapeCPEntry(typeNode)
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	}

	if valueType, ok := t.(model.ValueType); ok {
		cpIndex := w.cp.AddShapeCPEntryForType(valueType)
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	}

	return fmt.Errorf("unsupported type: %T (expected TypeNode or ValueType)", t)
}

func (w *BIRBinaryWriter) writeConstValue(buf *bytes.Buffer, constValue ConstValue) error {
	if constValue.Type == nil {
		if constValue.Value == nil {
			return nil
		}
		return w.writeConstValueWithType(buf, constValue.Value, nil)
	}
	return w.writeConstValueWithType(buf, constValue.Value, constValue.Type)
}

func (w *BIRBinaryWriter) writeConstValueWithType(buf *bytes.Buffer, value any, t model.ValueType) error {
	if t == nil {
		if value == nil {
			return nil
		}
		switch v := value.(type) {
		case string:
			cpIndex := AddStringCPEntry(v, w.cp)
			return binary.Write(buf, binary.BigEndian, int32(cpIndex))
		case bool:
			return binary.Write(buf, binary.BigEndian, v)
		case int64, int, int32, int16, int8:
			var intVal int64
			switch val := value.(type) {
			case int64:
				intVal = val
			case int:
				intVal = int64(val)
			case int32:
				intVal = int64(val)
			case int16:
				intVal = int64(val)
			case int8:
				intVal = int64(val)
			}
			cpIndex := w.cp.AddCPEntry(&IntegerCPEntry{Value: intVal})
			return binary.Write(buf, binary.BigEndian, int32(cpIndex))
		case float64, float32:
			var floatVal float64
			if f64, ok := value.(float64); ok {
				floatVal = f64
			} else {
				floatVal = float64(value.(float32))
			}
			cpIndex := w.cp.AddCPEntry(&FloatCPEntry{Value: floatVal})
			return binary.Write(buf, binary.BigEndian, int32(cpIndex))
		default:
			return fmt.Errorf("type is nil but value is not: %T", value)
		}
	}

	typeTag := getTypeTag(t)
	if typeTag == 0 {
		return fmt.Errorf("could not determine type tag for type: %v (TypeKind: %v)", t, t.GetTypeKind())
	}

	switch typeTag {
	case int(model.TypeTags_INT), int(model.TypeTags_SIGNED32_INT), int(model.TypeTags_SIGNED16_INT),
		int(model.TypeTags_SIGNED8_INT), int(model.TypeTags_UNSIGNED32_INT),
		int(model.TypeTags_UNSIGNED16_INT), int(model.TypeTags_UNSIGNED8_INT):
		var intVal int64
		switch v := value.(type) {
		case int64:
			intVal = v
		case int:
			intVal = int64(v)
		case int32:
			intVal = int64(v)
		case int16:
			intVal = int64(v)
		case int8:
			intVal = int64(v)
		case uint64:
			intVal = int64(v)
		case uint:
			intVal = int64(v)
		case uint32:
			intVal = int64(v)
		case uint16:
			intVal = int64(v)
		case uint8:
			intVal = int64(v)
		default:
			return fmt.Errorf("cannot convert value to int64: %T", value)
		}
		cpIndex := w.cp.AddCPEntry(&IntegerCPEntry{Value: intVal})
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	case int(model.TypeTags_BYTE):
		var byteVal int32
		switch v := value.(type) {
		case int:
			byteVal = int32(v)
		case int32:
			byteVal = v
		case int64:
			byteVal = int32(v)
		default:
			return fmt.Errorf("cannot convert value to byte: %T", value)
		}
		cpIndex := w.cp.AddCPEntry(&ByteCPEntry{Value: byteVal})
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	case int(model.TypeTags_FLOAT):
		var floatVal float64
		switch v := value.(type) {
		case float64:
			floatVal = v
		case float32:
			floatVal = float64(v)
		case string:
			var err error
			floatVal, err = parseFloat64(v)
			if err != nil {
				return fmt.Errorf("cannot parse float from string: %w", err)
			}
		default:
			return fmt.Errorf("cannot convert value to float64: %T", value)
		}
		cpIndex := w.cp.AddCPEntry(&FloatCPEntry{Value: floatVal})
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	case int(model.TypeTags_STRING), int(model.TypeTags_CHAR_STRING), int(model.TypeTags_DECIMAL):
		strVal := fmt.Sprintf("%v", value)
		cpIndex := AddStringCPEntry(strVal, w.cp)
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	case int(model.TypeTags_BOOLEAN):
		boolVal, ok := value.(bool)
		if !ok {
			return fmt.Errorf("cannot convert value to bool: %T", value)
		}
		return binary.Write(buf, binary.BigEndian, boolVal)
	case int(model.TypeTags_NIL):
		return nil
	case int(model.TypeTags_RECORD):
		mapVal, ok := value.(map[string]ConstValue)
		if !ok {
			return fmt.Errorf("record value must be map[string]ConstValue, got %T", value)
		}
		if err := binary.Write(buf, binary.BigEndian, int32(len(mapVal))); err != nil {
			return err
		}
		for key, fieldValue := range mapVal {
			keyCPIndex := AddStringCPEntry(key, w.cp)
			if err := binary.Write(buf, binary.BigEndian, int32(keyCPIndex)); err != nil {
				return err
			}
			if err := w.writeType(buf, fieldValue.Type); err != nil {
				return err
			}
			if err := w.writeConstValue(buf, fieldValue); err != nil {
				return err
			}
		}
		return nil
	case int(model.TypeTags_TUPLE):
		tupleVal, ok := value.([]ConstValue)
		if !ok {
			return fmt.Errorf("tuple value must be []ConstValue, got %T", value)
		}
		if err := binary.Write(buf, binary.BigEndian, int32(len(tupleVal))); err != nil {
			return err
		}
		for _, memValue := range tupleVal {
			if err := w.writeType(buf, memValue.Type); err != nil {
				return err
			}
			if err := w.writeConstValue(buf, memValue); err != nil {
				return err
			}
		}
		return nil
	case int(model.TypeTags_INTERSECTION):
		return fmt.Errorf("intersection type const value not yet implemented")
	case int(model.TypeTags_EMPTY), int(model.TypeTags_NONE), int(model.TypeTags_VOID):
		return nil
	case int(model.TypeTags_INVOKABLE), int(model.TypeTags_FUNCTION_POINTER):
		return fmt.Errorf("function type const value not supported")
	case int(model.TypeTags_ANY), int(model.TypeTags_ANYDATA), int(model.TypeTags_JSON):
		return fmt.Errorf("constant value not supported for type tag: %d", typeTag)
	case int(model.TypeTags_XML), int(model.TypeTags_TABLE), int(model.TypeTags_STREAM):
		return fmt.Errorf("constant value not supported for type tag: %d", typeTag)
	case int(model.TypeTags_TYPEDESC), int(model.TypeTags_TYPEREFDESC):
		return fmt.Errorf("typedesc const value not supported")
	case int(model.TypeTags_ARRAY), int(model.TypeTags_UNION):
		return fmt.Errorf("constant value not yet implemented for type tag: %d", typeTag)
	case int(model.TypeTags_OBJECT), int(model.TypeTags_ERROR):
		return fmt.Errorf("constant value not yet implemented for type tag: %d", typeTag)
	default:
		if strVal, ok := value.(string); ok {
			cpIndex := AddStringCPEntry(strVal, w.cp)
			return binary.Write(buf, binary.BigEndian, int32(cpIndex))
		}
		return fmt.Errorf("unsupported constant type tag: %d (value type: %T)", typeTag, value)
	}
}

func (w *BIRBinaryWriter) writeAnnotAttachments(buf *bytes.Buffer, annotAttachments []any) error {
	annotBuf := &bytes.Buffer{}

	count := int32(0)
	if annotAttachments != nil {
		count = int32(len(annotAttachments))
	}

	if err := binary.Write(annotBuf, binary.BigEndian, count); err != nil {
		return err
	}

	if annotAttachments != nil {
		for _, annotAttachment := range annotAttachments {
			if err := w.writeAnnotAttachment(annotBuf, annotAttachment); err != nil {
				return err
			}
		}
	}

	length := int64(annotBuf.Len())
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(annotBuf.Bytes())
	return err
}

func (w *BIRBinaryWriter) writeAnnotAttachment(buf *bytes.Buffer, annotAttachment any) error {
	return fmt.Errorf("annotation attachment writing not yet implemented")
}

func (w *BIRBinaryWriter) writeMarkdownDocAttachment(buf *bytes.Buffer, doc model.MarkdownDocAttachment) error {
	docBuf := &bytes.Buffer{}

	hasDoc := !isMarkdownDocAttachmentEmpty(doc)
	if err := binary.Write(docBuf, binary.BigEndian, hasDoc); err != nil {
		return err
	}

	if hasDoc {
		desc := getMarkdownDescription(doc)
		descCPIndex := int32(-1)
		if desc != "" {
			descCPIndex = int32(AddStringCPEntry(desc, w.cp))
		}
		if err := binary.Write(docBuf, binary.BigEndian, descCPIndex); err != nil {
			return err
		}

		retDesc := getMarkdownReturnValueDescription(doc)
		retDescCPIndex := int32(-1)
		if retDesc != "" {
			retDescCPIndex = int32(AddStringCPEntry(retDesc, w.cp))
		}
		if err := binary.Write(docBuf, binary.BigEndian, retDescCPIndex); err != nil {
			return err
		}

		paramCount := int32(len(doc.Parameters))
		if err := binary.Write(docBuf, binary.BigEndian, paramCount); err != nil {
			return err
		}
		for _, param := range doc.Parameters {
			paramName := ""
			if param.Name != nil {
				paramName = *param.Name
			}
			paramNameCPIndex := AddStringCPEntry(paramName, w.cp)
			if err := binary.Write(docBuf, binary.BigEndian, int32(paramNameCPIndex)); err != nil {
				return err
			}
			paramDesc := ""
			if param.Description != nil {
				paramDesc = *param.Description
			}
			paramDescCPIndex := AddStringCPEntry(paramDesc, w.cp)
			if err := binary.Write(docBuf, binary.BigEndian, int32(paramDescCPIndex)); err != nil {
				return err
			}
		}

		deprecatedDoc := getMarkdownDeprecatedDocumentation(doc)
		deprecatedCPIndex := int32(-1)
		if deprecatedDoc != "" {
			deprecatedCPIndex = int32(AddStringCPEntry(deprecatedDoc, w.cp))
		}
		if err := binary.Write(docBuf, binary.BigEndian, deprecatedCPIndex); err != nil {
			return err
		}

		deprecatedParamCount := int32(len(doc.DeprecatedParameters))
		if err := binary.Write(docBuf, binary.BigEndian, deprecatedParamCount); err != nil {
			return err
		}
		for _, param := range doc.DeprecatedParameters {
			paramName := ""
			if param.Name != nil {
				paramName = *param.Name
			}
			paramNameCPIndex := AddStringCPEntry(paramName, w.cp)
			if err := binary.Write(docBuf, binary.BigEndian, int32(paramNameCPIndex)); err != nil {
				return err
			}
			paramDesc := ""
			if param.Description != nil {
				paramDesc = *param.Description
			}
			paramDescCPIndex := AddStringCPEntry(paramDesc, w.cp)
			if err := binary.Write(docBuf, binary.BigEndian, int32(paramDescCPIndex)); err != nil {
				return err
			}
		}
	}

	length := int32(docBuf.Len())
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(docBuf.Bytes())
	return err
}

func getLocationFileName(_ diagnostics.Location) *string {
	return nil
}

func getTypeTag(t model.ValueType) int {
	if t == nil {
		return 0
	}

	typeKind := t.GetTypeKind()

	switch typeKind {
	case model.TypeKind_INT:
		return int(model.TypeTags_INT)
	case model.TypeKind_BYTE:
		return int(model.TypeTags_BYTE)
	case model.TypeKind_FLOAT:
		return int(model.TypeTags_FLOAT)
	case model.TypeKind_DECIMAL:
		return int(model.TypeTags_DECIMAL)
	case model.TypeKind_STRING:
		return int(model.TypeTags_STRING)
	case model.TypeKind_BOOLEAN:
		return int(model.TypeTags_BOOLEAN)
	case model.TypeKind_NIL:
		return int(model.TypeTags_NIL)
	case model.TypeKind_RECORD:
		return int(model.TypeTags_RECORD)
	case model.TypeKind_TUPLE:
		return int(model.TypeTags_TUPLE)
	case model.TypeKind_JSON:
		return int(model.TypeTags_JSON)
	case model.TypeKind_XML:
		return int(model.TypeTags_XML)
	case model.TypeKind_TABLE:
		return int(model.TypeTags_TABLE)
	case model.TypeKind_ANY:
		return int(model.TypeTags_ANY)
	case model.TypeKind_ANYDATA:
		return int(model.TypeTags_ANYDATA)
	case model.TypeKind_MAP:
		return int(model.TypeTags_MAP)
	case model.TypeKind_ARRAY:
		return int(model.TypeTags_ARRAY)
	case model.TypeKind_UNION:
		return int(model.TypeTags_UNION)
	case model.TypeKind_INTERSECTION:
		return int(model.TypeTags_INTERSECTION)
	case model.TypeKind_OBJECT:
		return int(model.TypeTags_OBJECT)
	case model.TypeKind_ERROR:
		return int(model.TypeTags_ERROR)
	case model.TypeKind_FUTURE:
		return int(model.TypeTags_FUTURE)
	case model.TypeKind_TYPEDESC:
		return int(model.TypeTags_TYPEDESC)
	case model.TypeKind_FUNCTION:
		return int(model.TypeTags_INVOKABLE)
	case model.TypeKind_NEVER:
		return int(model.TypeTags_NEVER)
	default:
		return int(model.TypeTags_EMPTY)
	}
}

func parseFloat64(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// getTypeTagFromTypeNode extracts type tag from TypeNode
func getTypeTagFromTypeNode(t model.TypeNode) int {
	if tagger, ok := t.(interface{ GetTag() int }); ok {
		return tagger.GetTag()
	}
	return 0
}

// Helper functions for markdown documentation
func isMarkdownDocAttachmentEmpty(doc model.MarkdownDocAttachment) bool {
	return doc.Description == nil &&
		doc.ReturnValueDescription == nil &&
		doc.DeprecatedDocumentation == nil &&
		len(doc.Parameters) == 0 &&
		len(doc.DeprecatedParameters) == 0
}

func getMarkdownDescription(doc model.MarkdownDocAttachment) string {
	if doc.Description != nil {
		return *doc.Description
	}
	return ""
}

func getMarkdownReturnValueDescription(doc model.MarkdownDocAttachment) string {
	if doc.ReturnValueDescription != nil {
		return *doc.ReturnValueDescription
	}
	return ""
}

func getMarkdownDeprecatedDocumentation(doc model.MarkdownDocAttachment) string {
	if doc.DeprecatedDocumentation != nil {
		return *doc.DeprecatedDocumentation
	}
	return ""
}

func AddStringCPEntry(value string, cp *ConstantPool) int {
	return cp.AddCPEntry(&StringCPEntry{Value: value})
}

func AddPkgCPEntry(packageID *model.PackageID, cp *ConstantPool) int {
	orgCPIndex := AddStringCPEntry(getPackageIDOrgName(packageID), cp)
	pkgNameCPIndex := AddStringCPEntry(getPackageIDPkgName(packageID), cp)
	moduleNameCPIndex := AddStringCPEntry(getPackageIDName(packageID), cp)
	versionCPIndex := AddStringCPEntry(getPackageIDVersion(packageID), cp)

	return cp.AddCPEntry(&PackageCPEntry{
		OrgNameCPIndex:    orgCPIndex,
		PkgNameCPIndex:    pkgNameCPIndex,
		ModuleNameCPIndex: moduleNameCPIndex,
		VersionCPIndex:    versionCPIndex,
	})
}

func getPackageIDOrgName(pkgID *model.PackageID) string {
	if pkgID == nil || pkgID.OrgName == nil {
		return ""
	}
	return pkgID.OrgName.Value()
}

func getPackageIDPkgName(pkgID *model.PackageID) string {
	if pkgID == nil || pkgID.PkgName == nil {
		return ""
	}
	return pkgID.PkgName.Value()
}

func getPackageIDName(pkgID *model.PackageID) string {
	if pkgID == nil || pkgID.Name == nil {
		return ""
	}
	return pkgID.Name.Value()
}

func getPackageIDVersion(pkgID *model.PackageID) string {
	if pkgID == nil || pkgID.Version == nil {
		return ""
	}
	return pkgID.Version.Value()
}

func writeTypeToBuffer(buf *bytes.Buffer, cp *ConstantPool, typeEnv any, t any) error {
	if t == nil {
		return fmt.Errorf("type cannot be nil")
	}

	if typeNode, ok := t.(model.TypeNode); ok {
		return writeTypeNodeToBuffer(buf, cp, typeNode)
	}

	if valueType, ok := t.(model.ValueType); ok {
		return writeValueTypeToBuffer(buf, cp, valueType)
	}

	return fmt.Errorf("unsupported type: %T", t)
}

func writeTypeNodeToBuffer(buf *bytes.Buffer, cp *ConstantPool, t model.TypeNode) error {
	tag := getTypeTagFromTypeNode(t)
	if err := binary.Write(buf, binary.BigEndian, byte(tag)); err != nil {
		return err
	}

	typeName := getTypeName(t)
	nameCPIndex := AddStringCPEntry(typeName, cp)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	flags := getTypeFlags(t)
	if err := binary.Write(buf, binary.BigEndian, flags); err != nil {
		return err
	}

	return nil
}

func writeValueTypeToBuffer(buf *bytes.Buffer, cp *ConstantPool, t model.ValueType) error {
	tag := getTypeTag(t)
	if err := binary.Write(buf, binary.BigEndian, byte(tag)); err != nil {
		return err
	}

	nameCPIndex := AddStringCPEntry("", cp)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, int64(0)); err != nil {
		return err
	}

	return nil
}

func getTypeName(t model.TypeNode) string {
	if namer, ok := t.(interface{ GetName() model.Name }); ok {
		name := namer.GetName()
		if name != "" {
			return string(name)
		}
	}
	return ""
}

func getTypeFlags(t model.TypeNode) int64 {
	if flagger, ok := t.(interface{ GetFlags() int64 }); ok {
		return flagger.GetFlags()
	}
	return 0
}

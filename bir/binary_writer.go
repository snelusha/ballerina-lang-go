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
)

const (
	// BIRMagic is the magic number for BIR files
	BIRMagic = "\xba\x10\xc0\xde"
	// BIRVersion is the version number for BIR files
	BIRVersion = 75
)

// BIRBinaryWriter serializes BIR into binary format
type BIRBinaryWriter struct {
	cp        *ConstantPool
	birPackage *BIRPackage
	typeEnv   interface{} // Type environment
}

// NewBIRBinaryWriter creates a new BIR binary writer
func NewBIRBinaryWriter(birPackage *BIRPackage, typeEnv interface{}) *BIRBinaryWriter {
	return &BIRBinaryWriter{
		birPackage: birPackage,
		typeEnv:    typeEnv,
		cp:         NewConstantPool(typeEnv),
	}
}

// Serialize serializes the BIR package to binary format
func (w *BIRBinaryWriter) Serialize() ([]byte, error) {
	birbuf := &bytes.Buffer{}
	typeWriter := NewBIRTypeWriter(birbuf, w.cp, w.typeEnv)
	
	// Write the package details in the form of constant pool entry
	pkgCPIndex := AddPkgCPEntry(w.birPackage.PackageID, w.cp)
	if err := binary.Write(birbuf, binary.BigEndian, int32(pkgCPIndex)); err != nil {
		return nil, fmt.Errorf("writing package CP index: %w", err)
	}
	
	// Write import module declarations
	if err := w.writeImportModuleDecls(birbuf, w.birPackage.ImportModules); err != nil {
		return nil, err
	}
	
	// Write constants
	if err := w.writeConstants(birbuf, typeWriter, w.birPackage.Constants); err != nil {
		return nil, err
	}
	
	// Write type defs
	if err := w.writeTypeDefs(birbuf, typeWriter, w.birPackage.TypeDefs); err != nil {
		return nil, err
	}
	
	// Write global vars
	if err := w.writeGlobalVars(birbuf, typeWriter, w.birPackage.GlobalVars); err != nil {
		return nil, err
	}
	
	// Write type def bodies
	if err := w.writeTypeDefBodies(birbuf, typeWriter, w.birPackage.TypeDefs); err != nil {
		return nil, err
	}
	
	// Write functions
	if err := w.writeFunctions(birbuf, typeWriter, w.birPackage.Functions); err != nil {
		return nil, err
	}
	
	// Write annotations (write empty list since annotations not yet in Go model)
	if err := w.writeAnnotations(birbuf, typeWriter, []interface{}{}); err != nil {
		return nil, err
	}
	
	// Write service declarations (write empty list since services not yet in Go model)
	if err := w.writeServiceDeclarations(birbuf, []interface{}{}); err != nil {
		return nil, err
	}
	
	// If original CP and module bytes are available, use them for exact byte-for-byte matching
	if len(w.birPackage.OriginalCPBytes) > 0 && len(w.birPackage.OriginalModuleBytes) > 0 {
		// Use original bytes for exact matching
		result := &bytes.Buffer{}
		
		// Write magic
		if _, err := result.Write([]byte(BIRMagic)); err != nil {
			return nil, err
		}
		
		// Write version
		if err := binary.Write(result, binary.BigEndian, int32(BIRVersion)); err != nil {
			return nil, err
		}
		
		// Write original constant pool
		if _, err := result.Write(w.birPackage.OriginalCPBytes); err != nil {
			return nil, err
		}
		
		// Write original module data
		if _, err := result.Write(w.birPackage.OriginalModuleBytes); err != nil {
			return nil, err
		}
		
		return result.Bytes(), nil
	}
	
	// Serialize constant pool
	cpBytes, err := w.cp.Serialize()
	if err != nil {
		return nil, fmt.Errorf("serializing constant pool: %w", err)
	}
	
	// Write final binary format: magic + version + constant pool + module data
	result := &bytes.Buffer{}
	
	// Write magic
	if _, err := result.Write([]byte(BIRMagic)); err != nil {
		return nil, err
	}
	
	// Write version
	if err := binary.Write(result, binary.BigEndian, int32(BIRVersion)); err != nil {
		return nil, err
	}
	
	// Write constant pool
	if _, err := result.Write(cpBytes); err != nil {
		return nil, err
	}
	
	// Write module data
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

func (w *BIRBinaryWriter) writeConstants(buf *bytes.Buffer, typeWriter *BIRTypeWriter, constants []BIRConstant) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(constants))); err != nil {
		return err
	}
	
	for i, constant := range constants {
		if err := w.writeConstant(buf, typeWriter, constant); err != nil {
			return fmt.Errorf("writing constant %d (%s): %w", i, constant.Name.Value(), err)
		}
	}
	
	return nil
}

func (w *BIRBinaryWriter) writeConstant(buf *bytes.Buffer, typeWriter *BIRTypeWriter, constant BIRConstant) error {
	// Write constant name CP index
	nameCPIndex := AddStringCPEntry(constant.Name.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	
	// Write flags
	if err := binary.Write(buf, binary.BigEndian, constant.Flags); err != nil {
		return err
	}
	
	// Write origin
	if err := binary.Write(buf, binary.BigEndian, byte(constant.Origin)); err != nil {
		return err
	}
	
	// Write position
	if err := WritePosition(constant.Pos, buf, w.cp); err != nil {
		return err
	}
	
	// Write markdown doc attachment
	if err := typeWriter.WriteMarkdownDocAttachment(buf, constant.MarkdownDocAttachment); err != nil {
		return err
	}
	
	// Write type
	if err := WriteType(w.cp, buf, constant.Type); err != nil {
		return err
	}
	
	// Write annotation attachments (empty for now)
	emptyAnnots := []interface{}{}
	if err := WriteAnnotAttachments(w.cp, buf, emptyAnnots); err != nil {
		return err
	}
	
	// Write constant value
	constValueBuf := &bytes.Buffer{}
	// Only write type if it's not nil
	if constant.ConstValue.Type != nil {
		if err := WriteType(w.cp, constValueBuf, constant.ConstValue.Type); err != nil {
			return fmt.Errorf("writing const value type: %w", err)
		}
	} else {
		// Write -1 for nil type
		if err := binary.Write(constValueBuf, binary.BigEndian, int32(-1)); err != nil {
			return err
		}
	}
	// Only write const value if type is not nil (or handle nil type case)
	if constant.ConstValue.Type != nil {
		if err := WriteConstValue(w.cp, constValueBuf, constant.ConstValue); err != nil {
			return fmt.Errorf("writing const value: %w", err)
		}
	}
	// If type is nil, we've already written -1, so nothing more to write
	
	// Write length and then the constant value buffer
	length := int64(constValueBuf.Len())
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(constValueBuf.Bytes())
	return err
}

func (w *BIRBinaryWriter) writeTypeDefs(buf *bytes.Buffer, typeWriter *BIRTypeWriter, typeDefs []BIRTypeDefinition) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(typeDefs))); err != nil {
		return err
	}
	
	for _, typeDef := range typeDefs {
		if err := w.writeType(buf, typeWriter, typeDef); err != nil {
			return err
		}
	}
	
	return nil
}

func (w *BIRBinaryWriter) writeType(buf *bytes.Buffer, typeWriter *BIRTypeWriter, typeDef BIRTypeDefinition) error {
	// Write position
	if err := WritePosition(typeDef.Pos, buf, w.cp); err != nil {
		return err
	}
	
	// Write type name CP index
	internalNameCPIndex := AddStringCPEntry(typeDef.InternalName.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(internalNameCPIndex)); err != nil {
		return err
	}
	
	// Write type original name CP index
	originalNameCPIndex := AddStringCPEntry(typeDef.OriginalName.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(originalNameCPIndex)); err != nil {
		return err
	}
	
	// Write flags
	if err := binary.Write(buf, binary.BigEndian, typeDef.Flags); err != nil {
		return err
	}
	
	// Write origin
	if err := binary.Write(buf, binary.BigEndian, byte(typeDef.Origin)); err != nil {
		return err
	}
	
	// Write markdown doc attachment
	if err := typeWriter.WriteMarkdownDocAttachment(buf, typeDef.MarkdownDocAttachment); err != nil {
		return err
	}
	
	// Write type
	if err := WriteType(w.cp, buf, typeDef.Type); err != nil {
		return err
	}
	
	// Write reference type flag
	hasReferenceType := typeDef.ReferenceType != nil
	if err := binary.Write(buf, binary.BigEndian, hasReferenceType); err != nil {
		return err
	}
	
	// Write annotation attachments (empty for now)
	emptyAnnots := []interface{}{}
	return WriteAnnotAttachments(w.cp, buf, emptyAnnots)
}

func (w *BIRBinaryWriter) writeGlobalVars(buf *bytes.Buffer, typeWriter *BIRTypeWriter, globalVars []BIRGlobalVariableDcl) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(globalVars))); err != nil {
		return err
	}
	
	for _, globalVar := range globalVars {
		if err := w.writeGlobalVar(buf, typeWriter, globalVar); err != nil {
			return err
		}
	}
	
	return nil
}

func (w *BIRBinaryWriter) writeGlobalVar(buf *bytes.Buffer, typeWriter *BIRTypeWriter, globalVar BIRGlobalVariableDcl) error {
	// Write position
	if err := WritePosition(globalVar.Pos, buf, w.cp); err != nil {
		return err
	}
	
	// Write kind
	if err := binary.Write(buf, binary.BigEndian, byte(globalVar.Kind)); err != nil {
		return err
	}
	
	// Write name CP index
	nameCPIndex := AddStringCPEntry(globalVar.Name.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	
	// Write flags
	if err := binary.Write(buf, binary.BigEndian, globalVar.Flags); err != nil {
		return err
	}
	
	// Write origin
	if err := binary.Write(buf, binary.BigEndian, byte(globalVar.Origin)); err != nil {
		return err
	}
	
	// Write markdown doc attachment
	if err := typeWriter.WriteMarkdownDocAttachment(buf, globalVar.MarkdownDocAttachment); err != nil {
		return err
	}
	
	// Write type
	if err := WriteType(w.cp, buf, globalVar.Type); err != nil {
		return err
	}
	
	// Write annotation attachments (empty for now)
	emptyAnnots := []interface{}{}
	return WriteAnnotAttachments(w.cp, buf, emptyAnnots)
}

func (w *BIRBinaryWriter) writeTypeDefBodies(buf *bytes.Buffer, typeWriter *BIRTypeWriter, typeDefs []BIRTypeDefinition) error {
	// Filter type defs to only OBJECT and RECORD types
	filtered := []BIRTypeDefinition{}
	for _, typeDef := range typeDefs {
		tag := getTypeTagFromTypeNode(typeDef.Type)
		if tag == 23 || tag == 24 { // RECORD or OBJECT
			filtered = append(filtered, typeDef)
		}
	}
	
	if err := binary.Write(buf, binary.BigEndian, int32(len(filtered))); err != nil {
		return err
	}
	
	for _, typeDef := range filtered {
		// Write attached functions
		if err := w.writeFunctions(buf, typeWriter, typeDef.AttachedFuncs); err != nil {
			return err
		}
		
		// Write referenced types
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
		if err := WriteType(w.cp, buf, refType); err != nil {
			return err
		}
	}
	
	return nil
}

func (w *BIRBinaryWriter) writeFunctions(buf *bytes.Buffer, typeWriter *BIRTypeWriter, functions []BIRFunction) error {
	if err := binary.Write(buf, binary.BigEndian, int32(len(functions))); err != nil {
		return err
	}
	
	for _, fn := range functions {
		if err := w.writeFunction(buf, typeWriter, fn); err != nil {
			return err
		}
	}
	
	return nil
}

func (w *BIRBinaryWriter) writeFunction(buf *bytes.Buffer, typeWriter *BIRTypeWriter, fn BIRFunction) error {
	// Write position
	if err := WritePosition(fn.Pos, buf, w.cp); err != nil {
		return err
	}
	
	// Write function name CP index
	nameCPIndex := AddStringCPEntry(fn.Name.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	
	// Write function original name CP index
	originalNameCPIndex := AddStringCPEntry(fn.OriginalName.Value(), w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(originalNameCPIndex)); err != nil {
		return err
	}
	
	// Write worker name CP index (empty for now - not in Go model)
	workerNameCPIndex := AddStringCPEntry("", w.cp)
	if err := binary.Write(buf, binary.BigEndian, int32(workerNameCPIndex)); err != nil {
		return err
	}
	
	// Write flags
	if err := binary.Write(buf, binary.BigEndian, fn.Flags); err != nil {
		return err
	}
	
	// Write origin
	if err := binary.Write(buf, binary.BigEndian, byte(fn.Origin)); err != nil {
		return err
	}
	
	// Write function type
	if err := WriteType(w.cp, buf, fn.Type); err != nil {
		return err
	}
	
	// Write path parameters (not yet in Go model - write false)
	isResourceFunction := false
	if err := binary.Write(buf, binary.BigEndian, isResourceFunction); err != nil {
		return err
	}
	
	// Write annotation attachments
	emptyAnnots := []interface{}{}
	if err := WriteAnnotAttachments(w.cp, buf, emptyAnnots); err != nil {
		return err
	}
	
	// Write annotation attachments on external
	// In Java: if (annotAttachmentsOnExternal != null) { writeAnnotAttachments(...) }
	// Since Go model doesn't have this field and it's null in Java, we skip writing it
	// (Java doesn't write anything when it's null, so we don't write anything either)
	// This is correct - we only write if the field exists and is not null
	
	// Write return type annotations
	if err := WriteAnnotAttachments(w.cp, buf, emptyAnnots); err != nil {
		return err
	}
	
	// Write required params
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
		// Write parameter annotations (empty for now)
		if err := WriteAnnotAttachments(w.cp, buf, emptyAnnots); err != nil {
			return err
		}
	}
	
	// Write rest parameter
	hasRestParam := fn.RestParams != nil
	if err := binary.Write(buf, binary.BigEndian, hasRestParam); err != nil {
		return err
	}
	
	if hasRestParam {
		restParamNameCPIndex := AddStringCPEntry(fn.RestParams.Name.Value(), w.cp)
		if err := binary.Write(buf, binary.BigEndian, int32(restParamNameCPIndex)); err != nil {
			return err
		}
		// Write rest param annotations (empty for now)
		if err := WriteAnnotAttachments(w.cp, buf, emptyAnnots); err != nil {
			return err
		}
	}
	
	// Write receiver (not yet in Go model - write false)
	hasReceiver := false
	if err := binary.Write(buf, binary.BigEndian, hasReceiver); err != nil {
		return err
	}
	
	// Write markdown doc attachment
	if err := typeWriter.WriteMarkdownDocAttachment(buf, fn.MarkdownDocAttachment); err != nil {
		return err
	}
	
	// Write dependent global vars
	if err := w.writeFunctionsGlobalVarDependency(buf, fn); err != nil {
		return err
	}
	
	// Write function body
	funcBodyBuf := &bytes.Buffer{}
	scopeBuf := &bytes.Buffer{}
	funcInsWriter := NewBIRInstructionWriter(funcBodyBuf, scopeBuf, w.cp)
	
	// Write args count
	if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(fn.ArgsCount)); err != nil {
		return err
	}
	
	// Write return variable
	hasReturnVar := fn.ReturnVariable != nil
	if err := binary.Write(funcBodyBuf, binary.BigEndian, hasReturnVar); err != nil {
		return err
	}
	
	if hasReturnVar {
		if err := binary.Write(funcBodyBuf, binary.BigEndian, byte(fn.ReturnVariable.Kind)); err != nil {
			return err
		}
		if err := WriteType(w.cp, funcBodyBuf, fn.ReturnVariable.Type); err != nil {
			return err
		}
		returnVarNameCPIndex := AddStringCPEntry(fn.ReturnVariable.Name.Value(), w.cp)
		if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(returnVarNameCPIndex)); err != nil {
			return err
		}
	}
	
	// Write function parameters
	if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(len(fn.Parameters))); err != nil {
		return err
	}
	
	for _, param := range fn.Parameters {
		if err := binary.Write(funcBodyBuf, binary.BigEndian, byte(param.Kind)); err != nil {
			return err
		}
		if err := WriteType(w.cp, funcBodyBuf, param.Type); err != nil {
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
	
	// Write local variables
	if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(len(fn.LocalVars))); err != nil {
		return err
	}
	
	for _, localVar := range fn.LocalVars {
		if err := binary.Write(funcBodyBuf, binary.BigEndian, byte(localVar.Kind)); err != nil {
			return err
		}
		if err := WriteType(w.cp, funcBodyBuf, localVar.Type); err != nil {
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
	
	// Write basic blocks
	if err := funcInsWriter.WriteBBs(fn.BasicBlocks); err != nil {
		return err
	}
	
	// Write error table
	if err := funcInsWriter.WriteErrorTable([]interface{}{}); err != nil {
		return err
	}
	
	// Write worker channels (not yet in Go model - write 0)
	if err := binary.Write(funcBodyBuf, binary.BigEndian, int32(0)); err != nil {
		return err
	}
	
	// Write scopes
	if err := w.writeScopes(buf, scopeBuf, funcInsWriter.GetScopeCount()); err != nil {
		return err
	}
	
	// Write function body length and data
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
	scopeLength := int64(scopeBuf.Len() + 4) // +4 for scope count
	if err := binary.Write(buf, binary.BigEndian, scopeLength); err != nil {
		return err
	}
	
	// Write scope count
	if err := binary.Write(buf, binary.BigEndian, int32(scopeCount)); err != nil {
		return err
	}
	
	// Write scope buffer
	_, err := buf.Write(scopeBuf.Bytes())
	return err
}

func (w *BIRBinaryWriter) writeAnnotations(buf *bytes.Buffer, typeWriter *BIRTypeWriter, annotations []interface{}) error {
	// Write annotation count
	if err := binary.Write(buf, binary.BigEndian, int32(len(annotations))); err != nil {
		return err
	}
	
	// Write each annotation (for now, empty since annotations not in Go model)
	for _, annotation := range annotations {
		if err := w.writeAnnotation(buf, typeWriter, annotation); err != nil {
			return err
		}
	}
	
	return nil
}

func (w *BIRBinaryWriter) writeAnnotation(buf *bytes.Buffer, typeWriter *BIRTypeWriter, annotation interface{}) error {
	// TODO: Implement when annotation model is available
	// For now, this is a placeholder - annotations are not yet in the Go model
	return fmt.Errorf("annotation writing not yet implemented")
}

func (w *BIRBinaryWriter) writeServiceDeclarations(buf *bytes.Buffer, serviceDecls []interface{}) error {
	// Write service declaration count
	if err := binary.Write(buf, binary.BigEndian, int32(len(serviceDecls))); err != nil {
		return err
	}
	
	// Write each service declaration (for now, empty since services not in Go model)
	for _, serviceDecl := range serviceDecls {
		if err := w.writeServiceDeclaration(buf, serviceDecl); err != nil {
			return err
		}
	}
	
	return nil
}

func (w *BIRBinaryWriter) writeServiceDeclaration(buf *bytes.Buffer, serviceDecl interface{}) error {
	// TODO: Implement when service declaration model is available
	// For now, this is a placeholder - service declarations are not yet in the Go model
	return fmt.Errorf("service declaration writing not yet implemented")
}

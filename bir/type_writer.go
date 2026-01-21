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

// BIRTypeWriter writes types to a buffer in binary format
type BIRTypeWriter struct {
	buf    *bytes.Buffer
	cp     *ConstantPool
	typeEnv interface{} // Type environment
}

// NewBIRTypeWriter creates a new type writer
func NewBIRTypeWriter(buf *bytes.Buffer, cp *ConstantPool, typeEnv interface{}) *BIRTypeWriter {
	return &BIRTypeWriter{
		buf:     buf,
		cp:      cp,
		typeEnv: typeEnv,
	}
}

// VisitType writes a type to the buffer
// Accepts both TypeNode and ValueType
func (tw *BIRTypeWriter) VisitType(t interface{}) error {
	if t == nil {
		return fmt.Errorf("type cannot be nil")
	}
	
	// Handle TypeNode
	if typeNode, ok := t.(model.TypeNode); ok {
		return tw.visitTypeNode(typeNode)
	}
	
	// Handle ValueType (runtime type)
	if valueType, ok := t.(model.ValueType); ok {
		return tw.visitValueType(valueType)
	}
	
	return fmt.Errorf("unsupported type: %T", t)
}

func (tw *BIRTypeWriter) visitTypeNode(t model.TypeNode) error {
	// Get type tag
	tag := getTypeTagFromTypeNode(t)
	if err := binary.Write(tw.buf, binary.BigEndian, byte(tag)); err != nil {
		return err
	}
	
	// Get type name
	typeName := getTypeName(t)
	nameCPIndex := AddStringCPEntry(typeName, tw.cp)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	
	// Get type flags
	flags := getTypeFlags(t)
	if err := binary.Write(tw.buf, binary.BigEndian, flags); err != nil {
		return err
	}
	
	// Write type-specific data based on tag
	return tw.writeTypeSpecificData(t, tag)
}

func (tw *BIRTypeWriter) visitValueType(t model.ValueType) error {
	// For ValueType, we need to serialize it differently
	// This is a placeholder - proper implementation requires understanding the ValueType structure
	// For now, write minimal type information
	typeKind := t.GetTypeKind()
	
	// Write a placeholder tag (this needs to be mapped from TypeKind to type tag)
	tag := int(typeKindToTag(typeKind))
	if err := binary.Write(tw.buf, binary.BigEndian, byte(tag)); err != nil {
		return err
	}
	
	// Write empty name for now
	nameCPIndex := AddStringCPEntry("", tw.cp)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	
	// Write zero flags
	if err := binary.Write(tw.buf, binary.BigEndian, int64(0)); err != nil {
		return err
	}
	
	// Type-specific data would go here
	return nil
}

// typeKindToTag converts TypeKind to a type tag (simplified mapping)
func typeKindToTag(kind model.TypeKind) int {
	// This is a simplified mapping - actual mapping should be based on TypeKind enum values
	return 0
}

// writeTypeSpecificData writes type-specific data based on the type tag
func (tw *BIRTypeWriter) writeTypeSpecificData(t model.TypeNode, tag int) error {
	// This is a simplified implementation
	// The full implementation would need to handle all type kinds
	// For now, we'll implement the most common cases
	
	switch tag {
	case 1: // ARRAY
		return tw.writeArrayType(t)
	case 2: // ANY
		// Nothing to write for ANY type
		return nil
	case 3: // ERROR
		return tw.writeErrorType(t)
	case 4: // FINITE
		return tw.writeFiniteType(t)
	case 5: // INVOKABLE
		return tw.writeInvokableType(t)
	case 6: // JSON
		// Nothing to write
		return nil
	case 7: // MAP
		return tw.writeMapType(t)
	case 8: // STREAM
		return tw.writeStreamType(t)
	case 9: // TYPEDESC
		return tw.writeTypedescType(t)
	case 10: // TYPE_REFERENCE
		return tw.writeTypeReferenceType(t)
	case 11: // PARAMETERIZED
		return tw.writeParameterizedType(t)
	case 12: // FUTURE
		return tw.writeFutureType(t)
	case 13: // HANDLE
		// Nothing to write
		return nil
	case 14: // NEVER
		// Nothing to write
		return nil
	case 15: // NIL
		// Nothing to write
		return nil
	case 16: // NO_TYPE
		// Nothing to write
		return nil
	case 17: // ANYDATA
		// Nothing to write
		return nil
	case 18: // PACKAGE
		return fmt.Errorf("package type serialization not implemented")
	case 19: // STRUCTURE
		return fmt.Errorf("structure type serialization not implemented")
	case 20: // TUPLE
		return tw.writeTupleType(t)
	case 21: // UNION
		return tw.writeUnionType(t)
	case 22: // INTERSECTION
		return tw.writeIntersectionType(t)
	case 23: // RECORD
		return tw.writeRecordType(t)
	case 24: // OBJECT
		return tw.writeObjectType(t)
	case 25: // XML
		return tw.writeXMLType(t)
	case 26: // TABLE
		return tw.writeTableType(t)
	default:
		return fmt.Errorf("unsupported type tag: %d", tag)
	}
}

// Helper methods for writing specific types
func (tw *BIRTypeWriter) writeArrayType(t model.TypeNode) error {
	// Get array state and size
	// This requires access to array-specific fields
	// For now, write defaults
	state := byte(0) // Default state
	if err := binary.Write(tw.buf, binary.BigEndian, state); err != nil {
		return err
	}
	
	size := int32(-1) // -1 means variable size
	if err := binary.Write(tw.buf, binary.BigEndian, size); err != nil {
		return err
	}
	
	// Get element type
	elemType := getArrayElementType(t)
	return WriteType(tw.cp, tw.buf, elemType)
}

func (tw *BIRTypeWriter) writeErrorType(t model.TypeNode) error {
	// Write error package and type name
	// This requires access to error-specific fields
	// For now, write placeholder values
	pkgIndex := int32(-1)
	if err := binary.Write(tw.buf, binary.BigEndian, pkgIndex); err != nil {
		return err
	}
	
	typeNameCPIndex := int32(-1)
	if err := binary.Write(tw.buf, binary.BigEndian, typeNameCPIndex); err != nil {
		return err
	}
	
	// Write detail type
	detailType := getErrorDetailType(t)
	return WriteType(tw.cp, tw.buf, detailType)
}

func (tw *BIRTypeWriter) writeFiniteType(t model.TypeNode) error {
	// Write finite type information
	// This requires access to finite type-specific fields
	typeName := getTypeName(t)
	nameCPIndex := AddStringCPEntry(typeName, tw.cp)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}
	
	flags := getTypeFlags(t)
	if err := binary.Write(tw.buf, binary.BigEndian, flags); err != nil {
		return err
	}
	
	// Write value space count (simplified)
	valueSpaceCount := int32(0)
	return binary.Write(tw.buf, binary.BigEndian, valueSpaceCount)
}

func (tw *BIRTypeWriter) writeInvokableType(t model.TypeNode) error {
	// Check if it's an any function
	isAnyFunction := false
	if err := binary.Write(tw.buf, binary.BigEndian, isAnyFunction); err != nil {
		return err
	}
	
	if isAnyFunction {
		return nil
	}
	
	// Write parameter types
	paramTypes := getInvokableParamTypes(t)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(len(paramTypes))); err != nil {
		return err
	}
	for _, paramType := range paramTypes {
		if err := WriteType(tw.cp, tw.buf, paramType); err != nil {
			return err
		}
	}
	
	// Write rest type
	restType := getInvokableRestType(t)
	hasRestType := restType != nil
	if err := binary.Write(tw.buf, binary.BigEndian, hasRestType); err != nil {
		return err
	}
	if hasRestType {
		if err := WriteType(tw.cp, tw.buf, restType); err != nil {
			return err
		}
	}
	
	// Write return type
	retType := getInvokableReturnType(t)
	return WriteType(tw.cp, tw.buf, retType)
}

func (tw *BIRTypeWriter) writeMapType(t model.TypeNode) error {
	constraint := getMapConstraintType(t)
	return WriteType(tw.cp, tw.buf, constraint)
}

func (tw *BIRTypeWriter) writeStreamType(t model.TypeNode) error {
	constraint := getStreamConstraintType(t)
	if err := WriteType(tw.cp, tw.buf, constraint); err != nil {
		return err
	}
	completionType := getStreamCompletionType(t)
	return WriteType(tw.cp, tw.buf, completionType)
}

func (tw *BIRTypeWriter) writeTypedescType(t model.TypeNode) error {
	constraint := getTypedescConstraintType(t)
	return WriteType(tw.cp, tw.buf, constraint)
}

func (tw *BIRTypeWriter) writeTypeReferenceType(t model.TypeNode) error {
	// Write package index
	pkgIndex := int32(-1)
	if err := binary.Write(tw.buf, binary.BigEndian, pkgIndex); err != nil {
		return err
	}
	
	// Write definition name
	defName := getTypeReferenceDefinitionName(t)
	defNameCPIndex := AddStringCPEntry(defName, tw.cp)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(defNameCPIndex)); err != nil {
		return err
	}
	
	// Write referred type
	referredType := getTypeReferenceReferredType(t)
	return WriteType(tw.cp, tw.buf, referredType)
}

func (tw *BIRTypeWriter) writeParameterizedType(t model.TypeNode) error {
	paramValueType := getParameterizedParamValueType(t)
	if err := WriteType(tw.cp, tw.buf, paramValueType); err != nil {
		return err
	}
	paramIndex := getParameterizedParamIndex(t)
	return binary.Write(tw.buf, binary.BigEndian, int32(paramIndex))
}

func (tw *BIRTypeWriter) writeFutureType(t model.TypeNode) error {
	constraint := getFutureConstraintType(t)
	return WriteType(tw.cp, tw.buf, constraint)
}

func (tw *BIRTypeWriter) writeTupleType(t model.TypeNode) error {
	members := getTupleMembers(t)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(len(members))); err != nil {
		return err
	}
	
	for i, member := range members {
		// Write member index as string
		indexStr := fmt.Sprintf("%d", i)
		indexCPIndex := AddStringCPEntry(indexStr, tw.cp)
		if err := binary.Write(tw.buf, binary.BigEndian, int32(indexCPIndex)); err != nil {
			return err
		}
		
		// Write member flags
		memberFlags := getTupleMemberFlags(member)
		if err := binary.Write(tw.buf, binary.BigEndian, memberFlags); err != nil {
			return err
		}
		
		// Write member type
		memberType := getTupleMemberType(member)
		if err := WriteType(tw.cp, tw.buf, memberType); err != nil {
			return err
		}
		
		// Write annotations (empty for now)
		emptyAnnots := []interface{}{}
		if err := WriteAnnotAttachments(tw.cp, tw.buf, emptyAnnots); err != nil {
			return err
		}
	}
	
	// Write rest type
	restType := getTupleRestType(t)
	hasRestType := restType != nil
	return binary.Write(tw.buf, binary.BigEndian, hasRestType)
}

func (tw *BIRTypeWriter) writeUnionType(t model.TypeNode) error {
	// Write cyclic flag
	isCyclic := false
	if err := binary.Write(tw.buf, binary.BigEndian, isCyclic); err != nil {
		return err
	}
	
	// Write symbol info if cyclic
	hasSymbol := false
	if err := binary.Write(tw.buf, binary.BigEndian, hasSymbol); err != nil {
		return err
	}
	
	// Write members
	memberTypes := getUnionMemberTypes(t)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(len(memberTypes))); err != nil {
		return err
	}
	for _, memberType := range memberTypes {
		if err := WriteType(tw.cp, tw.buf, memberType); err != nil {
			return err
		}
	}
	
	// Write original members
	originalMemberTypes := getUnionOriginalMemberTypes(t)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(len(originalMemberTypes))); err != nil {
		return err
	}
	for _, memberType := range originalMemberTypes {
		if err := WriteType(tw.cp, tw.buf, memberType); err != nil {
			return err
		}
	}
	
	// Write enum symbol info (false for now)
	hasEnumSymbol := false
	return binary.Write(tw.buf, binary.BigEndian, hasEnumSymbol)
}

func (tw *BIRTypeWriter) writeIntersectionType(t model.TypeNode) error {
	// Write constituent types
	constituentTypes := getIntersectionConstituentTypes(t)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(len(constituentTypes))); err != nil {
		return err
	}
	for _, constituentType := range constituentTypes {
		if err := WriteType(tw.cp, tw.buf, constituentType); err != nil {
			return err
		}
	}
	
	// Write effective type
	effectiveType := getIntersectionEffectiveType(t)
	return WriteType(tw.cp, tw.buf, effectiveType)
}

func (tw *BIRTypeWriter) writeRecordType(t model.TypeNode) error {
	// Write package index
	pkgIndex := int32(-1)
	if err := binary.Write(tw.buf, binary.BigEndian, pkgIndex); err != nil {
		return err
	}
	
	// Write type definition name
	defName := getTypeName(t)
	defNameCPIndex := AddStringCPEntry(defName, tw.cp)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(defNameCPIndex)); err != nil {
		return err
	}
	
	// Write sealed flag
	sealed := false
	if err := binary.Write(tw.buf, binary.BigEndian, sealed); err != nil {
		return err
	}
	
	// Write rest field type
	restFieldType := getRecordRestFieldType(t)
	if err := WriteType(tw.cp, tw.buf, restFieldType); err != nil {
		return err
	}
	
	// Write fields
	fields := getRecordFields(t)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(len(fields))); err != nil {
		return err
	}
	
	for _, field := range fields {
		fieldName := getRecordFieldName(field)
		fieldNameCPIndex := AddStringCPEntry(fieldName, tw.cp)
		if err := binary.Write(tw.buf, binary.BigEndian, int32(fieldNameCPIndex)); err != nil {
			return err
		}
		
		fieldFlags := getRecordFieldFlags(field)
		if err := binary.Write(tw.buf, binary.BigEndian, fieldFlags); err != nil {
			return err
		}
		
		// Write markdown doc (empty for now)
		emptyDoc := []byte{0, 0, 0, 0}
		if _, err := tw.buf.Write(emptyDoc); err != nil {
			return err
		}
		
		fieldType := getRecordFieldType(field)
		if err := WriteType(tw.cp, tw.buf, fieldType); err != nil {
			return err
		}
		
		// Write annotations (empty for now)
		emptyAnnots := []interface{}{}
		if err := WriteAnnotAttachments(tw.cp, tw.buf, emptyAnnots); err != nil {
			return err
		}
	}
	
	// Write type inclusions (empty for now)
	inclusionsCount := int32(0)
	return binary.Write(tw.buf, binary.BigEndian, inclusionsCount)
}

func (tw *BIRTypeWriter) writeObjectType(t model.TypeNode) error {
	// Write package index
	pkgIndex := int32(-1)
	if err := binary.Write(tw.buf, binary.BigEndian, pkgIndex); err != nil {
		return err
	}
	
	// Write type definition name
	defName := getTypeName(t)
	defNameCPIndex := AddStringCPEntry(defName, tw.cp)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(defNameCPIndex)); err != nil {
		return err
	}
	
	// Write flags
	flags := getTypeFlags(t)
	if err := binary.Write(tw.buf, binary.BigEndian, flags); err != nil {
		return err
	}
	
	// Write fields
	fields := getObjectFields(t)
	if err := binary.Write(tw.buf, binary.BigEndian, int32(len(fields))); err != nil {
		return err
	}
	
	for _, field := range fields {
		fieldName := getObjectFieldName(field)
		fieldNameCPIndex := AddStringCPEntry(fieldName, tw.cp)
		if err := binary.Write(tw.buf, binary.BigEndian, int32(fieldNameCPIndex)); err != nil {
			return err
		}
		
		fieldFlags := getObjectFieldFlags(field)
		if err := binary.Write(tw.buf, binary.BigEndian, fieldFlags); err != nil {
			return err
		}
		
		isDefaultable := getObjectFieldIsDefaultable(field)
		if err := binary.Write(tw.buf, binary.BigEndian, isDefaultable); err != nil {
			return err
		}
		
		// Write markdown doc (empty for now)
		emptyDoc := []byte{0, 0, 0, 0}
		if _, err := tw.buf.Write(emptyDoc); err != nil {
			return err
		}
		
		fieldType := getObjectFieldType(field)
		if err := WriteType(tw.cp, tw.buf, fieldType); err != nil {
			return err
		}
	}
	
	// Write initializer and constructor functions (simplified)
	hasGeneratedInitializer := false
	if err := binary.Write(tw.buf, binary.BigEndian, hasGeneratedInitializer); err != nil {
		return err
	}
	hasInitializer := false
	if err := binary.Write(tw.buf, binary.BigEndian, hasInitializer); err != nil {
		return err
	}
	
	// Write attached functions count
	attachedFuncsCount := int32(0)
	if err := binary.Write(tw.buf, binary.BigEndian, attachedFuncsCount); err != nil {
		return err
	}
	
	// Write type inclusions (empty for now)
	inclusionsCount := int32(0)
	return binary.Write(tw.buf, binary.BigEndian, inclusionsCount)
}

func (tw *BIRTypeWriter) writeXMLType(t model.TypeNode) error {
	constraint := getXMLConstraintType(t)
	return WriteType(tw.cp, tw.buf, constraint)
}

func (tw *BIRTypeWriter) writeTableType(t model.TypeNode) error {
	constraint := getTableConstraintType(t)
	if err := WriteType(tw.cp, tw.buf, constraint); err != nil {
		return err
	}
	
	// Write field name list
	hasFieldNameList := false
	if err := binary.Write(tw.buf, binary.BigEndian, hasFieldNameList); err != nil {
		return err
	}
	
	// Write key type constraint
	hasKeyTypeConstraint := false
	if err := binary.Write(tw.buf, binary.BigEndian, hasKeyTypeConstraint); err != nil {
		return err
	}
	
	return nil
}

// WriteMarkdownDocAttachment writes markdown documentation attachment
func (tw *BIRTypeWriter) WriteMarkdownDocAttachment(buf *bytes.Buffer, doc model.MarkdownDocAttachment) error {
	docBuf := &bytes.Buffer{}
	
	// MarkdownDocAttachment is a struct, check if it's empty
	hasDoc := !isMarkdownDocAttachmentEmpty(doc)
	if err := binary.Write(docBuf, binary.BigEndian, hasDoc); err != nil {
		return err
	}
	
	if hasDoc {
		// Write description
		desc := getMarkdownDescription(doc)
		descCPIndex := int32(-1)
		if desc != "" {
			descCPIndex = int32(AddStringCPEntry(desc, tw.cp))
		}
		if err := binary.Write(docBuf, binary.BigEndian, descCPIndex); err != nil {
			return err
		}
		
		// Write return value description
		retDesc := getMarkdownReturnValueDescription(doc)
		retDescCPIndex := int32(-1)
		if retDesc != "" {
			retDescCPIndex = int32(AddStringCPEntry(retDesc, tw.cp))
		}
		if err := binary.Write(docBuf, binary.BigEndian, retDescCPIndex); err != nil {
			return err
		}
		
		// Write parameters
		paramCount := int32(len(doc.Parameters))
		if err := binary.Write(docBuf, binary.BigEndian, paramCount); err != nil {
			return err
		}
		for _, param := range doc.Parameters {
			paramName := ""
			if param.Name != nil {
				paramName = *param.Name
			}
			paramNameCPIndex := AddStringCPEntry(paramName, tw.cp)
			if err := binary.Write(docBuf, binary.BigEndian, int32(paramNameCPIndex)); err != nil {
				return err
			}
			paramDesc := ""
			if param.Description != nil {
				paramDesc = *param.Description
			}
			paramDescCPIndex := AddStringCPEntry(paramDesc, tw.cp)
			if err := binary.Write(docBuf, binary.BigEndian, int32(paramDescCPIndex)); err != nil {
				return err
			}
		}
		
		// Write deprecated documentation
		deprecatedDoc := getMarkdownDeprecatedDocumentation(doc)
		deprecatedCPIndex := int32(-1)
		if deprecatedDoc != "" {
			deprecatedCPIndex = int32(AddStringCPEntry(deprecatedDoc, tw.cp))
		}
		if err := binary.Write(docBuf, binary.BigEndian, deprecatedCPIndex); err != nil {
			return err
		}
		
		// Write deprecated params
		deprecatedParamCount := int32(len(doc.DeprecatedParameters))
		if err := binary.Write(docBuf, binary.BigEndian, deprecatedParamCount); err != nil {
			return err
		}
		for _, param := range doc.DeprecatedParameters {
			paramName := ""
			if param.Name != nil {
				paramName = *param.Name
			}
			paramNameCPIndex := AddStringCPEntry(paramName, tw.cp)
			if err := binary.Write(docBuf, binary.BigEndian, int32(paramNameCPIndex)); err != nil {
				return err
			}
			paramDesc := ""
			if param.Description != nil {
				paramDesc = *param.Description
			}
			paramDescCPIndex := AddStringCPEntry(paramDesc, tw.cp)
			if err := binary.Write(docBuf, binary.BigEndian, int32(paramDescCPIndex)); err != nil {
				return err
			}
		}
	}
	
	// Write length and then the doc buffer
	length := int32(docBuf.Len())
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(docBuf.Bytes())
	return err
}

// Helper functions to extract type information
// These are placeholders that should be implemented based on actual type interfaces

func getTypeTagFromTypeNode(t model.TypeNode) int {
	if tagger, ok := t.(interface{ GetTag() int }); ok {
		return tagger.GetTag()
	}
	return 0
}

func getTypeName(t model.TypeNode) string {
	if namer, ok := t.(interface{ GetName() model.Name }); ok {
		name := namer.GetName()
		// Name is a string type alias, not a pointer
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

// Type-specific getters (simplified - should be implemented based on actual type interfaces)
func getArrayElementType(t model.TypeNode) model.TypeNode { return nil }
func getErrorDetailType(t model.TypeNode) model.TypeNode  { return nil }
func getInvokableParamTypes(t model.TypeNode) []model.TypeNode {
	// Try to get param types from invokable type
	// This depends on the actual InvokableType interface
	// For now, return empty slice
	return []model.TypeNode{}
}
func getInvokableRestType(t model.TypeNode) model.TypeNode { return nil }
func getInvokableReturnType(t model.TypeNode) model.TypeNode { return nil }
func getMapConstraintType(t model.TypeNode) model.TypeNode { return nil }
func getStreamConstraintType(t model.TypeNode) model.TypeNode { return nil }
func getStreamCompletionType(t model.TypeNode) model.TypeNode { return nil }
func getTypedescConstraintType(t model.TypeNode) model.TypeNode { return nil }
func getTypeReferenceDefinitionName(t model.TypeNode) string { return "" }
func getTypeReferenceReferredType(t model.TypeNode) model.TypeNode { return nil }
func getParameterizedParamValueType(t model.TypeNode) model.TypeNode { return nil }
func getParameterizedParamIndex(t model.TypeNode) int { return 0 }
func getFutureConstraintType(t model.TypeNode) model.TypeNode { return nil }
func getTupleMembers(t model.TypeNode) []interface{} { return []interface{}{} }
func getTupleMemberFlags(member interface{}) int64 { return 0 }
func getTupleMemberType(member interface{}) model.TypeNode { return nil }
func getTupleRestType(t model.TypeNode) model.TypeNode { return nil }
func getUnionMemberTypes(t model.TypeNode) []model.TypeNode { return []model.TypeNode{} }
func getUnionOriginalMemberTypes(t model.TypeNode) []model.TypeNode { return []model.TypeNode{} }
func getIntersectionConstituentTypes(t model.TypeNode) []model.TypeNode { return []model.TypeNode{} }
func getIntersectionEffectiveType(t model.TypeNode) model.TypeNode { return nil }
func getRecordRestFieldType(t model.TypeNode) model.TypeNode { return nil }
func getRecordFields(t model.TypeNode) []interface{} { return []interface{}{} }
func getRecordFieldName(field interface{}) string { return "" }
func getRecordFieldFlags(field interface{}) int64 { return 0 }
func getRecordFieldType(field interface{}) model.TypeNode { return nil }
func getObjectFields(t model.TypeNode) []interface{} { return []interface{}{} }
func getObjectFieldName(field interface{}) string { return "" }
func getObjectFieldFlags(field interface{}) int64 { return 0 }
func getObjectFieldIsDefaultable(field interface{}) bool { return false }
func getObjectFieldType(field interface{}) model.TypeNode { return nil }
func getXMLConstraintType(t model.TypeNode) model.TypeNode { return nil }
func getTableConstraintType(t model.TypeNode) model.TypeNode { return nil }
// isMarkdownDocAttachmentEmpty checks if a markdown doc attachment is empty
func isMarkdownDocAttachmentEmpty(doc model.MarkdownDocAttachment) bool {
	// MarkdownDocAttachment is a struct, check if all fields are empty
	// Check if description is nil/empty
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

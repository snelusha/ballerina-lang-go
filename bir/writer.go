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
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package bir

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"ballerina-lang-go/model"
	"ballerina-lang-go/tools/diagnostics"
)

const (
	BIR_MAGIC_0 = 0xba
	BIR_MAGIC_1 = 0x10
	BIR_MAGIC_2 = 0xc0
	BIR_MAGIC_3 = 0xde

	BIR_VERSION = 75
)

// CPEntryType represents the type of a constant pool entry
type CPEntryType byte

const (
	CP_ENTRY_INTEGER CPEntryType = iota + 1
	CP_ENTRY_FLOAT
	CP_ENTRY_BOOLEAN
	CP_ENTRY_STRING
	CP_ENTRY_PACKAGE
	CP_ENTRY_BYTE
	CP_ENTRY_SHAPE
)

// CPEntry represents a constant pool entry
type CPEntry interface {
	EntryType() CPEntryType
}

type (
	StringCPEntry struct {
		value string
	}

	IntegerCPEntry struct {
		value int64
	}

	FloatCPEntry struct {
		value float64
	}

	BooleanCPEntry struct {
		value bool
	}

	ByteCPEntry struct {
		value int32
	}

	PackageCPEntry struct {
		orgNameCPIndex    int
		pkgNameCPIndex    int
		moduleNameCPIndex int
		versionCPIndex    int
	}

	ShapeCPEntry struct {
		shape model.ValueType
	}
)

func (e *StringCPEntry) EntryType() CPEntryType  { return CP_ENTRY_STRING }
func (e *IntegerCPEntry) EntryType() CPEntryType { return CP_ENTRY_INTEGER }
func (e *FloatCPEntry) EntryType() CPEntryType   { return CP_ENTRY_FLOAT }
func (e *BooleanCPEntry) EntryType() CPEntryType { return CP_ENTRY_BOOLEAN }
func (e *ByteCPEntry) EntryType() CPEntryType    { return CP_ENTRY_BYTE }
func (e *PackageCPEntry) EntryType() CPEntryType { return CP_ENTRY_PACKAGE }
func (e *ShapeCPEntry) EntryType() CPEntryType   { return CP_ENTRY_SHAPE }

// ConstantPool manages constant pool entries
type ConstantPool struct {
	entries     []CPEntry
	entriesMap  map[string]int // key -> index for deduplication
	shapeWriter *BIRTypeWriter
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		entries:    make([]CPEntry, 0),
		entriesMap: make(map[string]int),
	}
}

func (cp *ConstantPool) addCPEntry(entry CPEntry) int {
	// Create a key for deduplication based on entry type and value
	key := cp.entryKey(entry)
	if idx, exists := cp.entriesMap[key]; exists {
		return idx
	}

	idx := len(cp.entries)
	cp.entries = append(cp.entries, entry)
	cp.entriesMap[key] = idx
	return idx
}

func (cp *ConstantPool) entryKey(entry CPEntry) string {
	switch e := entry.(type) {
	case *StringCPEntry:
		return fmt.Sprintf("str:%s", e.value)
	case *IntegerCPEntry:
		return fmt.Sprintf("int:%d", e.value)
	case *FloatCPEntry:
		return fmt.Sprintf("float:%g", e.value)
	case *BooleanCPEntry:
		return fmt.Sprintf("bool:%v", e.value)
	case *ByteCPEntry:
		return fmt.Sprintf("byte:%d", e.value)
	case *PackageCPEntry:
		return fmt.Sprintf("pkg:%d:%d:%d:%d", e.orgNameCPIndex, e.pkgNameCPIndex, e.moduleNameCPIndex, e.versionCPIndex)
	case *ShapeCPEntry:
		// For shape entries, we use a simple approach - in practice, type equality
		// would need proper type comparison
		return fmt.Sprintf("shape:%p", e.shape)
	default:
		return fmt.Sprintf("unknown:%p", entry)
	}
}

func (cp *ConstantPool) addStringCPEntry(value string) int {
	return cp.addCPEntry(&StringCPEntry{value: value})
}

func (cp *ConstantPool) addPkgCPEntry(pkgID *model.PackageID) int {
	if pkgID == nil {
		// Return empty package CP entry
		empty := ""
		orgCPIndex := cp.addStringCPEntry(empty)
		return cp.addCPEntry(&PackageCPEntry{
			orgNameCPIndex:    orgCPIndex,
			pkgNameCPIndex:    orgCPIndex,
			moduleNameCPIndex: orgCPIndex,
			versionCPIndex:    orgCPIndex,
		})
	}

	orgName := ""
	if pkgID.OrgName != nil {
		orgName = pkgID.OrgName.Value()
	}
	pkgName := ""
	if pkgID.PkgName != nil {
		pkgName = pkgID.PkgName.Value()
	}
	name := ""
	if pkgID.Name != nil {
		name = pkgID.Name.Value()
	}
	version := ""
	if pkgID.Version != nil {
		version = pkgID.Version.Value()
	}

	orgCPIndex := cp.addStringCPEntry(orgName)
	pkgNameCPIndex := cp.addStringCPEntry(pkgName)
	moduleNameCPIndex := cp.addStringCPEntry(name)
	versionCPIndex := cp.addStringCPEntry(version)
	return cp.addCPEntry(&PackageCPEntry{
		orgNameCPIndex:    orgCPIndex,
		pkgNameCPIndex:    pkgNameCPIndex,
		moduleNameCPIndex: moduleNameCPIndex,
		versionCPIndex:    versionCPIndex,
	})
}

func (cp *ConstantPool) addShapeCPEntry(shape model.ValueType) int {
	return cp.addCPEntry(&ShapeCPEntry{shape: shape})
}

func (cp *ConstantPool) serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// First write -1 as placeholder for size (will be overwritten)
	sizePlaceholder := make([]byte, 4)
	binary.BigEndian.PutUint32(sizePlaceholder, 0xFFFFFFFF)
	buf.Write(sizePlaceholder)

	// Write all entries
	for _, entry := range cp.entries {
		if err := cp.writeCPEntry(buf, entry); err != nil {
			return nil, err
		}
	}

	// Overwrite size at the beginning
	result := buf.Bytes()
	size := int32(len(cp.entries))
	result[0] = byte((uint32(size) >> 24) & 0xFF)
	result[1] = byte((uint32(size) >> 16) & 0xFF)
	result[2] = byte((uint32(size) >> 8) & 0xFF)
	result[3] = byte(uint32(size) & 0xFF)

	return result, nil
}

func (cp *ConstantPool) writeCPEntry(w io.Writer, entry CPEntry) error {
	// Write entry type
	if _, err := w.Write([]byte{byte(entry.EntryType())}); err != nil {
		return err
	}

	switch e := entry.(type) {
	case *StringCPEntry:
		// In Go, strings can't be nil, so we check for empty string
		// But we should never create null string entries in practice
		// For now, treat empty string as a valid string (not null)
		bytes := []byte(e.value)
		if err := binary.Write(w, binary.BigEndian, int32(len(bytes))); err != nil {
			return err
		}
		_, err := w.Write(bytes)
		return err

	case *IntegerCPEntry:
		return binary.Write(w, binary.BigEndian, e.value)

	case *FloatCPEntry:
		return binary.Write(w, binary.BigEndian, e.value)

	case *BooleanCPEntry:
		val := byte(0)
		if e.value {
			val = 1
		}
		_, err := w.Write([]byte{val})
		return err

	case *ByteCPEntry:
		return binary.Write(w, binary.BigEndian, e.value)

	case *PackageCPEntry:
		if err := binary.Write(w, binary.BigEndian, int32(e.orgNameCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(w, binary.BigEndian, int32(e.pkgNameCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(w, binary.BigEndian, int32(e.moduleNameCPIndex)); err != nil {
			return err
		}
		return binary.Write(w, binary.BigEndian, int32(e.versionCPIndex))

	case *ShapeCPEntry:
		// Write shape type - this requires type serialization
		// For now, we'll write a placeholder - full type serialization would require
		// implementing BIRTypeWriter similar to Java
		typeBuf := new(bytes.Buffer)
		if err := writeTypeToBuffer(typeBuf, cp, e.shape); err != nil {
			return err
		}
		typeBytes := typeBuf.Bytes()
		if err := binary.Write(w, binary.BigEndian, int32(len(typeBytes))); err != nil {
			return err
		}
		_, err := w.Write(typeBytes)
		return err

	default:
		return fmt.Errorf("unsupported CP entry type: %T", entry)
	}
}

// BIRBinaryWriter serializes BIRPackage to binary format
type BIRBinaryWriter struct {
	cp         *ConstantPool
	birPackage *BIRPackage
}

func NewBIRBinaryWriter(birPackage *BIRPackage) *BIRBinaryWriter {
	return &BIRBinaryWriter{
		cp:         NewConstantPool(),
		birPackage: birPackage,
	}
}

// Serialize writes the BIR package to binary format
func (w *BIRBinaryWriter) Serialize() ([]byte, error) {
	if w.birPackage == nil {
		return nil, fmt.Errorf("BIR package is nil")
	}
	if w.birPackage.PackageID == nil {
		return nil, fmt.Errorf("BIR package PackageID is nil")
	}

	birbuf := new(bytes.Buffer)

	// Write package ID as CP entry
	pkgCPIndex := w.cp.addPkgCPEntry(w.birPackage.PackageID)
	if err := binary.Write(birbuf, binary.BigEndian, int32(pkgCPIndex)); err != nil {
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

	// Write type defs (container only)
	if err := w.writeTypeDefs(birbuf); err != nil {
		return nil, err
	}

	// Write global vars
	if err := w.writeGlobalVars(birbuf); err != nil {
		return nil, err
	}

	// Write type def bodies
	if err := w.writeTypeDefBodies(birbuf); err != nil {
		return nil, err
	}

	// Write functions
	if err := w.writeFunctions(birbuf); err != nil {
		return nil, err
	}

	// Write annotations (if supported)
	// Note: Currently not fully supported in Go model
	if err := w.writeAnnotations(birbuf); err != nil {
		return nil, err
	}

	// Write service declarations (if supported)
	// Note: Currently not fully supported in Go model
	if err := w.writeServiceDeclarations(birbuf); err != nil {
		return nil, err
	}

	// Serialize constant pool
	cpBytes, err := w.cp.serialize()
	if err != nil {
		return nil, err
	}

	// Write final binary: magic + version + CP + module data
	result := new(bytes.Buffer)
	result.Write([]byte{BIR_MAGIC_0, BIR_MAGIC_1, BIR_MAGIC_2, BIR_MAGIC_3})
	if err := binary.Write(result, binary.BigEndian, int32(BIR_VERSION)); err != nil {
		return nil, err
	}
	result.Write(cpBytes)
	result.Write(birbuf.Bytes())

	return result.Bytes(), nil
}

func (w *BIRBinaryWriter) writeImportModuleDecls(buf *bytes.Buffer) error {
	imports := w.birPackage.ImportModules
	
	// Count valid imports (non-nil PackageID)
	validImports := 0
	for _, imp := range imports {
		if imp.PackageID != nil {
			validImports++
		}
	}
	
	if err := binary.Write(buf, binary.BigEndian, int32(validImports)); err != nil {
		return err
	}

	for _, imp := range imports {
		if imp.PackageID == nil {
			continue
		}
		// Use addPkgCPEntry which handles nil checks
		_ = w.cp.addPkgCPEntry(imp.PackageID)

		// Write individual indices (for import format)
		orgName := ""
		if imp.PackageID.OrgName != nil {
			orgName = imp.PackageID.OrgName.Value()
		}
		pkgName := ""
		if imp.PackageID.PkgName != nil {
			pkgName = imp.PackageID.PkgName.Value()
		}
		name := ""
		if imp.PackageID.Name != nil {
			name = imp.PackageID.Name.Value()
		}
		version := ""
		if imp.PackageID.Version != nil {
			version = imp.PackageID.Version.Value()
		}

		orgCPIndex := w.cp.addStringCPEntry(orgName)
		pkgNameCPIndex := w.cp.addStringCPEntry(pkgName)
		nameCPIndex := w.cp.addStringCPEntry(name)
		versionCPIndex := w.cp.addStringCPEntry(version)

		if err := binary.Write(buf, binary.BigEndian, int32(orgCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, int32(pkgNameCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, int32(versionCPIndex)); err != nil {
			return err
		}
	}
	return nil
}

func (w *BIRBinaryWriter) writeConstants(buf *bytes.Buffer) error {
	constants := w.birPackage.Constants
	if err := binary.Write(buf, binary.BigEndian, int32(len(constants))); err != nil {
		return err
	}

	for _, constant := range constants {
		if err := w.writeConstant(buf, constant); err != nil {
			return err
		}
	}
	return nil
}

func (w *BIRBinaryWriter) writeConstant(buf *bytes.Buffer, constant BIRConstant) error {
	// Name CP Index
	name := constant.Name.Value()
	nameCPIndex := w.cp.addStringCPEntry(name)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	// Flags
	if err := binary.Write(buf, binary.BigEndian, constant.Flags); err != nil {
		return err
	}

	// Origin
	if err := binary.Write(buf, binary.BigEndian, byte(constant.Origin)); err != nil {
		return err
	}

	// Position
	if err := w.writePosition(buf, constant.Pos); err != nil {
		return err
	}

	// Markdown doc attachment
	// Format: length (s4), has_doc (u1), content (if has_doc)
	docBuf := new(bytes.Buffer)
	hasDoc := false
	if err := binary.Write(docBuf, binary.BigEndian, byte(boolToByte(hasDoc))); err != nil {
		return err
	}
	docBytes := docBuf.Bytes()
	// Write length (including has_doc byte)
	if err := binary.Write(buf, binary.BigEndian, int32(len(docBytes))); err != nil {
		return err
	}
	buf.Write(docBytes)

	// Type
	if err := writeType(buf, w.cp, constant.Type); err != nil {
		return err
	}

	// Annotation attachments (written to main buffer, before constant value)
	// TODO: Implement annotation attachments writing
	annotBuf := new(bytes.Buffer)
	if err := binary.Write(annotBuf, binary.BigEndian, int32(0)); err != nil { // 0 annotations
		return err
	}
	annotBytes := annotBuf.Bytes()
	if err := binary.Write(buf, binary.BigEndian, int64(len(annotBytes))); err != nil {
		return err
	}
	buf.Write(annotBytes)

	// Constant value (written to separate buffer, then length + data to main buffer)
	// Always write type CP index (write -1 if type is nil, matching Java behavior)
	constValueBuf := new(bytes.Buffer)
	if err := writeType(constValueBuf, w.cp, constant.ConstValue.Type); err != nil {
		return err
	}
	if err := w.writeConstValue(constValueBuf, constant.ConstValue); err != nil {
		return err
	}

	// Constant value length and data
	constValueBytes := constValueBuf.Bytes()
	if err := binary.Write(buf, binary.BigEndian, int64(len(constValueBytes))); err != nil {
		return err
	}
	buf.Write(constValueBytes)

	return nil
}

func (w *BIRBinaryWriter) writeTypeDefs(buf *bytes.Buffer) error {
	typeDefs := w.birPackage.TypeDefs
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
	// Position
	if err := w.writePosition(buf, typeDef.Pos); err != nil {
		return err
	}

	// Internal name CP Index
	internalName := typeDef.InternalName.Value()
	internalNameCPIndex := w.cp.addStringCPEntry(internalName)
	if err := binary.Write(buf, binary.BigEndian, int32(internalNameCPIndex)); err != nil {
		return err
	}

	// Original name CP Index
	originalName := typeDef.OriginalName.Value()
	originalNameCPIndex := w.cp.addStringCPEntry(originalName)
	if err := binary.Write(buf, binary.BigEndian, int32(originalNameCPIndex)); err != nil {
		return err
	}

	// Flags
	if err := binary.Write(buf, binary.BigEndian, typeDef.Flags); err != nil {
		return err
	}

	// Origin
	if err := binary.Write(buf, binary.BigEndian, byte(typeDef.Origin)); err != nil {
		return err
	}

	// Markdown doc attachment
	// Format: length (s4), has_doc (u1), content (if has_doc)
	docBuf := new(bytes.Buffer)
	hasDoc := false
	if err := binary.Write(docBuf, binary.BigEndian, byte(boolToByte(hasDoc))); err != nil {
		return err
	}
	docBytes := docBuf.Bytes()
	// Write length (including has_doc byte)
	if err := binary.Write(buf, binary.BigEndian, int32(len(docBytes))); err != nil {
		return err
	}
	buf.Write(docBytes)

	// Type (TypeNode needs to be converted to ValueType if possible)
	// For now, write -1 as placeholder since TypeNode != ValueType
	if err := binary.Write(buf, binary.BigEndian, int32(-1)); err != nil {
		return err
	}

	// Has reference type
	hasReferenceType := typeDef.ReferenceType != nil
	if err := binary.Write(buf, binary.BigEndian, hasReferenceType); err != nil {
		return err
	}

	// Annotation attachments
	annotBuf := new(bytes.Buffer)
	if err := binary.Write(annotBuf, binary.BigEndian, int32(0)); err != nil { // 0 annotations
		return err
	}
	annotBytes := annotBuf.Bytes()
	if err := binary.Write(buf, binary.BigEndian, int64(len(annotBytes))); err != nil {
		return err
	}
	buf.Write(annotBytes)

	return nil
}

func (w *BIRBinaryWriter) writeTypeDefBodies(buf *bytes.Buffer) error {
	// Filter type defs to only OBJECT and RECORD types
	filtered := make([]BIRTypeDefinition, 0)
	for _, typeDef := range w.birPackage.TypeDefs {
		// Check if type is OBJECT or RECORD
		// TODO: Implement proper type tag checking
		// For now, we'll write empty list
		_ = typeDef
	}

	if err := binary.Write(buf, binary.BigEndian, int32(len(filtered))); err != nil {
		return err
	}

	for _, typeDef := range filtered {
		// Write attached functions
		if err := w.writeFunctions(buf); err != nil {
			return err
		}

		// Write referenced types
		if err := binary.Write(buf, binary.BigEndian, int32(len(typeDef.ReferencedTypes))); err != nil {
			return err
		}
		for range typeDef.ReferencedTypes {
			// TypeNode needs to be converted to ValueType if possible
			// For now, write -1 as placeholder
			if err := binary.Write(buf, binary.BigEndian, int32(-1)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *BIRBinaryWriter) writeGlobalVars(buf *bytes.Buffer) error {
	globalVars := w.birPackage.GlobalVars
	if err := binary.Write(buf, binary.BigEndian, int32(len(globalVars))); err != nil {
		return err
	}

	for _, gv := range globalVars {
		if err := w.writeGlobalVar(buf, gv); err != nil {
			return err
		}
	}
	return nil
}

func (w *BIRBinaryWriter) writeGlobalVar(buf *bytes.Buffer, gv BIRGlobalVariableDcl) error {
	// Position
	if err := w.writePosition(buf, gv.Pos); err != nil {
		return err
	}

	// Kind
	if err := binary.Write(buf, binary.BigEndian, byte(gv.Kind)); err != nil {
		return err
	}

	// Name CP Index
	name := gv.Name.Value()
	nameCPIndex := w.cp.addStringCPEntry(name)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	// Flags
	if err := binary.Write(buf, binary.BigEndian, gv.Flags); err != nil {
		return err
	}

	// Origin
	if err := binary.Write(buf, binary.BigEndian, byte(gv.Origin)); err != nil {
		return err
	}

	// Markdown doc attachment
	// Format: length (s4), has_doc (u1), content (if has_doc)
	docBuf := new(bytes.Buffer)
	hasDoc := false
	if err := binary.Write(docBuf, binary.BigEndian, byte(boolToByte(hasDoc))); err != nil {
		return err
	}
	docBytes := docBuf.Bytes()
	// Write length (including has_doc byte)
	if err := binary.Write(buf, binary.BigEndian, int32(len(docBytes))); err != nil {
		return err
	}
	buf.Write(docBytes)

	// Type
	if err := writeType(buf, w.cp, gv.Type); err != nil {
		return err
	}

	// Annotation attachments
	annotBuf := new(bytes.Buffer)
	if err := binary.Write(annotBuf, binary.BigEndian, int32(0)); err != nil { // 0 annotations
		return err
	}
	annotBytes := annotBuf.Bytes()
	if err := binary.Write(buf, binary.BigEndian, int64(len(annotBytes))); err != nil {
		return err
	}
	buf.Write(annotBytes)

	return nil
}

func (w *BIRBinaryWriter) writeFunctions(buf *bytes.Buffer) error {
	functions := w.birPackage.Functions
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
	// Position
	if err := w.writePosition(buf, fn.Pos); err != nil {
		return err
	}

	// Function name CP Index
	name := fn.Name.Value()
	nameCPIndex := w.cp.addStringCPEntry(name)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	// Function original name CP Index
	originalName := fn.OriginalName.Value()
	originalNameCPIndex := w.cp.addStringCPEntry(originalName)
	if err := binary.Write(buf, binary.BigEndian, int32(originalNameCPIndex)); err != nil {
		return err
	}

	// Worker name CP Index (empty for now)
	workerNameCPIndex := w.cp.addStringCPEntry("")
	if err := binary.Write(buf, binary.BigEndian, int32(workerNameCPIndex)); err != nil {
		return err
	}

	// Flags
	if err := binary.Write(buf, binary.BigEndian, fn.Flags); err != nil {
		return err
	}

	// Origin
	if err := binary.Write(buf, binary.BigEndian, byte(fn.Origin)); err != nil {
		return err
	}

	// Function type (InvokableType needs special handling)
	// For now, write -1 as placeholder
	if err := binary.Write(buf, binary.BigEndian, int32(-1)); err != nil {
		return err
	}

	// Path parameters (resource function)
	// TODO: Implement path parameters
	isResourceFunction := false
	if err := binary.Write(buf, binary.BigEndian, isResourceFunction); err != nil {
		return err
	}

	// Annotation attachments
	annotBuf := new(bytes.Buffer)
	if err := binary.Write(annotBuf, binary.BigEndian, int32(0)); err != nil { // 0 annotations
		return err
	}
	annotBytes := annotBuf.Bytes()
	if err := binary.Write(buf, binary.BigEndian, int64(len(annotBytes))); err != nil {
		return err
	}
	buf.Write(annotBytes)

	// Annotation attachments on external
	annotBuf2 := new(bytes.Buffer)
	if err := binary.Write(annotBuf2, binary.BigEndian, int32(0)); err != nil { // 0 annotations
		return err
	}
	annotBytes2 := annotBuf2.Bytes()
	if err := binary.Write(buf, binary.BigEndian, int64(len(annotBytes2))); err != nil {
		return err
	}
	buf.Write(annotBytes2)

	// Return type annotations
	annotBuf3 := new(bytes.Buffer)
	if err := binary.Write(annotBuf3, binary.BigEndian, int32(0)); err != nil { // 0 annotations
		return err
	}
	annotBytes3 := annotBuf3.Bytes()
	if err := binary.Write(buf, binary.BigEndian, int64(len(annotBytes3))); err != nil {
		return err
	}
	buf.Write(annotBytes3)

	// Required params
	if err := binary.Write(buf, binary.BigEndian, int32(len(fn.RequiredParams))); err != nil {
		return err
	}
	for _, param := range fn.RequiredParams {
		paramName := param.Name.Value()
		paramNameCPIndex := w.cp.addStringCPEntry(paramName)
		if err := binary.Write(buf, binary.BigEndian, int32(paramNameCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, param.Flags); err != nil {
			return err
		}
		// Parameter annotation attachments
		paramAnnotBuf := new(bytes.Buffer)
		if err := binary.Write(paramAnnotBuf, binary.BigEndian, int32(0)); err != nil { // 0 annotations
			return err
		}
		paramAnnotBytes := paramAnnotBuf.Bytes()
		if err := binary.Write(buf, binary.BigEndian, int64(len(paramAnnotBytes))); err != nil {
			return err
		}
		buf.Write(paramAnnotBytes)
	}

	// Rest parameter
	hasRestParam := fn.RestParams != nil
	if err := binary.Write(buf, binary.BigEndian, hasRestParam); err != nil {
		return err
	}
	if hasRestParam {
		restParamName := fn.RestParams.Name.Value()
		restParamNameCPIndex := w.cp.addStringCPEntry(restParamName)
		if err := binary.Write(buf, binary.BigEndian, int32(restParamNameCPIndex)); err != nil {
			return err
		}
		// Rest param annotation attachments
		restAnnotBuf := new(bytes.Buffer)
		if err := binary.Write(restAnnotBuf, binary.BigEndian, int32(0)); err != nil { // 0 annotations
			return err
		}
		restAnnotBytes := restAnnotBuf.Bytes()
		if err := binary.Write(buf, binary.BigEndian, int64(len(restAnnotBytes))); err != nil {
			return err
		}
		buf.Write(restAnnotBytes)
	}

	// Receiver
	hasReceiver := false // TODO: Implement receiver support
	if err := binary.Write(buf, binary.BigEndian, hasReceiver); err != nil {
		return err
	}

	// Markdown doc attachment
	// Format: length (s4), has_doc (u1), content (if has_doc)
	docBuf := new(bytes.Buffer)
	hasDoc := false
	if err := binary.Write(docBuf, binary.BigEndian, byte(boolToByte(hasDoc))); err != nil {
		return err
	}
	docBytes := docBuf.Bytes()
	// Write length (including has_doc byte)
	if err := binary.Write(buf, binary.BigEndian, int32(len(docBytes))); err != nil {
		return err
	}
	buf.Write(docBytes)

	// Dependent global vars
	if err := binary.Write(buf, binary.BigEndian, int32(len(fn.DependentGlobalVars))); err != nil {
		return err
	}
	for _, depVar := range fn.DependentGlobalVars {
		varName := depVar.Name.Value()
		varNameCPIndex := w.cp.addStringCPEntry(varName)
		if err := binary.Write(buf, binary.BigEndian, int32(varNameCPIndex)); err != nil {
			return err
		}
	}

	// Function body and scope buffer
	// In Java, both birbuf and scopebuf are created here and passed to BIRInstructionWriter
	funcBodyBuf := new(bytes.Buffer)
	scopeBuf := new(bytes.Buffer)
	
	// Create instruction writer with both buffers
	insWriter := NewBIRInstructionWriter(funcBodyBuf, scopeBuf, w.cp)
	
	if err := w.writeFunctionBody(funcBodyBuf, insWriter, fn); err != nil {
		return err
	}

	// Write scope length and data (BEFORE function body, matching Java line 316)
	// Format: scope_table_length (s8), scope_entry_count (s4), scope_entries
	scopeBytes := scopeBuf.Bytes()
	scopeLength := int64(len(scopeBytes) + 4) // +4 for scope count that we write next
	if err := binary.Write(buf, binary.BigEndian, scopeLength); err != nil {
		return err
	}
	scopeCount := insWriter.GetScopeCount()
	if err := binary.Write(buf, binary.BigEndian, int32(scopeCount)); err != nil {
		return err
	}
	buf.Write(scopeBytes)

	// Write function body length and data (AFTER scope table, matching Java line 319-321)
	funcBodyBytes := funcBodyBuf.Bytes()
	if err := binary.Write(buf, binary.BigEndian, int64(len(funcBodyBytes))); err != nil {
		return err
	}
	buf.Write(funcBodyBytes)

	return nil
}

func (w *BIRBinaryWriter) writeFunctionBody(buf *bytes.Buffer, insWriter *BIRInstructionWriter, fn BIRFunction) error {
	// Args count
	if err := binary.Write(buf, binary.BigEndian, int32(fn.ArgsCount)); err != nil {
		return err
	}

	// Return variable
	hasReturnVar := fn.ReturnVariable != nil
	if err := binary.Write(buf, binary.BigEndian, hasReturnVar); err != nil {
		return err
	}
	if hasReturnVar {
		if err := binary.Write(buf, binary.BigEndian, byte(fn.ReturnVariable.Kind)); err != nil {
			return err
		}
		if err := writeType(buf, w.cp, fn.ReturnVariable.Type); err != nil {
			return err
		}
		returnVarName := fn.ReturnVariable.Name.Value()
		returnVarNameCPIndex := w.cp.addStringCPEntry(returnVarName)
		if err := binary.Write(buf, binary.BigEndian, int32(returnVarNameCPIndex)); err != nil {
			return err
		}
	}

	// Function parameters (default parameters)
	if err := binary.Write(buf, binary.BigEndian, int32(len(fn.Parameters))); err != nil {
		return err
	}
	for _, param := range fn.Parameters {
		if err := binary.Write(buf, binary.BigEndian, byte(param.Kind)); err != nil {
			return err
		}
		if err := writeType(buf, w.cp, param.Type); err != nil {
			return err
		}
		paramName := param.Name.Value()
		paramNameCPIndex := w.cp.addStringCPEntry(paramName)
		if err := binary.Write(buf, binary.BigEndian, int32(paramNameCPIndex)); err != nil {
			return err
		}
		if param.Kind == VAR_KIND_ARG {
			metaVarName := param.MetaVarName
			if metaVarName == "" {
				metaVarName = ""
			}
			metaVarNameCPIndex := w.cp.addStringCPEntry(metaVarName)
			if err := binary.Write(buf, binary.BigEndian, int32(metaVarNameCPIndex)); err != nil {
				return err
			}
		}
		if err := binary.Write(buf, binary.BigEndian, param.HasDefaultExpr); err != nil {
			return err
		}
	}

	// Local variables
	if err := binary.Write(buf, binary.BigEndian, int32(len(fn.LocalVars))); err != nil {
		return err
	}
	for _, localVar := range fn.LocalVars {
		if err := binary.Write(buf, binary.BigEndian, byte(localVar.Kind)); err != nil {
			return err
		}
		if err := writeType(buf, w.cp, localVar.Type); err != nil {
			return err
		}
		localVarName := localVar.Name.Value()
		localVarNameCPIndex := w.cp.addStringCPEntry(localVarName)
		if err := binary.Write(buf, binary.BigEndian, int32(localVarNameCPIndex)); err != nil {
			return err
		}
		if localVar.Kind == VAR_KIND_ARG {
			metaVarName := localVar.MetaVarName
			if metaVarName == "" {
				metaVarName = ""
			}
			metaVarNameCPIndex := w.cp.addStringCPEntry(metaVarName)
			if err := binary.Write(buf, binary.BigEndian, int32(metaVarNameCPIndex)); err != nil {
				return err
			}
		}
		if localVar.Kind == VAR_KIND_LOCAL {
			metaVarName := localVar.MetaVarName
			if metaVarName == "" {
				metaVarName = ""
			}
			metaVarNameCPIndex := w.cp.addStringCPEntry(metaVarName)
			if err := binary.Write(buf, binary.BigEndian, int32(metaVarNameCPIndex)); err != nil {
				return err
			}
			endBBName := ""
			if localVar.EndBB != nil {
				endBBName = localVar.EndBB.Id.Value()
			}
			endBBNameCPIndex := w.cp.addStringCPEntry(endBBName)
			if err := binary.Write(buf, binary.BigEndian, int32(endBBNameCPIndex)); err != nil {
				return err
			}
			startBBName := ""
			if localVar.StartBB != nil {
				startBBName = localVar.StartBB.Id.Value()
			}
			startBBNameCPIndex := w.cp.addStringCPEntry(startBBName)
			if err := binary.Write(buf, binary.BigEndian, int32(startBBNameCPIndex)); err != nil {
				return err
			}
			if err := binary.Write(buf, binary.BigEndian, int32(localVar.InsOffset)); err != nil {
				return err
			}
		}
	}

	// Basic blocks (insWriter is passed in and already initialized)
	if err := insWriter.writeBBs(fn.BasicBlocks); err != nil {
		return err
	}

	// Error table
	// TODO: Implement error table writing
	if err := binary.Write(buf, binary.BigEndian, int32(0)); err != nil { // 0 error entries
		return err
	}

	// Worker channels
	// TODO: Implement worker channels
	if err := binary.Write(buf, binary.BigEndian, int32(0)); err != nil { // 0 channels
		return err
	}

	return nil
}

func (w *BIRBinaryWriter) writeAnnotations(buf *bytes.Buffer) error {
	// Annotations are not fully supported in Go model yet
	// Write empty list
	if err := binary.Write(buf, binary.BigEndian, int32(0)); err != nil {
		return err
	}
	return nil
}

func (w *BIRBinaryWriter) writeServiceDeclarations(buf *bytes.Buffer) error {
	// Service declarations are not fully supported in Go model yet
	// Write empty list
	if err := binary.Write(buf, binary.BigEndian, int32(0)); err != nil {
		return err
	}
	return nil
}

func (w *BIRBinaryWriter) writePosition(buf *bytes.Buffer, pos diagnostics.Location) error {
	sLine := int32(-2147483648) // Integer.MIN_VALUE
	eLine := int32(-2147483648)
	sCol := int32(-2147483648)
	eCol := int32(-2147483648)
	sourceFileName := ""

	if pos != nil {
		// TODO: Extract line/column information from Location
		// For now, use default values
	}

	sourceFileNameCPIndex := w.cp.addStringCPEntry(sourceFileName)
	if err := binary.Write(buf, binary.BigEndian, int32(sourceFileNameCPIndex)); err != nil {
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

func (w *BIRBinaryWriter) writeConstValue(buf *bytes.Buffer, constValue ConstValue) error {
	// The constant_value_type_cp_index has already been written before calling this function
	// Now we need to write constant_value_info based on the type tag
	// Since we're using type_tag_nil (10) as placeholder, constant_value_info should be nil_constant_info (size 0)
	// So we write nothing for type_tag_nil
	
	// If type is nil, write nil_constant_info (size 0, so nothing to write)
	if constValue.Type == nil {
		// For nil type, write nil_constant_info (size 0, so nothing to write)
		return nil
	}

	// Since we're using type_tag_nil (10) as placeholder for all types,
	// constant_value_info should be nil_constant_info (size 0)
	// So we write nothing
	// TODO: When we implement proper type tag extraction, we should write the appropriate constant_value_info
	// based on the actual type tag (int_constant_info, string_constant_info, etc.)
	return nil
}

// Helper functions

func writeType(buf *bytes.Buffer, cp *ConstantPool, t model.ValueType) error {
	if t == nil {
		return binary.Write(buf, binary.BigEndian, int32(-1))
	}
	cpIndex := cp.addShapeCPEntry(t)
	return binary.Write(buf, binary.BigEndian, int32(cpIndex))
}

func writeTypeToBuffer(buf *bytes.Buffer, cp *ConstantPool, t model.ValueType) error {
	// Simplified type writing - write a minimal valid type structure
	// Full type serialization would require implementing BIRTypeWriter
	// For now, we write type_tag_nil (10) as a placeholder since:
	// 1. It's in the constant_value_info switch cases (nil_constant_info has size 0)
	// 2. It has no type_structure (nothing to write after flags)
	// 3. It's a valid, simple type that won't cause parsing issues
	
	// Write tag - use type_tag_nil (10) as a safe placeholder
	// Note: This is not correct but allows the binary to be parseable
	// TODO: Implement proper type tag extraction from ValueType
	tag := byte(10) // type_tag_nil
	if err := binary.Write(buf, binary.BigEndian, tag); err != nil {
		return err
	}

	// Write name CP index (ValueType doesn't have GetName, use empty for now)
	name := ""
	nameCPIndex := cp.addStringCPEntry(name)
	if err := binary.Write(buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	// Write flags (ValueType doesn't have GetFlags, use 0)
	flags := int64(0)
	if err := binary.Write(buf, binary.BigEndian, flags); err != nil {
		return err
	}

	// type_tag_nil (10) has no type_structure in the switch cases,
	// so we don't need to write any additional structure data
	return nil
}

// BIRTypeWriter handles type serialization (simplified version)
type BIRTypeWriter struct {
	buf *bytes.Buffer
	cp  *ConstantPool
}

func NewBIRTypeWriter(buf *bytes.Buffer, cp *ConstantPool) *BIRTypeWriter {
	return &BIRTypeWriter{
		buf: buf,
		cp:  cp,
	}
}

// BIRInstructionWriter handles instruction and basic block serialization
type BIRInstructionWriter struct {
	buf        *bytes.Buffer
	scopeBuf   *bytes.Buffer
	cp         *ConstantPool
	insOffset  int
	scopeCount int
}

func NewBIRInstructionWriter(buf *bytes.Buffer, scopeBuf *bytes.Buffer, cp *ConstantPool) *BIRInstructionWriter {
	return &BIRInstructionWriter{
		buf:       buf,
		scopeBuf:  scopeBuf,
		cp:        cp,
		insOffset: 0,
	}
}

func (w *BIRInstructionWriter) GetScopeCount() int {
	return w.scopeCount
}

func (w *BIRInstructionWriter) writeBBs(basicBlocks []BIRBasicBlock) error {
	if err := binary.Write(w.buf, binary.BigEndian, int32(len(basicBlocks))); err != nil {
		return err
	}

	for _, bb := range basicBlocks {
		if err := w.writeBasicBlock(bb); err != nil {
			return err
		}
	}
	return nil
}

func (w *BIRInstructionWriter) writeBasicBlock(bb BIRBasicBlock) error {
	// BB name CP index
	bbName := bb.Id.Value()
	bbNameCPIndex := w.cp.addStringCPEntry(bbName)
	if err := binary.Write(w.buf, binary.BigEndian, int32(bbNameCPIndex)); err != nil {
		return err
	}

	// Number of instructions (including terminator)
	insCount := len(bb.Instructions) + 1
	if err := binary.Write(w.buf, binary.BigEndian, int32(insCount)); err != nil {
		return err
	}

	// Write non-terminator instructions
	for _, ins := range bb.Instructions {
		// Get position from instruction - all instructions embed BIRInstructionBase which has Pos
		var pos diagnostics.Location
		switch i := ins.(type) {
		case *Move:
			pos = i.Pos
		case *BinaryOp:
			pos = i.Pos
		case *UnaryOp:
			pos = i.Pos
		case *ConstantLoad:
			pos = i.Pos
		}
		if err := w.writePosition(pos); err != nil {
			return err
		}
		// TODO: Write scope information
		if err := binary.Write(w.buf, binary.BigEndian, byte(ins.GetKind())); err != nil {
			return err
		}
		if err := w.writeNonTerminator(ins); err != nil {
			return err
		}
		w.insOffset++
	}

	// Write terminator
	if bb.Terminator == nil {
		return fmt.Errorf("basic block without terminator: %s", bb.Id.Value())
	}

	// Get position from terminator - all terminators embed BIRTerminatorBase which has BIRInstructionBase with Pos
	var termPos diagnostics.Location
	switch t := bb.Terminator.(type) {
	case *Goto:
		termPos = t.Pos
	case *Return:
		termPos = t.Pos
	case *Branch:
		termPos = t.Pos
	case *Call:
		termPos = t.Pos
	}
	if err := w.writePosition(termPos); err != nil {
		return err
	}
	// TODO: Write scope information
	if err := binary.Write(w.buf, binary.BigEndian, byte(bb.Terminator.GetKind())); err != nil {
		return err
	}
	return w.writeTerminator(bb.Terminator)
}

func (w *BIRInstructionWriter) writeNonTerminator(ins BIRNonTerminator) error {
	switch i := ins.(type) {
	case *Move:
		return w.writeMove(i)
	case *BinaryOp:
		return w.writeBinaryOp(i)
	case *UnaryOp:
		return w.writeUnaryOp(i)
	case *ConstantLoad:
		return w.writeConstantLoad(i)
	default:
		return fmt.Errorf("unsupported non-terminator instruction: %T", ins)
	}
}

func (w *BIRInstructionWriter) writeTerminator(term BIRTerminator) error {
	switch t := term.(type) {
	case *Goto:
		return w.writeGoto(t)
	case *Return:
		return w.writeReturn(t)
	case *Branch:
		return w.writeBranch(t)
	case *Call:
		return w.writeCall(t)
	default:
		return fmt.Errorf("unsupported terminator instruction: %T", term)
	}
}

func (w *BIRInstructionWriter) writeMove(move *Move) error {
	if err := w.writeOperand(move.RhsOp); err != nil {
		return err
	}
	return w.writeOperand(move.LhsOp)
}

func (w *BIRInstructionWriter) writeBinaryOp(binOp *BinaryOp) error {
	if err := w.writeOperand(&binOp.RhsOp1); err != nil {
		return err
	}
	if err := w.writeOperand(&binOp.RhsOp2); err != nil {
		return err
	}
	return w.writeOperand(binOp.LhsOp)
}

func (w *BIRInstructionWriter) writeUnaryOp(unaryOp *UnaryOp) error {
	if err := w.writeOperand(unaryOp.RhsOp); err != nil {
		return err
	}
	return w.writeOperand(unaryOp.LhsOp)
}

func (w *BIRInstructionWriter) writeConstantLoad(constLoad *ConstantLoad) error {
	if err := writeType(w.buf, w.cp, constLoad.Type); err != nil {
		return err
	}
	if err := w.writeOperand(constLoad.LhsOp); err != nil {
		return err
	}

	// Write constant_value_info based on type tag
	// Since we're using type_tag_nil (10) as placeholder for all types,
	// constant_value_info should be nil_constant_info (size 0)
	// So we write nothing
	// TODO: When we implement proper type tag extraction, we should write the appropriate constant_value_info
	// based on the actual type tag (int_constant_info, string_constant_info, etc.)
	return nil
}

func (w *BIRInstructionWriter) writeGoto(gotoIns *Goto) error {
	if gotoIns.ThenBB == nil {
		return fmt.Errorf("goto instruction without target BB")
	}
	bbName := ""
	if gotoIns.ThenBB != nil {
		bbName = gotoIns.ThenBB.Id.Value()
	}
	bbNameCPIndex := w.cp.addStringCPEntry(bbName)
	return binary.Write(w.buf, binary.BigEndian, int32(bbNameCPIndex))
}

func (w *BIRInstructionWriter) writeReturn(returnIns *Return) error {
	// Return has no operands
	return nil
}

func (w *BIRInstructionWriter) writeBranch(branch *Branch) error {
	if err := w.writeOperand(branch.Op); err != nil {
		return err
	}
	if branch.TrueBB == nil {
		return fmt.Errorf("branch instruction without true BB")
	}
	trueBBName := ""
	if branch.TrueBB != nil {
		trueBBName = branch.TrueBB.Id.Value()
	}
	trueBBNameCPIndex := w.cp.addStringCPEntry(trueBBName)
	if err := binary.Write(w.buf, binary.BigEndian, int32(trueBBNameCPIndex)); err != nil {
		return err
	}
	if branch.FalseBB == nil {
		return fmt.Errorf("branch instruction without false BB")
	}
	falseBBName := ""
	if branch.FalseBB != nil {
		falseBBName = branch.FalseBB.Id.Value()
	}
	falseBBNameCPIndex := w.cp.addStringCPEntry(falseBBName)
	return binary.Write(w.buf, binary.BigEndian, int32(falseBBNameCPIndex))
}

func (w *BIRInstructionWriter) writeCall(call *Call) error {
	// Is virtual
	if err := binary.Write(w.buf, binary.BigEndian, call.IsVirtual); err != nil {
		return err
	}

	// Package index
	pkgCPIndex := w.cp.addPkgCPEntry(&call.CalleePkg)
	if err := binary.Write(w.buf, binary.BigEndian, int32(pkgCPIndex)); err != nil {
		return err
	}

	// Name CP index
	callName := call.Name.Value()
	nameCPIndex := w.cp.addStringCPEntry(callName)
	if err := binary.Write(w.buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	// Args count
	if err := binary.Write(w.buf, binary.BigEndian, int32(len(call.Args))); err != nil {
		return err
	}

	// Args
	for _, arg := range call.Args {
		if err := w.writeOperand(&arg); err != nil {
			return err
		}
	}

	// Has LHS operand
	hasLhsOp := call.LhsOp != nil
	if err := binary.Write(w.buf, binary.BigEndian, byte(boolToByte(hasLhsOp))); err != nil {
		return err
	}
	if hasLhsOp {
		if err := w.writeOperand(call.LhsOp); err != nil {
			return err
		}
	}

	// Then BB
	if call.ThenBB == nil {
		return fmt.Errorf("call instruction without then BB")
	}
	thenBBName := ""
	if call.ThenBB != nil {
		thenBBName = call.ThenBB.Id.Value()
	}
	thenBBNameCPIndex := w.cp.addStringCPEntry(thenBBName)
	return binary.Write(w.buf, binary.BigEndian, int32(thenBBNameCPIndex))
}

func (w *BIRInstructionWriter) writeOperand(op *BIROperand) error {
	if op == nil || op.VariableDcl == nil {
		return fmt.Errorf("nil operand")
	}

	// Check if ignored variable
	if op.VariableDcl.IgnoreVariable {
		if err := binary.Write(w.buf, binary.BigEndian, true); err != nil {
			return err
		}
		return writeType(w.buf, w.cp, op.VariableDcl.Type)
	}

	if err := binary.Write(w.buf, binary.BigEndian, false); err != nil {
		return err
	}

	// Kind
	if err := binary.Write(w.buf, binary.BigEndian, byte(op.VariableDcl.Kind)); err != nil {
		return err
	}

	// Scope
	if err := binary.Write(w.buf, binary.BigEndian, byte(op.VariableDcl.Scope)); err != nil {
		return err
	}

	// Name CP index
	varName := op.VariableDcl.Name.Value()
	nameCPIndex := w.cp.addStringCPEntry(varName)
	if err := binary.Write(w.buf, binary.BigEndian, int32(nameCPIndex)); err != nil {
		return err
	}

	// For GLOBAL and CONSTANT variables, write package index and type
	if op.VariableDcl.Kind == VAR_KIND_GLOBAL || op.VariableDcl.Kind == VAR_KIND_CONSTANT {
		// Get package ID from global var if available
		// TODO: Extract package ID properly
		// For now, use empty package ID
		emptyOrgName := model.Name("")
		emptyPkgName := model.Name("")
		emptyName := model.Name("")
		emptyVersion := model.Name("")
		pkgID := &model.PackageID{
			OrgName: &emptyOrgName,
			PkgName: &emptyPkgName,
			Name:    &emptyName,
			Version: &emptyVersion,
		}
		pkgCPIndex := w.cp.addPkgCPEntry(pkgID)
		if err := binary.Write(w.buf, binary.BigEndian, int32(pkgCPIndex)); err != nil {
			return err
		}
		return writeType(w.buf, w.cp, op.VariableDcl.Type)
	}

	return nil
}

func (w *BIRInstructionWriter) writePosition(pos diagnostics.Location) error {
	sLine := int32(-2147483648) // Integer.MIN_VALUE
	eLine := int32(-2147483648)
	sCol := int32(-2147483648)
	eCol := int32(-2147483648)
	sourceFileName := ""

	if pos != nil {
		// TODO: Extract line/column information from Location
	}

	sourceFileNameCPIndex := w.cp.addStringCPEntry(sourceFileName)
	if err := binary.Write(w.buf, binary.BigEndian, int32(sourceFileNameCPIndex)); err != nil {
		return err
	}
	if err := binary.Write(w.buf, binary.BigEndian, sLine); err != nil {
		return err
	}
	if err := binary.Write(w.buf, binary.BigEndian, sCol); err != nil {
		return err
	}
	if err := binary.Write(w.buf, binary.BigEndian, eLine); err != nil {
		return err
	}
	return binary.Write(w.buf, binary.BigEndian, eCol)
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

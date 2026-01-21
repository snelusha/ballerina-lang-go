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

// WritePosition writes a position (Location) to the buffer
func WritePosition(pos diagnostics.Location, buf *bytes.Buffer, cp *ConstantPool) error {
	var sLine, eLine, sCol, eCol int32 = -2147483648, -2147483648, -2147483648, -2147483648 // Integer.MIN_VALUE
	sourceFileName := ""
	
	if pos != nil {
		// Extract position information from Location
		// Note: This is a simplified implementation - actual Location interface may differ
		// The Java code uses pos.lineRange().startLine().line(), etc.
		// For now, we'll write default values - this should be adjusted based on actual Location interface
		if fileName := getLocationFileName(pos); fileName != nil {
			sourceFileName = *fileName
		}
	}
	
	// Write source file name CP index
	fileNameCPIndex := AddStringCPEntry(sourceFileName, cp)
	if err := binary.Write(buf, binary.BigEndian, int32(fileNameCPIndex)); err != nil {
		return err
	}
	
	// Write line and column information
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

// getLocationFileName extracts the file name from a Location
// This is a helper function that should be adjusted based on the actual Location interface
func getLocationFileName(pos diagnostics.Location) *string {
	// TODO: Implement based on actual Location interface
	// For now, return nil
	return nil
}

// AddStringCPEntry adds a string to the constant pool and returns its index
func AddStringCPEntry(value string, cp *ConstantPool) int {
	return cp.AddCPEntry(&StringCPEntry{Value: value})
}

// AddPkgCPEntry adds a package ID to the constant pool and returns its index
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

// Helper functions to extract package ID components
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

// WriteType writes a type to the buffer as a CP index
// Accepts both TypeNode (for AST types) and ValueType (for runtime types)
// If type is nil, writes -1 as the CP index (indicating no type)
func WriteType(cp *ConstantPool, buf *bytes.Buffer, t interface{}) error {
	if t == nil {
		// Write -1 to indicate nil/no type
		return binary.Write(buf, binary.BigEndian, int32(-1))
	}
	
	// Handle TypeNode (AST types)
	if typeNode, ok := t.(model.TypeNode); ok {
		cpIndex := cp.AddShapeCPEntry(typeNode)
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	}
	
	// Handle ValueType (runtime types - this is what BIR uses)
	if valueType, ok := t.(model.ValueType); ok {
		cpIndex := cp.AddShapeCPEntryForType(valueType)
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	}
	
	return fmt.Errorf("unsupported type: %T (expected TypeNode or ValueType)", t)
}

// WriteConstValue writes a constant value to the buffer
func WriteConstValue(cp *ConstantPool, buf *bytes.Buffer, constValue ConstValue) error {
	// Handle case where type is nil
	if constValue.Type == nil {
		// If type is nil but value exists, try to infer type from value
		// Otherwise, write nothing
		if constValue.Value == nil {
			return nil // Nothing to write
		}
		// Try to infer type from value
		return WriteConstValueWithType(cp, buf, constValue.Value, nil)
	}
	return WriteConstValueWithType(cp, buf, constValue.Value, constValue.Type)
}

// WriteConstValueWithType writes a constant value with its type to the buffer
func WriteConstValueWithType(cp *ConstantPool, buf *bytes.Buffer, value interface{}, t model.ValueType) error {
	// Handle nil type - write nothing or handle gracefully
	if t == nil {
		// For nil type, check if value is also nil
		if value == nil {
			return nil // Nothing to write
		}
		// If value exists but type is nil, try to infer from value type
		// This is a fallback for cases where type information is missing
		switch v := value.(type) {
		case string:
			// Treat as string
			cpIndex := AddStringCPEntry(v, cp)
			return binary.Write(buf, binary.BigEndian, int32(cpIndex))
		case bool:
			return binary.Write(buf, binary.BigEndian, v)
		case int64, int, int32, int16, int8:
			// Try to convert to int64
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
			cpIndex := cp.AddCPEntry(&IntegerCPEntry{Value: intVal})
			return binary.Write(buf, binary.BigEndian, int32(cpIndex))
		case float64, float32:
			var floatVal float64
			if f64, ok := value.(float64); ok {
				floatVal = f64
			} else {
				floatVal = float64(value.(float32))
			}
			cpIndex := cp.AddCPEntry(&FloatCPEntry{Value: floatVal})
			return binary.Write(buf, binary.BigEndian, int32(cpIndex))
		default:
			return fmt.Errorf("type is nil but value is not: %T", value)
		}
	}
	
	// Get the implied type tag
	typeTag := getTypeTag(t)
	
	// Debug: log the type tag if it's unexpected
	if typeTag == 0 {
		return fmt.Errorf("could not determine type tag for type: %v (TypeKind: %v)", t, t.GetTypeKind())
	}
	
	switch typeTag {
	case int(model.TypeTags_INT), int(model.TypeTags_SIGNED32_INT), int(model.TypeTags_SIGNED16_INT), 
		 int(model.TypeTags_SIGNED8_INT), int(model.TypeTags_UNSIGNED32_INT), 
		 int(model.TypeTags_UNSIGNED16_INT), int(model.TypeTags_UNSIGNED8_INT): // INT types
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
		cpIndex := cp.AddCPEntry(&IntegerCPEntry{Value: intVal})
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	case int(model.TypeTags_BYTE): // BYTE
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
		cpIndex := cp.AddCPEntry(&ByteCPEntry{Value: byteVal})
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	case int(model.TypeTags_FLOAT): // FLOAT
		var floatVal float64
		switch v := value.(type) {
		case float64:
			floatVal = v
		case float32:
			floatVal = float64(v)
		case string:
			// Parse string to float
			var err error
			floatVal, err = parseFloat64(v)
			if err != nil {
				return fmt.Errorf("cannot parse float from string: %w", err)
			}
		default:
			return fmt.Errorf("cannot convert value to float64: %T", value)
		}
		cpIndex := cp.AddCPEntry(&FloatCPEntry{Value: floatVal})
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	case int(model.TypeTags_STRING), int(model.TypeTags_CHAR_STRING), int(model.TypeTags_DECIMAL): // STRING, CHAR_STRING, DECIMAL
		strVal := fmt.Sprintf("%v", value)
		cpIndex := AddStringCPEntry(strVal, cp)
		return binary.Write(buf, binary.BigEndian, int32(cpIndex))
	case int(model.TypeTags_BOOLEAN): // BOOLEAN
		boolVal, ok := value.(bool)
		if !ok {
			return fmt.Errorf("cannot convert value to bool: %T", value)
		}
		return binary.Write(buf, binary.BigEndian, boolVal)
	case int(model.TypeTags_NIL): // NIL
		// Nothing to write for nil
		return nil
	case int(model.TypeTags_RECORD): // RECORD
		// Map of string to ConstValue
		mapVal, ok := value.(map[string]ConstValue)
		if !ok {
			return fmt.Errorf("record value must be map[string]ConstValue, got %T", value)
		}
		if err := binary.Write(buf, binary.BigEndian, int32(len(mapVal))); err != nil {
			return err
		}
		for key, fieldValue := range mapVal {
			keyCPIndex := AddStringCPEntry(key, cp)
			if err := binary.Write(buf, binary.BigEndian, int32(keyCPIndex)); err != nil {
				return err
			}
			if err := WriteType(cp, buf, fieldValue.Type); err != nil {
				return err
			}
			if err := WriteConstValue(cp, buf, fieldValue); err != nil {
				return err
			}
		}
		return nil
	case int(model.TypeTags_TUPLE): // TUPLE
		// Array of ConstValue
		tupleVal, ok := value.([]ConstValue)
		if !ok {
			return fmt.Errorf("tuple value must be []ConstValue, got %T", value)
		}
		if err := binary.Write(buf, binary.BigEndian, int32(len(tupleVal))); err != nil {
			return err
		}
		for _, memValue := range tupleVal {
			if err := WriteType(cp, buf, memValue.Type); err != nil {
				return err
			}
			if err := WriteConstValue(cp, buf, memValue); err != nil {
				return err
			}
		}
		return nil
	case int(model.TypeTags_INTERSECTION): // INTERSECTION
		// Get effective type and write recursively
		// This requires access to the type's effective type
		// For now, we'll need to handle this based on the actual type implementation
		return fmt.Errorf("intersection type const value not yet implemented")
	case int(model.TypeTags_EMPTY), int(model.TypeTags_NONE), int(model.TypeTags_VOID):
		// Empty/None/Void types - no constant value
		return nil
	case int(model.TypeTags_INVOKABLE), int(model.TypeTags_FUNCTION_POINTER):
		// Function types as constants - not typically used, but handle gracefully
		return fmt.Errorf("function type const value not supported")
	case int(model.TypeTags_ANY), int(model.TypeTags_ANYDATA), int(model.TypeTags_JSON):
		// These types don't typically have constant values
		return fmt.Errorf("constant value not supported for type tag: %d", typeTag)
	case int(model.TypeTags_XML), int(model.TypeTags_TABLE), int(model.TypeTags_STREAM):
		// These types don't typically have constant values
		return fmt.Errorf("constant value not supported for type tag: %d", typeTag)
	case int(model.TypeTags_TYPEDESC), int(model.TypeTags_TYPEREFDESC):
		// Typedesc constants - not typically used
		return fmt.Errorf("typedesc const value not supported")
	case int(model.TypeTags_ARRAY), int(model.TypeTags_UNION):
		// These might have constant values in some cases
		return fmt.Errorf("constant value not yet implemented for type tag: %d", typeTag)
	case int(model.TypeTags_OBJECT), int(model.TypeTags_ERROR):
		// These might have constant values
		return fmt.Errorf("constant value not yet implemented for type tag: %d", typeTag)
	default:
		// For unknown types, try to handle as string if it's a string-like value
		if strVal, ok := value.(string); ok {
			cpIndex := AddStringCPEntry(strVal, cp)
			return binary.Write(buf, binary.BigEndian, int32(cpIndex))
		}
		return fmt.Errorf("unsupported constant type tag: %d (value type: %T)", typeTag, value)
	}
}

// getTypeTagFromValueType is an alias for getTypeTag for use in instruction_writer
func getTypeTagFromValueType(t model.ValueType) int {
	return getTypeTag(t)
}

// getTypeTag extracts the type tag from a ValueType
// Maps TypeKind to TypeTags
func getTypeTag(t model.ValueType) int {
	if t == nil {
		return 0
	}
	
	// Get TypeKind from the type
	typeKind := t.GetTypeKind()
	
	// Map TypeKind to TypeTags
	// This is a simplified mapping - adjust based on actual mappings
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
		// Default to EMPTY or return 0 for unknown types
		return int(model.TypeTags_EMPTY)
	}
}

// parseFloat64 parses a string to float64
func parseFloat64(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// WriteAnnotAttachments writes annotation attachments to the buffer
// If annotAttachments is nil, writes length 0 (empty annotations)
func WriteAnnotAttachments(cp *ConstantPool, buf *bytes.Buffer, annotAttachments []interface{}) error {
	// Create a temporary buffer for annotations
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
			if err := WriteAnnotAttachment(cp, annotBuf, annotAttachment); err != nil {
				return err
			}
		}
	}
	
	// Write length and then the annotation buffer
	length := int64(annotBuf.Len())
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(annotBuf.Bytes())
	return err
}

// WriteAnnotAttachment writes a single annotation attachment
// This is a placeholder - actual implementation depends on annotation attachment structure
func WriteAnnotAttachment(cp *ConstantPool, buf *bytes.Buffer, annotAttachment interface{}) error {
	// TODO: Implement based on actual annotation attachment structure
	// This should match the Java implementation in BIRWriterUtils.writeAnnotAttachment
	return fmt.Errorf("annotation attachment writing not yet implemented")
}

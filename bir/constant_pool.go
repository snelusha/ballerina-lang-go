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
	"hash/fnv"

	"ballerina-lang-go/model"
)

// CPEntryType represents the type of a constant pool entry
type CPEntryType byte

const (
	CPEntryInteger CPEntryType = iota + 1
	CPEntryFloat
	CPEntryBoolean
	CPEntryString
	CPEntryPackage
	CPEntryByte
	CPEntryShape
)

// CPEntry represents a constant pool entry
type CPEntry interface {
	EntryType() CPEntryType
	hashKey() string
}

// IntegerCPEntry represents an integer constant pool entry
type IntegerCPEntry struct {
	Value int64
}

func (e *IntegerCPEntry) EntryType() CPEntryType { return CPEntryInteger }
func (e *IntegerCPEntry) hashKey() string        { return fmt.Sprintf("int:%d", e.Value) }

// FloatCPEntry represents a float constant pool entry
type FloatCPEntry struct {
	Value float64
}

func (e *FloatCPEntry) EntryType() CPEntryType { return CPEntryFloat }
func (e *FloatCPEntry) hashKey() string        { return fmt.Sprintf("float:%g", e.Value) }

// BooleanCPEntry represents a boolean constant pool entry
type BooleanCPEntry struct {
	Value bool
}

func (e *BooleanCPEntry) EntryType() CPEntryType { return CPEntryBoolean }
func (e *BooleanCPEntry) hashKey() string        { return fmt.Sprintf("bool:%v", e.Value) }

// StringCPEntry represents a string constant pool entry
type StringCPEntry struct {
	Value string
}

func (e *StringCPEntry) EntryType() CPEntryType { return CPEntryString }
func (e *StringCPEntry) hashKey() string        { return fmt.Sprintf("str:%s", e.Value) }

// ByteCPEntry represents a byte constant pool entry
type ByteCPEntry struct {
	Value int32
}

func (e *ByteCPEntry) EntryType() CPEntryType { return CPEntryByte }
func (e *ByteCPEntry) hashKey() string        { return fmt.Sprintf("byte:%d", e.Value) }

// PackageCPEntry represents a package constant pool entry
type PackageCPEntry struct {
	OrgNameCPIndex    int
	PkgNameCPIndex    int
	ModuleNameCPIndex int
	VersionCPIndex    int
}

func (e *PackageCPEntry) EntryType() CPEntryType { return CPEntryPackage }
func (e *PackageCPEntry) hashKey() string {
	return fmt.Sprintf("pkg:%d:%d:%d:%d", e.OrgNameCPIndex, e.PkgNameCPIndex, e.ModuleNameCPIndex, e.VersionCPIndex)
}

// ShapeCPEntry represents a type shape constant pool entry (for TypeNode)
type ShapeCPEntry struct {
	Shape model.TypeNode
}

func (e *ShapeCPEntry) EntryType() CPEntryType { return CPEntryShape }
func (e *ShapeCPEntry) hashKey() string {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%p", e.Shape)))
	return fmt.Sprintf("shape:%d", h.Sum32())
}

// ShapeCPEntryForType represents a type shape constant pool entry (for ValueType/Type)
type ShapeCPEntryForType struct {
	Type model.ValueType
}

func (e *ShapeCPEntryForType) EntryType() CPEntryType { return CPEntryShape }
func (e *ShapeCPEntryForType) hashKey() string {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%p", e.Type)))
	return fmt.Sprintf("shapetype:%d", h.Sum32())
}

type ConstantPool struct {
	entries  []CPEntry
	entryMap map[string]int
	typeEnv  any
}

func (cp *ConstantPool) PrintEntries() {
	for i, entry := range cp.entries {
		fmt.Printf("Entry %d: Type %d, Key %+v\n", i, entry.EntryType(), entry.hashKey())
	}

	// for key, idx := range cp.entryMap {
	// 	fmt.Printf("Map Key: %s, Index: %d\n", key, idx)
	// }
}

func NewConstantPool(typeEnv any) *ConstantPool {
	return &ConstantPool{
		entries:  make([]CPEntry, 0),
		entryMap: make(map[string]int),
		typeEnv:  typeEnv,
	}
}

func (cp *ConstantPool) AddCPEntry(entry CPEntry) int {
	key := entry.hashKey()
	if idx, exists := cp.entryMap[key]; exists {
		return idx
	}
	idx := len(cp.entries)
	cp.entries = append(cp.entries, entry)
	cp.entryMap[key] = idx
	return idx
}

func (cp *ConstantPool) AddShapeCPEntry(shape model.TypeNode) int {
	return cp.AddCPEntry(&ShapeCPEntry{Shape: shape})
}

func (cp *ConstantPool) AddShapeCPEntryForType(shape model.ValueType) int {
	return cp.AddCPEntry(&ShapeCPEntryForType{Type: shape})
}

func (cp *ConstantPool) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	countPos := buf.Len()
	if err := binary.Write(&buf, binary.BigEndian, int32(-1)); err != nil {
		return nil, err
	}

	for _, entry := range cp.entries {
		if err := cp.writeEntry(&buf, entry); err != nil {
			return nil, fmt.Errorf("writing CP entry: %w", err)
		}
	}

	data := buf.Bytes()
	count := int32(len(cp.entries))
	binary.BigEndian.PutUint32(data[countPos:countPos+4], uint32(count))

	return data, nil
}

func (cp *ConstantPool) writeEntry(buf *bytes.Buffer, entry CPEntry) error {
	if err := binary.Write(buf, binary.BigEndian, byte(entry.EntryType())); err != nil {
		return err
	}

	switch e := entry.(type) {
	case *IntegerCPEntry:
		return binary.Write(buf, binary.BigEndian, e.Value)
	case *FloatCPEntry:
		return binary.Write(buf, binary.BigEndian, e.Value)
	case *BooleanCPEntry:
		var val byte
		if e.Value {
			val = 1
		}
		return binary.Write(buf, binary.BigEndian, val)
	case *StringCPEntry:
		if e.Value == "" {
			return binary.Write(buf, binary.BigEndian, int16(-1))
		}
		strBytes := []byte(e.Value)
		if err := binary.Write(buf, binary.BigEndian, int32(len(strBytes))); err != nil {
			return err
		}
		_, err := buf.Write(strBytes)
		return err
	case *ByteCPEntry:
		return binary.Write(buf, binary.BigEndian, e.Value)
	case *PackageCPEntry:
		if err := binary.Write(buf, binary.BigEndian, int32(e.OrgNameCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, int32(e.PkgNameCPIndex)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, int32(e.ModuleNameCPIndex)); err != nil {
			return err
		}
		return binary.Write(buf, binary.BigEndian, int32(e.VersionCPIndex))
	case *ShapeCPEntry:
		typeBuf := &bytes.Buffer{}
		if err := writeTypeToBuffer(typeBuf, cp, cp.typeEnv, e.Shape); err != nil {
			return fmt.Errorf("writing shape type: %w", err)
		}
		shapeBytes := typeBuf.Bytes()
		if err := binary.Write(buf, binary.BigEndian, int32(len(shapeBytes))); err != nil {
			return err
		}
		_, err := buf.Write(shapeBytes)
		return err
	case *ShapeCPEntryForType:
		typeBuf := &bytes.Buffer{}
		if err := writeTypeToBuffer(typeBuf, cp, cp.typeEnv, e.Type); err != nil {
			return fmt.Errorf("writing shape ValueType: %w", err)
		}
		shapeBytes := typeBuf.Bytes()
		if err := binary.Write(buf, binary.BigEndian, int32(len(shapeBytes))); err != nil {
			return err
		}
		_, err := buf.Write(shapeBytes)
		return err
	default:
		return fmt.Errorf("unknown CP entry type: %T", entry)
	}
}

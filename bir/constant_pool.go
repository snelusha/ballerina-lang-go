package bir

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"ballerina-lang-go/model"
)

type CPEntryType uint8

const (
	CP_ENTRY_INTEGER CPEntryType = iota + 1
	CP_ENTRY_FLOAT
	CP_ENTRY_BOOLEAN
	CP_ENTRY_STRING
	CP_ENTRY_PACKAGE
	CP_ENTRY_BYTE
	CP_ENTRY_SHAPE
)

type CPEntry interface {
	EntryType() CPEntryType
}

type (
	IntegerCPEntry struct {
		Value int64
	}
	FloatCPEntry struct {
		Value float32
	}
	BooleanCPEntry struct {
		Value bool
	}
	StringCPEntry struct {
		Value *string
	}
	ByteCPEntry struct {
		Value byte
	}
	PackageCPEntry struct {
		OrgNameCPIndex    int
		PkgNameCPIndex    int
		ModuleNameCPIndex int
		VersionCPIndex    int
	}
	ShapeCPEntry struct {
		Shape model.ValueType
	}
)

func (e *IntegerCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_INTEGER
}

func (e *FloatCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_FLOAT
}

func (e *BooleanCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_BOOLEAN
}

func (e *StringCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_STRING
}

func (e *ByteCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_BYTE
}

func (e *PackageCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_PACKAGE
}

func (e *ShapeCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_SHAPE
}

var (
	_ CPEntry = &IntegerCPEntry{}
	_ CPEntry = &FloatCPEntry{}
	_ CPEntry = &BooleanCPEntry{}
	_ CPEntry = &StringCPEntry{}
	_ CPEntry = &ByteCPEntry{}
	_ CPEntry = &PackageCPEntry{}
	_ CPEntry = &ShapeCPEntry{}
)

type ConstantPool struct {
	entries  []CPEntry
	entryMap map[string]int
}

func (cp *ConstantPool) EntryKey(entry CPEntry) string {
	switch e := entry.(type) {
	case *StringCPEntry:
		return fmt.Sprintf("str:%s", *e.Value)
	case *IntegerCPEntry:
		return fmt.Sprintf("int:%d", e.Value)
	case *FloatCPEntry:
		return fmt.Sprintf("float:%g", e.Value)
	case *BooleanCPEntry:
		return fmt.Sprintf("bool:%v", e.Value)
	case *ByteCPEntry:
		return fmt.Sprintf("byte:%d", e.Value)
	case *PackageCPEntry:
		return fmt.Sprintf("pkg:%d:%d:%d:%d", e.OrgNameCPIndex, e.PkgNameCPIndex, e.ModuleNameCPIndex, e.VersionCPIndex)
	case *ShapeCPEntry:
		// For shape entries, we use a simple approach - in practice, type equality
		// would need proper type comparison
		return fmt.Sprintf("shape:%p", e.Shape)
	default:
		return fmt.Sprintf("unknown:%p", entry)
	}
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		entries:  make([]CPEntry, 0),
		entryMap: make(map[string]int),
	}
}

func (cp *ConstantPool) AddCPEntry(entry CPEntry) int {
	key := cp.EntryKey(entry)
	if index, exists := cp.entryMap[key]; exists {
		return index
	}

	index := len(cp.entries)
	cp.entries = append(cp.entries, entry)
	cp.entryMap[key] = index
	return index
}

func (cp *ConstantPool) AddShapeCPEntry(shape model.ValueType) int {
	return cp.AddCPEntry(&ShapeCPEntry{Shape: shape})
}

func (cp *ConstantPool) WriteCPEntry(buf *bytes.Buffer, entry CPEntry) error {
	if err := binary.Write(buf, binary.BigEndian, entry.EntryType()); err != nil {
		return err
	}

	switch e := entry.(type) {
	case *IntegerCPEntry:
		if err := binary.Write(buf, binary.BigEndian, e.Value); err != nil {
			return err
		}
	case *FloatCPEntry:
		if err := binary.Write(buf, binary.BigEndian, e.Value); err != nil {
			return err
		}
	case *BooleanCPEntry:
		var boolByte byte
		if e.Value {
			boolByte = 1
		}
		if err := binary.Write(buf, binary.BigEndian, boolByte); err != nil {
			return err
		}
	case *StringCPEntry:
		if e.Value == nil {
			if err := binary.Write(buf, binary.BigEndian, -1); err != nil {
				return err
			}
		} else {
			strBytes := []byte(*e.Value)
			strLen := int32(len(strBytes))
			if err := binary.Write(buf, binary.BigEndian, strLen); err != nil {
				return err
			}
			if _, err := buf.Write(strBytes); err != nil {
				return err
			}
		}
	case *ByteCPEntry:
		if err := binary.Write(buf, binary.BigEndian, e.Value); err != nil {
			return err
		}
	case *PackageCPEntry:
		if err := binary.Write(buf, binary.BigEndian, e.OrgNameCPIndex); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, e.PkgNameCPIndex); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, e.ModuleNameCPIndex); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, e.VersionCPIndex); err != nil {
			return err
		}
	}

	return nil
}

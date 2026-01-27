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
		Value float64
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
		OrgNameCPIndex    int32
		PkgNameCPIndex    int32
		ModuleNameCPIndex int32
		VersionCPIndex    int32
	}
	ShapeCPEntry struct {
		Shape *minimalBType
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

func (e *PackageCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_PACKAGE
}

func (e *ByteCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_BYTE
}

func (e *ShapeCPEntry) EntryType() CPEntryType {
	return CP_ENTRY_SHAPE
}

var (
	_ CPEntry = &IntegerCPEntry{}
	_ CPEntry = &FloatCPEntry{}
	_ CPEntry = &BooleanCPEntry{}
	_ CPEntry = &StringCPEntry{}
	_ CPEntry = &PackageCPEntry{}
	_ CPEntry = &ByteCPEntry{}
	_ CPEntry = &ShapeCPEntry{}
)

type ConstantPool struct {
	entries  []CPEntry
	entryMap map[string]int
}

func (cp *ConstantPool) EntryKey(entry CPEntry) string {
	switch e := entry.(type) {
	case *IntegerCPEntry:
		return fmt.Sprintf("int:%d", e.Value)
	case *FloatCPEntry:
		return fmt.Sprintf("float:%g", e.Value)
	case *BooleanCPEntry:
		return fmt.Sprintf("bool:%v", e.Value)
	case *StringCPEntry:
		return fmt.Sprintf("str:%s", *e.Value)
	case *PackageCPEntry:
		return fmt.Sprintf("pkg:%d:%d:%d:%d", e.OrgNameCPIndex, e.PkgNameCPIndex, e.ModuleNameCPIndex, e.VersionCPIndex)
	case *ByteCPEntry:
		return fmt.Sprintf("byte:%d", e.Value)
	case *ShapeCPEntry:
		// TODO: Proper shape key generation
		if e.Shape == nil {
			return "shape:nil"
		}
		return fmt.Sprintf("shape:t%d:n%s:f%d", e.Shape.GetTag(), e.Shape.GetName(), e.Shape.GetFlags())
	default:
		panic("unknown CPEntry type")
	}
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		entries:  make([]CPEntry, 0),
		entryMap: make(map[string]int),
	}
}

func (cp *ConstantPool) AddEntry(entry CPEntry) int32 {
	key := cp.EntryKey(entry)
	if index, exists := cp.entryMap[key]; exists {
		return int32(index)
	}

	index := len(cp.entries)
	cp.entries = append(cp.entries, entry)
	cp.entryMap[key] = index
	return int32(index)
}

func (cp *ConstantPool) AddIntegerCPEntry(value int64) int32 {
	return cp.AddEntry(&IntegerCPEntry{Value: value})
}

func (cp *ConstantPool) AddFloatCPEntry(value float64) int32 {
	return cp.AddEntry(&FloatCPEntry{Value: value})
}

func (cp *ConstantPool) AddBooleanCPEntry(value bool) int32 {
	return cp.AddEntry(&BooleanCPEntry{Value: value})
}

func (cp *ConstantPool) AddStringCPEntry(value *string) int32 {
	return cp.AddEntry(&StringCPEntry{Value: value})
}

func (cp *ConstantPool) AddPackageCPEntry(pkg *model.PackageID) int32 {
	// TODO: Handle empty pkg

	OrgName := pkg.OrgName.Value()
	PkgName := pkg.PkgName.Value()
	ModuleName := pkg.Name.Value()
	Version := pkg.Version.Value()

	return cp.AddEntry(&PackageCPEntry{
		OrgNameCPIndex:    cp.AddStringCPEntry(&OrgName),
		PkgNameCPIndex:    cp.AddStringCPEntry(&PkgName),
		ModuleNameCPIndex: cp.AddStringCPEntry(&ModuleName),
		VersionCPIndex:    cp.AddStringCPEntry(&Version),
	})
}

func (cp *ConstantPool) AddByteCPEntry(value byte) int32 {
	return cp.AddEntry(&ByteCPEntry{Value: value})
}

func (cp *ConstantPool) AddShapeCPEntry(shape *minimalBType) int32 {
	// Pre-add the name to the constant pool before adding the shape
	// This ensures the CP doesn't grow during serialization
	name := shape.name.Value()
	cp.AddStringCPEntry(&name)
	return cp.AddEntry(&ShapeCPEntry{Shape: shape})
}

func (cp *ConstantPool) WriteCPEntry(buf *bytes.Buffer, entry CPEntry) error {
	entryType := entry.EntryType()
	if err := binary.Write(buf, binary.BigEndian, int8(entryType)); err != nil {
		return err
	}

	switch e := entry.(type) {
	case *IntegerCPEntry:
		return binary.Write(buf, binary.BigEndian, e.Value)
	case *FloatCPEntry:
		return binary.Write(buf, binary.BigEndian, e.Value)
	case *BooleanCPEntry:
		var b byte
		if e.Value {
			b = 1
		}
		return binary.Write(buf, binary.BigEndian, uint8(b))
	case *StringCPEntry:
		if e.Value == nil {
			return binary.Write(buf, binary.BigEndian, int32(-1))
		}
		strBytes := []byte(*e.Value)
		if err := binary.Write(buf, binary.BigEndian, int32(len(strBytes))); err != nil {
			return err
		}
		_, err := buf.Write(strBytes)
		return err
	case *ByteCPEntry:
		return binary.Write(buf, binary.BigEndian, int32(e.Value))
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

		// tag
		if err := binary.Write(typeBuf, binary.BigEndian, uint8(e.Shape.tag)); err != nil {
			return err
		}

		// name
		name := e.Shape.name.Value()
		if err := binary.Write(typeBuf, binary.BigEndian, cp.AddStringCPEntry(&name)); err != nil {
			return err
		}

		// flags
		if err := binary.Write(typeBuf, binary.BigEndian, e.Shape.flags); err != nil {
			return err
		}

		if err := binary.Write(buf, binary.BigEndian, int32(typeBuf.Len())); err != nil {
			return err
		}
		_, err := buf.Write(typeBuf.Bytes())
		return err
	default:
		return fmt.Errorf("unsupported constant pool entry type: %T", entry)
	}
}

func (cp *ConstantPool) Serialize() ([]byte, error) {
	buf := &bytes.Buffer{}

	fmt.Printf("[CP] Serializing %d entries\n", len(cp.entries))
	for i, entry := range cp.entries {
		fmt.Printf("[CP] Entry %d: %T\n", i, entry)
	}

	if err := binary.Write(buf, binary.BigEndian, int32(-1)); err != nil {
		return nil, err
	}

	for _, entry := range cp.entries {
		if err := cp.WriteCPEntry(buf, entry); err != nil {
			return nil, err
		}
	}

	bytes := buf.Bytes()
	entryCount := int32(len(cp.entries))
	binary.BigEndian.PutUint32(bytes[0:4], uint32(entryCount))

	fmt.Printf("[CP] Total serialized size: %d bytes\n", len(bytes))
	return bytes, nil
}

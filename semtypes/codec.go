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

package semtypes

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"
)

// Binary codec version for semtype section. Bump when format changes.
const SemTypeCodecVersion = 1

// SemTypeSectionMagic is the 4-byte magic at the start of the BIR semtype section.
const SemTypeSectionMagic = "SEMT"

// SemType wire tags
const (
	semTypeTagBasic   = 0
	semTypeTagComplex = 1
)

// SubtypeData wire tags
const (
	subtypeTagAll     = 0
	subtypeTagNothing = 1
	subtypeTagBoolean = 2
	subtypeTagInt     = 3
	// subtypeTagBdd and more: reserved for future (list/mapping/function/object/table/cell/error/xml)
)

// refSemType is a forward reference to a semtype in the decoded table.
// It implements SemType by delegating to the resolved type after decode.
type refSemType struct {
	index int
	table *[]SemType
}

func (r *refSemType) get() SemType {
	tbl := *r.table
	if r.index < 0 || r.index >= len(tbl) {
		panic(fmt.Sprintf("semtype ref index %d out of range (len=%d)", r.index, len(tbl)))
	}
	return tbl[r.index]
}

// SemTypeRef is used during decode for forward refs; it implements ComplexSemType
// by delegating to the referenced type. Type-assert to *refSemType to detect refs.
var _ ComplexSemType = (*refSemType)(nil)

func (r *refSemType) All() int {
	return r.get().All()
}

func (r *refSemType) Some() int {
	return r.get().(ComplexSemType).Some()
}

func (r *refSemType) SubtypeDataList() []ProperSubtypeData {
	return r.get().(ComplexSemType).SubtypeDataList()
}

func (r *refSemType) String() string {
	return r.get().String()
}

// WriteSemType encodes one SemType to w. lookupIndex must return the index for
// any SemType that appears as a nested reference (e.g. in list rest type).
func WriteSemType(w io.Writer, lookupIndex func(SemType) int32, t SemType) error {
	if t == nil {
		return fmt.Errorf("cannot write nil SemType")
	}
	if b, ok := t.(*BasicTypeBitSet); ok {
		if err := writeUint8(w, semTypeTagBasic); err != nil {
			return err
		}
		return writeInt32(w, int32(b.bitset))
	}
	ct, ok := t.(ComplexSemType)
	if !ok {
		return fmt.Errorf("unsupported SemType %T", t)
	}

	// If this complex semtype uses subtype data that the binary codec does not
	// yet understand (e.g. FloatSubtype, StringSubtype, Bdd/list/mapping/etc.),
	// widen it to a BasicTypeBitSet so we at least preserve the set of basic
	// tags (INT|FLOAT|LIST|MAPPING|...) without failing serialization.
	if !isCodecSupportedComplex(ct) {
		widened := WidenToBasicTypes(t)
		if err := writeUint8(w, semTypeTagBasic); err != nil {
			return err
		}
		return writeInt32(w, int32(widened.bitset))
	}

	if err := writeUint8(w, semTypeTagComplex); err != nil {
		return err
	}
	if err := writeInt32(w, int32(ct.All())); err != nil {
		return err
	}
	if err := writeInt32(w, int32(ct.Some())); err != nil {
		return err
	}
	list := ct.SubtypeDataList()
	if err := writeInt32(w, int32(len(list))); err != nil {
		return err
	}
	some := ct.Some()
	for _, data := range list {
		code := 0
		for (some & (1 << code)) == 0 {
			code++
		}
		some &^= 1 << code
		if err := writeUint8(w, uint8(code)); err != nil {
			return err
		}
		if err := writeSubtypeData(w, lookupIndex, BasicTypeCodeFrom(code), data); err != nil {
			return err
		}
	}
	return nil
}

// isCodecSupportedComplex reports whether the codec knows how to encode all
// subtype data for this complex semtype without widening. Currently we only
// handle:
//   - AllOrNothingSubtype for any basic type
//   - IntSubtype for BT_INT
//   - BooleanSubtype for BT_BOOLEAN
// Everything else is widened to a BasicTypeBitSet via WidenToBasicTypes.
func isCodecSupportedComplex(ct ComplexSemType) bool {
	some := ct.Some()
	dataList := ct.SubtypeDataList()
	if some == 0 || len(dataList) == 0 {
		// No "some" subtypes; representable using just the all-bitset.
		return true
	}

	idx := 0
	for some != 0 && idx < len(dataList) {
		// Find next basic type code present in the "some" bitset.
		tz := bits.TrailingZeros(uint(some))
		if tz < 0 {
			break
		}
		code := BasicTypeCodeFrom(tz)
		data := dataList[idx]
		idx++
		some &^= 1 << tz

		// If this is an all/nothing subtype, we can always encode it.
		if _, ok := data.(AllOrNothingSubtype); ok {
			continue
		}
		if _, ok := data.(*AllOrNothingSubtype); ok {
			continue
		}

		switch code {
		case BT_INT:
			if _, ok := data.(IntSubtype); ok {
				continue
			}
			if _, ok := data.(*IntSubtype); ok {
				continue
			}
			return false
		case BT_BOOLEAN:
			if _, ok := data.(BooleanSubtype); ok {
				continue
			}
			if _, ok := data.(*BooleanSubtype); ok {
				continue
			}
			return false
		default:
			// Any other basic type with non-trivial subtype data (e.g. FLOAT,
			// STRING, LIST, MAPPING, TABLE, OBJECT, ...) is not yet supported.
			return false
		}
	}
	return true
}

func writeSubtypeData(w io.Writer, lookupIndex func(SemType) int32, code BasicTypeCode, data SubtypeData) error {
	switch a := data.(type) {
	case *AllOrNothingSubtype:
		if a.IsAllSubtype() {
			return writeUint8(w, subtypeTagAll)
		}
		return writeUint8(w, subtypeTagNothing)
	case AllOrNothingSubtype:
		if a.IsAllSubtype() {
			return writeUint8(w, subtypeTagAll)
		}
		return writeUint8(w, subtypeTagNothing)
	}
	if _, ok := data.(BooleanSubtype); ok {
		opt := BooleanSubtypeSingleValue(data)
		if opt.IsEmpty() {
			return fmt.Errorf("boolean subtype has no single value")
		}
		if err := writeUint8(w, subtypeTagBoolean); err != nil {
			return err
		}
		return writeBool(w, opt.Get())
	}
	if i, ok := data.(*IntSubtype); ok {
		if err := writeUint8(w, subtypeTagInt); err != nil {
			return err
		}
		return writeIntSubtype(w, *i)
	}
	if i, ok := data.(IntSubtype); ok {
		if err := writeUint8(w, subtypeTagInt); err != nil {
			return err
		}
		return writeIntSubtype(w, i)
	}
	// TODO: Bdd, list/mapping/function/object subtypes
	return fmt.Errorf("unsupported SubtypeData for code %v: %T", code, data)
}

func writeIntSubtype(w io.Writer, i IntSubtype) error {
	if err := writeInt32(w, int32(len(i.Ranges))); err != nil {
		return err
	}
	for _, r := range i.Ranges {
		if err := writeInt64(w, r.Min); err != nil {
			return err
		}
		if err := writeInt64(w, r.Max); err != nil {
			return err
		}
	}
	return nil
}

func writeUint8(w io.Writer, v uint8) error {
	_, err := w.Write([]byte{v})
	return err
}

func writeBool(w io.Writer, v bool) error {
	var b uint8
	if v {
		b = 1
	}
	return writeUint8(w, b)
}

func writeInt32(w io.Writer, v int32) error {
	return binary.Write(w, binary.BigEndian, v)
}

func writeInt64(w io.Writer, v int64) error {
	return binary.Write(w, binary.BigEndian, v)
}

// ReadSemType decodes one SemType from r. semTypes is the slice being filled
// (indices 0..current are set). When a ref to index j is read, a refSemType
// pointing into semTypes is returned so that after all decode refs resolve.
func ReadSemType(r io.Reader, env Env, semTypes *[]SemType) (SemType, error) {
	var tag uint8
	if err := readUint8(r, &tag); err != nil {
		return nil, err
	}
	switch tag {
	case semTypeTagBasic:
		var bitset int32
		if err := readInt32(r, &bitset); err != nil {
			return nil, err
		}
		b := BasicTypeBitSetFrom(int(bitset))
		return &b, nil
	case semTypeTagComplex:
		var all, some int32
		if err := readInt32(r, &all); err != nil {
			return nil, err
		}
		if err := readInt32(r, &some); err != nil {
			return nil, err
		}
		var n int32
		if err := readInt32(r, &n); err != nil {
			return nil, err
		}
		var subtypeList []BasicSubtype
		for i := int32(0); i < n; i++ {
			var code uint8
			if err := readUint8(r, &code); err != nil {
				return nil, err
			}
			data, err := readSubtypeData(r, env, semTypes, BasicTypeCodeFrom(int(code)))
			if err != nil {
				return nil, err
			}
			subtypeList = append(subtypeList, BasicSubtypeFrom(BasicTypeCodeFrom(int(code)), data))
		}
		return CreateComplexSemTypeWithAllBitSetSubtypeList(int(all), subtypeList), nil
	default:
		return nil, fmt.Errorf("unknown semtype tag %d", tag)
	}
}

func readSubtypeData(r io.Reader, env Env, semTypes *[]SemType, code BasicTypeCode) (ProperSubtypeData, error) {
	var tag uint8
	if err := readUint8(r, &tag); err != nil {
		return nil, err
	}
	switch tag {
	case subtypeTagAll:
		return &AllOrNothingSubtypeAll, nil
	case subtypeTagNothing:
		return &AllOrNothingSubtypeNothing, nil
	case subtypeTagBoolean:
		var b uint8
		if err := readUint8(r, &b); err != nil {
			return nil, err
		}
		return BooleanSubtypeFrom(b != 0), nil
	case subtypeTagInt:
		v, err := readIntSubtype(r)
		if err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported subtype data tag %d for code %v", tag, code)
	}
}

func readIntSubtype(r io.Reader) (IntSubtype, error) {
	var n int32
	if err := readInt32(r, &n); err != nil {
		return IntSubtype{}, err
	}
	ranges := make([]Range, 0, n)
	for i := int32(0); i < n; i++ {
		var min, max int64
		if err := readInt64(r, &min); err != nil {
			return IntSubtype{}, err
		}
		if err := readInt64(r, &max); err != nil {
			return IntSubtype{}, err
		}
		ranges = append(ranges, RangeFrom(min, max))
	}
	return CreateIntSubtype(ranges...), nil
}

func readUint8(r io.Reader, v *uint8) error {
	var buf [1]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}
	*v = buf[0]
	return nil
}

func readInt32(r io.Reader, v *int32) error {
	return binary.Read(r, binary.BigEndian, v)
}

func readInt64(r io.Reader, v *int64) error {
	return binary.Read(r, binary.BigEndian, v)
}

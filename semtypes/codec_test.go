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
	"bytes"
	"testing"
)

func TestCodecBasicTypeBitSet(t *testing.T) {
	env := CreateTypeEnv()
	ctx := TypeCheckContext(env)

	for _, ty := range []SemType{&NEVER, &NIL, &INT, &FLOAT, &STRING, &BOOLEAN, &VAL} {
		var buf bytes.Buffer
		lookup := func(SemType) int32 { return -1 }
		if err := WriteSemType(&buf, lookup, ty); err != nil {
			t.Fatalf("write %v: %v", ty, err)
		}
		semTypes := &[]SemType{}
		got, err := ReadSemType(&buf, env, semTypes)
		if err != nil {
			t.Fatalf("read %v: %v", ty, err)
		}
		if !IsSameType(ctx, ty, got) {
			t.Errorf("roundtrip %v: got %v", ty, got)
		}
	}
}

func TestCodecIntSubtype(t *testing.T) {
	env := CreateTypeEnv()
	ctx := TypeCheckContext(env)

	ty := IntConst(42)
	var buf bytes.Buffer
	lookup := func(SemType) int32 { return -1 }
	if err := WriteSemType(&buf, lookup, ty); err != nil {
		t.Fatal(err)
	}
	semTypes := &[]SemType{}
	got, err := ReadSemType(&buf, env, semTypes)
	if err != nil {
		t.Fatal(err)
	}
	if !IsSameType(ctx, ty, got) {
		t.Errorf("roundtrip IntConst(42): got %v", got)
	}
}

func TestCodecBooleanSubtype(t *testing.T) {
	env := CreateTypeEnv()
	ctx := TypeCheckContext(env)

	ty := BooleanConst(true)
	var buf bytes.Buffer
	lookup := func(SemType) int32 { return -1 }
	if err := WriteSemType(&buf, lookup, ty); err != nil {
		t.Fatal(err)
	}
	semTypes := &[]SemType{}
	got, err := ReadSemType(&buf, env, semTypes)
	if err != nil {
		t.Fatal(err)
	}
	if !IsSameType(ctx, ty, got) {
		t.Errorf("roundtrip BooleanConst(true): got %v", got)
	}
}

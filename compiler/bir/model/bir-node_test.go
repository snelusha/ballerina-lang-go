package model

import (
	"ballerina-lang-go/compiler/semantics/model/elements"
	"ballerina-lang-go/compiler/semantics/model/symbols"
	"ballerina-lang-go/compiler/util"
	"testing"
)

func TestBIRPackageCreation(t *testing.T) {
	org := util.NewName("myorg")
	pkgName := util.NewName("mypackage")
	name := util.NewName("mymodule")
	version := util.NewName("1.0.0")
	sourceFileName := util.NewName("main.bal")

	pkg := NewBIRPackage(nil, org, pkgName, name, version, sourceFileName, "/src", false)

	if pkg == nil {
		t.Fatal("Package should not be nil")
	}

	if pkg.PackageID.Org.Value != "myorg" {
		t.Errorf("Expected org 'myorg', got '%s'", pkg.PackageID.Org.Value)
	}

	if len(pkg.Functions) != 0 {
		t.Errorf("Expected 0 functions, got %d", len(pkg.Functions))
	}
}

func TestBIRVariableDclCreation(t *testing.T) {
	name := util.NewName("testVar")
	varDecl := NewBIRVariableDclMinimal(nil, name, VarScopeFunction, VarKindLocal)

	if varDecl == nil {
		t.Fatal("Variable declaration should not be nil")
	}

	if varDecl.Name.Value != "testVar" {
		t.Errorf("Expected name 'testVar', got '%s'", varDecl.Name.Value)
	}

	if varDecl.Kind != VarKindLocal {
		t.Errorf("Expected kind VarKindLocal, got %v", varDecl.Kind)
	}

	if varDecl.Scope != VarScopeFunction {
		t.Errorf("Expected scope VarScopeFunction, got %v", varDecl.Scope)
	}
}

func TestBIRFunctionCreation(t *testing.T) {
	name := util.NewName("testFunction")
	fn := NewBIRFunctionMinimal(nil, name, 0, nil, util.NewName("worker"), 0, symbols.SymbolOriginSource)

	if fn == nil {
		t.Fatal("Function should not be nil")
	}

	if fn.Name.Value != "testFunction" {
		t.Errorf("Expected name 'testFunction', got '%s'", fn.Name.Value)
	}

	if len(fn.BasicBlocks) != 0 {
		t.Errorf("Expected 0 basic blocks, got %d", len(fn.BasicBlocks))
	}
}

func TestBIRBasicBlockCreation(t *testing.T) {
	bb := NewBIRBasicBlockWithNumber(5)

	if bb == nil {
		t.Fatal("Basic block should not be nil")
	}

	if bb.Number != 5 {
		t.Errorf("Expected number 5, got %d", bb.Number)
	}

	expectedID := "bb5"
	if bb.ID.Value != expectedID {
		t.Errorf("Expected ID '%s', got '%s'", expectedID, bb.ID.Value)
	}
}

func TestVarKindValue(t *testing.T) {
	tests := []struct {
		kind     VarKind
		expected byte
	}{
		{VarKindLocal, 1},
		{VarKindArg, 2},
		{VarKindTemp, 3},
		{VarKindReturn, 4},
		{VarKindGlobal, 5},
		{VarKindSelf, 6},
		{VarKindConstant, 7},
		{VarKindSynthetic, 8},
	}

	for _, tt := range tests {
		if got := tt.kind.Value(); got != tt.expected {
			t.Errorf("VarKind.Value() = %d, want %d for %v", got, tt.expected, tt.kind)
		}
	}
}

func TestVarScopeValue(t *testing.T) {
	tests := []struct {
		scope    VarScope
		expected byte
	}{
		{VarScopeFunction, 1},
		{VarScopeGlobal, 2},
	}

	for _, tt := range tests {
		if got := tt.scope.Value(); got != tt.expected {
			t.Errorf("VarScope.Value() = %d, want %d for %v", got, tt.expected, tt.scope)
		}
	}
}

func TestBIRImportModuleEquality(t *testing.T) {
	org := util.NewName("org")
	name1 := util.NewName("module1")
	name2 := util.NewName("module2")
	version := util.NewName("1.0.0")

	mod1 := NewBIRImportModule(nil, org, name1, version)
	mod2 := NewBIRImportModule(nil, org, name1, version)
	mod3 := NewBIRImportModule(nil, org, name2, version)

	if mod1.PackageID != mod2.PackageID {
		t.Error("Expected modules with same package ID to have equal package IDs")
	}

	if mod1.PackageID == mod3.PackageID {
		t.Error("Expected modules with different names to have different package IDs")
	}
}

func TestBIRAnnotationCreation(t *testing.T) {
	name := util.NewName("TestAnnot")
	originalName := util.NewName("TestAnnot")
	attachPoints := make(map[elements.AttachPoint]struct{})

	annot := NewBIRAnnotation(nil, name, originalName, 0, attachPoints, nil, symbols.SymbolOriginSource)

	if annot == nil {
		t.Fatal("Annotation should not be nil")
	}

	if annot.Name.Value != "TestAnnot" {
		t.Errorf("Expected name 'TestAnnot', got '%s'", annot.Name.Value)
	}
}

func TestConstValueCreation(t *testing.T) {
	val := NewConstValue(42, nil)

	if val == nil {
		t.Fatal("ConstValue should not be nil")
	}

	if intVal, ok := val.Value.(int); !ok || intVal != 42 {
		t.Errorf("Expected value 42, got %v", val.Value)
	}
}

func TestChannelDetailsCreation(t *testing.T) {
	details := NewChannelDetails("testChannel", true, false)

	if details == nil {
		t.Fatal("ChannelDetails should not be nil")
	}

	if details.Name != "testChannel" {
		t.Errorf("Expected name 'testChannel', got '%s'", details.Name)
	}

	if !details.ChannelInSameStrand {
		t.Error("Expected ChannelInSameStrand to be true")
	}

	if details.Send {
		t.Error("Expected Send to be false")
	}
}

func TestBIRLockDetailsHolder(t *testing.T) {
	holder := NewBIRLockDetailsHolder()

	if holder == nil {
		t.Fatal("LockDetailsHolder should not be nil")
	}

	if !holder.IsEmpty() {
		t.Error("Expected holder to be empty initially")
	}

	if holder.Size() != 0 {
		t.Errorf("Expected size 0, got %d", holder.Size())
	}
}

func TestBirScopeCreation(t *testing.T) {
	parentScope := &BirScope{ID: 1, Parent: nil}
	childScope := &BirScope{ID: 2, Parent: parentScope}

	if childScope.Parent != parentScope {
		t.Error("Expected child scope to have correct parent")
	}

	if childScope.ID != 2 {
		t.Errorf("Expected ID 2, got %d", childScope.ID)
	}
}

func TestBIRMappingConstructorEntry(t *testing.T) {
	keyValueEntry := NewBIRMappingConstructorKeyValueEntry(nil, nil)
	if !keyValueEntry.IsKeyValuePair() {
		t.Error("KeyValueEntry should return true for IsKeyValuePair")
	}

	spreadEntry := NewBIRMappingConstructorSpreadFieldEntry(nil)
	if spreadEntry.IsKeyValuePair() {
		t.Error("SpreadFieldEntry should return false for IsKeyValuePair")
	}
}

func TestBIRFunctionDuplicate(t *testing.T) {
	name := util.NewName("original")
	original := NewBIRFunctionMinimal(nil, name, 0, nil, util.NewName("worker"), 0, symbols.SymbolOriginSource)

	duplicate := original.Duplicate()

	if duplicate == nil {
		t.Fatal("Duplicate should not be nil")
	}

	if duplicate.Name.Value != original.Name.Value {
		t.Errorf("Expected duplicate name to match original")
	}

	if duplicate.Flags != original.Flags {
		t.Errorf("Expected duplicate flags to match original")
	}
}

func TestBIRTypeDefinitionGetName(t *testing.T) {
	name := util.NewName("TestType")
	typeDef := NewBIRTypeDefinitionSimple(nil, name, name, 0, false, nil, nil, symbols.SymbolOriginSource)

	if typeDef.GetName().Value != "TestType" {
		t.Errorf("Expected name 'TestType', got '%s'", typeDef.GetName().Value)
	}
}

func TestBIRFunctionGetName(t *testing.T) {
	name := util.NewName("testFunc")
	fn := NewBIRFunctionMinimal(nil, name, 0, nil, util.NewName("worker"), 0, symbols.SymbolOriginSource)

	if fn.GetName().Value != "testFunc" {
		t.Errorf("Expected name 'testFunc', got '%s'", fn.GetName().Value)
	}
}

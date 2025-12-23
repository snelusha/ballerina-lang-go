package bir

import (
	"testing"
)

func TestNewName(t *testing.T) {
	name := NewName("testName")
	if name.Value != "testName" {
		t.Errorf("Expected 'testName', got '%s'", name.Value)
	}
}

func TestNameEquals(t *testing.T) {
	name1 := NewName("test")
	name2 := NewName("test")
	name3 := NewName("other")

	if !name1.Equals(name2) {
		t.Error("Expected names to be equal")
	}

	if name1.Equals(name3) {
		t.Error("Expected names to be different")
	}
}

func TestPackageID(t *testing.T) {
	org := NewName("ballerina")
	name := NewName("lang.int")
	version := NewName("1.0.0")

	pkgID := NewPackageID(org, name, name, version)

	if !pkgID.GetOrgName().Equals(org) {
		t.Error("Organization name mismatch")
	}

	if !pkgID.GetName().Equals(name) {
		t.Error("Package name mismatch")
	}
}

func TestPackageIDEquals(t *testing.T) {
	org := NewName("ballerina")
	name := NewName("lang.int")
	version := NewName("1.0.0")

	pkgID1 := NewPackageID(org, name, name, version)
	pkgID2 := NewPackageID(org, name, name, version)

	if !pkgID1.Equals(pkgID2) {
		t.Error("Expected package IDs to be equal")
	}
}

func TestIsLangLibPackageID(t *testing.T) {
	org := NewName("ballerina")
	langName := NewName("lang")
	intName := NewName("int")
	version := NewName("1.0.0")

	nameComps := []*Name{langName, intName}
	pkgID := NewPackageIDWithNameComps(org, nameComps, version)

	if !IsLangLibPackageID(pkgID) {
		t.Error("Expected package to be identified as lang lib")
	}
}

func TestVarKind(t *testing.T) {
	if VarKindLocal.GetValue() != 1 {
		t.Errorf("Expected VarKindLocal value to be 1, got %d", VarKindLocal.GetValue())
	}

	if VarKindGlobal.GetValue() != 5 {
		t.Errorf("Expected VarKindGlobal value to be 5, got %d", VarKindGlobal.GetValue())
	}
}

func TestVarScope(t *testing.T) {
	if VarScopeFunction.GetValue() != 1 {
		t.Errorf("Expected VarScopeFunction value to be 1, got %d", VarScopeFunction.GetValue())
	}
}

func TestSymbolOrigin(t *testing.T) {
	if SymbolOriginBuiltin.Value() != 1 {
		t.Errorf("Expected SymbolOriginBuiltin value to be 1, got %d", SymbolOriginBuiltin.Value())
	}

	origin := SymbolOriginFromValue(3)
	if origin != SymbolOriginCompiledSource {
		t.Error("Expected SymbolOriginCompiledSource")
	}
}

func TestSymbolOriginToBIROrigin(t *testing.T) {
	if SymbolOriginSource.ToBIROrigin() != SymbolOriginCompiledSource {
		t.Error("Expected SOURCE to convert to COMPILED_SOURCE")
	}
}

func TestBIRPackage(t *testing.T) {
	org := NewName("testOrg")
	pkgName := NewName("testPkg")
	name := NewName("testModule")
	version := NewName("1.0.0")
	sourceFileName := NewName("main.bal")

	pkg := NewBIRPackage(nil, org, pkgName, name, version, sourceFileName, "/test/path", false)

	if pkg.PackageID == nil {
		t.Error("Expected package ID to be set")
	}

	if !pkg.PackageID.GetOrgName().Equals(org) {
		t.Error("Organization name mismatch in package")
	}
}

func TestBIRImportModule(t *testing.T) {
	org := NewName("ballerina")
	name := NewName("io")
	version := NewName("1.0.0")

	importModule := NewBIRImportModule(nil, org, name, version)

	if importModule.PackageID == nil {
		t.Error("Expected package ID to be set")
	}

	if !importModule.PackageID.GetName().Equals(name) {
		t.Error("Import module name mismatch")
	}
}

func TestBIRVariableDcl(t *testing.T) {
	name := NewName("testVar")
	varDcl := NewBIRVariableDclBasic(nil, name, VarScopeFunction, VarKindLocal)

	if !varDcl.Name.Equals(name) {
		t.Error("Variable name mismatch")
	}

	if varDcl.Kind != VarKindLocal {
		t.Error("Variable kind mismatch")
	}
}

func TestBIRVariableDclEquals(t *testing.T) {
	name := NewName("testVar")
	varDcl1 := NewBIRVariableDclBasic(nil, name, VarScopeFunction, VarKindLocal)
	varDcl2 := NewBIRVariableDclBasic(nil, name, VarScopeFunction, VarKindLocal)

	if !varDcl1.Equals(varDcl2) {
		t.Error("Expected variable declarations to be equal")
	}
}

func TestBIRBasicBlock(t *testing.T) {
	bb := NewBIRBasicBlockSimple(0)

	if bb.Number != 0 {
		t.Errorf("Expected basic block number to be 0, got %d", bb.Number)
	}

	if bb.ID.Value != "bb0" {
		t.Errorf("Expected basic block ID to be 'bb0', got '%s'", bb.ID.Value)
	}
}

func TestBIRParameter(t *testing.T) {
	name := NewName("param1")
	param := NewBIRParameter(nil, name, 0)

	if !param.Name.Equals(name) {
		t.Error("Parameter name mismatch")
	}

	if param.Flags != 0 {
		t.Error("Parameter flags mismatch")
	}
}

func TestBIROperand(t *testing.T) {
	name := NewName("operand")
	varDcl := NewBIRVariableDclBasic(nil, name, VarScopeFunction, VarKindLocal)
	operand := NewBIROperand(varDcl)

	if operand.VariableDcl != varDcl {
		t.Error("Operand variable declaration mismatch")
	}
}

func TestBIROperandEquals(t *testing.T) {
	name := NewName("operand")
	varDcl := NewBIRVariableDclBasic(nil, name, VarScopeFunction, VarKindLocal)
	operand1 := NewBIROperand(varDcl)
	operand2 := NewBIROperand(varDcl)

	if !operand1.Equals(operand2) {
		t.Error("Expected operands to be equal")
	}
}

func TestChannelDetails(t *testing.T) {
	details := NewChannelDetails("channel1", true, false)

	if details.Name != "channel1" {
		t.Error("Channel name mismatch")
	}

	if !details.ChannelInSameStrand {
		t.Error("ChannelInSameStrand should be true")
	}

	if details.Send {
		t.Error("Send should be false")
	}
}

func TestAttachPoint(t *testing.T) {
	ap := GetAttachmentPoint("function", true)

	if ap == nil {
		t.Error("Expected attach point to be found")
	}

	if ap.Point != PointFunction {
		t.Error("Attach point type mismatch")
	}

	if !ap.Source {
		t.Error("Source should be true")
	}
}

func TestConstValue(t *testing.T) {
	cv := NewConstValue(42, nil)

	if cv.Value != 42 {
		t.Errorf("Expected value to be 42, got %v", cv.Value)
	}
}

func TestBIRLockDetailsHolder(t *testing.T) {
	holder := NewBIRLockDetailsHolder()

	if !holder.IsEmpty() {
		t.Error("Expected holder to be empty")
	}

	if holder.Size() != 0 {
		t.Errorf("Expected size to be 0, got %d", holder.Size())
	}
}

func TestBIRMappingConstructorKeyValueEntry(t *testing.T) {
	name1 := NewName("key")
	varDcl1 := NewBIRVariableDclBasic(nil, name1, VarScopeFunction, VarKindLocal)
	keyOp := NewBIROperand(varDcl1)

	name2 := NewName("value")
	varDcl2 := NewBIRVariableDclBasic(nil, name2, VarScopeFunction, VarKindLocal)
	valueOp := NewBIROperand(varDcl2)

	entry := NewBIRMappingConstructorKeyValueEntry(keyOp, valueOp)

	if !entry.IsKeyValuePair() {
		t.Error("Expected entry to be a key-value pair")
	}
}

func TestBIRMappingConstructorSpreadFieldEntry(t *testing.T) {
	name := NewName("spread")
	varDcl := NewBIRVariableDclBasic(nil, name, VarScopeFunction, VarKindLocal)
	exprOp := NewBIROperand(varDcl)

	entry := NewBIRMappingConstructorSpreadFieldEntry(exprOp)

	if entry.IsKeyValuePair() {
		t.Error("Expected entry not to be a key-value pair")
	}
}

func TestMarkdownDocAttachment(t *testing.T) {
	doc := NewMarkdownDocAttachment(2)

	if doc == nil {
		t.Error("Expected markdown doc attachment to be created")
	}

	if len(doc.Parameters) != 0 {
		t.Error("Expected parameters to be empty initially")
	}
}

func TestParameter(t *testing.T) {
	param := NewParameter("arg1", "First argument")

	if param.GetName() != "arg1" {
		t.Error("Parameter name mismatch")
	}

	if param.GetDescription() != "First argument" {
		t.Error("Parameter description mismatch")
	}
}

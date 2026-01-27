package bir

import (
	"fmt"
	"os"
	"testing"

	"ballerina-lang-go/context"
)

func TestWriteReadRoundTrip(t *testing.T) {
	testdataDir := "./testdata/bir"
	birFiles := getCorpusFiles(t, testdataDir)
	for _, birFile := range birFiles {
		t.Run(birFile, func(t *testing.T) {
			t.Parallel()
			testWriteReadRoundTrip(t, birFile)
		})
	}
}

func testWriteReadRoundTrip(t *testing.T, birFile string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while loading BIR package from %s: %v", birFile, r)
		}
	}()

	// Load BIR package
	file, err := os.Open(birFile)
	if err != nil {
		t.Fatalf("failed to open test BIR file: %v", err)
	}
	defer file.Close()
	cx := context.NewCompilerContext()
	fmt.Println("here")
	pkg, err := LoadBIRPackageFromReader(cx, file)
	if err != nil {
		t.Errorf("error loading BIR package from %s: %v", birFile, err)
		return
	}
	fmt.Println("here2")

	if pkg == nil {
		t.Errorf("BIR package is nil for %s", birFile)
		return
	}

	writer := NewBIRWriter(pkg)
	serializedData, err := writer.Serialize()
	if err != nil {
		t.Errorf("error serializing BIR package for %s: %v", birFile, err)
		return
	}

	reader := NewBIRReader(serializedData)
	loadedPkg, err := reader.LoadBIRPackage()
	if err != nil {
		t.Errorf("error re-loading serialized BIR package for %s: %v", birFile, err)
		return
	}

	if loadedPkg == nil {
		t.Errorf("re-loaded BIR package is nil for %s", birFile)
		return
	}

	if pkg.PackageID.OrgName.Value() != loadedPkg.PackageID.OrgName.Value() ||
		pkg.PackageID.PkgName.Value() != loadedPkg.PackageID.PkgName.Value() ||
		pkg.PackageID.Version.Value() != loadedPkg.PackageID.Version.Value() {
		t.Errorf("mismatch in package IDs after round-trip for %s", birFile)
	}

	if len(pkg.ImportModules) != len(loadedPkg.ImportModules) {
		t.Errorf("mismatch in number of import modules after round-trip for %s", birFile)
		return
	}

	for i, imp := range pkg.ImportModules {
		loadedImp := loadedPkg.ImportModules[i]
		if imp.PackageID.OrgName.Value() != loadedImp.PackageID.OrgName.Value() ||
			imp.PackageID.PkgName.Value() != loadedImp.PackageID.PkgName.Value() ||
			imp.PackageID.Version.Value() != loadedImp.PackageID.Version.Value() {
			t.Errorf("mismatch in import module %d after round-trip for %s", i, birFile)
		}
	}

	for i, constant := range pkg.Constants {
		loadedConstant := loadedPkg.Constants[i]
		if constant.Name.Value() != loadedConstant.Name.Value() {
			t.Errorf("mismatch in constant name %d after round-trip for %s", i, birFile)
		}
		if constant.Flags != loadedConstant.Flags {
			t.Errorf("mismatch in constant flags %d after round-trip for %s", i, birFile)
		}
		if constant.Origin != loadedConstant.Origin {
			t.Errorf("mismatch in constant origin %d after round-trip for %s", i, birFile)
		}

		if constant.Type.tag != loadedConstant.Type.tag {
			t.Errorf("mismatch in constant type tag %d after round-trip for %s", i, birFile)
		}
		if constant.Type.name.Value() != loadedConstant.Type.name.Value() {
			t.Errorf("mismatch in constant type name %d after round-trip for %s", i, birFile)
		}
		if constant.Type.flags != loadedConstant.Type.flags {
			t.Errorf("mismatch in constant type flags %d after round-trip for %s", i, birFile)
		}
		if constant.ConstValue.Type.tag != loadedConstant.ConstValue.Type.tag {
			t.Errorf("mismatch in constant value type tag %d after round-trip for %s", i, birFile)
		}
		if constant.ConstValue.Type.name.Value() != loadedConstant.ConstValue.Type.name.Value() {
			t.Errorf("mismatch in constant value type name %d after round-trip for %s", i, birFile)
		}
		if constant.ConstValue.Type.flags != loadedConstant.ConstValue.Type.flags {
			t.Errorf("mismatch in constant value type flags %d after round-trip for %s", i, birFile)
		}
		if constant.ConstValue.Value != loadedConstant.ConstValue.Value {
			t.Errorf("mismatch in constant value %d after round-trip for %s", i, birFile)
		}
	}

	for i, globalVar := range pkg.GlobalVars {
		loadedGlobalVar := loadedPkg.GlobalVars[i]
		if globalVar.Kind != loadedGlobalVar.Kind {
			t.Errorf("mismatch in global variable kind %d after round-trip for %s", i, birFile)
		}
		if globalVar.Name.Value() != loadedGlobalVar.Name.Value() {
			t.Errorf("mismatch in global variable name %d after round-trip for %s", i, birFile)
		}
		if globalVar.Flags != loadedGlobalVar.Flags {
			t.Errorf("mismatch in global variable flags %d after round-trip for %s", i, birFile)
		}
		if globalVar.Origin != loadedGlobalVar.Origin {
			t.Errorf("mismatch in global variable origin %d after round-trip for %s", i, birFile)
		}

		if globalVar.Type.tag != loadedGlobalVar.Type.tag {
			t.Errorf("mismatch in global variable type tag %d after round-trip for %s", i, birFile)
		}
		if globalVar.Type.name.Value() != loadedGlobalVar.Type.name.Value() {
			t.Errorf("mismatch in global variable type name %d after round-trip for %s", i, birFile)
		}
		if globalVar.Type.flags != loadedGlobalVar.Type.flags {
			t.Errorf("mismatch in global variable type flags %d after round-trip for %s", i, birFile)
		}
	}
}

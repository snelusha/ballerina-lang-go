package birutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	"ballerina-lang-go/context"
)

var supportedSubsets = []string{"subset1"}

func getCorpusFiles(t *testing.T, baseDir string) []string {
	// Find all .bir files
	var birFiles []string
	for _, subset := range supportedSubsets {
		dirPath := filepath.Join(baseDir, subset)
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".bir") {
				birFiles = append(birFiles, path)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Error walking corpus/bir/%s directory: %v", subset, err)
		}
	}

	if len(birFiles) == 0 {
		t.Fatalf("No .bir files found in %s", baseDir)
	}
	return birFiles
}

func TestWriteReadRoundTrip(t *testing.T) {
	testdataDir := "../bir/testdata/bir"
	// testdataDir := "./testdata/bir"
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
	pkg, err := bir.LoadBIRPackageFromReader(cx, file)
	if err != nil {
		t.Errorf("error loading BIR package from %s: %v", birFile, err)
		return
	}

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

	prettyPrinter := bir.PrettyPrinter{}

	transformFunctions := func(fns []bir.BIRFunction) []bir.BIRFunction {
		var transformedFns []bir.BIRFunction
		for _, f := range fns {
			transformedFns = append(transformedFns, bir.BIRFunction{
				Name:           f.Name,
				OriginalName:   f.OriginalName,
				Flags:          f.Flags,
				Origin:         f.Origin,
				ReturnVariable: f.ReturnVariable,
				LocalVars:      f.LocalVars,
			})
		}
		return transformedFns
	}

	mPkg := &bir.BIRPackage{
		PackageID:     pkg.PackageID,
		ImportModules: pkg.ImportModules,
		Constants:     pkg.Constants,
		GlobalVars:    pkg.GlobalVars,
		Functions:     transformFunctions(pkg.Functions),
	}

	orgPkgStr := prettyPrinter.Print(*mPkg)
	loadedOrgPkgStr := prettyPrinter.Print(*loadedPkg)

	fmt.Println(orgPkgStr)
	fmt.Println("-----")
	fmt.Println(loadedOrgPkgStr)

	if orgPkgStr != loadedOrgPkgStr {
		t.Errorf("mismatch in BIR package after round-trip for %s\nOriginal:\n%s\n\nLoaded:\n%s", birFile, orgPkgStr, loadedOrgPkgStr)
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

		ct := constant.Type.(ast.BType)
		loadedCt := loadedConstant.Type.(ast.BType)

		if ct.BTypeGetTag() != loadedCt.BTypeGetTag() {
			t.Errorf("mismatch in constant type tag %d after round-trip for %s", i, birFile)
		}
		if string(ct.BTypeGetName()) != string(loadedCt.BTypeGetName()) {
			t.Errorf("mismatch in constant type name %d after round-trip for %s", i, birFile)
		}
		if ct.BTypeGetFlags() != loadedCt.BTypeGetFlags() {
			t.Errorf("mismatch in constant type flags %d after round-trip for %s", i, birFile)
		}

		cvt := constant.ConstValue.Type.(ast.BType)
		loadedCvt := loadedConstant.ConstValue.Type.(ast.BType)

		if cvt.BTypeGetTag() != loadedCvt.BTypeGetTag() {
			t.Errorf("mismatch in constant value type tag %d after round-trip for %s", i, birFile)
		}
		if string(cvt.BTypeGetName()) != string(loadedCvt.BTypeGetName()) {
			t.Errorf("mismatch in constant value type name %d after round-trip for %s", i, birFile)
		}
		if cvt.BTypeGetFlags() != loadedCvt.BTypeGetFlags() {
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

		globalVarType := globalVar.Type.(ast.BType)
		loadedGlobalVarType := loadedGlobalVar.Type.(ast.BType)

		if globalVarType.BTypeGetTag() != loadedGlobalVarType.BTypeGetTag() {
			t.Errorf("mismatch in global variable type tag %d after round-trip for %s", i, birFile)
		}
		if string(globalVarType.BTypeGetName()) != string(loadedGlobalVarType.BTypeGetName()) {
			t.Errorf("mismatch in global variable type name %d after round-trip for %s", i, birFile)
		}
		if globalVarType.BTypeGetFlags() != loadedGlobalVarType.BTypeGetFlags() {
			t.Errorf("mismatch in global variable type flags %d after round-trip for %s", i, birFile)
		}

		for i, fn := range pkg.Functions {
			loadedFn := loadedPkg.Functions[i]
			if fn.Name != loadedFn.Name {
				t.Errorf("mismatch in function name %d after round-trip for %s", i, birFile)
			}
			if fn.OriginalName != loadedFn.OriginalName {
				t.Errorf("mismatch in function original name %d after round-trip for %s", i, birFile)
			}
			if fn.Flags != loadedFn.Flags {
				t.Errorf("mismatch in function flags %d after round-trip for %s", i, birFile)
			}
			if fn.Origin != loadedFn.Origin {
				t.Errorf("mismatch in function origin %d after round-trip for %s", i, birFile)
			}

			for rpi, rp := range fn.RequiredParams {
				loadedRp := loadedFn.RequiredParams[rpi]
				if rp.Name.Value() != loadedRp.Name.Value() {
					t.Logf("Expected: %s, Got: %s", rp.Name.Value(), loadedRp.Name.Value())
					t.Errorf("mismatch in required param name %d of function %d after round-trip for %s", rpi, i, birFile)
				}
				if rp.Flags != loadedRp.Flags {
					t.Logf("Expected: %d, Got: %d", rp.Flags, loadedRp.Flags)
					t.Errorf("mismatch in required param flags %d of function %d after round-trip for %s", rpi, i, birFile)
				}
			}

			if fn.ReturnVariable != nil && loadedFn.ReturnVariable != nil {
				if fn.ReturnVariable.Kind != loadedFn.ReturnVariable.Kind {
					t.Errorf("mismatch in return variable kind of function %d after round-trip for %s", i, birFile)
				}

				returnVarType := fn.ReturnVariable.Type.(ast.BType)
				loadedReturnVarType := loadedFn.ReturnVariable.Type.(ast.BType)

				if returnVarType.BTypeGetTag() != loadedReturnVarType.BTypeGetTag() {
					t.Errorf("mismatch in return variable type tag of function %d after round-trip for %s", i, birFile)
				}
				if string(returnVarType.BTypeGetName()) != string(loadedReturnVarType.BTypeGetName()) {
					t.Errorf("mismatch in return variable type name of function %d after round-trip for %s", i, birFile)
				}
				if returnVarType.BTypeGetFlags() != loadedReturnVarType.BTypeGetFlags() {
					t.Errorf("mismatch in return variable type flags of function %d after round-trip for %s", i, birFile)
				}

				if fn.ReturnVariable.Name.Value() != loadedFn.ReturnVariable.Name.Value() {
					t.Errorf("mismatch in return variable name of function %d after round-trip for %s", i, birFile)
				}
			}
		}
	}
}

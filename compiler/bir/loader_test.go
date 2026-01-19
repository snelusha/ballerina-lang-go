/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 */

package bir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func TestLoadBIRPackageFromReader(t *testing.T) {
	// Use the same sample.bir file from testdata
	file, err := os.Open(filepath.Join("testdata", "sample.bir"))
	if err != nil {
		t.Fatalf("Failed to open sample.bir: %v", err)
	}
	defer file.Close()

	// Load using our loader
	pkg, err := LoadBIRPackageFromReader(file)
	if err != nil {
		t.Fatalf("LoadBIRPackageFromReader failed: %v", err)
	}

	if pkg == nil {
		t.Fatal("LoadBIRPackageFromReader returned nil package")
	}

	// Verify package ID
	pkgID := pkg.GetPackageID()
	if pkgID.OrgName.GetValue() == "" && pkgID.Name.GetValue() == "" {
		t.Error("PackageID should have at least org or name set")
	}

	// Verify imports
	imports := pkg.GetImportModules()
	if imports == nil {
		t.Error("Imports should not be nil")
	}
	t.Logf("Loaded %d imports", len(*imports))

	// Verify functions
	functions := pkg.GetFunctions()
	if functions == nil {
		t.Error("Functions should not be nil")
	}
	t.Logf("Loaded %d functions", len(*functions))

	// Print function names (similar to the example)
	for _, fn := range *functions {
		t.Logf("Function: %s (original: %s, worker: %s)",
			fn.GetName().GetValue(),
			fn.GetOriginalName().GetValue(),
			fn.GetWorkerName().GetValue())

		// Check required params
		params := fn.GetRequiredParams()
		if params != nil {
			t.Logf("  Required params: %d", len(*params))
			for _, p := range *params {
				t.Logf("    - %s (flags: %d)", p.GetName().GetValue(), p.GetFlags())
			}
		}
	}

	// Verify other collections exist (even if empty)
	typeDefs := pkg.GetTypeDefs()
	if typeDefs == nil {
		t.Error("TypeDefs should not be nil")
	}
	t.Logf("Type definitions: %d", len(*typeDefs))

	globalVars := pkg.GetGlobalVars()
	if globalVars == nil {
		t.Error("GlobalVars should not be nil")
	}
	t.Logf("Global vars: %d", len(*globalVars))

	annotations := pkg.GetAnnotations()
	if annotations == nil {
		t.Error("Annotations should not be nil")
	}
	t.Logf("Annotations: %d", len(*annotations))

	constants := pkg.GetConstants()
	if constants == nil {
		t.Error("Constants should not be nil")
	}
	t.Logf("Constants: %d", len(*constants))

	serviceDecls := pkg.GetServiceDecls()
	if serviceDecls == nil {
		t.Error("ServiceDecls should not be nil")
	}
	t.Logf("Service declarations: %d", len(*serviceDecls))
}

func TestLoadBIRPackageFromReader_CompareWithDirectRead(t *testing.T) {
	// Read directly using Kaitai (like the example)
	file1, err := os.Open(filepath.Join("testdata", "sample.bir"))
	if err != nil {
		t.Fatalf("Failed to open sample.bir: %v", err)
	}
	defer file1.Close()

	birModel := NewBir()
	err = birModel.Read(kaitai.NewStream(file1), nil, birModel)
	if err != nil {
		t.Fatalf("Direct Kaitai read failed: %v", err)
	}

	// Read using our loader
	file2, err := os.Open(filepath.Join("testdata", "sample.bir"))
	if err != nil {
		t.Fatalf("Failed to open sample.bir: %v", err)
	}
	defer file2.Close()

	pkg, err := LoadBIRPackageFromReader(file2)
	if err != nil {
		t.Fatalf("LoadBIRPackageFromReader failed: %v", err)
	}

	// Compare function counts
	directFuncCount := len(birModel.Module.Functions)
	loadedFuncs := pkg.GetFunctions()
	loadedFuncCount := 0
	if loadedFuncs != nil {
		loadedFuncCount = len(*loadedFuncs)
	}

	if directFuncCount != loadedFuncCount {
		t.Errorf("Function count mismatch: direct read=%d, loader=%d",
			directFuncCount, loadedFuncCount)
	}

	// Compare import counts
	directImportCount := len(birModel.Module.Imports)
	loadedImports := pkg.GetImportModules()
	loadedImportCount := 0
	if loadedImports != nil {
		loadedImportCount = len(*loadedImports)
	}

	if directImportCount != loadedImportCount {
		t.Errorf("Import count mismatch: direct read=%d, loader=%d",
			directImportCount, loadedImportCount)
	}

	// Compare function names
	if directFuncCount > 0 && loadedFuncCount > 0 {
		for i, directFunc := range birModel.Module.Functions {
			if i >= loadedFuncCount {
				break
			}

			// Get name from direct read
			var directName string
			if nameCp, ok := birModel.ConstantPool.ConstantPoolEntries[directFunc.NameCpIndex].CpInfo.(*Bir_StringCpInfo); ok {
				directName = nameCp.Value
			}

			// Get name from loaded model
			loadedFunc := (*loadedFuncs)[i]
			loadedName := loadedFunc.GetName().GetValue()

			if directName != loadedName {
				t.Errorf("Function name mismatch at index %d: direct=%s, loaded=%s",
					i, directName, loadedName)
			} else {
				t.Logf("✓ Function %d: %s", i, directName)
			}
		}
	}
}

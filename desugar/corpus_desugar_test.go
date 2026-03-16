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

package desugar_test

import (
	"flag"
	"fmt"
	"testing"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/context"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/test_util/testphases"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update expected desugared AST files")

type positionTestVisitor struct {
	missingPositions []string
	missingCount     int
}

func (v *positionTestVisitor) Visit(node ast.BLangNode) ast.Visitor {
	if node == nil {
		return nil
	}
	if node.GetPosition() == nil {
		v.missingCount++
		if len(v.missingPositions) < 25 {
			v.missingPositions = append(v.missingPositions, fmt.Sprintf("%T", node))
		}
	}
	return v
}

func (v *positionTestVisitor) VisitTypeData(typeData *model.TypeData) ast.Visitor {
	return v
}

func TestDesugar(t *testing.T) {
	flag.Parse()

	testPairs := test_util.GetValidAndPanicTests(t, test_util.Desugar)

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testDesugar(t, testPair)
		})
	}
}

func testDesugar(t *testing.T, testCase test_util.TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Desugar panicked for %s: %v", testCase.InputPath, r)
		}
	}()

	env := context.NewCompilerEnvironment(semtypes.CreateTypeEnv())
	cx := context.NewCompilerContext(env)
	result, err := testphases.RunPipeline(cx, testphases.PhaseDesugar, testCase.InputPath)
	if err != nil {
		t.Errorf("pipeline failed for %s: %v", testCase.InputPath, err)
		return
	}

	// Ensure every visited AST node in the desugared tree has position info.
	posVisitor := &positionTestVisitor{}
	ast.Walk(posVisitor, result.Package)
	if posVisitor.missingCount > 0 {
		t.Errorf(
			"found %d desugared AST nodes without position for %s (showing up to %d types): %v",
			posVisitor.missingCount,
			testCase.InputPath,
			len(posVisitor.missingPositions),
			posVisitor.missingPositions,
		)
	}

	// Serialize AST after desugaring
	prettyPrinter := ast.PrettyPrinter{}
	actualAST := prettyPrinter.Print(result.Package)

	// If update flag is set, update expected file
	if *update {
		if test_util.UpdateIfNeeded(t, testCase.ExpectedPath, actualAST) {
			t.Errorf("updated expected desugared AST file: %s", testCase.ExpectedPath)
		}
		return
	}

	// Read expected AST file
	expectedAST := test_util.ReadExpectedFile(t, testCase.ExpectedPath)

	// Compare AST strings exactly
	if actualAST != expectedAST {
		t.Errorf("Desugared AST mismatch for %s\nExpected file: %s\n%s",
			testCase.InputPath, testCase.ExpectedPath, getDiff(expectedAST, actualAST))
		return
	}

	t.Logf("Desugar completed successfully for %s", testCase.InputPath)
}

func getDiff(expectedAST, actualAST string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedAST, actualAST, false)
	return dmp.DiffPrettyText(diffs)
}

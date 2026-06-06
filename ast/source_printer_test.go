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

package ast

import (
	"testing"

	"ballerina-lang-go/model"
)

func TestSourcePrinter(t *testing.T) {
	intType := &BLangValueType{TypeKind: TypeKind_INT}
	stringType := &BLangValueType{TypeKind: TypeKind_STRING}

	pkg := &BLangPackage{
		TopLevelNodes: []TopLevelNode{
			&BLangImportPackage{
				OrgName:      &BLangIdentifier{Value: "ballerina"},
				PkgNameComps: []BLangIdentifier{{Value: "io"}},
			},
			&BLangSimpleVariable{
				BLangVariableBase: BLangVariableBase{typeNode: stringType, Expr: &BLangLiteral{Value: "hello"}},
				Name:              &BLangIdentifier{Value: "greeting"},
			},
			&BLangFunction{
				bLangInvokableNodeBase: bLangInvokableNodeBase{
					Name: BLangIdentifier{Value: "add"},
					RequiredParams: []BLangSimpleVariable{
						{BLangVariableBase: BLangVariableBase{typeNode: intType}, Name: &BLangIdentifier{Value: "x"}},
						{BLangVariableBase: BLangVariableBase{typeNode: intType}, Name: &BLangIdentifier{Value: "y"}},
					},
					returnTypeDescriptor: intType,
					Body: &BLangBlockFunctionBody{Stmts: []StatementNode{
						&BLangReturn{Expr: &BLangBinaryExpr{
							LhsExpr: &BLangSimpleVarRef{VariableName: &BLangIdentifier{Value: "x"}},
							RhsExpr: &BLangSimpleVarRef{VariableName: &BLangIdentifier{Value: "y"}},
							OpKind:  model.OperatorKind_ADD,
						}},
					}},
					flags: model.FlagPublic,
				},
			},
		},
	}

	actual := ToSource(pkg)
	expected := `import ballerina/io;

string greeting = "hello";

public function add(int x, int y) returns int {
	return x + y;
}`
	if actual != expected {
		t.Fatalf("source mismatch\nexpected:\n%s\nactual:\n%s", expected, actual)
	}
}

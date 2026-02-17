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

package context

import (
	"fmt"

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

type CompilerContext struct {
	env *CompilerEnvironment
}

func NewCompilerContext(env *CompilerEnvironment) *CompilerContext {
	return &CompilerContext{
		env: env,
	}
}

func (this *CompilerContext) NewSymbolSpace(packageId model.PackageID) *model.SymbolSpace {
	return this.env.NewSymbolSpace(packageId)
}

func (this *CompilerContext) NewFunctionScope(parent model.Scope, pkg model.PackageID) *model.FunctionScope {
	return this.env.NewFunctionScope(parent, pkg)
}

func (this *CompilerContext) NewBlockScope(parent model.Scope, pkg model.PackageID) *model.BlockScope {
	return this.env.NewBlockScope(parent, pkg)
}

func (this *CompilerContext) GetSymbol(symbol model.Symbol) model.Symbol {
	return this.env.GetSymbol(symbol)
}

func (this *CompilerContext) RefSymbol(symbol model.Symbol) model.SymbolRef {
	return this.env.RefSymbol(symbol)
}

func (this *CompilerContext) SymbolName(symbol model.Symbol) string {
	return this.GetSymbol(symbol).Name()
}

func (this CompilerContext) SymbolType(symbol model.Symbol) semtypes.SemType {
	return this.GetSymbol(symbol).Type()
}

func (this CompilerContext) SymbolKind(symbol model.Symbol) model.SymbolKind {
	return this.GetSymbol(symbol).Kind()
}

func (this CompilerContext) SymbolIsPublic(symbol model.Symbol) bool {
	return this.GetSymbol(symbol).IsPublic()
}

func (this CompilerContext) SetSymbolType(symbol model.Symbol, ty semtypes.SemType) {
	this.GetSymbol(symbol).SetType(ty)
}

func (this CompilerContext) GetDefaultPackage() *model.PackageID {
	return this.env.GetDefaultPackage()
}

func (this CompilerContext) NewPackageID(orgName model.Name, nameComps []model.Name, version model.Name) *model.PackageID {
	return this.env.NewPackageID(orgName, nameComps, version)
}

func (this *CompilerEnvironment) Unimplemented(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Unimplemented: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Unimplemented: %s", message))
}

func (this *CompilerEnvironment) SemanticError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Semantic error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Semantic error: %s", message))
}

// TODO: implement these properly
func (this *CompilerEnvironment) SyntaxError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Syntax error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Syntax error: %s", message))
}

func (this *CompilerEnvironment) InternalError(message string, pos diagnostics.Location) {
	if pos != nil {
		panic(fmt.Sprintf("Internal error: %s at %s", message, pos))
	}
	panic(fmt.Sprintf("Internal error: %s", message))
}

// GetTypeEnv returns the type environment for this context
func (this CompilerContext) GetTypeEnv() semtypes.Env {
	return this.env.GetTypeEnv()
}

func (this CompilerContext) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	return this.env.GetNextAnonymousTypeKey(packageID)
}

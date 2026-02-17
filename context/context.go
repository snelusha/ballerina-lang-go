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
	env         *CompilerEnvironment
	diagnostics []diagnostics.Diagnostic
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
	// If this happen that's a bug in SymbolResolver
	if symbol == nil {
		this.InternalError("RefSymbol called with nil symbol", nil)
	}
	if refSymbol, ok := symbol.(*model.SymbolRef); ok {
		return *refSymbol
	}
	// This should never happen because we should never store actual symbols in the AST
	this.InternalError(fmt.Sprintf("Symbol is not a SymbolRef: type=%T, name=%s, kind=%v", symbol, symbol.Name(), symbol.Kind()), nil)
	return model.SymbolRef{}
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

func (this *CompilerContext) Unimplemented(message string, pos diagnostics.Location) {
	panic(fmt.Sprintf("Unimplemented: %s at %s", message, pos))
}

func (this *CompilerContext) SemanticError(message string, pos diagnostics.Location) {
	code := "SEMANTIC_ERROR"
	diagnosticInfo := diagnostics.NewDiagnosticInfo(&code, message, diagnostics.Error)
	dignostic := diagnostics.CreateDiagnostic(diagnosticInfo, pos)
	this.diagnostics = append(this.diagnostics, dignostic)
}

func (this *CompilerContext) SyntaxError(message string, pos diagnostics.Location) {
	code := "SYNTAX_ERROR"
	diagnosticInfo := diagnostics.NewDiagnosticInfo(&code, message, diagnostics.Error)
	dignostic := diagnostics.CreateDiagnostic(diagnosticInfo, pos)
	this.diagnostics = append(this.diagnostics, dignostic)
}

func (this *CompilerContext) InternalError(message string, pos diagnostics.Location) {
	panic(fmt.Sprintf("Internal error: %s at %s", message, pos))
}

func (this *CompilerContext) Diagnostics() []diagnostics.Diagnostic {
	return this.diagnostics
}

func (this *CompilerContext) HasDiagnostics() bool {
	return len(this.diagnostics) > 0
}

// GetTypeEnv returns the type environment for this context
func (this CompilerContext) GetTypeEnv() semtypes.Env {
	return this.env.GetTypeEnv()
}

func (this CompilerContext) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	return this.env.GetNextAnonymousTypeKey(packageID)
}

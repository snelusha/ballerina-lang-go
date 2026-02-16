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
	"strconv"

	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"
)

// TODO: consider moving type resolution env in to this
type CompilerContext struct {
	anonTypeCount   map[*model.PackageID]int
	packageInterner *model.PackageIDInterner
	symbolSpaces    []*model.SymbolSpace
	typeEnv         semtypes.Env
	diagnostics     []diagnostics.Diagnostic
}

func (this *CompilerContext) NewSymbolSpace(packageId model.PackageID) *model.SymbolSpace {
	space := model.NewSymbolSpaceInner(packageId, len(this.symbolSpaces))
	this.symbolSpaces = append(this.symbolSpaces, space)
	return space
}

func (this *CompilerContext) NewFunctionScope(parent model.Scope, pkg model.PackageID) *model.FunctionScope {
	return &model.FunctionScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerContext) NewBlockScope(parent model.Scope, pkg model.PackageID) *model.BlockScope {
	return &model.BlockScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerContext) GetSymbol(symbol model.Symbol) model.Symbol {
	if refSymbol, ok := symbol.(*model.SymbolRef); ok {
		symbolSpace := this.symbolSpaces[refSymbol.SpaceIndex]
		return symbolSpace.Symbols[refSymbol.Index]
	}
	return symbol
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

func (this *CompilerContext) SymbolType(symbol model.Symbol) semtypes.SemType {
	return this.GetSymbol(symbol).Type()
}

func (this *CompilerContext) SymbolKind(symbol model.Symbol) model.SymbolKind {
	return this.GetSymbol(symbol).Kind()
}

func (this *CompilerContext) SymbolIsPublic(symbol model.Symbol) bool {
	return this.GetSymbol(symbol).IsPublic()
}

func (this *CompilerContext) SetSymbolType(symbol model.Symbol, ty semtypes.SemType) {
	this.GetSymbol(symbol).SetType(ty)
}

func (this *CompilerContext) GetDefaultPackage() *model.PackageID {
	return this.packageInterner.GetDefaultPackage()
}

func (this *CompilerContext) NewPackageID(orgName model.Name, nameComps []model.Name, version model.Name) *model.PackageID {
	return model.NewPackageID(this.packageInterner, orgName, nameComps, version)
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

func NewCompilerContext(typeEnv semtypes.Env) *CompilerContext {
	return &CompilerContext{
		anonTypeCount:   make(map[*model.PackageID]int),
		packageInterner: model.DefaultPackageIDInterner,
		typeEnv:         typeEnv,
	}
}

// GetTypeEnv returns the type environment for this context
func (this *CompilerContext) GetTypeEnv() semtypes.Env {
	return this.typeEnv
}

const (
	ANON_PREFIX       = "$anon"
	BUILTIN_ANON_TYPE = ANON_PREFIX + "Type$builtin$"
	ANON_TYPE         = ANON_PREFIX + "Type$"
)

func (this *CompilerContext) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	nextValue := this.anonTypeCount[packageID]
	this.anonTypeCount[packageID] = nextValue + 1
	if packageID != nil && model.ANNOTATIONS_PKG != packageID {
		return BUILTIN_ANON_TYPE + "_" + strconv.Itoa(nextValue)
	}
	return ANON_TYPE + "_" + strconv.Itoa(nextValue)
}

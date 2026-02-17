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
)

// TODO: consider moving type resolution env in to this
type CompilerEnvironment struct {
	anonTypeCount   map[*model.PackageID]int
	packageInterner *model.PackageIDInterner
	symbolSpaces    []*model.SymbolSpace
	typeEnv         semtypes.Env
}

func (this *CompilerEnvironment) NewSymbolSpace(packageId model.PackageID) *model.SymbolSpace {
	space := model.NewSymbolSpaceInner(packageId, len(this.symbolSpaces))
	this.symbolSpaces = append(this.symbolSpaces, space)
	return space
}

func (this *CompilerEnvironment) NewFunctionScope(parent model.Scope, pkg model.PackageID) *model.FunctionScope {
	return &model.FunctionScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerEnvironment) NewBlockScope(parent model.Scope, pkg model.PackageID) *model.BlockScope {
	return &model.BlockScope{
		BlockScopeBase: model.BlockScopeBase{
			Parent: parent,
			Main:   this.NewSymbolSpace(pkg),
		},
	}
}

func (this *CompilerEnvironment) GetSymbol(symbol model.Symbol) model.Symbol {
	if refSymbol, ok := symbol.(*model.SymbolRef); ok {
		symbolSpace := this.symbolSpaces[refSymbol.SpaceIndex]
		return symbolSpace.Symbols[refSymbol.Index]
	}
	return symbol
}

func (this *CompilerEnvironment) RefSymbol(symbol model.Symbol) model.SymbolRef {
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

func (this *CompilerEnvironment) SymbolName(symbol model.Symbol) string {
	return this.GetSymbol(symbol).Name()
}

func (this *CompilerEnvironment) SymbolType(symbol model.Symbol) semtypes.SemType {
	return this.GetSymbol(symbol).Type()
}

func (this *CompilerEnvironment) SymbolKind(symbol model.Symbol) model.SymbolKind {
	return this.GetSymbol(symbol).Kind()
}

func (this *CompilerEnvironment) SymbolIsPublic(symbol model.Symbol) bool {
	return this.GetSymbol(symbol).IsPublic()
}

func (this *CompilerEnvironment) SetSymbolType(symbol model.Symbol, ty semtypes.SemType) {
	this.GetSymbol(symbol).SetType(ty)
}

func (this *CompilerEnvironment) GetDefaultPackage() *model.PackageID {
	return this.packageInterner.GetDefaultPackage()
}

func (this *CompilerEnvironment) NewPackageID(orgName model.Name, nameComps []model.Name, version model.Name) *model.PackageID {
	return model.NewPackageID(this.packageInterner, orgName, nameComps, version)
}

func NewCompilerEnvironment(typeEnv semtypes.Env) *CompilerEnvironment {
	return &CompilerEnvironment{
		anonTypeCount:   make(map[*model.PackageID]int),
		packageInterner: model.DefaultPackageIDInterner,
		typeEnv:         typeEnv,
	}
}

// GetTypeEnv returns the type environment for this context
func (this *CompilerEnvironment) GetTypeEnv() semtypes.Env {
	return this.typeEnv
}

const (
	ANON_PREFIX       = "$anon"
	BUILTIN_ANON_TYPE = ANON_PREFIX + "Type$builtin$"
	ANON_TYPE         = ANON_PREFIX + "Type$"
)

func (this *CompilerEnvironment) GetNextAnonymousTypeKey(packageID *model.PackageID) string {
	nextValue := this.anonTypeCount[packageID]
	this.anonTypeCount[packageID] = nextValue + 1
	if packageID != nil && model.ANNOTATIONS_PKG != packageID {
		return BUILTIN_ANON_TYPE + "_" + strconv.Itoa(nextValue)
	}
	return ANON_TYPE + "_" + strconv.Itoa(nextValue)
}

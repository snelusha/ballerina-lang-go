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

package compile

import (
	"ballerina-lang-go/context"
	libcommon "ballerina-lang-go/lib/common"
	"ballerina-lang-go/model"
	"ballerina-lang-go/semtypes"
)

func GetIoSymbols(ctx *context.CompilerContext) model.ExportedSymbolSpace {
	pkg := model.NewPackageID(
		model.DefaultPackageIDInterner,
		model.Name("ballerina"),
		[]model.Name{model.Name("io")},
		model.Name("0.0.1"),
	)
	space := ctx.NewSymbolSpace(*pkg)
	printLnSignature := model.FunctionSignature{
		RestParamType: &semtypes.ANY,
		ReturnType:    &semtypes.NIL,
	}
	printLnSymbol := model.NewFunctionSymbol("println", printLnSignature, true)
	ctx.SetSymbolType(printLnSymbol, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &printLnSignature))

	fileReadStringSignature := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{&semtypes.ANY},
		ReturnType: &semtypes.STRING,
	}
	fileReadStringSymbol := model.NewFunctionSymbol("fileReadString", fileReadStringSignature, true)
	ctx.SetSymbolType(fileReadStringSymbol, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &fileReadStringSignature))

	fileWriteStringSignature := model.FunctionSignature{
		ParamTypes: []semtypes.SemType{&semtypes.STRING, &semtypes.STRING},
		ReturnType: semtypes.Union(&semtypes.NIL, &semtypes.ERROR),
	}
	fileWriteStringSymbol := model.NewFunctionSymbol("fileWriteString", fileWriteStringSignature, true)
	ctx.SetSymbolType(fileWriteStringSymbol, libcommon.FunctionSignatureToSemType(ctx.GetTypeEnv(), &fileWriteStringSignature))

	space.AddSymbol("println", printLnSymbol)
	space.AddSymbol("fileReadString", fileReadStringSymbol)
	space.AddSymbol("fileWriteString", fileWriteStringSymbol)
	return model.ExportedSymbolSpace{
		Main: space,
	}
}

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

package io

import (
	"fmt"
	"os"

	"ballerina-lang-go/runtime"
)

const (
	orgName    = "ballerina"
	moduleName = "io"
	funcName   = "println"
)

func Println(values ...any) {
	fmt.Fprintln(os.Stdout, values...)
}

func printlnExtern(args []any) (any, error) {
	Println(args...)
	return nil, nil
}

func fileReadStringExtern(args []any) (any, error) {
	path := args[0].(string)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return string(content), nil
}

func fileWriteStringExtern(args []any) (any, error) {
	path := args[0].(string)
	content := args[1].(string)
	err := os.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func initIOModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "println", printlnExtern)
	runtime.RegisterExternFunction(rt, orgName, moduleName, "fileReadString", fileReadStringExtern)
	runtime.RegisterExternFunction(rt, orgName, moduleName, "fileWriteString", fileWriteStringExtern)
}

func init() {
	runtime.RegisterModuleInitializer(initIOModule)
}

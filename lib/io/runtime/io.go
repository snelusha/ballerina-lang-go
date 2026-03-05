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
	"strings"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/values"
)

const (
	orgName    = "ballerina"
	moduleName = "io"
	funcName   = "println"
)

func Println(vals ...values.BalValue) {
	parts := make([]string, len(vals))
	visited := make(map[uintptr]bool)
	for i, v := range vals {
		parts[i] = values.String(v, visited)
	}
	fmt.Fprintln(os.Stdout, strings.Join(parts, ""))
}

func printlnExtern(args []values.BalValue) (values.BalValue, error) {
	Println(args...)
	return nil, nil
}

func FileReadString(path values.BalValue) values.BalValue {
	content, err := os.ReadFile(values.String(path, nil))
	if err != nil {
		panic(err)
	}
	return string(content)
}

func fileReadStringExtern(args []values.BalValue) (values.BalValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("fileReadString expects exactly one argument")
	}
	return FileReadString(args[0]), nil
}

func initIOModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "println", printlnExtern)
	runtime.RegisterExternFunction(rt, orgName, moduleName, "fileReadString", fileReadStringExtern)
}

func init() {
	runtime.RegisterModuleInitializer(initIOModule)
}

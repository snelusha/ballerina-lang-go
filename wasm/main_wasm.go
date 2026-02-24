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

package main

import (
	"syscall/js"

	"ballerina-lang-go/projects/directory"
)

func main() {
	js.Global().Set("run", js.FuncOf(run))

	select {}
}

func run(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return map[string]any{
			"error": "expected at least 2 arguments: (fsProxy, path)",
		}
	}

	proxy := args[0]
	path := args[1].String()

	fsys := NewLocalStorageFS(proxy)

	result, err := directory.LoadProject(fsys, path)
	if err != nil {
		return map[string]any{
			"error": err.Error(),
		}
	}

	_ = result

	return nil
}

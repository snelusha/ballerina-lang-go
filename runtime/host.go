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

package runtime

import (
	"io"
	"io/fs"
	"os"
)

// Host carries platform capabilities injected by the embedder (CLI, tests, WASM).
// Paths read from Host.FS are relative to that filesystem's root.
type Host struct {
	FS     fs.FS
	Stdout io.Writer
}

// DefaultHost returns a host rooted at the current working directory with
// process standard output.
func DefaultHost() Host {
	return Host{
		FS:     os.DirFS("."),
		Stdout: os.Stdout,
	}
}

func normalizeHost(h Host) Host {
	if h.FS == nil {
		h.FS = os.DirFS(".")
	}
	if h.Stdout == nil {
		h.Stdout = os.Stdout
	}
	return h
}

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

package bfs

import (
	"io/fs"
	"path"
)

// PathFS extends fs.FS with path manipulation methods.
type PathFS interface {
	fs.FS

	Join(elem ...string) string

	Dir(p string) string
	Base(p string) string
	Abs(p string) (string, error)
}

func Join(fsys fs.FS, elem ...string) string {
	if pfs, ok := fsys.(PathFS); ok {
		return pfs.Join(elem...)
	}
	return path.Join(elem...)
}

func Dir(fsys fs.FS, p string) string {
	if pfs, ok := fsys.(PathFS); ok {
		return pfs.Dir(p)
	}
	return path.Dir(p)
}

func Base(fsys fs.FS, p string) string {
	if pfs, ok := fsys.(PathFS); ok {
		return pfs.Base(p)
	}
	return path.Base(p)
}

func Abs(fsys fs.FS, p string) (string, error) {
	if pfs, ok := fsys.(PathFS); ok {
		return pfs.Abs(p)
	}
	return p, nil
}

// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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
	"os"
	"path/filepath"
)

// dirFS wraps os.DirFS and implements PathFS and ReadDirFS interfaces.
type dirFS struct {
	fs.FS
	baseDir string
}

func NewDirFS(baseDir string) (fs.FS, error) {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, err
	}
	return &dirFS{
		FS:      os.DirFS(absBaseDir),
		baseDir: absBaseDir,
	}, nil
}

func (d *dirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(d.FS, name)
}

func (d *dirFS) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (d *dirFS) Dir(p string) string {
	return filepath.Dir(p)
}

func (d *dirFS) Base(p string) string {
	return filepath.Base(p)
}

func (d *dirFS) Abs(p string) (string, error) {
	if filepath.IsAbs(p) {
		return filepath.Clean(p), nil
	}
	abs := filepath.Join(d.baseDir, p)
	return filepath.Clean(abs), nil
}

func (d *dirFS) BaseDir() string {
	return d.baseDir
}

// DirFSForPath resolves path (file or directory) to an absolute path, then returns
// an fs.FS rooted at the containing directory and the load path relative to that root.
// Use this with directory.LoadProject(fsys, loadPath, ...). Uses filepath and os.
func DirFSForPath(path string) (fsys fs.FS, loadPath string, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, "", err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, "", err
	}
	var fsBase string
	if info.IsDir() {
		fsBase = absPath
		loadPath = "."
	} else {
		fsBase = filepath.Dir(absPath)
		loadPath = filepath.Base(absPath)
	}
	fsys, err = NewDirFS(fsBase)
	if err != nil {
		return nil, "", err
	}
	return fsys, loadPath, nil
}

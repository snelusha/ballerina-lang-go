package bfs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"
)

type memFS struct {
	files map[string]*memFile
}

func (mfs *memFS) Create(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "create", Path: name, Err: fs.ErrInvalid}
	}

	file := &memFile{
		name: name,
	}
	mfs.files[name] = file

	return &openMemFile{
		memFile: file,
		Reader:  bytes.NewReader(file.data),
	}, nil
}

func (mfs *memFS) MkdirAll(path string, perm fs.FileMode) error {
	return nil
}

func (mfs *memFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "openfile", Path: name, Err: fs.ErrInvalid}
	}

	file, ok := mfs.files[name]
	if !ok {
		file = &memFile{
			name: name,
			mode: perm,
		}
		mfs.files[name] = file
	}

	return &openMemFile{
		memFile: file,
		Reader:  bytes.NewReader(file.data),
	}, nil
}

func (mfs *memFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	// if !fs.ValidPath(name) {
	// 	return &fs.PathError{Op: "writefile", Path: name, Err: fs.ErrInvalid}
	// }

	mfs.files[name] = &memFile{
		name: name,
		data: data,
		mode: perm,
	}

	return nil
}

func (mfs *memFS) Remove(name string) error {
	// Try exact match first
	if _, ok := mfs.files[name]; ok {
		delete(mfs.files, name)
		return nil
	}

	// Try directory removal (prefix match)
	dir := name
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	found := false
	for fname := range mfs.files {
		if strings.HasPrefix(fname, dir) {
			found = true
			delete(mfs.files, fname)
		}
	}

	if !found {
		return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrNotExist}
	}
	return nil
}

func (mfs *memFS) Rename(oldpath string, newpath string) error {
	// Try exact match first
	if file, ok := mfs.files[oldpath]; ok {
		delete(mfs.files, oldpath)
		file.name = newpath
		mfs.files[newpath] = file
		return nil
	}

	// Try directory rename (prefix match)
	oldDir := oldpath
	if !strings.HasSuffix(oldDir, "/") {
		oldDir += "/"
	}

	found := false
	for name, file := range mfs.files {
		if strings.HasPrefix(name, oldDir) {
			found = true
			newName := newpath + strings.TrimPrefix(name, oldpath)
			delete(mfs.files, name)
			file.name = newName
			mfs.files[newName] = file
		}
	}

	if !found {
		return &fs.PathError{Op: "rename", Path: oldpath, Err: fs.ErrNotExist}
	}
	return nil
}

type memFile struct {
	name  string
	data  []byte
	mode  fs.FileMode
	isDir bool
}

type openMemFile struct {
	*memFile
	*bytes.Reader
}

func (o *openMemFile) IsDir() bool {
	return o.isDir
}

func (o *openMemFile) ModTime() time.Time {
	return time.Time{}
}

func (o *openMemFile) Mode() fs.FileMode {
	return o.mode
}

func (o *openMemFile) Name() string {
	return o.name
}

func (o *openMemFile) Size() int64 {
	return int64(len(o.data))
}

func (o *openMemFile) Sys() any {
	return nil
}

func (o *openMemFile) Close() error {
	return nil
}

func (o *openMemFile) Read(p []byte) (int, error) {
	return o.Reader.Read(p)
}

func (o *openMemFile) Stat() (fs.FileInfo, error) {
	return o, nil
}

func (mfs *memFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	file, ok := mfs.files[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	return &openMemFile{
		memFile: file,
		Reader:  bytes.NewReader(file.data),
	}, nil
}

func CopyToFile(fsys fs.FS, name string, src io.Reader) error {
	file, err := Create(fsys, name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, ok := file.(io.Writer)
	if !ok {
		return &fs.PathError{Op: "copy", Path: name, Err: fs.ErrInvalid}
	}

	_, err = io.Copy(writer, src)
	return err
}

func CopyFromFile(fsys fs.FS, name string, dst io.Writer) error {
	file, err := fsys.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(dst, file)
	return err
}

func PrintFiles(fsys fs.FS) {
	fmt.Printf("%+v\n", fsys.(*memFS).files)
}

func NewMemFS() fs.FS {
	return &memFS{
		files: make(map[string]*memFile),
	}
}

package bfs

import "io/fs"

type Writable interface {
	fs.FS

	WriteFile(name string, data []byte, perm fs.FileMode) error
}

type Mutable interface {
	fs.FS

	Create(name string) (fs.File, error)
	Open(name string) (fs.File, error)
	OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error)
	Remove(name string) error
	Rename(oldpath, newpath string) error
	MkdirAll(path string, perm fs.FileMode) error
}

func WriteFile(fsys fs.FS, name string, data []byte, perm fs.FileMode) error {
	if wfs, ok := fsys.(Writable); ok {
		return wfs.WriteFile(name, data, perm)
	}
	return &fs.PathError{Op: "writefile", Path: name, Err: fs.ErrInvalid}
}

func Create(fsys fs.FS, name string) (fs.File, error) {
	mfs, ok := fsys.(Mutable)
	if !ok {
		return nil, &fs.PathError{Op: "create", Path: name, Err: fs.ErrInvalid}
	}
	return mfs.Create(name)
}

func MkdirAll(fsys fs.FS, path string, perm fs.FileMode) error {
	mfs, ok := fsys.(Mutable)
	if !ok {
		return &fs.PathError{Op: "mkdirall", Path: path, Err: fs.ErrInvalid}
	}
	return mfs.MkdirAll(path, perm)
}

func OpenFile(fsys fs.FS, name string, flag int, perm fs.FileMode) (fs.File, error) {
	mfs, ok := fsys.(Mutable)
	if !ok {
		return nil, &fs.PathError{Op: "openfile", Path: name, Err: fs.ErrInvalid}
	}
	return mfs.OpenFile(name, flag, perm)
}

func Remove(fsys fs.FS, name string) error {
	mfs, ok := fsys.(Mutable)
	if !ok {
		return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrInvalid}
	}
	return mfs.Remove(name)
}

func Rename(fsys fs.FS, oldpath, newpath string) error {
	mfs, ok := fsys.(Mutable)
	if !ok {
		return &fs.PathError{Op: "rename", Path: oldpath + "->" + newpath, Err: fs.ErrInvalid}
	}
	return mfs.Rename(oldpath, newpath)
}

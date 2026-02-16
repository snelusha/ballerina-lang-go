package projects

import (
	"io/fs"

	"ballerina-lang-go/context"
)

type Environment struct {
	fsys        fs.FS
	compilerCtx *context.CompilerContext
}

func NewEnvironment(fsys fs.FS, cx *context.CompilerContext) *Environment {
	return &Environment{
		fsys:        fsys,
		compilerCtx: cx,
	}
}

func (e *Environment) compilerContext() *context.CompilerContext {
	return e.compilerCtx
}

func (e *Environment) fs() fs.FS {
	return e.fsys
}

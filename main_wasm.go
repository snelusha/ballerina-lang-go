//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"syscall/js"

	"ballerina-lang-go/common/bfs"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/runtime"

	_ "ballerina-lang-go/lib/io/runtime"
)

// var memFsys *bfs.MemFS
var memFsys fs.FS

func main() {
	memFsys = bfs.NewMemFS()

	c := make(chan struct{}, 0)

	js.Global().Set("updateFile", js.FuncOf(updateFile))
	js.Global().Set("runProject", js.FuncOf(runProject))
	js.Global().Set("readDir", js.FuncOf(readDir))
	js.Global().Set("readFile", js.FuncOf(readFile))
	js.Global().Set("mkdir", js.FuncOf(mkdir))

	fmt.Println("Ballerina WASM initialized!")

	<-c
}

func updateFile(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return "Invalid arguments: expected path and content"
	}

	path := args[0].String()
	content := args[1].String()

	err := bfs.WriteFile(memFsys, path, []byte(content), 0o644)
	if err != nil {
		return fmt.Sprintf("Error writing file: %v", err)
	}

	// err := memFsys.WriteFile(path, []byte(content), 0o644)
	// if err != nil {
	// 	return fmt.Sprintf("Error writing file: %v", err)
	// }

	return nil
}

func runProject(this js.Value, args []js.Value) (returnResult any) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v\n", r)
		}
	}()

	if len(args) < 1 {
		// return "Invalid arguments: expected project path"
		return nil
	}

	path := args[0].String()

	buildOpts := projects.NewBuildOptionsBuilder().
		Build()

	result, err := directory.LoadProject(memFsys, path, directory.ProjectLoadConfig{
		BuildOptions: &buildOpts,
	})
	if err != nil {
		// return fmt.Sprintf("Error loading project: %v", err)
		return nil
	}

	diagResult := result.Diagnostics()
	if diagResult.HasErrors() {
		// return fmt.Sprintf("Project loading errors: %v", diagResult.Diagnostics())
		return nil
	}

	project := result.Project()
	pkg := project.CurrentPackage()

	compilation := pkg.Compilation()
	compilationDiags := compilation.DiagnosticResult()
	if compilationDiags.HasErrors() {
		return nil
		// return fmt.Sprintf("Compilation errors: %v", compilationDiags.Diagnostics())
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	if birPkg == nil {
		return nil
		// return "BIR generation failed"
	}

	rt := runtime.NewRuntime()
	if err := rt.Interpret(*birPkg); err != nil {
		return fmt.Sprintf("Runtime error: %v", err)
	}

	return nil
}

func readDir(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "Invalid arguments: expected path"
	}
	path := args[0].String()

	// entries, err := memFsys.ReadDir(path)
	entries, err := fs.ReadDir(memFsys, path)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}

	type DirEntry struct {
		Name  string `json:"name"`
		IsDir bool   `json:"isDir"`
	}

	var result []DirEntry
	for _, entry := range entries {
		result = append(result, DirEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		})
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("Error marshalling entries: %v", err)
	}

	return string(jsonBytes)
}

func readFile(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "Invalid arguments: expected path"
	}
	path := args[0].String()

	f, err := memFsys.Open(path)
	if err != nil {
		return fmt.Sprintf("Error opening file: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}

	return string(content)
}

func mkdir(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "Invalid arguments: expected path"
	}
	path := args[0].String()

	// err := memFsys.MkdirAll(path, 0o755)
	err := bfs.MkdirAll(memFsys, path, 0o755)
	if err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}

	return nil
}

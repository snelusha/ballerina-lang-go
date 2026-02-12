//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"syscall/js"

	"ballerina-lang-go/common/bfs"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/runtime"
)

var memFsys *bfs.MemFS

func main() {
	// Initialize the in-memory file system
	memFsys = bfs.NewMemFS()

	c := make(chan struct{}, 0)

	// Expose Go functions to JavaScript
	js.Global().Set("updateFile", js.FuncOf(updateFile))
	js.Global().Set("runProject", js.FuncOf(runProject))
	js.Global().Set("readDir", js.FuncOf(readDir))
	js.Global().Set("readFile", js.FuncOf(readFile))
	js.Global().Set("mkdir", js.FuncOf(mkdir))

	fmt.Println("Wasm module initialized. Ready to run Ballerina projects.")
	<-c
}

// updateFile adds or updates a file in the in-memory file system.
// Arguments:
// 0: path (string) - The path of the file to create/update
// 1: content (string) - The content of the file
func updateFile(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return "Invalid arguments: expected path and content"
	}

	path := args[0].String()
	content := args[1].String()

	err := memFsys.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Sprintf("Error writing file: %v", err)
	}

	return nil
}

// runProject runs the Ballerina project from the in-memory file system.
// Arguments:
// 0: path (string) - The path to the project root or file to run
func runProject(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "Invalid arguments: expected project path"
	}

	path := args[0].String()

	buildOpts := projects.NewBuildOptionsBuilder().
		Build()

	// Load the project using the in-memory file system
	result, err := directory.LoadProject(memFsys, path, directory.ProjectLoadConfig{
		BuildOptions: &buildOpts,
	})
	if err != nil {
		return fmt.Sprintf("Error loading project: %v", err)
	}

	// Check for loading errors
	diagResult := result.Diagnostics()
	if diagResult.HasErrors() {
		// return diagnostics as string for now
		return fmt.Sprintf("Project loading errors: %v", diagResult.Diagnostics())
	}

	project := result.Project()
	pkg := project.CurrentPackage()

	// Compilation
	compilation := pkg.Compilation()
	compilationDiags := compilation.DiagnosticResult()
	if compilationDiags.HasErrors() {
		return fmt.Sprintf("Compilation errors: %v", compilationDiags.Diagnostics())
	}

	// Code generation / Interpretation
	// Assuming interpretation for now as per main.go logic
	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	if birPkg == nil {
		return "BIR generation failed"
	}

	// Runtime execution
	rt := runtime.NewRuntime()
	if err := rt.Interpret(*birPkg); err != nil {
		return fmt.Sprintf("Runtime error: %v", err)
	}

	return "Project executed successfully!"
}

// readDir returns a JSON string of directory entries.
// Arguments:
// 0: path (string) - The path of the directory to read
func readDir(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "Invalid arguments: expected path"
	}
	path := args[0].String()

	entries, err := memFsys.ReadDir(path)
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

// readFile returns the content of a file.
// Arguments:
// 0: path (string) - The path of the file to read
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

// mkdir creates a directory.
// Arguments:
// 0: path (string) - The path of the directory to create
func mkdir(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "Invalid arguments: expected path"
	}
	path := args[0].String()

	err := memFsys.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}

	return nil
}

//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"syscall/js"

	"ballerina-lang-go/common/bfs"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/tools/diagnostics"

	_ "ballerina-lang-go/lib/rt"
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
		printDiagnostics(memFsys, os.Stderr, diagResult, false)
		return
	}

	project := result.Project()
	pkg := project.CurrentPackage()

	compilation := pkg.Compilation()
	compilationDiags := compilation.DiagnosticResult()
	if compilationDiags.HasErrors() {
		printDiagnostics(memFsys, os.Stderr, compilationDiags, false)
		return
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

type outputStyle struct {
	reset, red, yellow, cyan, bold string
}

func (s outputStyle) severityColor(severity diagnostics.DiagnosticSeverity) string {
	if severity == diagnostics.Warning {
		return s.yellow
	}
	return s.red
}

func outputStyleFor(noColors bool) outputStyle {
	s := outputStyle{}
	if !noColors {
		s.reset = "\033[0m"
		s.red = "\033[31m"
		s.yellow = "\033[33m"
		s.cyan = "\033[36m"
		s.bold = "\033[1m"
	}
	return s
}

type diagnosticLocation struct {
	filePath            string
	startLine, startCol int
	endLine, endCol     int
	lineNumStr          string
	numWidth            int
}

func buildDiagnosticLocation(filePath string, startLine, startCol, endLine, endCol int) diagnosticLocation {
	lineNumStr := fmt.Sprintf("%d", startLine+1)
	return diagnosticLocation{
		filePath:   filePath,
		startLine:  startLine,
		startCol:   startCol,
		endLine:    endLine,
		endCol:     endCol,
		lineNumStr: lineNumStr,
		numWidth:   len(lineNumStr),
	}
}

func printDiagnostics(fsys fs.FS, w io.Writer, diagResult projects.DiagnosticResult, noColors bool) {
	for _, d := range diagResult.Diagnostics() {
		printDiagnostic(fsys, w, d, noColors)
	}
}

func printDiagnostic(fsys fs.FS, w io.Writer, d diagnostics.Diagnostic, noColors bool) {
	s := outputStyleFor(noColors)
	printDiagnosticHeader(w, s, d)

	location := d.Location()
	if location == nil {
		fmt.Fprintln(w)
		return
	}

	lineRange := location.LineRange()
	loc := buildDiagnosticLocation(
		lineRange.FileName(),
		lineRange.StartLine().Line(), lineRange.StartLine().Offset(),
		lineRange.EndLine().Line(), lineRange.EndLine().Offset(),
	)
	printDiagnosticLocation(w, s, loc)
	printSourceSnippet(w, s, loc, fsys, s.severityColor(d.DiagnosticInfo().Severity()))
	fmt.Fprintln(w)
}

func printDiagnosticHeader(w io.Writer, s outputStyle, d diagnostics.Diagnostic) {
	info := d.DiagnosticInfo()
	codeStr := ""
	if c := info.Code(); c != "" {
		codeStr = fmt.Sprintf("[%s]", c)
	}
	fmt.Fprintf(w, "%s%s%s%s%s: %s%s%s\n",
		s.bold, s.severityColor(info.Severity()), strings.ToLower(info.Severity().String()), codeStr, s.reset,
		s.bold, d.Message(), s.reset,
	)
}

func printDiagnosticLocation(w io.Writer, s outputStyle, loc diagnosticLocation) {
	fmt.Fprintf(w, "%*s%s-->%s %s:%d:%d\n",
		loc.numWidth, "", s.cyan, s.reset, loc.filePath, loc.startLine+1, loc.startCol+1,
	)
	if loc.filePath != "" {
		fmt.Fprintf(w, "%*s %s|%s\n", loc.numWidth, "", s.cyan, s.reset)
	}
}

func printSourceSnippet(w io.Writer, s outputStyle, loc diagnosticLocation, fsys fs.FS, severityColor string) {
	content, err := fs.ReadFile(fsys, loc.filePath)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	if loc.startLine >= len(lines) {
		return
	}
	lineContent := lines[loc.startLine]
	fmt.Fprintf(w, "%s%s |%s %s\n", s.cyan, loc.lineNumStr, s.reset, lineContent)
	highlightLen := loc.endCol - loc.startCol
	if loc.startLine != loc.endLine {
		highlightLen = len(lineContent) - loc.startCol
	}
	if highlightLen < 1 {
		highlightLen = 1
	}
	pointer := buildPointer(lineContent, loc.startCol, highlightLen)
	fmt.Fprintf(w, "%*s %s| %s%s%s\n", loc.numWidth, "", s.cyan, severityColor, pointer, s.reset)
}

func buildPointer(lineContent string, startCol, highlightLen int) string {
	var b strings.Builder
	for i := 0; i < startCol && i < len(lineContent); i++ {
		if lineContent[i] == '\t' {
			b.WriteByte('\t')
		} else {
			b.WriteByte(' ')
		}
	}
	for range highlightLen {
		b.WriteByte('^')
	}
	return b.String()
}

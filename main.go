//go:build !js || !wasm

package main

import (
	"fmt"
	"os"
	"strings"

	"ballerina-lang-go/bir"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/projects/directory"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/tools/diagnostics"
)

var runOpts struct {
	dumpTokens    bool
	dumpST        bool
	dumpAST       bool
	dumpCFG       bool
	dumpBIR       bool
	traceRecovery bool
	logFile       string
	format        string // Output format (dot, etc.)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: go run main.go <source-file>")
	}

	path := args[0]

	buildOpts := projects.NewBuildOptionsBuilder().
		WithDumpAST(runOpts.dumpAST).
		WithDumpBIR(runOpts.dumpBIR).
		WithDumpCFG(runOpts.dumpCFG).
		WithDumpCFGFormat(projects.ParseCFGFormat(runOpts.format)).
		WithDumpTokens(runOpts.dumpTokens).
		WithDumpST(runOpts.dumpST).
		WithTraceRecovery(runOpts.traceRecovery).
		Build()

	// Create a file system that points to the root of the OS.
	// This allows us to use absolute paths as relative paths to this FS.
	// Note: On Windows, this would need a different approach (like multiple drives).
	// For now assuming Unix-like environment for "bal run".
	fsys := os.DirFS(".")

	// When using os.DirFS("/"), absolute paths like /a/b/c become a/b/c
	// We need to ensure the path passed to LoadProject matches this expectation
	// or relies on LoadProject to handle cleaning.
	// However, standard os.DirFS usually expects relative paths.
	// Let's pass the OS root FS and let the path be handled.
	// Actually, if we use os.DirFS("/"), then "/Users/..." becomes "Users/...".
	// It's safer to pass the CWD or just use the path as is if we can mock root.

	// A better approach for the CLI tool which has real OS access is to use a
	// wrapper around os functions that implements fs.FS but allows absolute paths,
	// OR just use os.DirFS("/") and strip the leading slash from absolute paths.

	// Let's try os.DirFS("/") and trimming leading slash if needed.
	path = strings.TrimPrefix(path, "/")

	result, err := directory.LoadProject(fsys, path, directory.ProjectLoadConfig{
		BuildOptions: &buildOpts,
	})
	if err != nil {
		fmt.Printf("Error loading project: %v\n", err)
		return
	}

	// Check for loading errors
	diagResult := result.Diagnostics()
	if diagResult.HasErrors() {
		printDiagnostics(diagResult)
		fmt.Println("project loading contains errors")
		return
	}

	project := result.Project()
	pkg := project.CurrentPackage()

	// Get package compilation (triggers parsing, type checking, semantic analysis, CFG analysis)
	compilation := pkg.Compilation()

	// Check for compilation errors
	compilationDiags := compilation.DiagnosticResult()
	if compilationDiags.HasErrors() {
		printDiagnostics(compilationDiags)
		fmt.Println("compilation failed with errors")
		return
	}

	// Create backend and generate BIR
	backend := projects.NewBallerinaBackend(compilation)
	birPkg := backend.BIR()

	if birPkg == nil {
		fmt.Println("BIR generation failed: no BIR package produced")
		return
	}

	// Dump BIR if requested
	if buildOpts.DumpBIR() {
		prettyPrinter := bir.PrettyPrinter{}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "==================BEGIN BIR==================")
		fmt.Println(strings.TrimSpace(prettyPrinter.Print(*birPkg)))
		fmt.Fprintln(os.Stderr, "===================END BIR===================")
	}

	rt := runtime.NewRuntime()
	if err := rt.Interpret(*birPkg); err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		return
	}
}

func printDiagnostics(diagResult projects.DiagnosticResult) {
	for _, d := range diagResult.Diagnostics() {
		fmt.Fprintln(os.Stderr, formatDiagnostic(d))
	}
}

func formatDiagnostic(d diagnostics.Diagnostic) string {
	loc := d.Location()
	info := d.DiagnosticInfo()

	// Format: filepath:line:col: severity: message
	if loc != nil {
		lineRange := loc.LineRange()
		return fmt.Sprintf("%s:%d:%d: %s: %s",
			lineRange.FileName(),
			lineRange.StartLine().Line(),
			lineRange.StartLine().Offset(),
			info.Severity().String(),
			d.Message())
	}
	return fmt.Sprintf("%s: %s", info.Severity().String(), d.Message())
}

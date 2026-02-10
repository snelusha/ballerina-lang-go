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

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/context"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/tools/diagnostics"

	"github.com/spf13/cobra"
)

var runOpts struct {
	dumpTokens    bool
	dumpST        bool
	dumpAST       bool
	dumpBIR       bool
	traceRecovery bool
	logFile       string
}

var runCmd = &cobra.Command{
	Use:   "run <source-file.bal>",
	Short: "Compile and run the current package or a Ballerina source file",
	Long: `	Compile the current package and run it.

	The 'run' command compiles and executes the given Ballerina source file.

	A Ballerina program consists of one or more modules; one of these modules
	is distinguished as the root module, which is the default module of
	current package.

	Ballerina program execution consists of two consecutive phases.
	The initialization phase initializes all modules of a program one after
	another. If a module defines a function named 'init()', it will be
	invoked during this phase. If the root module of the program defines a
	public function named 'main()', then it will be invoked.

	If the initialization phase of program execution completes successfully,
	then execution proceeds to the listening phase. If there are no module
	listeners, then the listening phase immediately terminates successfully.
	Otherwise, the listening phase initializes the module listeners.

	A service declaration is the syntactic sugar for creating a service object
	and attaching it to the module listener specified in the service
	declaration.`,
	Args: validateSourceFile,
	RunE: runBallerina,
}

func init() {
	runCmd.Flags().BoolVar(&runOpts.dumpTokens, "dump-tokens", false, "Dump lexer tokens")
	runCmd.Flags().BoolVar(&runOpts.dumpST, "dump-st", false, "Dump syntax tree")
	runCmd.Flags().BoolVar(&runOpts.dumpAST, "dump-ast", false, "Dump abstract syntax tree")
	runCmd.Flags().BoolVar(&runOpts.dumpBIR, "dump-bir", false, "Dump Ballerina Intermediate Representation")
	runCmd.Flags().BoolVar(&runOpts.traceRecovery, "trace-recovery", false, "Enable error recovery tracing")
	runCmd.Flags().StringVar(&runOpts.logFile, "log-file", "", "Write debug output to specified file")
}

func runBallerina(cmd *cobra.Command, args []string) error {
	fileName := args[0]

	var debugCtx *debugcommon.DebugContext
	var wg sync.WaitGroup
	flags := uint16(0)

	if runOpts.dumpTokens {
		flags |= debugcommon.DUMP_TOKENS
	}
	if runOpts.dumpST {
		flags |= debugcommon.DUMP_ST
	}
	if runOpts.traceRecovery {
		flags |= debugcommon.DEBUG_ERROR_RECOVERY
	}

	if flags != 0 {
		debugcommon.Init(flags)
		debugCtx = &debugcommon.DebugCtx

		var logWriter *os.File
		var err error
		if runOpts.logFile != "" {
			logWriter, err = os.Create(runOpts.logFile)
			if err != nil {
				cmdErr := fmt.Errorf("error creating log file %s: %w", runOpts.logFile, err)
				printError(cmdErr, "", false)
				return cmdErr
			}
		} else {
			logWriter = os.Stderr
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if runOpts.logFile != "" {
				defer logWriter.Close()
			}
			for msg := range debugCtx.Channel {
				fmt.Fprintf(logWriter, "%s\n", msg)
			}
		}()

		// Ensure debug context cleanup on any exit path
		defer func() {
			if debugCtx != nil {
				close(debugCtx.Channel)
				wg.Wait()
			}
		}()
	}

	// Compile the source
	fmt.Fprintln(os.Stderr, "Compiling source")
	fmt.Fprintf(os.Stderr, "\t%s\n", filepath.Base(fileName))

	cx := context.NewCompilerContext(semtypes.CreateTypeEnv())

	syntaxTree, err := parser.GetSyntaxTree(debugCtx, fileName)
	if err != nil {
		printError(fmt.Errorf("compilation failed: %w", err), "", false)
		return fmt.Errorf("compilation failed: %w", err)
	}

	if cx.HasErrors() {
		printErrors(cx)
		return nil
	}

	compilationUnit := ast.GetCompilationUnit(cx, syntaxTree)
	if runOpts.dumpAST {
		prettyPrinter := ast.PrettyPrinter{}
		fmt.Println(prettyPrinter.Print(compilationUnit))
	}

	if cx.HasErrors() {
		printErrors(cx)
		return nil
	}

	pkg := ast.ToPackage(compilationUnit)

	if cx.HasErrors() {
		printErrors(cx)
		return nil
	}

	// Resolve symbols (imports) before type resolution
	importedSymbols := semantics.ResolveImports(cx, pkg)
	semantics.ResolveSymbols(cx, pkg, importedSymbols)
	// Add type resolution step
	typeResolver := semantics.NewTypeResolver(cx)
	typeResolver.ResolveTypes(cx, pkg)

	// if cx.HasErrors() {
	// 	printErrors(cx)
	// 	return nil
	// }

	// Run semantic analysis after type resolution
	semanticAnalyzer := semantics.NewSemanticAnalyzer(cx)
	semanticAnalyzer.Analyze(pkg)

	if cx.HasErrors() {
		printErrors(cx)
		return nil
	}

	birPkg := bir.GenBir(cx, pkg)
	if runOpts.dumpBIR {
		prettyPrinter := bir.PrettyPrinter{}
		// Print the BIR with separators
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "==================BEGIN BIR==================")
		fmt.Println(strings.TrimSpace(prettyPrinter.Print(*birPkg)))
		fmt.Fprintln(os.Stderr, "===================END BIR===================")
	}

	// Run the executable
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Running executable")
	fmt.Fprintln(os.Stderr)

	rt := runtime.NewRuntime()
	if err := rt.Interpret(*birPkg); err != nil {
		return err
	}
	return nil
}

func printErrors(cx *context.CompilerContext) {
	if cx.HasErrors() {
		fmt.Fprintln(os.Stderr, "\nCompilation failed with the following errors:")
		for _, diagnostic := range cx.GetDiagnostics() {
			printDiagnostic(diagnostic)
		}
	}
}

func printDiagnostic(d diagnostics.Diagnostic) {
	location := d.Location()
	lineRange := location.LineRange()
	fileName := lineRange.FileName()
	startLine := lineRange.StartLine().Line()
	startCol := lineRange.StartLine().Offset()

	// Color codes
	reset := "\033[0m"
	red := "\033[31m"
	yellow := "\033[33m"
	cyan := "\033[36m"
	bold := "\033[1m"

	severity := d.DiagnosticInfo().Severity()
	severityStr := strings.ToLower(severity.String())
	severityColor := red
	if severity == diagnostics.Warning {
		severityColor = yellow
	}

	code := d.DiagnosticInfo().Code()
	codeStr := ""
	if code != "" {
		codeStr = fmt.Sprintf("[%s]", code)
	}

	// 1. severity[CODE]: MESSAGE
	fmt.Fprintf(os.Stderr, "%s%s%s%s%s: %s%s%s\n",
		bold, severityColor, severityStr, codeStr, reset,
		bold, d.Message(), reset,
	)

	lineNumStr := fmt.Sprintf("%d", startLine+1)
	numWidth := len(lineNumStr)

	// 2.  --> FILE:LINE:COL
	fmt.Fprintf(os.Stderr, "%*s%s-->%s %s:%d:%d\n",
		numWidth, "", cyan, reset, fileName, startLine+1, startCol+1,
	)

	// Print source snippet if available
	if fileName != "" {
		file, err := os.Open(fileName)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			currentLine := 0

			fmt.Fprintf(os.Stderr, "%*s %s|%s\n", numWidth, "", cyan, reset)

			for scanner.Scan() {
				if currentLine == startLine {
					lineContent := scanner.Text()

					// 3. LINE | CONTENT
					fmt.Fprintf(os.Stderr, "%s%s %s| %s\n", cyan, lineNumStr, reset, lineContent)

					// 4.   | POINTER
					pointer := ""
					for i := range startCol {
						if len(lineContent) > i && lineContent[i] == '\t' {
							pointer += "\t"
						} else {
							pointer += " "
						}
					}
					pointer += "^"
					fmt.Fprintf(os.Stderr, "%*s %s| %s%s%s\n", numWidth, "", cyan, severityColor, pointer, reset)
					break
				}
				currentLine++
			}
		}
	}
	fmt.Fprintln(os.Stderr)
}

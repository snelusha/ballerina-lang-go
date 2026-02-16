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

package projects

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"ballerina-lang-go/tools/diagnostics"

	"golang.org/x/term"
)

// PrintDiagnostics prints all diagnostics from a DiagnosticResult to the given writer.
func PrintDiagnostics(w io.Writer, diagResult DiagnosticResult) {
	for _, d := range diagResult.Diagnostics() {
		PrintDiagnostic(w, d)
	}
}

// PrintDiagnostic prints a single diagnostic for CLI output.
func PrintDiagnostic(w io.Writer, d diagnostics.Diagnostic) {
	useColor := isTerminal(w)

	reset := ""
	red := ""
	yellow := ""
	cyan := ""
	bold := ""

	if useColor {
		reset = "\033[0m"
		red = "\033[31m"
		yellow = "\033[33m"
		cyan = "\033[36m"
		bold = "\033[1m"
	}

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

	// severity[CODE]: MESSAGE
	fmt.Fprintf(w, "%s%s%s%s%s: %s%s%s\n",
		bold, severityColor, severityStr, codeStr, reset,
		bold, d.Message(), reset,
	)

	location := d.Location()
	if location == nil {
		fmt.Fprintln(w)
		return
	}

	lineRange := location.LineRange()
	fileName := lineRange.FileName()
	startLine := lineRange.StartLine().Line()
	startCol := lineRange.StartLine().Offset()

	lineNumStr := fmt.Sprintf("%d", startLine+1)
	numWidth := len(lineNumStr)

	// --> FILE:LINE:COL
	fmt.Fprintf(w, "%*s%s-->%s %s:%d:%d\n",
		numWidth, "", cyan, reset, fileName, startLine+1, startCol+1,
	)

	// Print source snippet if available
	if fileName != "" {
		fmt.Println(fileName)
		file, err := os.Open(fileName)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			currentLine := 0

			fmt.Fprintf(w, "%*s %s|%s\n", numWidth, "", cyan, reset)

			for scanner.Scan() {
				if currentLine == startLine {
					lineContent := scanner.Text()

					// LINE | CONTENT
					fmt.Fprintf(w, "%s%s %s| %s\n", cyan, lineNumStr, reset, lineContent)

					endLine := lineRange.EndLine().Line()
					endCol := lineRange.EndLine().Offset()

					// | POINTER
					pointer := ""
					for i := range startCol {
						if len(lineContent) > i && lineContent[i] == '\t' {
							pointer += "\t"
						} else {
							pointer += " "
						}
					}

					highlightLen := 1
					if startLine == endLine {
						highlightLen = endCol - startCol
					} else if startLine < endLine {
						highlightLen = len(lineContent) - startCol
					}
					if highlightLen < 1 {
						highlightLen = 1
					}

					for range highlightLen {
						pointer += "^"
					}
					fmt.Fprintf(w, "%*s %s| %s%s%s\n", numWidth, "", cyan, severityColor, pointer, reset)
					break
				}
				currentLine++
			}
		}
	}
	fmt.Fprintln(w)
}

func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

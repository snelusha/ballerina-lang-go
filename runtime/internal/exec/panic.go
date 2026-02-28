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

package exec

import (
	"fmt"

	"ballerina-lang-go/tools/diagnostics"
)

// RuntimePanic carries a runtime panic message and optional source location from BIR.
// When recovered, it is turned into a formatted diagnostic.
type RuntimePanic struct {
	Loc     diagnostics.Location
	Message string
}

func (p RuntimePanic) Error() string {
	return p.Message
}

// panicAt panics with a RuntimePanic so the executor can recover and format a diagnostic.
// loc may be nil (e.g. stack overflow); the diagnostic will then have no location.
func panicAt(loc diagnostics.Location, format string, args ...any) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	panic(RuntimePanic{Loc: loc, Message: msg})
}

// FormatRuntimeDiagnostic builds a diagnostic from a RuntimePanic and returns its formatted string.
func FormatRuntimeDiagnostic(p RuntimePanic) string {
	info := diagnostics.NewDiagnosticInfo(nil, "%s", diagnostics.Error)
	loc := p.Loc
	if loc == nil {
		loc = diagnostics.NewBLangDiagnosticLocation("", 0, 0, 0, 0, 0, 0)
	}
	d := diagnostics.CreateDiagnostic(info, loc, p.Message)
	return d.String()
}

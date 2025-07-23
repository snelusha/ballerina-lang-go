package tree

import (
	"ballerina-lang-go/compiler/internal/diagnostics"
	"ballerina-lang-go/tools/diagnostics"
)

// STNodeDiagnostic defines the interface for internal representation of diagnostic
// that is related to an internal syntax node.
type STNodeDiagnostic interface {
	diagnostics.IRDiagnostic
	DiagnosticCode() diagnostics.DiagnosticCode
	Args() []interface{}
}

// stNodeDiagnosticImpl is the concrete implementation of STNodeDiagnostic.
type stNodeDiagnosticImpl struct {
	diagnosticCode diagnostics.DiagnosticCode
	args           []interface{}
}

// NewSTNodeDiagnostic constructs a new STNodeDiagnostic with the given diagnostic code and arguments.
func NewSTNodeDiagnostic(diagnosticCode diagnostics.DiagnosticCode, args ...interface{}) STNodeDiagnostic {
	return &stNodeDiagnosticImpl{
		diagnosticCode: diagnosticCode,
		args:           args,
	}
}

// STNodeDiagnosticFromCode creates a new STNodeDiagnostic from the given diagnostic code and arguments.
func STNodeDiagnosticFromCode(diagnosticCode diagnostics.DiagnosticCode, args ...interface{}) STNodeDiagnostic {
	return NewSTNodeDiagnostic(diagnosticCode, args...)
}

// DiagnosticCode returns the diagnostic code.
func (s stNodeDiagnosticImpl) DiagnosticCode() diagnostics.DiagnosticCode {
	return s.diagnosticCode
}

// Args returns the arguments array.
func (s stNodeDiagnosticImpl) Args() []interface{} {
	return s.args
}

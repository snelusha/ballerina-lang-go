package tree

import (
	common "ballerina-lang-go/compiler"
	"ballerina-lang-go/compiler/internal/diagnostics"
	diagnosticsTools "ballerina-lang-go/tools/diagnostics"
)

// stNodeDiagnosticImpl is the concrete implementation of STNodeDiagnostic.
type stNodeDiagnosticImpl struct {
	diagnostics.IRDiagnosticImpl
	diagnosticCode diagnosticsTools.DiagnosticCode
	args           []interface{}
}

// NewSTNodeDiagnostic constructs an STNodeDiagnostic with the given diagnostic code and arguments.
func NewSTNodeDiagnostic(diagnosticCode diagnosticsTools.DiagnosticCode, args ...interface{}) common.STNodeDiagnostic {
	return &stNodeDiagnosticImpl{
		diagnosticCode: diagnosticCode,
		args:           args,
	}
}

// FromDiagnosticCode creates an STNodeDiagnostic from the given diagnostic code and arguments.
func FromDiagnosticCode(diagnosticCode diagnosticsTools.DiagnosticCode, args ...interface{}) common.STNodeDiagnostic {
	return NewSTNodeDiagnostic(diagnosticCode, args...)
}

// DiagnosticCode returns the diagnostic code.
func (d stNodeDiagnosticImpl) DiagnosticCode() diagnosticsTools.DiagnosticCode {
	return d.diagnosticCode
}

// Args returns the arguments associated with this diagnostic.
func (d stNodeDiagnosticImpl) Args() []interface{} {
	return d.args
}

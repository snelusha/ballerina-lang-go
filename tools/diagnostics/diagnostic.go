package diagnostics

import "fmt"

// Diagnostic represents a diagnostic message (error, warning, etc.) with location information.
// A diagnostic represents a compiler error, a warning or a message at a specific location in the source file.
type Diagnostic interface {
	Location() Location
	DiagnosticInfo() DiagnosticInfo
	Message() string
	Properties() []DiagnosticProperty[any]
	String() string
}

type diagnosticBase struct{}

// String returns a string representation of the diagnostic.
// This is the default implementation from the abstract Diagnostic class.
func (db diagnosticBase) String(d Diagnostic) string {
	var location string
	if d.Location().LineRange().FileName() == "" {
		location = ""
	} else {
		location = fmt.Sprintf(" [%s:%s]",
			d.Location().LineRange().FileName(),
			d.Location().LineRange().String())
	}
	return fmt.Sprintf("%s%s %s",
		d.DiagnosticInfo().Severity().String(),
		location,
		d.Message())
}

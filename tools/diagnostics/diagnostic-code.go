package diagnostics

// DiagnosticCode represents a diagnostic code.
// Diagnostic code uniquely identifies a diagnostic.
type DiagnosticCode interface {
	Severity() DiagnosticSeverity
	DiagnosticId() string
	MessageKey() string
}

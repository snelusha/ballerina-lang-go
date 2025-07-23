package diagnostics

// DiagnosticSeverity represents a severity of a Diagnostic.
type DiagnosticSeverity uint8

const (
	// Internal represents an internal diagnostic severity.
	Internal DiagnosticSeverity = iota
	// Hint represents a hint diagnostic severity.
	Hint
	// Info represents an info diagnostic severity.
	Info
	// Warning represents a warning diagnostic severity.
	Warning
	// Error represents an error diagnostic severity.
	Error
)

// String returns the string representation of the diagnostic severity.
func (ds DiagnosticSeverity) String() string {
	switch ds {
	case Internal:
		return "INTERNAL"
	case Hint:
		return "HINT"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

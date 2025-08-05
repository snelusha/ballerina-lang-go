package diagnostics

// DiagnosticSeverity represents a severity of a Diagnostic.
type DiagnosticSeverity uint8

const (
	Internal DiagnosticSeverity = iota
	Hint
	Info
	Warning
	Error
)

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

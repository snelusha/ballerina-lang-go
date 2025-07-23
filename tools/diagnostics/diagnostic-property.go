package diagnostics

// DiagnosticProperty represents properties passed when diagnostic logging.
type DiagnosticProperty[T any] interface {
	Kind() DiagnosticPropertyKind
	Value() T
}

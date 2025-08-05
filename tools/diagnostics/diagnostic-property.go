package diagnostics

// DiagnosticProperty represents properties passed when diagnostic logging.
type DiagnosticProperty[T interface{}] interface {
	Kind() DiagnosticPropertyKind
	Value() T
}

package diagnostics

// DiagnosticPropertyKind represents the kind of the diagnostic property.
type DiagnosticPropertyKind uint8

const (
	Symbolic DiagnosticPropertyKind = iota
	String
	Numeric
	Collection
	Other
)

// String returns the string representation of the diagnostic property kind.
func (dpk DiagnosticPropertyKind) String() string {
	switch dpk {
	case Symbolic:
		return "SYMBOLIC"
	case String:
		return "STRING"
	case Numeric:
		return "NUMERIC"
	case Collection:
		return "COLLECTION"
	case Other:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

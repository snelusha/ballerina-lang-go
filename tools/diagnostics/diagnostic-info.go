package diagnostics

// DiagnosticInfo represents an abstract shape of a Diagnostic that is independent of
// the location and message arguments.
type DiagnosticInfo interface {
	Code() string
	MessageFormat() string
	Severity() DiagnosticSeverity
	DiagnosticInfoLookupKey() DiagnosticInfoLookupKey
}

// DiagnosticInfoLookupKey represents the comparable fields of DiagnosticInfo for equality/hashing.
type DiagnosticInfoLookupKey struct {
	Code          *string // pointer to handle nil values
	MessageFormat string
	Severity      DiagnosticSeverity
}

type diagnosticInfoImpl struct {
	code          *string // pointer to handle nil values
	messageFormat string
	severity      DiagnosticSeverity
}

// NewDiagnosticInfo constructs an abstract shape of a Diagnostic.
//
// Parameters:
//   - code: a code that can be used to uniquely identify a diagnostic category
//   - messageFormat: a pattern that can be formatted with message formatting utilities
//   - severity: the severity of the diagnostic
func NewDiagnosticInfo(code *string, messageFormat string, severity DiagnosticSeverity) DiagnosticInfo {
	return &diagnosticInfoImpl{
		code:          code,
		messageFormat: messageFormat,
		severity:      severity,
	}
}

func (di diagnosticInfoImpl) Code() string {
	if di.code == nil {
		return ""
	}
	return *di.code
}

func (di diagnosticInfoImpl) MessageFormat() string {
	return di.messageFormat
}

func (di diagnosticInfoImpl) Severity() DiagnosticSeverity {
	return di.severity
}

// DiagnosticInfoLookupKey returns the lookup key for equality comparisons.
func (di diagnosticInfoImpl) DiagnosticInfoLookupKey() DiagnosticInfoLookupKey {
	return DiagnosticInfoLookupKey{
		Code:          di.code,
		MessageFormat: di.messageFormat,
		Severity:      di.severity,
	}
}

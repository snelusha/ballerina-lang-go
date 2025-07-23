package diagnostics

// CreateDiagnostic creates a Diagnostic instance from the given details.
//
// Parameters:
//   - diagnosticInfo: static diagnostic information
//   - location: the location of the diagnostic
//   - args: arguments to diagnostic message format
//
// Returns a Diagnostic instance.
func CreateDiagnostic(diagnosticInfo DiagnosticInfo, location Location, args ...interface{}) Diagnostic {
	return NewDefaultDiagnostic(diagnosticInfo, location, []DiagnosticProperty[any]{}, args...)
}

// CreateDiagnosticWithProperties creates a Diagnostic instance from the given details.
//
// Parameters:
//   - diagnosticInfo: static diagnostic information
//   - location: the location of the diagnostic
//   - properties: properties associated with the diagnostic
//   - args: arguments to diagnostic message format
//
// Returns a Diagnostic instance.
func CreateDiagnosticWithProperties(diagnosticInfo DiagnosticInfo, location Location, properties []DiagnosticProperty[any], args ...interface{}) Diagnostic {
	return NewDefaultDiagnostic(diagnosticInfo, location, properties, args...)
}

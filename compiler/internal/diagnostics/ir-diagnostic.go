package diagnostics

// IRDiagnostic defines the interface for internal representation of a Diagnostic.
type IRDiagnostic interface {
	// This interface is currently empty as the Java class has no methods
	// Methods will be added here as the Java class evolves
}

// irDiagnosticImpl is the concrete implementation of IRDiagnostic.
type irDiagnosticImpl struct {
	// Fields will be added here as the Java class evolves
}

// NewIRDiagnostic constructs a new IRDiagnostic.
func NewIRDiagnostic() IRDiagnostic {
	return &irDiagnosticImpl{}
}

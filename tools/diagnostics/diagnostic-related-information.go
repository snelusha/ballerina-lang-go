package diagnostics

// DiagnosticRelatedInformation represents a message and location related to a particular Diagnostic.
// A sample usage would be to record all symbol information related to duplicate symbol error.
type DiagnosticRelatedInformation interface {
	Location() Location
	Message() string
}

// diagnosticRelatedInformationImpl is the concrete implementation of DiagnosticRelatedInformation.
type diagnosticRelatedInformationImpl struct {
	location Location
	message  string
}

// NewDiagnosticRelatedInformation constructs a DiagnosticRelatedInformation with the given location and message.
func NewDiagnosticRelatedInformation(location Location, message string) DiagnosticRelatedInformation {
	return &diagnosticRelatedInformationImpl{
		location: location,
		message:  message,
	}
}

// Location returns the location of the related information.
func (dri diagnosticRelatedInformationImpl) Location() Location {
	return dri.location
}

// Message returns the message of the related information.
func (dri diagnosticRelatedInformationImpl) Message() string {
	return dri.message
}

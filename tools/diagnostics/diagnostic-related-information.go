package diagnostics

// DiagnosticRelatedInformation represents a message and location related to a particular Diagnostic.
// A sample usage would be to record all symbol information related to duplicate symbol error.
type DiagnosticRelatedInformation interface {
	Location() Location
	Message() string
}

type diagnosticRelatedInformationImpl struct {
	location Location
	message  string
}

func NewDiagnosticRelatedInformation(location Location, message string) DiagnosticRelatedInformation {
	return &diagnosticRelatedInformationImpl{
		location: location,
		message:  message,
	}
}

func (dri diagnosticRelatedInformationImpl) Location() Location {
	return dri.location
}

func (dri diagnosticRelatedInformationImpl) Message() string {
	return dri.message
}

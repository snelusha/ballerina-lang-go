package diagnostics

import (
	"ballerina-lang-go/tools/text"
	"fmt"
)

// DefaultDiagnostic is an internal implementation of the Diagnostic interface that is used by the DiagnosticFactory
// to create diagnostics.
type DefaultDiagnostic interface {
	Diagnostic
}

// defaultDiagnosticImpl is the concrete implementation of DefaultDiagnostic.
type defaultDiagnosticImpl struct {
	diagnosticBase
	diagnosticInfo DiagnosticInfo
	location       Location
	properties     []DiagnosticProperty[any]
	message        string
}

// NewDefaultDiagnostic constructs a DefaultDiagnostic with the given parameters.
func NewDefaultDiagnostic(diagnosticInfo DiagnosticInfo, location Location, properties []DiagnosticProperty[any], args ...interface{}) DefaultDiagnostic {
	message := formatMessage(diagnosticInfo.MessageFormat(), args...)
	return &defaultDiagnosticImpl{
		diagnosticBase: diagnosticBase{},
		diagnosticInfo: diagnosticInfo,
		location:       location,
		properties:     properties,
		message:        message,
	}
}

// Location returns the location of the diagnostic.
func (dd *defaultDiagnosticImpl) Location() Location {
	return dd.location
}

// DiagnosticInfo returns the diagnostic information.
func (dd *defaultDiagnosticImpl) DiagnosticInfo() DiagnosticInfo {
	return dd.diagnosticInfo
}

// Message returns the formatted diagnostic message.
func (dd *defaultDiagnosticImpl) Message() string {
	return dd.message
}

// Properties returns the diagnostic properties.
func (dd *defaultDiagnosticImpl) Properties() []DiagnosticProperty[any] {
	return dd.properties
}

// String returns a string representation of the diagnostic.
// This overrides the base implementation with custom formatting.
func (dd *defaultDiagnosticImpl) String() string {
	lineRange := dd.location.LineRange()
	filePath := lineRange.FileName()

	// Create one-based line range (convert from zero-based to one-based)
	startLine := lineRange.StartLine()
	endLine := lineRange.EndLine()

	oneBasedStartLine := text.LinePositionFromLineAndOffset(startLine.Line()+1, startLine.Offset()+1)
	oneBasedEndLine := text.LinePositionFromLineAndOffset(endLine.Line()+1, endLine.Offset()+1)
	oneBasedLineRange := text.LineRangeFromFileNameAndLinePositions(filePath, oneBasedStartLine, oneBasedEndLine)

	return fmt.Sprintf("%s [%s:%s] %s",
		dd.diagnosticInfo.Severity().String(),
		filePath,
		oneBasedLineRange.String(),
		dd.Message())
}

// formatMessage formats the message using the provided format string and arguments.
// This is a simplified version of Java's MessageFormat.format().
func formatMessage(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

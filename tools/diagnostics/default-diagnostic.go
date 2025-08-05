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

type defaultDiagnosticImpl struct {
	diagnosticBase
	diagnosticInfo DiagnosticInfo
	location       Location
	properties     []DiagnosticProperty[interface{}]
	message        string
}

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

func (dd *defaultDiagnosticImpl) Location() Location {
	return dd.location
}

func (dd *defaultDiagnosticImpl) DiagnosticInfo() DiagnosticInfo {
	return dd.diagnosticInfo
}

func (dd *defaultDiagnosticImpl) Message() string {
	return dd.message
}

func (dd *defaultDiagnosticImpl) Properties() []DiagnosticProperty[any] {
	return dd.properties
}

func (dd *defaultDiagnosticImpl) String() string {
	lineRange := dd.location.LineRange()
	filePath := lineRange.FileName()

	startLine := lineRange.StartLine()
	endLine := lineRange.EndLine()

	oneBasedStartLine := text.LinePositionFromLineAndOffset(startLine.Line()+1, startLine.Offset()+1)
	oneBasedEndLine := text.LinePositionFromLineAndOffset(endLine.Line()+1, endLine.Offset()+1)
	oneBasedLineRange := text.LineRangeFromLinePositions(filePath, oneBasedStartLine, oneBasedEndLine)

	return fmt.Sprintf("%s [%s:%s] %s",
		dd.diagnosticInfo.Severity().String(),
		filePath,
		oneBasedLineRange.String(),
		dd.Message())
}

func formatMessage(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

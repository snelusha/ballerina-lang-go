package elements

// MarkdownDocAttachment represents a markdown doc attachment
type MarkdownDocAttachment struct {
	Description             string
	Parameters              []*MarkdownDocParameter
	ReturnValueDescription  string
	DeprecatedDocumentation string
	DeprecatedParams        []*MarkdownDocParameter
}

// NewMarkdownDocAttachment creates a new MarkdownDocAttachment with the given parameter capacity
func NewMarkdownDocAttachment(paramCount int) *MarkdownDocAttachment {
	return &MarkdownDocAttachment{
		Parameters:       make([]*MarkdownDocParameter, 0, paramCount),
		DeprecatedParams: make([]*MarkdownDocParameter, 0),
	}
}

// MarkdownDocParameter represents a parameter in markdown documentation
type MarkdownDocParameter struct {
	Name        string
	Description string
}

// NewMarkdownDocParameter creates a new MarkdownDocParameter
func NewMarkdownDocParameter(name, description string) *MarkdownDocParameter {
	return &MarkdownDocParameter{
		Name:        name,
		Description: description,
	}
}

// GetName returns the parameter name
func (p *MarkdownDocParameter) GetName() string {
	return p.Name
}

// GetDescription returns the parameter description
func (p *MarkdownDocParameter) GetDescription() string {
	return p.Description
}

type AttachPoint struct {
	Point Point
}

type Point int

const (
	PointService Point = iota
	PointResource
	PointFunction
	PointObject
	PointRecord
	PointField
	PointParameter
	PointReturn
	PointType
	PointClass
	PointListener
	PointWorker
	PointConst
	PointAnnotation
	PointExternal
	PointVar
	PointLet
	PointSource
)

package elements

type MarkdownDocAttachment interface {
	GetDescription() string
	SetDescription(desc string)
	GetParameters() []*MarkdownParameter
	SetParameters(params []*MarkdownParameter)
	AddParameter(param *MarkdownParameter)
	GetReturnValueDescription() string
	SetReturnValueDescription(desc string)
	GetDeprecatedDocumentation() string
	SetDeprecatedDocumentation(doc string)
	GetDeprecatedParams() []*MarkdownParameter
	SetDeprecatedParams(params []*MarkdownParameter)
	AddDeprecatedParam(param *MarkdownParameter)
}

type markdownDocAttachmentImpl struct {
	description             string
	parameters              []*MarkdownParameter
	returnValueDescription  string
	deprecatedDocumentation string
	deprecatedParams        []*MarkdownParameter
}

func NewMarkdownDocAttachment(paramCount int) MarkdownDocAttachment {
	return &markdownDocAttachmentImpl{
		parameters:       make([]*MarkdownParameter, 0, paramCount),
		deprecatedParams: make([]*MarkdownParameter, 0),
	}
}

func (m *markdownDocAttachmentImpl) GetDescription() string {
	return m.description
}

func (m *markdownDocAttachmentImpl) SetDescription(desc string) {
	m.description = desc
}

func (m *markdownDocAttachmentImpl) GetParameters() []*MarkdownParameter {
	return m.parameters
}

func (m *markdownDocAttachmentImpl) SetParameters(params []*MarkdownParameter) {
	m.parameters = params
}

func (m *markdownDocAttachmentImpl) AddParameter(param *MarkdownParameter) {
	m.parameters = append(m.parameters, param)
}

func (m *markdownDocAttachmentImpl) GetReturnValueDescription() string {
	return m.returnValueDescription
}

func (m *markdownDocAttachmentImpl) SetReturnValueDescription(desc string) {
	m.returnValueDescription = desc
}

func (m *markdownDocAttachmentImpl) GetDeprecatedDocumentation() string {
	return m.deprecatedDocumentation
}

func (m *markdownDocAttachmentImpl) SetDeprecatedDocumentation(doc string) {
	m.deprecatedDocumentation = doc
}

func (m *markdownDocAttachmentImpl) GetDeprecatedParams() []*MarkdownParameter {
	return m.deprecatedParams
}

func (m *markdownDocAttachmentImpl) SetDeprecatedParams(params []*MarkdownParameter) {
	m.deprecatedParams = params
}

func (m *markdownDocAttachmentImpl) AddDeprecatedParam(param *MarkdownParameter) {
	m.deprecatedParams = append(m.deprecatedParams, param)
}

type MarkdownParameter struct {
	Name        string
	Description string
}

func NewMarkdownParameter(name, description string) *MarkdownParameter {
	return &MarkdownParameter{
		Name:        name,
		Description: description,
	}
}

func (p *MarkdownParameter) GetName() string {
	return p.Name
}

func (p *MarkdownParameter) GetDescription() string {
	return p.Description
}

package elements

type Parameter interface {
	GetName() string
	GetDescription() string
}

type parameterImpl struct {
	name        string
	description string
}

func NewParameter(name, description string) Parameter {
	return &parameterImpl{
		name:        name,
		description: description,
	}
}

func (p *parameterImpl) GetName() string {
	return p.name
}

func (p *parameterImpl) GetDescription() string {
	return p.description
}

type MarkdownDocAttachment interface {
	GetDescription() string
	GetParameters() []Parameter
	GetReturnValueDescription() string
	GetDeprecatedDocumentation() string
	GetDeprecatedParams() []Parameter
	SetDescription(description string)
	AddParameter(param Parameter)
	SetReturnValueDescription(description string)
	SetDeprecatedDocumentation(doc string)
	AddDeprecatedParam(param Parameter)
}

type markdownDocAttachmentImpl struct {
	description             string
	parameters              []Parameter
	returnValueDescription  string
	deprecatedDocumentation string
	deprecatedParams        []Parameter
}

func NewMarkdownDocAttachment(paramCount int) MarkdownDocAttachment {
	return &markdownDocAttachmentImpl{
		parameters:       make([]Parameter, 0, paramCount),
		deprecatedParams: make([]Parameter, 0),
	}
}

func (m *markdownDocAttachmentImpl) GetDescription() string {
	return m.description
}

func (m *markdownDocAttachmentImpl) GetParameters() []Parameter {
	return m.parameters
}

func (m *markdownDocAttachmentImpl) GetReturnValueDescription() string {
	return m.returnValueDescription
}

func (m *markdownDocAttachmentImpl) GetDeprecatedDocumentation() string {
	return m.deprecatedDocumentation
}

func (m *markdownDocAttachmentImpl) GetDeprecatedParams() []Parameter {
	return m.deprecatedParams
}

func (m *markdownDocAttachmentImpl) SetDescription(description string) {
	m.description = description
}

func (m *markdownDocAttachmentImpl) AddParameter(param Parameter) {
	m.parameters = append(m.parameters, param)
}

func (m *markdownDocAttachmentImpl) SetReturnValueDescription(description string) {
	m.returnValueDescription = description
}

func (m *markdownDocAttachmentImpl) SetDeprecatedDocumentation(doc string) {
	m.deprecatedDocumentation = doc
}

func (m *markdownDocAttachmentImpl) AddDeprecatedParam(param Parameter) {
	m.deprecatedParams = append(m.deprecatedParams, param)
}

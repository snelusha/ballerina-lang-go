package elements

type MarkdownDocAttachment interface {
	GetDescription() string
	GetParameters() []Parameter
	GetReturnValueDescription() string
	GetDeprecatedDocumentation() string
	GetDeprecatedParameters() []Parameter
}

type markdownDocAttachmentImpl struct {
	description             string
	parameters              []Parameter
	returnValueDescription  string
	deprecatedDocumentation string
	deprecatedParams        []Parameter
}

func NewMarkdownDocAttachment(description string, parameters []Parameter, returnValueDescription string, deprecatedDocumentation string, deprecatedParams []Parameter) MarkdownDocAttachment {
	return &markdownDocAttachmentImpl{
		description:             description,
		parameters:              parameters,
		returnValueDescription:  returnValueDescription,
		deprecatedDocumentation: deprecatedDocumentation,
		deprecatedParams:        deprecatedParams,
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

func (m *markdownDocAttachmentImpl) GetDeprecatedParameters() []Parameter {
	return m.deprecatedParams
}

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

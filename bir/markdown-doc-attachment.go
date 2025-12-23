package bir

type MarkdownDocAttachment struct {
	Description             string
	Parameters              []*Parameter
	ReturnValueDescription  string
	DeprecatedDocumentation string
	DeprecatedParams        []*Parameter
}

type Parameter struct {
	Name        string
	Description string
}

func NewMarkdownDocAttachment(paramCount int) *MarkdownDocAttachment {
	return &MarkdownDocAttachment{
		Parameters:       make([]*Parameter, 0, paramCount),
		DeprecatedParams: make([]*Parameter, 0),
	}
}

func NewParameter(name, description string) *Parameter {
	return &Parameter{
		Name:        name,
		Description: description,
	}
}

func (p *Parameter) GetName() string {
	return p.Name
}

func (p *Parameter) GetDescription() string {
	return p.Description
}

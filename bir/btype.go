package bir

type BType interface {
	Accept(visitor BTypeVisitor)
	String() string
}

type BInvokableType interface {
	BType
	GetReturnType() BType
	GetParameterTypes() []BType
}

type BTypeVisitor interface {
	Visit(btype BType)
}

type NamedNode interface {
	GetName() *Name
}

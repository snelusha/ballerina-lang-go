package bir

type BIRTerminator interface {
	Accept(visitor BIRVisitor)
	GetPos() Location
}

type BIRNonTerminator interface {
	Accept(visitor BIRVisitor)
	GetPos() Location
}

type BIRInstruction interface {
	Accept(visitor BIRVisitor)
	GetPos() Location
}

type Lock interface {
	BIRTerminator
}

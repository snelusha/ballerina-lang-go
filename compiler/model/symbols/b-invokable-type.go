package symbols

type BInvokableType interface {
	BType
	GetParameterTypes() []BType
	GetRestType() BType
	GetRetType() BType
	AddParamType(bType BType)
	SetParamTypes(paramTypes []BType)
	GetTypeSignature() string
}

type bInvokableTypeImpl struct {
	*bTypeImpl
	paramTypes []BType
	restType   BType
	retType    BType
}

func NewBInvokableType(paramTypes []BType, restType, retType BType, tsymbol BTypeSymbol) BInvokableType {
	return &bInvokableTypeImpl{
		bTypeImpl:  NewBTypeWithName(TypeTagInvokable, tsymbol, nil, FlagReadOnly).(*bTypeImpl),
		paramTypes: paramTypes,
		restType:   restType,
		retType:    retType,
	}
}

func NewBInvokableTypeNoRest(paramTypes []BType, retType BType, tsymbol BTypeSymbol) BInvokableType {
	return NewBInvokableType(paramTypes, nil, retType, tsymbol)
}

func (b *bInvokableTypeImpl) GetParameterTypes() []BType {
	return b.paramTypes
}

func (b *bInvokableTypeImpl) GetRestType() BType {
	return b.restType
}

func (b *bInvokableTypeImpl) GetRetType() BType {
	return b.retType
}

func (b *bInvokableTypeImpl) GetReturnType() BType {
	return b.retType
}

func (b *bInvokableTypeImpl) AddParamType(bType BType) {
	b.paramTypes = append(b.paramTypes, bType)
}

func (b *bInvokableTypeImpl) SetParamTypes(paramTypes []BType) {
	b.paramTypes = paramTypes
}

func (b *bInvokableTypeImpl) GetTypeSignature() string {
	return "function signature"
}

const (
	TypeTagInvokable = 1
	FlagReadOnly     = 1 << 0
)

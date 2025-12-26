package symbols

type BTypeSymbol interface{}

type bTypeSymbolImpl struct {
	isTypeParamResolved bool
	typeParamTSymbol    BTypeSymbol
	annotations         BVarSymbol
}

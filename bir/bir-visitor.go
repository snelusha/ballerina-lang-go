package bir

type BIRVisitor interface {
	VisitPackage(pkg *BIRPackage)
	VisitImportModule(importModule *BIRImportModule)
	VisitTypeDefinition(typeDef *BIRTypeDefinition)
	VisitVariableDcl(varDcl *BIRVariableDcl)
	VisitFunctionParameter(funcParam *BIRFunctionParameter)
	VisitFunction(function *BIRFunction)
	VisitBasicBlock(basicBlock *BIRBasicBlock)
	VisitParameter(param *BIRParameter)
	VisitAnnotation(annotation *BIRAnnotation)
	VisitConstant(constant *BIRConstant)
	VisitAnnotationAttachment(annotAttach *BIRAnnotationAttachment)
	VisitConstAnnotationAttachment(constAnnotAttach *BIRConstAnnotationAttachment)
	VisitErrorEntry(errorEntry *BIRErrorEntry)
	VisitServiceDeclaration(serviceDecl *BIRServiceDeclaration)
	VisitOperand(operand *BIROperand)
	VisitTerminator(terminator BIRTerminator)
	VisitNonTerminator(nonTerminator BIRNonTerminator)
}

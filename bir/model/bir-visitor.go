package model

type BIRVisitor interface {
	VisitBIRPackage(pkg BIRPackage)
	VisitBIRImportModule(module BIRImportModule)
	VisitBIRTypeDefinition(typeDef BIRTypeDefinition)
	VisitBIRVariableDcl(varDecl BIRVariableDcl)
	VisitBIRFunctionParameter(param BIRFunctionParameter)
	VisitBIRFunction(function BIRFunction)
	VisitBIRBasicBlock(bb BIRBasicBlock)
	VisitBIRParameter(param BIRParameter)
	VisitBIRAnnotation(annot BIRAnnotation)
	VisitBIRConstant(constant BIRConstant)
	VisitBIRAnnotationAttachment(attach BIRAnnotationAttachment)
	VisitBIRConstAnnotationAttachment(attach BIRConstAnnotationAttachment)
	VisitBIRErrorEntry(entry BIRErrorEntry)
	VisitBIRServiceDeclaration(decl BIRServiceDeclaration)
	VisitBIROperand(operand BIROperand)
}

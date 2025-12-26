package model

type BIRVisitor interface {
	VisitPackage(pkg BIRPackage)
	VisitImportModule(importModule BIRImportModule)
	VisitTypeDefinition(typeDef BIRTypeDefinition)
	VisitVariableDcl(varDcl BIRVariableDcl)
	VisitFunctionParameter(funcParam BIRFunctionParameter)
	VisitFunction(function BIRFunction)
	VisitBasicBlock(basicBlock BIRBasicBlock)
	VisitParameter(parameter BIRParameter)
	VisitAnnotation(annotation BIRAnnotation)
	VisitConstant(constant BIRConstant)
	VisitAnnotationAttachment(annotAttach BIRAnnotationAttachment)
	VisitConstAnnotationAttachment(constAnnotAttach BIRConstAnnotationAttachment)
	VisitErrorEntry(errorEntry BIRErrorEntry)
	VisitServiceDeclaration(serviceDecl BIRServiceDeclaration)
	VisitGlobalVariableDcl(globalVarDcl BIRGlobalVariableDcl)
	VisitOperand(operand BIROperand)
}

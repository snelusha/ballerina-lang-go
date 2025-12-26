package model

import (
	"fmt"

	"ballerina-lang-go/compiler/model/elements"
	"ballerina-lang-go/compiler/model/symbols"
	"ballerina-lang-go/compiler/semantics/model/types"
	"ballerina-lang-go/diagnostics"
)

type BIRNode interface {
	GetPos() diagnostics.Location
	Accept(visitor BIRVisitor)
}

type birNodeImpl struct {
	pos diagnostics.Location
}

func newBIRNode(pos diagnostics.Location) *birNodeImpl {
	return &birNodeImpl{pos: pos}
}

func (b *birNodeImpl) GetPos() diagnostics.Location {
	return b.pos
}

type BIRPackage interface {
	BIRNode
	GetPackageID() elements.PackageID
	GetImportModules() []BIRImportModule
	GetTypeDefs() []BIRTypeDefinition
	GetGlobalVars() []BIRGlobalVariableDcl
	GetImportedGlobalVarsDummyVarDcls() []BIRGlobalVariableDcl
	GetFunctions() []BIRFunction
	GetAnnotations() []BIRAnnotation
	GetConstants() []BIRConstant
	GetServiceDecls() []BIRServiceDeclaration
	IsListenerAvailable() bool
	SetListenerAvailable(available bool)
	GetRecordDefaultValueMap() map[string]map[string]string
	AddImportModule(module BIRImportModule)
	AddTypeDef(typeDef BIRTypeDefinition)
	AddGlobalVar(globalVar BIRGlobalVariableDcl)
	AddFunction(function BIRFunction)
	AddAnnotation(annotation BIRAnnotation)
	AddConstant(constant BIRConstant)
	AddServiceDecl(serviceDecl BIRServiceDeclaration)
}

type birPackageImpl struct {
	*birNodeImpl
	packageID                      elements.PackageID
	importModules                  []BIRImportModule
	typeDefs                       []BIRTypeDefinition
	globalVars                     []BIRGlobalVariableDcl
	importedGlobalVarsDummyVarDcls []BIRGlobalVariableDcl
	functions                      []BIRFunction
	annotations                    []BIRAnnotation
	constants                      []BIRConstant
	serviceDecls                   []BIRServiceDeclaration
	isListenerAvailable            bool
	recordDefaultValueMap          map[string]map[string]string
}

func NewBIRPackage(pos diagnostics.Location, org, pkgName, name, version, sourceFileName Name, sourceRoot string, skipTest, isTestPkg bool) BIRPackage {
	pkgID := elements.NewPackageIDFull(org, pkgName, name, version, sourceFileName, sourceRoot, isTestPkg, skipTest)
	return &birPackageImpl{
		birNodeImpl:                    newBIRNode(pos),
		packageID:                      pkgID,
		importModules:                  make([]BIRImportModule, 0),
		typeDefs:                       make([]BIRTypeDefinition, 0),
		globalVars:                     make([]BIRGlobalVariableDcl, 0),
		importedGlobalVarsDummyVarDcls: make([]BIRGlobalVariableDcl, 0),
		functions:                      make([]BIRFunction, 0),
		annotations:                    make([]BIRAnnotation, 0),
		constants:                      make([]BIRConstant, 0),
		serviceDecls:                   make([]BIRServiceDeclaration, 0),
		recordDefaultValueMap:          make(map[string]map[string]string),
	}
}

func (b *birPackageImpl) GetPackageID() elements.PackageID {
	return b.packageID
}

func (b *birPackageImpl) GetImportModules() []BIRImportModule {
	return b.importModules
}

func (b *birPackageImpl) GetTypeDefs() []BIRTypeDefinition {
	return b.typeDefs
}

func (b *birPackageImpl) GetGlobalVars() []BIRGlobalVariableDcl {
	return b.globalVars
}

func (b *birPackageImpl) GetImportedGlobalVarsDummyVarDcls() []BIRGlobalVariableDcl {
	return b.importedGlobalVarsDummyVarDcls
}

func (b *birPackageImpl) GetFunctions() []BIRFunction {
	return b.functions
}

func (b *birPackageImpl) GetAnnotations() []BIRAnnotation {
	return b.annotations
}

func (b *birPackageImpl) GetConstants() []BIRConstant {
	return b.constants
}

func (b *birPackageImpl) GetServiceDecls() []BIRServiceDeclaration {
	return b.serviceDecls
}

func (b *birPackageImpl) IsListenerAvailable() bool {
	return b.isListenerAvailable
}

func (b *birPackageImpl) SetListenerAvailable(available bool) {
	b.isListenerAvailable = available
}

func (b *birPackageImpl) GetRecordDefaultValueMap() map[string]map[string]string {
	return b.recordDefaultValueMap
}

func (b *birPackageImpl) AddImportModule(module BIRImportModule) {
	b.importModules = append(b.importModules, module)
}

func (b *birPackageImpl) AddTypeDef(typeDef BIRTypeDefinition) {
	b.typeDefs = append(b.typeDefs, typeDef)
}

func (b *birPackageImpl) AddGlobalVar(globalVar BIRGlobalVariableDcl) {
	b.globalVars = append(b.globalVars, globalVar)
}

func (b *birPackageImpl) AddFunction(function BIRFunction) {
	b.functions = append(b.functions, function)
}

func (b *birPackageImpl) AddAnnotation(annotation BIRAnnotation) {
	b.annotations = append(b.annotations, annotation)
}

func (b *birPackageImpl) AddConstant(constant BIRConstant) {
	b.constants = append(b.constants, constant)
}

func (b *birPackageImpl) AddServiceDecl(serviceDecl BIRServiceDeclaration) {
	b.serviceDecls = append(b.serviceDecls, serviceDecl)
}

func (b *birPackageImpl) Accept(visitor BIRVisitor) {
	visitor.VisitPackage(b)
}

type BIRImportModule interface {
	BIRNode
	GetPackageID() elements.PackageID
}

type birImportModuleImpl struct {
	*birNodeImpl
	packageID elements.PackageID
}

func NewBIRImportModule(pos diagnostics.Location, org, name, version Name) BIRImportModule {
	return &birImportModuleImpl{
		birNodeImpl: newBIRNode(pos),
		packageID:   elements.NewPackageID(org, name, version),
	}
}

func (b *birImportModuleImpl) GetPackageID() elements.PackageID {
	return b.packageID
}

func (b *birImportModuleImpl) Accept(visitor BIRVisitor) {
	visitor.VisitImportModule(b)
}

type BIRDocumentableNode interface {
	BIRNode
	GetMarkdownDocAttachment() elements.MarkdownDocAttachment
	SetMarkdownDocAttachment(attachment elements.MarkdownDocAttachment)
}

type birDocumentableNodeImpl struct {
	*birNodeImpl
	markdownDocAttachment elements.MarkdownDocAttachment
}

func newBIRDocumentableNode(pos diagnostics.Location) *birDocumentableNodeImpl {
	return &birDocumentableNodeImpl{
		birNodeImpl: newBIRNode(pos),
	}
}

func (b *birDocumentableNodeImpl) GetMarkdownDocAttachment() elements.MarkdownDocAttachment {
	return b.markdownDocAttachment
}

func (b *birDocumentableNodeImpl) SetMarkdownDocAttachment(attachment elements.MarkdownDocAttachment) {
	b.markdownDocAttachment = attachment
}

type BIRVariableDcl interface {
	BIRDocumentableNode
	GetType() types.BType
	SetType(bType types.BType)
	GetName() Name
	GetOriginalName() Name
	GetMetaVarName() string
	GetJvmVarName() string
	GetKind() VarKind
	GetScope() VarScope
	IsIgnoreVariable() bool
	SetIgnoreVariable(ignore bool)
	GetEndBB() BIRBasicBlock
	SetEndBB(bb BIRBasicBlock)
	GetStartBB() BIRBasicBlock
	SetStartBB(bb BIRBasicBlock)
	GetInsOffset() int
	SetInsOffset(offset int)
	IsOnlyUsedInSingleBB() bool
	SetOnlyUsedInSingleBB(singleBB bool)
	IsInitialized() bool
	SetInitialized(initialized bool)
	GetInsScope() BirScope
	SetInsScope(scope BirScope)
	String() string
}

type birVariableDclImpl struct {
	*birDocumentableNodeImpl
	bType              types.BType
	name               Name
	originalName       Name
	metaVarName        string
	jvmVarName         string
	kind               VarKind
	scope              VarScope
	ignoreVariable     bool
	endBB              BIRBasicBlock
	startBB            BIRBasicBlock
	insOffset          int
	onlyUsedInSingleBB bool
	initialized        bool
	insScope           BirScope
}

func NewBIRVariableDcl(pos diagnostics.Location, bType types.BType, name, originalName Name, scope VarScope, kind VarKind, metaVarName string) BIRVariableDcl {
	jvmVarName := name.GetValue()
	// TODO: Implement proper replacement logic for %
	return &birVariableDclImpl{
		birDocumentableNodeImpl: newBIRDocumentableNode(pos),
		bType:                   bType,
		name:                    name,
		originalName:            originalName,
		scope:                   scope,
		kind:                    kind,
		metaVarName:             metaVarName,
		jvmVarName:              jvmVarName,
	}
}

func NewBIRVariableDclSimple(pos diagnostics.Location, bType types.BType, name Name, scope VarScope, kind VarKind, metaVarName string) BIRVariableDcl {
	return NewBIRVariableDcl(pos, bType, name, name, scope, kind, metaVarName)
}

func (b *birVariableDclImpl) GetType() types.BType {
	return b.bType
}

func (b *birVariableDclImpl) SetType(bType types.BType) {
	b.bType = bType
}

func (b *birVariableDclImpl) GetName() Name {
	return b.name
}

func (b *birVariableDclImpl) GetOriginalName() Name {
	return b.originalName
}

func (b *birVariableDclImpl) GetMetaVarName() string {
	return b.metaVarName
}

func (b *birVariableDclImpl) GetJvmVarName() string {
	return b.jvmVarName
}

func (b *birVariableDclImpl) GetKind() VarKind {
	return b.kind
}

func (b *birVariableDclImpl) GetScope() VarScope {
	return b.scope
}

func (b *birVariableDclImpl) IsIgnoreVariable() bool {
	return b.ignoreVariable
}

func (b *birVariableDclImpl) SetIgnoreVariable(ignore bool) {
	b.ignoreVariable = ignore
}

func (b *birVariableDclImpl) GetEndBB() BIRBasicBlock {
	return b.endBB
}

func (b *birVariableDclImpl) SetEndBB(bb BIRBasicBlock) {
	b.endBB = bb
}

func (b *birVariableDclImpl) GetStartBB() BIRBasicBlock {
	return b.startBB
}

func (b *birVariableDclImpl) SetStartBB(bb BIRBasicBlock) {
	b.startBB = bb
}

func (b *birVariableDclImpl) GetInsOffset() int {
	return b.insOffset
}

func (b *birVariableDclImpl) SetInsOffset(offset int) {
	b.insOffset = offset
}

func (b *birVariableDclImpl) IsOnlyUsedInSingleBB() bool {
	return b.onlyUsedInSingleBB
}

func (b *birVariableDclImpl) SetOnlyUsedInSingleBB(singleBB bool) {
	b.onlyUsedInSingleBB = singleBB
}

func (b *birVariableDclImpl) IsInitialized() bool {
	return b.initialized
}

func (b *birVariableDclImpl) SetInitialized(initialized bool) {
	b.initialized = initialized
}

func (b *birVariableDclImpl) GetInsScope() BirScope {
	return b.insScope
}

func (b *birVariableDclImpl) SetInsScope(scope BirScope) {
	b.insScope = scope
}

func (b *birVariableDclImpl) String() string {
	return b.name.GetValue()
}

func (b *birVariableDclImpl) Accept(visitor BIRVisitor) {
	visitor.VisitVariableDcl(b)
}

type BIRParameter interface {
	BIRNode
	GetName() Name
	GetFlags() int64
	GetAnnotAttachments() []BIRAnnotationAttachment
	AddAnnotAttachment(attachment BIRAnnotationAttachment)
}

type birParameterImpl struct {
	*birNodeImpl
	name             Name
	flags            int64
	annotAttachments []BIRAnnotationAttachment
}

func NewBIRParameter(pos diagnostics.Location, name Name, flags int64) BIRParameter {
	return &birParameterImpl{
		birNodeImpl:      newBIRNode(pos),
		name:             name,
		flags:            flags,
		annotAttachments: make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birParameterImpl) GetName() Name {
	return b.name
}

func (b *birParameterImpl) GetFlags() int64 {
	return b.flags
}

func (b *birParameterImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birParameterImpl) AddAnnotAttachment(attachment BIRAnnotationAttachment) {
	b.annotAttachments = append(b.annotAttachments, attachment)
}

func (b *birParameterImpl) Accept(visitor BIRVisitor) {
	visitor.VisitParameter(b)
}

type BIRGlobalVariableDcl interface {
	BIRVariableDcl
	GetFlags() int64
	GetPkgID() elements.PackageID
	GetOrigin() symbols.SymbolOrigin
	GetAnnotAttachments() []BIRAnnotationAttachment
	AddAnnotAttachment(attachment BIRAnnotationAttachment)
}

type birGlobalVariableDclImpl struct {
	*birVariableDclImpl
	flags            int64
	pkgID            elements.PackageID
	origin           symbols.SymbolOrigin
	annotAttachments []BIRAnnotationAttachment
}

func NewBIRGlobalVariableDcl(pos diagnostics.Location, flags int64, bType types.BType, pkgID elements.PackageID, name, originalName Name, scope VarScope, kind VarKind, metaVarName string, origin symbols.SymbolOrigin) BIRGlobalVariableDcl {
	return &birGlobalVariableDclImpl{
		birVariableDclImpl: NewBIRVariableDcl(pos, bType, name, originalName, scope, kind, metaVarName).(*birVariableDclImpl),
		flags:              flags,
		pkgID:              pkgID,
		origin:             origin,
		annotAttachments:   make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birGlobalVariableDclImpl) GetFlags() int64 {
	return b.flags
}

func (b *birGlobalVariableDclImpl) GetPkgID() elements.PackageID {
	return b.pkgID
}

func (b *birGlobalVariableDclImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birGlobalVariableDclImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birGlobalVariableDclImpl) AddAnnotAttachment(attachment BIRAnnotationAttachment) {
	b.annotAttachments = append(b.annotAttachments, attachment)
}

func (b *birGlobalVariableDclImpl) Accept(visitor BIRVisitor) {
	visitor.VisitGlobalVariableDcl(b)
}

type BIRFunctionParameter interface {
	BIRVariableDcl
	HasDefaultExpr() bool
	IsPathParameter() bool
	SetPathParameter(isPath bool)
}

type birFunctionParameterImpl struct {
	*birVariableDclImpl
	hasDefaultExpr  bool
	isPathParameter bool
}

func NewBIRFunctionParameter(pos diagnostics.Location, bType types.BType, name Name, scope VarScope, kind VarKind, metaVarName string, hasDefaultExpr, isPathParameter bool) BIRFunctionParameter {
	return &birFunctionParameterImpl{
		birVariableDclImpl: NewBIRVariableDclSimple(pos, bType, name, scope, kind, metaVarName).(*birVariableDclImpl),
		hasDefaultExpr:     hasDefaultExpr,
		isPathParameter:    isPathParameter,
	}
}

func (b *birFunctionParameterImpl) HasDefaultExpr() bool {
	return b.hasDefaultExpr
}

func (b *birFunctionParameterImpl) IsPathParameter() bool {
	return b.isPathParameter
}

func (b *birFunctionParameterImpl) SetPathParameter(isPath bool) {
	b.isPathParameter = isPath
}

func (b *birFunctionParameterImpl) Accept(visitor BIRVisitor) {
	visitor.VisitFunctionParameter(b)
}

type BIRFunction interface {
	BIRDocumentableNode
	types.NamedNode
	GetName() Name
	GetOriginalName() Name
	GetFlags() int64
	GetOrigin() symbols.SymbolOrigin
	GetType() types.BInvokableType
	GetRequiredParams() []BIRParameter
	GetReceiver() BIRVariableDcl
	SetReceiver(receiver BIRVariableDcl)
	GetRestParam() BIRParameter
	GetArgsCount() int
	GetLocalVars() []BIRVariableDcl
	GetReturnVariable() BIRVariableDcl
	GetParameters() []BIRFunctionParameter
	GetBasicBlocks() []BIRBasicBlock
	GetErrorTable() []BIRErrorEntry
	GetWorkerName() Name
	GetWorkerChannels() []ChannelDetails
	GetAnnotAttachments() []BIRAnnotationAttachment
	GetAnnotAttachmentsOnExternal() []BIRAnnotationAttachment
	SetAnnotAttachmentsOnExternal(attachments []BIRAnnotationAttachment)
	GetReturnTypeAnnots() []BIRAnnotationAttachment
	GetDependentGlobalVars() []BIRGlobalVariableDcl
	GetPathParams() []BIRVariableDcl
	GetRestPathParam() BIRVariableDcl
	GetResourcePath() []Name
	GetResourcePathSegmentPosList() []diagnostics.Location
	GetAccessor() Name
	GetPathSegmentTypeList() []types.BType
	HasWorkers() bool
	SetHasWorkers(hasWorkers bool)
	AddLocalVar(localVar BIRVariableDcl)
	AddParameter(param BIRFunctionParameter)
	AddBasicBlock(bb BIRBasicBlock)
	AddErrorEntry(entry BIRErrorEntry)
	AddAnnotAttachment(attachment BIRAnnotationAttachment)
	AddReturnTypeAnnot(annot BIRAnnotationAttachment)
}

type birFunctionImpl struct {
	*birDocumentableNodeImpl
	name                       Name
	originalName               Name
	flags                      int64
	origin                     symbols.SymbolOrigin
	funcType                   types.BInvokableType
	requiredParams             []BIRParameter
	receiver                   BIRVariableDcl
	restParam                  BIRParameter
	argsCount                  int
	localVars                  []BIRVariableDcl
	returnVariable             BIRVariableDcl
	parameters                 []BIRFunctionParameter
	basicBlocks                []BIRBasicBlock
	errorTable                 []BIRErrorEntry
	workerName                 Name
	workerChannels             []ChannelDetails
	annotAttachments           []BIRAnnotationAttachment
	annotAttachmentsOnExternal []BIRAnnotationAttachment
	returnTypeAnnots           []BIRAnnotationAttachment
	dependentGlobalVars        []BIRGlobalVariableDcl
	pathParams                 []BIRVariableDcl
	restPathParam              BIRVariableDcl
	resourcePath               []Name
	resourcePathSegmentPosList []diagnostics.Location
	accessor                   Name
	pathSegmentTypeList        []types.BType
	hasWorkers                 bool
}

func NewBIRFunction(pos diagnostics.Location, name, originalName Name, flags int64, funcType types.BInvokableType, workerName Name, sendInsCount int, origin symbols.SymbolOrigin) BIRFunction {
	return &birFunctionImpl{
		birDocumentableNodeImpl: newBIRDocumentableNode(pos),
		name:                    name,
		originalName:            originalName,
		flags:                   flags,
		origin:                  origin,
		funcType:                funcType,
		localVars:               make([]BIRVariableDcl, 0),
		parameters:              make([]BIRFunctionParameter, 0),
		requiredParams:          make([]BIRParameter, 0),
		basicBlocks:             make([]BIRBasicBlock, 0),
		errorTable:              make([]BIRErrorEntry, 0),
		workerName:              workerName,
		workerChannels:          make([]ChannelDetails, sendInsCount),
		annotAttachments:        make([]BIRAnnotationAttachment, 0),
		returnTypeAnnots:        make([]BIRAnnotationAttachment, 0),
		dependentGlobalVars:     make([]BIRGlobalVariableDcl, 0),
	}
}

func (b *birFunctionImpl) GetName() Name {
	return b.name
}

func (b *birFunctionImpl) GetOriginalName() Name {
	return b.originalName
}

func (b *birFunctionImpl) GetFlags() int64 {
	return b.flags
}

func (b *birFunctionImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birFunctionImpl) GetType() types.BInvokableType {
	return b.funcType
}

func (b *birFunctionImpl) GetRequiredParams() []BIRParameter {
	return b.requiredParams
}

func (b *birFunctionImpl) GetReceiver() BIRVariableDcl {
	return b.receiver
}

func (b *birFunctionImpl) SetReceiver(receiver BIRVariableDcl) {
	b.receiver = receiver
}

func (b *birFunctionImpl) GetRestParam() BIRParameter {
	return b.restParam
}

func (b *birFunctionImpl) GetArgsCount() int {
	return b.argsCount
}

func (b *birFunctionImpl) GetLocalVars() []BIRVariableDcl {
	return b.localVars
}

func (b *birFunctionImpl) GetReturnVariable() BIRVariableDcl {
	return b.returnVariable
}

func (b *birFunctionImpl) GetParameters() []BIRFunctionParameter {
	return b.parameters
}

func (b *birFunctionImpl) GetBasicBlocks() []BIRBasicBlock {
	return b.basicBlocks
}

func (b *birFunctionImpl) GetErrorTable() []BIRErrorEntry {
	return b.errorTable
}

func (b *birFunctionImpl) GetWorkerName() Name {
	return b.workerName
}

func (b *birFunctionImpl) GetWorkerChannels() []ChannelDetails {
	return b.workerChannels
}

func (b *birFunctionImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birFunctionImpl) GetAnnotAttachmentsOnExternal() []BIRAnnotationAttachment {
	return b.annotAttachmentsOnExternal
}

func (b *birFunctionImpl) SetAnnotAttachmentsOnExternal(attachments []BIRAnnotationAttachment) {
	b.annotAttachmentsOnExternal = attachments
}

func (b *birFunctionImpl) GetReturnTypeAnnots() []BIRAnnotationAttachment {
	return b.returnTypeAnnots
}

func (b *birFunctionImpl) GetDependentGlobalVars() []BIRGlobalVariableDcl {
	return b.dependentGlobalVars
}

func (b *birFunctionImpl) GetPathParams() []BIRVariableDcl {
	return b.pathParams
}

func (b *birFunctionImpl) GetRestPathParam() BIRVariableDcl {
	return b.restPathParam
}

func (b *birFunctionImpl) GetResourcePath() []Name {
	return b.resourcePath
}

func (b *birFunctionImpl) GetResourcePathSegmentPosList() []diagnostics.Location {
	return b.resourcePathSegmentPosList
}

func (b *birFunctionImpl) GetAccessor() Name {
	return b.accessor
}

func (b *birFunctionImpl) GetPathSegmentTypeList() []types.BType {
	return b.pathSegmentTypeList
}

func (b *birFunctionImpl) HasWorkers() bool {
	return b.hasWorkers
}

func (b *birFunctionImpl) SetHasWorkers(hasWorkers bool) {
	b.hasWorkers = hasWorkers
}

func (b *birFunctionImpl) AddLocalVar(localVar BIRVariableDcl) {
	b.localVars = append(b.localVars, localVar)
}

func (b *birFunctionImpl) AddParameter(param BIRFunctionParameter) {
	b.parameters = append(b.parameters, param)
}

func (b *birFunctionImpl) AddBasicBlock(bb BIRBasicBlock) {
	b.basicBlocks = append(b.basicBlocks, bb)
}

func (b *birFunctionImpl) AddErrorEntry(entry BIRErrorEntry) {
	b.errorTable = append(b.errorTable, entry)
}

func (b *birFunctionImpl) AddAnnotAttachment(attachment BIRAnnotationAttachment) {
	b.annotAttachments = append(b.annotAttachments, attachment)
}

func (b *birFunctionImpl) AddReturnTypeAnnot(annot BIRAnnotationAttachment) {
	b.returnTypeAnnots = append(b.returnTypeAnnots, annot)
}

func (b *birFunctionImpl) Accept(visitor BIRVisitor) {
	visitor.VisitFunction(b)
}

type BIRBasicBlock interface {
	BIRNode
	GetNumber() int
	GetID() Name
	GetInstructions() []BIRNonTerminator
	GetTerminator() BIRTerminator
	SetTerminator(terminator BIRTerminator)
	AddInstruction(instruction BIRNonTerminator)
	String() string
}

type birBasicBlockImpl struct {
	*birNodeImpl
	number       int
	id           Name
	instructions []BIRNonTerminator
	terminator   BIRTerminator
}

const BIRBasicBlockPrefix = "bb"

func NewBIRBasicBlock(id Name, number int) BIRBasicBlock {
	return &birBasicBlockImpl{
		birNodeImpl:  newBIRNode(nil),
		number:       number,
		id:           id,
		instructions: make([]BIRNonTerminator, 0),
	}
}

func NewBIRBasicBlockSimple(number int) BIRBasicBlock {
	return NewBIRBasicBlock(NewName(fmt.Sprintf("%s%d", BIRBasicBlockPrefix, number)), number)
}

func NewBIRBasicBlockWithPrefix(idPrefix string, number int) BIRBasicBlock {
	return NewBIRBasicBlock(NewName(fmt.Sprintf("%s%d", idPrefix, number)), number)
}

func (b *birBasicBlockImpl) GetNumber() int {
	return b.number
}

func (b *birBasicBlockImpl) GetID() Name {
	return b.id
}

func (b *birBasicBlockImpl) GetInstructions() []BIRNonTerminator {
	return b.instructions
}

func (b *birBasicBlockImpl) GetTerminator() BIRTerminator {
	return b.terminator
}

func (b *birBasicBlockImpl) SetTerminator(terminator BIRTerminator) {
	b.terminator = terminator
}

func (b *birBasicBlockImpl) AddInstruction(instruction BIRNonTerminator) {
	b.instructions = append(b.instructions, instruction)
}

func (b *birBasicBlockImpl) String() string {
	return b.id.GetValue()
}

func (b *birBasicBlockImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBasicBlock(b)
}

type BIRTypeDefinition interface {
	BIRDocumentableNode
	types.NamedNode
	GetName() Name
	GetOriginalName() Name
	GetInternalName() Name
	GetAttachedFuncs() []BIRFunction
	GetFlags() int64
	GetType() types.BType
	IsBuiltin() bool
	GetReferencedTypes() []types.BType
	GetReferenceType() types.BType
	SetReferenceType(refType types.BType)
	GetOrigin() symbols.SymbolOrigin
	GetAnnotAttachments() []BIRAnnotationAttachment
	GetIndex() int
	SetIndex(index int)
	AddAttachedFunc(function BIRFunction)
	AddReferencedType(bType types.BType)
	AddAnnotAttachment(attachment BIRAnnotationAttachment)
}

type birTypeDefinitionImpl struct {
	*birDocumentableNodeImpl
	name             Name
	originalName     Name
	internalName     Name
	attachedFuncs    []BIRFunction
	flags            int64
	bType            types.BType
	isBuiltin        bool
	referencedTypes  []types.BType
	referenceType    types.BType
	origin           symbols.SymbolOrigin
	annotAttachments []BIRAnnotationAttachment
	index            int
}

func NewBIRTypeDefinition(pos diagnostics.Location, internalName Name, flags int64, isBuiltin bool, bType types.BType, attachedFuncs []BIRFunction, origin symbols.SymbolOrigin, name, originalName Name) BIRTypeDefinition {
	return &birTypeDefinitionImpl{
		birDocumentableNodeImpl: newBIRDocumentableNode(pos),
		internalName:            internalName,
		flags:                   flags,
		isBuiltin:               isBuiltin,
		bType:                   bType,
		attachedFuncs:           attachedFuncs,
		referencedTypes:         make([]types.BType, 0),
		origin:                  origin,
		name:                    name,
		originalName:            originalName,
		annotAttachments:        make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birTypeDefinitionImpl) GetName() Name {
	return b.name
}

func (b *birTypeDefinitionImpl) GetOriginalName() Name {
	return b.originalName
}

func (b *birTypeDefinitionImpl) GetInternalName() Name {
	return b.internalName
}

func (b *birTypeDefinitionImpl) GetAttachedFuncs() []BIRFunction {
	return b.attachedFuncs
}

func (b *birTypeDefinitionImpl) GetFlags() int64 {
	return b.flags
}

func (b *birTypeDefinitionImpl) GetType() types.BType {
	return b.bType
}

func (b *birTypeDefinitionImpl) IsBuiltin() bool {
	return b.isBuiltin
}

func (b *birTypeDefinitionImpl) GetReferencedTypes() []types.BType {
	return b.referencedTypes
}

func (b *birTypeDefinitionImpl) GetReferenceType() types.BType {
	return b.referenceType
}

func (b *birTypeDefinitionImpl) SetReferenceType(refType types.BType) {
	b.referenceType = refType
}

func (b *birTypeDefinitionImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birTypeDefinitionImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birTypeDefinitionImpl) GetIndex() int {
	return b.index
}

func (b *birTypeDefinitionImpl) SetIndex(index int) {
	b.index = index
}

func (b *birTypeDefinitionImpl) AddAttachedFunc(function BIRFunction) {
	b.attachedFuncs = append(b.attachedFuncs, function)
}

func (b *birTypeDefinitionImpl) AddReferencedType(bType types.BType) {
	b.referencedTypes = append(b.referencedTypes, bType)
}

func (b *birTypeDefinitionImpl) AddAnnotAttachment(attachment BIRAnnotationAttachment) {
	b.annotAttachments = append(b.annotAttachments, attachment)
}

func (b *birTypeDefinitionImpl) Accept(visitor BIRVisitor) {
	visitor.VisitTypeDefinition(b)
}

type BIRErrorEntry interface {
	BIRNode
	GetTrapBB() BIRBasicBlock
	GetEndBB() BIRBasicBlock
	GetErrorOp() BIROperand
	GetTargetBB() BIRBasicBlock
}

type birErrorEntryImpl struct {
	*birNodeImpl
	trapBB   BIRBasicBlock
	endBB    BIRBasicBlock
	errorOp  BIROperand
	targetBB BIRBasicBlock
}

func NewBIRErrorEntry(trapBB, endBB BIRBasicBlock, errorOp BIROperand, targetBB BIRBasicBlock) BIRErrorEntry {
	return &birErrorEntryImpl{
		birNodeImpl: newBIRNode(nil),
		trapBB:      trapBB,
		endBB:       endBB,
		errorOp:     errorOp,
		targetBB:    targetBB,
	}
}

func (b *birErrorEntryImpl) GetTrapBB() BIRBasicBlock {
	return b.trapBB
}

func (b *birErrorEntryImpl) GetEndBB() BIRBasicBlock {
	return b.endBB
}

func (b *birErrorEntryImpl) GetErrorOp() BIROperand {
	return b.errorOp
}

func (b *birErrorEntryImpl) GetTargetBB() BIRBasicBlock {
	return b.targetBB
}

func (b *birErrorEntryImpl) Accept(visitor BIRVisitor) {
	visitor.VisitErrorEntry(b)
}

type ChannelDetails interface {
	GetName() string
	IsChannelInSameStrand() bool
	IsSend() bool
	String() string
}

type channelDetailsImpl struct {
	name                string
	channelInSameStrand bool
	send                bool
}

func NewChannelDetails(name string, channelInSameStrand, send bool) ChannelDetails {
	return &channelDetailsImpl{
		name:                name,
		channelInSameStrand: channelInSameStrand,
		send:                send,
	}
}

func (c *channelDetailsImpl) GetName() string {
	return c.name
}

func (c *channelDetailsImpl) IsChannelInSameStrand() bool {
	return c.channelInSameStrand
}

func (c *channelDetailsImpl) IsSend() bool {
	return c.send
}

func (c *channelDetailsImpl) String() string {
	return c.name
}

type BIRAnnotation interface {
	BIRDocumentableNode
	GetName() Name
	GetOriginalName() Name
	GetFlags() int64
	GetOrigin() symbols.SymbolOrigin
	GetAttachPoints() []elements.AttachPoint
	GetAnnotationType() types.BType
	GetPackageID() elements.PackageID
	SetPackageID(pkgID elements.PackageID)
	GetAnnotAttachments() []BIRAnnotationAttachment
	AddAnnotAttachment(attachment BIRAnnotationAttachment)
}

type birAnnotationImpl struct {
	*birDocumentableNodeImpl
	name             Name
	originalName     Name
	flags            int64
	origin           symbols.SymbolOrigin
	attachPoints     []elements.AttachPoint
	annotationType   types.BType
	packageID        elements.PackageID
	annotAttachments []BIRAnnotationAttachment
}

func NewBIRAnnotation(pos diagnostics.Location, name, originalName Name, flags int64, points []elements.AttachPoint, annotationType types.BType, origin symbols.SymbolOrigin) BIRAnnotation {
	return &birAnnotationImpl{
		birDocumentableNodeImpl: newBIRDocumentableNode(pos),
		name:                    name,
		originalName:            originalName,
		flags:                   flags,
		attachPoints:            points,
		annotationType:          annotationType,
		origin:                  origin,
		annotAttachments:        make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birAnnotationImpl) GetName() Name {
	return b.name
}

func (b *birAnnotationImpl) GetOriginalName() Name {
	return b.originalName
}

func (b *birAnnotationImpl) GetFlags() int64 {
	return b.flags
}

func (b *birAnnotationImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birAnnotationImpl) GetAttachPoints() []elements.AttachPoint {
	return b.attachPoints
}

func (b *birAnnotationImpl) GetAnnotationType() types.BType {
	return b.annotationType
}

func (b *birAnnotationImpl) GetPackageID() elements.PackageID {
	return b.packageID
}

func (b *birAnnotationImpl) SetPackageID(pkgID elements.PackageID) {
	b.packageID = pkgID
}

func (b *birAnnotationImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birAnnotationImpl) AddAnnotAttachment(attachment BIRAnnotationAttachment) {
	b.annotAttachments = append(b.annotAttachments, attachment)
}

func (b *birAnnotationImpl) Accept(visitor BIRVisitor) {
	visitor.VisitAnnotation(b)
}

type ConstValue interface {
	GetType() types.BType
	GetValue() interface{}
}

type constValueImpl struct {
	bType types.BType
	value interface{}
}

func NewConstValue(value interface{}, bType types.BType) ConstValue {
	return &constValueImpl{
		value: value,
		bType: bType,
	}
}

func (c *constValueImpl) GetType() types.BType {
	return c.bType
}

func (c *constValueImpl) GetValue() interface{} {
	return c.value
}

type BIRConstant interface {
	BIRDocumentableNode
	GetName() Name
	GetFlags() int64
	GetType() types.BType
	GetConstValue() ConstValue
	GetOrigin() symbols.SymbolOrigin
	GetAnnotAttachments() []BIRAnnotationAttachment
	AddAnnotAttachment(attachment BIRAnnotationAttachment)
}

type birConstantImpl struct {
	*birDocumentableNodeImpl
	name             Name
	flags            int64
	bType            types.BType
	constValue       ConstValue
	origin           symbols.SymbolOrigin
	annotAttachments []BIRAnnotationAttachment
}

func NewBIRConstant(pos diagnostics.Location, name Name, flags int64, bType types.BType, constValue ConstValue, origin symbols.SymbolOrigin) BIRConstant {
	return &birConstantImpl{
		birDocumentableNodeImpl: newBIRDocumentableNode(pos),
		name:                    name,
		flags:                   flags,
		bType:                   bType,
		constValue:              constValue,
		origin:                  origin,
		annotAttachments:        make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birConstantImpl) GetName() Name {
	return b.name
}

func (b *birConstantImpl) GetFlags() int64 {
	return b.flags
}

func (b *birConstantImpl) GetType() types.BType {
	return b.bType
}

func (b *birConstantImpl) GetConstValue() ConstValue {
	return b.constValue
}

func (b *birConstantImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birConstantImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birConstantImpl) AddAnnotAttachment(attachment BIRAnnotationAttachment) {
	b.annotAttachments = append(b.annotAttachments, attachment)
}

func (b *birConstantImpl) Accept(visitor BIRVisitor) {
	visitor.VisitConstant(b)
}

type BIRAnnotationAttachment interface {
	BIRNode
	GetAnnotPkgID() elements.PackageID
	GetAnnotTagRef() Name
}

type birAnnotationAttachmentImpl struct {
	*birNodeImpl
	annotPkgID  elements.PackageID
	annotTagRef Name
}

func NewBIRAnnotationAttachment(pos diagnostics.Location, annotPkgID elements.PackageID, annotTagRef Name) BIRAnnotationAttachment {
	return &birAnnotationAttachmentImpl{
		birNodeImpl: newBIRNode(pos),
		annotPkgID:  annotPkgID,
		annotTagRef: annotTagRef,
	}
}

func (b *birAnnotationAttachmentImpl) GetAnnotPkgID() elements.PackageID {
	return b.annotPkgID
}

func (b *birAnnotationAttachmentImpl) GetAnnotTagRef() Name {
	return b.annotTagRef
}

func (b *birAnnotationAttachmentImpl) Accept(visitor BIRVisitor) {
	visitor.VisitAnnotationAttachment(b)
}

type BIRConstAnnotationAttachment interface {
	BIRAnnotationAttachment
	GetAnnotValue() ConstValue
}

type birConstAnnotationAttachmentImpl struct {
	*birAnnotationAttachmentImpl
	annotValue ConstValue
}

func NewBIRConstAnnotationAttachment(pos diagnostics.Location, annotPkgID elements.PackageID, annotTagRef Name, annotValue ConstValue) BIRConstAnnotationAttachment {
	return &birConstAnnotationAttachmentImpl{
		birAnnotationAttachmentImpl: NewBIRAnnotationAttachment(pos, annotPkgID, annotTagRef).(*birAnnotationAttachmentImpl),
		annotValue:                  annotValue,
	}
}

func (b *birConstAnnotationAttachmentImpl) GetAnnotValue() ConstValue {
	return b.annotValue
}

func (b *birConstAnnotationAttachmentImpl) Accept(visitor BIRVisitor) {
	visitor.VisitConstAnnotationAttachment(b)
}

type BIRLockDetailsHolder interface {
	IsEmpty() bool
	RemoveLastLock()
	GetLock(index int) BIRTerminatorLock
	AddLock(lock BIRTerminatorLock)
	Size() int
}

type birLockDetailsHolderImpl struct {
	locks []BIRTerminatorLock
}

func NewBIRLockDetailsHolder() BIRLockDetailsHolder {
	return &birLockDetailsHolderImpl{
		locks: make([]BIRTerminatorLock, 0),
	}
}

func (b *birLockDetailsHolderImpl) IsEmpty() bool {
	return len(b.locks) == 0
}

func (b *birLockDetailsHolderImpl) RemoveLastLock() {
	if len(b.locks) > 0 {
		b.locks = b.locks[:len(b.locks)-1]
	}
}

func (b *birLockDetailsHolderImpl) GetLock(index int) BIRTerminatorLock {
	return b.locks[index]
}

func (b *birLockDetailsHolderImpl) AddLock(lock BIRTerminatorLock) {
	b.locks = append(b.locks, lock)
}

func (b *birLockDetailsHolderImpl) Size() int {
	return len(b.locks)
}

type BIRMappingConstructorEntry interface {
	IsKeyValuePair() bool
}

type BIRMappingConstructorKeyValueEntry interface {
	BIRMappingConstructorEntry
	GetKeyOp() BIROperand
	GetValueOp() BIROperand
}

type birMappingConstructorKeyValueEntryImpl struct {
	keyOp   BIROperand
	valueOp BIROperand
}

func NewBIRMappingConstructorKeyValueEntry(keyOp, valueOp BIROperand) BIRMappingConstructorKeyValueEntry {
	return &birMappingConstructorKeyValueEntryImpl{
		keyOp:   keyOp,
		valueOp: valueOp,
	}
}

func (b *birMappingConstructorKeyValueEntryImpl) IsKeyValuePair() bool {
	return true
}

func (b *birMappingConstructorKeyValueEntryImpl) GetKeyOp() BIROperand {
	return b.keyOp
}

func (b *birMappingConstructorKeyValueEntryImpl) GetValueOp() BIROperand {
	return b.valueOp
}

type BIRMappingConstructorSpreadFieldEntry interface {
	BIRMappingConstructorEntry
	GetExprOp() BIROperand
}

type birMappingConstructorSpreadFieldEntryImpl struct {
	exprOp BIROperand
}

func NewBIRMappingConstructorSpreadFieldEntry(exprOp BIROperand) BIRMappingConstructorSpreadFieldEntry {
	return &birMappingConstructorSpreadFieldEntryImpl{
		exprOp: exprOp,
	}
}

func (b *birMappingConstructorSpreadFieldEntryImpl) IsKeyValuePair() bool {
	return false
}

func (b *birMappingConstructorSpreadFieldEntryImpl) GetExprOp() BIROperand {
	return b.exprOp
}

type BIRListConstructorEntry interface {
	GetExprOp() BIROperand
}

type BIRListConstructorSpreadMemberEntry interface {
	BIRListConstructorEntry
}

type birListConstructorSpreadMemberEntryImpl struct {
	exprOp BIROperand
}

func NewBIRListConstructorSpreadMemberEntry(exprOp BIROperand) BIRListConstructorSpreadMemberEntry {
	return &birListConstructorSpreadMemberEntryImpl{
		exprOp: exprOp,
	}
}

func (b *birListConstructorSpreadMemberEntryImpl) GetExprOp() BIROperand {
	return b.exprOp
}

type BIRListConstructorExprEntry interface {
	BIRListConstructorEntry
}

type birListConstructorExprEntryImpl struct {
	exprOp BIROperand
}

func NewBIRListConstructorExprEntry(exprOp BIROperand) BIRListConstructorExprEntry {
	return &birListConstructorExprEntryImpl{
		exprOp: exprOp,
	}
}

func (b *birListConstructorExprEntryImpl) GetExprOp() BIROperand {
	return b.exprOp
}

type BIRServiceDeclaration interface {
	BIRDocumentableNode
	GetAttachPoint() []string
	GetAttachPointLiteral() string
	GetListenerTypes() []types.BType
	GetGeneratedName() Name
	GetAssociatedClassName() Name
	GetType() types.BType
	GetOrigin() symbols.SymbolOrigin
	GetFlags() int64
}

type birServiceDeclarationImpl struct {
	*birDocumentableNodeImpl
	attachPoint         []string
	attachPointLiteral  string
	listenerTypes       []types.BType
	generatedName       Name
	associatedClassName Name
	bType               types.BType
	origin              symbols.SymbolOrigin
	flags               int64
}

func NewBIRServiceDeclaration(attachPoint []string, attachPointLiteral string, listenerTypes []types.BType, generatedName, associatedClassName Name, bType types.BType, origin symbols.SymbolOrigin, flags int64, location diagnostics.Location) BIRServiceDeclaration {
	return &birServiceDeclarationImpl{
		birDocumentableNodeImpl: newBIRDocumentableNode(location),
		attachPoint:             attachPoint,
		attachPointLiteral:      attachPointLiteral,
		listenerTypes:           listenerTypes,
		generatedName:           generatedName,
		associatedClassName:     associatedClassName,
		bType:                   bType,
		origin:                  origin,
		flags:                   flags,
	}
}

func (b *birServiceDeclarationImpl) GetAttachPoint() []string {
	return b.attachPoint
}

func (b *birServiceDeclarationImpl) GetAttachPointLiteral() string {
	return b.attachPointLiteral
}

func (b *birServiceDeclarationImpl) GetListenerTypes() []types.BType {
	return b.listenerTypes
}

func (b *birServiceDeclarationImpl) GetGeneratedName() Name {
	return b.generatedName
}

func (b *birServiceDeclarationImpl) GetAssociatedClassName() Name {
	return b.associatedClassName
}

func (b *birServiceDeclarationImpl) GetType() types.BType {
	return b.bType
}

func (b *birServiceDeclarationImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birServiceDeclarationImpl) GetFlags() int64 {
	return b.flags
}

func (b *birServiceDeclarationImpl) Accept(visitor BIRVisitor) {
	visitor.VisitServiceDeclaration(b)
}

type BIRNonTerminator interface {
	BIRNode
}

type BIRTerminator interface {
	BIRNode
}

type BIRTerminatorLock interface {
	BIRTerminator
}

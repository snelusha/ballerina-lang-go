package model

import (
	"ballerina-lang-go/compiler/util"
	"ballerina-lang-go/model/elements"
	"ballerina-lang-go/model/symbols"
	"ballerina-lang-go/semantics/model/types"
	"ballerina-lang-go/tools/diagnostics"
)

type BIRNode interface {
	Accept(visitor BIRVisitor)
}

type BIRPackage interface {
	BIRNode
	GetPackageID() *elements.PackageID
	GetImportModules() []BIRImportModule
	GetTypeDefs() []BIRTypeDefinition
	GetGlobalVars() []BIRGlobalVariableDcl
	GetImportedGlobalVarsDummyVarDcls() []BIRGlobalVariableDcl
	GetFunctions() []BIRFunction
	GetAnnotations() []BIRAnnotation
	GetConstants() []BIRConstant
	GetServiceDecls() []BIRServiceDeclaration
	GetIsListenerAvailable() bool
	SetIsListenerAvailable(available bool)
	GetRecordDefaultValueMap() map[string]map[string]string
	SetRecordDefaultValueMap(m map[string]map[string]string)
}

type birPackageImpl struct {
	pos                            diagnostics.Location
	packageID                      *elements.PackageID
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

func NewBIRPackage(pos diagnostics.Location, org, pkgName, name, version, sourceFileName *util.Name,
	sourceRoot string, skipTest, isTestPkg bool) BIRPackage {
	return &birPackageImpl{
		pos: pos,
		packageID: elements.NewPackageID(
			org.Value,
			pkgName.Value,
			name.Value,
			version.Value,
			sourceFileName.Value,
			sourceRoot,
			isTestPkg,
			skipTest,
		),
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

func (b *birPackageImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRPackage(b)
}

func (b *birPackageImpl) GetPackageID() *elements.PackageID {
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

func (b *birPackageImpl) GetIsListenerAvailable() bool {
	return b.isListenerAvailable
}

func (b *birPackageImpl) SetIsListenerAvailable(available bool) {
	b.isListenerAvailable = available
}

func (b *birPackageImpl) GetRecordDefaultValueMap() map[string]map[string]string {
	return b.recordDefaultValueMap
}

func (b *birPackageImpl) SetRecordDefaultValueMap(m map[string]map[string]string) {
	b.recordDefaultValueMap = m
}

type BIRImportModule interface {
	BIRNode
	GetPackageID() *elements.PackageID
}

type birImportModuleImpl struct {
	pos       diagnostics.Location
	packageID *elements.PackageID
}

func NewBIRImportModule(pos diagnostics.Location, org, name, version *util.Name) BIRImportModule {
	return &birImportModuleImpl{
		pos:       pos,
		packageID: elements.NewPackageIDShort(org.Value, name.Value, version.Value),
	}
}

func (b *birImportModuleImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRImportModule(b)
}

func (b *birImportModuleImpl) GetPackageID() *elements.PackageID {
	return b.packageID
}

type BIRVariableDcl interface {
	BIRDocumentableNode
	GetType() types.BType
	SetType(t types.BType)
	GetName() *util.Name
	GetOriginalName() *util.Name
	GetMetaVarName() string
	GetJvmVarName() string
	GetKind() VarKind
	GetScope() VarScope
	GetIgnoreVariable() bool
	SetIgnoreVariable(ignore bool)
	GetEndBB() BIRBasicBlock
	SetEndBB(bb BIRBasicBlock)
	GetStartBB() BIRBasicBlock
	SetStartBB(bb BIRBasicBlock)
	GetInsOffset() int
	SetInsOffset(offset int)
	GetOnlyUsedInSingleBB() bool
	SetOnlyUsedInSingleBB(single bool)
	GetInitialized() bool
	SetInitialized(init bool)
	GetInsScope() *BirScope
	SetInsScope(scope *BirScope)
}

type birVariableDclImpl struct {
	pos                   diagnostics.Location
	typ                   types.BType
	name                  *util.Name
	originalName          *util.Name
	metaVarName           string
	jvmVarName            string
	kind                  VarKind
	scope                 VarScope
	ignoreVariable        bool
	endBB                 BIRBasicBlock
	startBB               BIRBasicBlock
	insOffset             int
	onlyUsedInSingleBB    bool
	initialized           bool
	insScope              *BirScope
	markdownDocAttachment elements.MarkdownDocAttachment
}

func NewBIRVariableDcl(pos diagnostics.Location, typ types.BType, name, originalName *util.Name,
	scope VarScope, kind VarKind, metaVarName string) BIRVariableDcl {
	jvmVarName := name.Value
	return &birVariableDclImpl{
		pos:          pos,
		typ:          typ,
		name:         name,
		originalName: originalName,
		scope:        scope,
		kind:         kind,
		metaVarName:  metaVarName,
		jvmVarName:   jvmVarName,
	}
}

func NewBIRVariableDclSimple(pos diagnostics.Location, typ types.BType, name *util.Name,
	scope VarScope, kind VarKind, metaVarName string) BIRVariableDcl {
	return NewBIRVariableDcl(pos, typ, name, name, scope, kind, metaVarName)
}

func (b *birVariableDclImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRVariableDcl(b)
}

func (b *birVariableDclImpl) GetType() types.BType {
	return b.typ
}

func (b *birVariableDclImpl) SetType(t types.BType) {
	b.typ = t
}

func (b *birVariableDclImpl) GetName() *util.Name {
	return b.name
}

func (b *birVariableDclImpl) GetOriginalName() *util.Name {
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

func (b *birVariableDclImpl) GetIgnoreVariable() bool {
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

func (b *birVariableDclImpl) GetOnlyUsedInSingleBB() bool {
	return b.onlyUsedInSingleBB
}

func (b *birVariableDclImpl) SetOnlyUsedInSingleBB(single bool) {
	b.onlyUsedInSingleBB = single
}

func (b *birVariableDclImpl) GetInitialized() bool {
	return b.initialized
}

func (b *birVariableDclImpl) SetInitialized(init bool) {
	b.initialized = init
}

func (b *birVariableDclImpl) GetInsScope() *BirScope {
	return b.insScope
}

func (b *birVariableDclImpl) SetInsScope(scope *BirScope) {
	b.insScope = scope
}

func (b *birVariableDclImpl) GetMarkdownDocAttachment() elements.MarkdownDocAttachment {
	return b.markdownDocAttachment
}

func (b *birVariableDclImpl) SetMarkdownDocAttachment(doc elements.MarkdownDocAttachment) {
	b.markdownDocAttachment = doc
}

type BIRParameter interface {
	BIRNode
	GetName() *util.Name
	GetFlags() int64
	GetAnnotAttachments() []BIRAnnotationAttachment
	SetAnnotAttachments(attachments []BIRAnnotationAttachment)
}

type birParameterImpl struct {
	pos              diagnostics.Location
	name             *util.Name
	flags            int64
	annotAttachments []BIRAnnotationAttachment
}

func NewBIRParameter(pos diagnostics.Location, name *util.Name, flags int64) BIRParameter {
	return &birParameterImpl{
		pos:              pos,
		name:             name,
		flags:            flags,
		annotAttachments: make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birParameterImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRParameter(b)
}

func (b *birParameterImpl) GetName() *util.Name {
	return b.name
}

func (b *birParameterImpl) GetFlags() int64 {
	return b.flags
}

func (b *birParameterImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birParameterImpl) SetAnnotAttachments(attachments []BIRAnnotationAttachment) {
	b.annotAttachments = attachments
}

type BIRGlobalVariableDcl interface {
	BIRVariableDcl
	GetFlags() int64
	GetPkgID() *elements.PackageID
	GetOrigin() symbols.SymbolOrigin
	GetAnnotAttachments() []BIRAnnotationAttachment
	SetAnnotAttachments(attachments []BIRAnnotationAttachment)
}

type birGlobalVariableDclImpl struct {
	*birVariableDclImpl
	flags            int64
	pkgID            *elements.PackageID
	origin           symbols.SymbolOrigin
	annotAttachments []BIRAnnotationAttachment
}

func NewBIRGlobalVariableDcl(pos diagnostics.Location, flags int64, typ types.BType,
	pkgID *elements.PackageID, name, originalName *util.Name, scope VarScope, kind VarKind,
	metaVarName string, origin symbols.SymbolOrigin) BIRGlobalVariableDcl {
	return &birGlobalVariableDclImpl{
		birVariableDclImpl: NewBIRVariableDcl(pos, typ, name, originalName, scope, kind, metaVarName).(*birVariableDclImpl),
		flags:              flags,
		pkgID:              pkgID,
		origin:             origin,
		annotAttachments:   make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birGlobalVariableDclImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRVariableDcl(b)
}

func (b *birGlobalVariableDclImpl) GetFlags() int64 {
	return b.flags
}

func (b *birGlobalVariableDclImpl) GetPkgID() *elements.PackageID {
	return b.pkgID
}

func (b *birGlobalVariableDclImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birGlobalVariableDclImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birGlobalVariableDclImpl) SetAnnotAttachments(attachments []BIRAnnotationAttachment) {
	b.annotAttachments = attachments
}

type BIRFunctionParameter interface {
	BIRVariableDcl
	GetHasDefaultExpr() bool
	GetIsPathParameter() bool
	SetIsPathParameter(isPath bool)
}

type birFunctionParameterImpl struct {
	*birVariableDclImpl
	hasDefaultExpr  bool
	isPathParameter bool
}

func NewBIRFunctionParameter(pos diagnostics.Location, typ types.BType, name *util.Name,
	scope VarScope, kind VarKind, metaVarName string, hasDefaultExpr, isPathParameter bool) BIRFunctionParameter {
	return &birFunctionParameterImpl{
		birVariableDclImpl: NewBIRVariableDcl(pos, typ, name, name, scope, kind, metaVarName).(*birVariableDclImpl),
		hasDefaultExpr:     hasDefaultExpr,
		isPathParameter:    isPathParameter,
	}
}

func (b *birFunctionParameterImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRFunctionParameter(b)
}

func (b *birFunctionParameterImpl) GetHasDefaultExpr() bool {
	return b.hasDefaultExpr
}

func (b *birFunctionParameterImpl) GetIsPathParameter() bool {
	return b.isPathParameter
}

func (b *birFunctionParameterImpl) SetIsPathParameter(isPath bool) {
	b.isPathParameter = isPath
}

type BIRFunction interface {
	BIRDocumentableNode
	types.NamedNode
	GetOriginalName() *util.Name
	GetFlags() int64
	GetOrigin() symbols.SymbolOrigin
	GetType() types.BInvokableType
	GetRequiredParams() []BIRParameter
	SetRequiredParams(params []BIRParameter)
	GetReceiver() BIRVariableDcl
	SetReceiver(receiver BIRVariableDcl)
	GetRestParam() BIRParameter
	SetRestParam(param BIRParameter)
	GetArgsCount() int
	SetArgsCount(count int)
	GetLocalVars() []BIRVariableDcl
	SetLocalVars(vars []BIRVariableDcl)
	GetReturnVariable() BIRVariableDcl
	SetReturnVariable(v BIRVariableDcl)
	GetParameters() []BIRFunctionParameter
	SetParameters(params []BIRFunctionParameter)
	GetBasicBlocks() []BIRBasicBlock
	SetBasicBlocks(blocks []BIRBasicBlock)
	GetErrorTable() []BIRErrorEntry
	SetErrorTable(table []BIRErrorEntry)
	GetWorkerName() *util.Name
	GetWorkerChannels() []*ChannelDetails
	SetWorkerChannels(channels []*ChannelDetails)
	GetAnnotAttachments() []BIRAnnotationAttachment
	SetAnnotAttachments(attachments []BIRAnnotationAttachment)
	GetAnnotAttachmentsOnExternal() []BIRAnnotationAttachment
	SetAnnotAttachmentsOnExternal(attachments []BIRAnnotationAttachment)
	GetReturnTypeAnnots() []BIRAnnotationAttachment
	SetReturnTypeAnnots(annots []BIRAnnotationAttachment)
	GetDependentGlobalVars() []BIRGlobalVariableDcl
	SetDependentGlobalVars(vars []BIRGlobalVariableDcl)
	GetPathParams() []BIRVariableDcl
	SetPathParams(params []BIRVariableDcl)
	GetRestPathParam() BIRVariableDcl
	SetRestPathParam(param BIRVariableDcl)
	GetResourcePath() []*util.Name
	SetResourcePath(path []*util.Name)
	GetResourcePathSegmentPosList() []diagnostics.Location
	SetResourcePathSegmentPosList(posList []diagnostics.Location)
	GetAccessor() *util.Name
	SetAccessor(accessor *util.Name)
	GetPathSegmentTypeList() []types.BType
	SetPathSegmentTypeList(typeList []types.BType)
	GetHasWorkers() bool
	SetHasWorkers(hasWorkers bool)
}

type birFunctionImpl struct {
	pos                        diagnostics.Location
	name                       *util.Name
	originalName               *util.Name
	flags                      int64
	origin                     symbols.SymbolOrigin
	typ                        types.BInvokableType
	requiredParams             []BIRParameter
	receiver                   BIRVariableDcl
	restParam                  BIRParameter
	argsCount                  int
	localVars                  []BIRVariableDcl
	returnVariable             BIRVariableDcl
	parameters                 []BIRFunctionParameter
	basicBlocks                []BIRBasicBlock
	errorTable                 []BIRErrorEntry
	workerName                 *util.Name
	workerChannels             []*ChannelDetails
	annotAttachments           []BIRAnnotationAttachment
	annotAttachmentsOnExternal []BIRAnnotationAttachment
	returnTypeAnnots           []BIRAnnotationAttachment
	dependentGlobalVars        []BIRGlobalVariableDcl
	pathParams                 []BIRVariableDcl
	restPathParam              BIRVariableDcl
	resourcePath               []*util.Name
	resourcePathSegmentPosList []diagnostics.Location
	accessor                   *util.Name
	pathSegmentTypeList        []types.BType
	hasWorkers                 bool
	markdownDocAttachment      elements.MarkdownDocAttachment
}

func NewBIRFunction(pos diagnostics.Location, name, originalName *util.Name, flags int64,
	typ types.BInvokableType, workerName *util.Name, sendInsCount int, origin symbols.SymbolOrigin) BIRFunction {
	return &birFunctionImpl{
		pos:                 pos,
		name:                name,
		originalName:        originalName,
		flags:               flags,
		typ:                 typ,
		localVars:           make([]BIRVariableDcl, 0),
		parameters:          make([]BIRFunctionParameter, 0),
		requiredParams:      make([]BIRParameter, 0),
		basicBlocks:         make([]BIRBasicBlock, 0),
		errorTable:          make([]BIRErrorEntry, 0),
		workerName:          workerName,
		workerChannels:      make([]*ChannelDetails, sendInsCount),
		annotAttachments:    make([]BIRAnnotationAttachment, 0),
		returnTypeAnnots:    make([]BIRAnnotationAttachment, 0),
		origin:              origin,
		dependentGlobalVars: make([]BIRGlobalVariableDcl, 0),
	}
}

func (b *birFunctionImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRFunction(b)
}

func (b *birFunctionImpl) GetName() string {
	return b.name.Value
}

func (b *birFunctionImpl) GetOriginalName() *util.Name {
	return b.originalName
}

func (b *birFunctionImpl) GetFlags() int64 {
	return b.flags
}

func (b *birFunctionImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birFunctionImpl) GetType() types.BInvokableType {
	return b.typ
}

func (b *birFunctionImpl) GetRequiredParams() []BIRParameter {
	return b.requiredParams
}

func (b *birFunctionImpl) SetRequiredParams(params []BIRParameter) {
	b.requiredParams = params
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

func (b *birFunctionImpl) SetRestParam(param BIRParameter) {
	b.restParam = param
}

func (b *birFunctionImpl) GetArgsCount() int {
	return b.argsCount
}

func (b *birFunctionImpl) SetArgsCount(count int) {
	b.argsCount = count
}

func (b *birFunctionImpl) GetLocalVars() []BIRVariableDcl {
	return b.localVars
}

func (b *birFunctionImpl) SetLocalVars(vars []BIRVariableDcl) {
	b.localVars = vars
}

func (b *birFunctionImpl) GetReturnVariable() BIRVariableDcl {
	return b.returnVariable
}

func (b *birFunctionImpl) SetReturnVariable(v BIRVariableDcl) {
	b.returnVariable = v
}

func (b *birFunctionImpl) GetParameters() []BIRFunctionParameter {
	return b.parameters
}

func (b *birFunctionImpl) SetParameters(params []BIRFunctionParameter) {
	b.parameters = params
}

func (b *birFunctionImpl) GetBasicBlocks() []BIRBasicBlock {
	return b.basicBlocks
}

func (b *birFunctionImpl) SetBasicBlocks(blocks []BIRBasicBlock) {
	b.basicBlocks = blocks
}

func (b *birFunctionImpl) GetErrorTable() []BIRErrorEntry {
	return b.errorTable
}

func (b *birFunctionImpl) SetErrorTable(table []BIRErrorEntry) {
	b.errorTable = table
}

func (b *birFunctionImpl) GetWorkerName() *util.Name {
	return b.workerName
}

func (b *birFunctionImpl) GetWorkerChannels() []*ChannelDetails {
	return b.workerChannels
}

func (b *birFunctionImpl) SetWorkerChannels(channels []*ChannelDetails) {
	b.workerChannels = channels
}

func (b *birFunctionImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birFunctionImpl) SetAnnotAttachments(attachments []BIRAnnotationAttachment) {
	b.annotAttachments = attachments
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

func (b *birFunctionImpl) SetReturnTypeAnnots(annots []BIRAnnotationAttachment) {
	b.returnTypeAnnots = annots
}

func (b *birFunctionImpl) GetDependentGlobalVars() []BIRGlobalVariableDcl {
	return b.dependentGlobalVars
}

func (b *birFunctionImpl) SetDependentGlobalVars(vars []BIRGlobalVariableDcl) {
	b.dependentGlobalVars = vars
}

func (b *birFunctionImpl) GetPathParams() []BIRVariableDcl {
	return b.pathParams
}

func (b *birFunctionImpl) SetPathParams(params []BIRVariableDcl) {
	b.pathParams = params
}

func (b *birFunctionImpl) GetRestPathParam() BIRVariableDcl {
	return b.restPathParam
}

func (b *birFunctionImpl) SetRestPathParam(param BIRVariableDcl) {
	b.restPathParam = param
}

func (b *birFunctionImpl) GetResourcePath() []*util.Name {
	return b.resourcePath
}

func (b *birFunctionImpl) SetResourcePath(path []*util.Name) {
	b.resourcePath = path
}

func (b *birFunctionImpl) GetResourcePathSegmentPosList() []diagnostics.Location {
	return b.resourcePathSegmentPosList
}

func (b *birFunctionImpl) SetResourcePathSegmentPosList(posList []diagnostics.Location) {
	b.resourcePathSegmentPosList = posList
}

func (b *birFunctionImpl) GetAccessor() *util.Name {
	return b.accessor
}

func (b *birFunctionImpl) SetAccessor(accessor *util.Name) {
	b.accessor = accessor
}

func (b *birFunctionImpl) GetPathSegmentTypeList() []types.BType {
	return b.pathSegmentTypeList
}

func (b *birFunctionImpl) SetPathSegmentTypeList(typeList []types.BType) {
	b.pathSegmentTypeList = typeList
}

func (b *birFunctionImpl) GetHasWorkers() bool {
	return b.hasWorkers
}

func (b *birFunctionImpl) SetHasWorkers(hasWorkers bool) {
	b.hasWorkers = hasWorkers
}

func (b *birFunctionImpl) GetMarkdownDocAttachment() elements.MarkdownDocAttachment {
	return b.markdownDocAttachment
}

func (b *birFunctionImpl) SetMarkdownDocAttachment(doc elements.MarkdownDocAttachment) {
	b.markdownDocAttachment = doc
}

type BIRBasicBlock interface {
	BIRNode
	GetNumber() int
	GetID() *util.Name
	GetInstructions() []BIRNonTerminator
	SetInstructions(instructions []BIRNonTerminator)
	GetTerminator() BIRTerminator
	SetTerminator(terminator BIRTerminator)
}

type birBasicBlockImpl struct {
	pos          diagnostics.Location
	number       int
	id           *util.Name
	instructions []BIRNonTerminator
	terminator   BIRTerminator
}

const birBasicBlockPrefix = "bb"

func NewBIRBasicBlock(id *util.Name, number int) BIRBasicBlock {
	return &birBasicBlockImpl{
		number:       number,
		id:           id,
		instructions: make([]BIRNonTerminator, 0),
	}
}

func NewBIRBasicBlockWithNumber(number int) BIRBasicBlock {
	return NewBIRBasicBlock(util.NewName(birBasicBlockPrefix+string(rune(number))), number)
}

func NewBIRBasicBlockWithPrefix(idPrefix string, number int) BIRBasicBlock {
	return NewBIRBasicBlock(util.NewName(idPrefix+string(rune(number))), number)
}

func (b *birBasicBlockImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRBasicBlock(b)
}

func (b *birBasicBlockImpl) GetNumber() int {
	return b.number
}

func (b *birBasicBlockImpl) GetID() *util.Name {
	return b.id
}

func (b *birBasicBlockImpl) GetInstructions() []BIRNonTerminator {
	return b.instructions
}

func (b *birBasicBlockImpl) SetInstructions(instructions []BIRNonTerminator) {
	b.instructions = instructions
}

func (b *birBasicBlockImpl) GetTerminator() BIRTerminator {
	return b.terminator
}

func (b *birBasicBlockImpl) SetTerminator(terminator BIRTerminator) {
	b.terminator = terminator
}

type BIRTypeDefinition interface {
	BIRDocumentableNode
	types.NamedNode
	GetOriginalName() *util.Name
	GetInternalName() *util.Name
	GetAttachedFuncs() []BIRFunction
	SetAttachedFuncs(funcs []BIRFunction)
	GetFlags() int64
	GetType() types.BType
	GetIsBuiltin() bool
	GetReferencedTypes() []types.BType
	SetReferencedTypes(types []types.BType)
	GetReferenceType() types.BType
	SetReferenceType(typ types.BType)
	GetOrigin() symbols.SymbolOrigin
	GetAnnotAttachments() []BIRAnnotationAttachment
	SetAnnotAttachments(attachments []BIRAnnotationAttachment)
	GetIndex() int
	SetIndex(index int)
}

type birTypeDefinitionImpl struct {
	pos                   diagnostics.Location
	name                  *util.Name
	originalName          *util.Name
	internalName          *util.Name
	attachedFuncs         []BIRFunction
	flags                 int64
	typ                   types.BType
	isBuiltin             bool
	referencedTypes       []types.BType
	referenceType         types.BType
	origin                symbols.SymbolOrigin
	annotAttachments      []BIRAnnotationAttachment
	index                 int
	markdownDocAttachment elements.MarkdownDocAttachment
}

func NewBIRTypeDefinition(pos diagnostics.Location, internalName *util.Name, flags int64, isBuiltin bool,
	typ types.BType, attachedFuncs []BIRFunction, origin symbols.SymbolOrigin,
	name, originalName *util.Name) BIRTypeDefinition {
	return &birTypeDefinitionImpl{
		pos:              pos,
		internalName:     internalName,
		flags:            flags,
		isBuiltin:        isBuiltin,
		typ:              typ,
		attachedFuncs:    attachedFuncs,
		referencedTypes:  make([]types.BType, 0),
		origin:           origin,
		name:             name,
		originalName:     originalName,
		annotAttachments: make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birTypeDefinitionImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRTypeDefinition(b)
}

func (b *birTypeDefinitionImpl) GetName() string {
	return b.name.Value
}

func (b *birTypeDefinitionImpl) GetOriginalName() *util.Name {
	return b.originalName
}

func (b *birTypeDefinitionImpl) GetInternalName() *util.Name {
	return b.internalName
}

func (b *birTypeDefinitionImpl) GetAttachedFuncs() []BIRFunction {
	return b.attachedFuncs
}

func (b *birTypeDefinitionImpl) SetAttachedFuncs(funcs []BIRFunction) {
	b.attachedFuncs = funcs
}

func (b *birTypeDefinitionImpl) GetFlags() int64 {
	return b.flags
}

func (b *birTypeDefinitionImpl) GetType() types.BType {
	return b.typ
}

func (b *birTypeDefinitionImpl) GetIsBuiltin() bool {
	return b.isBuiltin
}

func (b *birTypeDefinitionImpl) GetReferencedTypes() []types.BType {
	return b.referencedTypes
}

func (b *birTypeDefinitionImpl) SetReferencedTypes(types []types.BType) {
	b.referencedTypes = types
}

func (b *birTypeDefinitionImpl) GetReferenceType() types.BType {
	return b.referenceType
}

func (b *birTypeDefinitionImpl) SetReferenceType(typ types.BType) {
	b.referenceType = typ
}

func (b *birTypeDefinitionImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birTypeDefinitionImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birTypeDefinitionImpl) SetAnnotAttachments(attachments []BIRAnnotationAttachment) {
	b.annotAttachments = attachments
}

func (b *birTypeDefinitionImpl) GetIndex() int {
	return b.index
}

func (b *birTypeDefinitionImpl) SetIndex(index int) {
	b.index = index
}

func (b *birTypeDefinitionImpl) GetMarkdownDocAttachment() elements.MarkdownDocAttachment {
	return b.markdownDocAttachment
}

func (b *birTypeDefinitionImpl) SetMarkdownDocAttachment(doc elements.MarkdownDocAttachment) {
	b.markdownDocAttachment = doc
}

type BIRErrorEntry interface {
	BIRNode
	GetTrapBB() BIRBasicBlock
	GetEndBB() BIRBasicBlock
	GetErrorOp() BIROperand
	GetTargetBB() BIRBasicBlock
}

type birErrorEntryImpl struct {
	trapBB   BIRBasicBlock
	endBB    BIRBasicBlock
	errorOp  BIROperand
	targetBB BIRBasicBlock
}

func NewBIRErrorEntry(trapBB, endBB BIRBasicBlock, errorOp BIROperand, targetBB BIRBasicBlock) BIRErrorEntry {
	return &birErrorEntryImpl{
		trapBB:   trapBB,
		endBB:    endBB,
		errorOp:  errorOp,
		targetBB: targetBB,
	}
}

func (b *birErrorEntryImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRErrorEntry(b)
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

type ChannelDetails struct {
	Name                string
	ChannelInSameStrand bool
	Send                bool
}

func NewChannelDetails(name string, channelInSameStrand, send bool) *ChannelDetails {
	return &ChannelDetails{
		Name:                name,
		ChannelInSameStrand: channelInSameStrand,
		Send:                send,
	}
}

type BIRAnnotation interface {
	BIRDocumentableNode
	GetName() *util.Name
	GetOriginalName() *util.Name
	GetFlags() int64
	GetOrigin() symbols.SymbolOrigin
	GetAttachPoints() []elements.AttachPoint
	SetAttachPoints(points []elements.AttachPoint)
	GetAnnotationType() types.BType
	GetPackageID() *elements.PackageID
	SetPackageID(pkgID *elements.PackageID)
	GetAnnotAttachments() []BIRAnnotationAttachment
	SetAnnotAttachments(attachments []BIRAnnotationAttachment)
}

type birAnnotationImpl struct {
	pos                   diagnostics.Location
	name                  *util.Name
	originalName          *util.Name
	flags                 int64
	origin                symbols.SymbolOrigin
	attachPoints          []elements.AttachPoint
	annotationType        types.BType
	packageID             *elements.PackageID
	annotAttachments      []BIRAnnotationAttachment
	markdownDocAttachment elements.MarkdownDocAttachment
}

func NewBIRAnnotation(pos diagnostics.Location, name, originalName *util.Name, flags int64,
	points []elements.AttachPoint, annotationType types.BType, origin symbols.SymbolOrigin) BIRAnnotation {
	return &birAnnotationImpl{
		pos:              pos,
		name:             name,
		originalName:     originalName,
		flags:            flags,
		attachPoints:     points,
		annotationType:   annotationType,
		origin:           origin,
		annotAttachments: make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birAnnotationImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRAnnotation(b)
}

func (b *birAnnotationImpl) GetName() *util.Name {
	return b.name
}

func (b *birAnnotationImpl) GetOriginalName() *util.Name {
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

func (b *birAnnotationImpl) SetAttachPoints(points []elements.AttachPoint) {
	b.attachPoints = points
}

func (b *birAnnotationImpl) GetAnnotationType() types.BType {
	return b.annotationType
}

func (b *birAnnotationImpl) GetPackageID() *elements.PackageID {
	return b.packageID
}

func (b *birAnnotationImpl) SetPackageID(pkgID *elements.PackageID) {
	b.packageID = pkgID
}

func (b *birAnnotationImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birAnnotationImpl) SetAnnotAttachments(attachments []BIRAnnotationAttachment) {
	b.annotAttachments = attachments
}

func (b *birAnnotationImpl) GetMarkdownDocAttachment() elements.MarkdownDocAttachment {
	return b.markdownDocAttachment
}

func (b *birAnnotationImpl) SetMarkdownDocAttachment(doc elements.MarkdownDocAttachment) {
	b.markdownDocAttachment = doc
}

type BIRConstant interface {
	BIRDocumentableNode
	GetName() *util.Name
	GetFlags() int64
	GetType() types.BType
	GetConstValue() *ConstValue
	GetOrigin() symbols.SymbolOrigin
	GetAnnotAttachments() []BIRAnnotationAttachment
	SetAnnotAttachments(attachments []BIRAnnotationAttachment)
}

type birConstantImpl struct {
	pos                   diagnostics.Location
	name                  *util.Name
	flags                 int64
	typ                   types.BType
	constValue            *ConstValue
	origin                symbols.SymbolOrigin
	annotAttachments      []BIRAnnotationAttachment
	markdownDocAttachment elements.MarkdownDocAttachment
}

func NewBIRConstant(pos diagnostics.Location, name *util.Name, flags int64, typ types.BType,
	constValue *ConstValue, origin symbols.SymbolOrigin) BIRConstant {
	return &birConstantImpl{
		pos:              pos,
		name:             name,
		flags:            flags,
		typ:              typ,
		constValue:       constValue,
		origin:           origin,
		annotAttachments: make([]BIRAnnotationAttachment, 0),
	}
}

func (b *birConstantImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRConstant(b)
}

func (b *birConstantImpl) GetName() *util.Name {
	return b.name
}

func (b *birConstantImpl) GetFlags() int64 {
	return b.flags
}

func (b *birConstantImpl) GetType() types.BType {
	return b.typ
}

func (b *birConstantImpl) GetConstValue() *ConstValue {
	return b.constValue
}

func (b *birConstantImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birConstantImpl) GetAnnotAttachments() []BIRAnnotationAttachment {
	return b.annotAttachments
}

func (b *birConstantImpl) SetAnnotAttachments(attachments []BIRAnnotationAttachment) {
	b.annotAttachments = attachments
}

func (b *birConstantImpl) GetMarkdownDocAttachment() elements.MarkdownDocAttachment {
	return b.markdownDocAttachment
}

func (b *birConstantImpl) SetMarkdownDocAttachment(doc elements.MarkdownDocAttachment) {
	b.markdownDocAttachment = doc
}

type BIRAnnotationAttachment interface {
	BIRNode
	GetAnnotPkgID() *elements.PackageID
	GetAnnotTagRef() *util.Name
}

type birAnnotationAttachmentImpl struct {
	pos         diagnostics.Location
	annotPkgID  *elements.PackageID
	annotTagRef *util.Name
}

func NewBIRAnnotationAttachment(pos diagnostics.Location, annotPkgID *elements.PackageID,
	annotTagRef *util.Name) BIRAnnotationAttachment {
	return &birAnnotationAttachmentImpl{
		pos:         pos,
		annotPkgID:  annotPkgID,
		annotTagRef: annotTagRef,
	}
}

func (b *birAnnotationAttachmentImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRAnnotationAttachment(b)
}

func (b *birAnnotationAttachmentImpl) GetAnnotPkgID() *elements.PackageID {
	return b.annotPkgID
}

func (b *birAnnotationAttachmentImpl) GetAnnotTagRef() *util.Name {
	return b.annotTagRef
}

type BIRConstAnnotationAttachment interface {
	BIRAnnotationAttachment
	GetAnnotValue() *ConstValue
}

type birConstAnnotationAttachmentImpl struct {
	*birAnnotationAttachmentImpl
	annotValue *ConstValue
}

func NewBIRConstAnnotationAttachment(pos diagnostics.Location, annotPkgID *elements.PackageID,
	annotTagRef *util.Name, annotValue *ConstValue) BIRConstAnnotationAttachment {
	return &birConstAnnotationAttachmentImpl{
		birAnnotationAttachmentImpl: NewBIRAnnotationAttachment(pos, annotPkgID, annotTagRef).(*birAnnotationAttachmentImpl),
		annotValue:                  annotValue,
	}
}

func (b *birConstAnnotationAttachmentImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRConstAnnotationAttachment(b)
}

func (b *birConstAnnotationAttachmentImpl) GetAnnotValue() *ConstValue {
	return b.annotValue
}

type ConstValue struct {
	Type  types.BType
	Value any
}

func NewConstValue(value any, typ types.BType) *ConstValue {
	return &ConstValue{
		Type:  typ,
		Value: value,
	}
}

type BIRDocumentableNode interface {
	BIRNode
	GetMarkdownDocAttachment() elements.MarkdownDocAttachment
	SetMarkdownDocAttachment(doc elements.MarkdownDocAttachment)
}

type BIRLockDetailsHolder struct {
	locks []BIRTerminator
}

func NewBIRLockDetailsHolder() *BIRLockDetailsHolder {
	return &BIRLockDetailsHolder{
		locks: make([]BIRTerminator, 0),
	}
}

func (h *BIRLockDetailsHolder) IsEmpty() bool {
	return len(h.locks) == 0
}

func (h *BIRLockDetailsHolder) RemoveLastLock() {
	if len(h.locks) > 0 {
		h.locks = h.locks[:len(h.locks)-1]
	}
}

func (h *BIRLockDetailsHolder) GetLock(index int) BIRTerminator {
	return h.locks[index]
}

func (h *BIRLockDetailsHolder) AddLock(lock BIRTerminator) {
	h.locks = append(h.locks, lock)
}

func (h *BIRLockDetailsHolder) Size() int {
	return len(h.locks)
}

type BIRMappingConstructorEntry interface {
	IsKeyValuePair() bool
}

type BIRMappingConstructorKeyValueEntry struct {
	KeyOp   BIROperand
	ValueOp BIROperand
}

func NewBIRMappingConstructorKeyValueEntry(keyOp, valueOp BIROperand) *BIRMappingConstructorKeyValueEntry {
	return &BIRMappingConstructorKeyValueEntry{
		KeyOp:   keyOp,
		ValueOp: valueOp,
	}
}

func (e *BIRMappingConstructorKeyValueEntry) IsKeyValuePair() bool {
	return true
}

type BIRMappingConstructorSpreadFieldEntry struct {
	ExprOp BIROperand
}

func NewBIRMappingConstructorSpreadFieldEntry(exprOp BIROperand) *BIRMappingConstructorSpreadFieldEntry {
	return &BIRMappingConstructorSpreadFieldEntry{
		ExprOp: exprOp,
	}
}

func (e *BIRMappingConstructorSpreadFieldEntry) IsKeyValuePair() bool {
	return false
}

type BIRListConstructorEntry interface {
	GetExprOp() BIROperand
}

type BIRListConstructorSpreadMemberEntry struct {
	ExprOp BIROperand
}

func NewBIRListConstructorSpreadMemberEntry(exprOp BIROperand) *BIRListConstructorSpreadMemberEntry {
	return &BIRListConstructorSpreadMemberEntry{
		ExprOp: exprOp,
	}
}

func (e *BIRListConstructorSpreadMemberEntry) GetExprOp() BIROperand {
	return e.ExprOp
}

type BIRListConstructorExprEntry struct {
	ExprOp BIROperand
}

func NewBIRListConstructorExprEntry(exprOp BIROperand) *BIRListConstructorExprEntry {
	return &BIRListConstructorExprEntry{
		ExprOp: exprOp,
	}
}

func (e *BIRListConstructorExprEntry) GetExprOp() BIROperand {
	return e.ExprOp
}

type BIRServiceDeclaration interface {
	BIRDocumentableNode
	GetAttachPoint() []string
	GetAttachPointLiteral() string
	GetListenerTypes() []types.BType
	GetGeneratedName() *util.Name
	GetAssociatedClassName() *util.Name
	GetType() types.BType
	GetOrigin() symbols.SymbolOrigin
	GetFlags() int64
}

type birServiceDeclarationImpl struct {
	pos                   diagnostics.Location
	attachPoint           []string
	attachPointLiteral    string
	listenerTypes         []types.BType
	generatedName         *util.Name
	associatedClassName   *util.Name
	typ                   types.BType
	origin                symbols.SymbolOrigin
	flags                 int64
	markdownDocAttachment elements.MarkdownDocAttachment
}

func NewBIRServiceDeclaration(attachPoint []string, attachPointLiteral string, listenerTypes []types.BType,
	generatedName, associatedClassName *util.Name, typ types.BType, origin symbols.SymbolOrigin,
	flags int64, location diagnostics.Location) BIRServiceDeclaration {
	return &birServiceDeclarationImpl{
		pos:                 location,
		attachPoint:         attachPoint,
		attachPointLiteral:  attachPointLiteral,
		listenerTypes:       listenerTypes,
		generatedName:       generatedName,
		associatedClassName: associatedClassName,
		typ:                 typ,
		origin:              origin,
		flags:               flags,
	}
}

func (b *birServiceDeclarationImpl) Accept(visitor BIRVisitor) {
	visitor.VisitBIRServiceDeclaration(b)
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

func (b *birServiceDeclarationImpl) GetGeneratedName() *util.Name {
	return b.generatedName
}

func (b *birServiceDeclarationImpl) GetAssociatedClassName() *util.Name {
	return b.associatedClassName
}

func (b *birServiceDeclarationImpl) GetType() types.BType {
	return b.typ
}

func (b *birServiceDeclarationImpl) GetOrigin() symbols.SymbolOrigin {
	return b.origin
}

func (b *birServiceDeclarationImpl) GetFlags() int64 {
	return b.flags
}

func (b *birServiceDeclarationImpl) GetMarkdownDocAttachment() elements.MarkdownDocAttachment {
	return b.markdownDocAttachment
}

func (b *birServiceDeclarationImpl) SetMarkdownDocAttachment(doc elements.MarkdownDocAttachment) {
	b.markdownDocAttachment = doc
}

type BIRNonTerminator interface {
	BIRAbstractInstruction
}

type BIRTerminator interface {
	BIRAbstractInstruction
	GetNextBasicBlocks() []BIRBasicBlock
}

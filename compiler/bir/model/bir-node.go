package model

import (
	"ballerina-lang-go/compiler/semantics/model/elements"
	"ballerina-lang-go/compiler/semantics/model/symbols"
	"ballerina-lang-go/compiler/semantics/model/types"
	"ballerina-lang-go/compiler/util"
	"ballerina-lang-go/tools/diagnostics"
	"fmt"
)

type BIRNode interface {
	GetPos() diagnostics.Location
}

type birNodeBase struct {
	pos diagnostics.Location
}

func (b *birNodeBase) GetPos() diagnostics.Location {
	return b.pos
}

type BIRDocumentableNode interface {
	BIRNode
	GetMarkdownDocAttachment() *elements.MarkdownDocAttachment
	SetMarkdownDocAttachment(doc *elements.MarkdownDocAttachment)
}

type birDocumentableNodeBase struct {
	birNodeBase
	markdownDocAttachment *elements.MarkdownDocAttachment
}

func (b *birDocumentableNodeBase) GetMarkdownDocAttachment() *elements.MarkdownDocAttachment {
	return b.markdownDocAttachment
}

func (b *birDocumentableNodeBase) SetMarkdownDocAttachment(doc *elements.MarkdownDocAttachment) {
	b.markdownDocAttachment = doc
}

type BIRPackage struct {
	birNodeBase
	PackageID                      elements.PackageID
	ImportModules                  map[elements.PackageID]*BIRImportModule
	TypeDefs                       []*BIRTypeDefinition
	GlobalVars                     []*BIRGlobalVariableDcl
	ImportedGlobalVarsDummyVarDcls map[*BIRGlobalVariableDcl]struct{}
	Functions                      []*BIRFunction
	Annotations                    []*BIRAnnotation
	Constants                      []*BIRConstant
	ServiceDecls                   []*BIRServiceDeclaration
	IsListenerAvailable            bool
	RecordDefaultValueMap          map[string]map[string]string
}

func NewBIRPackage(pos diagnostics.Location, org, pkgName, name, version, sourceFileName util.Name, sourceRoot string, skipTest bool) *BIRPackage {
	return NewBIRPackageFull(pos, org, pkgName, name, version, sourceFileName, sourceRoot, skipTest, false)
}

func NewBIRPackageFull(pos diagnostics.Location, org, pkgName, name, version, sourceFileName util.Name, sourceRoot string, skipTest, isTestPkg bool) *BIRPackage {
	return &BIRPackage{
		birNodeBase:                    birNodeBase{pos: pos},
		PackageID:                      elements.NewPackageIDFull(org, pkgName, name, version, sourceFileName, sourceRoot, isTestPkg, skipTest),
		ImportModules:                  make(map[elements.PackageID]*BIRImportModule),
		TypeDefs:                       make([]*BIRTypeDefinition, 0),
		GlobalVars:                     make([]*BIRGlobalVariableDcl, 0),
		ImportedGlobalVarsDummyVarDcls: make(map[*BIRGlobalVariableDcl]struct{}),
		Functions:                      make([]*BIRFunction, 0),
		Annotations:                    make([]*BIRAnnotation, 0),
		Constants:                      make([]*BIRConstant, 0),
		ServiceDecls:                   make([]*BIRServiceDeclaration, 0),
		RecordDefaultValueMap:          make(map[string]map[string]string),
	}
}

// Helper methods for adding items to BIRPackage

func (p *BIRPackage) AddImportModule(module *BIRImportModule) {
	p.ImportModules[module.PackageID] = module
}

func (p *BIRPackage) AddTypeDef(typeDef *BIRTypeDefinition) {
	p.TypeDefs = append(p.TypeDefs, typeDef)
}

func (p *BIRPackage) AddGlobalVar(globalVar *BIRGlobalVariableDcl) {
	p.GlobalVars = append(p.GlobalVars, globalVar)
}

func (p *BIRPackage) AddImportedGlobalVarDummy(globalVar *BIRGlobalVariableDcl) {
	p.ImportedGlobalVarsDummyVarDcls[globalVar] = struct{}{}
}

func (p *BIRPackage) AddFunction(function *BIRFunction) {
	p.Functions = append(p.Functions, function)
}

func (p *BIRPackage) AddAnnotation(annotation *BIRAnnotation) {
	p.Annotations = append(p.Annotations, annotation)
}

func (p *BIRPackage) AddConstant(constant *BIRConstant) {
	p.Constants = append(p.Constants, constant)
}

func (p *BIRPackage) AddServiceDecl(serviceDecl *BIRServiceDeclaration) {
	p.ServiceDecls = append(p.ServiceDecls, serviceDecl)
}

type BIRImportModule struct {
	birNodeBase
	PackageID elements.PackageID
}

func NewBIRImportModule(pos diagnostics.Location, org, name, version util.Name) *BIRImportModule {
	return &BIRImportModule{
		birNodeBase: birNodeBase{pos: pos},
		PackageID:   elements.NewPackageID(org, name, version),
	}
}

type BIRVariableDcl struct {
	birDocumentableNodeBase
	Type               types.BType
	Name               util.Name
	OriginalName       util.Name
	MetaVarName        string
	JvmVarName         string
	Kind               VarKind
	Scope              VarScope
	IgnoreVariable     bool
	EndBB              *BIRBasicBlock
	StartBB            *BIRBasicBlock
	InsOffset          int
	OnlyUsedInSingleBB bool
	Initialized        bool
	InsScope           *BirScope
}

func NewBIRVariableDcl(pos diagnostics.Location, typ types.BType, name, originalName util.Name, scope VarScope, kind VarKind, metaVarName string) *BIRVariableDcl {
	jvmVarName := name.Value
	return &BIRVariableDcl{
		birDocumentableNodeBase: birDocumentableNodeBase{birNodeBase: birNodeBase{pos: pos}},
		Type:                    typ,
		Name:                    name,
		OriginalName:            originalName,
		MetaVarName:             metaVarName,
		JvmVarName:              jvmVarName,
		Kind:                    kind,
		Scope:                   scope,
	}
}

func NewBIRVariableDclSimple(pos diagnostics.Location, typ types.BType, name util.Name, scope VarScope, kind VarKind, metaVarName string) *BIRVariableDcl {
	return NewBIRVariableDcl(pos, typ, name, name, scope, kind, metaVarName)
}

func NewBIRVariableDclMinimal(typ types.BType, name util.Name, scope VarScope, kind VarKind) *BIRVariableDcl {
	return NewBIRVariableDclSimple(nil, typ, name, scope, kind, "")
}

func (v *BIRVariableDcl) String() string {
	return v.Name.String()
}

type BIRParameter struct {
	birNodeBase
	Name             util.Name
	Flags            int64
	AnnotAttachments []*BIRAnnotationAttachment
}

func NewBIRParameter(pos diagnostics.Location, name util.Name, flags int64) *BIRParameter {
	return &BIRParameter{
		birNodeBase:      birNodeBase{pos: pos},
		Name:             name,
		Flags:            flags,
		AnnotAttachments: make([]*BIRAnnotationAttachment, 0),
	}
}

type BIRGlobalVariableDcl struct {
	*BIRVariableDcl
	Flags            int64
	PkgID            elements.PackageID
	Origin           symbols.SymbolOrigin
	AnnotAttachments []*BIRAnnotationAttachment
}

func NewBIRGlobalVariableDcl(pos diagnostics.Location, flags int64, typ types.BType, pkgID elements.PackageID, name, originalName util.Name, scope VarScope, kind VarKind, metaVarName string, origin symbols.SymbolOrigin) *BIRGlobalVariableDcl {
	return &BIRGlobalVariableDcl{
		BIRVariableDcl:   NewBIRVariableDcl(pos, typ, name, originalName, scope, kind, metaVarName),
		Flags:            flags,
		PkgID:            pkgID,
		Origin:           origin,
		AnnotAttachments: make([]*BIRAnnotationAttachment, 0),
	}
}

type BIRFunctionParameter struct {
	*BIRVariableDcl
	HasDefaultExpr  bool
	IsPathParameter bool
}

func NewBIRFunctionParameter(pos diagnostics.Location, typ types.BType, name util.Name, scope VarScope, kind VarKind, metaVarName string, hasDefaultExpr bool) *BIRFunctionParameter {
	return &BIRFunctionParameter{
		BIRVariableDcl: NewBIRVariableDclSimple(pos, typ, name, scope, kind, metaVarName),
		HasDefaultExpr: hasDefaultExpr,
	}
}

func NewBIRFunctionParameterWithPathParam(pos diagnostics.Location, typ types.BType, name util.Name, scope VarScope, kind VarKind, metaVarName string, hasDefaultExpr, isPathParameter bool) *BIRFunctionParameter {
	return &BIRFunctionParameter{
		BIRVariableDcl:  NewBIRVariableDclSimple(pos, typ, name, scope, kind, metaVarName),
		HasDefaultExpr:  hasDefaultExpr,
		IsPathParameter: isPathParameter,
	}
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

func (c *ChannelDetails) String() string {
	return c.Name
}

type BIRFunction struct {
	birDocumentableNodeBase
	Name                       util.Name
	OriginalName               util.Name
	Flags                      int64
	Origin                     symbols.SymbolOrigin
	Type                       types.BInvokableType
	RequiredParams             []*BIRParameter
	Receiver                   *BIRVariableDcl
	RestParam                  *BIRParameter
	ArgsCount                  int
	LocalVars                  []*BIRVariableDcl
	ReturnVariable             *BIRVariableDcl
	Parameters                 []*BIRFunctionParameter
	BasicBlocks                []*BIRBasicBlock
	ErrorTable                 []*BIRErrorEntry
	WorkerName                 util.Name
	WorkerChannels             []*ChannelDetails
	AnnotAttachments           []*BIRAnnotationAttachment
	AnnotAttachmentsOnExternal []*BIRAnnotationAttachment
	ReturnTypeAnnots           []*BIRAnnotationAttachment
	DependentGlobalVars        map[*BIRGlobalVariableDcl]struct{}
	PathParams                 []*BIRVariableDcl
	RestPathParam              *BIRVariableDcl
	ResourcePath               []util.Name
	ResourcePathSegmentPosList []diagnostics.Location
	Accessor                   util.Name
	PathSegmentTypeList        []types.BType
	HasWorkers                 bool
}

func NewBIRFunction(pos diagnostics.Location, name, originalName util.Name, flags int64, origin symbols.SymbolOrigin, typ types.BInvokableType, requiredParams []*BIRParameter, receiver *BIRVariableDcl, restParam *BIRParameter, argsCount int, localVars []*BIRVariableDcl, returnVariable *BIRVariableDcl, parameters []*BIRFunctionParameter, basicBlocks []*BIRBasicBlock, errorTable []*BIRErrorEntry, workerName util.Name, workerChannels []*ChannelDetails, annotAttachments, returnTypeAnnots []*BIRAnnotationAttachment, dependentGlobalVars map[*BIRGlobalVariableDcl]struct{}) *BIRFunction {
	return &BIRFunction{
		birDocumentableNodeBase: birDocumentableNodeBase{birNodeBase: birNodeBase{pos: pos}},
		Name:                    name,
		OriginalName:            originalName,
		Flags:                   flags,
		Origin:                  origin,
		Type:                    typ,
		RequiredParams:          requiredParams,
		Receiver:                receiver,
		RestParam:               restParam,
		ArgsCount:               argsCount,
		LocalVars:               localVars,
		ReturnVariable:          returnVariable,
		Parameters:              parameters,
		BasicBlocks:             basicBlocks,
		ErrorTable:              errorTable,
		WorkerName:              workerName,
		WorkerChannels:          workerChannels,
		AnnotAttachments:        annotAttachments,
		ReturnTypeAnnots:        returnTypeAnnots,
		DependentGlobalVars:     dependentGlobalVars,
	}
}

func NewBIRFunctionSimple(pos diagnostics.Location, name, originalName util.Name, flags int64, typ types.BInvokableType, workerName util.Name, sendInsCount int, origin symbols.SymbolOrigin) *BIRFunction {
	return &BIRFunction{
		birDocumentableNodeBase: birDocumentableNodeBase{birNodeBase: birNodeBase{pos: pos}},
		Name:                    name,
		OriginalName:            originalName,
		Flags:                   flags,
		Type:                    typ,
		LocalVars:               make([]*BIRVariableDcl, 0),
		Parameters:              make([]*BIRFunctionParameter, 0),
		RequiredParams:          make([]*BIRParameter, 0),
		BasicBlocks:             make([]*BIRBasicBlock, 0),
		ErrorTable:              make([]*BIRErrorEntry, 0),
		WorkerName:              workerName,
		WorkerChannels:          make([]*ChannelDetails, sendInsCount),
		AnnotAttachments:        make([]*BIRAnnotationAttachment, 0),
		ReturnTypeAnnots:        make([]*BIRAnnotationAttachment, 0),
		Origin:                  origin,
	}
}

func NewBIRFunctionMinimal(pos diagnostics.Location, name util.Name, flags int64, typ types.BInvokableType, workerName util.Name, sendInsCount int, origin symbols.SymbolOrigin) *BIRFunction {
	return NewBIRFunctionSimple(pos, name, name, flags, typ, workerName, sendInsCount, origin)
}

func (f *BIRFunction) Duplicate() *BIRFunction {
	dup := NewBIRFunctionSimple(f.pos, f.Name, f.OriginalName, f.Flags, f.Type, f.WorkerName, 0, f.Origin)
	dup.LocalVars = f.LocalVars
	dup.Parameters = f.Parameters
	dup.RequiredParams = f.RequiredParams
	dup.BasicBlocks = f.BasicBlocks
	dup.ErrorTable = f.ErrorTable
	dup.WorkerChannels = f.WorkerChannels
	dup.AnnotAttachments = f.AnnotAttachments
	dup.AnnotAttachmentsOnExternal = f.AnnotAttachmentsOnExternal
	dup.ReturnTypeAnnots = f.ReturnTypeAnnots
	return dup
}

func (f *BIRFunction) GetName() util.Name {
	return f.Name
}

type BIRBasicBlock struct {
	birNodeBase
	Number       int
	ID           util.Name
	Instructions []BIRNonTerminator
	Terminator   BIRTerminator
}

const BIRBasicBlockPrefix = "bb"

func NewBIRBasicBlock(id util.Name, number int) *BIRBasicBlock {
	return &BIRBasicBlock{
		birNodeBase:  birNodeBase{pos: nil},
		Number:       number,
		ID:           id,
		Instructions: make([]BIRNonTerminator, 0),
	}
}

func NewBIRBasicBlockWithNumber(number int) *BIRBasicBlock {
	return NewBIRBasicBlock(util.NewName(BIRBasicBlockPrefix+fmt.Sprintf("%d", number)), number)
}

func NewBIRBasicBlockWithPrefix(idPrefix string, number int) *BIRBasicBlock {
	return NewBIRBasicBlock(util.NewName(idPrefix+fmt.Sprintf("%d", number)), number)
}

func (b *BIRBasicBlock) String() string {
	return b.ID.Value
}

type BIRTypeDefinition struct {
	birDocumentableNodeBase
	Name             util.Name
	OriginalName     util.Name
	InternalName     util.Name
	AttachedFuncs    []*BIRFunction
	Flags            int64
	Type             types.BType
	IsBuiltin        bool
	ReferencedTypes  []types.BType
	ReferenceType    types.BType
	Origin           symbols.SymbolOrigin
	AnnotAttachments []*BIRAnnotationAttachment
	Index            int
}

func NewBIRTypeDefinition(pos diagnostics.Location, internalName util.Name, flags int64, isBuiltin bool, typ types.BType, attachedFuncs []*BIRFunction, origin symbols.SymbolOrigin, name, originalName util.Name) *BIRTypeDefinition {
	return &BIRTypeDefinition{
		birDocumentableNodeBase: birDocumentableNodeBase{birNodeBase: birNodeBase{pos: pos}},
		InternalName:            internalName,
		Flags:                   flags,
		IsBuiltin:               isBuiltin,
		Type:                    typ,
		AttachedFuncs:           attachedFuncs,
		ReferencedTypes:         make([]types.BType, 0),
		Origin:                  origin,
		Name:                    name,
		OriginalName:            originalName,
		AnnotAttachments:        make([]*BIRAnnotationAttachment, 0),
	}
}

func NewBIRTypeDefinitionSimple(pos diagnostics.Location, name, originalName util.Name, flags int64, isBuiltin bool, typ types.BType, attachedFuncs []*BIRFunction, origin symbols.SymbolOrigin) *BIRTypeDefinition {
	return NewBIRTypeDefinition(pos, name, flags, isBuiltin, typ, attachedFuncs, origin, name, originalName)
}

func (t *BIRTypeDefinition) GetName() util.Name {
	return t.Name
}

type BIRErrorEntry struct {
	birNodeBase
	TrapBB   *BIRBasicBlock
	EndBB    *BIRBasicBlock
	ErrorOp  *BIROperand
	TargetBB *BIRBasicBlock
}

func NewBIRErrorEntry(trapBB, endBB *BIRBasicBlock, errorOp *BIROperand, targetBB *BIRBasicBlock) *BIRErrorEntry {
	return &BIRErrorEntry{
		birNodeBase: birNodeBase{pos: nil},
		TrapBB:      trapBB,
		EndBB:       endBB,
		ErrorOp:     errorOp,
		TargetBB:    targetBB,
	}
}

type BIRAnnotation struct {
	birDocumentableNodeBase
	Name             util.Name
	OriginalName     util.Name
	Flags            int64
	Origin           symbols.SymbolOrigin
	AttachPoints     map[elements.AttachPoint]struct{}
	AnnotationType   types.BType
	PackageID        elements.PackageID
	AnnotAttachments []*BIRAnnotationAttachment
}

func NewBIRAnnotation(pos diagnostics.Location, name, originalName util.Name, flags int64, points map[elements.AttachPoint]struct{}, annotationType types.BType, origin symbols.SymbolOrigin) *BIRAnnotation {
	return &BIRAnnotation{
		birDocumentableNodeBase: birDocumentableNodeBase{birNodeBase: birNodeBase{pos: pos}},
		Name:                    name,
		OriginalName:            originalName,
		Flags:                   flags,
		AttachPoints:            points,
		AnnotationType:          annotationType,
		Origin:                  origin,
		AnnotAttachments:        make([]*BIRAnnotationAttachment, 0),
	}
}

type ConstValue struct {
	Type  types.BType
	Value interface{}
}

func NewConstValue(value interface{}, typ types.BType) *ConstValue {
	return &ConstValue{
		Value: value,
		Type:  typ,
	}
}

type BIRConstant struct {
	birDocumentableNodeBase
	Name             util.Name
	Flags            int64
	Type             types.BType
	ConstValue       *ConstValue
	Origin           symbols.SymbolOrigin
	AnnotAttachments []*BIRAnnotationAttachment
}

func NewBIRConstant(pos diagnostics.Location, name util.Name, flags int64, typ types.BType, constValue *ConstValue, origin symbols.SymbolOrigin) *BIRConstant {
	return &BIRConstant{
		birDocumentableNodeBase: birDocumentableNodeBase{birNodeBase: birNodeBase{pos: pos}},
		Name:                    name,
		Flags:                   flags,
		Type:                    typ,
		ConstValue:              constValue,
		Origin:                  origin,
		AnnotAttachments:        make([]*BIRAnnotationAttachment, 0),
	}
}

type BIRAnnotationAttachment struct {
	birNodeBase
	AnnotPkgID  elements.PackageID
	AnnotTagRef util.Name
}

func NewBIRAnnotationAttachment(pos diagnostics.Location, annotPkgID elements.PackageID, annotTagRef util.Name) *BIRAnnotationAttachment {
	return &BIRAnnotationAttachment{
		birNodeBase: birNodeBase{pos: pos},
		AnnotPkgID:  annotPkgID,
		AnnotTagRef: annotTagRef,
	}
}

type BIRConstAnnotationAttachment struct {
	*BIRAnnotationAttachment
	AnnotValue *ConstValue
}

func NewBIRConstAnnotationAttachment(pos diagnostics.Location, annotPkgID elements.PackageID, annotTagRef util.Name, annotValue *ConstValue) *BIRConstAnnotationAttachment {
	return &BIRConstAnnotationAttachment{
		BIRAnnotationAttachment: NewBIRAnnotationAttachment(pos, annotPkgID, annotTagRef),
		AnnotValue:              annotValue,
	}
}

type BIRLockDetailsHolder struct {
	locks []BIRLock
}

func NewBIRLockDetailsHolder() *BIRLockDetailsHolder {
	return &BIRLockDetailsHolder{
		locks: make([]BIRLock, 0),
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

func (h *BIRLockDetailsHolder) GetLock(index int) BIRLock {
	return h.locks[index]
}

func (h *BIRLockDetailsHolder) AddLock(lock BIRLock) {
	h.locks = append(h.locks, lock)
}

func (h *BIRLockDetailsHolder) Size() int {
	return len(h.locks)
}

type BIRMappingConstructorEntry interface {
	IsKeyValuePair() bool
}

type BIRMappingConstructorKeyValueEntry struct {
	KeyOp   *BIROperand
	ValueOp *BIROperand
}

func NewBIRMappingConstructorKeyValueEntry(keyOp, valueOp *BIROperand) *BIRMappingConstructorKeyValueEntry {
	return &BIRMappingConstructorKeyValueEntry{
		KeyOp:   keyOp,
		ValueOp: valueOp,
	}
}

func (e *BIRMappingConstructorKeyValueEntry) IsKeyValuePair() bool {
	return true
}

type BIRMappingConstructorSpreadFieldEntry struct {
	ExprOp *BIROperand
}

func NewBIRMappingConstructorSpreadFieldEntry(exprOp *BIROperand) *BIRMappingConstructorSpreadFieldEntry {
	return &BIRMappingConstructorSpreadFieldEntry{
		ExprOp: exprOp,
	}
}

func (e *BIRMappingConstructorSpreadFieldEntry) IsKeyValuePair() bool {
	return false
}

type BIRListConstructorEntry interface {
	GetExprOp() *BIROperand
}

type birListConstructorEntryBase struct {
	exprOp *BIROperand
}

func (b *birListConstructorEntryBase) GetExprOp() *BIROperand {
	return b.exprOp
}

type BIRListConstructorSpreadMemberEntry struct {
	birListConstructorEntryBase
}

func NewBIRListConstructorSpreadMemberEntry(exprOp *BIROperand) *BIRListConstructorSpreadMemberEntry {
	return &BIRListConstructorSpreadMemberEntry{
		birListConstructorEntryBase: birListConstructorEntryBase{exprOp: exprOp},
	}
}

type BIRListConstructorExprEntry struct {
	birListConstructorEntryBase
}

func NewBIRListConstructorExprEntry(exprOp *BIROperand) *BIRListConstructorExprEntry {
	return &BIRListConstructorExprEntry{
		birListConstructorEntryBase: birListConstructorEntryBase{exprOp: exprOp},
	}
}

type BIRServiceDeclaration struct {
	birDocumentableNodeBase
	AttachPoint         []string
	AttachPointLiteral  string
	ListenerTypes       []types.BType
	GeneratedName       util.Name
	AssociatedClassName util.Name
	Type                types.BType
	Origin              symbols.SymbolOrigin
	Flags               int64
}

func NewBIRServiceDeclaration(attachPoint []string, attachPointLiteral string, listenerTypes []types.BType, generatedName, associatedClassName util.Name, typ types.BType, origin symbols.SymbolOrigin, flags int64, location diagnostics.Location) *BIRServiceDeclaration {
	return &BIRServiceDeclaration{
		birDocumentableNodeBase: birDocumentableNodeBase{birNodeBase: birNodeBase{pos: location}},
		AttachPoint:             attachPoint,
		AttachPointLiteral:      attachPointLiteral,
		ListenerTypes:           listenerTypes,
		GeneratedName:           generatedName,
		AssociatedClassName:     associatedClassName,
		Type:                    typ,
		Origin:                  origin,
		Flags:                   flags,
	}
}

type BIROperand = BIROperandImpl

// BIRLock represents lock information (placeholder)
type BIRLock interface {
}

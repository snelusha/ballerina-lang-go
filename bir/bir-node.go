package bir

import (
	"fmt"
	"strings"
)

type BIRNode interface {
	Accept(visitor BIRVisitor)
	GetPos() Location
}

type BIRPackage struct {
	Pos                            Location
	PackageID                      *PackageID
	ImportModules                  []*BIRImportModule
	TypeDefs                       []*BIRTypeDefinition
	GlobalVars                     []*BIRGlobalVariableDcl
	ImportedGlobalVarsDummyVarDcls []*BIRGlobalVariableDcl
	Functions                      []*BIRFunction
	Annotations                    []*BIRAnnotation
	Constants                      []*BIRConstant
	ServiceDecls                   []*BIRServiceDeclaration
	IsListenerAvailable            bool
	RecordDefaultValueMap          map[string]map[string]string
}

func NewBIRPackage(pos Location, org, pkgName, name, version, sourceFileName *Name, sourceRoot string, skipTest bool) *BIRPackage {
	return NewBIRPackageWithTestPkg(pos, org, pkgName, name, version, sourceFileName, sourceRoot, skipTest, false)
}

func NewBIRPackageWithTestPkg(pos Location, org, pkgName, name, version, sourceFileName *Name, sourceRoot string, skipTest, isTestPkg bool) *BIRPackage {
	return &BIRPackage{
		Pos:                            pos,
		PackageID:                      NewPackageIDFull(org, pkgName, name, version, sourceFileName, sourceRoot, isTestPkg, skipTest),
		ImportModules:                  make([]*BIRImportModule, 0),
		TypeDefs:                       make([]*BIRTypeDefinition, 0),
		GlobalVars:                     make([]*BIRGlobalVariableDcl, 0),
		ImportedGlobalVarsDummyVarDcls: make([]*BIRGlobalVariableDcl, 0),
		Functions:                      make([]*BIRFunction, 0),
		Annotations:                    make([]*BIRAnnotation, 0),
		Constants:                      make([]*BIRConstant, 0),
		ServiceDecls:                   make([]*BIRServiceDeclaration, 0),
		RecordDefaultValueMap:          make(map[string]map[string]string),
	}
}

func (p *BIRPackage) Accept(visitor BIRVisitor) {
	visitor.VisitPackage(p)
}

func (p *BIRPackage) GetPos() Location {
	return p.Pos
}

type BIRImportModule struct {
	Pos       Location
	PackageID *PackageID
}

func NewBIRImportModule(pos Location, org, name, version *Name) *BIRImportModule {
	return &BIRImportModule{
		Pos:       pos,
		PackageID: NewPackageID(org, name, name, version),
	}
}

func (m *BIRImportModule) Accept(visitor BIRVisitor) {
	visitor.VisitImportModule(m)
}

func (m *BIRImportModule) GetPos() Location {
	return m.Pos
}

func (m *BIRImportModule) Equals(other *BIRImportModule) bool {
	if m == other {
		return true
	}
	if other == nil {
		return false
	}
	return m.PackageID.Equals(other.PackageID)
}

type BIRScope interface {
	GetStartBB() *BIRBasicBlock
	GetEndBB() *BIRBasicBlock
}

type BIRVariableDcl struct {
	Pos                   Location
	Type                  BType
	Name                  *Name
	OriginalName          *Name
	MetaVarName           string
	JvmVarName            string
	Kind                  VarKind
	Scope                 VarScope
	IgnoreVariable        bool
	EndBB                 *BIRBasicBlock
	StartBB               *BIRBasicBlock
	InsOffset             int
	OnlyUsedInSingleBB    bool
	Initialized           bool
	InsScope              BIRScope
	MarkdownDocAttachment *MarkdownDocAttachment
}

func NewBIRVariableDcl(pos Location, typ BType, name, originalName *Name, scope VarScope, kind VarKind, metaVarName string) *BIRVariableDcl {
	jvmVarName := strings.ReplaceAll(name.Value, "%", "_")
	return &BIRVariableDcl{
		Pos:          pos,
		Type:         typ,
		Name:         name,
		OriginalName: originalName,
		Scope:        scope,
		Kind:         kind,
		MetaVarName:  metaVarName,
		JvmVarName:   jvmVarName,
		Initialized:  false,
	}
}

func NewBIRVariableDclSimple(pos Location, typ BType, name *Name, scope VarScope, kind VarKind, metaVarName string) *BIRVariableDcl {
	return NewBIRVariableDcl(pos, typ, name, name, scope, kind, metaVarName)
}

func NewBIRVariableDclBasic(typ BType, name *Name, scope VarScope, kind VarKind) *BIRVariableDcl {
	return NewBIRVariableDclSimple(nil, typ, name, scope, kind, "")
}

func (v *BIRVariableDcl) Accept(visitor BIRVisitor) {
	visitor.VisitVariableDcl(v)
}

func (v *BIRVariableDcl) GetPos() Location {
	return v.Pos
}

func (v *BIRVariableDcl) SetMarkdownDocAttachment(markdownDocAttachment *MarkdownDocAttachment) {
	v.MarkdownDocAttachment = markdownDocAttachment
}

func (v *BIRVariableDcl) Equals(other *BIRVariableDcl) bool {
	if v == other {
		return true
	}
	if other == nil {
		return false
	}
	return v.Name.Equals(other.Name)
}

func (v *BIRVariableDcl) String() string {
	return v.Name.String()
}

type BIRParameter struct {
	Pos              Location
	Name             *Name
	Flags            int64
	AnnotAttachments []*BIRAnnotationAttachment
}

func NewBIRParameter(pos Location, name *Name, flags int64) *BIRParameter {
	return &BIRParameter{
		Pos:              pos,
		Name:             name,
		Flags:            flags,
		AnnotAttachments: make([]*BIRAnnotationAttachment, 0),
	}
}

func (p *BIRParameter) Accept(visitor BIRVisitor) {
	visitor.VisitParameter(p)
}

func (p *BIRParameter) GetPos() Location {
	return p.Pos
}

type BIRGlobalVariableDcl struct {
	*BIRVariableDcl
	Flags            int64
	PkgId            *PackageID
	Origin           SymbolOrigin
	AnnotAttachments []*BIRAnnotationAttachment
}

func NewBIRGlobalVariableDcl(pos Location, flags int64, typ BType, pkgId *PackageID, name, originalName *Name, scope VarScope, kind VarKind, metaVarName string, origin SymbolOrigin) *BIRGlobalVariableDcl {
	return &BIRGlobalVariableDcl{
		BIRVariableDcl:   NewBIRVariableDcl(pos, typ, name, originalName, scope, kind, metaVarName),
		Flags:            flags,
		PkgId:            pkgId,
		Origin:           origin,
		AnnotAttachments: make([]*BIRAnnotationAttachment, 0),
	}
}

func (g *BIRGlobalVariableDcl) Accept(visitor BIRVisitor) {
	visitor.VisitVariableDcl(g.BIRVariableDcl)
}

type BIRFunctionParameter struct {
	*BIRVariableDcl
	HasDefaultExpr  bool
	IsPathParameter bool
}

func NewBIRFunctionParameter(pos Location, typ BType, name *Name, scope VarScope, kind VarKind, metaVarName string, hasDefaultExpr bool) *BIRFunctionParameter {
	return &BIRFunctionParameter{
		BIRVariableDcl: NewBIRVariableDclSimple(pos, typ, name, scope, kind, metaVarName),
		HasDefaultExpr: hasDefaultExpr,
	}
}

func NewBIRFunctionParameterWithPath(pos Location, typ BType, name *Name, scope VarScope, kind VarKind, metaVarName string, hasDefaultExpr, isPathParameter bool) *BIRFunctionParameter {
	return &BIRFunctionParameter{
		BIRVariableDcl:  NewBIRVariableDclSimple(pos, typ, name, scope, kind, metaVarName),
		HasDefaultExpr:  hasDefaultExpr,
		IsPathParameter: isPathParameter,
	}
}

func (f *BIRFunctionParameter) Accept(visitor BIRVisitor) {
	visitor.VisitFunctionParameter(f)
}

type BIRFunction struct {
	Pos                        Location
	Name                       *Name
	OriginalName               *Name
	Flags                      int64
	Origin                     SymbolOrigin
	Type                       BInvokableType
	RequiredParams             []*BIRParameter
	Receiver                   *BIRVariableDcl
	RestParam                  *BIRParameter
	ArgsCount                  int
	LocalVars                  []*BIRVariableDcl
	ReturnVariable             *BIRVariableDcl
	Parameters                 []*BIRFunctionParameter
	BasicBlocks                []*BIRBasicBlock
	ErrorTable                 []*BIRErrorEntry
	WorkerName                 *Name
	WorkerChannels             []*ChannelDetails
	AnnotAttachments           []*BIRAnnotationAttachment
	AnnotAttachmentsOnExternal []*BIRAnnotationAttachment
	ReturnTypeAnnots           []*BIRAnnotationAttachment
	DependentGlobalVars        []*BIRGlobalVariableDcl
	PathParams                 []*BIRVariableDcl
	RestPathParam              *BIRVariableDcl
	ResourcePath               []*Name
	ResourcePathSegmentPosList []Location
	Accessor                   *Name
	PathSegmentTypeList        []BType
	HasWorkers                 bool
	MarkdownDocAttachment      *MarkdownDocAttachment
}

func NewBIRFunction(pos Location, name, originalName *Name, flags int64, origin SymbolOrigin, typ BInvokableType, requiredParams []*BIRParameter, receiver *BIRVariableDcl, restParam *BIRParameter, argsCount int, localVars []*BIRVariableDcl, returnVariable *BIRVariableDcl, parameters []*BIRFunctionParameter, basicBlocks []*BIRBasicBlock, errorTable []*BIRErrorEntry, workerName *Name, workerChannels []*ChannelDetails, annotAttachments, returnTypeAnnots []*BIRAnnotationAttachment, dependentGlobalVars []*BIRGlobalVariableDcl) *BIRFunction {
	return &BIRFunction{
		Pos:                 pos,
		Name:                name,
		OriginalName:        originalName,
		Flags:               flags,
		Origin:              origin,
		Type:                typ,
		RequiredParams:      requiredParams,
		Receiver:            receiver,
		RestParam:           restParam,
		ArgsCount:           argsCount,
		LocalVars:           localVars,
		ReturnVariable:      returnVariable,
		Parameters:          parameters,
		BasicBlocks:         basicBlocks,
		ErrorTable:          errorTable,
		WorkerName:          workerName,
		WorkerChannels:      workerChannels,
		AnnotAttachments:    annotAttachments,
		ReturnTypeAnnots:    returnTypeAnnots,
		DependentGlobalVars: dependentGlobalVars,
	}
}

func NewBIRFunctionSimple(pos Location, name, originalName *Name, flags int64, typ BInvokableType, workerName *Name, sendInsCount int, origin SymbolOrigin) *BIRFunction {
	return &BIRFunction{
		Pos:                 pos,
		Name:                name,
		OriginalName:        originalName,
		Flags:               flags,
		Type:                typ,
		LocalVars:           make([]*BIRVariableDcl, 0),
		Parameters:          make([]*BIRFunctionParameter, 0),
		RequiredParams:      make([]*BIRParameter, 0),
		BasicBlocks:         make([]*BIRBasicBlock, 0),
		ErrorTable:          make([]*BIRErrorEntry, 0),
		WorkerName:          workerName,
		WorkerChannels:      make([]*ChannelDetails, sendInsCount),
		AnnotAttachments:    make([]*BIRAnnotationAttachment, 0),
		ReturnTypeAnnots:    make([]*BIRAnnotationAttachment, 0),
		Origin:              origin,
		DependentGlobalVars: make([]*BIRGlobalVariableDcl, 0),
	}
}

func NewBIRFunctionBasic(pos Location, name *Name, flags int64, typ BInvokableType, workerName *Name, sendInsCount int, origin SymbolOrigin) *BIRFunction {
	return NewBIRFunctionSimple(pos, name, name, flags, typ, workerName, sendInsCount, origin)
}

func (f *BIRFunction) Accept(visitor BIRVisitor) {
	visitor.VisitFunction(f)
}

func (f *BIRFunction) GetPos() Location {
	return f.Pos
}

func (f *BIRFunction) SetMarkdownDocAttachment(markdownDocAttachment *MarkdownDocAttachment) {
	f.MarkdownDocAttachment = markdownDocAttachment
}

func (f *BIRFunction) GetName() *Name {
	return f.Name
}

func (f *BIRFunction) Duplicate() *BIRFunction {
	newFunc := NewBIRFunctionSimple(f.Pos, f.Name, f.OriginalName, f.Flags, f.Type, f.WorkerName, 0, f.Origin)
	newFunc.LocalVars = f.LocalVars
	newFunc.Parameters = f.Parameters
	newFunc.RequiredParams = f.RequiredParams
	newFunc.BasicBlocks = f.BasicBlocks
	newFunc.ErrorTable = f.ErrorTable
	newFunc.WorkerChannels = f.WorkerChannels
	newFunc.AnnotAttachments = f.AnnotAttachments
	newFunc.AnnotAttachmentsOnExternal = f.AnnotAttachmentsOnExternal
	newFunc.ReturnTypeAnnots = f.ReturnTypeAnnots
	return newFunc
}

type BIRBasicBlock struct {
	Pos          Location
	Number       int
	ID           *Name
	Instructions []BIRNonTerminator
	Terminator   BIRTerminator
}

const BIRBasicBlockPrefix = "bb"

func NewBIRBasicBlock(id *Name, number int) *BIRBasicBlock {
	return &BIRBasicBlock{
		Pos:          nil,
		Number:       number,
		ID:           id,
		Instructions: make([]BIRNonTerminator, 0),
		Terminator:   nil,
	}
}

func NewBIRBasicBlockSimple(number int) *BIRBasicBlock {
	return NewBIRBasicBlock(NewName(fmt.Sprintf("%s%d", BIRBasicBlockPrefix, number)), number)
}

func NewBIRBasicBlockWithPrefix(idPrefix string, number int) *BIRBasicBlock {
	return NewBIRBasicBlock(NewName(fmt.Sprintf("%s%d", idPrefix, number)), number)
}

func (b *BIRBasicBlock) Accept(visitor BIRVisitor) {
	visitor.VisitBasicBlock(b)
}

func (b *BIRBasicBlock) GetPos() Location {
	return b.Pos
}

func (b *BIRBasicBlock) String() string {
	return b.ID.Value
}

type BIRTypeDefinition struct {
	Pos                   Location
	Name                  *Name
	OriginalName          *Name
	InternalName          *Name
	AttachedFuncs         []*BIRFunction
	Flags                 int64
	Type                  BType
	IsBuiltin             bool
	ReferencedTypes       []BType
	ReferenceType         BType
	Origin                SymbolOrigin
	AnnotAttachments      []*BIRAnnotationAttachment
	Index                 int
	MarkdownDocAttachment *MarkdownDocAttachment
}

func NewBIRTypeDefinition(pos Location, internalName *Name, flags int64, isBuiltin bool, typ BType, attachedFuncs []*BIRFunction, origin SymbolOrigin, name, originalName *Name) *BIRTypeDefinition {
	return &BIRTypeDefinition{
		Pos:              pos,
		InternalName:     internalName,
		Flags:            flags,
		IsBuiltin:        isBuiltin,
		Type:             typ,
		AttachedFuncs:    attachedFuncs,
		ReferencedTypes:  make([]BType, 0),
		Origin:           origin,
		Name:             name,
		OriginalName:     originalName,
		AnnotAttachments: make([]*BIRAnnotationAttachment, 0),
	}
}

func NewBIRTypeDefinitionSimple(pos Location, name, originalName *Name, flags int64, isBuiltin bool, typ BType, attachedFuncs []*BIRFunction, origin SymbolOrigin) *BIRTypeDefinition {
	return NewBIRTypeDefinition(pos, name, flags, isBuiltin, typ, attachedFuncs, origin, name, originalName)
}

func (t *BIRTypeDefinition) Accept(visitor BIRVisitor) {
	visitor.VisitTypeDefinition(t)
}

func (t *BIRTypeDefinition) GetPos() Location {
	return t.Pos
}

func (t *BIRTypeDefinition) SetMarkdownDocAttachment(markdownDocAttachment *MarkdownDocAttachment) {
	t.MarkdownDocAttachment = markdownDocAttachment
}

func (t *BIRTypeDefinition) String() string {
	return fmt.Sprintf("%v %s", t.Type, t.InternalName)
}

func (t *BIRTypeDefinition) GetName() *Name {
	return t.Name
}

type BIRErrorEntry struct {
	Pos      Location
	TrapBB   *BIRBasicBlock
	EndBB    *BIRBasicBlock
	ErrorOp  *BIROperand
	TargetBB *BIRBasicBlock
}

func NewBIRErrorEntry(trapBB, endBB *BIRBasicBlock, errorOp *BIROperand, targetBB *BIRBasicBlock) *BIRErrorEntry {
	return &BIRErrorEntry{
		Pos:      nil,
		TrapBB:   trapBB,
		EndBB:    endBB,
		ErrorOp:  errorOp,
		TargetBB: targetBB,
	}
}

func (e *BIRErrorEntry) Accept(visitor BIRVisitor) {
	visitor.VisitErrorEntry(e)
}

func (e *BIRErrorEntry) GetPos() Location {
	return e.Pos
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

type BIRAnnotation struct {
	Pos                   Location
	Name                  *Name
	OriginalName          *Name
	Flags                 int64
	Origin                SymbolOrigin
	AttachPoints          []*AttachPoint
	AnnotationType        BType
	PackageID             *PackageID
	AnnotAttachments      []*BIRAnnotationAttachment
	MarkdownDocAttachment *MarkdownDocAttachment
}

func NewBIRAnnotation(pos Location, name, originalName *Name, flags int64, points []*AttachPoint, annotationType BType, origin SymbolOrigin) *BIRAnnotation {
	return &BIRAnnotation{
		Pos:              pos,
		Name:             name,
		OriginalName:     originalName,
		Flags:            flags,
		AttachPoints:     points,
		AnnotationType:   annotationType,
		Origin:           origin,
		AnnotAttachments: make([]*BIRAnnotationAttachment, 0),
	}
}

func (a *BIRAnnotation) Accept(visitor BIRVisitor) {
	visitor.VisitAnnotation(a)
}

func (a *BIRAnnotation) GetPos() Location {
	return a.Pos
}

func (a *BIRAnnotation) SetMarkdownDocAttachment(markdownDocAttachment *MarkdownDocAttachment) {
	a.MarkdownDocAttachment = markdownDocAttachment
}

type BIRConstant struct {
	Pos                   Location
	Name                  *Name
	Flags                 int64
	Type                  BType
	ConstValue            *ConstValue
	Origin                SymbolOrigin
	AnnotAttachments      []*BIRAnnotationAttachment
	MarkdownDocAttachment *MarkdownDocAttachment
}

func NewBIRConstant(pos Location, name *Name, flags int64, typ BType, constValue *ConstValue, origin SymbolOrigin) *BIRConstant {
	return &BIRConstant{
		Pos:              pos,
		Name:             name,
		Flags:            flags,
		Type:             typ,
		ConstValue:       constValue,
		Origin:           origin,
		AnnotAttachments: make([]*BIRAnnotationAttachment, 0),
	}
}

func (c *BIRConstant) Accept(visitor BIRVisitor) {
	visitor.VisitConstant(c)
}

func (c *BIRConstant) GetPos() Location {
	return c.Pos
}

func (c *BIRConstant) SetMarkdownDocAttachment(markdownDocAttachment *MarkdownDocAttachment) {
	c.MarkdownDocAttachment = markdownDocAttachment
}

type BIRAnnotationAttachment struct {
	Pos         Location
	AnnotPkgId  *PackageID
	AnnotTagRef *Name
}

func NewBIRAnnotationAttachment(pos Location, annotPkgId *PackageID, annotTagRef *Name) *BIRAnnotationAttachment {
	return &BIRAnnotationAttachment{
		Pos:         pos,
		AnnotPkgId:  annotPkgId,
		AnnotTagRef: annotTagRef,
	}
}

func (a *BIRAnnotationAttachment) Accept(visitor BIRVisitor) {
	visitor.VisitAnnotationAttachment(a)
}

func (a *BIRAnnotationAttachment) GetPos() Location {
	return a.Pos
}

type BIRConstAnnotationAttachment struct {
	*BIRAnnotationAttachment
	AnnotValue *ConstValue
}

func NewBIRConstAnnotationAttachment(pos Location, annotPkgId *PackageID, annotTagRef *Name, annotValue *ConstValue) *BIRConstAnnotationAttachment {
	return &BIRConstAnnotationAttachment{
		BIRAnnotationAttachment: NewBIRAnnotationAttachment(pos, annotPkgId, annotTagRef),
		AnnotValue:              annotValue,
	}
}

func (c *BIRConstAnnotationAttachment) Accept(visitor BIRVisitor) {
	visitor.VisitConstAnnotationAttachment(c)
}

type ConstValue struct {
	Type  BType
	Value interface{}
}

func NewConstValue(value interface{}, typ BType) *ConstValue {
	return &ConstValue{
		Value: value,
		Type:  typ,
	}
}

type BIRLockDetailsHolder struct {
	locks []Lock
}

func NewBIRLockDetailsHolder() *BIRLockDetailsHolder {
	return &BIRLockDetailsHolder{
		locks: make([]Lock, 0),
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

func (h *BIRLockDetailsHolder) GetLock(index int) Lock {
	return h.locks[index]
}

func (h *BIRLockDetailsHolder) AddLock(lock Lock) {
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

type BIRListConstructorSpreadMemberEntry struct {
	ExprOp *BIROperand
}

func NewBIRListConstructorSpreadMemberEntry(exprOp *BIROperand) *BIRListConstructorSpreadMemberEntry {
	return &BIRListConstructorSpreadMemberEntry{
		ExprOp: exprOp,
	}
}

func (e *BIRListConstructorSpreadMemberEntry) GetExprOp() *BIROperand {
	return e.ExprOp
}

type BIRListConstructorExprEntry struct {
	ExprOp *BIROperand
}

func NewBIRListConstructorExprEntry(exprOp *BIROperand) *BIRListConstructorExprEntry {
	return &BIRListConstructorExprEntry{
		ExprOp: exprOp,
	}
}

func (e *BIRListConstructorExprEntry) GetExprOp() *BIROperand {
	return e.ExprOp
}

type BIRServiceDeclaration struct {
	Pos                   Location
	AttachPoint           []string
	AttachPointLiteral    string
	ListenerTypes         []BType
	GeneratedName         *Name
	AssociatedClassName   *Name
	Type                  BType
	Origin                SymbolOrigin
	Flags                 int64
	MarkdownDocAttachment *MarkdownDocAttachment
}

func NewBIRServiceDeclaration(attachPoint []string, attachPointLiteral string, listenerTypes []BType, generatedName, associatedClassName *Name, typ BType, origin SymbolOrigin, flags int64, location Location) *BIRServiceDeclaration {
	return &BIRServiceDeclaration{
		Pos:                 location,
		AttachPoint:         attachPoint,
		AttachPointLiteral:  attachPointLiteral,
		ListenerTypes:       listenerTypes,
		GeneratedName:       generatedName,
		AssociatedClassName: associatedClassName,
		Type:                typ,
		Origin:              origin,
		Flags:               flags,
	}
}

func (s *BIRServiceDeclaration) Accept(visitor BIRVisitor) {
	visitor.VisitServiceDeclaration(s)
}

func (s *BIRServiceDeclaration) GetPos() Location {
	return s.Pos
}

func (s *BIRServiceDeclaration) SetMarkdownDocAttachment(markdownDocAttachment *MarkdownDocAttachment) {
	s.MarkdownDocAttachment = markdownDocAttachment
}

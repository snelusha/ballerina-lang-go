package compiler

import (
	"ballerina-lang-go/compiler/syntax/tree"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/tools/text"
	"strings"
)

type STNodeDiagnostic interface {
	DiagnosticCode() diagnostics.DiagnosticCode
	Args() []interface{}
}

type STNodeVisitor interface {
	Visit(node STNode)
}

type STNodeTransformer interface {
	Transform(node STNode) interface{}
}

type STNodeList interface {
	STNode
	Size() int
	Get(index int) STNode
	Set(index int, node STNode) STNodeList
	Add(node STNode) STNodeList
	AddAll(nodes []STNode) STNodeList
	Remove(index int) STNodeList
	IsEmpty() bool
}

type STNode interface {
	ChildInBucket(bucket int) STNode
	WidthWithMinutiae() int
	Width() int
	WidthWithLeadingMinutiae() int
	WidthWithTrailingMinutiae() int
	LeadingMinutiae() STNode
	TrailingMinutiae() STNode
	HasDiagnostics() bool
	Diagnostics() []STNodeDiagnostic
	BucketCount() int
	IsMissing() bool
	Tokens() []STToken
	FirstToken() STToken
	LastToken() STToken
	ModifyWith(diagnostics []STNodeDiagnostic) STNode
	Replace(target, replacement STNode) STNode
	CreateUnlinkedFacade() Node
	CreateFacade(position int, parent NonTerminalNode) Node
	Accept(visitor STNodeVisitor)
	Apply(transformer STNodeTransformer) interface{}
	String() string
	WriteTo(builder *strings.Builder)
	ToSourceCode() string

	tokensInternal(tokens *[]STToken)
	firstTokenInternal() STNode
	lastTokenInternal() STNode
	checkForReferenceEquality(children ...STNode) bool
	updateWithChildren(children []STNode)
	updateWithChildrenAndIndexes(children []STNode, firstChildIndex, lastChildIndex int)
	getFirstChildIndex(children ...STNode) int
	getLastChildIndex(children ...STNode) int
	updateDiagnostics(children []STNode)
}

type STToken interface {
	STNode
	Text() string
	HasTrailingNewline() bool
	LookbackTokenCount() int
}

type NonTerminalNode interface {
	Node
	Children() ChildNodeList
	ChildEntries() []ChildNodeEntry
	FindTokenWithPosition(position int) Token
	FindTokenWithPositionAndInsideMinutiae(position int, insideMinutiae bool) Token
	FindNodeWithTextRange(textRange text.TextRange) NonTerminalNode
	FindNodeWithTextRangeIncludeStartOffset(textRange text.TextRange, includeStartOffset bool) NonTerminalNode
	Replace(target, replacement Node) NonTerminalNode
}

type ChildNodeList interface {
	Get(childIndex int) Node
	Size() int
	Iterator() <-chan *Node
}

type ChildNodeEntry interface {
	Name() string
	Node() Node
	NodeList() NodeList[Node]
	IsList() bool
}

type NodeList[T interface{}] interface {
	Node
	Get() T
	Add(node T) NodeList[T]
	AddAt(index int, node T) NodeList[T]
	AddAll(nodes []T) NodeList[T]
	Set(index int, node T) NodeList[T]
	Remove(node T) NodeList[T]
	RemoveAt(index int) NodeList[T]
	RemoveAll() NodeList[T]
	Size() int
	IsEmpty() bool
	Stream() <-chan T
}

type NodeLocation interface {
	diagnostics.Location
}

type InvalidTokenMinutiaeNode interface {
	NonTerminalNode
	InvalidToken() Token
	ParentMinutiae() Minutiae
}

type Minutiae interface {
	Text() string
	Kind() tree.SyntaxKind
	ParentToken() Token
	IsInvalidNodeMinutiae() bool
	InvalidTokenMinutiaeNode() InvalidTokenMinutiaeNode
	TextRange() text.TextRange
	LineRange() text.LineRange
}

type MinutiaeList interface {
	Get(index int) Minutiae
	Add(minutiae Minutiae) MinutiaeList
	AddAt(index int, minutiae Minutiae) MinutiaeList
	AddAll(minutiaeList []Minutiae) MinutiaeList
	Set(index int, minutiae Minutiae) MinutiaeList
	Remove(minutiae Minutiae) MinutiaeList
	RemoveAt(index int) MinutiaeList
	RemoveAll() MinutiaeList
	Size() int
	IsEmpty() bool
	String() string
}

type Token interface {
	Node
	Text() string
	ContainsLeadingMinutiae() bool
	ContainsTrailingMinutiae() bool
	Modify(leadingMinutiae, trailingMinutiae MinutiaeList) Token
}

type Node interface {
	Position() int
	Parent() NonTerminalNode
	Ancestors() []NonTerminalNode
	AncestorWithFilter(filter func(Node) bool) *NonTerminalNode
	TextRange() text.TextRange
	TextRangeWithMinutiae() text.TextRange
	Kind() tree.SyntaxKind
	Location() NodeLocation
	Diagnostics() []diagnostics.Diagnostic
	HasDiagnostics() bool
	IsMissing() bool
	SyntaxTree() SyntaxTree
	LineRange() text.LineRange
	LeadingMinutiae() MinutiaeList
	TrailingMinutiae() MinutiaeList
	LeadingInvalidTokens() []Token
	TrailingInvalidTokens() []Token
	Accept(visitor NodeVisitor)
	Apply(transformer NodeTransformer) interface{}
	String() string
	ToSourceCode() string
}

type SyntaxTree interface {
	TextDocument() text.TextDocument
	ContainsModulePart() bool
	RootNode() Node
	FilePath() string
	ModifyWith(rootNode Node) SyntaxTree
	ReplaceNode(target, replacement Node) SyntaxTree
	HasDiagnostics() bool
	String() string
	Diagnostics() []diagnostics.Diagnostic
	ToSourceCode() string
}

type NodeVisitor interface{}

type NodeTransformer interface{}

type Iterator[T interface{}] interface {
	Next() T
	HasNext() bool
}

type iteratorImpl[T interface{}] struct {
	data []T
	pos  int
}

func (it *iteratorImpl[T]) HasNext() bool {
	return it.pos < len(it.data)
}

func (it *iteratorImpl[T]) Next() T {
	val := it.data[it.pos]
	it.pos = it.pos + 1

	return val
}

func NewIterator[T interface{}](data []T) Iterator[T] {
	return &iteratorImpl[T]{data: data, pos: 0}
}

type MapIterator[K comparable, V interface{}] interface {
	Next() (K, V)
	HasNext() bool
}

type mapIteratorImpl[K comparable, V interface{}] struct {
	keys []K
	data map[K]V
	pos  int
}

func (it *mapIteratorImpl[K, V]) HasNext() bool {
	return it.pos < len(it.keys)
}

func (it *mapIteratorImpl[K, V]) Next() (K, V) {
	key := it.keys[it.pos]
	val := it.data[key]
	it.pos = it.pos + 1

	return key, val
}

func NewMapIterator[K comparable, V interface{}](m map[K]V) MapIterator[K, V] {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	return &mapIteratorImpl[K, V]{keys: keys, data: m, pos: 0}
}

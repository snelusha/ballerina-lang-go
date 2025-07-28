package compiler

import (
	"ballerina-lang-go/compiler/internal/diagnostics"
	"ballerina-lang-go/compiler/syntax/tree"
	diagnosticsTools "ballerina-lang-go/tools/diagnostics"
	"strings"
)

// STNodeDiagnostic represents the internal representation of diagnostic that is related to an internal syntax node.
type STNodeDiagnostic interface {
	diagnostics.IRDiagnostic
	DiagnosticCode() diagnosticsTools.DiagnosticCode
	Args() []interface{}
}

// STNode is the interface for the abstract base class for all tree nodes in the internal syntax tree.
type STNode interface {
	Kind() tree.SyntaxKind
	Diagnostics() []STNodeDiagnostic
	Width() int
	WidthWithMinutiae() int
	WidthWithLeadingMinutiae() int
	WidthWithTrailingMinutiae() int
	LeadingMinutiae() STNode
	TrailingMinutiae() STNode
	HasDiagnostics() bool
	BucketCount() int
	ChildInBucket(bucket int) STNode
	IsMissing() bool
	Tokens() []STToken
	FirstToken() STToken
	LastToken() STToken
	ModifyWith(diagnostics []STNodeDiagnostic) STNode
	Replace(target STNode, replacement STNode) STNode
	CreateUnlinkedFacade() tree.Node
	CreateFacade(position int, parent tree.NonTerminalNode) tree.Node
	Accept(visitor STNodeVisitor)
	Apply(transformer STNodeTransformer) interface{}
	ToString() string
	WriteTo(builder *strings.Builder)
	ToSourceCode() string
}

// STToken represents a terminal node in the internal syntax tree.
type STToken interface {
	STNode
	Text() string
	HasTrailingNewline() bool
	LookbackTokenCount() int
	ModifyWithMinutiae(leadingMinutiae STNode, trailingMinutiae STNode) STToken
}

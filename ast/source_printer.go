// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ast

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"ballerina-lang-go/model"
)

// SourcePrinter renders BLang AST nodes back to formatted Ballerina source.
type SourcePrinter struct {
	indentLevel int
	buffer      strings.Builder
}

func ToSource(node BLangNode) string {
	printer := SourcePrinter{}
	return printer.Print(node)
}

func (p *SourcePrinter) Print(node BLangNode) string {
	p.printNode(node)
	return strings.TrimRight(p.buffer.String(), "\n")
}

func (p *SourcePrinter) printNode(node BLangNode) {
	switch t := node.(type) {
	case *BLangPackage:
		p.printPackage(t)
	case *BLangCompilationUnit:
		p.printCompilationUnit(t)
	case *BLangImportPackage:
		p.printImportPackage(t)
	case *BLangXMLNS:
		p.printXMLNS(t)
	case *BLangConstant:
		p.printConstant(t)
	case *BLangSimpleVariable:
		p.printSimpleVariable(t, true)
	case *BLangTypeDefinition:
		p.printTypeDefinition(t)
	case *BLangClassDefinition:
		p.printClassDefinition(t)
	case *BLangFunction:
		p.printFunction(t)
	case *BLangBlockFunctionBody:
		p.printBlockFunctionBody(t)
	case *BLangBlockStmt:
		p.printBlockStmt(t)
	case *BLangSimpleVariableDef:
		p.printSimpleVariableDef(t)
	case *BLangAssignment:
		p.printAssignment(t)
	case *BLangCompoundAssignment:
		p.printCompoundAssignment(t)
	case *BLangExpressionStmt:
		p.write(p.expr(t.Expr))
		p.write(";")
	case *BLangReturn:
		p.write("return")
		if t.Expr != nil {
			p.write(" ")
			p.write(p.expr(t.Expr))
		}
		p.write(";")
	case *BLangPanic:
		p.write("panic ")
		p.write(p.expr(t.Expr))
		p.write(";")
	case *BLangBreak:
		p.write("break;")
	case *BLangContinue:
		p.write("continue;")
	case *BLangIf:
		p.printIf(t)
	case *BLangWhile:
		p.printWhile(t)
	case *BLangForeach:
		p.printForeach(t)
	case *BLangLock:
		p.write("lock ")
		p.printBlockStmt(&t.Body)
	case *BLangMatchStatement:
		p.printMatchStatement(t)
	default:
		panic("unsupported AST node for source printing: " + reflect.TypeOf(t).String())
	}
}

func (p *SourcePrinter) printPackage(node *BLangPackage) {
	for i := range node.TopLevelNodes {
		if i > 0 {
			p.newline()
			p.newline()
		}
		p.printNode(node.TopLevelNodes[i].(BLangNode))
	}
	if len(node.TopLevelNodes) > 0 {
		return
	}
	first := true
	printTop := func(node BLangNode) {
		if !first {
			p.newline()
			p.newline()
		}
		first = false
		p.printNode(node)
	}
	for i := range node.Imports {
		printTop(&node.Imports[i])
	}
	for i := range node.Constants {
		printTop(&node.Constants[i])
	}
	for i := range node.GlobalVars {
		printTop(&node.GlobalVars[i])
	}
	for i := range node.TypeDefinitions {
		printTop(&node.TypeDefinitions[i])
	}
	for i := range node.ClassDefinitions {
		printTop(&node.ClassDefinitions[i])
	}
	for i := range node.Functions {
		printTop(&node.Functions[i])
	}
}

func (p *SourcePrinter) printCompilationUnit(node *BLangCompilationUnit) {
	for i, topLevelNode := range node.TopLevelNodes {
		if i > 0 {
			p.newline()
			p.newline()
		}
		p.printNode(topLevelNode.(BLangNode))
	}
}

func (p *SourcePrinter) printImportPackage(node *BLangImportPackage) {
	p.write("import ")
	if node.OrgName != nil && node.OrgName.Value != "" {
		p.write(p.identifier(node.OrgName))
		p.write("/")
	}
	parts := make([]string, 0, len(node.PkgNameComps))
	for i := range node.PkgNameComps {
		parts = append(parts, p.identifier(&node.PkgNameComps[i]))
	}
	p.write(strings.Join(parts, "."))
	if node.Alias != nil && node.Alias.Value != "" {
		p.write(" as ")
		p.write(p.identifier(node.Alias))
	}
	p.write(";")
}

func (p *SourcePrinter) printXMLNS(node *BLangXMLNS) {
	p.write("xmlns ")
	p.write(p.expr(node.GetNamespaceURI()))
	if prefix := node.GetPrefix(); prefix != nil && prefix.Value != "" {
		p.write(" as ")
		p.write(p.identifier(prefix))
	}
	p.write(";")
}

func (p *SourcePrinter) printConstant(node *BLangConstant) {
	p.printDocumentationNode(node.MarkdownDocumentationAttachment)
	p.writeQualifiers(node.Flags(), model.FlagPublic)
	p.write("const ")
	p.printVariableParts(node.TypeNode(), node.Name, node.Expr)
	p.write(";")
}

func (p *SourcePrinter) printSimpleVariable(node *BLangSimpleVariable, semicolon bool) {
	p.printDocumentationNode(node.MarkdownDocumentationAttachment)
	p.writeQualifiers(node.Flags(), model.FlagPublic, model.FlagFinal, model.FlagConfigurable, model.FlagIsolated)
	p.printVariableParts(node.TypeNode(), node.Name, node.Expr)
	if semicolon {
		p.write(";")
	}
}

func (p *SourcePrinter) printVariableParts(typeNode BType, name *BLangIdentifier, expr BLangActionOrExpression) {
	if typeNode == nil {
		p.write("var")
	} else {
		p.write(p.typ(typeNode))
	}
	if name != nil && name.Value != "" {
		p.write(" ")
		p.write(p.identifier(name))
	}
	if expr != nil {
		p.write(" = ")
		p.write(p.expr(expr))
	}
}

func (p *SourcePrinter) printTypeDefinition(node *BLangTypeDefinition) {
	p.printDocumentation(node.markdownDocumentationAttachment)
	p.writeQualifiers(node.flags, model.FlagPublic)
	p.write("type ")
	p.write(p.identifier(node.Name))
	p.write(" ")
	p.write(p.typeData(node.typeData))
	p.write(";")
}

func (p *SourcePrinter) printClassDefinition(node *BLangClassDefinition) {
	p.printDocumentation(node.MarkdownDocumentationAttachment)
	p.writeQualifiers(node.Flags(), model.FlagPublic, model.FlagIsolated, model.FlagDistinct, model.FlagClient, model.FlagService, model.FlagReadonly)
	p.write("class ")
	p.write(p.identifier(node.Name))
	p.write(" ")
	p.writeBlock(func() {
		for i, field := range node.Fields {
			if i > 0 {
				p.newline()
			}
			p.indent()
			p.printSimpleVariable(field.(*BLangSimpleVariable), true)
		}
		methodNames := make([]string, 0, len(node.Methods))
		for name := range node.Methods {
			methodNames = append(methodNames, name)
		}
		sort.Strings(methodNames)
		for _, name := range methodNames {
			if len(node.Fields) > 0 || name != methodNames[0] {
				p.newline()
			}
			p.indent()
			p.printFunction(node.Methods[name])
		}
	})
}

func (p *SourcePrinter) printFunction(node *BLangFunction) {
	p.printDocumentation(node.MarkdownDocumentationAttachment)
	p.writeQualifiers(node.Flags(), model.FlagPublic, model.FlagIsolated, model.FlagTransactional, model.FlagRemote, model.FlagResource)
	p.write("function ")
	p.write(p.identifier(&node.Name))
	p.write("(")
	params := make([]string, 0, len(node.RequiredParams)+1)
	for i := range node.RequiredParams {
		params = append(params, p.param(&node.RequiredParams[i]))
	}
	if node.RestParam != nil {
		params = append(params, p.param(node.RestParam.(*BLangSimpleVariable)))
	}
	p.write(strings.Join(params, ", "))
	p.write(")")
	if p.shouldPrintReturnType(node.GetReturnTypeDescriptor()) {
		p.write(" returns ")
		p.write(p.typ(node.GetReturnTypeDescriptor().(BType)))
	}
	if node.Body == nil {
		p.write(";")
		return
	}
	p.write(" ")
	p.printNode(node.Body.(BLangNode))
}

func (p *SourcePrinter) printBlockFunctionBody(node *BLangBlockFunctionBody) {
	p.writeBlock(func() {
		for i, stmt := range node.Stmts {
			if i > 0 {
				p.newline()
			}
			p.indent()
			p.printNode(stmt.(BLangNode))
		}
	})
}

func (p *SourcePrinter) printBlockStmt(node *BLangBlockStmt) {
	p.writeBlock(func() {
		for i, stmt := range node.Stmts {
			if i > 0 {
				p.newline()
			}
			p.indent()
			p.printNode(stmt.(BLangNode))
		}
	})
}

func (p *SourcePrinter) printSimpleVariableDef(node *BLangSimpleVariableDef) {
	p.printSimpleVariable(node.Var, true)
}

func (p *SourcePrinter) printAssignment(node *BLangAssignment) {
	p.write(p.expr(node.VarRef.(BLangActionOrExpression)))
	p.write(" = ")
	p.write(p.expr(node.Expr))
	p.write(";")
}

func (p *SourcePrinter) printCompoundAssignment(node *BLangCompoundAssignment) {
	p.write(p.expr(node.VarRef.(BLangActionOrExpression)))
	p.write(" ")
	p.write(string(node.OpKind))
	p.write("= ")
	p.write(p.expr(node.Expr))
	p.write(";")
}

func (p *SourcePrinter) printIf(node *BLangIf) {
	p.write("if ")
	p.write(p.expr(node.Expr))
	p.write(" ")
	p.printBlockStmt(&node.Body)
	if node.ElseStmt != nil {
		p.write(" else ")
		p.printNode(node.ElseStmt.(BLangNode))
	}
}

func (p *SourcePrinter) printWhile(node *BLangWhile) {
	p.write("while ")
	p.write(p.expr(node.Expr))
	p.write(" ")
	p.printBlockStmt(&node.Body)
}

func (p *SourcePrinter) printForeach(node *BLangForeach) {
	p.write("foreach ")
	if node.VariableDef != nil && node.VariableDef.Var != nil {
		p.printVariableParts(node.VariableDef.Var.TypeNode(), node.VariableDef.Var.Name, nil)
	}
	p.write(" in ")
	p.write(p.expr(node.Collection))
	p.write(" ")
	p.printBlockStmt(&node.Body)
}

func (p *SourcePrinter) printMatchStatement(node *BLangMatchStatement) {
	p.write("match ")
	p.write(p.expr(node.Expr))
	p.write(" ")
	p.writeBlock(func() {
		for i := range node.MatchClauses {
			if i > 0 {
				p.newline()
			}
			p.indent()
			patterns := make([]string, 0, len(node.MatchClauses[i].Patterns))
			for _, pattern := range node.MatchClauses[i].Patterns {
				patterns = append(patterns, p.matchPattern(pattern))
			}
			p.write(strings.Join(patterns, " | "))
			if node.MatchClauses[i].Guard != nil {
				p.write(" if ")
				p.write(p.expr(node.MatchClauses[i].Guard))
			}
			p.write(" => ")
			p.printBlockStmt(&node.MatchClauses[i].Body)
		}
	})
}

func (p *SourcePrinter) expr(node BLangActionOrExpression) string {
	if node == nil {
		return ""
	}
	switch t := node.(type) {
	case *BLangSimpleVarRef:
		if t.PkgAlias != nil && t.PkgAlias.Value != "" {
			return p.identifier(t.PkgAlias) + ":" + p.identifier(t.VariableName)
		}
		return p.identifier(t.VariableName)
	case *BLangLocalVarRef:
		return p.expr(&t.BLangSimpleVarRef)
	case *BLangConstRef:
		if t.OriginalValue != "" {
			return t.OriginalValue
		}
		return p.expr(&t.BLangSimpleVarRef)
	case *BLangLiteral:
		return p.literal(t)
	case *BLangNumericLiteral:
		return p.literal(&t.BLangLiteral)
	case *BLangBinaryExpr:
		return p.expr(t.LhsExpr) + " " + string(t.OpKind) + " " + p.expr(t.RhsExpr)
	case *BLangUnaryExpr:
		return string(t.Operator) + p.expr(t.Expr)
	case *BLangGroupExpr:
		return "(" + p.expr(t.Expression) + ")"
	case *BLangInvocation:
		return p.invocation(&t.bLangInvocationBase, t.PkgAlias, t.Async, false)
	case *BLangRemoteMethodCallAction:
		return p.invocation(&t.bLangInvocationBase, nil, false, true)
	case *BLangNamedArgsExpression:
		return p.identifier(&t.Name) + " = " + p.expr(t.Expr)
	case *BLangIndexBasedAccess:
		return p.expr(t.Expr) + "[" + p.expr(t.IndexExpr) + "]"
	case *BLangFieldBaseAccess:
		return p.expr(t.Expr) + "." + p.identifier(&t.Field)
	case *BLangListConstructorExpr:
		parts := make([]string, 0, len(t.Exprs))
		for i, expr := range t.Exprs {
			part := p.expr(expr)
			if i < len(t.SpreadMembers) && t.SpreadMembers[i] {
				part = "..." + part
			}
			parts = append(parts, part)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case *BLangMappingConstructorExpr:
		parts := make([]string, 0, len(t.Fields))
		for _, field := range t.Fields {
			parts = append(parts, p.mappingField(field))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case *BLangTypeConversionExpr:
		return "<" + p.typ(t.TypeDescriptor) + ">" + p.expr(t.Expression)
	case *BLangTypeTestExpr:
		op := " is "
		if t.isNegation {
			op = " !is "
		}
		return p.expr(t.Expr) + op + p.typeData(t.Type)
	case *BLangCheckedExpr:
		return "check " + p.expr(t.Expr)
	case *BLangCheckPanickedExpr:
		return "checkpanic " + p.expr(t.Expr)
	case *BLangTrapExpr:
		return "trap " + p.expr(t.Expr)
	case *BLangNewExpression:
		return p.newExpression(t)
	case *BLangLambdaFunction:
		return p.lambda(t)
	case *BLangErrorConstructorExpr:
		return p.errorConstructor(t)
	case *BLangTemplateExpr:
		return p.templateExpr(t)
	default:
		panic("unsupported expression for source printing: " + reflect.TypeOf(t).String())
	}
}

func (p *SourcePrinter) invocation(node *bLangInvocationBase, pkgAlias *BLangIdentifier, async bool, remote bool) string {
	var b strings.Builder
	if async {
		b.WriteString("start ")
	}
	if node.Expr != nil {
		b.WriteString(p.expr(node.Expr))
		if remote {
			b.WriteString("->")
		} else {
			b.WriteString(".")
		}
	} else if pkgAlias != nil && pkgAlias.Value != "" {
		b.WriteString(p.identifier(pkgAlias))
		b.WriteString(":")
	}
	b.WriteString(p.identifier(node.Name))
	b.WriteString("(")
	args := make([]string, 0, len(node.ArgExprs))
	for _, arg := range node.ArgExprs {
		args = append(args, p.expr(arg))
	}
	b.WriteString(strings.Join(args, ", "))
	b.WriteString(")")
	return b.String()
}

func (p *SourcePrinter) newExpression(node *BLangNewExpression) string {
	args := make([]string, 0, len(node.ArgsExprs))
	for _, arg := range node.ArgsExprs {
		args = append(args, p.expr(arg))
	}
	if node.TypeDescriptor != nil {
		return "new " + p.typ(node.TypeDescriptor) + "(" + strings.Join(args, ", ") + ")"
	}
	return "new (" + strings.Join(args, ", ") + ")"
}

func (p *SourcePrinter) lambda(node *BLangLambdaFunction) string {
	if node.Function == nil {
		return "function () {}"
	}
	printer := SourcePrinter{}
	printer.printFunction(node.Function)
	return printer.buffer.String()
}

func (p *SourcePrinter) errorConstructor(node *BLangErrorConstructorExpr) string {
	args := make([]string, 0, len(node.PositionalArgs)+len(node.NamedArgs))
	for _, arg := range node.PositionalArgs {
		args = append(args, p.expr(arg))
	}
	for i := range node.NamedArgs {
		args = append(args, p.expr(&node.NamedArgs[i]))
	}
	name := "error"
	if node.ErrorTypeRef != nil {
		name = p.typ(node.ErrorTypeRef)
	}
	return name + "(" + strings.Join(args, ", ") + ")"
}

func (p *SourcePrinter) templateExpr(node *BLangTemplateExpr) string {
	if node.Kind != TemplateExprKindString {
		panic("unsupported template expression kind for source printing")
	}
	var b strings.Builder
	b.WriteByte('`')
	for i, s := range node.Strings {
		b.WriteString(strings.ReplaceAll(s, "`", "\\`"))
		if i < len(node.Insertions) {
			b.WriteString("${")
			b.WriteString(p.expr(node.Insertions[i]))
			b.WriteByte('}')
		}
	}
	b.WriteByte('`')
	return b.String()
}

func (p *SourcePrinter) mappingField(field MappingField) string {
	switch t := field.(type) {
	case *BLangMappingKeyValueField:
		key := p.expr(t.Key.Expr)
		if t.Key.Kind == MappingKeyComputed {
			key = "[" + key + "]"
		}
		return key + ": " + p.expr(t.ValueExpr)
	case *BLangConstRef:
		return p.expr(t)
	default:
		panic("unsupported mapping field for source printing: " + reflect.TypeOf(t).String())
	}
}

func (p *SourcePrinter) typ(node BType) string {
	if node == nil {
		return ""
	}
	var source string
	switch t := node.(type) {
	case *BLangValueType:
		if t.TypeKind == TypeKind_NIL {
			source = "()"
		} else {
			source = string(t.TypeKind)
		}
	case *BLangBuiltInRefTypeNode:
		source = string(t.TypeKind)
	case *BLangUserDefinedType:
		if t.PkgAlias.Value != "" {
			source = p.identifier(&t.PkgAlias) + ":" + p.identifier(&t.TypeName)
		} else {
			source = p.identifier(&t.TypeName)
		}
	case *BLangArrayType:
		source = p.typeData(t.Elemtype)
		for i := 0; i < t.Dimensions; i++ {
			if i < len(t.Sizes) && t.Sizes[i] != nil {
				source += "[" + p.expr(t.Sizes[i]) + "]"
			} else {
				source += "[]"
			}
		}
	case *BLangUnionTypeNode:
		source = p.typeData(t.lhs) + "|" + p.typeData(t.rhs)
	case *BLangIntersectionTypeNode:
		source = p.typeData(t.lhs) + " & " + p.typeData(t.rhs)
	case *BLangErrorTypeNode:
		source = "error"
		if t.DetailType.TypeDescriptor != nil {
			source += "<" + p.typeData(t.DetailType) + ">"
		}
	case *BLangConstrainedType:
		source = p.typeData(t.Type) + "<" + p.typeData(t.Constraint) + ">"
	case *BLangStreamType:
		source = "stream<" + p.typeData(t.ValueType)
		if t.CompletionType.TypeDescriptor != nil {
			source += ", " + p.typeData(t.CompletionType)
		}
		source += ">"
	case *BLangTupleTypeNode:
		parts := make([]string, 0, len(t.Members)+1)
		for i := range t.Members {
			parts = append(parts, p.typeDescriptor(t.Members[i].TypeDesc))
		}
		if t.Rest != nil {
			parts = append(parts, p.typ(t.Rest)+"...")
		}
		source = "[" + strings.Join(parts, ", ") + "]"
	case *BLangRecordType:
		source = p.recordType(t)
	case *BLangObjectType:
		source = p.objectType(t)
	case *BLangFiniteTypeNode:
		parts := make([]string, 0, len(t.ValueSpace))
		for _, value := range t.ValueSpace {
			parts = append(parts, p.expr(value))
		}
		source = strings.Join(parts, "|")
	case *BLangFunctionType:
		source = p.functionType(t)
	default:
		panic("unsupported type for source printing: " + reflect.TypeOf(t).String())
	}
	if base, ok := any(node).(interface{ bTypeGetFlags() model.Flag }); ok && base.bTypeGetFlags().Has(model.FlagReadonly) {
		source += " & readonly"
	}
	return source
}

func (p *SourcePrinter) typeData(typeData TypeData) string {
	return p.typeDescriptor(typeData.TypeDescriptor)
}

func (p *SourcePrinter) typeDescriptor(typeDescriptor TypeDescriptor) string {
	if typeDescriptor == nil {
		return ""
	}
	return p.typ(typeDescriptor.(BType))
}

func (p *SourcePrinter) recordType(node *BLangRecordType) string {
	var b strings.Builder
	b.WriteString("record {| ")
	first := true
	for _, inclusion := range node.TypeInclusions {
		if !first {
			b.WriteByte(' ')
		}
		first = false
		b.WriteString("*")
		b.WriteString(p.typ(inclusion))
		b.WriteString(";")
	}
	for _, field := range node.fields {
		if !first {
			b.WriteByte(' ')
		}
		first = false
		b.WriteString(p.typ(field.Type))
		b.WriteByte(' ')
		b.WriteString(string(field.Name))
		if field.DefaultExpr != nil {
			b.WriteString(" = ")
			b.WriteString(p.expr(field.DefaultExpr))
		}
		b.WriteString(";")
	}
	if node.RestType != nil {
		if !first {
			b.WriteByte(' ')
		}
		b.WriteString(p.typ(node.RestType))
		b.WriteString("...;")
	} else if node.IsOpen {
		if !first {
			b.WriteByte(' ')
		}
		b.WriteString("anydata...;")
	}
	b.WriteString(" |}")
	return b.String()
}

func (p *SourcePrinter) objectType(node *BLangObjectType) string {
	var b strings.Builder
	b.WriteString("object {")
	first := true
	for member := range node.Members() {
		if !first {
			b.WriteByte(' ')
		}
		first = false
		switch m := member.(type) {
		case *BObjectField:
			b.WriteString(p.typ(m.Ty))
			b.WriteByte(' ')
			b.WriteString(m.Name())
			b.WriteByte(';')
		case *BMethodDecl:
			b.WriteString(p.functionType(&m.BLangFunctionType))
			b.WriteByte(' ')
			b.WriteString(m.Name())
			b.WriteByte(';')
		}
	}
	b.WriteString("}")
	return b.String()
}

func (p *SourcePrinter) functionType(node *BLangFunctionType) string {
	params := make([]string, 0, len(node.RequiredParams)+1)
	for i := range node.RequiredParams {
		param := node.RequiredParams[i]
		if param.Name != nil && param.Name.Value != "" {
			params = append(params, p.typ(param.TypeDesc)+" "+p.identifier(param.Name))
		} else {
			params = append(params, p.typ(param.TypeDesc))
		}
	}
	if node.RestParam != nil {
		params = append(params, p.typ(node.RestParam.TypeDesc)+"...")
	}
	source := "function (" + strings.Join(params, ", ") + ")"
	if p.shouldPrintReturnType(node.ReturnTypeDescriptor) {
		source += " returns " + p.typ(node.ReturnTypeDescriptor)
	}
	return source
}

func (p *SourcePrinter) shouldPrintReturnType(typeDescriptor TypeDescriptor) bool {
	if typeDescriptor == nil {
		return false
	}
	valueType, ok := typeDescriptor.(*BLangValueType)
	return !ok || (valueType.TypeKind != TypeKind_NIL && valueType.TypeKind != TypeKind_VOID && valueType.TypeKind != TypeKind_NONE)
}

func (p *SourcePrinter) param(node *BLangSimpleVariable) string {
	prefix := ""
	if node.IsRestParam() {
		prefix = "..."
	}
	text := p.typ(node.TypeNode()) + prefix
	if node.Name != nil && node.Name.Value != "" {
		text += " " + p.identifier(node.Name)
	}
	if node.Expr != nil {
		text += " = " + p.expr(node.Expr)
	}
	return text
}

func (p *SourcePrinter) matchPattern(node MatchPatternNode) string {
	switch t := node.(type) {
	case *BLangConstPattern:
		return p.expr(t.Expr)
	case *BLangWildCardMatchPattern:
		return "_"
	default:
		panic("unsupported match pattern for source printing: " + reflect.TypeOf(t).String())
	}
}

func (p *SourcePrinter) literal(node *BLangLiteral) string {
	if node.OriginalValue != "" {
		return node.OriginalValue
	}
	switch value := node.Value.(type) {
	case nil:
		return "()"
	case string:
		return strconv.Quote(value)
	case bool:
		return strconv.FormatBool(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (p *SourcePrinter) identifier(node *BLangIdentifier) string {
	if node == nil {
		return ""
	}
	if node.OriginalValue != "" {
		return node.OriginalValue
	}
	return node.Value
}

func (p *SourcePrinter) printDocumentationNode(node MarkdownDocumentationNode) {
	if node == nil {
		return
	}
	doc, ok := node.(*BLangMarkdownDocumentation)
	if !ok {
		return
	}
	p.printDocumentation(doc)
}

func (p *SourcePrinter) printDocumentation(node *BLangMarkdownDocumentation) {
	if node == nil {
		return
	}
	for _, line := range node.DocumentationLines {
		p.write("#")
		if line.Text != "" {
			p.write(" ")
			p.write(line.Text)
		}
		p.newline()
		p.indent()
	}
}

func (p *SourcePrinter) writeQualifiers(flags model.Flag, allowed ...model.Flag) {
	qualifiers := []struct {
		flag model.Flag
		text string
	}{
		{model.FlagPublic, "public"},
		{model.FlagIsolated, "isolated"},
		{model.FlagTransactional, "transactional"},
		{model.FlagRemote, "remote"},
		{model.FlagResource, "resource"},
		{model.FlagFinal, "final"},
		{model.FlagConfigurable, "configurable"},
		{model.FlagDistinct, "distinct"},
		{model.FlagClient, "client"},
		{model.FlagService, "service"},
		{model.FlagReadonly, "readonly"},
	}
	allowedSet := map[model.Flag]bool{}
	for _, flag := range allowed {
		allowedSet[flag] = true
	}
	for _, qualifier := range qualifiers {
		if allowedSet[qualifier.flag] && flags.Has(qualifier.flag) {
			p.write(qualifier.text)
			p.write(" ")
		}
	}
}

func (p *SourcePrinter) writeBlock(body func()) {
	p.write("{")
	p.indentLevel++
	if body != nil {
		p.newline()
		body()
	}
	p.indentLevel--
	p.newline()
	p.indent()
	p.write("}")
}

func (p *SourcePrinter) write(str string) {
	p.buffer.WriteString(str)
}

func (p *SourcePrinter) newline() {
	p.buffer.WriteByte('\n')
}

func (p *SourcePrinter) indent() {
	for i := 0; i < p.indentLevel; i++ {
		p.buffer.WriteByte('\t')
	}
}

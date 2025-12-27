# BIR Model Package

This package contains the Ballerina Intermediate Representation (BIR) model types, migrated from Java to Go.

## Overview

The BIR (Ballerina Intermediate Representation) is an intermediate representation of Ballerina programs used in the compilation pipeline. This package defines the core data structures that represent BIR programs.

## Package Structure

```
compiler/bir/model/
├── bir-node.go       # Main BIR node types
├── bir-scope.go      # Scope tracking
├── var-kind.go       # Variable kind enumeration
├── var-scope.go      # Variable scope enumeration
└── bir-node_test.go  # Unit tests
```

## Main Types

### Base Types

- **BIRNode**: Base interface for all BIR nodes, defines visitor pattern support
- **BIRDocumentableNode**: Interface for nodes that can have markdown documentation
- **BIRVisitor**: Visitor interface for traversing BIR trees

### Core BIR Types

- **BIRPackage**: Represents a Ballerina package with all its components
- **BIRFunction**: Function definition with parameters, local variables, and basic blocks
- **BIRBasicBlock**: A basic block containing instructions and a terminator
- **BIRTypeDefinition**: Type definitions with attached functions
- **BIRVariableDcl**: Variable declaration with scope and kind information

### Variable Types

- **BIRParameter**: Function parameter
- **BIRFunctionParameter**: Function parameter with default expression support
- **BIRGlobalVariableDcl**: Global variable declaration

### Annotation Types

- **BIRAnnotation**: Annotation definition
- **BIRAnnotationAttachment**: Annotation usage/attachment
- **BIRConstAnnotationAttachment**: Constant annotation attachment with value

### Supporting Types

- **BIRImportModule**: Imported module reference
- **BIRErrorEntry**: Error handling entry in error table
- **BIRConstant**: Constant definition
- **BIRServiceDeclaration**: Service declaration
- **ChannelDetails**: Worker channel information
- **ConstValue**: Constant value container

### Constructor Entry Types

For mapping and list constructors:
- **BIRMappingConstructorEntry**: Interface for mapping constructor entries
- **BIRMappingConstructorKeyValueEntry**: Key-value pair entry
- **BIRMappingConstructorSpreadFieldEntry**: Spread field entry
- **BIRListConstructorEntry**: Interface for list constructor entries
- **BIRListConstructorSpreadMemberEntry**: Spread member entry
- **BIRListConstructorExprEntry**: Expression entry

## Enumerations

### VarKind
Variable kinds in BIR:
- `VarKindLocal`: User-defined local variable
- `VarKindArg`: Function argument
- `VarKindTemp`: Temporary variable for sub-expressions
- `VarKindReturn`: Special variable for return value
- `VarKindGlobal`: User-defined global variable
- `VarKindSelf`: Self-referencing variable
- `VarKindConstant`: Constant variable
- `VarKindSynthetic`: Compiler-generated variable

### VarScope
Variable scopes:
- `VarScopeFunction`: Function scope
- `VarScopeGlobal`: Global scope

## Usage Examples

### Creating a Package

```go
import (
    "ballerina-lang-go/compiler/bir/model"
    "ballerina-lang-go/compiler/util"
)

pkg := model.NewBIRPackage(
    nil,
    util.NewName("myorg"),
    util.NewName("mypackage"),
    util.NewName("mymodule"),
    util.NewName("1.0.0"),
    util.NewName("main.bal"),
    "/src",
    false,
)
```

### Creating a Function

```go
fn := model.NewBIRFunctionMinimal(
    nil,
    util.NewName("myFunction"),
    0,
    nil,
    util.NewName("worker"),
    0,
    symbols.SymbolOriginSource,
)
```

### Creating a Variable Declaration

```go
varDecl := model.NewBIRVariableDclMinimal(
    nil,
    util.NewName("myVar"),
    model.VarScopeFunction,
    model.VarKindLocal,
)
```

### Using the Visitor Pattern

```go
type MyVisitor struct{}

func (v *MyVisitor) VisitPackage(pkg *model.BIRPackage) {
    // Visit package
}

func (v *MyVisitor) VisitFunction(fn *model.BIRFunction) {
    // Visit function
}

// ... implement other visitor methods

visitor := &MyVisitor{}
pkg.Accept(visitor)
```

## Memory Optimization

This package is designed with memory efficiency in mind:

1. **Value Types**: Small immutable types (`Name`, `VarKind`, `VarScope`) use value semantics to reduce heap allocations
2. **Judicious Pointers**: Pointers are used only when necessary for mutability, optional fields, or large structures
3. **Struct Embedding**: Base types are embedded by value to reduce pointer indirection
4. **Pre-allocated Collections**: Collections are pre-allocated in constructors to avoid nil checks

## Testing

Run tests with:
```bash
go test ./compiler/bir/model -v
```

All tests should pass with zero failures.

## Dependencies

This package depends on:
- `ballerina-lang-go/compiler/util` - Name type
- `ballerina-lang-go/compiler/semantics/model/elements` - PackageID, AttachPoint, MarkdownDocAttachment
- `ballerina-lang-go/compiler/semantics/model/symbols` - SymbolOrigin
- `ballerina-lang-go/compiler/semantics/model/types` - BType, BInvokableType
- `ballerina-lang-go/tools/diagnostics` - Location

## Migration Notes

This package is a direct migration from the Java class `org.wso2.ballerinalang.compiler.bir.model.BIRNode`.

Key migration decisions:
- Java classes → Go interfaces + unexported structs
- Java inheritance → Go struct embedding
- Java enums → Go typed constants with iota
- Java collections → Go slices and maps

See [docs/migration/BIRNode-migration.md](../../../docs/migration/BIRNode-migration.md) for detailed migration documentation.

## Future Work

The following referenced types are defined as placeholder interfaces and need full implementation:
- `BIROperand`
- `BIRNonTerminator`
- `BIRTerminator`
- `BIRLock`

These will be implemented when the instruction classes are migrated.

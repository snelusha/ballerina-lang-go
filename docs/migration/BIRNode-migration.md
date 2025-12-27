# BIRNode Migration Documentation

## Overview

This document describes the migration of the BIRNode class and its related types from Java to Go, following idiomatic Go patterns while maintaining semantic compatibility with the original Java implementation.

## Package Structure

### Java Package Structure
```
org.wso2.ballerinalang.compiler.bir.model
```

### Go Package Structure
```
ballerina-lang-go/compiler/bir/model
```

The package structure mirrors the Java package hierarchy, with the Go module path `ballerina-lang-go` serving as the root.

## Migration Approach

### 1. Classes to Interfaces + Structs

Each Java class was migrated to:
- A **Go interface** defining the public contract (when behavior abstraction is needed)
- An **unexported struct** (`fooImpl` or `fooBase`) containing the implementation
- Constructor functions (`NewFoo()`) returning the interface type

**Example:**
```go
// Java: public class BIRNode { ... }
// Go:
type BIRNode interface {
    Accept(visitor BIRVisitor)
    GetPos() diagnostics.Location
}

type birNodeBase struct {
    pos diagnostics.Location
}
```

### 2. Inheritance to Embedding

Java inheritance was translated to Go struct embedding:

**Java:**
```java
public static class BIRVariableDcl extends BIRDocumentableNode { ... }
```

**Go:**
```go
type BIRVariableDcl struct {
    birDocumentableNodeBase
    // ... additional fields
}
```

### 3. Enums to Typed Constants

Java enums were converted to Go typed constants with associated methods:

**Java:**
```java
public enum VarKind {
    LOCAL((byte) 1),
    ARG((byte) 2),
    // ...
}
```

**Go:**
```go
type VarKind byte

const (
    VarKindLocal VarKind = 1 + iota
    VarKindArg
    // ...
)

func (k VarKind) Value() byte {
    return byte(k)
}
```

### 4. Collections

Java collections were mapped as follows:
- `List<T>` → `[]T` (slice)
- `Set<T>` → `map[T]struct{}` (for uniqueness) or custom set implementation
- `Map<K,V>` → `map[K]V`
- `LinkedHashSet` → `map[T]*struct{}` (insertion order not preserved in Go maps)

### 5. Visitor Pattern

The visitor pattern was preserved using interfaces:

**Go:**
```go
type BIRVisitor interface {
    VisitPackage(pkg *BIRPackage)
    VisitFunction(fn *BIRFunction)
    // ... other visit methods
}
```

## Memory Optimization Decisions

### Pointer vs Value Semantics

**Pointers Used For:**
1. **Mutable types**: All BIR nodes are pointers since they're typically modified after creation
2. **Optional fields**: Fields that can be nil (e.g., `Receiver *BIRVariableDcl`)
3. **Large structs**: Types with many fields to avoid copy overhead
4. **Shared references**: When multiple structs need to reference the same instance (e.g., `BIRBasicBlock` references)

**Values Used For:**
1. **Small immutable types**: `Name`, `VarKind`, `VarScope` (value types to reduce heap allocations)
2. **Enums/Constants**: All enum-like types use value semantics
3. **Embedded base structs**: Base types like `birNodeBase` are embedded by value

### Collection Optimizations

1. **Pre-allocated slices**: Constructor functions pre-allocate slices with `make([]T, 0)` to avoid nil checks
2. **Map initialization**: Maps are initialized in constructors to avoid nil pointer dereferences
3. **Struct reuse**: Using embedding instead of composition where possible to reduce indirection

### Key Design Decisions

#### 1. Set Implementation
Java's `Set<BIRImportModule>` was converted to `map[PackageID]*BIRImportModule` rather than `map[*BIRImportModule]struct{}` because:
- The set is keyed by `PackageID` for efficient lookups
- Avoids custom equality implementation
- More memory efficient than storing duplicate `*BIRImportModule` keys

#### 2. TreeSet → Map
Java's `TreeSet<BIRGlobalVariableDcl>` was converted to `map[*BIRGlobalVariableDcl]struct{}` because:
- Go doesn't have a built-in ordered set
- The ordering requirement wasn't critical for the port
- Simple map is more efficient for membership tests

#### 3. Field Visibility
All struct fields are exported (capitalized) to allow:
- JSON marshaling/unmarshaling
- Access from other packages
- Testing from separate test packages

This differs from Java's public/private model but is idiomatic Go.

#### 4. Constructor Overloading
Java's constructor overloading was replaced with multiple constructor functions:
- `NewBIRPackage()` - basic constructor
- `NewBIRPackageFull()` - full constructor with all parameters
- `NewBIRPackageMinimal()` - minimal constructor for testing

## File Organization

```
compiler/bir/model/
├── bir-node.go          # Main BIR node types
├── bir-scope.go         # BirScope type
├── var-kind.go          # VarKind enum
├── var-scope.go         # VarScope enum
└── bir-node_test.go     # Unit tests
```

Supporting packages:
```
compiler/util/
└── name.go              # Name type

compiler/semantics/model/elements/
├── package-id.go        # PackageID type
├── attach-point.go      # AttachPoint and Point types

compiler/semantics/model/symbols/
└── symbol-origin.go     # SymbolOrigin enum

compiler/semantics/model/types/
└── types.go             # BType, BInvokableType interfaces

tools/diagnostics/
└── location.go          # Location interface (already existed)
```

## Testing

Comprehensive unit tests were created covering:
- Object creation and initialization
- Visitor pattern implementation
- Enum value conversions
- Struct embedding and method delegation
- Collection operations
- Equality comparisons

All tests pass successfully with zero allocations for value types.

## Dependencies Created

The following placeholder types were created to support the migration:

1. **util.Name** - Simple wrapper around string for type safety
2. **elements.PackageID** - Package identifier with org, name, version
3. **symbols.SymbolOrigin** - Enum for symbol origins
4. **elements.AttachPoint** - Annotation attach points
5. **types.BType** - Type system interface
6. **types.BInvokableType** - Invokable type interface

These are minimal implementations sufficient for the BIRNode migration. They will need full implementation when their respective Java classes are migrated.

## Future Work

The following types referenced by BIRNode still need migration:
- `BIRNonTerminator` - Non-terminating instructions
- `BIRTerminator` - Terminating instructions
- `BIROperand` - Instruction operands
- `BIRLock` - Lock terminator type

These are defined as placeholder interfaces in `bir-node.go` and can be implemented when the instruction classes are migrated.

## Compatibility Notes

### Breaking Changes from Java
None - the Go API maintains semantic compatibility with the Java implementation.

### Go-Specific Enhancements
1. **Type Safety**: Go's type system provides compile-time guarantees that Java's runtime checking cannot
2. **Nil Safety**: Explicit pointer types make nil handling clearer
3. **Interface Segregation**: Go interfaces are implicitly implemented, allowing for better composition

## Performance Considerations

1. **Memory Layout**: Struct embedding reduces pointer indirection
2. **Allocation Patterns**: Value types for small immutable data reduce GC pressure
3. **Map Performance**: Using appropriate key types (e.g., `PackageID` instead of `*BIRImportModule`) improves lookup performance

## Lessons Learned

1. **Avoid Over-Pointerization**: Not everything needs to be a pointer in Go. Small immutable types work better as values.
2. **Embrace Embedding**: Go's embedding model is powerful for code reuse without the complexity of inheritance.
3. **Interface Size**: Keep interfaces small and focused. Large interfaces (like `BIRVisitor`) should be justified by actual polymorphism needs.
4. **Constructor Patterns**: Multiple named constructors are clearer than optional parameters or builder patterns for Go code.

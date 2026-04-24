## Module enumeration declaration

```
module-enum-decl :=
metadata
[public] enum identifier { enum-member (, enum-member)\* } [;]
enum-member := metadata identifier [= const-expr]
```

A module-enum-decl provides a convenient syntax for declaring a union of string constants.

Each enum-member is defined as compile-time constant in the same way as if it had been defined using a module-const-decl. The result of evaluating the const-expr must be a string. If the const-expr is omitted, it defaults to be the same as the identifier.

The identifier is defined as a type in the same was as if it had been defined by a module-type-defn, with the type-descriptor being the union of the constants defined by the members.

If the module-enum-decl is public, then both the type and the constants are public.

So for example:

```
public enum Color {
RED,
GREEN,
BLUE
}
```

is exactly equivalent to:

```
public const RED = "RED";
public const GREEN = "GREEN";
public const BLUE = "BLUE";
public type Color RED|GREEN|BLUE;
```

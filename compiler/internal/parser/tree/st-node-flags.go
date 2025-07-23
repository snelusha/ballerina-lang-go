package tree

// STNodeFlags represents a set of flags that can be attached to an internal syntax node.

const (
	STNodeFlagsHasDiagnostic uint8 = 1 << 0x1
	STNodeFlagsIsMissing     uint8 = 1 << 0x2
)

// STNodeFlagsIsFlagSet checks whether the given flag is set in the given flags.
func STNodeFlagsIsFlagSet(flags uint8, flag uint8) bool {
	return (flags & flag) != 0
}

// STNodeFlagsWithFlag sets a flag in the given flags.
func STNodeFlagsWithFlag(flags uint8, flag uint8) uint8 {
	return flags | flag
}

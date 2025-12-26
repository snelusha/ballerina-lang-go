package util

import "ballerina-lang-go/compiler/model/elements"

const (
	PUBLIC uint64 = 1 << iota
	NATIVE
	FINAL
	ATTACHED
	DEPRECATED
	READONLY
	FUNCTION_FINAL
	INTERFACE
	REQUIRED
	RECORD
	PRIVATE
	ANONYMOUS
	OPTIONAL
	TESTABLE
	CONSTANT
	REMOTE
	CLIENT
	RESOURCE
	SERVICE
	LISTENER
	LAMBDA
	TYPE_PARAM
	LANG_LIB
	WORKER
	FORKED
	TRANSACTIONAL
	PARAMETERIZED
	DISTINCT
	CLASS
	ISOLATED
	ISOLATED_PARAM
	CONFIGURABLE
	OBJECT_CTOR
	ENUM
	INCLUDED
	REQUIRED_PARAM
	DEFAULTABLE_PARAM
	REST_PARAM
	FIELD
	ANY_FUNCTION
	INFER
	ENUM_MEMBER
	QUERY_LAMBDA
	EFFECTIVE_TYPE_DEF
	SOURCE_ANNOTATION
)

var flagToMask = map[elements.Flag]uint64{
	elements.PUBLIC:            PUBLIC,
	elements.PRIVATE:           PRIVATE,
	elements.REMOTE:            REMOTE,
	elements.TRANSACTIONAL:     TRANSACTIONAL,
	elements.NATIVE:            NATIVE,
	elements.FINAL:             FINAL,
	elements.ATTACHED:          ATTACHED,
	elements.LAMBDA:            LAMBDA,
	elements.WORKER:            WORKER,
	elements.LISTENER:          LISTENER,
	elements.READONLY:          READONLY,
	elements.FUNCTION_FINAL:    FUNCTION_FINAL,
	elements.INTERFACE:         INTERFACE,
	elements.REQUIRED:          REQUIRED,
	elements.RECORD:            RECORD,
	elements.ANONYMOUS:         ANONYMOUS,
	elements.OPTIONAL:          OPTIONAL,
	elements.TESTABLE:          TESTABLE,
	elements.CLIENT:            CLIENT,
	elements.RESOURCE:          RESOURCE,
	elements.ISOLATED:          ISOLATED,
	elements.SERVICE:           SERVICE,
	elements.CONSTANT:          CONSTANT,
	elements.TYPE_PARAM:        TYPE_PARAM,
	elements.LANG_LIB:          LANG_LIB,
	elements.FORKED:            FORKED,
	elements.DISTINCT:          DISTINCT,
	elements.CLASS:             CLASS,
	elements.CONFIGURABLE:      CONFIGURABLE,
	elements.OBJECT_CTOR:       OBJECT_CTOR,
	elements.ENUM:              ENUM,
	elements.INCLUDED:          INCLUDED,
	elements.REQUIRED_PARAM:    REQUIRED_PARAM,
	elements.DEFAULTABLE_PARAM: DEFAULTABLE_PARAM,
	elements.REST_PARAM:        REST_PARAM,
	elements.FIELD:             FIELD,
	elements.ANY_FUNCTION:      ANY_FUNCTION,
	elements.ENUM_MEMBER:       ENUM_MEMBER,
	elements.QUERY_LAMBDA:      QUERY_LAMBDA,
}

func AsMask(flagSet map[elements.Flag]any) uint64 {
	var mask uint64
	for flag := range flagSet {
		if flagVal, ok := flagToMask[flag]; ok {
			mask |= flagVal
		}
	}
	return mask
}

func AsMaskFromSlice(flags []elements.Flag) uint64 {
	var mask uint64
	for _, flag := range flags {
		if flagVal, ok := flagToMask[flag]; ok {
			mask |= flagVal
		}
	}
	return mask
}

func UnMask(mask uint64) map[elements.Flag]any {
	flagSet := make(map[elements.Flag]any)

	for flag, flagVal := range flagToMask {
		if (mask & flagVal) == flagVal {
			flagSet[flag] = struct{}{}
		}
	}

	return flagSet
}

func UnMaskToSlice(mask uint64) []elements.Flag {
	flags := make([]elements.Flag, 0)

	for flag, flagVal := range flagToMask {
		if (mask & flagVal) == flagVal {
			flags = append(flags, flag)
		}
	}

	return flags
}

func Unset(mask, flag uint64) uint64 {
	return mask & (^flag)
}

func AddIfFlagOn(flagSet map[elements.Flag]any, mask, flagVal uint64, flag elements.Flag) {
	if (mask & flagVal) == flagVal {
		flagSet[flag] = struct{}{}
	}
}

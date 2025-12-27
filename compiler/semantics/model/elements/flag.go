package elements

type Flag int

const (
	FlagPublic Flag = 1 << iota
	FlagPrivate
	FlagRemote
	FlagTransactional
	FlagNative
	FlagFinal
	FlagAttached
	FlagLambda
	FlagWorker
	FlagParallel
	FlagListener
	FlagReadonly
	FlagFunctionFinal
	FlagInterface
	FlagRequired
	FlagRecord
	FlagAnonymous
	FlagOptional
	FlagTestable
	FlagClient
	FlagResource
	FlagIsolated
	FlagService
	FlagConstant
	FlagTypeParam
	FlagLangLib
	FlagForked
	FlagDistinct
	FlagClass
	FlagConfigurable
	FlagObjectCtor
	FlagEnum
	FlagIncluded
	FlagRequiredParam
	FlagDefaultableParam
	FlagRestParam
	FlagField
	FlagAnyFunction
	FlagNeverAllowed
	FlagEnumMember
	FlagQueryLambda
)

func UnmaskFlags(flags int64) map[Flag]struct{} {
	result := make(map[Flag]struct{})

	allFlags := []Flag{
		FlagPublic, FlagPrivate, FlagRemote, FlagTransactional, FlagNative,
		FlagFinal, FlagAttached, FlagLambda, FlagWorker, FlagParallel,
		FlagListener, FlagReadonly, FlagFunctionFinal, FlagInterface, FlagRequired,
		FlagRecord, FlagAnonymous, FlagOptional, FlagTestable, FlagClient,
		FlagResource, FlagIsolated, FlagService, FlagConstant, FlagTypeParam,
		FlagLangLib, FlagForked, FlagDistinct, FlagClass, FlagConfigurable,
		FlagObjectCtor, FlagEnum, FlagIncluded, FlagRequiredParam, FlagDefaultableParam,
		FlagRestParam, FlagField, FlagAnyFunction, FlagNeverAllowed, FlagEnumMember,
		FlagQueryLambda,
	}

	for _, flag := range allFlags {
		if flags&int64(flag) != 0 {
			result[flag] = struct{}{}
		}
	}

	return result
}

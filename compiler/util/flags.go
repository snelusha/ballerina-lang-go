package util

const (
	FlagPublic           int64 = 1 << 0
	FlagNative           int64 = 1 << 1
	FlagFinal            int64 = 1 << 2
	FlagAttached         int64 = 1 << 3
	FlagDeprecated       int64 = 1 << 4
	FlagReadonly         int64 = 1 << 5
	FlagFunctionFinal    int64 = 1 << 6
	FlagInterface        int64 = 1 << 7
	FlagRequired         int64 = 1 << 8
	FlagRecord           int64 = 1 << 9
	FlagPrivate          int64 = 1 << 10
	FlagAnonymous        int64 = 1 << 11
	FlagOptional         int64 = 1 << 12
	FlagTestable         int64 = 1 << 13
	FlagConstant         int64 = 1 << 14
	FlagRemote           int64 = 1 << 15
	FlagClient           int64 = 1 << 16
	FlagResource         int64 = 1 << 17
	FlagService          int64 = 1 << 18
	FlagListener         int64 = 1 << 19
	FlagLambda           int64 = 1 << 20
	FlagTypeParam        int64 = 1 << 21
	FlagLangLib          int64 = 1 << 22
	FlagWorker           int64 = 1 << 23
	FlagForked           int64 = 1 << 24
	FlagTransactional    int64 = 1 << 25
	FlagParameterized    int64 = 1 << 26
	FlagDistinct         int64 = 1 << 27
	FlagClass            int64 = 1 << 28
	FlagIsolated         int64 = 1 << 29
	FlagIsolatedParam    int64 = 1 << 30
	FlagConfigurable     int64 = 1 << 31
	FlagObjectCtor       int64 = 1 << 32
	FlagEnum             int64 = 1 << 33
	FlagIncluded         int64 = 1 << 34
	FlagRequiredParam    int64 = 1 << 35
	FlagDefaultableParam int64 = 1 << 36
	FlagRestParam        int64 = 1 << 37
	FlagField            int64 = 1 << 38
	FlagAnyFunction      int64 = 1 << 39
	FlagInfer            int64 = 1 << 40
	FlagEnumMember       int64 = 1 << 41
	FlagQueryLambda      int64 = 1 << 42
	FlagEffectiveTypeDef int64 = 1 << 43
	FlagSourceAnnotation int64 = 1 << 44
)

func IsFlagOn(mask, flag int64) bool {
	return (mask & flag) == flag
}

func Unset(mask, flag int64) int64 {
	return mask & (^flag)
}

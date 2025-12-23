package model

type InstructionKind byte

const (
	InstructionKindGoto InstructionKind = 1 + iota
	InstructionKindCall
	InstructionKindBranch
	InstructionKindReturn
	InstructionKindAsyncCall
	InstructionKindWait
	InstructionKindFPCall
	InstructionKindWKReceive
	InstructionKindWKSend
	InstructionKindFlush
	InstructionKindLock
	InstructionKindFieldLock
	InstructionKindUnlock
	InstructionKindWaitAll
	InstructionKindWKAltReceive
	InstructionKindWKMultipleReceive

	InstructionKindMove              InstructionKind = 20
	InstructionKindConstLoad         InstructionKind = 21
	InstructionKindNewStructure      InstructionKind = 22
	InstructionKindMapStore          InstructionKind = 23
	InstructionKindMapLoad           InstructionKind = 24
	InstructionKindNewArray          InstructionKind = 25
	InstructionKindArrayStore        InstructionKind = 26
	InstructionKindArrayLoad         InstructionKind = 27
	InstructionKindNewError          InstructionKind = 28
	InstructionKindTypeCast          InstructionKind = 29
	InstructionKindIsLike            InstructionKind = 30
	InstructionKindTypeTest          InstructionKind = 31
	InstructionKindNewInstance       InstructionKind = 32
	InstructionKindObjectStore       InstructionKind = 33
	InstructionKindObjectLoad        InstructionKind = 34
	InstructionKindPanic             InstructionKind = 35
	InstructionKindFPLoad            InstructionKind = 36
	InstructionKindStringLoad        InstructionKind = 37
	InstructionKindNewXMLElement     InstructionKind = 38
	InstructionKindNewXMLText        InstructionKind = 39
	InstructionKindNewXMLComment     InstructionKind = 40
	InstructionKindNewXMLPI          InstructionKind = 41
	InstructionKindNewXMLSequence    InstructionKind = 42
	InstructionKindNewXMLQName       InstructionKind = 43
	InstructionKindNewStringXMLQName InstructionKind = 44
	InstructionKindXMLSeqStore       InstructionKind = 45
	InstructionKindXMLSeqLoad        InstructionKind = 46
	InstructionKindXMLLoad           InstructionKind = 47
	InstructionKindXMLLoadAll        InstructionKind = 48
	InstructionKindXMLAttributeLoad  InstructionKind = 49
	InstructionKindXMLAttributeStore InstructionKind = 50
	InstructionKindNewTable          InstructionKind = 51
	InstructionKindNewTypedesc       InstructionKind = 52
	InstructionKindNewStream         InstructionKind = 53
	InstructionKindTableStore        InstructionKind = 54
	InstructionKindTableLoad         InstructionKind = 55

	InstructionKindAdd           InstructionKind = 61
	InstructionKindSub           InstructionKind = 62
	InstructionKindMul           InstructionKind = 63
	InstructionKindDiv           InstructionKind = 64
	InstructionKindMod           InstructionKind = 65
	InstructionKindEqual         InstructionKind = 66
	InstructionKindNotEqual      InstructionKind = 67
	InstructionKindGreaterThan   InstructionKind = 68
	InstructionKindGreaterEqual  InstructionKind = 69
	InstructionKindLessThan      InstructionKind = 70
	InstructionKindLessEqual     InstructionKind = 71
	InstructionKindAnd           InstructionKind = 72
	InstructionKindOr            InstructionKind = 73
	InstructionKindRefEqual      InstructionKind = 74
	InstructionKindRefNotEqual   InstructionKind = 75
	InstructionKindClosedRange   InstructionKind = 76
	InstructionKindHalfOpenRange InstructionKind = 77
	InstructionKindAnnotAccess   InstructionKind = 78

	InstructionKindTypeof                    InstructionKind = 80
	InstructionKindNot                       InstructionKind = 81
	InstructionKindNegate                    InstructionKind = 82
	InstructionKindBitwiseAnd                InstructionKind = 83
	InstructionKindBitwiseOr                 InstructionKind = 84
	InstructionKindBitwiseXor                InstructionKind = 85
	InstructionKindBitwiseLeftShift          InstructionKind = 86
	InstructionKindBitwiseRightShift         InstructionKind = 87
	InstructionKindBitwiseUnsignedRightShift InstructionKind = 88

	InstructionKindNewRegExp              InstructionKind = 89
	InstructionKindNewREDisjunction       InstructionKind = 90
	InstructionKindNewRESequence          InstructionKind = 91
	InstructionKindNewREAssertion         InstructionKind = 92
	InstructionKindNewREAtomQuantifier    InstructionKind = 93
	InstructionKindNewRELiteralCharEscape InstructionKind = 94
	InstructionKindNewRECharClass         InstructionKind = 95
	InstructionKindNewRECharSet           InstructionKind = 96
	InstructionKindNewRECharSetRange      InstructionKind = 97
	InstructionKindNewRECapturingGroup    InstructionKind = 98
	InstructionKindNewREFlagExpr          InstructionKind = 99
	InstructionKindNewREFlagOnOff         InstructionKind = 100
	InstructionKindNewREQuantifier        InstructionKind = 101
	InstructionKindRecordDefaultFPLoad    InstructionKind = 102
	InstructionKindPlatform               InstructionKind = 128
)

func (i InstructionKind) Value() byte {
	return byte(i)
}

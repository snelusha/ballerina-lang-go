package tree

// SyntaxKind defines various kinds of syntax tree nodes, tokens and minutiae.
type SyntaxKind uint16

const (
	// Special values (in the order from Java enum end)
	None       SyntaxKind = 0
	List       SyntaxKind = 1
	EofToken   SyntaxKind = 2
	ModulePart SyntaxKind = 3
	Invalid    SyntaxKind = 4

	// Keywords (50-399)
	PublicKeyword     SyntaxKind = 50
	PrivateKeyword    SyntaxKind = 51
	RemoteKeyword     SyntaxKind = 52
	AbstractKeyword   SyntaxKind = 53
	ClientKeyword     SyntaxKind = 54
	ImportKeyword     SyntaxKind = 100
	FunctionKeyword   SyntaxKind = 101
	ConstKeyword      SyntaxKind = 102
	ListenerKeyword   SyntaxKind = 103
	ServiceKeyword    SyntaxKind = 104
	XmlnsKeyword      SyntaxKind = 105
	AnnotationKeyword SyntaxKind = 106
	TypeKeyword       SyntaxKind = 107
	RecordKeyword     SyntaxKind = 108
	ObjectKeyword     SyntaxKind = 109
	AsKeyword         SyntaxKind = 111
	OnKeyword         SyntaxKind = 112
	ResourceKeyword   SyntaxKind = 113
	FinalKeyword      SyntaxKind = 114
	SourceKeyword     SyntaxKind = 115
	WorkerKeyword     SyntaxKind = 117
	ParameterKeyword  SyntaxKind = 118
	FieldKeyword      SyntaxKind = 119
	IsolatedKeyword   SyntaxKind = 120

	ReturnsKeyword       SyntaxKind = 200
	ReturnKeyword        SyntaxKind = 201
	ExternalKeyword      SyntaxKind = 202
	TrueKeyword          SyntaxKind = 203
	FalseKeyword         SyntaxKind = 204
	IfKeyword            SyntaxKind = 205
	ElseKeyword          SyntaxKind = 206
	WhileKeyword         SyntaxKind = 207
	CheckKeyword         SyntaxKind = 208
	CheckpanicKeyword    SyntaxKind = 209
	PanicKeyword         SyntaxKind = 210
	ContinueKeyword      SyntaxKind = 211
	BreakKeyword         SyntaxKind = 212
	TypeofKeyword        SyntaxKind = 213
	IsKeyword            SyntaxKind = 214
	NullKeyword          SyntaxKind = 215
	LockKeyword          SyntaxKind = 216
	ForkKeyword          SyntaxKind = 217
	TrapKeyword          SyntaxKind = 218
	InKeyword            SyntaxKind = 219
	ForeachKeyword       SyntaxKind = 220
	TableKeyword         SyntaxKind = 221
	KeyKeyword           SyntaxKind = 222
	LetKeyword           SyntaxKind = 223
	NewKeyword           SyntaxKind = 224
	FromKeyword          SyntaxKind = 225
	WhereKeyword         SyntaxKind = 226
	SelectKeyword        SyntaxKind = 227
	StartKeyword         SyntaxKind = 228
	FlushKeyword         SyntaxKind = 229
	ConfigurableKeyword  SyntaxKind = 230
	WaitKeyword          SyntaxKind = 231
	DoKeyword            SyntaxKind = 232
	TransactionKeyword   SyntaxKind = 233
	TransactionalKeyword SyntaxKind = 234
	CommitKeyword        SyntaxKind = 235
	RollbackKeyword      SyntaxKind = 236
	RetryKeyword         SyntaxKind = 237
	EnumKeyword          SyntaxKind = 238
	Base16Keyword        SyntaxKind = 239
	Base64Keyword        SyntaxKind = 240
	MatchKeyword         SyntaxKind = 241
	ConflictKeyword      SyntaxKind = 242
	LimitKeyword         SyntaxKind = 243
	JoinKeyword          SyntaxKind = 244
	OuterKeyword         SyntaxKind = 245
	EqualsKeyword        SyntaxKind = 246
	ClassKeyword         SyntaxKind = 247
	OrderKeyword         SyntaxKind = 248
	ByKeyword            SyntaxKind = 249
	AscendingKeyword     SyntaxKind = 250
	DescendingKeyword    SyntaxKind = 251
	UnderscoreKeyword    SyntaxKind = 252
	NotIsKeyword         SyntaxKind = 253
	NaturalKeyword       SyntaxKind = 254

	// Type keywords (300-319)
	IntKeyword      SyntaxKind = 300
	ByteKeyword     SyntaxKind = 301
	FloatKeyword    SyntaxKind = 302
	DecimalKeyword  SyntaxKind = 303
	StringKeyword   SyntaxKind = 304
	BooleanKeyword  SyntaxKind = 305
	XmlKeyword      SyntaxKind = 306
	JsonKeyword     SyntaxKind = 307
	HandleKeyword   SyntaxKind = 308
	AnyKeyword      SyntaxKind = 309
	AnydataKeyword  SyntaxKind = 310
	NeverKeyword    SyntaxKind = 311
	VarKeyword      SyntaxKind = 312
	MapKeyword      SyntaxKind = 313
	FutureKeyword   SyntaxKind = 314
	TypedescKeyword SyntaxKind = 315
	ErrorKeyword    SyntaxKind = 316
	StreamKeyword   SyntaxKind = 317
	ReadonlyKeyword SyntaxKind = 318
	DistinctKeyword SyntaxKind = 319
	FailKeyword     SyntaxKind = 320

	// Contextual keywords (400-402)
	ReKeyword      SyntaxKind = 400
	GroupKeyword   SyntaxKind = 401
	CollectKeyword SyntaxKind = 402

	// Separators (500-519)
	OpenBraceToken      SyntaxKind = 500
	CloseBraceToken     SyntaxKind = 501
	OpenParenToken      SyntaxKind = 502
	CloseParenToken     SyntaxKind = 503
	OpenBracketToken    SyntaxKind = 504
	CloseBracketToken   SyntaxKind = 505
	SemicolonToken      SyntaxKind = 506
	DotToken            SyntaxKind = 507
	ColonToken          SyntaxKind = 508
	CommaToken          SyntaxKind = 509
	EllipsisToken       SyntaxKind = 510
	OpenBracePipeToken  SyntaxKind = 511
	CloseBracePipeToken SyntaxKind = 512
	AtToken             SyntaxKind = 513
	HashToken           SyntaxKind = 514
	BacktickToken       SyntaxKind = 515
	DoubleQuoteToken    SyntaxKind = 516
	SingleQuoteToken    SyntaxKind = 517
	DoubleBacktickToken SyntaxKind = 518
	TripleBacktickToken SyntaxKind = 519

	// Operators (550-595)
	EqualToken                       SyntaxKind = 550
	DoubleEqualToken                 SyntaxKind = 551
	TripleEqualToken                 SyntaxKind = 552
	PlusToken                        SyntaxKind = 553
	MinusToken                       SyntaxKind = 554
	SlashToken                       SyntaxKind = 555
	PercentToken                     SyntaxKind = 556
	AsteriskToken                    SyntaxKind = 557
	LtToken                          SyntaxKind = 558
	LtEqualToken                     SyntaxKind = 559
	GtToken                          SyntaxKind = 560
	RightDoubleArrowToken            SyntaxKind = 561
	QuestionMarkToken                SyntaxKind = 562
	PipeToken                        SyntaxKind = 563
	GtEqualToken                     SyntaxKind = 564
	ExclamationMarkToken             SyntaxKind = 565
	NotEqualToken                    SyntaxKind = 566
	NotDoubleEqualToken              SyntaxKind = 567
	BitwiseAndToken                  SyntaxKind = 568
	BitwiseXorToken                  SyntaxKind = 569
	LogicalAndToken                  SyntaxKind = 570
	LogicalOrToken                   SyntaxKind = 571
	NegationToken                    SyntaxKind = 572
	RightArrowToken                  SyntaxKind = 573
	InterpolationStartToken          SyntaxKind = 574
	XmlPiStartToken                  SyntaxKind = 575
	XmlPiEndToken                    SyntaxKind = 576
	XmlCommentStartToken             SyntaxKind = 577
	XmlCommentEndToken               SyntaxKind = 578
	SyncSendToken                    SyntaxKind = 579
	LeftArrowToken                   SyntaxKind = 580
	DoubleDotLtToken                 SyntaxKind = 580 // Note: Same value as LeftArrowToken in Java
	DoubleLtToken                    SyntaxKind = 581
	AnnotChainingToken               SyntaxKind = 582
	OptionalChainingToken            SyntaxKind = 583
	ElvisToken                       SyntaxKind = 584
	DotLtToken                       SyntaxKind = 585
	SlashLtToken                     SyntaxKind = 586
	DoubleSlashDoubleAsteriskLtToken SyntaxKind = 587
	SlashAsteriskToken               SyntaxKind = 588
	DoubleGtToken                    SyntaxKind = 589
	TripleGtToken                    SyntaxKind = 590
	XmlCdataStartToken               SyntaxKind = 591
	XmlCdataEndToken                 SyntaxKind = 592
	BackSlashToken                   SyntaxKind = 593
	DollarToken                      SyntaxKind = 594
	EscapedMinusToken                SyntaxKind = 595

	// Documentation reference types (900-908)
	TypeDocReferenceToken       SyntaxKind = 900
	ServiceDocReferenceToken    SyntaxKind = 901
	VariableDocReferenceToken   SyntaxKind = 902
	VarDocReferenceToken        SyntaxKind = 903
	AnnotationDocReferenceToken SyntaxKind = 904
	ModuleDocReferenceToken     SyntaxKind = 905
	FunctionDocReferenceToken   SyntaxKind = 906
	ParameterDocReferenceToken  SyntaxKind = 907
	ConstDocReferenceToken      SyntaxKind = 908

	// Literal tokens (1000-1007)
	IdentifierToken                  SyntaxKind = 1000
	StringLiteralToken               SyntaxKind = 1001
	DecimalIntegerLiteralToken       SyntaxKind = 1002
	HexIntegerLiteralToken           SyntaxKind = 1003
	DecimalFloatingPointLiteralToken SyntaxKind = 1004
	HexFloatingPointLiteralToken     SyntaxKind = 1005
	XmlTextContent                   SyntaxKind = 1006
	TemplateString                   SyntaxKind = 1007
	PromptContent                    SyntaxKind = 1007 // Note: Same value as TemplateString in Java

	// Documentation (1100-1104)
	DocumentationDescription SyntaxKind = 1100
	ParameterName            SyntaxKind = 1101
	CodeContent              SyntaxKind = 1102
	DeprecationLiteral       SyntaxKind = 1103
	DocumentationString      SyntaxKind = 1104

	// Other
	InvalidToken SyntaxKind = 1191

	// Statements (1200-1225)
	BlockStatement               SyntaxKind = 1200
	LocalVarDecl                 SyntaxKind = 1201
	AssignmentStatement          SyntaxKind = 1202
	IfElseStatement              SyntaxKind = 1203
	ElseBlock                    SyntaxKind = 1204
	WhileStatement               SyntaxKind = 1205
	CallStatement                SyntaxKind = 1206
	PanicStatement               SyntaxKind = 1207
	ReturnStatement              SyntaxKind = 1208
	ContinueStatement            SyntaxKind = 1209
	BreakStatement               SyntaxKind = 1210
	CompoundAssignmentStatement  SyntaxKind = 1211
	LocalTypeDefinitionStatement SyntaxKind = 1212
	ActionStatement              SyntaxKind = 1213
	LockStatement                SyntaxKind = 1214
	NamedWorkerDeclaration       SyntaxKind = 1215
	ForkStatement                SyntaxKind = 1216
	ForeachStatement             SyntaxKind = 1217
	TransactionStatement         SyntaxKind = 1218
	RollbackStatement            SyntaxKind = 1219
	RetryStatement               SyntaxKind = 1220
	XmlNamespaceDeclaration      SyntaxKind = 1221
	MatchStatement               SyntaxKind = 1222
	InvalidExpressionStatement   SyntaxKind = 1223
	DoStatement                  SyntaxKind = 1224
	FailStatement                SyntaxKind = 1225

	// Expressions (1300-1348)
	BinaryExpression                    SyntaxKind = 1300
	BracedExpression                    SyntaxKind = 1301
	FunctionCall                        SyntaxKind = 1302
	QualifiedNameReference              SyntaxKind = 1303
	IndexedExpression                   SyntaxKind = 1304
	FieldAccess                         SyntaxKind = 1305
	MethodCall                          SyntaxKind = 1306
	CheckExpression                     SyntaxKind = 1307
	MappingConstructor                  SyntaxKind = 1308
	TypeofExpression                    SyntaxKind = 1309
	UnaryExpression                     SyntaxKind = 1310
	TypeTestExpression                  SyntaxKind = 1311
	SimpleNameReference                 SyntaxKind = 1313
	TrapExpression                      SyntaxKind = 1314
	ListConstructor                     SyntaxKind = 1315
	TypeCastExpression                  SyntaxKind = 1316
	TableConstructor                    SyntaxKind = 1317
	LetExpression                       SyntaxKind = 1318
	XmlTemplateExpression               SyntaxKind = 1319
	RawTemplateExpression               SyntaxKind = 1320
	StringTemplateExpression            SyntaxKind = 1321
	ImplicitNewExpression               SyntaxKind = 1322
	ExplicitNewExpression               SyntaxKind = 1323
	ParenthesizedArgList                SyntaxKind = 1324
	ExplicitAnonymousFunctionExpression SyntaxKind = 1325
	ImplicitAnonymousFunctionExpression SyntaxKind = 1326
	QueryExpression                     SyntaxKind = 1327
	AnnotAccess                         SyntaxKind = 1328
	OptionalFieldAccess                 SyntaxKind = 1329
	ConditionalExpression               SyntaxKind = 1330
	TransactionalExpression             SyntaxKind = 1331
	ObjectConstructor                   SyntaxKind = 1332
	XmlFilterExpression                 SyntaxKind = 1333
	XmlStepExpression                   SyntaxKind = 1334
	XmlNamePatternChain                 SyntaxKind = 1335
	XmlAtomicNamePattern                SyntaxKind = 1336
	StringLiteral                       SyntaxKind = 1337
	NumericLiteral                      SyntaxKind = 1338
	BooleanLiteral                      SyntaxKind = 1339
	NilLiteral                          SyntaxKind = 1340
	NullLiteral                         SyntaxKind = 1341
	ByteArrayLiteral                    SyntaxKind = 1342
	AsteriskLiteral                     SyntaxKind = 1343
	RequiredExpression                  SyntaxKind = 1344
	ErrorConstructor                    SyntaxKind = 1345
	RegexTemplateExpression             SyntaxKind = 1346 // Note: Out of sequence in Java
	XmlStepMethodCallExtend             SyntaxKind = 1346 // Note: Same value as RegexTemplateExpression in Java
	XmlStepIndexedExtend                SyntaxKind = 1347
	NaturalExpression                   SyntaxKind = 1348

	// Minutiae kinds (1500-1503)
	WhitespaceMinutiae  SyntaxKind = 1500
	EndOfLineMinutiae   SyntaxKind = 1501
	CommentMinutiae     SyntaxKind = 1502
	InvalidNodeMinutiae SyntaxKind = 1503

	// Invalid nodes (1601)
	InvalidTokenMinutiaeNode SyntaxKind = 1601

	// Module-level declarations (2000-2010)
	ImportDeclaration             SyntaxKind = 2000
	FunctionDefinition            SyntaxKind = 2001
	TypeDefinition                SyntaxKind = 2002
	ServiceDeclaration            SyntaxKind = 2003
	ModuleVarDecl                 SyntaxKind = 2004
	ListenerDeclaration           SyntaxKind = 2005
	ConstDeclaration              SyntaxKind = 2006
	AnnotationDeclaration         SyntaxKind = 2007
	ModuleXmlNamespaceDeclaration SyntaxKind = 2008
	EnumDeclaration               SyntaxKind = 2009
	ClassDefinition               SyntaxKind = 2010

	// Type descriptors (2000-2034) - Note: OVERLAPS with module declarations in original Java!
	TypeDesc              SyntaxKind = 2000 // Same as ImportDeclaration
	RecordTypeDesc        SyntaxKind = 2001 // Same as FunctionDefinition
	ObjectTypeDesc        SyntaxKind = 2002 // Same as TypeDefinition
	NilTypeDesc           SyntaxKind = 2003 // Same as ServiceDeclaration
	OptionalTypeDesc      SyntaxKind = 2004 // Same as ModuleVarDecl
	ArrayTypeDesc         SyntaxKind = 2005 // Same as ListenerDeclaration
	IntTypeDesc           SyntaxKind = 2006 // Same as ConstDeclaration
	ByteTypeDesc          SyntaxKind = 2007 // Same as AnnotationDeclaration
	FloatTypeDesc         SyntaxKind = 2008 // Same as ModuleXmlNamespaceDeclaration
	DecimalTypeDesc       SyntaxKind = 2009 // Same as EnumDeclaration
	StringTypeDesc        SyntaxKind = 2010 // Same as ClassDefinition
	BooleanTypeDesc       SyntaxKind = 2011
	XmlTypeDesc           SyntaxKind = 2012
	JsonTypeDesc          SyntaxKind = 2013
	HandleTypeDesc        SyntaxKind = 2014
	AnyTypeDesc           SyntaxKind = 2015
	AnydataTypeDesc       SyntaxKind = 2016
	NeverTypeDesc         SyntaxKind = 2017
	VarTypeDesc           SyntaxKind = 2018
	ServiceTypeDesc       SyntaxKind = 2019
	MapTypeDesc           SyntaxKind = 2020
	UnionTypeDesc         SyntaxKind = 2021
	ErrorTypeDesc         SyntaxKind = 2022
	StreamTypeDesc        SyntaxKind = 2023
	TableTypeDesc         SyntaxKind = 2024
	FunctionTypeDesc      SyntaxKind = 2025
	TupleTypeDesc         SyntaxKind = 2026
	ParenthesisedTypeDesc SyntaxKind = 2027
	ReadonlyTypeDesc      SyntaxKind = 2028
	DistinctTypeDesc      SyntaxKind = 2029
	IntersectionTypeDesc  SyntaxKind = 2030
	SingletonTypeDesc     SyntaxKind = 2031
	TypeReferenceTypeDesc SyntaxKind = 2032
	TypedescTypeDesc      SyntaxKind = 2033
	FutureTypeDesc        SyntaxKind = 2034

	// Actions (2500-2512)
	RemoteMethodCallAction     SyntaxKind = 2500
	BracedAction               SyntaxKind = 2501
	CheckAction                SyntaxKind = 2502
	StartAction                SyntaxKind = 2503
	TrapAction                 SyntaxKind = 2504
	FlushAction                SyntaxKind = 2505
	AsyncSendAction            SyntaxKind = 2506
	SyncSendAction             SyntaxKind = 2507
	ReceiveAction              SyntaxKind = 2508
	WaitAction                 SyntaxKind = 2509
	QueryAction                SyntaxKind = 2510
	CommitAction               SyntaxKind = 2511
	ClientResourceAccessAction SyntaxKind = 2512

	// Other syntax elements (3000-3096)
	ReturnTypeDescriptor          SyntaxKind = 3000
	RequiredParam                 SyntaxKind = 3001
	DefaultableParam              SyntaxKind = 3002
	RestParam                     SyntaxKind = 3003
	ExternalFunctionBody          SyntaxKind = 3004
	RecordField                   SyntaxKind = 3005
	RecordFieldWithDefaultValue   SyntaxKind = 3006
	TypeReference                 SyntaxKind = 3007
	RecordRestType                SyntaxKind = 3008
	PositionalArg                 SyntaxKind = 3009
	NamedArg                      SyntaxKind = 3010
	RestArg                       SyntaxKind = 3011
	ObjectField                   SyntaxKind = 3012
	ImportOrgName                 SyntaxKind = 3013
	ModuleName                    SyntaxKind = 3014
	SubModuleName                 SyntaxKind = 3015
	ImportVersion                 SyntaxKind = 3016
	OrderByClause                 SyntaxKind = 3017
	ImportPrefix                  SyntaxKind = 3018
	SpecificField                 SyntaxKind = 3019
	ComputedNameField             SyntaxKind = 3020
	SpreadField                   SyntaxKind = 3021
	OrderKey                      SyntaxKind = 3022
	ResourceAccessorDefinition    SyntaxKind = 3023
	Annotation                    SyntaxKind = 3024
	Metadata                      SyntaxKind = 3025
	ArrayDimension                SyntaxKind = 3026
	AnnotationAttachPoint         SyntaxKind = 3028
	FunctionBodyBlock             SyntaxKind = 3029
	NamedWorkerDeclarator         SyntaxKind = 3030
	ExpressionFunctionBody        SyntaxKind = 3031
	TypeCastParam                 SyntaxKind = 3032
	KeySpecifier                  SyntaxKind = 3033
	ExplicitTypeParams            SyntaxKind = 3034
	LetVarDecl                    SyntaxKind = 3035
	StreamTypeParams              SyntaxKind = 3036
	FunctionSignature             SyntaxKind = 3037
	InferParamList                SyntaxKind = 3038
	TypeParameter                 SyntaxKind = 3039
	KeyTypeConstraint             SyntaxKind = 3040
	QueryConstructType            SyntaxKind = 3041
	FromClause                    SyntaxKind = 3042
	WhereClause                   SyntaxKind = 3043
	LetClause                     SyntaxKind = 3044
	QueryPipeline                 SyntaxKind = 3045
	SelectClause                  SyntaxKind = 3046
	MethodDeclaration             SyntaxKind = 3047
	TypedBindingPattern           SyntaxKind = 3048
	BindingPattern                SyntaxKind = 3049
	CaptureBindingPattern         SyntaxKind = 3050
	RestBindingPattern            SyntaxKind = 3051
	ListBindingPattern            SyntaxKind = 3052
	ReceiveFields                 SyntaxKind = 3053
	RestType                      SyntaxKind = 3054
	WaitFieldsList                SyntaxKind = 3055
	WaitField                     SyntaxKind = 3056
	EnumMember                    SyntaxKind = 3057
	BracketedList                 SyntaxKind = 3058
	ListBpOrListConstructor       SyntaxKind = 3059
	MappingBindingPattern         SyntaxKind = 3060
	FieldBindingPattern           SyntaxKind = 3061
	MappingBpOrMappingConstructor SyntaxKind = 3062
	WildcardBindingPattern        SyntaxKind = 3063
	MatchClause                   SyntaxKind = 3064
	MatchGuard                    SyntaxKind = 3065
	ObjectMethodDefinition        SyntaxKind = 3066
	OnConflictClause              SyntaxKind = 3067
	LimitClause                   SyntaxKind = 3068
	JoinClause                    SyntaxKind = 3069
	OnClause                      SyntaxKind = 3070
	ListMatchPattern              SyntaxKind = 3071
	RestMatchPattern              SyntaxKind = 3072
	MappingMatchPattern           SyntaxKind = 3073
	FieldMatchPattern             SyntaxKind = 3074
	ErrorMatchPattern             SyntaxKind = 3075
	NamedArgMatchPattern          SyntaxKind = 3076
	ErrorBindingPattern           SyntaxKind = 3077
	NamedArgBindingPattern        SyntaxKind = 3078
	TupleTypeDescOrListConst      SyntaxKind = 3079
	OnFailClause                  SyntaxKind = 3080
	ResourceAccessorDeclaration   SyntaxKind = 3081
	ResourcePathSegmentParam      SyntaxKind = 3082
	ResourcePathRestParam         SyntaxKind = 3083
	IncludedRecordParam           SyntaxKind = 3084
	ArrayTypeDescOrMemberAccess   SyntaxKind = 3085
	InferredTypedescDefault       SyntaxKind = 3086
	SpreadMember                  SyntaxKind = 3087
	ComputedResourceAccessSegment SyntaxKind = 3088
	ResourceAccessRestSegment     SyntaxKind = 3089
	MemberTypeDesc                SyntaxKind = 3090
	GroupingKeyVarDeclaration     SyntaxKind = 3091
	GroupingKeyVarName            SyntaxKind = 3092
	GroupByClause                 SyntaxKind = 3093
	CollectClause                 SyntaxKind = 3094
	AlternateReceive              SyntaxKind = 3095
	ReceiveField                  SyntaxKind = 3096

	// XML elements (4000-4012)
	XmlElement         SyntaxKind = 4000
	XmlEmptyElement    SyntaxKind = 4001
	XmlText            SyntaxKind = 4002
	XmlComment         SyntaxKind = 4003
	XmlPi              SyntaxKind = 4004
	XmlElementStartTag SyntaxKind = 4005
	XmlElementEndTag   SyntaxKind = 4006
	XmlSimpleName      SyntaxKind = 4007
	XmlQualifiedName   SyntaxKind = 4008
	XmlAttribute       SyntaxKind = 4009
	XmlAttributeValue  SyntaxKind = 4010
	Interpolation      SyntaxKind = 4011
	XmlCdata           SyntaxKind = 4012

	// Regular expressions (4013-4048)
	ReSequence                             SyntaxKind = 4013
	ReAtomQuantifier                       SyntaxKind = 4014
	ReAssertion                            SyntaxKind = 4015
	ReLiteralCharDotOrEscape               SyntaxKind = 4016
	ReQuoteEscape                          SyntaxKind = 4017
	ReSimpleCharClassEscape                SyntaxKind = 4018
	ReUnicodePropertyEscape                SyntaxKind = 4019
	ReUnicodeScript                        SyntaxKind = 4020
	ReUnicodeGeneralCategory               SyntaxKind = 4021
	ReCharacterClass                       SyntaxKind = 4022
	ReCharSetAtomWithReCharSetNoDash       SyntaxKind = 4023
	ReCharSetAtomNoDashWithReCharSetNoDash SyntaxKind = 4024
	ReCharSetRange                         SyntaxKind = 4025
	ReCharSetRangeNoDash                   SyntaxKind = 4026
	ReCharSetRangeWithReCharSet            SyntaxKind = 4027
	ReCharSetRangeNoDashWithReCharSet      SyntaxKind = 4028
	ReCapturingGroup                       SyntaxKind = 4029
	ReFlagExpr                             SyntaxKind = 4030
	ReFlagsOnOff                           SyntaxKind = 4031
	ReFlags                                SyntaxKind = 4032
	ReQuantifier                           SyntaxKind = 4033
	ReBracedQuantifier                     SyntaxKind = 4034
	ReAssertionValue                       SyntaxKind = 4035
	ReLiteralChar                          SyntaxKind = 4036
	ReNumericEscape                        SyntaxKind = 4037
	ReControlEscape                        SyntaxKind = 4038
	ReSimpleCharClassCode                  SyntaxKind = 4039
	ReProperty                             SyntaxKind = 4040
	ReUnicodeScriptStart                   SyntaxKind = 4041
	ReUnicodePropertyValue                 SyntaxKind = 4042
	ReUnicodeGeneralCategoryStart          SyntaxKind = 4043
	ReUnicodeGeneralCategoryName           SyntaxKind = 4044
	ReCharSetAtomNoDash                    SyntaxKind = 4045
	ReFlagsValue                           SyntaxKind = 4046
	ReBaseQuantifierValue                  SyntaxKind = 4047
	Digit                                  SyntaxKind = 4048

	// Documentation (4500-4509)
	MarkdownDocumentation                    SyntaxKind = 4500
	MarkdownDocumentationLine                SyntaxKind = 4501
	MarkdownReferenceDocumentationLine       SyntaxKind = 4502
	MarkdownParameterDocumentationLine       SyntaxKind = 4503
	MarkdownReturnParameterDocumentationLine SyntaxKind = 4504
	MarkdownDeprecationDocumentationLine     SyntaxKind = 4505
	MarkdownCodeLine                         SyntaxKind = 4506
	BallerinaNameReference                   SyntaxKind = 4507
	MarkdownCodeBlock                        SyntaxKind = 4508
	InlineCodeReference                      SyntaxKind = 4509
)

// String returns the string representation of the syntax kind.
func (sk SyntaxKind) String() string {
	switch sk {
	case None:
		return "NONE"
	case List:
		return "LIST"
	case EofToken:
		return "EOF_TOKEN"
	case ModulePart:
		return "MODULE_PART"
	case Invalid:
		return "INVALID"

	// Keywords
	case PublicKeyword:
		return "PUBLIC_KEYWORD"
	case PrivateKeyword:
		return "PRIVATE_KEYWORD"
	case RemoteKeyword:
		return "REMOTE_KEYWORD"
	case AbstractKeyword:
		return "ABSTRACT_KEYWORD"
	case ClientKeyword:
		return "CLIENT_KEYWORD"
	case ImportKeyword:
		return "IMPORT_KEYWORD"
	case FunctionKeyword:
		return "FUNCTION_KEYWORD"
	case ConstKeyword:
		return "CONST_KEYWORD"
	case ListenerKeyword:
		return "LISTENER_KEYWORD"
	case ServiceKeyword:
		return "SERVICE_KEYWORD"
	case XmlnsKeyword:
		return "XMLNS_KEYWORD"
	case AnnotationKeyword:
		return "ANNOTATION_KEYWORD"
	case TypeKeyword:
		return "TYPE_KEYWORD"
	case RecordKeyword:
		return "RECORD_KEYWORD"
	case ObjectKeyword:
		return "OBJECT_KEYWORD"
	case AsKeyword:
		return "AS_KEYWORD"
	case OnKeyword:
		return "ON_KEYWORD"
	case ResourceKeyword:
		return "RESOURCE_KEYWORD"
	case FinalKeyword:
		return "FINAL_KEYWORD"
	case SourceKeyword:
		return "SOURCE_KEYWORD"
	case WorkerKeyword:
		return "WORKER_KEYWORD"
	case ParameterKeyword:
		return "PARAMETER_KEYWORD"
	case FieldKeyword:
		return "FIELD_KEYWORD"
	case IsolatedKeyword:
		return "ISOLATED_KEYWORD"

	// Add more cases as needed...
	default:
		return "UNKNOWN"
	}
}

// StringValue returns the string value associated with the syntax kind.
func (sk SyntaxKind) StringValue() string {
	switch sk {
	// Keywords with string values
	case PublicKeyword:
		return "public"
	case PrivateKeyword:
		return "private"
	case RemoteKeyword:
		return "remote"
	case AbstractKeyword:
		return "abstract"
	case ClientKeyword:
		return "client"
	case ImportKeyword:
		return "import"
	case FunctionKeyword:
		return "function"
	case ConstKeyword:
		return "const"
	case ListenerKeyword:
		return "listener"
	case ServiceKeyword:
		return "service"
	case XmlnsKeyword:
		return "xmlns"
	case AnnotationKeyword:
		return "annotation"
	case TypeKeyword:
		return "type"
	case RecordKeyword:
		return "record"
	case ObjectKeyword:
		return "object"
	case AsKeyword:
		return "as"
	case OnKeyword:
		return "on"
	case ResourceKeyword:
		return "resource"
	case FinalKeyword:
		return "final"
	case SourceKeyword:
		return "source"
	case WorkerKeyword:
		return "worker"
	case ParameterKeyword:
		return "parameter"
	case FieldKeyword:
		return "field"
	case IsolatedKeyword:
		return "isolated"

	// Separators
	case OpenBraceToken:
		return "{"
	case CloseBraceToken:
		return "}"
	case OpenParenToken:
		return "("
	case CloseParenToken:
		return ")"
	case OpenBracketToken:
		return "["
	case CloseBracketToken:
		return "]"
	case SemicolonToken:
		return ";"
	case DotToken:
		return "."
	case ColonToken:
		return ":"
	case CommaToken:
		return ","
	case EllipsisToken:
		return "..."

	// Add more cases as needed...
	default:
		return ""
	}
}

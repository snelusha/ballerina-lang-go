package tree

type SyntaxKind struct {
	tag      uint16
	strValue string
}

func (sk SyntaxKind) StringValue() string {
	return sk.strValue
}

var (
	PUBLIC_KEYWORD     = SyntaxKind{tag: 50, strValue: "public"}
	PRIVATE_KEYWORD    = SyntaxKind{tag: 51, strValue: "private"}
	REMOTE_KEYWORD     = SyntaxKind{tag: 52, strValue: "remote"}
	ABSTRACT_KEYWORD   = SyntaxKind{tag: 53, strValue: "abstract"}
	CLIENT_KEYWORD     = SyntaxKind{tag: 54, strValue: "client"}
	IMPORT_KEYWORD     = SyntaxKind{tag: 100, strValue: "import"}
	FUNCTION_KEYWORD   = SyntaxKind{tag: 101, strValue: "function"}
	CONST_KEYWORD      = SyntaxKind{tag: 102, strValue: "const"}
	LISTENER_KEYWORD   = SyntaxKind{tag: 103, strValue: "listener"}
	SERVICE_KEYWORD    = SyntaxKind{tag: 104, strValue: "service"}
	XMLNS_KEYWORD      = SyntaxKind{tag: 105, strValue: "xmlns"}
	ANNOTATION_KEYWORD = SyntaxKind{tag: 106, strValue: "annotation"}
	TYPE_KEYWORD       = SyntaxKind{tag: 107, strValue: "type"}
	RECORD_KEYWORD     = SyntaxKind{tag: 108, strValue: "record"}
	OBJECT_KEYWORD     = SyntaxKind{tag: 109, strValue: "object"}
	AS_KEYWORD         = SyntaxKind{tag: 111, strValue: "as"}
	ON_KEYWORD         = SyntaxKind{tag: 112, strValue: "on"}
	RESOURCE_KEYWORD   = SyntaxKind{tag: 113, strValue: "resource"}
	FINAL_KEYWORD      = SyntaxKind{tag: 114, strValue: "final"}
	SOURCE_KEYWORD     = SyntaxKind{tag: 115, strValue: "source"}
	WORKER_KEYWORD     = SyntaxKind{tag: 117, strValue: "worker"}
	PARAMETER_KEYWORD  = SyntaxKind{tag: 118, strValue: "parameter"}
	FIELD_KEYWORD      = SyntaxKind{tag: 119, strValue: "field"}
	ISOLATED_KEYWORD   = SyntaxKind{tag: 120, strValue: "isolated"}

	RETURNS_KEYWORD       = SyntaxKind{tag: 200, strValue: "returns"}
	RETURN_KEYWORD        = SyntaxKind{tag: 201, strValue: "return"}
	EXTERNAL_KEYWORD      = SyntaxKind{tag: 202, strValue: "external"}
	TRUE_KEYWORD          = SyntaxKind{tag: 203, strValue: "true"}
	FALSE_KEYWORD         = SyntaxKind{tag: 204, strValue: "false"}
	IF_KEYWORD            = SyntaxKind{tag: 205, strValue: "if"}
	ELSE_KEYWORD          = SyntaxKind{tag: 206, strValue: "else"}
	WHILE_KEYWORD         = SyntaxKind{tag: 207, strValue: "while"}
	CHECK_KEYWORD         = SyntaxKind{tag: 208, strValue: "check"}
	CHECKPANIC_KEYWORD    = SyntaxKind{tag: 209, strValue: "checkpanic"}
	PANIC_KEYWORD         = SyntaxKind{tag: 210, strValue: "panic"}
	CONTINUE_KEYWORD      = SyntaxKind{tag: 211, strValue: "continue"}
	BREAK_KEYWORD         = SyntaxKind{tag: 212, strValue: "break"}
	TYPEOF_KEYWORD        = SyntaxKind{tag: 213, strValue: "typeof"}
	IS_KEYWORD            = SyntaxKind{tag: 214, strValue: "is"}
	NULL_KEYWORD          = SyntaxKind{tag: 215, strValue: "null"}
	LOCK_KEYWORD          = SyntaxKind{tag: 216, strValue: "lock"}
	FORK_KEYWORD          = SyntaxKind{tag: 217, strValue: "fork"}
	TRAP_KEYWORD          = SyntaxKind{tag: 218, strValue: "trap"}
	IN_KEYWORD            = SyntaxKind{tag: 219, strValue: "in"}
	FOREACH_KEYWORD       = SyntaxKind{tag: 220, strValue: "foreach"}
	TABLE_KEYWORD         = SyntaxKind{tag: 221, strValue: "table"}
	KEY_KEYWORD           = SyntaxKind{tag: 222, strValue: "key"}
	LET_KEYWORD           = SyntaxKind{tag: 223, strValue: "let"}
	NEW_KEYWORD           = SyntaxKind{tag: 224, strValue: "new"}
	FROM_KEYWORD          = SyntaxKind{tag: 225, strValue: "from"}
	WHERE_KEYWORD         = SyntaxKind{tag: 226, strValue: "where"}
	SELECT_KEYWORD        = SyntaxKind{tag: 227, strValue: "select"}
	START_KEYWORD         = SyntaxKind{tag: 228, strValue: "start"}
	FLUSH_KEYWORD         = SyntaxKind{tag: 229, strValue: "flush"}
	CONFIGURABLE_KEYWORD  = SyntaxKind{tag: 230, strValue: "configurable"}
	WAIT_KEYWORD          = SyntaxKind{tag: 231, strValue: "wait"}
	DO_KEYWORD            = SyntaxKind{tag: 232, strValue: "do"}
	TRANSACTION_KEYWORD   = SyntaxKind{tag: 233, strValue: "transaction"}
	TRANSACTIONAL_KEYWORD = SyntaxKind{tag: 234, strValue: "transactional"}
	COMMIT_KEYWORD        = SyntaxKind{tag: 235, strValue: "commit"}
	ROLLBACK_KEYWORD      = SyntaxKind{tag: 236, strValue: "rollback"}
	RETRY_KEYWORD         = SyntaxKind{tag: 237, strValue: "retry"}
	ENUM_KEYWORD          = SyntaxKind{tag: 238, strValue: "enum"}
	BASE16_KEYWORD        = SyntaxKind{tag: 239, strValue: "base16"}
	BASE64_KEYWORD        = SyntaxKind{tag: 240, strValue: "base64"}
	MATCH_KEYWORD         = SyntaxKind{tag: 241, strValue: "match"}
	CONFLICT_KEYWORD      = SyntaxKind{tag: 242, strValue: "conflict"}
	LIMIT_KEYWORD         = SyntaxKind{tag: 243, strValue: "limit"}
	JOIN_KEYWORD          = SyntaxKind{tag: 244, strValue: "join"}
	OUTER_KEYWORD         = SyntaxKind{tag: 245, strValue: "outer"}
	EQUALS_KEYWORD        = SyntaxKind{tag: 246, strValue: "equals"}
	CLASS_KEYWORD         = SyntaxKind{tag: 247, strValue: "class"}
	ORDER_KEYWORD         = SyntaxKind{tag: 248, strValue: "order"}
	BY_KEYWORD            = SyntaxKind{tag: 249, strValue: "by"}
	ASCENDING_KEYWORD     = SyntaxKind{tag: 250, strValue: "ascending"}
	DESCENDING_KEYWORD    = SyntaxKind{tag: 251, strValue: "descending"}
	UNDERSCORE_KEYWORD    = SyntaxKind{tag: 252, strValue: "_"}
	NOT_IS_KEYWORD        = SyntaxKind{tag: 253, strValue: "!is"}
	NATURAL_KEYWORD       = SyntaxKind{tag: 254, strValue: "natural"}

	// Type keywords
	INT_KEYWORD      = SyntaxKind{tag: 300, strValue: "int"}
	BYTE_KEYWORD     = SyntaxKind{tag: 301, strValue: "byte"}
	FLOAT_KEYWORD    = SyntaxKind{tag: 302, strValue: "float"}
	DECIMAL_KEYWORD  = SyntaxKind{tag: 303, strValue: "decimal"}
	STRING_KEYWORD   = SyntaxKind{tag: 304, strValue: "string"}
	BOOLEAN_KEYWORD  = SyntaxKind{tag: 305, strValue: "boolean"}
	XML_KEYWORD      = SyntaxKind{tag: 306, strValue: "xml"}
	JSON_KEYWORD     = SyntaxKind{tag: 307, strValue: "json"}
	HANDLE_KEYWORD   = SyntaxKind{tag: 308, strValue: "handle"}
	ANY_KEYWORD      = SyntaxKind{tag: 309, strValue: "any"}
	ANYDATA_KEYWORD  = SyntaxKind{tag: 310, strValue: "anydata"}
	NEVER_KEYWORD    = SyntaxKind{tag: 311, strValue: "never"}
	VAR_KEYWORD      = SyntaxKind{tag: 312, strValue: "var"}
	MAP_KEYWORD      = SyntaxKind{tag: 313, strValue: "map"}
	FUTURE_KEYWORD   = SyntaxKind{tag: 314, strValue: "future"}
	TYPEDESC_KEYWORD = SyntaxKind{tag: 315, strValue: "typedesc"}
	ERROR_KEYWORD    = SyntaxKind{tag: 316, strValue: "error"}
	STREAM_KEYWORD   = SyntaxKind{tag: 317, strValue: "stream"}
	READONLY_KEYWORD = SyntaxKind{tag: 318, strValue: "readonly"}
	DISTINCT_KEYWORD = SyntaxKind{tag: 319, strValue: "distinct"}
	FAIL_KEYWORD     = SyntaxKind{tag: 320, strValue: "fail"}

	// Contextual keywords
	RE_KEYWORD      = SyntaxKind{tag: 400, strValue: "re"} // Any kind above this is considered as a keyword
	GROUP_KEYWORD   = SyntaxKind{tag: 401, strValue: "group"}
	COLLECT_KEYWORD = SyntaxKind{tag: 402, strValue: "collect"}

	// Separators
	OPEN_BRACE_TOKEN       = SyntaxKind{tag: 500, strValue: "{"}
	CLOSE_BRACE_TOKEN      = SyntaxKind{tag: 501, strValue: "}"}
	OPEN_PAREN_TOKEN       = SyntaxKind{tag: 502, strValue: "("}
	CLOSE_PAREN_TOKEN      = SyntaxKind{tag: 503, strValue: ")"}
	OPEN_BRACKET_TOKEN     = SyntaxKind{tag: 504, strValue: "["}
	CLOSE_BRACKET_TOKEN    = SyntaxKind{tag: 505, strValue: "]"}
	SEMICOLON_TOKEN        = SyntaxKind{tag: 506, strValue: ";"}
	DOT_TOKEN              = SyntaxKind{tag: 507, strValue: "."}
	COLON_TOKEN            = SyntaxKind{tag: 508, strValue: ":"}
	COMMA_TOKEN            = SyntaxKind{tag: 509, strValue: ","}
	ELLIPSIS_TOKEN         = SyntaxKind{tag: 510, strValue: "..."}
	OPEN_BRACE_PIPE_TOKEN  = SyntaxKind{tag: 511, strValue: "{|"}
	CLOSE_BRACE_PIPE_TOKEN = SyntaxKind{tag: 512, strValue: "|}"}
	AT_TOKEN               = SyntaxKind{tag: 513, strValue: "@"}
	HASH_TOKEN             = SyntaxKind{tag: 514, strValue: "#"}
	BACKTICK_TOKEN         = SyntaxKind{tag: 515, strValue: "`"}
	DOUBLE_QUOTE_TOKEN     = SyntaxKind{tag: 516, strValue: "\""}
	SINGLE_QUOTE_TOKEN     = SyntaxKind{tag: 517, strValue: "'"}
	DOUBLE_BACKTICK_TOKEN  = SyntaxKind{tag: 518, strValue: "``"}
	TRIPLE_BACKTICK_TOKEN  = SyntaxKind{tag: 519, strValue: "```"}

	// Operators
	EQUAL_TOKEN                           = SyntaxKind{tag: 550, strValue: "="}
	DOUBLE_EQUAL_TOKEN                    = SyntaxKind{tag: 551, strValue: "=="}
	TRIPPLE_EQUAL_TOKEN                   = SyntaxKind{tag: 552, strValue: "==="}
	PLUS_TOKEN                            = SyntaxKind{tag: 553, strValue: "+"}
	MINUS_TOKEN                           = SyntaxKind{tag: 554, strValue: "-"}
	SLASH_TOKEN                           = SyntaxKind{tag: 555, strValue: "/"}
	PERCENT_TOKEN                         = SyntaxKind{tag: 556, strValue: "%"}
	ASTERISK_TOKEN                        = SyntaxKind{tag: 557, strValue: "*"}
	LT_TOKEN                              = SyntaxKind{tag: 558, strValue: "<"}
	LT_EQUAL_TOKEN                        = SyntaxKind{tag: 559, strValue: "<="}
	GT_TOKEN                              = SyntaxKind{tag: 560, strValue: ">"}
	RIGHT_DOUBLE_ARROW_TOKEN              = SyntaxKind{tag: 561, strValue: "=>"}
	QUESTION_MARK_TOKEN                   = SyntaxKind{tag: 562, strValue: "?"}
	PIPE_TOKEN                            = SyntaxKind{tag: 563, strValue: "|"}
	GT_EQUAL_TOKEN                        = SyntaxKind{tag: 564, strValue: ">="}
	EXCLAMATION_MARK_TOKEN                = SyntaxKind{tag: 565, strValue: "!"}
	NOT_EQUAL_TOKEN                       = SyntaxKind{tag: 566, strValue: "!="}
	NOT_DOUBLE_EQUAL_TOKEN                = SyntaxKind{tag: 567, strValue: "!=="}
	BITWISE_AND_TOKEN                     = SyntaxKind{tag: 568, strValue: "&"}
	BITWISE_XOR_TOKEN                     = SyntaxKind{tag: 569, strValue: "^"}
	LOGICAL_AND_TOKEN                     = SyntaxKind{tag: 570, strValue: "&&"}
	LOGICAL_OR_TOKEN                      = SyntaxKind{tag: 571, strValue: "||"}
	NEGATION_TOKEN                        = SyntaxKind{tag: 572, strValue: "~"}
	RIGHT_ARROW_TOKEN                     = SyntaxKind{tag: 573, strValue: "->"}
	INTERPOLATION_START_TOKEN             = SyntaxKind{tag: 574, strValue: "${"}
	XML_PI_START_TOKEN                    = SyntaxKind{tag: 575, strValue: "<?"}
	XML_PI_END_TOKEN                      = SyntaxKind{tag: 576, strValue: "?>"}
	XML_COMMENT_START_TOKEN               = SyntaxKind{tag: 577, strValue: "<!--"}
	XML_COMMENT_END_TOKEN                 = SyntaxKind{tag: 578, strValue: "-->"}
	SYNC_SEND_TOKEN                       = SyntaxKind{tag: 579, strValue: "->>"}
	LEFT_ARROW_TOKEN                      = SyntaxKind{tag: 580, strValue: "<-"}
	DOUBLE_DOT_LT_TOKEN                   = SyntaxKind{tag: 580, strValue: "..<"}
	DOUBLE_LT_TOKEN                       = SyntaxKind{tag: 581, strValue: "<<"}
	ANNOT_CHAINING_TOKEN                  = SyntaxKind{tag: 582, strValue: ".@"}
	OPTIONAL_CHAINING_TOKEN               = SyntaxKind{tag: 583, strValue: "?."}
	ELVIS_TOKEN                           = SyntaxKind{tag: 584, strValue: "?:"}
	DOT_LT_TOKEN                          = SyntaxKind{tag: 585, strValue: ".<"}
	SLASH_LT_TOKEN                        = SyntaxKind{tag: 586, strValue: "/<"}
	DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN = SyntaxKind{tag: 587, strValue: "/**/<"}
	SLASH_ASTERISK_TOKEN                  = SyntaxKind{tag: 588, strValue: "/*"}
	DOUBLE_GT_TOKEN                       = SyntaxKind{tag: 589, strValue: ">>"}
	TRIPPLE_GT_TOKEN                      = SyntaxKind{tag: 590, strValue: ">>>"}
	XML_CDATA_START_TOKEN                 = SyntaxKind{tag: 591, strValue: "<![CDATA["}
	XML_CDATA_END_TOKEN                   = SyntaxKind{tag: 592, strValue: "]]>"}
	BACK_SLASH_TOKEN                      = SyntaxKind{tag: 593, strValue: "\\"}
	DOLLAR_TOKEN                          = SyntaxKind{tag: 594, strValue: "$"}
	ESCAPED_MINUS_TOKEN                   = SyntaxKind{tag: 595, strValue: "\\-"}

	// Documentation reference types
	TYPE_DOC_REFERENCE_TOKEN       = SyntaxKind{tag: 900, strValue: "type"}
	SERVICE_DOC_REFERENCE_TOKEN    = SyntaxKind{tag: 901, strValue: "service"}
	VARIABLE_DOC_REFERENCE_TOKEN   = SyntaxKind{tag: 902, strValue: "variable"}
	VAR_DOC_REFERENCE_TOKEN        = SyntaxKind{tag: 903, strValue: "var"}
	ANNOTATION_DOC_REFERENCE_TOKEN = SyntaxKind{tag: 904, strValue: "annotation"}
	MODULE_DOC_REFERENCE_TOKEN     = SyntaxKind{tag: 905, strValue: "module"}
	FUNCTION_DOC_REFERENCE_TOKEN   = SyntaxKind{tag: 906, strValue: "function"}
	PARAMETER_DOC_REFERENCE_TOKEN  = SyntaxKind{tag: 907, strValue: "parameter"}
	CONST_DOC_REFERENCE_TOKEN      = SyntaxKind{tag: 908, strValue: "const"}

	// Literal tokens
	IDENTIFIER_TOKEN                     = SyntaxKind{tag: 1000, strValue: ""}
	STRING_LITERAL_TOKEN                 = SyntaxKind{tag: 1001, strValue: ""}
	DECIMAL_INTEGER_LITERAL_TOKEN        = SyntaxKind{tag: 1002, strValue: ""}
	HEX_INTEGER_LITERAL_TOKEN            = SyntaxKind{tag: 1003, strValue: ""}
	DECIMAL_FLOATING_POINT_LITERAL_TOKEN = SyntaxKind{tag: 1004, strValue: ""}
	HEX_FLOATING_POINT_LITERAL_TOKEN     = SyntaxKind{tag: 1005, strValue: ""}
	XML_TEXT_CONTENT                     = SyntaxKind{tag: 1006, strValue: ""}
	TEMPLATE_STRING                      = SyntaxKind{tag: 1007, strValue: ""}
	PROMPT_CONTENT                       = SyntaxKind{tag: 1007, strValue: ""}

	// Documentation
	DOCUMENTATION_DESCRIPTION = SyntaxKind{tag: 1100, strValue: ""}
	PARAMETER_NAME            = SyntaxKind{tag: 1101, strValue: ""}
	CODE_CONTENT              = SyntaxKind{tag: 1102, strValue: ""}
	DEPRECATION_LITERAL       = SyntaxKind{tag: 1103, strValue: ""}
	DOCUMENTATION_STRING      = SyntaxKind{tag: 1104, strValue: ""}

	// Other
	INVALID_TOKEN = SyntaxKind{tag: 1191, strValue: ""}

	//-----------------------------------------------non-terminal-kinds-----------------------------------------------

	// Minutiae kinds
	WHITESPACE_MINUTIAE   = SyntaxKind{tag: 1500, strValue: ""}
	END_OF_LINE_MINUTIAE  = SyntaxKind{tag: 1501, strValue: ""}
	COMMENT_MINUTIAE      = SyntaxKind{tag: 1502, strValue: ""}
	INVALID_NODE_MINUTIAE = SyntaxKind{tag: 1503, strValue: ""}

	// Invalid nodes
	INVALID_TOKEN_MINUTIAE_NODE = SyntaxKind{tag: 1601, strValue: ""}

	// module-level declarations
	IMPORT_DECLARATION               = SyntaxKind{tag: 2000, strValue: ""}
	FUNCTION_DEFINITION              = SyntaxKind{tag: 2001, strValue: ""}
	TYPE_DEFINITION                  = SyntaxKind{tag: 2002, strValue: ""}
	SERVICE_DECLARATION              = SyntaxKind{tag: 2003, strValue: ""}
	MODULE_VAR_DECL                  = SyntaxKind{tag: 2004, strValue: ""}
	LISTENER_DECLARATION             = SyntaxKind{tag: 2005, strValue: ""}
	CONST_DECLARATION                = SyntaxKind{tag: 2006, strValue: ""}
	ANNOTATION_DECLARATION           = SyntaxKind{tag: 2007, strValue: ""}
	MODULE_XML_NAMESPACE_DECLARATION = SyntaxKind{tag: 2008, strValue: ""}
	ENUM_DECLARATION                 = SyntaxKind{tag: 2009, strValue: ""}
	CLASS_DEFINITION                 = SyntaxKind{tag: 2010, strValue: ""}

	// Statements
	BLOCK_STATEMENT                 = SyntaxKind{tag: 1200, strValue: ""}
	LOCAL_VAR_DECL                  = SyntaxKind{tag: 1201, strValue: ""}
	ASSIGNMENT_STATEMENT            = SyntaxKind{tag: 1202, strValue: ""}
	IF_ELSE_STATEMENT               = SyntaxKind{tag: 1203, strValue: ""}
	ELSE_BLOCK                      = SyntaxKind{tag: 1204, strValue: ""}
	WHILE_STATEMENT                 = SyntaxKind{tag: 1205, strValue: ""}
	CALL_STATEMENT                  = SyntaxKind{tag: 1206, strValue: ""}
	PANIC_STATEMENT                 = SyntaxKind{tag: 1207, strValue: ""}
	RETURN_STATEMENT                = SyntaxKind{tag: 1208, strValue: ""}
	CONTINUE_STATEMENT              = SyntaxKind{tag: 1209, strValue: ""}
	BREAK_STATEMENT                 = SyntaxKind{tag: 1210, strValue: ""}
	COMPOUND_ASSIGNMENT_STATEMENT   = SyntaxKind{tag: 1211, strValue: ""}
	LOCAL_TYPE_DEFINITION_STATEMENT = SyntaxKind{tag: 1212, strValue: ""}
	ACTION_STATEMENT                = SyntaxKind{tag: 1213, strValue: ""}
	LOCK_STATEMENT                  = SyntaxKind{tag: 1214, strValue: ""}
	NAMED_WORKER_DECLARATION        = SyntaxKind{tag: 1215, strValue: ""}
	FORK_STATEMENT                  = SyntaxKind{tag: 1216, strValue: ""}
	FOREACH_STATEMENT               = SyntaxKind{tag: 1217, strValue: ""}
	TRANSACTION_STATEMENT           = SyntaxKind{tag: 1218, strValue: ""}
	ROLLBACK_STATEMENT              = SyntaxKind{tag: 1219, strValue: ""}
	RETRY_STATEMENT                 = SyntaxKind{tag: 1220, strValue: ""}
	XML_NAMESPACE_DECLARATION       = SyntaxKind{tag: 1221, strValue: ""}
	MATCH_STATEMENT                 = SyntaxKind{tag: 1222, strValue: ""}
	INVALID_EXPRESSION_STATEMENT    = SyntaxKind{tag: 1223, strValue: ""}
	DO_STATEMENT                    = SyntaxKind{tag: 1224, strValue: ""}
	FAIL_STATEMENT                  = SyntaxKind{tag: 1225, strValue: ""}

	// Expressions
	BINARY_EXPRESSION                      = SyntaxKind{tag: 1300, strValue: ""}
	BRACED_EXPRESSION                      = SyntaxKind{tag: 1301, strValue: ""}
	FUNCTION_CALL                          = SyntaxKind{tag: 1302, strValue: ""}
	QUALIFIED_NAME_REFERENCE               = SyntaxKind{tag: 1303, strValue: ""}
	INDEXED_EXPRESSION                     = SyntaxKind{tag: 1304, strValue: ""}
	FIELD_ACCESS                           = SyntaxKind{tag: 1305, strValue: ""}
	METHOD_CALL                            = SyntaxKind{tag: 1306, strValue: ""}
	CHECK_EXPRESSION                       = SyntaxKind{tag: 1307, strValue: ""}
	MAPPING_CONSTRUCTOR                    = SyntaxKind{tag: 1308, strValue: ""}
	TYPEOF_EXPRESSION                      = SyntaxKind{tag: 1309, strValue: ""}
	UNARY_EXPRESSION                       = SyntaxKind{tag: 1310, strValue: ""}
	TYPE_TEST_EXPRESSION                   = SyntaxKind{tag: 1311, strValue: ""}
	SIMPLE_NAME_REFERENCE                  = SyntaxKind{tag: 1313, strValue: ""}
	TRAP_EXPRESSION                        = SyntaxKind{tag: 1314, strValue: ""}
	LIST_CONSTRUCTOR                       = SyntaxKind{tag: 1315, strValue: ""}
	TYPE_CAST_EXPRESSION                   = SyntaxKind{tag: 1316, strValue: ""}
	TABLE_CONSTRUCTOR                      = SyntaxKind{tag: 1317, strValue: ""}
	LET_EXPRESSION                         = SyntaxKind{tag: 1318, strValue: ""}
	XML_TEMPLATE_EXPRESSION                = SyntaxKind{tag: 1319, strValue: ""}
	REGEX_TEMPLATE_EXPRESSION              = SyntaxKind{tag: 1346, strValue: ""}
	RAW_TEMPLATE_EXPRESSION                = SyntaxKind{tag: 1320, strValue: ""}
	STRING_TEMPLATE_EXPRESSION             = SyntaxKind{tag: 1321, strValue: ""}
	IMPLICIT_NEW_EXPRESSION                = SyntaxKind{tag: 1322, strValue: ""}
	EXPLICIT_NEW_EXPRESSION                = SyntaxKind{tag: 1323, strValue: ""}
	PARENTHESIZED_ARG_LIST                 = SyntaxKind{tag: 1324, strValue: ""}
	EXPLICIT_ANONYMOUS_FUNCTION_EXPRESSION = SyntaxKind{tag: 1325, strValue: ""}
	IMPLICIT_ANONYMOUS_FUNCTION_EXPRESSION = SyntaxKind{tag: 1326, strValue: ""}
	QUERY_EXPRESSION                       = SyntaxKind{tag: 1327, strValue: ""}
	ANNOT_ACCESS                           = SyntaxKind{tag: 1328, strValue: ""}
	OPTIONAL_FIELD_ACCESS                  = SyntaxKind{tag: 1329, strValue: ""}
	CONDITIONAL_EXPRESSION                 = SyntaxKind{tag: 1330, strValue: ""}
	TRANSACTIONAL_EXPRESSION               = SyntaxKind{tag: 1331, strValue: ""}
	OBJECT_CONSTRUCTOR                     = SyntaxKind{tag: 1332, strValue: ""}
	XML_FILTER_EXPRESSION                  = SyntaxKind{tag: 1333, strValue: ""}
	XML_STEP_EXPRESSION                    = SyntaxKind{tag: 1334, strValue: ""}
	XML_NAME_PATTERN_CHAIN                 = SyntaxKind{tag: 1335, strValue: ""}
	XML_ATOMIC_NAME_PATTERN                = SyntaxKind{tag: 1336, strValue: ""}
	STRING_LITERAL                         = SyntaxKind{tag: 1337, strValue: ""}
	NUMERIC_LITERAL                        = SyntaxKind{tag: 1338, strValue: ""}
	BOOLEAN_LITERAL                        = SyntaxKind{tag: 1339, strValue: ""}
	NIL_LITERAL                            = SyntaxKind{tag: 1340, strValue: ""}
	NULL_LITERAL                           = SyntaxKind{tag: 1341, strValue: ""}
	BYTE_ARRAY_LITERAL                     = SyntaxKind{tag: 1342, strValue: ""}
	ASTERISK_LITERAL                       = SyntaxKind{tag: 1343, strValue: ""}
	REQUIRED_EXPRESSION                    = SyntaxKind{tag: 1344, strValue: ""}
	ERROR_CONSTRUCTOR                      = SyntaxKind{tag: 1345, strValue: ""}
	XML_STEP_METHOD_CALL_EXTEND            = SyntaxKind{tag: 1346, strValue: ""}
	XML_STEP_INDEXED_EXTEND                = SyntaxKind{tag: 1347, strValue: ""}
	NATURAL_EXPRESSION                     = SyntaxKind{tag: 1348, strValue: ""}

	// Type descriptors
	TYPE_DESC                = SyntaxKind{tag: 2000, strValue: ""}
	RECORD_TYPE_DESC         = SyntaxKind{tag: 2001, strValue: ""}
	OBJECT_TYPE_DESC         = SyntaxKind{tag: 2002, strValue: ""}
	NIL_TYPE_DESC            = SyntaxKind{tag: 2003, strValue: ""}
	OPTIONAL_TYPE_DESC       = SyntaxKind{tag: 2004, strValue: ""}
	ARRAY_TYPE_DESC          = SyntaxKind{tag: 2005, strValue: ""}
	INT_TYPE_DESC            = SyntaxKind{tag: 2006, strValue: ""}
	BYTE_TYPE_DESC           = SyntaxKind{tag: 2007, strValue: ""}
	FLOAT_TYPE_DESC          = SyntaxKind{tag: 2008, strValue: ""}
	DECIMAL_TYPE_DESC        = SyntaxKind{tag: 2009, strValue: ""}
	STRING_TYPE_DESC         = SyntaxKind{tag: 2010, strValue: ""}
	BOOLEAN_TYPE_DESC        = SyntaxKind{tag: 2011, strValue: ""}
	XML_TYPE_DESC            = SyntaxKind{tag: 2012, strValue: ""}
	JSON_TYPE_DESC           = SyntaxKind{tag: 2013, strValue: ""}
	HANDLE_TYPE_DESC         = SyntaxKind{tag: 2014, strValue: ""}
	ANY_TYPE_DESC            = SyntaxKind{tag: 2015, strValue: ""}
	ANYDATA_TYPE_DESC        = SyntaxKind{tag: 2016, strValue: ""}
	NEVER_TYPE_DESC          = SyntaxKind{tag: 2017, strValue: ""}
	VAR_TYPE_DESC            = SyntaxKind{tag: 2018, strValue: ""}
	SERVICE_TYPE_DESC        = SyntaxKind{tag: 2019, strValue: ""}
	MAP_TYPE_DESC            = SyntaxKind{tag: 2020, strValue: ""}
	UNION_TYPE_DESC          = SyntaxKind{tag: 2021, strValue: ""}
	ERROR_TYPE_DESC          = SyntaxKind{tag: 2022, strValue: ""}
	STREAM_TYPE_DESC         = SyntaxKind{tag: 2023, strValue: ""}
	TABLE_TYPE_DESC          = SyntaxKind{tag: 2024, strValue: ""}
	FUNCTION_TYPE_DESC       = SyntaxKind{tag: 2025, strValue: ""}
	TUPLE_TYPE_DESC          = SyntaxKind{tag: 2026, strValue: ""}
	PARENTHESISED_TYPE_DESC  = SyntaxKind{tag: 2027, strValue: ""}
	READONLY_TYPE_DESC       = SyntaxKind{tag: 2028, strValue: ""}
	DISTINCT_TYPE_DESC       = SyntaxKind{tag: 2029, strValue: ""}
	INTERSECTION_TYPE_DESC   = SyntaxKind{tag: 2030, strValue: ""}
	SINGLETON_TYPE_DESC      = SyntaxKind{tag: 2031, strValue: ""}
	TYPE_REFERENCE_TYPE_DESC = SyntaxKind{tag: 2032, strValue: ""}
	TYPEDESC_TYPE_DESC       = SyntaxKind{tag: 2033, strValue: ""}
	FUTURE_TYPE_DESC         = SyntaxKind{tag: 2034, strValue: ""}

	// Actions
	REMOTE_METHOD_CALL_ACTION     = SyntaxKind{tag: 2500, strValue: ""}
	BRACED_ACTION                 = SyntaxKind{tag: 2501, strValue: ""}
	CHECK_ACTION                  = SyntaxKind{tag: 2502, strValue: ""}
	START_ACTION                  = SyntaxKind{tag: 2503, strValue: ""}
	TRAP_ACTION                   = SyntaxKind{tag: 2504, strValue: ""}
	FLUSH_ACTION                  = SyntaxKind{tag: 2505, strValue: ""}
	ASYNC_SEND_ACTION             = SyntaxKind{tag: 2506, strValue: ""}
	SYNC_SEND_ACTION              = SyntaxKind{tag: 2507, strValue: ""}
	RECEIVE_ACTION                = SyntaxKind{tag: 2508, strValue: ""}
	WAIT_ACTION                   = SyntaxKind{tag: 2509, strValue: ""}
	QUERY_ACTION                  = SyntaxKind{tag: 2510, strValue: ""}
	COMMIT_ACTION                 = SyntaxKind{tag: 2511, strValue: ""}
	CLIENT_RESOURCE_ACCESS_ACTION = SyntaxKind{tag: 2512, strValue: ""}

	// Other
	RETURN_TYPE_DESCRIPTOR            = SyntaxKind{tag: 3000, strValue: ""}
	REQUIRED_PARAM                    = SyntaxKind{tag: 3001, strValue: ""}
	DEFAULTABLE_PARAM                 = SyntaxKind{tag: 3002, strValue: ""}
	REST_PARAM                        = SyntaxKind{tag: 3003, strValue: ""}
	EXTERNAL_FUNCTION_BODY            = SyntaxKind{tag: 3004, strValue: ""}
	RECORD_FIELD                      = SyntaxKind{tag: 3005, strValue: ""}
	RECORD_FIELD_WITH_DEFAULT_VALUE   = SyntaxKind{tag: 3006, strValue: ""}
	TYPE_REFERENCE                    = SyntaxKind{tag: 3007, strValue: ""}
	RECORD_REST_TYPE                  = SyntaxKind{tag: 3008, strValue: ""}
	POSITIONAL_ARG                    = SyntaxKind{tag: 3009, strValue: ""}
	NAMED_ARG                         = SyntaxKind{tag: 3010, strValue: ""}
	REST_ARG                          = SyntaxKind{tag: 3011, strValue: ""}
	OBJECT_FIELD                      = SyntaxKind{tag: 3012, strValue: ""}
	IMPORT_ORG_NAME                   = SyntaxKind{tag: 3013, strValue: ""}
	MODULE_NAME                       = SyntaxKind{tag: 3014, strValue: ""}
	SUB_MODULE_NAME                   = SyntaxKind{tag: 3015, strValue: ""}
	IMPORT_VERSION                    = SyntaxKind{tag: 3016, strValue: ""}
	ORDER_BY_CLAUSE                   = SyntaxKind{tag: 3017, strValue: ""}
	IMPORT_PREFIX                     = SyntaxKind{tag: 3018, strValue: ""}
	SPECIFIC_FIELD                    = SyntaxKind{tag: 3019, strValue: ""}
	COMPUTED_NAME_FIELD               = SyntaxKind{tag: 3020, strValue: ""}
	SPREAD_FIELD                      = SyntaxKind{tag: 3021, strValue: ""}
	ORDER_KEY                         = SyntaxKind{tag: 3022, strValue: ""}
	RESOURCE_ACCESSOR_DEFINITION      = SyntaxKind{tag: 3023, strValue: ""}
	ANNOTATION                        = SyntaxKind{tag: 3024, strValue: ""}
	METADATA                          = SyntaxKind{tag: 3025, strValue: ""}
	ARRAY_DIMENSION                   = SyntaxKind{tag: 3026, strValue: ""}
	ANNOTATION_ATTACH_POINT           = SyntaxKind{tag: 3028, strValue: ""}
	FUNCTION_BODY_BLOCK               = SyntaxKind{tag: 3029, strValue: ""}
	NAMED_WORKER_DECLARATOR           = SyntaxKind{tag: 3030, strValue: ""}
	EXPRESSION_FUNCTION_BODY          = SyntaxKind{tag: 3031, strValue: ""}
	TYPE_CAST_PARAM                   = SyntaxKind{tag: 3032, strValue: ""}
	KEY_SPECIFIER                     = SyntaxKind{tag: 3033, strValue: ""}
	EXPLICIT_TYPE_PARAMS              = SyntaxKind{tag: 3034, strValue: ""}
	LET_VAR_DECL                      = SyntaxKind{tag: 3035, strValue: ""}
	STREAM_TYPE_PARAMS                = SyntaxKind{tag: 3036, strValue: ""}
	FUNCTION_SIGNATURE                = SyntaxKind{tag: 3037, strValue: ""}
	INFER_PARAM_LIST                  = SyntaxKind{tag: 3038, strValue: ""}
	TYPE_PARAMETER                    = SyntaxKind{tag: 3039, strValue: ""}
	KEY_TYPE_CONSTRAINT               = SyntaxKind{tag: 3040, strValue: ""}
	QUERY_CONSTRUCT_TYPE              = SyntaxKind{tag: 3041, strValue: ""}
	FROM_CLAUSE                       = SyntaxKind{tag: 3042, strValue: ""}
	WHERE_CLAUSE                      = SyntaxKind{tag: 3043, strValue: ""}
	LET_CLAUSE                        = SyntaxKind{tag: 3044, strValue: ""}
	QUERY_PIPELINE                    = SyntaxKind{tag: 3045, strValue: ""}
	SELECT_CLAUSE                     = SyntaxKind{tag: 3046, strValue: ""}
	METHOD_DECLARATION                = SyntaxKind{tag: 3047, strValue: ""}
	TYPED_BINDING_PATTERN             = SyntaxKind{tag: 3048, strValue: ""}
	BINDING_PATTERN                   = SyntaxKind{tag: 3049, strValue: ""}
	CAPTURE_BINDING_PATTERN           = SyntaxKind{tag: 3050, strValue: ""}
	REST_BINDING_PATTERN              = SyntaxKind{tag: 3051, strValue: ""}
	LIST_BINDING_PATTERN              = SyntaxKind{tag: 3052, strValue: ""}
	RECEIVE_FIELDS                    = SyntaxKind{tag: 3053, strValue: ""}
	REST_TYPE                         = SyntaxKind{tag: 3054, strValue: ""}
	WAIT_FIELDS_LIST                  = SyntaxKind{tag: 3055, strValue: ""}
	WAIT_FIELD                        = SyntaxKind{tag: 3056, strValue: ""}
	ENUM_MEMBER                       = SyntaxKind{tag: 3057, strValue: ""}
	BRACKETED_LIST                    = SyntaxKind{tag: 3058, strValue: ""}
	LIST_BP_OR_LIST_CONSTRUCTOR       = SyntaxKind{tag: 3059, strValue: ""}
	MAPPING_BINDING_PATTERN           = SyntaxKind{tag: 3060, strValue: ""}
	FIELD_BINDING_PATTERN             = SyntaxKind{tag: 3061, strValue: ""}
	MAPPING_BP_OR_MAPPING_CONSTRUCTOR = SyntaxKind{tag: 3062, strValue: ""}
	WILDCARD_BINDING_PATTERN          = SyntaxKind{tag: 3063, strValue: ""}
	MATCH_CLAUSE                      = SyntaxKind{tag: 3064, strValue: ""}
	MATCH_GUARD                       = SyntaxKind{tag: 3065, strValue: ""}
	OBJECT_METHOD_DEFINITION          = SyntaxKind{tag: 3066, strValue: ""}
	ON_CONFLICT_CLAUSE                = SyntaxKind{tag: 3067, strValue: ""}
	LIMIT_CLAUSE                      = SyntaxKind{tag: 3068, strValue: ""}
	JOIN_CLAUSE                       = SyntaxKind{tag: 3069, strValue: ""}
	ON_CLAUSE                         = SyntaxKind{tag: 3070, strValue: ""}
	LIST_MATCH_PATTERN                = SyntaxKind{tag: 3071, strValue: ""}
	REST_MATCH_PATTERN                = SyntaxKind{tag: 3072, strValue: ""}
	MAPPING_MATCH_PATTERN             = SyntaxKind{tag: 3073, strValue: ""}
	FIELD_MATCH_PATTERN               = SyntaxKind{tag: 3074, strValue: ""}
	ERROR_MATCH_PATTERN               = SyntaxKind{tag: 3075, strValue: ""}
	NAMED_ARG_MATCH_PATTERN           = SyntaxKind{tag: 3076, strValue: ""}
	ERROR_BINDING_PATTERN             = SyntaxKind{tag: 3077, strValue: ""}
	NAMED_ARG_BINDING_PATTERN         = SyntaxKind{tag: 3078, strValue: ""}
	TUPLE_TYPE_DESC_OR_LIST_CONST     = SyntaxKind{tag: 3079, strValue: ""}
	ON_FAIL_CLAUSE                    = SyntaxKind{tag: 3080, strValue: ""}
	RESOURCE_ACCESSOR_DECLARATION     = SyntaxKind{tag: 3081, strValue: ""}
	RESOURCE_PATH_SEGMENT_PARAM       = SyntaxKind{tag: 3082, strValue: ""}
	RESOURCE_PATH_REST_PARAM          = SyntaxKind{tag: 3083, strValue: ""}
	INCLUDED_RECORD_PARAM             = SyntaxKind{tag: 3084, strValue: ""}
	ARRAY_TYPE_DESC_OR_MEMBER_ACCESS  = SyntaxKind{tag: 3085, strValue: ""}
	INFERRED_TYPEDESC_DEFAULT         = SyntaxKind{tag: 3086, strValue: ""}
	SPREAD_MEMBER                     = SyntaxKind{tag: 3087, strValue: ""}
	COMPUTED_RESOURCE_ACCESS_SEGMENT  = SyntaxKind{tag: 3088, strValue: ""}
	RESOURCE_ACCESS_REST_SEGMENT      = SyntaxKind{tag: 3089, strValue: ""}
	MEMBER_TYPE_DESC                  = SyntaxKind{tag: 3090, strValue: ""}
	GROUPING_KEY_VAR_DECLARATION      = SyntaxKind{tag: 3091, strValue: ""}
	GROUPING_KEY_VAR_NAME             = SyntaxKind{tag: 3092, strValue: ""}
	GROUP_BY_CLAUSE                   = SyntaxKind{tag: 3093, strValue: ""}
	COLLECT_CLAUSE                    = SyntaxKind{tag: 3094, strValue: ""}
	ALTERNATE_RECEIVE                 = SyntaxKind{tag: 3095, strValue: ""}
	RECEIVE_FIELD                     = SyntaxKind{tag: 3096, strValue: ""}

	// XML
	XML_ELEMENT           = SyntaxKind{tag: 4000, strValue: ""}
	XML_EMPTY_ELEMENT     = SyntaxKind{tag: 4001, strValue: ""}
	XML_TEXT              = SyntaxKind{tag: 4002, strValue: ""}
	XML_COMMENT           = SyntaxKind{tag: 4003, strValue: ""}
	XML_PI                = SyntaxKind{tag: 4004, strValue: ""}
	XML_ELEMENT_START_TAG = SyntaxKind{tag: 4005, strValue: ""}
	XML_ELEMENT_END_TAG   = SyntaxKind{tag: 4006, strValue: ""}
	XML_SIMPLE_NAME       = SyntaxKind{tag: 4007, strValue: ""}
	XML_QUALIFIED_NAME    = SyntaxKind{tag: 4008, strValue: ""}
	XML_ATTRIBUTE         = SyntaxKind{tag: 4009, strValue: ""}
	XML_ATTRIBUTE_VALUE   = SyntaxKind{tag: 4010, strValue: ""}
	INTERPOLATION         = SyntaxKind{tag: 4011, strValue: ""}
	XML_CDATA             = SyntaxKind{tag: 4012, strValue: ""}

	// Reg Exp
	RE_SEQUENCE                                       = SyntaxKind{tag: 4013, strValue: ""}
	RE_ATOM_QUANTIFIER                                = SyntaxKind{tag: 4014, strValue: ""}
	RE_ASSERTION                                      = SyntaxKind{tag: 4015, strValue: ""}
	RE_LITERAL_CHAR_DOT_OR_ESCAPE                     = SyntaxKind{tag: 4016, strValue: ""}
	RE_QUOTE_ESCAPE                                   = SyntaxKind{tag: 4017, strValue: ""}
	RE_SIMPLE_CHAR_CLASS_ESCAPE                       = SyntaxKind{tag: 4018, strValue: ""}
	RE_UNICODE_PROPERTY_ESCAPE                        = SyntaxKind{tag: 4019, strValue: ""}
	RE_UNICODE_SCRIPT                                 = SyntaxKind{tag: 4020, strValue: ""}
	RE_UNICODE_GENERAL_CATEGORY                       = SyntaxKind{tag: 4021, strValue: ""}
	RE_CHARACTER_CLASS                                = SyntaxKind{tag: 4022, strValue: ""}
	RE_CHAR_SET_ATOM_WITH_RE_CHAR_SET_NO_DASH         = SyntaxKind{tag: 4023, strValue: ""}
	RE_CHAR_SET_ATOM_NO_DASH_WITH_RE_CHAR_SET_NO_DASH = SyntaxKind{tag: 4024, strValue: ""}
	RE_CHAR_SET_RANGE                                 = SyntaxKind{tag: 4025, strValue: ""}
	RE_CHAR_SET_RANGE_NO_DASH                         = SyntaxKind{tag: 4026, strValue: ""}
	RE_CHAR_SET_RANGE_WITH_RE_CHAR_SET                = SyntaxKind{tag: 4027, strValue: ""}
	RE_CHAR_SET_RANGE_NO_DASH_WITH_RE_CHAR_SET        = SyntaxKind{tag: 4028, strValue: ""}
	RE_CAPTURING_GROUP                                = SyntaxKind{tag: 4029, strValue: ""}
	RE_FLAG_EXPR                                      = SyntaxKind{tag: 4030, strValue: ""}
	RE_FLAGS_ON_OFF                                   = SyntaxKind{tag: 4031, strValue: ""}
	RE_FLAGS                                          = SyntaxKind{tag: 4032, strValue: ""}
	RE_QUANTIFIER                                     = SyntaxKind{tag: 4033, strValue: ""}
	RE_BRACED_QUANTIFIER                              = SyntaxKind{tag: 4034, strValue: ""}

	RE_ASSERTION_VALUE                = SyntaxKind{tag: 4035, strValue: ""}
	RE_LITERAL_CHAR                   = SyntaxKind{tag: 4036, strValue: ""}
	RE_NUMERIC_ESCAPE                 = SyntaxKind{tag: 4037, strValue: ""}
	RE_CONTROL_ESCAPE                 = SyntaxKind{tag: 4038, strValue: ""}
	RE_SIMPLE_CHAR_CLASS_CODE         = SyntaxKind{tag: 4039, strValue: ""}
	RE_PROPERTY                       = SyntaxKind{tag: 4040, strValue: ""}
	RE_UNICODE_SCRIPT_START           = SyntaxKind{tag: 4041, strValue: ""}
	RE_UNICODE_PROPERTY_VALUE         = SyntaxKind{tag: 4042, strValue: ""}
	RE_UNICODE_GENERAL_CATEGORY_START = SyntaxKind{tag: 4043, strValue: ""}
	RE_UNICODE_GENERAL_CATEGORY_NAME  = SyntaxKind{tag: 4044, strValue: ""}
	RE_CHAR_SET_ATOM_NO_DASH          = SyntaxKind{tag: 4045, strValue: ""}
	RE_FLAGS_VALUE                    = SyntaxKind{tag: 4046, strValue: ""}
	RE_BASE_QUANTIFIER_VALUE          = SyntaxKind{tag: 4047, strValue: ""}
	DIGIT                             = SyntaxKind{tag: 4048, strValue: ""}

	// Documentation
	MARKDOWN_DOCUMENTATION                       = SyntaxKind{tag: 4500, strValue: ""}
	MARKDOWN_DOCUMENTATION_LINE                  = SyntaxKind{tag: 4501, strValue: ""}
	MARKDOWN_REFERENCE_DOCUMENTATION_LINE        = SyntaxKind{tag: 4502, strValue: ""}
	MARKDOWN_PARAMETER_DOCUMENTATION_LINE        = SyntaxKind{tag: 4503, strValue: ""}
	MARKDOWN_RETURN_PARAMETER_DOCUMENTATION_LINE = SyntaxKind{tag: 4504, strValue: ""}
	MARKDOWN_DEPRECATION_DOCUMENTATION_LINE      = SyntaxKind{tag: 4505, strValue: ""}
	MARKDOWN_CODE_LINE                           = SyntaxKind{tag: 4506, strValue: ""}
	BALLERINA_NAME_REFERENCE                     = SyntaxKind{tag: 4507, strValue: ""}
	MARKDOWN_CODE_BLOCK                          = SyntaxKind{tag: 4508, strValue: ""}
	INLINE_CODE_REFERENCE                        = SyntaxKind{tag: 4509, strValue: ""}

	INVALID     = SyntaxKind{tag: 4, strValue: ""}
	MODULE_PART = SyntaxKind{tag: 3, strValue: ""}
	EOF_TOKEN   = SyntaxKind{tag: 2, strValue: ""}
	LIST        = SyntaxKind{tag: 1, strValue: ""}
	NONE        = SyntaxKind{tag: 0, strValue: ""}
)

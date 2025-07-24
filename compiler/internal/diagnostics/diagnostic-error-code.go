package diagnostics

import (
	"ballerina-lang-go/tools/diagnostics"
)

type diagnosticErrorCode struct {
	diagnosticId string
	messageKey   string
}

// Implement DiagnosticCode interface
func (d diagnosticErrorCode) Severity() diagnostics.DiagnosticSeverity {
	return diagnostics.Error
}

func (d diagnosticErrorCode) DiagnosticId() string {
	return d.diagnosticId
}

func (d diagnosticErrorCode) MessageKey() string {
	return d.messageKey
}

// LookupKey for structural equality comparison
type diagnosticErrorCodeLookupKey struct {
	diagnosticId string
	messageKey   string
}

func (d diagnosticErrorCode) LookupKey() diagnosticErrorCodeLookupKey {
	return diagnosticErrorCodeLookupKey{
		diagnosticId: d.diagnosticId,
		messageKey:   d.messageKey,
	}
}

func (d diagnosticErrorCode) Equals(other diagnosticErrorCode) bool {
	return d.LookupKey() == other.LookupKey()
}

var (
	ERROR_SYNTAX_ERROR = diagnosticErrorCode{diagnosticId: "BCE0000", messageKey: "error.syntax.error"}

	// Missing tokens
	ERROR_MISSING_TOKEN                                 = diagnosticErrorCode{diagnosticId: "BCE0001", messageKey: "error.missing.token"}
	ERROR_MISSING_SEMICOLON_TOKEN                       = diagnosticErrorCode{diagnosticId: "BCE0002", messageKey: "error.missing.semicolon.token"}
	ERROR_MISSING_COLON_TOKEN                           = diagnosticErrorCode{diagnosticId: "BCE0003", messageKey: "error.missing.colon.token"}
	ERROR_MISSING_OPEN_PAREN_TOKEN                      = diagnosticErrorCode{diagnosticId: "BCE0004", messageKey: "error.missing.open.paren.token"}
	ERROR_MISSING_CLOSE_PAREN_TOKEN                     = diagnosticErrorCode{diagnosticId: "BCE0005", messageKey: "error.missing.close.paren.token"}
	ERROR_MISSING_OPEN_BRACE_TOKEN                      = diagnosticErrorCode{diagnosticId: "BCE0006", messageKey: "error.missing.open.brace.token"}
	ERROR_MISSING_CLOSE_BRACE_TOKEN                     = diagnosticErrorCode{diagnosticId: "BCE0007", messageKey: "error.missing.close.brace.token"}
	ERROR_MISSING_OPEN_BRACKET_TOKEN                    = diagnosticErrorCode{diagnosticId: "BCE0008", messageKey: "error.missing.open.bracket.token"}
	ERROR_MISSING_CLOSE_BRACKET_TOKEN                   = diagnosticErrorCode{diagnosticId: "BCE0009", messageKey: "error.missing.close.bracket.token"}
	ERROR_MISSING_EQUAL_TOKEN                           = diagnosticErrorCode{diagnosticId: "BCE0010", messageKey: "error.missing.equal.token"}
	ERROR_MISSING_COMMA_TOKEN                           = diagnosticErrorCode{diagnosticId: "BCE0011", messageKey: "error.missing.comma.token"}
	ERROR_MISSING_BINARY_OPERATOR                       = diagnosticErrorCode{diagnosticId: "BCE0012", messageKey: "error.missing.binary.operator"}
	ERROR_MISSING_SLASH_TOKEN                           = diagnosticErrorCode{diagnosticId: "BCE0013", messageKey: "error.missing.slash.token"}
	ERROR_MISSING_AT_TOKEN                              = diagnosticErrorCode{diagnosticId: "BCE0014", messageKey: "error.missing.at.token"}
	ERROR_MISSING_QUESTION_MARK_TOKEN                   = diagnosticErrorCode{diagnosticId: "BCE0015", messageKey: "error.missing.question.mark.token"}
	ERROR_MISSING_GT_TOKEN                              = diagnosticErrorCode{diagnosticId: "BCE0016", messageKey: "error.missing.gt.token"}
	ERROR_MISSING_GT_EQUAL_TOKEN                        = diagnosticErrorCode{diagnosticId: "BCE0017", messageKey: "error.missing.gt.equal.token"}
	ERROR_MISSING_LT_TOKEN                              = diagnosticErrorCode{diagnosticId: "BCE0018", messageKey: "error.missing.lt.token"}
	ERROR_MISSING_LT_EQUAL_TOKEN                        = diagnosticErrorCode{diagnosticId: "BCE0019", messageKey: "error.missing.lt.equal.token"}
	ERROR_MISSING_RIGHT_DOUBLE_ARROW_TOKEN              = diagnosticErrorCode{diagnosticId: "BCE0020", messageKey: "error.missing.right.double.arrow.token"}
	ERROR_MISSING_XML_COMMENT_END_TOKEN                 = diagnosticErrorCode{diagnosticId: "BCE0021", messageKey: "error.missing.xml.comment.end.token"}
	ERROR_MISSING_XML_PI_END_TOKEN                      = diagnosticErrorCode{diagnosticId: "BCE0022", messageKey: "error.missing.xml.pi.end.token"}
	ERROR_MISSING_DOUBLE_QUOTE_TOKEN                    = diagnosticErrorCode{diagnosticId: "BCE0023", messageKey: "error.missing.double.quote.token"}
	ERROR_MISSING_BACKTICK_TOKEN                        = diagnosticErrorCode{diagnosticId: "BCE0024", messageKey: "error.missing.backtick.token"}
	ERROR_MISSING_OPEN_BRACE_PIPE_TOKEN                 = diagnosticErrorCode{diagnosticId: "BCE0025", messageKey: "error.missing.open.brace.pipe.token"}
	ERROR_MISSING_CLOSE_BRACE_PIPE_TOKEN                = diagnosticErrorCode{diagnosticId: "BCE0026", messageKey: "error.missing.close.brace.pipe.token"}
	ERROR_MISSING_ASTERISK_TOKEN                        = diagnosticErrorCode{diagnosticId: "BCE0027", messageKey: "error.missing.asterisk.token"}
	ERROR_MISSING_PIPE_TOKEN                            = diagnosticErrorCode{diagnosticId: "BCE0028", messageKey: "error.missing.pipe.token"}
	ERROR_MISSING_DOT_TOKEN                             = diagnosticErrorCode{diagnosticId: "BCE0029", messageKey: "error.missing.dot.token"}
	ERROR_MISSING_ELLIPSIS_TOKEN                        = diagnosticErrorCode{diagnosticId: "BCE0030", messageKey: "error.missing.ellipsis.token"}
	ERROR_MISSING_HASH_TOKEN                            = diagnosticErrorCode{diagnosticId: "BCE0031", messageKey: "error.missing.hash.token"}
	ERROR_MISSING_SINGLE_QUOTE_TOKEN                    = diagnosticErrorCode{diagnosticId: "BCE0032", messageKey: "error.missing.single.quote.token"}
	ERROR_MISSING_DOUBLE_EQUAL_TOKEN                    = diagnosticErrorCode{diagnosticId: "BCE0033", messageKey: "error.missing.double.equal.token"}
	ERROR_MISSING_TRIPPLE_EQUAL_TOKEN                   = diagnosticErrorCode{diagnosticId: "BCE0034", messageKey: "error.missing.tripple.equal.token"}
	ERROR_MISSING_MINUS_TOKEN                           = diagnosticErrorCode{diagnosticId: "BCE0035", messageKey: "error.missing.minus.token"}
	ERROR_MISSING_PERCENT_TOKEN                         = diagnosticErrorCode{diagnosticId: "BCE0036", messageKey: "error.missing.percent.token"}
	ERROR_MISSING_EXCLAMATION_MARK_TOKEN                = diagnosticErrorCode{diagnosticId: "BCE0037", messageKey: "error.missing.exclamation.mark.token"}
	ERROR_MISSING_NOT_EQUAL_TOKEN                       = diagnosticErrorCode{diagnosticId: "BCE0038", messageKey: "error.missing.not.equal.token"}
	ERROR_MISSING_NOT_DOUBLE_EQUAL_TOKEN                = diagnosticErrorCode{diagnosticId: "BCE0039", messageKey: "error.missing.not.double.equal.token"}
	ERROR_MISSING_BITWISE_AND_TOKEN                     = diagnosticErrorCode{diagnosticId: "BCE0040", messageKey: "error.missing.bitwise.and.token"}
	ERROR_MISSING_BITWISE_XOR_TOKEN                     = diagnosticErrorCode{diagnosticId: "BCE0041", messageKey: "error.missing.bitwise.xor.token"}
	ERROR_MISSING_LOGICAL_AND_TOKEN                     = diagnosticErrorCode{diagnosticId: "BCE0042", messageKey: "error.missing.logical.and.token"}
	ERROR_MISSING_LOGICAL_OR_TOKEN                      = diagnosticErrorCode{diagnosticId: "BCE0043", messageKey: "error.missing.logical.or.token"}
	ERROR_MISSING_NEGATION_TOKEN                        = diagnosticErrorCode{diagnosticId: "BCE0044", messageKey: "error.missing.negation.token"}
	ERROR_MISSING_RIGHT_ARROW_TOKEN                     = diagnosticErrorCode{diagnosticId: "BCE0045", messageKey: "error.missing.right.arrow.token"}
	ERROR_MISSING_INTERPOLATION_START_TOKEN             = diagnosticErrorCode{diagnosticId: "BCE0046", messageKey: "error.missing.interpolation.start.token"}
	ERROR_MISSING_XML_PI_START_TOKEN                    = diagnosticErrorCode{diagnosticId: "BCE0047", messageKey: "error.missing.xml.pi.start.token"}
	ERROR_MISSING_XML_COMMENT_START_TOKEN               = diagnosticErrorCode{diagnosticId: "BCE0048", messageKey: "error.missing.xml.comment.start.token"}
	ERROR_MISSING_SYNC_SEND_TOKEN                       = diagnosticErrorCode{diagnosticId: "BCE0049", messageKey: "error.missing.sync.send.token"}
	ERROR_MISSING_LEFT_ARROW_TOKEN                      = diagnosticErrorCode{diagnosticId: "BCE0050", messageKey: "error.missing.left.arrow.token"}
	ERROR_MISSING_DOUBLE_DOT_LT_TOKEN                   = diagnosticErrorCode{diagnosticId: "BCE0051", messageKey: "error.missing.double.dot.lt.token"}
	ERROR_MISSING_DOUBLE_LT_TOKEN                       = diagnosticErrorCode{diagnosticId: "BCE0052", messageKey: "error.missing.double.lt.token"}
	ERROR_MISSING_ANNOT_CHAINING_TOKEN                  = diagnosticErrorCode{diagnosticId: "BCE0053", messageKey: "error.missing.annot.chaining.token"}
	ERROR_MISSING_OPTIONAL_CHAINING_TOKEN               = diagnosticErrorCode{diagnosticId: "BCE0054", messageKey: "error.missing.optional.chaining.token"}
	ERROR_MISSING_ELVIS_TOKEN                           = diagnosticErrorCode{diagnosticId: "BCE0055", messageKey: "error.missing.elvis.token"}
	ERROR_MISSING_DOT_LT_TOKEN                          = diagnosticErrorCode{diagnosticId: "BCE0056", messageKey: "error.missing.dot.lt.token"}
	ERROR_MISSING_SLASH_LT_TOKEN                        = diagnosticErrorCode{diagnosticId: "BCE0057", messageKey: "error.missing.slash.lt.token"}
	ERROR_MISSING_DOUBLE_SLASH_DOUBLE_ASTERISK_LT_TOKEN = diagnosticErrorCode{diagnosticId: "BCE0058", messageKey: "error.missing.double.slash.double.asterisk.lt.token"}
	ERROR_MISSING_SLASH_ASTERISK_TOKEN                  = diagnosticErrorCode{diagnosticId: "BCE0059", messageKey: "error.missing.slash.asterisk.token"}
	ERROR_MISSING_DOUBLE_GT_TOKEN                       = diagnosticErrorCode{diagnosticId: "BCE0060", messageKey: "error.missing.double.gt.token"}
	ERROR_MISSING_TRIPPLE_GT_TOKEN                      = diagnosticErrorCode{diagnosticId: "BCE0061", messageKey: "error.missing.tripple.gt.token"}
	ERROR_MISSING_XML_CDATA_END_TOKEN                   = diagnosticErrorCode{diagnosticId: "BCE0062", messageKey: "error.missing.xml.cdata.end.token"}

	// Missing keywords
	ERROR_MISSING_PUBLIC_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0200", messageKey: "error.missing.public.keyword"}
	ERROR_MISSING_PRIVATE_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0201", messageKey: "error.missing.private.keyword"}
	ERROR_MISSING_REMOTE_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0202", messageKey: "error.missing.remote.keyword"}
	ERROR_MISSING_ABSTRACT_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0203", messageKey: "error.missing.abstract.keyword"}
	ERROR_MISSING_CLIENT_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0204", messageKey: "error.missing.client.keyword"}
	ERROR_MISSING_LISTENER_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0205", messageKey: "error.missing.listener.keyword"}
	ERROR_MISSING_XMLNS_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0206", messageKey: "error.missing.xmlns.keyword"}
	ERROR_MISSING_RESOURCE_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0207", messageKey: "error.missing.resource.keyword"}
	ERROR_MISSING_FINAL_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0208", messageKey: "error.missing.final.keyword"}
	ERROR_MISSING_WORKER_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0209", messageKey: "error.missing.worker.keyword"}
	ERROR_MISSING_PARAMETER_KEYWORD     = diagnosticErrorCode{diagnosticId: "BCE0210", messageKey: "error.missing.parameter.keyword"}
	ERROR_MISSING_RETURNS_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0211", messageKey: "error.missing.returns.keyword"}
	ERROR_MISSING_RETURN_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0212", messageKey: "error.missing.return.keyword"}
	ERROR_MISSING_TRUE_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0213", messageKey: "error.missing.true.keyword"}
	ERROR_MISSING_FALSE_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0214", messageKey: "error.missing.false.keyword"}
	ERROR_MISSING_ELSE_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0215", messageKey: "error.missing.else.keyword"}
	ERROR_MISSING_WHILE_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0216", messageKey: "error.missing.while.keyword"}
	ERROR_MISSING_CHECK_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0217", messageKey: "error.missing.check.keyword"}
	ERROR_MISSING_CHECKPANIC_KEYWORD    = diagnosticErrorCode{diagnosticId: "BCE0218", messageKey: "error.missing.checkpanic.keyword"}
	ERROR_MISSING_PANIC_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0219", messageKey: "error.missing.panic.keyword"}
	ERROR_MISSING_CONTINUE_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0220", messageKey: "error.missing.continue.keyword"}
	ERROR_MISSING_BREAK_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0221", messageKey: "error.missing.break.keyword"}
	ERROR_MISSING_TYPEOF_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0222", messageKey: "error.missing.typeof.keyword"}
	ERROR_MISSING_IS_KEYWORD            = diagnosticErrorCode{diagnosticId: "BCE0223", messageKey: "error.missing.is.keyword"}
	ERROR_MISSING_NULL_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0224", messageKey: "error.missing.null.keyword"}
	ERROR_MISSING_LOCK_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0225", messageKey: "error.missing.lock.keyword"}
	ERROR_MISSING_FORK_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0226", messageKey: "error.missing.fork.keyword"}
	ERROR_MISSING_TRAP_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0227", messageKey: "error.missing.trap.keyword"}
	ERROR_MISSING_FOREACH_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0228", messageKey: "error.missing.foreach.keyword"}
	ERROR_MISSING_NEW_KEYWORD           = diagnosticErrorCode{diagnosticId: "BCE0229", messageKey: "error.missing.new.keyword"}
	ERROR_MISSING_WHERE_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0230", messageKey: "error.missing.where.keyword"}
	ERROR_MISSING_SELECT_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0231", messageKey: "error.missing.select.keyword"}
	ERROR_MISSING_START_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0232", messageKey: "error.missing.start.keyword"}
	ERROR_MISSING_FLUSH_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0233", messageKey: "error.missing.flush.keyword"}
	ERROR_MISSING_WAIT_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0234", messageKey: "error.missing.wait.keyword"}
	ERROR_MISSING_DO_KEYWORD            = diagnosticErrorCode{diagnosticId: "BCE0235", messageKey: "error.missing.do.keyword"}
	ERROR_MISSING_TRANSACTION_KEYWORD   = diagnosticErrorCode{diagnosticId: "BCE0236", messageKey: "error.missing.transaction.keyword"}
	ERROR_MISSING_TRANSACTIONAL_KEYWORD = diagnosticErrorCode{diagnosticId: "BCE0237", messageKey: "error.missing.transactional.keyword"}
	ERROR_MISSING_COMMIT_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0238", messageKey: "error.missing.commit.keyword"}
	ERROR_MISSING_ROLLBACK_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0239", messageKey: "error.missing.rollback.keyword"}
	ERROR_MISSING_RETRY_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0240", messageKey: "error.missing.retry.keyword"}
	ERROR_MISSING_BASE16_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0241", messageKey: "error.missing.base16.keyword"}
	ERROR_MISSING_BASE64_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0242", messageKey: "error.missing.base64.keyword"}
	ERROR_MISSING_MATCH_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0243", messageKey: "error.missing.match.keyword"}
	ERROR_MISSING_DEFAULT_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0244", messageKey: "error.missing.default.keyword"}
	ERROR_MISSING_TYPE_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0245", messageKey: "error.missing.type.keyword"}
	ERROR_MISSING_ON_KEYWORD            = diagnosticErrorCode{diagnosticId: "BCE0246", messageKey: "error.missing.on.keyword"}
	ERROR_MISSING_ANNOTATION_KEYWORD    = diagnosticErrorCode{diagnosticId: "BCE0247", messageKey: "error.missing.annotation.keyword"}
	ERROR_MISSING_FUNCTION_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0248", messageKey: "error.missing.function.keyword"}
	ERROR_MISSING_SOURCE_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0249", messageKey: "error.missing.source.keyword"}
	ERROR_MISSING_ENUM_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0250", messageKey: "error.missing.enum.keyword"}
	ERROR_MISSING_FIELD_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0251", messageKey: "error.missing.field.keyword"}
	ERROR_MISSING_VERSION_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0252", messageKey: "error.missing.version.keyword"}
	ERROR_MISSING_OBJECT_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0253", messageKey: "error.missing.object.keyword"}
	ERROR_MISSING_RECORD_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0254", messageKey: "error.missing.record.keyword"}
	ERROR_MISSING_SERVICE_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0255", messageKey: "error.missing.service.keyword"}
	ERROR_MISSING_AS_KEYWORD            = diagnosticErrorCode{diagnosticId: "BCE0256", messageKey: "error.missing.as.keyword"}
	ERROR_MISSING_LET_KEYWORD           = diagnosticErrorCode{diagnosticId: "BCE0257", messageKey: "error.missing.let.keyword"}
	ERROR_MISSING_TABLE_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0258", messageKey: "error.missing.table.keyword"}
	ERROR_MISSING_KEY_KEYWORD           = diagnosticErrorCode{diagnosticId: "BCE0259", messageKey: "error.missing.key.keyword"}
	ERROR_MISSING_FROM_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0260", messageKey: "error.missing.from.keyword"}
	ERROR_MISSING_IN_KEYWORD            = diagnosticErrorCode{diagnosticId: "BCE0261", messageKey: "error.missing.in.keyword"}
	ERROR_MISSING_IF_KEYWORD            = diagnosticErrorCode{diagnosticId: "BCE0262", messageKey: "error.missing.if.keyword"}
	ERROR_MISSING_IMPORT_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0263", messageKey: "error.missing.import.keyword"}
	ERROR_MISSING_CONST_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0264", messageKey: "error.missing.const.keyword"}
	ERROR_MISSING_EXTERNAL_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0265", messageKey: "error.missing.external.keyword"}
	ERROR_MISSING_ORDER_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0266", messageKey: "error.missing.order.keyword"}
	ERROR_MISSING_BY_KEYWORD            = diagnosticErrorCode{diagnosticId: "BCE0267", messageKey: "error.missing.by.keyword"}
	ERROR_MISSING_CONFLICT_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0268", messageKey: "error.missing.conflict.keyword"}
	ERROR_MISSING_LIMIT_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0269", messageKey: "error.missing.limit.keyword"}
	ERROR_MISSING_ASCENDING_KEYWORD     = diagnosticErrorCode{diagnosticId: "BCE0270", messageKey: "error.missing.ascending.keyword"}
	ERROR_MISSING_DESCENDING_KEYWORD    = diagnosticErrorCode{diagnosticId: "BCE0271", messageKey: "error.missing.descending.keyword"}
	ERROR_MISSING_JOIN_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0272", messageKey: "error.missing.join.keyword"}
	ERROR_MISSING_OUTER_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0273", messageKey: "error.missing.outer.keyword"}
	ERROR_MISSING_CLASS_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0274", messageKey: "error.missing.class.keyword"}
	ERROR_MISSING_FAIL_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0275", messageKey: "error.missing.fail.keyword"}
	ERROR_MISSING_EQUALS_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0276", messageKey: "error.missing.equals.keyword"}
	ERROR_MISSING_INT_KEYWORD           = diagnosticErrorCode{diagnosticId: "BCE0277", messageKey: "error.missing.int.keyword"}
	ERROR_MISSING_BYTE_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0278", messageKey: "error.missing.byte.keyword"}
	ERROR_MISSING_FLOAT_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0279", messageKey: "error.missing.float.keyword"}
	ERROR_MISSING_DECIMAL_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0280", messageKey: "error.missing.decimal.keyword"}
	ERROR_MISSING_STRING_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0281", messageKey: "error.missing.string.keyword"}
	ERROR_MISSING_BOOLEAN_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0282", messageKey: "error.missing.boolean.keyword"}
	ERROR_MISSING_XML_KEYWORD           = diagnosticErrorCode{diagnosticId: "BCE0283", messageKey: "error.missing.xml.keyword"}
	ERROR_MISSING_JSON_KEYWORD          = diagnosticErrorCode{diagnosticId: "BCE0284", messageKey: "error.missing.json.keyword"}
	ERROR_MISSING_HANDLE_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0285", messageKey: "error.missing.handle.keyword"}
	ERROR_MISSING_ANY_KEYWORD           = diagnosticErrorCode{diagnosticId: "BCE0286", messageKey: "error.missing.any.keyword"}
	ERROR_MISSING_ANYDATA_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0287", messageKey: "error.missing.anydata.keyword"}
	ERROR_MISSING_NEVER_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0288", messageKey: "error.missing.never.keyword"}
	ERROR_MISSING_VAR_KEYWORD           = diagnosticErrorCode{diagnosticId: "BCE0289", messageKey: "error.missing.var.keyword"}
	ERROR_MISSING_MAP_KEYWORD           = diagnosticErrorCode{diagnosticId: "BCE0290", messageKey: "error.missing.map.keyword"}
	ERROR_MISSING_ERROR_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0291", messageKey: "error.missing.error.keyword"}
	ERROR_MISSING_STREAM_KEYWORD        = diagnosticErrorCode{diagnosticId: "BCE0292", messageKey: "error.missing.stream.keyword"}
	ERROR_MISSING_READONLY_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0293", messageKey: "error.missing.readonly.keyword"}
	ERROR_MISSING_DISTINCT_KEYWORD      = diagnosticErrorCode{diagnosticId: "BCE0294", messageKey: "error.missing.distinct.keyword"}
	ERROR_MISSING_RE_KEYWORD            = diagnosticErrorCode{diagnosticId: "BCE0295", messageKey: "error.missing.re.keyword"}
	ERROR_MISSING_GROUP_KEYWORD         = diagnosticErrorCode{diagnosticId: "BCE0296", messageKey: "error.missing.group.keyword"}
	ERROR_MISSING_COLLECT_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0297", messageKey: "error.missing.collect.keyword"}
	ERROR_MISSING_NATURAL_KEYWORD       = diagnosticErrorCode{diagnosticId: "BCE0298", messageKey: "error.missing.natural.keyword"}

	// Missing other tokens
	ERROR_MISSING_IDENTIFIER                     = diagnosticErrorCode{diagnosticId: "BCE0400", messageKey: "error.missing.identifier"}
	ERROR_MISSING_STRING_LITERAL                 = diagnosticErrorCode{diagnosticId: "BCE0401", messageKey: "error.missing.string.literal"}
	ERROR_MISSING_DECIMAL_INTEGER_LITERAL        = diagnosticErrorCode{diagnosticId: "BCE0402", messageKey: "error.missing.decimal.integer.literal"}
	ERROR_MISSING_HEX_INTEGER_LITERAL            = diagnosticErrorCode{diagnosticId: "BCE0403", messageKey: "error.missing.hex.integer.literal"}
	ERROR_MISSING_DECIMAL_FLOATING_POINT_LITERAL = diagnosticErrorCode{diagnosticId: "BCE0404", messageKey: "error.missing.decimal.floating.point.literal"}
	ERROR_MISSING_HEX_FLOATING_POINT_LITERAL     = diagnosticErrorCode{diagnosticId: "BCE0405", messageKey: "error.missing.hex.floating.point.literal"}
	ERROR_MISSING_XML_TEXT_CONTENT               = diagnosticErrorCode{diagnosticId: "BCE0406", messageKey: "error.missing.xml.text.content"}
	ERROR_MISSING_TEMPLATE_STRING                = diagnosticErrorCode{diagnosticId: "BCE0407", messageKey: "error.missing.template.string"}
	ERROR_MISSING_BYTE_ARRAY_CONTENT             = diagnosticErrorCode{diagnosticId: "BCE0408", messageKey: "error.missing.byte.array.content"}
	ERROR_MISSING_DIGIT_AFTER_EXPONENT_INDICATOR = diagnosticErrorCode{diagnosticId: "BCE0409", messageKey: "error.missing.digit.after.exponent.indicator"}
	ERROR_MISSING_HEX_DIGIT_AFTER_DOT            = diagnosticErrorCode{diagnosticId: "BCE0410", messageKey: "error.missing.hex.digit.after.dot"}
	ERROR_MISSING_DOUBLE_QUOTE                   = diagnosticErrorCode{diagnosticId: "BCE0411", messageKey: "error.missing.double.quote"}
	ERROR_MISSING_ENTITY_REFERENCE_NAME          = diagnosticErrorCode{diagnosticId: "BCE0412", messageKey: "error.missing.entity.reference.name"}
	ERROR_MISSING_SEMICOLON_IN_XML_REFERENCE     = diagnosticErrorCode{diagnosticId: "BCE0413", messageKey: "error.missing.semicolon.in.xml.reference"}
	ERROR_MISSING_ATTACH_POINT_NAME              = diagnosticErrorCode{diagnosticId: "BCE0414", messageKey: "error.missing.attach.point.name"}
	ERROR_MISSING_HEX_NUMBER_AFTER_HEX_INDICATOR = diagnosticErrorCode{diagnosticId: "BCE0415", messageKey: "error.missing.hex.number.after.hex.indicator"}
	ERROR_MISSING_DIGIT_AFTER_DOT                = diagnosticErrorCode{diagnosticId: "BCE0416", messageKey: "error.missing.digit.after.dot"}
	ERROR_MISSING_RE_UNICODE_PROPERTY_VALUE      = diagnosticErrorCode{diagnosticId: "BCE0417", messageKey: "error.missing.unicode.property.value"}
	ERROR_MISSING_RE_QUANTIFIER_DIGIT            = diagnosticErrorCode{diagnosticId: "BCE0418", messageKey: "error.missing.digit.in.quantifier"}
	ERROR_MISSING_BACKSLASH                      = diagnosticErrorCode{diagnosticId: "BCE0420", messageKey: "error.missing.backslash"}

	// Missing non-terminal nodes
	ERROR_MISSING_FUNCTION_NAME                                  = diagnosticErrorCode{diagnosticId: "BCE0500", messageKey: "error.missing.function.name"}
	ERROR_MISSING_TYPE_DESC                                      = diagnosticErrorCode{diagnosticId: "BCE0501", messageKey: "error.missing.type.desc"}
	ERROR_MISSING_EXPRESSION                                     = diagnosticErrorCode{diagnosticId: "BCE0502", messageKey: "error.missing.expression"}
	ERROR_MISSING_SELECT_CLAUSE                                  = diagnosticErrorCode{diagnosticId: "BCE0503", messageKey: "error.missing.select.clause"}
	ERROR_MISSING_RECEIVE_FIELD_IN_RECEIVE_ACTION                = diagnosticErrorCode{diagnosticId: "BCE0504", messageKey: "error.missing.receive.field.in.receive.action"}
	ERROR_MISSING_WAIT_FIELD_IN_WAIT_ACTION                      = diagnosticErrorCode{diagnosticId: "BCE0505", messageKey: "error.missing.wait.field.in.wait.action"}
	ERROR_MISSING_WAIT_FUTURE_EXPRESSION                         = diagnosticErrorCode{diagnosticId: "BCE0506", messageKey: "error.missing.wait.future.expression"}
	ERROR_MISSING_ENUM_MEMBER                                    = diagnosticErrorCode{diagnosticId: "BCE0507", messageKey: "error.missing.enum.member"}
	ERROR_MISSING_XML_ATOMIC_NAME_PATTERN                        = diagnosticErrorCode{diagnosticId: "BCE0508", messageKey: "error.missing.xml.atomic.name.pattern"}
	ERROR_MISSING_TUPLE_MEMBER                                   = diagnosticErrorCode{diagnosticId: "BCE0509", messageKey: "error.missing.tuple.member"}
	ERROR_MISSING_ORDER_KEY                                      = diagnosticErrorCode{diagnosticId: "BCE0510", messageKey: "error.missing.order.key"}
	ERROR_MISSING_ANNOTATION_ATTACH_POINT                        = diagnosticErrorCode{diagnosticId: "BCE0511", messageKey: "error.missing.annotation.attach.point"}
	ERROR_MISSING_LET_VARIABLE_DECLARATION                       = diagnosticErrorCode{diagnosticId: "BCE0512", messageKey: "error.missing.let.variable.declaration"}
	ERROR_MISSING_NAMED_WORKER_DECLARATION_IN_FORK_STMT          = diagnosticErrorCode{diagnosticId: "BCE0513", messageKey: "error.missing.named.worker.declaration.in.fork.stmt"}
	ERROR_MISSING_KEY_EXPR_IN_MEMBER_ACCESS_EXPR                 = diagnosticErrorCode{diagnosticId: "BCE0514", messageKey: "error.missing.key.expr.in.member.access.expr"}
	ERROR_MISSING_ERROR_MESSAGE_BINDING_PATTERN                  = diagnosticErrorCode{diagnosticId: "BCE0515", messageKey: "error.missing.error.message.binding.pattern"}
	ERROR_CONFIGURABLE_VARIABLE_MUST_BE_INITIALIZED_OR_REQUIRED  = diagnosticErrorCode{diagnosticId: "BCE0516", messageKey: "error.configurable.variable.must.be.initialized.or.required"}
	ERROR_MISSING_RESOURCE_PATH_IN_RESOURCE_ACCESSOR_DEFINITION  = diagnosticErrorCode{diagnosticId: "BCE0517", messageKey: "error.missing.resource.path.in.resource.accessor.definition"}
	ERROR_MISSING_RESOURCE_PATH_IN_RESOURCE_ACCESSOR_DECLARATION = diagnosticErrorCode{diagnosticId: "BCE0518", messageKey: "error.missing.resource.path.in.resource.accessor.declaration"}
	ERROR_MISSING_ERROR_MESSAGE_IN_ERROR_CONSTRUCTOR             = diagnosticErrorCode{diagnosticId: "BCE0519", messageKey: "error.missing.error.message.in.error.constructor"}
	ERROR_MISSING_ARG_WITHIN_PARENTHESIS                         = diagnosticErrorCode{diagnosticId: "BCE0520", messageKey: "error.missing.arg.within.parenthesis"}
	ERROR_MISSING_VARIABLE_NAME                                  = diagnosticErrorCode{diagnosticId: "BCE0521", messageKey: "error.missing.variable.name"}
	ERROR_MISSING_FIELD_NAME                                     = diagnosticErrorCode{diagnosticId: "BCE0522", messageKey: "error.missing.field.name"}
	ERROR_MISSING_BUILTIN_TYPE                                   = diagnosticErrorCode{diagnosticId: "BCE0523", messageKey: "error.missing.builtin.type"}
	ERROR_ANNOTATION_NOT_ATTACHED_TO_A_CONSTRUCT                 = diagnosticErrorCode{diagnosticId: "BCE0524", messageKey: "error.annotation.not.attached.to.a.construct"}
	ERROR_DOCUMENTATION_NOT_ATTACHED_TO_A_CONSTRUCT              = diagnosticErrorCode{diagnosticId: "BCE0525", messageKey: "error.documentation.not.attached.to.a.construct"}
	ERROR_MISSING_MATCH_PATTERN                                  = diagnosticErrorCode{diagnosticId: "BCE0526", messageKey: "error.missing.match.pattern"}
	ERROR_MISSING_TYPE_REFERENCE                                 = diagnosticErrorCode{diagnosticId: "BCE0527", messageKey: "error.missing.type.reference"}
	ERROR_MISSING_BACKTICK_STRING                                = diagnosticErrorCode{diagnosticId: "BCE0528", messageKey: "error.missing.backtick.string"}
	ERROR_MISSING_NAMED_ARG                                      = diagnosticErrorCode{diagnosticId: "BCE0529", messageKey: "error.missing.named.arg"}
	ERROR_MISSING_FIELD_MATCH_PATTERN_MEMBER                     = diagnosticErrorCode{diagnosticId: "BCE0530", messageKey: "error.missing.field.match.pattern.member"}
	ERROR_MISSING_OBJECT_CONSTRUCTOR_EXPRESSION                  = diagnosticErrorCode{diagnosticId: "BCE0531", messageKey: "error.missing.object.constructor.expression"}
	ERROR_MISSING_GROUPING_KEY                                   = diagnosticErrorCode{diagnosticId: "BCE0532", messageKey: "error.missing.grouping.key"}
	ERROR_MISSING_NATURAL_PROMPT_BLOCK                           = diagnosticErrorCode{diagnosticId: "BCE0533", messageKey: "error.missing.natural.prompt.block"}

	// Invalid nodes
	ERROR_INVALID_TOKEN                                             = diagnosticErrorCode{diagnosticId: "BCE0600", messageKey: "error.invalid.token"}
	ERROR_EXPRESSION_EXPECTED_ACTION_FOUND                          = diagnosticErrorCode{diagnosticId: "BCE0601", messageKey: "error.expression.expected.action.found"}
	ERROR_ONLY_TYPE_REFERENCE_ALLOWED_AS_TYPE_INCLUSIONS            = diagnosticErrorCode{diagnosticId: "BCE0602", messageKey: "error.only.type.reference.allowed.as.type.inclusions"}
	ERROR_NAMED_WORKER_NOT_ALLOWED_HERE                             = diagnosticErrorCode{diagnosticId: "BCE0603", messageKey: "error.named.worker.not.allowed.here"}
	ERROR_ONLY_NAMED_WORKERS_ALLOWED_HERE                           = diagnosticErrorCode{diagnosticId: "BCE0604", messageKey: "error.only.named.workers.allowed.here"}
	ERROR_IMPORT_DECLARATION_AFTER_OTHER_DECLARATIONS               = diagnosticErrorCode{diagnosticId: "BCE0605", messageKey: "error.import.declaration.after.other.declarations"}
	ERROR_ANNOTATIONS_ATTACHED_TO_EXPRESSION                        = diagnosticErrorCode{diagnosticId: "BCE0606", messageKey: "error.annotations.attached.to.expression"}
	ERROR_INVALID_EXPRESSION_IN_START_ACTION                        = diagnosticErrorCode{diagnosticId: "BCE0607", messageKey: "error.invalid.expression.in.start.action"}
	ERROR_DUPLICATE_QUALIFIER                                       = diagnosticErrorCode{diagnosticId: "BCE0608", messageKey: "error.duplicate.qualifier"}
	ERROR_QUALIFIER_NOT_ALLOWED                                     = diagnosticErrorCode{diagnosticId: "BCE0609", messageKey: "error.qualifier.not.allowed"}
	ERROR_TYPE_INCLUSION_IN_OBJECT_CONSTRUCTOR                      = diagnosticErrorCode{diagnosticId: "BCE0610", messageKey: "error.type.inclusion.in.object.constructor"}
	ERROR_MAPPING_CONSTRUCTOR_EXPR_AS_A_WAIT_EXPR                   = diagnosticErrorCode{diagnosticId: "BCE0611", messageKey: "error.mapping.constructor.expr.as.a.wait.expr"}
	ERROR_INVALID_PARAM_LIST_IN_INFER_ANONYMOUS_FUNCTION_EXPR       = diagnosticErrorCode{diagnosticId: "BCE0612", messageKey: "error.invalid.param.list.in.infer.anonymous.function.expr"}
	ERROR_MORE_RECORD_FIELDS_AFTER_REST_FIELD                       = diagnosticErrorCode{diagnosticId: "BCE0613", messageKey: "error.more.record.fields.after.rest.field"}
	ERROR_INVALID_XML_NAMESPACE_URI                                 = diagnosticErrorCode{diagnosticId: "BCE0614", messageKey: "error.invalid.xml.namespace.uri"}
	ERROR_INTERPOLATION_IS_NOT_ALLOWED_FOR_XML_TAG_NAMES            = diagnosticErrorCode{diagnosticId: "BCE0615", messageKey: "error.interpolation.is.not.allowed.for.xml.tag.names"}
	ERROR_INTERPOLATION_IS_NOT_ALLOWED_WITHIN_ELEMENT_TAGS          = diagnosticErrorCode{diagnosticId: "BCE0616", messageKey: "error.interpolation.is.not.allowed.within.element.tags"}
	ERROR_INTERPOLATION_IS_NOT_ALLOWED_WITHIN_XML_COMMENTS          = diagnosticErrorCode{diagnosticId: "BCE0617", messageKey: "error.interpolation.is.not.allowed.within.xml.comments"}
	ERROR_INTERPOLATION_IS_NOT_ALLOWED_WITHIN_XML_PI                = diagnosticErrorCode{diagnosticId: "BCE0618", messageKey: "error.interpolation.is.not.allowed.within.xml.pi"}
	ERROR_INVALID_EXPR_IN_ASSIGNMENT_LHS                            = diagnosticErrorCode{diagnosticId: "BCE0619", messageKey: "error.invalid.expr.in.assignment.lhs"}
	ERROR_INVALID_EXPR_IN_COMPOUND_ASSIGNMENT_LHS                   = diagnosticErrorCode{diagnosticId: "BCE0620", messageKey: "error.invalid.expr.in.compound.assignment.lhs"}
	ERROR_INVALID_METADATA                                          = diagnosticErrorCode{diagnosticId: "BCE0621", messageKey: "error.invalid.metadata"}
	ERROR_INVALID_QUALIFIER                                         = diagnosticErrorCode{diagnosticId: "BCE0622", messageKey: "error.invalid.qualifier"}
	ERROR_ANNOTATIONS_ATTACHED_TO_STATEMENT                         = diagnosticErrorCode{diagnosticId: "BCE0623", messageKey: "error.annotations.attached.to.statement"}
	ERROR_ACTION_AS_A_WAIT_EXPR                                     = diagnosticErrorCode{diagnosticId: "BCE0625", messageKey: "error.action.as.a.wait.expr"}
	ERROR_INVALID_USAGE_OF_VAR                                      = diagnosticErrorCode{diagnosticId: "BCE0626", messageKey: "error.invalid.usage.of.var"}
	ERROR_MATCH_PATTERN_AFTER_REST_MATCH_PATTERN                    = diagnosticErrorCode{diagnosticId: "BCE0627", messageKey: "error.match.pattern.after.rest.match.pattern"}
	ERROR_MATCH_PATTERN_NOT_ALLOWED                                 = diagnosticErrorCode{diagnosticId: "BCE0628", messageKey: "error.match.pattern.not.allowed"}
	ERROR_MATCH_STATEMENT_SHOULD_HAVE_ONE_OR_MORE_MATCH_CLAUSES     = diagnosticErrorCode{diagnosticId: "BCE0629", messageKey: "error.match.statement.should.have.one.or.more.match.clauses"}
	ERROR_PARAMETER_AFTER_THE_REST_PARAMETER                        = diagnosticErrorCode{diagnosticId: "BCE0630", messageKey: "error.parameter.after.the.rest.parameter"}
	ERROR_REQUIRED_PARAMETER_AFTER_THE_DEFAULTABLE_PARAMETER        = diagnosticErrorCode{diagnosticId: "BCE0631", messageKey: "error.required.parameter.after.the.defaultable.parameter"}
	ERROR_NAMED_ARG_FOLLOWED_BY_POSITIONAL_ARG                      = diagnosticErrorCode{diagnosticId: "BCE0632", messageKey: "error.named.arg.followed.by.positional.arg"}
	ERROR_REST_ARG_FOLLOWED_BY_ANOTHER_ARG                          = diagnosticErrorCode{diagnosticId: "BCE0633", messageKey: "error.rest.arg.followed.by.another.arg"}
	ERROR_BINDING_PATTERN_NOT_ALLOWED                               = diagnosticErrorCode{diagnosticId: "BCE0634", messageKey: "error.binding.pattern.not.allowed"}
	ERROR_INVALID_BASE16_CONTENT_IN_BYTE_ARRAY_LITERAL              = diagnosticErrorCode{diagnosticId: "BCE0635", messageKey: "error.invalid.base16.content.in.byte.array.literal"}
	ERROR_INVALID_BASE64_CONTENT_IN_BYTE_ARRAY_LITERAL              = diagnosticErrorCode{diagnosticId: "BCE0636", messageKey: "error.invalid.base64.content.in.byte.array.literal"}
	ERROR_INVALID_CONTENT_IN_BYTE_ARRAY_LITERAL                     = diagnosticErrorCode{diagnosticId: "BCE0637", messageKey: "error.invalid.content.in.byte.array.literal"}
	ERROR_INVALID_EXPRESSION_STATEMENT                              = diagnosticErrorCode{diagnosticId: "BCE0638", messageKey: "error.invalid.expression.statement"}
	ERROR_INVALID_ARRAY_LENGTH                                      = diagnosticErrorCode{diagnosticId: "BCE0639", messageKey: "error.invalid.array.length"}
	ERROR_SELECT_CLAUSE_IN_QUERY_ACTION                             = diagnosticErrorCode{diagnosticId: "BCE0640", messageKey: "error.select.clause.in.query.action"}
	ERROR_MORE_CLAUSES_AFTER_SELECT_CLAUSE                          = diagnosticErrorCode{diagnosticId: "BCE0641", messageKey: "error.more.clauses.after.select.clause"}
	ERROR_QUERY_CONSTRUCT_TYPE_IN_QUERY_ACTION                      = diagnosticErrorCode{diagnosticId: "BCE0642", messageKey: "error.query.construct.type.in.query.action"}
	ERROR_NO_WHITESPACES_ALLOWED_IN_RIGHT_SHIFT_OP                  = diagnosticErrorCode{diagnosticId: "BCE0643", messageKey: "error.no.whitespaces.allowed.in.right.shift.op"}
	ERROR_NO_WHITESPACES_ALLOWED_IN_UNSIGNED_RIGHT_SHIFT_OP         = diagnosticErrorCode{diagnosticId: "BCE0644", messageKey: "error.no.whitespaces.allowed.in.unsigned.right.shift.op"}
	ERROR_INVALID_WHITESPACE_IN_SLASH_LT_TOKEN                      = diagnosticErrorCode{diagnosticId: "BCE0645", messageKey: "error.invalid.whitespace.in.slash.lt.token"}
	ERROR_LOCAL_TYPE_DEFINITION_NOT_ALLOWED                         = diagnosticErrorCode{diagnosticId: "BCE0646", messageKey: "error.local.type.definition.not.allowed"}
	ERROR_LEADING_ZEROS_IN_NUMERIC_LITERALS                         = diagnosticErrorCode{diagnosticId: "BCE0647", messageKey: "error.leading.zeros.in.numeric.literals"}
	ERROR_INVALID_STRING_NUMERIC_ESCAPE_SEQUENCE                    = diagnosticErrorCode{diagnosticId: "BCE0648", messageKey: "error.invalid.string.numeric.escape.sequence"}
	ERROR_INVALID_ESCAPE_SEQUENCE                                   = diagnosticErrorCode{diagnosticId: "BCE0649", messageKey: "error.invalid.escape.sequence"}
	ERROR_INVALID_WHITESPACE_BEFORE                                 = diagnosticErrorCode{diagnosticId: "BCE0650", messageKey: "error.invalid.whitespace.before"}
	ERROR_INVALID_WHITESPACE_AFTER                                  = diagnosticErrorCode{diagnosticId: "BCE0651", messageKey: "error.invalid.whitespace.after"}
	ERROR_INVALID_XML_NAME                                          = diagnosticErrorCode{diagnosticId: "BCE0652", messageKey: "error.invalid.xml.name"}
	ERROR_INVALID_CHARACTER_IN_XML_ATTRIBUTE_VALUE                  = diagnosticErrorCode{diagnosticId: "BCE0653", messageKey: "error.invalid.character.in.xml.attribute.value"}
	ERROR_INVALID_ENTITY_REFERENCE_NAME_START                       = diagnosticErrorCode{diagnosticId: "BCE0654", messageKey: "error.invalid.entity.reference.name.start"}
	ERROR_DOUBLE_HYPHEN_NOT_ALLOWED_WITHIN_XML_COMMENT              = diagnosticErrorCode{diagnosticId: "BCE0655", messageKey: "error.double.hyphen.not.allowed.within.xml.comment"}
	ERROR_MORE_THAN_ONE_OBJECT_NETWORK_QUALIFIERS                   = diagnosticErrorCode{diagnosticId: "BCE0657", messageKey: "error.more.than.one.object.network.qualifiers"}
	ERROR_REMOTE_METHOD_HAS_A_VISIBILITY_QUALIFIER                  = diagnosticErrorCode{diagnosticId: "BCE0658", messageKey: "error.remote.method.has.a.visibility.qualifier"}
	ERROR_PRIVATE_QUALIFIER_IN_OBJECT_MEMBER_DESCRIPTOR             = diagnosticErrorCode{diagnosticId: "BCE0659", messageKey: "error.private.qualifier.in.object.member.descriptor"}
	ERROR_RESOURCE_PATH_IN_FUNCTION_DEFINITION                      = diagnosticErrorCode{diagnosticId: "BCE0660", messageKey: "error.resource.path.in.function.definition"}
	ERROR_RESOURCE_PATH_SEGMENT_NOT_ALLOWED_AFTER_REST_PARAM        = diagnosticErrorCode{diagnosticId: "BCE0661", messageKey: "error.resource.path.segment.not.allowed.after.rest.param"}
	ERROR_REST_ARG_IN_ERROR_CONSTRUCTOR                             = diagnosticErrorCode{diagnosticId: "BCE0662", messageKey: "error.rest.arg.in.error.constructor"}
	ERROR_ADDITIONAL_POSITIONAL_ARG_IN_ERROR_CONSTRUCTOR            = diagnosticErrorCode{diagnosticId: "BCE0663", messageKey: "error.additional.positional.arg.in.error.constructor"}
	ERROR_DEFAULTABLE_PARAMETER_CANNOT_BE_INCLUDED_RECORD_PARAMETER = diagnosticErrorCode{diagnosticId: "BCE0664", messageKey: "error.defaultable.parameter.cannot.be.included.record.parameter"}
	ERROR_INCOMPLETE_QUOTED_IDENTIFIER                              = diagnosticErrorCode{diagnosticId: "BCE0665", messageKey: "error.incomplete.quoted.identifier"}
	ERROR_INCLUSIVE_RECORD_TYPE_CANNOT_CONTAIN_REST_FIELD           = diagnosticErrorCode{diagnosticId: "BCE0666", messageKey: "error.inclusive.record.type.cannot.contain.rest.field"}
	ERROR_VARIABLE_DECL_HAVING_BP_MUST_BE_INITIALIZED               = diagnosticErrorCode{diagnosticId: "BCE0667", messageKey: "error.variable.decl.having.bp.must.be.initialized"}
	ERROR_ISOLATED_VAR_CANNOT_BE_DECLARED_AS_PUBLIC                 = diagnosticErrorCode{diagnosticId: "BCE0668", messageKey: "error.isolated.var.cannot.be.declared.as.public"}
	ERROR_VARIABLE_DECLARED_WITH_VAR_CANNOT_BE_PUBLIC               = diagnosticErrorCode{diagnosticId: "BCE0669", messageKey: "error.variable.declared.with.var.cannot.be.public"}
	ERROR_FIELD_BP_INSIDE_LIST_BP                                   = diagnosticErrorCode{diagnosticId: "BCE0670", messageKey: "error.field.binding.pattern.inside.list.binding.pattern"}
	ERROR_INVALID_EXPRESSION_EXPECTED_CALL_EXPRESSION               = diagnosticErrorCode{diagnosticId: "BCE0671", messageKey: "error.invalid.expression.expected.a.call.expression"}
	ERROR_TYPE_DESC_AFTER_REST_DESCRIPTOR                           = diagnosticErrorCode{diagnosticId: "BCE0672", messageKey: "error.type.desc.after.rest.descriptor"}
	ERROR_CONFIGURABLE_VAR_IMPLICITLY_FINAL                         = diagnosticErrorCode{diagnosticId: "BCE0673", messageKey: "error.configurable.var.implicitly.final"}
	ERROR_LOCAL_CONST_DECL_NOT_ALLOWED                              = diagnosticErrorCode{diagnosticId: "BCE0674", messageKey: "error.local.const.decl.not.allowed"}
	ERROR_FIELD_INITIALIZATION_NOT_ALLOWED_IN_OBJECT_TYPE           = diagnosticErrorCode{diagnosticId: "BCE0675", messageKey: "error.field.initialization.not.allowed.in.object.type"}
	ERROR_INTERVENING_WHITESPACES_ARE_NOT_ALLOWED                   = diagnosticErrorCode{diagnosticId: "BCE0676", messageKey: "error.intervening.whitespaces.are.not.allowed"}
	ERROR_INVALID_BINDING_PATTERN                                   = diagnosticErrorCode{diagnosticId: "BCE0677", messageKey: "error.invalid.binding.pattern"}
	ERROR_RESOURCE_PATH_CANNOT_BEGIN_WITH_SLASH                     = diagnosticErrorCode{diagnosticId: "BCE0678", messageKey: "error.resource.path.cannot.begin.with.slash"}
	REST_PARAMETER_CANNOT_BE_INCLUDED_RECORD_PARAMETER              = diagnosticErrorCode{diagnosticId: "BCE0679", messageKey: "error.rest.parameter.cannot.be.included.record.parameter"}
	RESOURCE_ACCESS_SEGMENT_IS_NOT_ALLOWED_AFTER_REST_SEGMENT       = diagnosticErrorCode{diagnosticId: "BCE0680", messageKey: "error.resource.access.segment.is.not.allowed.after.rest.segment"}
	ERROR_INVALID_TOKEN_IN_REG_EXP                                  = diagnosticErrorCode{diagnosticId: "BCE0681", messageKey: "error.invalid.token.in.reg.exp"}
	ERROR_INVALID_FLAG_IN_REG_EXP                                   = diagnosticErrorCode{diagnosticId: "BCE0682", messageKey: "error.invalid.flag.in.reg.exp"}
	ERROR_INVALID_QUANTIFIER_IN_REG_EXP                             = diagnosticErrorCode{diagnosticId: "BCE0683", messageKey: "error.invalid.quantifier.in.reg.exp"}
	ERROR_ANNOTATIONS_NOT_ALLOWED_FOR_TUPLE_REST_DESCRIPTOR         = diagnosticErrorCode{diagnosticId: "BCE0684", messageKey: "error.annotations.not.allowed.for.tuple.rest.descriptor"}
	ERROR_INVALID_RE_SYNTAX_CHAR                                    = diagnosticErrorCode{diagnosticId: "BCE0685", messageKey: "error.invalid.syntax.char"}
	ERROR_MORE_CLAUSES_AFTER_COLLECT_CLAUSE                         = diagnosticErrorCode{diagnosticId: "BCE0686", messageKey: "error.more.clauses.after.collect.clause"}
	ERROR_COLLECT_CLAUSE_IN_QUERY_ACTION                            = diagnosticErrorCode{diagnosticId: "BCE0687", messageKey: "error.collect.clause.in.query.action"}
)

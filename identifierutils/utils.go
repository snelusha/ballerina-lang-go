package identifierutils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	charPrefix            = "&"
	escapePrefix          = "\\"
	generatedMethodPrefix = "$gen$"
)

var (
	unicodePattern              = regexp.MustCompile(`\\(\\*)u\{([a-fA-F0-9]+)\}`)
	unescapedSpecialCharPattern = regexp.MustCompile(`([$&+,:;=\?@#\\|/'\\ \[\}\]<>."^*{}~` + "`" + `()%!-])`)
)

type identifier struct {
	name      string
	isEncoded bool
}

func encodeSpecialCharacters(id string) string {
	var sb strings.Builder
	index := 0

	for index < len(id) {
		if id[index] == '\\' && index+1 < len(id) {
			if formatted := getFormattedStringForQuotedIdentifiers(rune(id[index+1])); formatted != "" {
				unicodePoint := charPrefix + formatted
				sb.WriteString(unicodePoint)
				index = index + 2
				continue
			}
			index = index + 1
		}
		sb.WriteByte(id[index])
		index = index + 1
	}
	return sb.String()
}

func EscapeSpecialCharacters(id string) string {
	return unescapedSpecialCharPattern.ReplaceAllString(id, "\\$1")
}

func encodeIdentifier(id string) string {
	if strings.Contains(id, escapePrefix) {
		id = encodeSpecialCharacters(id)
		return UnescapeJava(id)
	}
	return id
}

func UnescapeJava(str string) string {
	if str == "" {
		return str
	}

	var sb strings.Builder
	index := 0

	for index < len(str) {
		if str[index] == '\\' && index+1 < len(str) {
			next := str[index+1]
			switch next {
			case 'n':
				sb.WriteByte('\n')
				index = index + 2
			case 't':
				sb.WriteByte('\t')
				index = index + 2
			case 'r':
				sb.WriteByte('\r')
				index = index + 2
			case 'b':
				sb.WriteByte('\b')
				index = index + 2
			case 'f':
				sb.WriteByte('\f')
				index = index + 2
			case '\\':
				sb.WriteByte('\\')
				index = index + 2
			case '"':
				sb.WriteByte('"')
				index = index + 2
			case '\'':
				sb.WriteByte('\'')
				index = index + 2
			case 'u':
				if index+5 < len(str) {
					hexStr := str[index+2 : index+6]
					if codePoint, err := strconv.ParseInt(hexStr, 16, 32); err == nil {
						sb.WriteRune(rune(codePoint))
						index = index + 6
						continue
					}
				}
				sb.WriteByte(str[index])
				index = index + 1
			default:
				sb.WriteByte(str[index])
				index = index + 1
			}
		} else {
			sb.WriteByte(str[index])
			index = index + 1
		}
	}

	return sb.String()
}

func encodeGeneratedName(id string) identifier {
	var sb strings.Builder

	isEncoded := false

	for _, ch := range id {
		if formatted := getFormattedStringForJvmReservedSet(ch); formatted != "" {
			unicodePoint := charPrefix + formatted
			sb.WriteString(unicodePoint)
			isEncoded = true
		} else {
			sb.WriteRune(ch)
		}
	}
	return identifier{name: sb.String(), isEncoded: isEncoded}
}

func getFormattedStringForQuotedIdentifiers(c rune) string {
	if c == '$' {
		return "0036"
	}

	return getFormattedStringForJvmReservedSet(c)
}

func getFormattedStringForJvmReservedSet(c rune) string {
	switch c {
	case '\\':
		return "0092"
	case '.':
		return "0046"
	case ':':
		return "0058"
	case ';':
		return "0059"
	case '[':
		return "0091"
	case ']':
		return "0093"
	case '/':
		return "0047"
	case '<':
		return "0060"
	case '>':
		return "0062"
	default:
		return ""
	}
}

func DecodeIdentifier(encodedId string) string {
	if encodedId == "" {
		return ""
	}

	var sb strings.Builder
	index := 0

	for index < len(encodedId) {
		if encodedId[index] == '&' && index+4 < len(encodedId) {
			if isUnicodePoint(encodedId, index) {
				codePoint, _ := strconv.ParseInt(encodedId[index+1:index+5], 10, 32)
				sb.WriteRune(rune(codePoint))
				index = index + 5
			} else {
				sb.WriteByte(encodedId[index])
				index = index + 1
			}
		} else {
			sb.WriteByte(encodedId[index])
			index++
		}
	}

	return decodeGeneratedMethodName(sb.String())
}

func decodeGeneratedMethodName(decodedName string) string {
	if strings.HasPrefix(decodedName, generatedMethodPrefix) {
		return decodedName[len(generatedMethodPrefix):]
	}
	return decodedName
}

func UnescapeBallerina(text string) string {
	return UnescapeJava(UnescapeUnicodeCodepoints((text)))
}

func UnescapeUnicodeCodepoints(id string) string {
	result := unicodePattern.ReplaceAllStringFunc(id, func(match string) string {
		submatch := unicodePattern.FindStringSubmatch(match)
		if len(submatch) < 3 {
			return match
		}

		leadingSlashes := submatch[1]
		if IsEscapedNumericEscape(leadingSlashes) {
			return match
		}

		codePoint, err := strconv.ParseInt(submatch[2], 16, 32)
		if err != nil {
			return match
		}

		ch := string(rune(codePoint))

		if ch == "\\" {
			return leadingSlashes + "\\u005C"
		}

		return leadingSlashes + ch
	})

	return result
}

func IsEscapedNumericEscape(leadingSlashes string) bool {
	return !isEven(len(leadingSlashes))
}

func isEven(n int) bool {
	return (n & 1) == 0
}

func isUnicodePoint(encodedName string, index int) bool {
	if index+5 > len(encodedName) {
		return false
	}
	return containsOnlyDigits(encodedName[index+1 : index+5])
}

func containsOnlyDigits(digitString string) bool {
	for _, ch := range digitString {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

func EncodeFunctionIdentifier(functionName string) string {
	functionName = encodeIdentifier(functionName)

	switch functionName {
	case ".<init>":
		return "$gen$&0046&0060init&0062"
	case ".<start>":
		return "$gen$&0046&0060start&0062"
	case ".<stop>":
		return "$gen$&0046&0060stop&0062"
	case ".<testinit>":
		return "$gen$&0046&0060testinit&0062"
	}

	encodedName := encodeGeneratedName(functionName)
	if encodedName.isEncoded {
		return generatedMethodPrefix + encodedName.name
	}
	return functionName
}

func EncodeNonFunctionIdentifier(identifierString string) string {
	identifierString = encodeIdentifier(identifierString)
	encodedName := encodeGeneratedName(identifierString)
	return encodedName.name
}

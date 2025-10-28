package identifierutils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	charPrefix            = "&"
	escapePrefix          = "\\"
	generatedMethodPrefix = "$gen$"
)

var (
	unicodePattern          = regexp.MustCompile(`\\(\\*)u\{([a-fA-F0-9]+)\}`)
	unescapedSpecialCharSet = regexp.MustCompile(`([$&+,:;=\?@#\\|/'\\ \[\}\]<>."^*{}~` + "`" + `()%!-])`)
)

func encodeSpecialCharacters(identifier string) string {
	var sb strings.Builder
	sb.Grow(len(identifier) * 2)

	runes := []rune(identifier)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i+1 < len(runes) {
			if formatted := getFormattedStringForQuotedIdentifiers(runes[i+1]); formatted != "" {
				sb.WriteString(charPrefix)
				sb.WriteString(formatted)
				i++
				continue
			}
		}
		sb.WriteRune(runes[i])
	}
	return sb.String()
}

func EscapeSpecialCharacters(identifier string) string {
	return unescapedSpecialCharSet.ReplaceAllString(identifier, "\\$1")
}

func encodeIdentifier(identifier string) string {
	if strings.Contains(identifier, escapePrefix) {
		identifier = encodeSpecialCharacters(identifier)
		return UnescapeJava(identifier)
	}
	return identifier
}

func UnescapeJava(str string) string {
	if str == "" {
		return str
	}

	var sb strings.Builder
	sb.Grow(len(str))

	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		if runes[i] != '\\' {
			sb.WriteRune(runes[i])
			continue
		}

		if i+1 >= len(runes) {
			sb.WriteRune(runes[i])
			continue
		}

		i++
		ch := runes[i]

		switch ch {
		case 'n':
			sb.WriteRune('\n')
		case 't':
			sb.WriteRune('\t')
		case 'r':
			sb.WriteRune('\r')
		case 'f':
			sb.WriteRune('\f')
		case 'b':
			sb.WriteRune('\b')
		case '\\':
			sb.WriteRune('\\')
		case '\'':
			sb.WriteRune('\'')
		case '"':
			sb.WriteRune('"')
		case 'u':
			if i+4 < len(runes) {
				hexStr := string(runes[i+1 : i+5])
				if codePoint, err := strconv.ParseInt(hexStr, 16, 32); err == nil {
					sb.WriteRune(rune(codePoint))
					i += 4
				} else {
					sb.WriteRune('u')
				}
			} else {
				sb.WriteRune('u')
			}
		default:
			sb.WriteRune(runes[i])
		}
	}

	return sb.String()
}

type identifier struct {
	name      string
	isEncoded bool
}

func encodeGeneratedName(identifierStr string) identifier {
	var sb strings.Builder
	sb.Grow(len(identifierStr) * 2)

	isEncoded := false
	runes := []rune(identifierStr)

	for _, r := range runes {
		if formatted := getFormattedStringForJvmReservedSet(r); formatted != "" {
			sb.WriteString(charPrefix)
			sb.WriteString(formatted)
			isEncoded = true
		} else {
			sb.WriteRune(r)
		}
	}

	return identifier{
		name:      sb.String(),
		isEncoded: isEncoded,
	}
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

func DecodeIdentifier(encodedIdentifier string) string {
	if encodedIdentifier == "" {
		return ""
	}

	var sb strings.Builder
	sb.Grow(len(encodedIdentifier))

	runes := []rune(encodedIdentifier)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '&' && i+4 < len(runes) {
			if isUnicodePoint(runes[i+1 : i+5]) {
				unicodeStr := string(runes[i+1 : i+5])
				if codePoint, err := strconv.ParseInt(unicodeStr, 10, 32); err == nil {
					sb.WriteRune(rune(codePoint))
					i += 4
					continue
				}
			}
		}
		sb.WriteRune(runes[i])
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
	return UnescapeJava(UnescapeUnicodeCodepoints(text))
}

func UnescapeUnicodeCodepoints(identifierStr string) string {
	result := unicodePattern.ReplaceAllStringFunc(identifierStr, func(match string) string {
		submatches := unicodePattern.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		leadingSlashes := submatches[1]
		if IsEscapedNumericEscape(leadingSlashes) {
			return match
		}

		codePointHex := submatches[2]
		codePoint, err := strconv.ParseInt(codePointHex, 16, 32)
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
	return !isEven(utf8.RuneCountInString(leadingSlashes))
}

func isEven(n int) bool {
	return (n & 1) == 0
}

func isUnicodePoint(runes []rune) bool {
	if len(runes) != 4 {
		return false
	}
	for _, r := range runes {
		if !unicode.IsDigit(r) {
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

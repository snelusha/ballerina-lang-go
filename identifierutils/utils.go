package identifierutils

import (
	"regexp"
	"strings"
)

const (
	unicodeRegex          = `\\(\\*)u\{([a-fA-F0-9]+)\}`
	charPrefix            = "&"
	escapePrefix          = "\\"
	generatedMethodPrefix = "$gen$"
)

var (
	unicodePattern          = regexp.MustCompile(unicodeRegex)
	unescapedSpecialCharSet = regexp.MustCompile(`([$&+,:;=\?@#\\|/'\\ \[\}\]<>.\"^*{}~` + "`" + `()%!-])`)
)

type identifier struct {
	Name      string
	IsEncoded bool
}

func encodeSpecialCharacters(ident string) string {
	var sb strings.Builder

	index := 0

	for index < len(ident) {
		if ident[index] == '\\' && (index+1 < len(ident)) {
			formattedStr := getFormattedStringForQuotedIdentifiers(rune(ident[index+1]))
			if formattedStr != "" {
				unicodePoint := charPrefix + formattedStr
				sb.WriteString(unicodePoint)
				index = index + 2
				continue
			}
		}
		sb.WriteByte(ident[index])
		index = index + 1
	}

	return sb.String()
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

/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package identifier

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	// pattern: \(\*)u\{([a-fA-F0-9]+)\}
	unicodePattern = regexp.MustCompile(`\\(\\*)u\{([a-fA-F0-9]+)\}`)
)

const (
	charPrefix            = "&"
	escapePrefix          = "\\"
	generatedMethodPrefix = "$gen$"
)

// list of special characters to escape in identifiers (from Java version)
var specialChars = func() string {
	// include the backtick using hex code \x60
	return "$&+,:;=?@#\\|/' []<>.\"^*{}~\x60()%!-" 
}()

type Identifier struct {
	Name     string
	IsEncoded bool
}

func isEven(n int) bool { return (n & 1) == 0 }

func isEscapedNumericEscape(leadingSlashes string) bool {
	return !isEven(len(leadingSlashes))
}

// EscapeSpecialCharacters escapes a set of special characters by prefixing with a backslash.
func EscapeSpecialCharacters(identifier string) string {
	var b strings.Builder
	for _, r := range identifier {
		if strings.ContainsRune(specialChars, r) {
			b.WriteRune('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// UnescapeJava attempts to unescape common Java string escapes.
// Uses strconv.Unquote on a quoted version of the string as a best-effort.
func UnescapeJava(str string) string {
	if str == "" {
		return str
	}
	// Try strconv.Unquote by wrapping with double quotes.
	if unq, err := strconv.Unquote("\"" + str + "\""); err == nil {
		return unq
	}
	// Fallback: replace common escape sequences.
	replacer := strings.NewReplacer(
		"\\b", "\b",
		"\\t", "\t",
		"\\n", "\n",
		"\\f", "\f",
		"\\r", "\r",
		"\\\"", "\"",
		"\\'", "'",
		"\\\\", "\\",
	)
	return replacer.Replace(str)
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

func getFormattedStringForQuotedIdentifiers(c rune) string {
	if c == '$' {
		return "0036"
	}
	return getFormattedStringForJvmReservedSet(c)
}

func encodeSpecialCharacters(identifier string) string {
	var b strings.Builder
	runes := []rune(identifier)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i+1 < len(runes) {
			formatted := getFormattedStringForQuotedIdentifiers(runes[i+1])
			if formatted != "" {
				b.WriteString(charPrefix)
				b.WriteString(formatted)
				i++
				continue
			}
		}
		b.WriteRune(runes[i])
	}
	return b.String()
}

func encodeIdentifier(identifier string) string {
	if strings.Contains(identifier, escapePrefix) {
		identifier = encodeSpecialCharacters(identifier)
		return UnescapeJava(identifier)
	}
	return identifier
}

func encodeGeneratedName(identifier string) Identifier {
	var b strings.Builder
	isEncoded := false
	for _, r := range identifier {
		formatted := getFormattedStringForJvmReservedSet(r)
		if formatted != "" {
			b.WriteString(charPrefix)
			b.WriteString(formatted)
			isEncoded = true
		} else {
			b.WriteRune(r)
		}
	}
	return Identifier{Name: b.String(), IsEncoded: isEncoded}
}

func DecodeGeneratedMethodName(decodedName string) string {
	if strings.HasPrefix(decodedName, generatedMethodPrefix) {
		return decodedName[len(generatedMethodPrefix):]
	}
	return decodedName
}

func containsOnlyDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func isUnicodePoint(encodedName string, index int) bool {
	if index+5 > len(encodedName) {
		return false
	}
	return containsOnlyDigits(encodedName[index+1:index+5])
}

// DecodeIdentifier decodes encoded names (e.g. &0046 -> '.') and strips generated prefix.
func DecodeIdentifier(encodedIdentifier string) string {
	if encodedIdentifier == "" {
		return encodedIdentifier
	}
	var b strings.Builder
	i := 0
	for i < len(encodedIdentifier) {
		if encodedIdentifier[i] == '&' && i+4 < len(encodedIdentifier) {
			if isUnicodePoint(encodedIdentifier, i) {
				numStr := encodedIdentifier[i+1 : i+5]
				n, err := strconv.Atoi(numStr)
				if err == nil {
					b.WriteRune(rune(n))
					i += 5
					continue
				}
			}
			b.WriteByte(encodedIdentifier[i])
			i++
		} else {
			b.WriteByte(encodedIdentifier[i])
			i++
		}
	}
	return DecodeGeneratedMethodName(b.String())
}

// UnescapeUnicodeCodepoints replaces numeric unicode escapes of the form \u{XXXX} with actual characters.
func UnescapeUnicodeCodepoints(identifier string) string {
	if identifier == "" {
		return identifier
	}
	var out strings.Builder
	lastIndex := 0
	matches := unicodePattern.FindAllStringSubmatchIndex(identifier, -1)
	for _, m := range matches {
		// m: pairs of indices: full match [0,1], group1 [2,3], group2 [4,5] etc.
		fullStart, fullEnd := m[0], m[1]
		grp1Start, grp1End := m[2], m[3]
		grp2Start, grp2End := m[4], m[5]

		// Append text before match
		out.WriteString(identifier[lastIndex:fullStart])

		leadingSlashes := identifier[grp1Start:grp1End]
		if isEscapedNumericEscape(leadingSlashes) {
			// leave as-is
			out.WriteString(identifier[fullStart:fullEnd])
			lastIndex = fullEnd
			continue
		}

		hexStr := identifier[grp2Start:grp2End]
		cp, err := strconv.ParseInt(hexStr, 16, 32)
		if err != nil {
			// on error, keep original
			out.WriteString(identifier[fullStart:fullEnd])
			lastIndex = fullEnd
			continue
		}
		r := rune(cp)
		ch := string(r)
		if ch == "\\" {
			// special-case backslash -> replace with \u005C after leading slashes
			out.WriteString(leadingSlashes + "\\u005C")
		} else {
			out.WriteString(leadingSlashes + ch)
		}
		lastIndex = fullEnd
	}
	out.WriteString(identifier[lastIndex:])
	return out.String()
}

func UnescapeBallerina(text string) string {
	return UnescapeJava(UnescapeUnicodeCodepoints(text))
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
	encoded := encodeGeneratedName(functionName)
	if encoded.IsEncoded {
		return generatedMethodPrefix + encoded.Name
	}
	return functionName
}

func EncodeNonFunctionIdentifier(identifierString string) string {
	identifierString = encodeIdentifier(identifierString)
	encoded := encodeGeneratedName(identifierString)
	return encoded.Name
}

// EscapeSpecialCharactersForIdentifier is a small helper used by callers that want the same Java semantics
func EscapeSpecialCharactersForIdentifier(identifier string) string {
	// mirrors Java's escapeSpecialCharacters
	return EscapeSpecialCharacters(identifier)
}

// end of file

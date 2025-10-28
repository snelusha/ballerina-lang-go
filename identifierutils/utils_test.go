package identifierutils

import (
	"strings"
	"testing"
)

func TestEscapeSpecialCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"test$name", "test\\$name"},
		{"test&name", "test\\&name"},
		{"test+name", "test\\+name"},
		{"test,name", "test\\,name"},
		{"test:name", "test\\:name"},
		{"test;name", "test\\;name"},
		{"test=name", "test\\=name"},
		{"test?name", "test\\?name"},
		{"test@name", "test\\@name"},
		{"test#name", "test\\#name"},
		{"test\\name", "test\\\\name"},
		{"test|name", "test\\|name"},
		{"test/name", "test\\/name"},
		{"test'name", "test\\'name"},
		{"test name", "test\\ name"},
		{"test[name", "test\\[name"},
		{"test}name", "test\\}name"},
		{"test]name", "test\\]name"},
		{"test<name", "test\\<name"},
		{"test>name", "test\\>name"},
		{"test.name", "test\\.name"},
		{"test\"name", "test\\\"name"},
		{"test^name", "test\\^name"},
		{"test*name", "test\\*name"},
		{"test{name", "test\\{name"},
		{"test~name", "test\\~name"},
		{"test`name", "test\\`name"},
		{"test(name", "test\\(name"},
		{"test)name", "test\\)name"},
		{"test%name", "test\\%name"},
		{"test!name", "test\\!name"},
		{"test-name", "test\\-name"},
		{"$&+,:;=?@#\\|/' []<>.\"^*{}~`()%!-", "\\$\\&\\+\\,\\:\\;\\=\\?\\@\\#\\\\\\|\\/\\'\\ \\[\\]\\<\\>\\.\\\"\\^\\*\\{\\}\\~\\`\\(\\)\\%\\!\\-"},
		{"", ""},
		{"noSpecialChars", "noSpecialChars"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EscapeSpecialCharacters(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeSpecialCharacters(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnescapeJava(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"test\\nname", "test\nname"},
		{"test\\tname", "test\tname"},
		{"test\\rname", "test\rname"},
		{"test\\\\name", "test\\name"},
		{"test\\'name", "test'name"},
		{"test\\\"name", "test\"name"},
		{"\\u0041", "A"},
		{"\\u0061", "a"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnescapeJava(tt.input)
			if result != tt.expected {
				t.Errorf("UnescapeJava(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDecodeIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"&0036", "$"},
		{"&0092", "\\"},
		{"&0046", "."},
		{"&0058", ":"},
		{"&0059", ";"},
		{"&0091", "["},
		{"&0093", "]"},
		{"&0047", "/"},
		{"&0060", "<"},
		{"&0062", ">"},
		{"test&0046name", "test.name"},
		{"&0060init&0062", "<init>"},
		{"&0046&0060init&0062", ".<init>"},
		{"$gen$test", "test"},
		{"$gen$&0046&0060init&0062", ".<init>"},
		{"$gen$&0046&0060start&0062", ".<start>"},
		{"$gen$&0046&0060stop&0062", ".<stop>"},
		{"$gen$&0046&0060testinit&0062", ".<testinit>"},
		{"regularMethod", "regularMethod"},
		{"&abcd", "&abcd"},
		{"&12ab", "&12ab"},
		{"&", "&"},
		{"&123", "&123"},
		{"&12345", "Ӓ5"},
		{"", ""},
		{"test&0046name&0058value", "test.name:value"},
		{"prefix&0091index&0093", "prefix[index]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DecodeIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnescapeBallerina(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"test\\nname", "test\nname"},
		{"test\\u{0041}", "testA"},
		{"test\\u{61}", "testa"},
		{"test\\u{1F600}", "test😀"},
		{"\\u{0048}ello", "Hello"},
		{"test\\\\u{61}", "test\\u{61}"},
		{"test\\\\\\u{61}", "test\\a"},
		{"\\\\u{61}", "\\u{61}"},
		{"test\\u{5C}", "test\\"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnescapeBallerina(tt.input)
			if result != tt.expected {
				t.Errorf("UnescapeBallerina(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnescapeUnicodeCodepoints(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"\\u{0041}", "A"},
		{"\\u{61}", "a"},
		{"\\u{0048}ello", "Hello"},
		{"test\\u{1F600}emoji", "test😀emoji"},
		{"\\u{1F44D}", "👍"},
		{"\\\\u{61}", "\\\\u{61}"},
		{"\\\\\\u{61}", "\\\\a"},
		{"\\\\\\\\u{61}", "\\\\\\\\u{61}"},
		{"\\u{5C}", "\\u005C"},
		{"\\u{48}\\u{65}\\u{6C}\\u{6C}\\u{6F}", "Hello"},
		{"", ""},
		{"no unicode here", "no unicode here"},
		{"prefix\\u{41}middle\\u{42}suffix", "prefixAmiddleBsuffix"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnescapeUnicodeCodepoints(tt.input)
			if result != tt.expected {
				t.Errorf("UnescapeUnicodeCodepoints(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsEscapedNumericEscape(t *testing.T) {
	tests := []struct {
		leadingSlashes     string
		expectedNotEscaped bool
	}{
		{"", true},
		{"\\", false},
		{"\\\\", true},
		{"\\\\\\", false},
		{"\\\\\\\\", true},
	}

	for _, tt := range tests {
		t.Run(tt.leadingSlashes, func(t *testing.T) {
			result := IsEscapedNumericEscape(tt.leadingSlashes)
			expected := !tt.expectedNotEscaped
			if result != expected {
				t.Errorf("IsEscapedNumericEscape(%q) = %v, want %v", tt.leadingSlashes, result, expected)
			}
		})
	}
}

func TestEncodeFunctionIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{".<init>", "$gen$&0046&0060init&0062"},
		{".<start>", "$gen$&0046&0060start&0062"},
		{".<stop>", "$gen$&0046&0060stop&0062"},
		{".<testinit>", "$gen$&0046&0060testinit&0062"},
		{"normalFunction", "normalFunction"},
		{"myFunction", "myFunction"},
		{"test.method", "$gen$test&0046method"},
		{"test:method", "$gen$test&0058method"},
		{"test;method", "$gen$test&0059method"},
		{"test[method", "$gen$test&0091method"},
		{"test]method", "$gen$test&0093method"},
		{"test/method", "$gen$test&0047method"},
		{"test<method", "$gen$test&0060method"},
		{"test>method", "$gen$test&0062method"},
		{"test\\method", "testmethod"},
		{"test\\$method", "test&0036method"},
		{"test.method:value", "$gen$test&0046method&0058value"},
		{"", ""},
		{"simpleMethod", "simpleMethod"},
		{"method_name", "method_name"},
		{"method123", "method123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EncodeFunctionIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeFunctionIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeNonFunctionIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"myVariable", "myVariable"},
		{"test.field", "test&0046field"},
		{"test:field", "test&0058field"},
		{"test;field", "test&0059field"},
		{"test[field", "test&0091field"},
		{"test]field", "test&0093field"},
		{"test/field", "test&0047field"},
		{"test<field", "test&0060field"},
		{"test>field", "test&0062field"},
		{"test\\field", "test\ffield"},
		{"test\\$field", "test&0036field"},
		{"test.field:value", "test&0046field&0058value"},
		{"", ""},
		{"simpleField", "simpleField"},
		{"field_name", "field_name"},
		{"field123", "field123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EncodeNonFunctionIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeNonFunctionIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRoundTripEncodeDecode(t *testing.T) {
	testCases := []string{
		".<init>",
		".<start>",
		".<stop>",
		".<testinit>",
		"test.method",
		"test:field",
		"simple",
	}

	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			encoded := EncodeFunctionIdentifier(testCase)
			decoded := DecodeIdentifier(encoded)
			if decoded != testCase {
				t.Errorf("Round trip failed for %q: encoded=%q, decoded=%q", testCase, encoded, decoded)
			}
		})
	}
}

func TestRoundTripNonFunctionEncodeDecode(t *testing.T) {
	testCases := []string{
		"test.field",
		"test:value",
		"simple",
		"field[index]",
	}

	for _, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			encoded := EncodeNonFunctionIdentifier(testCase)
			decoded := DecodeIdentifier(encoded)
			if decoded != testCase {
				t.Errorf("Round trip failed for %q: encoded=%q, decoded=%q", testCase, encoded, decoded)
			}
		})
	}
}

func TestComplexUnicodeEscapeSequences(t *testing.T) {
	input := "\\u{48}\\u{65}\\u{6C}\\u{6C}\\u{6F} \\u{1F600}"
	expected := "Hello 😀"
	result := UnescapeBallerina(input)
	if result != expected {
		t.Errorf("UnescapeBallerina(%q) = %q, want %q", input, result, expected)
	}
}

func TestMixedEscapeSequences(t *testing.T) {
	input := "Line1\\nLine2\\u{0009}Tab"
	expected := "Line1\nLine2\tTab"
	result := UnescapeBallerina(input)
	if result != expected {
		t.Errorf("UnescapeBallerina(%q) = %q, want %q", input, result, expected)
	}
}

func TestEscapedBackslashBeforeUnicode(t *testing.T) {
	input := "test\\\\u{61}end"
	result := UnescapeUnicodeCodepoints(input)
	expected := "test\\\\u{61}end"
	if result != expected {
		t.Errorf("UnescapeUnicodeCodepoints(%q) = %q, want %q", input, result, expected)
	}
}

func TestMultipleBackslashesBeforeUnicode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\\u{41}", "A"},
		{"\\\\u{41}", "\\\\u{41}"},
		{"\\\\\\u{41}", "\\\\A"},
		{"\\\\\\\\u{41}", "\\\\\\\\u{41}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := UnescapeUnicodeCodepoints(tt.input)
			if result != tt.expected {
				t.Errorf("UnescapeUnicodeCodepoints(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSpecialCharacterEdgeCases(t *testing.T) {
	input := "$&+,:;=?@#\\|/' []<>.\"^*{}~`()%!-"
	escaped := EscapeSpecialCharacters(input)

	checks := []string{"\\$", "\\&", "\\+", "\\,", "\\:", "\\;"}
	for _, check := range checks {
		if !strings.Contains(escaped, check) {
			t.Errorf("EscapeSpecialCharacters(%q) should contain %q but got %q", input, check, escaped)
		}
	}
}

func TestDecodeWithPartialUnicodePoint(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"&123", "&123"},
		{"&abc", "&abc"},
		{"&12a4", "&12a4"},
		{"test&0", "test&0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DecodeIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeSpecialCharactersInQuotedIdentifier(t *testing.T) {
	input := "test\\$name"
	encoded := EncodeFunctionIdentifier(input)
	if !strings.Contains(encoded, "&0036") {
		t.Errorf("EncodeFunctionIdentifier(%q) = %q, should contain &0036", input, encoded)
	}
}

func TestEmptyAndNullInputs(t *testing.T) {
	if EscapeSpecialCharacters("") != "" {
		t.Error("EscapeSpecialCharacters(\"\") should return \"\"")
	}
	if EncodeFunctionIdentifier("") != "" {
		t.Error("EncodeFunctionIdentifier(\"\") should return \"\"")
	}
	if EncodeNonFunctionIdentifier("") != "" {
		t.Error("EncodeNonFunctionIdentifier(\"\") should return \"\"")
	}
	if DecodeIdentifier("") != "" {
		t.Error("DecodeIdentifier(\"\") should return \"\"")
	}
	if UnescapeBallerina("") != "" {
		t.Error("UnescapeBallerina(\"\") should return \"\"")
	}
	if UnescapeUnicodeCodepoints("") != "" {
		t.Error("UnescapeUnicodeCodepoints(\"\") should return \"\"")
	}
	if UnescapeJava("") != "" {
		t.Error("UnescapeJava(\"\") should return \"\"")
	}
}

func TestHighCodePointUnicodeCharacters(t *testing.T) {
	input := "\\u{1F600}\\u{1F44D}\\u{1F389}"
	expected := "😀👍🎉"
	result := UnescapeBallerina(input)
	if result != expected {
		t.Errorf("UnescapeBallerina(%q) = %q, want %q", input, result, expected)
	}
}

func TestConsecutiveEncodedCharacters(t *testing.T) {
	encoded := "&0046&0058&0059"
	decoded := DecodeIdentifier(encoded)
	expected := ".:;"
	if decoded != expected {
		t.Errorf("DecodeIdentifier(%q) = %q, want %q", encoded, decoded, expected)
	}
}

func TestUnicodeBackslashSpecialHandling(t *testing.T) {
	input := "test\\u{5C}end"
	result := UnescapeUnicodeCodepoints(input)
	expected := "test\\u005Cend"
	if result != expected {
		t.Errorf("UnescapeUnicodeCodepoints(%q) = %q, want %q", input, result, expected)
	}

	finalResult := UnescapeJava(result)
	expectedFinal := "test\\end"
	if finalResult != expectedFinal {
		t.Errorf("UnescapeJava(%q) = %q, want %q", result, finalResult, expectedFinal)
	}
}

func TestEncodingPreservesNonReservedCharacters(t *testing.T) {
	input := "test_method$123"
	encoded := EncodeFunctionIdentifier(input)

	if !strings.Contains(encoded, "_") {
		t.Errorf("EncodeFunctionIdentifier(%q) = %q, should preserve _", input, encoded)
	}
	if !strings.Contains(encoded, "123") {
		t.Errorf("EncodeFunctionIdentifier(%q) = %q, should preserve 123", input, encoded)
	}
}

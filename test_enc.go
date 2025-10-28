package main

import (
	"ballerina-lang-go/identifierutils"
	"fmt"
)

func main() {
	input := "test\\field"
	fmt.Printf("Input: %q (len=%d)\n", input, len(input))
	for i, b := range []byte(input) {
		fmt.Printf("  [%d] = 0x%02x (%c)\n", i, b, b)
	}

	result := identifierutils.EncodeNonFunctionIdentifier(input)
	fmt.Printf("\nResult: %q (len=%d)\n", result, len(result))
	for i, b := range []byte(result) {
		if b >= 32 && b < 127 {
			fmt.Printf("  [%d] = 0x%02x (%c)\n", i, b, b)
		} else {
			fmt.Printf("  [%d] = 0x%02x\n", i, b)
		}
	}

	expected := "test\ffield"
	fmt.Printf("\nExpected: %q (len=%d)\n", expected, len(expected))
	for i, b := range []byte(expected) {
		if b >= 32 && b < 127 {
			fmt.Printf("  [%d] = 0x%02x (%c)\n", i, b, b)
		} else {
			fmt.Printf("  [%d] = 0x%02x\n", i, b)
		}
	}

	if result == expected {
		fmt.Println("\n✓ MATCH")
	} else {
		fmt.Println("\n✗ NO MATCH")
	}
}

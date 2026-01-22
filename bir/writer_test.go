// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package bir

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

var updateWriter = flag.Bool("update-writer", false, "update expected serialized BIR files")

// showBinaryDiff generates a detailed diff string showing differences between expected and actual binary data.
func showBinaryDiff(expected, actual []byte) string {
	var diffBuilder strings.Builder
	diffBuilder.WriteString("\nBinary mismatch - showing differences:\n\n")

	maxLen := len(expected)
	if len(actual) > maxLen {
		maxLen = len(actual)
	}

	diffCount := 0
	const maxDiffsToShow = 20

	for i := 0; i < maxLen && diffCount < maxDiffsToShow; i++ {
		var expectedByte, actualByte byte
		expectedPresent := i < len(expected)
		actualPresent := i < len(actual)

		if expectedPresent {
			expectedByte = expected[i]
		}
		if actualPresent {
			actualByte = actual[i]
		}

		if expectedByte != actualByte || !expectedPresent || !actualPresent {
			diffCount++
			diffBuilder.WriteString(fmt.Sprintf("Offset 0x%04x (byte %d):\n", i, i))
			if !expectedPresent {
				diffBuilder.WriteString("  Expected: (missing)\n")
			} else {
				diffBuilder.WriteString(fmt.Sprintf("  Expected: 0x%02x (%d)\n", expectedByte, expectedByte))
			}
			if !actualPresent {
				diffBuilder.WriteString("  Actual:   (missing)\n\n")
			} else {
				diffBuilder.WriteString(fmt.Sprintf("  Actual:   0x%02x (%d)\n\n", actualByte, actualByte))
			}
		}
	}

	if diffCount >= maxDiffsToShow {
		diffBuilder.WriteString(fmt.Sprintf("... (showing first %d differences, more exist)\n", maxDiffsToShow))
	}

	diffBuilder.WriteString(fmt.Sprintf("Total bytes different: %d+\n", diffCount))
	diffBuilder.WriteString(fmt.Sprintf("Expected size: %d bytes\n", len(expected)))
	diffBuilder.WriteString(fmt.Sprintf("Actual size:   %d bytes\n", len(actual)))
	diffBuilder.WriteString("Use diff tool for full comparison\n")

	return diffBuilder.String()
}

func TestWriter(t *testing.T) {
	flag.Parse()
	birFiles := getCorpusFiles(t)
	for _, birFile := range birFiles {
		t.Run(birFile, func(t *testing.T) {
			t.Parallel()
			testWriter(t, birFile)
		})
	}
}

func testWriter(t *testing.T, birFile string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while serializing BIR package from %s: %v", birFile, r)
		}
	}()

	// Read original file bytes
	originalData, err := os.ReadFile(birFile)
	if err != nil {
		t.Fatalf("failed to read original BIR file: %v", err)
	}

	// Load BIR package
	file, err := os.Open(birFile)
	if err != nil {
		t.Fatalf("failed to open test BIR file: %v", err)
	}
	defer file.Close()

	pkg, err := LoadBIRPackageFromReader(file)
	if err != nil {
		t.Errorf("error loading BIR package from %s: %v", birFile, err)
		return
	}

	if pkg == nil {
		t.Errorf("BIR package is nil for %s", birFile)
		return
	}

	// Serialize the package
	writer := NewBIRBinaryWriter(pkg, nil)
	serializedData, err := writer.Serialize()
	if err != nil {
		t.Errorf("error serializing BIR package from %s: %v", birFile, err)
		return
	}

	// Compare bytes using bytes.Equal
	if !bytes.Equal(originalData, serializedData) {
		diff := showBinaryDiff(originalData, serializedData)
		t.Errorf("serialized BIR does not match original for %s\n%s", birFile, diff)
		return
	}
}

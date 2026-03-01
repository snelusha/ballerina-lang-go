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
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package corpus

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseErrorLinesFromBal(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name     string
		content  string
		expected []int
	}{
		{
			name:     "single @error",
			content:  "line1\nline2\nline3\nfoo(); // @error\nline5",
			expected: []int{4},
		},
		{
			name:     "multiple @error on different lines",
			content:  "a\nb // @error\nc\nd // @error\ne",
			expected: []int{2, 4},
		},
		{
			name:     "no @error",
			content:  "line1\nline2\nline3",
			expected: nil,
		},
		{
			name:     "empty file",
			content:  "",
			expected: nil,
		},
		{
			name:     "comment without @error",
			content:  "x = 1; // some comment",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, "test.bal")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("write file: %v", err)
			}
			got, err := parseErrorLinesFromBal(path)
			if err != nil {
				t.Fatalf("parseErrorLinesFromBal: %v", err)
			}
			if len(got) != len(tt.expected) {
				t.Errorf("got %v, expected %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("got %v, expected %v", got, tt.expected)
					break
				}
			}
		})
	}
}

func TestParseErrorLineNumbersFromStderr(t *testing.T) {
	tests := []struct {
		name         string
		stderr       string
		baseFileName string
		expected     map[int]bool
	}{
		{
			name: "single diagnostic",
			stderr: `error[SEMANTIC_ERROR]: incompatible types
  --> 1-e.bal:21:12
   |
21 |     ignore(b == x); // @error
   |            ^^^^^^
`,
			baseFileName: "1-e.bal",
			expected:     map[int]bool{21: true},
		},
		{
			name: "multiple diagnostics same file",
			stderr: `error: msg1
  --> foo.bal:10:5
   |
error: msg2
  --> foo.bal:20:3
   |
`,
			baseFileName: "foo.bal",
			expected:     map[int]bool{10: true, 20: true},
		},
		{
			name: "wrong file ignored",
			stderr: `error: msg
  --> other.bal:5:1
   |
`,
			baseFileName: "mine.bal",
			expected:     map[int]bool{},
		},
		{
			name: "malformed lines ignored",
			stderr: `not a diagnostic
--> missing-colon
error: msg
  --> bar.bal:7:2
   |
`,
			baseFileName: "bar.bal",
			expected:     map[int]bool{7: true},
		},
		{
			name:         "empty stderr",
			stderr:       "",
			baseFileName: "x.bal",
			expected:     map[int]bool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseErrorLineNumbersFromStderr(tt.stderr, tt.baseFileName)
			if len(got) != len(tt.expected) {
				t.Errorf("got %v, expected %v", got, tt.expected)
				return
			}
			for line, want := range tt.expected {
				if !want {
					continue
				}
				if !got[line] {
					t.Errorf("line %d: got false, expected true", line)
				}
			}
		})
	}
}

func TestIsErrorTest(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"foo-e.bal", true},
		{"/path/to/subset1/01-function/call01-e.bal", true},
		{"1-e.bal", true},
		{"foo-v.bal", false},
		{"foo.bal", false},
		{"foo-e.bal.txt", false},
	}

	for _, tt := range tests {
		got := isErrorTest(tt.path)
		if got != tt.expected {
			t.Errorf("isErrorTest(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestVerifyErrorLocations(t *testing.T) {
	dir := t.TempDir()
	balPath := filepath.Join(dir, "test-e.bal")
	content := "line1\nline2\nfoo(); // @error\nline4"
	if err := os.WriteFile(balPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	t.Run("all expected lines reported", func(t *testing.T) {
		stderr := `error: msg
  --> test-e.bal:3:5
   |
`
		missing, err := verifyErrorLocations(balPath, stderr)
		if err != nil {
			t.Fatalf("verifyErrorLocations: %v", err)
		}
		if len(missing) != 0 {
			t.Errorf("got missing lines %v, expected none", missing)
		}
	})

	t.Run("expected line not reported", func(t *testing.T) {
		stderr := `error: msg
  --> test-e.bal:1:1
   |
`
		missing, err := verifyErrorLocations(balPath, stderr)
		if err != nil {
			t.Fatalf("verifyErrorLocations: %v", err)
		}
		if len(missing) != 1 || missing[0] != 3 {
			t.Errorf("got missing lines %v, expected [3]", missing)
		}
	})
}

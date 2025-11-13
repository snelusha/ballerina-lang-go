package centralclient

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var tempBalaCache string

func TestMain(t *testing.M) {
	tempBalaCache = filepath.Join("build", "temp-test-utils-bala-cache")
	if err := os.MkdirAll(tempBalaCache, 0o755); err != nil {
		panic(err)
	}

	code := t.Run()

	if err := os.RemoveAll(tempBalaCache); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestGetAsList(t *testing.T) {
	tests := []struct {
		name          string
		versionList   string
		expectedArray []string
	}{
		{
			name:          "empty array",
			versionList:   "[]",
			expectedArray: []string{},
		},
		{
			name:          "single version",
			versionList:   "[\"1.1.11\"]",
			expectedArray: []string{"1.1.11"},
		},
		{
			name:          "multiple versions",
			versionList:   "[\"1.0.0\", \"1.2.0\"]",
			expectedArray: []string{"1.0.0", "1.2.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versions, err := getAsList(tt.versionList)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(versions) != len(tt.expectedArray) {
				t.Errorf("expected %d versions, got %d", len(tt.expectedArray), len(versions))
				return
			}
			for i, v := range versions {
				if v != tt.expectedArray[i] {
					t.Errorf("at index %d: expected %s, got %s", i, tt.expectedArray[i], v)
				}
			}
		})
	}
}

func TestValidatePackageVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		isValid bool
	}{
		{"valid simple version", "1.0.0", true},
		{"valid multi-digit version", "1.1.11", true},
		{"valid snapshot version", "2.2.2-snapshot", true},
		{"valid snapshot with number", "2.2.2-snapshot-1", true},
		{"valid alpha version", "2.2.2-alpha", true},
		{"invalid short version", "200", false},
		{"invalid four-part version", "2.2.2.2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewLogFormatter(false)

			_, err := validatePackageVersion(tt.version, formatter)

			if tt.isValid {
				if err != nil {
					t.Errorf("expected valid version, got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error for invalid version: %s", tt.version)
				} else if !strings.Contains(err.Error(), "Invalid version:") {
					t.Errorf("error message should contain 'Invalid version:', got: %s", err.Error())
				}
			}
		})
	}
}

func TestJsonContentTypeChecker(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{"plain application/json", "application/json", true},
		{"json with charset", "application/json; charset=utf-8", true},
		{"octet-stream", "application/octet-stream", false},
		{"text/plain", "text/plain", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isApplicationJSONContentType(tt.contentType)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for content type: %s", tt.expected, result, tt.contentType)
			}
		})
	}
}

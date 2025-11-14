package centralclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

const (
	testBalVersion = "slp5"
	anyPlatform    = "any"
	testBalaName   = "sf-any.bala"
	winery         = "winery"
	accessToken    = "273cc9f6-c333-36ab-aa2q-f08e9513ff5y"
)

var utilTestResources = filepath.Join("testdata", "utils")

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type TestRunner = func(client CentralAPIClient) (string, string)

type TestCase struct {
	runner         TestRunner
	name           string
	filepath       string
	expectedOutput string
	expectedError  string
}

func parseTestCases(dir string) ([]TestCase, error) {
	var testCases []TestCase

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txtar") {
			filepath := filepath.Join(dir, file.Name())
			tc, err := txtar.ParseFile(filepath)
			if err != nil {
				return nil, err
			}
			testCase, err := parseTestCase(tc, filepath)
			if err != nil {
				return nil, err
			}

			testCases = append(testCases, testCase)
		}
	}

	return testCases, nil
}

func parseTestCase(archive *txtar.Archive, filepath string) (TestCase, error) {
	if archive == nil || len(archive.Files) <= 2 {
		return TestCase{}, fmt.Errorf("invalid test case archive")
	}

	tr := parseInput(archive.Files[0])

	return TestCase{
		runner:         tr,
		name:           strings.TrimSuffix(path.Base(filepath), ".txtar"),
		filepath:       filepath,
		expectedOutput: strings.TrimSpace(string(archive.Files[1].Data)),
		expectedError:  strings.TrimSpace(string(archive.Files[2].Data)),
	}, nil
}

func parseInput(data txtar.File) TestRunner {
	content := strings.ReplaceAll(string(data.Data), "\r\n", "\n")
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		panic(fmt.Sprintf("empty input data in test case: %s", data.Name))
	}
	switch lines[0] {
	case "GetPackageVersions":
		return func(client CentralAPIClient) (string, string) {
			versions, err := client.GetPackageVersions(lines[1], lines[2], lines[3], lines[4])
			if err != nil {
				return "", err.Error()
			}
			return fmt.Sprintf("%v", versions), ""
		}
	case "GetPackage":
		return func(client CentralAPIClient) (string, string) {
			pkg, err := client.GetPackage(lines[1], lines[2], lines[3], lines[4], lines[5])
			if err != nil {
				return "", err.Error()
			}
			return fmt.Sprintf("org=%s name=%s version=%s", pkg.Organization, pkg.Name, pkg.Version), ""
		}
	case "GetConnectors":
		return func(client CentralAPIClient) (string, string) {
			params := make(map[string]string)
			params["q"] = lines[1]
			connectors, err := client.GetConnectors(params, lines[2], lines[3])
			if err != nil {
				return "", err.Error()
			}
			return fmt.Sprintf("%v", connectors != nil), ""
		}
	case "GetConnector":
		return func(client CentralAPIClient) (string, string) {
			connector, err := client.GetConnector(lines[1], lines[2], lines[3])
			if err != nil {
				return "", err.Error()
			}
			return fmt.Sprintf("%v", connector), ""
		}
	case "GetTriggers":
		return func(client CentralAPIClient) (string, string) {
			params := make(map[string]string)
			params["q"] = lines[1]
			triggers, err := client.GetTriggers(params, lines[2], lines[3])
			if err != nil {
				return "", err.Error()
			}
			return fmt.Sprintf("%v", triggers), ""
		}
	case "GetTrigger":
		return func(client CentralAPIClient) (string, string) {
			trigger, err := client.GetTrigger(lines[1], lines[2], lines[3])
			if err != nil {
				return "", err.Error()
			}
			return fmt.Sprintf("%v", trigger), ""
		}
	default:
		panic(fmt.Sprintf("unsupported test case type: %s (file: %s)", lines[0], data.Name))
	}
}

func updateTestCase(tc TestCase, actualOutput, actualError string) error {
	archive, err := txtar.ParseFile(tc.filepath)
	if err != nil {
		return err
	}

	if len(archive.Files) < 3 {
		return fmt.Errorf("invalid archive structure")
	}

	if actualOutput != "" {
		archive.Files[1].Data = fmt.Appendf(nil, "%s\n\n", actualOutput)
	}

	if actualError != "" {
		archive.Files[2].Data = fmt.Appendf(nil, "%s\n\n", actualError)
	}

	return os.WriteFile(tc.filepath, txtar.Format(archive), 0o644)
}

func TestTxtarTestCases(t *testing.T) {
	bless := os.Getenv("BLESS") == "1" || os.Getenv("BLESS") == "true"

	testCases, err := parseTestCases("testdata")
	if err != nil {
		t.Fatalf("failed to parse test cases: %v", err)
	}

	packageJSONPath := filepath.Join(utilTestResources, "package.json")
	packageJSON, _ := os.ReadFile(packageJSONPath)
	packageSearchPath := filepath.Join(utilTestResources, "packageSearch.json")
	packageSearchJSON, _ := os.ReadFile(packageSearchPath)

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			// GetPackageVersions endpoints
			switch req.URL.Path {
			case "/registry/packages/wso2/sf":
				body := `["1.0.0", "1.1.0", "1.2.0"]`
				return NewJSONResponse(http.StatusOK, body, req), nil
			case "/registry/packages/unknown/package":
				body := `{"message":"package not found: unknown/package:*_any"}`
				return NewJSONResponse(http.StatusNotFound, body, req), nil
			case "/registry/packages/testorg/testpkg":
				body := `{"message":"unauthorized access token for organization: 'testorg'"}`
				return NewJSONResponse(http.StatusUnauthorized, body, req), nil
			case "/registry/packages/testorg/bad-pkg":
				body := `{"message":"invalid package name format"}`
				return NewJSONResponse(http.StatusBadRequest, body, req), nil
			case "/registry/packages/testorg/internalerror":
				body := `{"message":"internal server error occurred"}`
				return NewJSONResponse(http.StatusInternalServerError, body, req), nil
			case "/registry/packages/testorg/unavailable":
				body := `{"message":"service temporarily unavailable"}`
				return NewJSONResponse(http.StatusServiceUnavailable, body, req), nil
			case "/registry/packages/testorg/invalidjson":
				body := `invalid json response`
				return NewJSONResponse(http.StatusOK, body, req), nil

			// GetPackage endpoints
			case "/registry/packages/foo/winery/1.3.5":
				return NewJSONResponse(http.StatusOK, string(packageJSON), req), nil
			case "/registry/packages/unknown/notfound/1.0.0":
				body := `{"message":"package not found for: unknown/notfound:1.0.0"}`
				return NewJSONResponse(http.StatusNotFound, body, req), nil
			case "/registry/packages/testorg/unauthorized/1.0.0":
				body := `{"message":"unauthorized access token"}`
				return NewJSONResponse(http.StatusUnauthorized, body, req), nil
			case "/registry/packages/testorg/badrequest/1.0.0":
				body := `{"message":"invalid version format"}`
				return NewJSONResponse(http.StatusBadRequest, body, req), nil
			case "/registry/packages/testorg/servererror/1.0.0":
				body := `{"message":"database connection failed"}`
				return NewJSONResponse(http.StatusInternalServerError, body, req), nil

			// GetConnectors endpoints
			case "/registry/connectors":
				if req.URL.Query().Get("q") == "unauthorized" {
					body := `{"message":"unauthorized access"}`
					return NewJSONResponse(http.StatusUnauthorized, body, req), nil
				}
				if req.URL.Query().Get("q") == "badrequest" {
					body := `{"message":"invalid query parameter"}`
					return NewJSONResponse(http.StatusBadRequest, body, req), nil
				}
				if req.URL.Query().Get("q") == "servererror" {
					body := `{"message":"internal server error"}`
					return NewJSONResponse(http.StatusInternalServerError, body, req), nil
				}
				return NewJSONResponse(http.StatusOK, string(packageSearchJSON), req), nil

			// GetConnector endpoints
			case "/registry/connectors/123":
				body := `{"id": "123", "organization": "foo", "name": "winery", "version": "1.3.5"}`
				return NewJSONResponse(http.StatusOK, body, req), nil
			case "/registry/connectors/notfound":
				body := `{"message":"connector not found"}`
				return NewJSONResponse(http.StatusNotFound, body, req), nil
			case "/registry/connectors/unauthorized":
				body := `{"message":"unauthorized access"}`
				return NewJSONResponse(http.StatusUnauthorized, body, req), nil
			case "/registry/connectors/invalidjson":
				body := `invalid json`
				return NewJSONResponse(http.StatusOK, body, req), nil
			case "/registry/connectors/servererror":
				body := `{"message":"internal server error"}`
				return NewJSONResponse(http.StatusInternalServerError, body, req), nil

			// GetTriggers endpoints
			case "/registry/triggers":
				if req.URL.Query().Get("q") == "unauthorized" {
					body := `{"message":"unauthorized access"}`
					return NewJSONResponse(http.StatusUnauthorized, body, req), nil
				}
				if req.URL.Query().Get("q") == "badrequest" {
					body := `{"message":"invalid query parameter"}`
					return NewJSONResponse(http.StatusBadRequest, body, req), nil
				}
				if req.URL.Query().Get("q") == "servererror" {
					body := `{"message":"internal server error"}`
					return NewJSONResponse(http.StatusInternalServerError, body, req), nil
				}
				body := `{"count": 2, "triggers": [{"id": "1", "name": "trigger1"}]}`
				return NewJSONResponse(http.StatusOK, body, req), nil

			// GetTrigger endpoints
			case "/registry/triggers/456":
				body := `{"id": "456", "name": "http-trigger", "type": "http"}`
				return NewJSONResponse(http.StatusOK, body, req), nil
			case "/registry/triggers/notfound":
				body := `{"message":"trigger not found"}`
				return NewJSONResponse(http.StatusNotFound, body, req), nil
			case "/registry/triggers/unauthorized":
				body := `{"message":"unauthorized access"}`
				return NewJSONResponse(http.StatusUnauthorized, body, req), nil
			case "/registry/triggers/invalidjson":
				body := `invalid json`
				return NewJSONResponse(http.StatusOK, body, req), nil
			case "/registry/triggers/servererror":
				body := `{"message":"internal server error"}`
				return NewJSONResponse(http.StatusInternalServerError, body, req), nil

			default:
				return NewJSONResponse(http.StatusNotFound, ``, req), nil
			}
		}),
	}

	client := NewTestCentralAPIClient(mockClient)

	for _, tc := range testCases {
		output, errStr := tc.runner(client)

		if bless {
			if err := updateTestCase(tc, output, errStr); err != nil {
				t.Fatalf("failed to update test case %s: %v", tc.name, err)
			}
			t.Logf("blessed test case %s", tc.name)
			continue
		}

		if output != tc.expectedOutput {
			t.Errorf("test case %s: expected output '%s', got '%s'", tc.name, tc.expectedOutput, output)
		}
		if errStr != tc.expectedError {
			t.Errorf("test case %s: expected error '%s', got '%s'", tc.name, tc.expectedError, errStr)
		}
	}
}

func NewJSONResponse(status int, body string, req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Request:       req,
		Proto:         "HTTP/1.1",
		ProtoMinor:    1,
		ProtoMajor:    1,
		ContentLength: int64(len(body)),
		Close:         true,
	}
}

func NewTestCentralAPIClient(mockClient *http.Client) CentralAPIClient {
	return &centralAPIClientImpl{
		baseURL:     "https://localhost:9090/registry",
		httpClient:  mockClient,
		accessToken: accessToken,
		outStream:   os.Stdout,
		maxRetries:  2,
	}
}

func TestGetPackageVersions(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `["1.0.0", "1.1.0", "1.2.0"]`
			return NewJSONResponse(http.StatusOK, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	versions, err := tc.GetPackageVersions("wso2", "sf", anyPlatform, testBalVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedVersions := []string{"1.0.0", "1.1.0", "1.2.0"}
	if len(versions) != len(expectedVersions) {
		t.Fatalf("expected %d versions, got %d", len(expectedVersions), len(versions))
	}
	for i, v := range versions {
		if v != expectedVersions[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expectedVersions[i], v)
		}
	}
}

func TestGetPackageVersionsNotFound(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"package not found: unknown/package:*_any"}`
			return NewJSONResponse(http.StatusNotFound, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	versions, err := tc.GetPackageVersions("unknown", "package", anyPlatform, testBalVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(versions) != 0 {
		t.Fatalf("expected 0 versions, got %d", len(versions))
	}
}

func TestGetPackageVersionsUnauthorized(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"unauthorized access token for organization: 'testorg'"}`
			return NewJSONResponse(http.StatusUnauthorized, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackageVersions("testorg", "testpkg", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for unauthorized access, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "testorg") {
		t.Errorf("expected error message to contain 'testorg', got: %s", err.Error())
	}
}

func TestGetPackageVersionsBadRequest(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"invalid package name format"}`
			return NewJSONResponse(http.StatusBadRequest, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackageVersions("testorg", "bad-pkg", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for bad request, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "invalid package name format") {
		t.Errorf("expected error message to contain 'invalid package name format', got: %s", err.Error())
	}
}

func TestGetPackageVersionsInternalServerError(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"internal server error occurred"}`
			return NewJSONResponse(http.StatusInternalServerError, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackageVersions("testorg", "testpkg", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for internal server error, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "internal server error occurred") {
		t.Errorf("expected error message to contain 'internal server error occurred', got: %s", err.Error())
	}
}

func TestGetPackageVersionsServiceUnavailable(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"service temporarily unavailable"}`
			return NewJSONResponse(http.StatusServiceUnavailable, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackageVersions("testorg", "testpkg", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for service unavailable, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "service temporarily unavailable") {
		t.Errorf("expected error message to contain 'service temporarily unavailable', got: %s", err.Error())
	}
}

func TestGetPackageVersionsInvalidJSON(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `invalid json response`
			return NewJSONResponse(http.StatusOK, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackageVersions("testorg", "testpkg", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}

	fmt.Println(err.Error())

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), ErrCannotFindVersions) {
		t.Errorf("expected error message to contain '%s', got: %s", ErrCannotFindVersions, err.Error())
	}
}

func TestGetPackageVersionsNonJSONContentType(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Body:       io.NopCloser(strings.NewReader("plain text response")),
				Header: http.Header{
					"Content-Type": []string{"text/plain"},
				},
				Request:       req,
				Proto:         "HTTP/1.1",
				ProtoMinor:    1,
				ProtoMajor:    1,
				ContentLength: 19,
				Close:         true,
			}, nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackageVersions("testorg", "testpkg", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for non-JSON content type, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), ErrCannotFindVersions) {
		t.Errorf("expected error message to contain '%s', got: %s", ErrCannotFindVersions, err.Error())
	}
}

func TestGetPackageVersionsEmptyVersionsList(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `[]`
			return NewJSONResponse(http.StatusOK, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	versions, err := tc.GetPackageVersions("testorg", "testpkg", anyPlatform, testBalVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(versions) != 0 {
		t.Errorf("expected 0 versions for empty list, got %d", len(versions))
	}
}

func TestGetPackageVersionsNotFoundWithError(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"repository error: cannot access database"}`
			return NewJSONResponse(http.StatusNotFound, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackageVersions("testorg", "testpkg", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for 404 with non-package-not-found message, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "repository error") {
		t.Errorf("expected error message to contain 'repository error', got: %s", err.Error())
	}
}

func TestGetPackage(t *testing.T) {
	packageJSONPath := filepath.Join(utilTestResources, "package.json")
	packageJSON, err := os.ReadFile(packageJSONPath)
	if err != nil {
		t.Fatalf("failed to read package.json: %v", err)
	}

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return NewJSONResponse(http.StatusOK, string(packageJSON), req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	pkg, err := tc.GetPackage("foo", winery, "1.3.5", anyPlatform, testBalVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pkg.Organization != "foo" {
		t.Errorf("expected organization 'foo', got '%s'", pkg.Organization)
	}
	if pkg.Name != winery {
		t.Errorf("expected name 'winery', got '%s'", pkg.Name)
	}
	if pkg.Version != "1.3.5" {
		t.Errorf("expected version '1.3.5', got '%s'", pkg.Version)
	}
}

func TestGetPackageNotFound(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"package not found for: unknown/package:1.0.0"}`
			return NewJSONResponse(http.StatusNotFound, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackage("unknown", "package", "1.0.0", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for package not found, got nil")
	}

	if _, ok := err.(*NoPackageError); !ok {
		t.Errorf("expected NoPackageError, got %T", err)
	}
}

func TestGetPackageUnauthorized(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"unauthorized access token"}`
			return NewJSONResponse(http.StatusUnauthorized, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackage("foo", winery, "1.3.5", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for unauthorized access, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetPackageBadRequest(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"invalid version format"}`
			return NewJSONResponse(http.StatusBadRequest, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackage("foo", winery, "invalid", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for bad request, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "invalid version format") {
		t.Errorf("expected error message to contain 'invalid version format', got: %s", err.Error())
	}
}

func TestGetPackageInternalServerError(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"database connection failed"}`
			return NewJSONResponse(http.StatusInternalServerError, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetPackage("foo", winery, "1.3.5", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for internal server error, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetConnectors(t *testing.T) {
	packageSearchPath := filepath.Join(utilTestResources, "packageSearch.json")
	packageSearchJSON, err := os.ReadFile(packageSearchPath)
	if err != nil {
		t.Fatalf("failed to read packageSearch.json: %v", err)
	}

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return NewJSONResponse(http.StatusOK, string(packageSearchJSON), req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	params := map[string]string{
		"q": "winery",
	}

	connectors, err := tc.GetConnectors(params, anyPlatform, testBalVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if connectors == nil {
		t.Fatal("expected connectors, got nil")
	}
}

func TestGetConnectorsUnauthorized(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"unauthorized access"}`
			return NewJSONResponse(http.StatusUnauthorized, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	params := map[string]string{"q": "test"}

	_, err := tc.GetConnectors(params, anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for unauthorized access, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetConnectorsBadRequest(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"invalid query parameter"}`
			return NewJSONResponse(http.StatusBadRequest, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	params := map[string]string{"invalid": "param"}

	_, err := tc.GetConnectors(params, anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for bad request, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetConnector(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{
				"id": "123",
				"organization": "foo",
				"name": "winery",
				"version": "1.3.5"
			}`
			return NewJSONResponse(http.StatusOK, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	connector, err := tc.GetConnector("123", anyPlatform, testBalVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if connector == nil {
		t.Fatal("expected connector, got nil")
	}

	if id, ok := connector["id"].(string); !ok || id != "123" {
		t.Errorf("expected id '123', got '%v'", connector["id"])
	}
}

func TestGetConnectorNotFound(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"connector not found"}`
			return NewJSONResponse(http.StatusNotFound, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetConnector("999", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for connector not found, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetConnectorUnauthorized(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"unauthorized access"}`
			return NewJSONResponse(http.StatusUnauthorized, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetConnector("123", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for unauthorized access, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetConnectorInvalidJSON(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `invalid json`
			return NewJSONResponse(http.StatusOK, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetConnector("123", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetTriggers(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{
				"triggers": [
					{
						"organization": "ballerina",
						"name": "http",
						"version": "1.0.0"
					}
				],
				"count": 1
			}`
			return NewJSONResponse(http.StatusOK, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	params := map[string]string{"q": "http"}

	triggers, err := tc.GetTriggers(params, anyPlatform, testBalVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if triggers == nil {
		t.Fatal("expected triggers, got nil")
	}
}

func TestGetTriggersUnauthorized(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"unauthorized access"}`
			return NewJSONResponse(http.StatusUnauthorized, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	params := map[string]string{"q": "test"}

	_, err := tc.GetTriggers(params, anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for unauthorized access, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetTriggersBadRequest(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"invalid parameters"}`
			return NewJSONResponse(http.StatusBadRequest, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	params := map[string]string{"invalid": "param"}

	_, err := tc.GetTriggers(params, anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for bad request, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetTriggersInternalServerError(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"internal server error"}`
			return NewJSONResponse(http.StatusInternalServerError, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	params := map[string]string{"q": "test"}

	_, err := tc.GetTriggers(params, anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for internal server error, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetTrigger(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{
				"id": "456",
				"organization": "ballerina",
				"name": "http",
				"version": "1.0.0"
			}`
			return NewJSONResponse(http.StatusOK, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	trigger, err := tc.GetTrigger("456", anyPlatform, testBalVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if trigger == nil {
		t.Fatal("expected trigger, got nil")
	}

	if id, ok := trigger["id"].(string); !ok || id != "456" {
		t.Errorf("expected id '456', got '%v'", trigger["id"])
	}
}

func TestGetTriggerNotFound(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"trigger not found"}`
			return NewJSONResponse(http.StatusNotFound, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetTrigger("999", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for trigger not found, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetTriggerUnauthorized(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"unauthorized access"}`
			return NewJSONResponse(http.StatusUnauthorized, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetTrigger("456", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for unauthorized access, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetTriggerInvalidJSON(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `not valid json`
			return NewJSONResponse(http.StatusOK, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetTrigger("456", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestGetTriggerServiceUnavailable(t *testing.T) {
	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"service unavailable"}`
			return NewJSONResponse(http.StatusServiceUnavailable, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	_, err := tc.GetTrigger("456", anyPlatform, testBalVersion)
	if err == nil {
		t.Fatal("expected error for service unavailable, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestPullPackageSuccess(t *testing.T) {
	balaPath := filepath.Join(utilTestResources, testBalaName)
	balaContent, err := os.ReadFile(balaPath)
	if err != nil {
		t.Fatalf("failed to read bala file: %v", err)
	}

	tempDir := t.TempDir()

	expectedBalaFileName := "sf-2020r2-any-1.3.5.bala"

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "/registry/packages/") {
				return &http.Response{
					StatusCode: http.StatusFound,
					Status:     http.StatusText(http.StatusFound),
					Body:       io.NopCloser(strings.NewReader("")),
					Header: http.Header{
						"Location":            []string{"https://fileserver.dev-central.ballerina.io/2.0/wso2/sf/1.3.5/" + expectedBalaFileName},
						"Content-Disposition": []string{"attachment; filename=" + expectedBalaFileName},
					},
					Request:    req,
					Proto:      "HTTP/1.1",
					ProtoMinor: 1,
					ProtoMajor: 1,
					Close:      true,
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Body:       io.NopCloser(strings.NewReader(string(balaContent))),
				Header: http.Header{
					"Content-Type":   []string{ApplicationOctetStream},
					"Content-Length": []string{fmt.Sprintf("%d", len(balaContent))},
				},
				Request:       req,
				Proto:         "HTTP/1.1",
				ProtoMinor:    1,
				ProtoMajor:    1,
				ContentLength: int64(len(balaContent)),
				Close:         true,
			}, nil
		}),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	tc := NewTestCentralAPIClient(mockClient)

	err = tc.PullPackage("wso2", "sf", "1.3.5", tempDir, "2020r2-any", testBalVersion, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	balaDir := filepath.Join(tempDir, "1.3.5", "2020r2-any")

	if info, err := os.Stat(balaDir); err != nil || !info.IsDir() {
		t.Fatalf("bala directory does not exist: %s", balaDir)
	}

	requiredFiles := []string{
		"bala.json",
		"package.json",
		"dependency-graph.json",
	}

	for _, file := range requiredFiles {
		filePath := filepath.Join(balaDir, file)
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("required file does not exist: %s", file)
		}
	}

	requiredDirs := []string{
		"modules",
	}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(balaDir, dir)
		if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
			t.Errorf("required directory does not exist: %s", dir)
		}
	}
}

func TestPullPackageSuccessWithDeprecation(t *testing.T) {
	balaPath := filepath.Join(utilTestResources, testBalaName)
	balaContent, err := os.ReadFile(balaPath)
	if err != nil {
		t.Fatalf("failed to read bala file: %v", err)
	}

	tempDir := t.TempDir()

	expectedBalaFileName := "sf-2020r2-any-1.3.5.bala"

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "/registry/packages/") {
				return &http.Response{
					StatusCode: http.StatusFound,
					Status:     http.StatusText(http.StatusFound),
					Body:       io.NopCloser(strings.NewReader("")),
					Header: http.Header{
						"Location":            []string{"https://fileserver.dev-central.ballerina.io/2.0/wso2/sf/1.3.5/" + expectedBalaFileName},
						"Content-Disposition": []string{"attachment; filename=" + expectedBalaFileName},
						"Is-Deprecated":       []string{"true"},
						"Deprecate-Message":   []string{"This package is deprecated. Please use the new version."},
					},
					Request:    req,
					Proto:      "HTTP/1.1",
					ProtoMinor: 1,
					ProtoMajor: 1,
					Close:      true,
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Body:       io.NopCloser(strings.NewReader(string(balaContent))),
				Header: http.Header{
					"Content-Type":   []string{ApplicationOctetStream},
					"Content-Length": []string{fmt.Sprintf("%d", len(balaContent))},
				},
				Request:       req,
				Proto:         "HTTP/1.1",
				ProtoMinor:    1,
				ProtoMajor:    1,
				ContentLength: int64(len(balaContent)),
				Close:         true,
			}, nil
		}),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	tc := NewTestCentralAPIClient(mockClient)

	err = tc.PullPackage("wso2", "sf", "1.3.5", tempDir, "2020r2-any", testBalVersion, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	balaDir := filepath.Join(tempDir, "1.3.5", "2020r2-any")

	if info, err := os.Stat(balaDir); err != nil || !info.IsDir() {
		t.Fatalf("bala directory does not exist: %s", balaDir)
	}

	deprecatedFile := filepath.Join(balaDir, "deprecated.txt")
	if _, err := os.Stat(deprecatedFile); err != nil {
		t.Errorf("deprecated.txt file does not exist")
	} else {
		content, err := os.ReadFile(deprecatedFile)
		if err != nil {
			t.Errorf("failed to read deprecated.txt: %v", err)
		} else if !strings.Contains(string(content), "This package is deprecated") {
			t.Errorf("deprecated.txt does not contain expected message, got: %s", string(content))
		}
	}
}

func TestPullPackageNotFound(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"package not found: wso2/unknown:1.0.0"}`
			return NewJSONResponse(http.StatusNotFound, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "unknown", "1.0.0", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for package not found, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "package not found") {
		t.Errorf("expected error message to contain 'package not found', got: %s", err.Error())
	}
}

func TestPullPackageUnauthorized(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"unauthorized access token for organization: 'wso2'"}`
			return NewJSONResponse(http.StatusUnauthorized, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "sf", "1.3.5", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for unauthorized access, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "wso2") {
		t.Errorf("expected error message to contain 'wso2', got: %s", err.Error())
	}
}

func TestPullPackageBadRequest(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"invalid package version format"}`
			return NewJSONResponse(http.StatusBadRequest, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "sf", "invalid", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for bad request, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "invalid package version format") {
		t.Errorf("expected error message to contain 'invalid package version format', got: %s", err.Error())
	}
}

func TestPullPackageInternalServerError(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"internal server error occurred"}`
			return NewJSONResponse(http.StatusInternalServerError, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "sf", "1.3.5", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for internal server error, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "internal server error occurred") {
		t.Errorf("expected error message to contain 'internal server error occurred', got: %s", err.Error())
	}
}

func TestPullPackageServiceUnavailable(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			body := `{"message":"service temporarily unavailable"}`
			return NewJSONResponse(http.StatusServiceUnavailable, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "sf", "1.3.5", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for service unavailable, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "service temporarily unavailable") {
		t.Errorf("expected error message to contain 'service temporarily unavailable', got: %s", err.Error())
	}
}

func TestPullPackageMissingBalaURL(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusFound,
				Status:     http.StatusText(http.StatusFound),
				Body:       io.NopCloser(strings.NewReader("")),
				Header: http.Header{
					"Digest": []string{"sha-256=abc123"},
				},
				Request:    req,
				Proto:      "HTTP/1.1",
				ProtoMinor: 1,
				ProtoMajor: 1,
				Close:      true,
			}, nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "sf", "1.3.5", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for missing bala URL, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), ErrCannotPullPackage) {
		t.Errorf("expected error message to contain '%s', got: %s", ErrCannotPullPackage, err.Error())
	}
}

func TestPullPackageBalaDownloadFailed(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Accept") == ApplicationOctetStream && !strings.Contains(req.URL.String(), "fileserver") {
				return &http.Response{
					StatusCode: http.StatusFound,
					Status:     http.StatusText(http.StatusFound),
					Body:       io.NopCloser(strings.NewReader("")),
					Header: http.Header{
						"Location":            []string{"https://fileserver.dev-central.ballerina.io/2.0/wso2/sf/1.3.5/sf-any-1.3.5.bala"},
						"Content-Disposition": []string{testBalaName},
						"Digest":              []string{"sha-256=abc123"},
					},
					Request:    req,
					Proto:      "HTTP/1.1",
					ProtoMinor: 1,
					ProtoMajor: 1,
					Close:      true,
				}, nil
			}
			body := `{"message":"bala file not found"}`
			return NewJSONResponse(http.StatusNotFound, body, req), nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "sf", "1.3.5", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for bala download failure, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}

	if !strings.Contains(err.Error(), "bala file not found") {
		t.Errorf("expected error message to contain 'bala file not found', got: %s", err.Error())
	}
}

func TestPullPackageRedirectWithoutLocationHeader(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusFound,
				Status:     http.StatusText(http.StatusFound),
				Body:       io.NopCloser(strings.NewReader("")),
				Header: http.Header{
					"Content-Disposition": []string{testBalaName},
					"Digest":              []string{"sha-256=abc123"},
				},
				Request:    req,
				Proto:      "HTTP/1.1",
				ProtoMinor: 1,
				ProtoMajor: 1,
				Close:      true,
			}, nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "sf", "1.3.5", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for missing Location header, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestPullPackageRedirectWithoutContentDisposition(t *testing.T) {
	tempDir := t.TempDir()

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusFound,
				Status:     http.StatusText(http.StatusFound),
				Body:       io.NopCloser(strings.NewReader("")),
				Header: http.Header{
					"Location": []string{"https://fileserver.dev-central.ballerina.io/2.0/wso2/sf/1.3.5/sf-any-1.3.5.bala"},
					"Digest":   []string{"sha-256=abc123"},
				},
				Request:    req,
				Proto:      "HTTP/1.1",
				ProtoMinor: 1,
				ProtoMajor: 1,
				Close:      true,
			}, nil
		}),
	}

	tc := NewTestCentralAPIClient(mockClient)

	err := tc.PullPackage("wso2", "sf", "1.3.5", tempDir, anyPlatform, testBalVersion, false)
	if err == nil {
		t.Fatal("expected error for missing Content-Disposition header, got nil")
	}

	if _, ok := err.(*CentralClientError); !ok {
		t.Errorf("expected CentralClientError, got %T", err)
	}
}

func TestPullPackageConnectionResetRetry(t *testing.T) {
	balaPath := filepath.Join(utilTestResources, testBalaName)
	balaContent, err := os.ReadFile(balaPath)
	if err != nil {
		t.Fatalf("failed to read bala file: %v", err)
	}

	tempDir := t.TempDir()
	expectedBalaFileName := "sf-2020r2-any-1.3.5.bala"

	attemptCount := 0
	downloadAttempts := 0

	mockClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "/registry/packages/") {
				attemptCount++
				return &http.Response{
					StatusCode: http.StatusFound,
					Status:     http.StatusText(http.StatusFound),
					Body:       io.NopCloser(strings.NewReader("")),
					Header: http.Header{
						"Location":            []string{"https://fileserver.dev-central.ballerina.io/2.0/foo/sf/1.3.5/" + expectedBalaFileName},
						"Content-Disposition": []string{"attachment; filename=" + expectedBalaFileName},
					},
					Request:    req,
					Proto:      "HTTP/1.1",
					ProtoMinor: 1,
					ProtoMajor: 1,
					Close:      true,
				}, nil
			}

			downloadAttempts++

			if downloadAttempts <= 2 {
				return nil, fmt.Errorf("Connection reset by peer")
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Body:       io.NopCloser(strings.NewReader(string(balaContent))),
				Header: http.Header{
					"Content-Type":           []string{ApplicationOctetStream},
					"Content-Length":         []string{fmt.Sprintf("%d", len(balaContent))},
					"RESOLVED_REQUESTED_URI": []string{"https://fileserver.dev-central.ballerina.io/2.0/foo/sf/1.3.5/" + expectedBalaFileName},
				},
				Request:       req,
				Proto:         "HTTP/1.1",
				ProtoMinor:    1,
				ProtoMajor:    1,
				ContentLength: int64(len(balaContent)),
				Close:         true,
			}, nil
		}),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	tc := NewTestCentralAPIClient(mockClient)

	err = tc.PullPackage("foo", "sf", "1.3.5", tempDir, "2020r2-any", testBalVersion, false)
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}

	balaDir := filepath.Join(tempDir, "1.3.5", "2020r2-any")
	if info, err := os.Stat(balaDir); err != nil || !info.IsDir() {
		t.Fatalf("bala directory does not exist after successful retry: %s", balaDir)
	}

	if attemptCount != 3 {
		t.Errorf("expected 3 total attempts (2 failures + 1 success), got %d", attemptCount)
	}

	requiredFiles := []string{"bala.json", "package.json"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(balaDir, file)
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("required file does not exist after retry: %s", file)
		}
	}
}

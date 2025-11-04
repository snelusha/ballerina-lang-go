package centralclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ballerina-lang-go/centralclient/model"
)

type CentralAPIClient interface {
	GetPackage(orgNamePath, packageNamePath, version, supportedPlatform, ballerinaVersion string) (*model.Package, error)
	GetPackageVersions(orgNamePath, packageNamePath, supportedPlatform, ballerinaVersion string) ([]string, error)
	PushPackage(balaPath, org, name, version, supportedPlatform, ballerinaVersion string) error
	PullPackage(org, name, version, packagePathInBalaCache, supportedPlatform, ballerinaVersion string, isBuild bool) error
	PullTool(toolID, version, balaCacheDirPath, supportedPlatform, ballerinaVersion string, isBuild bool) ([]string, error)
	ResolvePackageNames(request *model.PackageNameResolutionRequest, supportedPlatform, ballerinaVersion string) (*model.PackageNameResolutionResponse, error)
	ResolveDependencies(request *model.PackageResolutionRequest, supportedPlatform, ballerinaVersion string) (*model.PackageResolutionResponse, error)
	ResolveToolDependencies(request *model.ToolResolutionCentralRequest, supportedPlatform, ballerinaVersion string) (*model.ToolResolutionCentralResponse, error)
	SearchPackage(query, supportedPlatform, ballerinaVersion string) (*model.PackageSearchResult, error)
	SearchTool(keyword, supportedPlatform, ballerinaVersion string) (*model.ToolSearchResult, error)
	DeprecatePackage(packageInfo, deprecationMsg, supportedPlatform, ballerinaVersion string, isUndo bool) error
	GetConnectors(params map[string]string, supportedPlatform, ballerinaVersion string) (interface{}, error)
	GetConnector(id, supportedPlatform, ballerinaVersion string) (map[string]interface{}, error)
	GetConnectorByInfo(connector *model.ConnectorInfo, supportedPlatform, ballerinaVersion string) (map[string]interface{}, error)
	GetTriggers(params map[string]string, supportedPlatform, ballerinaVersion string) (interface{}, error)
	GetTrigger(id, supportedPlatform, ballerinaVersion string) (map[string]interface{}, error)
	AccessToken() string
	SetAccessToken(token string)
}

type centralAPIClientImpl struct {
	baseURL        string
	proxyURL       *url.URL
	accessToken    string
	outStream      io.Writer
	verboseEnabled bool
	proxyUsername  string
	proxyPassword  string
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	callTimeout    time.Duration
	maxRetries     int
	httpClient     *http.Client
}

func NewCentralAPIClient(baseURL string, proxyURL *url.URL, accessToken string) CentralAPIClient {
	verboseEnabled := os.Getenv(SysPropCentralVerboseEnabled) == "true"
	return &centralAPIClientImpl{
		baseURL:        baseURL,
		proxyURL:       proxyURL,
		accessToken:    accessToken,
		outStream:      os.Stdout,
		verboseEnabled: verboseEnabled,
		proxyUsername:  "",
		proxyPassword:  "",
		connectTimeout: DefaultConnectTimeout * time.Second,
		readTimeout:    DefaultReadTimeout * time.Second,
		writeTimeout:   DefaultWriteTimeout * time.Second,
		callTimeout:    DefaultCallTimeout * time.Second,
		maxRetries:     MaxRetry,
	}
}

func NewCentralAPIClientWithOptions(baseURL string, proxyURL *url.URL, accessToken string, verboseEnabled bool, maxRetries int, outStream io.Writer) CentralAPIClient {
	return &centralAPIClientImpl{
		baseURL:        baseURL,
		proxyURL:       proxyURL,
		accessToken:    accessToken,
		outStream:      outStream,
		verboseEnabled: verboseEnabled,
		proxyUsername:  "",
		proxyPassword:  "",
		connectTimeout: DefaultConnectTimeout * time.Second,
		readTimeout:    DefaultReadTimeout * time.Second,
		writeTimeout:   DefaultWriteTimeout * time.Second,
		callTimeout:    DefaultCallTimeout * time.Second,
		maxRetries:     maxRetries,
	}
}

func NewCentralAPIClientFull(baseURL string, proxyURL *url.URL, proxyUsername, proxyPassword, accessToken string, connectTimeout, readTimeout, writeTimeout, callTimeout, maxRetries int) CentralAPIClient {
	verboseEnabled := os.Getenv(SysPropCentralVerboseEnabled) == "true"
	return &centralAPIClientImpl{
		baseURL:        baseURL,
		proxyURL:       proxyURL,
		accessToken:    accessToken,
		outStream:      os.Stdout,
		verboseEnabled: verboseEnabled,
		proxyUsername:  proxyUsername,
		proxyPassword:  proxyPassword,
		connectTimeout: time.Duration(connectTimeout) * time.Second,
		readTimeout:    time.Duration(readTimeout) * time.Second,
		writeTimeout:   time.Duration(writeTimeout) * time.Second,
		callTimeout:    time.Duration(callTimeout) * time.Second,
		maxRetries:     maxRetries,
	}
}

func (c *centralAPIClientImpl) getHTTPClient() *http.Client {
	if c.httpClient != nil {
		return c.httpClient
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	if c.proxyURL != nil {
		transport.Proxy = http.ProxyURL(c.proxyURL)
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   c.callTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	c.httpClient = client
	return client
}

func (c *centralAPIClientImpl) newRequest(method, urlStr, supportedPlatform, ballerinaVersion string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(BallerinaPlatform, supportedPlatform)
	req.Header.Set(UserAgent, ballerinaVersion)
	req.Header.Set(BallerinaCentralTelemetryDisabled, strconv.FormatBool(os.Getenv(TestModeActive) == "true"))

	if c.accessToken != "" {
		req.Header.Set(Authorization, getBearerToken(c.accessToken))
	}

	return req, nil
}

func (c *centralAPIClientImpl) GetPackage(orgNamePath, packageNamePath, version, supportedPlatform, ballerinaVersion string) (*model.Package, error) {
	packageSignature := orgNamePath + Separator + packageNamePath + ":" + version
	resourceURL := PackagePathPrefix + orgNamePath + Separator + packageNamePath

	urlStr := c.baseURL + resourceURL
	if version != "" {
		urlStr = urlStr + Separator + version
	}

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotFindPackage+packageSignature, err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotFindPackage+packageSignature, err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, resourceURL)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotFindPackage+packageSignature, err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var pkg model.Package
			if err := json.Unmarshal(bodyBytes, &pkg); err != nil {
				return nil, NewCentralClientErrorWithCause(ErrCannotFindPackage+packageSignature, err)
			}
			return &pkg, nil

		case http.StatusNotFound:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
				if strings.Contains(errResp.Message, "package not found for:") {
					return nil, NewNoPackageError(errResp.Message)
				}
				return nil, NewCentralClientError(ErrCannotFindPackage + packageSignature + ". reason: " + errResp.Message)
			}

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponseWithOrg(orgNamePath, bodyBytes)

		case http.StatusBadRequest, http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(errResp.Message)
			}
		}
	}

	return nil, NewCentralClientError(ErrCannotFindPackage + packageSignature)
}

func (c *centralAPIClientImpl) GetPackageVersions(orgNamePath, packageNamePath, supportedPlatform, ballerinaVersion string) ([]string, error) {
	packageSignature := orgNamePath + Separator + packageNamePath
	resourceURL := PackagePathPrefix + orgNamePath + Separator + packageNamePath
	urlStr := c.baseURL + resourceURL

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotFindVersions+packageSignature, err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotFindVersions+packageSignature, err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, resourceURL)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotFindVersions+packageSignature, err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var versions []string
			if err := json.Unmarshal(bodyBytes, &versions); err != nil {
				return nil, NewCentralClientErrorWithCause(ErrCannotFindVersions+packageSignature, err)
			}
			return versions, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponseWithOrg(orgNamePath, bodyBytes)

		case http.StatusNotFound:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
				if strings.Contains(errResp.Message, "package not found") {
					return []string{}, nil
				}
				return nil, NewCentralClientError(ErrCannotFindVersions + packageSignature + ". reason: " + errResp.Message)
			}

		case http.StatusBadRequest, http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
				return nil, NewCentralClientError(ErrCannotFindVersions + packageSignature + ". reason: " + errResp.Message)
			}
		}
	}

	return nil, NewCentralClientError(ErrCannotFindVersions + packageSignature + ".")
}

func (c *centralAPIClientImpl) PushPackage(balaPath, org, name, version, supportedPlatform, ballerinaVersion string) error {
	enableOutputStream := os.Getenv(EnableOutputStream) == "true"
	packageSignature := org + Separator + name + ":" + version
	urlStr := c.baseURL + Separator + "packages"

	fileName := org + "-" + name + "-" + version + ".bala"

	var body bytes.Buffer
	writer := io.MultiWriter(&body)

	balaFile, err := os.Open(balaPath)
	if err != nil {
		return NewCentralClientErrorWithCause(ErrCannotPush+"'"+packageSignature+"'", err)
	}
	defer balaFile.Close()

	if _, err := io.Copy(writer, balaFile); err != nil {
		return NewCentralClientErrorWithCause(ErrCannotPush+"'"+packageSignature+"'", err)
	}

	digestVal := SHA256 + checkHashInternal(balaPath)

	req, err := c.newRequest(http.MethodPost, urlStr, supportedPlatform, ballerinaVersion, &body)
	if err != nil {
		return NewCentralClientErrorWithCause(ErrCannotPush+"'"+packageSignature+"'", err)
	}

	req.Header.Set(ContentType, ApplicationOctetStream)
	req.Header.Set(Digest, digestVal)
	req.Header.Set(ContentDisposition, "attachment; filename="+fileName)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return NewCentralClientErrorWithCause(ErrCannotPush+"'"+packageSignature+"'", err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, Separator+"packages")

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewCentralClientErrorWithCause(ErrCannotPush+"'"+packageSignature+"'", err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	if resp.StatusCode == http.StatusNoContent {
		if enableOutputStream {
			fmt.Fprintln(c.outStream, packageSignature+" pushed to central successfully")
		}
		return nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return c.handleUnauthorizedResponseWithOrg(org, bodyBytes)
	}

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				if strings.Contains(errResp.Message, "subject claims missing in the user info repsonse") {
					errResp.Message = "unauthorized access token for organization: '" + org + "'. check access token set in 'Settings.toml' file."
				}
				return NewCentralClientError(errResp.Message)
			}

		case http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return NewCentralClientError(ErrCannotPush + "'" + packageSignature + "' reason:" + errResp.Message)
			}
		}
	}

	return NewCentralClientError(ErrCannotPush + "'" + packageSignature + "' to the remote repository '" + urlStr + "'.")
}

func (c *centralAPIClientImpl) PullPackage(org, name, version, packagePathInBalaCache, supportedPlatform, ballerinaVersion string, isBuild bool) error {
	retryCount := 0
	for retryCount <= c.maxRetries {
		err := c.pullPackageInternal(org, name, version, packagePathInBalaCache, supportedPlatform, ballerinaVersion, isBuild)
		if err != nil {
			if strings.Contains(err.Error(), ConnectionReset) && retryCount < c.maxRetries {
				if c.verboseEnabled {
					fmt.Fprintf(c.outStream, "* Retrying to pull the package: %s/%s:%s due to: %s. Retry attempt: %d\n\n",
						org, name, version, err.Error(), retryCount+1)
				}
				retryCount++
				continue
			}
			return err
		}
		break
	}
	return nil
}

func (c *centralAPIClientImpl) pullPackageInternal(org, name, version, packagePathInBalaCache, supportedPlatform, ballerinaVersion string, isBuild bool) error {
	resourceURL := PackagePathPrefix + org + Separator + name
	enableOutputStream := os.Getenv(EnableOutputStream) == "true"
	packageSignature := org + Separator + name
	urlStr := c.baseURL + resourceURL

	if version != "" {
		urlStr += Separator + version
		packageSignature += ":" + version
	} else {
		urlStr += "/*"
		packageSignature += ":*"
	}

	var logFormatter LogFormatter
	if isBuild {
		logFormatter = NewBuildLogFormatter()
	} else {
		logFormatter = NewLogFormatter()
	}

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+packageSignature+"'"), err)
	}

	req.Header.Set(AcceptEncoding, Identity)
	req.Header.Set(Accept, ApplicationOctetStream)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+packageSignature+"'"), err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, resourceURL)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+packageSignature+"'"), err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	if resp.StatusCode == http.StatusFound {
		balaURL := resp.Header.Get(Location)
		balaFileName := resp.Header.Get(ContentDisposition)
		deprecationFlag := resp.Header.Get("Is-Deprecated")
		deprecationMsg := resp.Header.Get("Deprecate-Message")
		digest := resp.Header.Get(Digest)

		if digest == "" {
			digest = ""
		}

		isDeprecated := deprecationFlag == "true"
		deprecationMessage := deprecationMsg

		if !isBuild && isDeprecated {
			fmt.Fprintf(c.outStream, "WARNING [%s] %s is deprecated: %s\n", name, packageSignature, deprecationMessage)
		}

		if balaURL != "" && balaFileName != "" {
			downloadReq, err := c.newRequest(http.MethodGet, balaURL, supportedPlatform, ballerinaVersion, nil)
			if err != nil {
				return NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+packageSignature+"'"), err)
			}

			downloadReq.Header.Set(AcceptEncoding, Identity)
			downloadReq.Header.Set(ContentDisposition, balaFileName)

			c.logRequestInitVerbose(downloadReq)

			downloadResp, err := client.Do(downloadReq)
			if err != nil {
				return NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+packageSignature+"'"), err)
			}
			defer downloadResp.Body.Close()

			c.logRequestConnectVerbose(downloadReq, balaURL)
			c.logResponseVerbose(downloadResp, "")

			if downloadResp.StatusCode == http.StatusOK {
				isNightlyBuild := strings.Contains(ballerinaVersion, "SNAPSHOT")
				var outStream io.Writer
				if enableOutputStream {
					outStream = c.outStream
				}

				deprecMsg := ""
				if isDeprecated {
					deprecMsg = deprecationMessage
				}

				return createBalaInHomeRepo(downloadResp, packagePathInBalaCache, org, name, isNightlyBuild,
					deprecMsg, balaURL, balaFileName, outStream, logFormatter, digest)
			}

			errorMessage := logFormatter.FormatLog(ErrCannotPullPackage + "'" + packageSignature +
				"'. BALA content download from '" + balaURL + "' failed.")
			return c.handleResponseErrors(downloadResp, errorMessage, bodyBytes)
		}

		errorMsg := logFormatter.FormatLog(ErrCannotPullPackage + "'" + packageSignature +
			"' from the remote repository '" + urlStr + "'. reason: bala file location is missing.")
		return NewCentralClientError(errorMsg)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return c.handleUnauthorizedResponseWithOrg(org, bodyBytes)
	}

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return NewCentralClientError("error: " + errResp.Message)
			}
		}

		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusServiceUnavailable {
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				errorMsg := logFormatter.FormatLog(ErrCannotPullPackage + "'" + packageSignature +
					"' from the remote repository '" + urlStr + "'. reason: " + errResp.Message)
				return NewCentralClientError(errorMsg)
			}
		}
	}

	errorMsg := logFormatter.FormatLog(ErrCannotPullPackage + "'" + packageSignature +
		"' from the remote repository '" + urlStr + "'.")
	return NewCentralClientError(errorMsg)
}

func (c *centralAPIClientImpl) PullTool(toolID, version, balaCacheDirPath, supportedPlatform, ballerinaVersion string, isBuild bool) ([]string, error) {
	retryCount := 0
	for retryCount <= c.maxRetries {
		result, err := c.pullToolInternal(toolID, version, balaCacheDirPath, supportedPlatform, ballerinaVersion, isBuild)
		if err != nil {
			if strings.Contains(err.Error(), ConnectionReset) && retryCount < c.maxRetries {
				if c.verboseEnabled {
					fmt.Fprintf(c.outStream, "* Retrying to pull the tool: %s:%s due to: %s. Retry attempt: %d\n\n",
						toolID, version, err.Error(), retryCount+1)
				}
				retryCount++
				continue
			}
			return nil, err
		}
		return result, nil
	}
	return nil, NewCentralClientError(fmt.Sprintf("Failed to pull the tool: %s:%s after %d attempts.", toolID, version, c.maxRetries))
}

func (c *centralAPIClientImpl) pullToolInternal(toolID, version, balaCacheDirPath, supportedPlatform, ballerinaVersion string, isBuild bool) ([]string, error) {
	resourceURL := ToolPathPrefix + toolID
	enableOutputStream := os.Getenv(EnableOutputStream) == "true"
	toolSignature := toolID
	urlStr := c.baseURL + resourceURL

	if version != "" {
		urlStr += Separator + version
		toolSignature += ":" + version
	}

	var logFormatter LogFormatter
	if isBuild {
		logFormatter = NewBuildLogFormatter()
	} else {
		logFormatter = NewLogFormatter()
	}

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+toolSignature+"'"), err)
	}

	req.Header.Set(ContentType, MediaTypeJSONContent)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+toolSignature+"'"), err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, resourceURL)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+toolSignature+"'"), err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	if resp.StatusCode == http.StatusOK {
		contentType := resp.Header.Get(ContentType)
		if isApplicationJSONContentType(contentType) {
			var jsonContent map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &jsonContent); err != nil {
				return nil, NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+toolSignature+"'"), err)
			}

			org, _ := jsonContent[Organization].(string)
			pkgName, _ := jsonContent[PkgName].(string)
			latestVersion, _ := jsonContent[Version].(string)
			balaURL, _ := jsonContent[BalaURL].(string)
			platform, ok := jsonContent[Platform].(string)
			if !ok {
				platform = AnyPlatform
			}
			digest, _ := jsonContent[Digest].(string)

			if balaURL != "" && org != "" && latestVersion != "" && pkgName != "" {
				balaFileName := "attachment; filename=" + org + "-" + pkgName + "-" + platform + "-" + latestVersion + ".bala"

				downloadReq, err := c.newRequest(http.MethodGet, balaURL, supportedPlatform, ballerinaVersion, nil)
				if err != nil {
					return nil, NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+toolSignature+"'"), err)
				}

				downloadReq.Header.Set(AcceptEncoding, Identity)
				downloadReq.Header.Set(ContentDisposition, balaFileName)

				c.logRequestInitVerbose(downloadReq)

				downloadResp, err := client.Do(downloadReq)
				if err != nil {
					return nil, NewCentralClientErrorWithCause(logFormatter.FormatLog(ErrCannotPullPackage+"'"+toolSignature+"'"), err)
				}
				defer downloadResp.Body.Close()

				c.logRequestConnectVerbose(downloadReq, balaURL)
				c.logResponseVerbose(downloadResp, "")

				packagePathInBalaCache := filepath.Join(balaCacheDirPath, org, pkgName)

				if downloadResp.StatusCode == http.StatusOK {
					isNightlyBuild := strings.Contains(ballerinaVersion, "SNAPSHOT")
					var outStream io.Writer
					if enableOutputStream {
						outStream = c.outStream
					}

					err := createBalaInHomeRepo(downloadResp, packagePathInBalaCache, org, pkgName, isNightlyBuild,
						"", balaURL, balaFileName, outStream, logFormatter, digest)

					if err != nil {
						if pkgExistsErr, ok := err.(*PackageAlreadyExistsError); ok {
							if enableOutputStream {
								fmt.Fprintf(c.outStream, "tool '%s:%s' is already available locally.\n", toolID, latestVersion)
							}
							return []string{org, pkgName, pkgExistsErr.Version()}, nil
						}
						return nil, err
					}

					return []string{org, pkgName, latestVersion}, nil
				}

				errorMessage := logFormatter.FormatLog(ErrCannotPullPackage + "'" + toolSignature +
					"'. BALA content download from '" + balaURL + "' failed.")
				return nil, c.handleResponseErrors(downloadResp, errorMessage, bodyBytes)
			}

			errorMsg := logFormatter.FormatLog(ErrCannotPullPackage + "'" + toolSignature +
				"' from the remote repository '" + urlStr + "'. reason: bala file location is missing.")
			return nil, NewCentralClientError(errorMsg)
		}
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, c.handleUnauthorizedResponse(bodyBytes)
	}

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError("error: " + errResp.Message)
			}
		}

		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusServiceUnavailable {
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				errorMsg := logFormatter.FormatLog(ErrCannotPullPackage + "'" + toolSignature +
					"' from the remote repository '" + urlStr + "'. reason: " + errResp.Message)
				return nil, NewCentralClientError(errorMsg)
			}
		}
	}

	errorMsg := logFormatter.FormatLog(ErrCannotPullPackage + "'" + toolSignature +
		"' from the remote repository '" + urlStr + "'.")
	return nil, NewCentralClientError(errorMsg)
}

func (c *centralAPIClientImpl) ResolvePackageNames(request *model.PackageNameResolutionRequest, supportedPlatform, ballerinaVersion string) (*model.PackageNameResolutionResponse, error) {
	urlStr := c.baseURL + PackagePathPrefix + ResolveModules

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrPackageResolution, err)
	}

	req, err := c.newRequest(http.MethodPost, urlStr, supportedPlatform, ballerinaVersion, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrPackageResolution, err)
	}

	req.Header.Set(ContentType, MediaTypeJSON)
	req.Header.Set(Accept, MediaTypeJSONContent)
	req.Header.Set(AcceptEncoding, Identity)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, PackagePathPrefix+ResolveModules)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
	}

	c.logResponseVerbose(resp, string(respBodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var response model.PackageNameResolutionResponse
			if err := json.Unmarshal(respBodyBytes, &response); err != nil {
				return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
			}
			return &response, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponse(respBodyBytes)

		case http.StatusBadRequest:
			var errResp model.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewConnectionError(errResp.Message)
			}

		case http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp model.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewConnectionError(ErrPackageResolution + " reason:" + errResp.Message)
			}
		}
	}

	return nil, NewConnectionError(ErrPackageResolution)
}

func (c *centralAPIClientImpl) ResolveDependencies(request *model.PackageResolutionRequest, supportedPlatform, ballerinaVersion string) (*model.PackageResolutionResponse, error) {
	urlStr := c.baseURL + PackagePathPrefix + ResolveDependencies

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrPackageResolution, err)
	}

	req, err := c.newRequest(http.MethodPost, urlStr, supportedPlatform, ballerinaVersion, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrPackageResolution, err)
	}

	req.Header.Set(ContentType, MediaTypeJSON)
	req.Header.Set(Accept, MediaTypeJSONContent)
	req.Header.Set(AcceptEncoding, Identity)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, PackagePathPrefix+ResolveDependencies)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
	}

	c.logResponseVerbose(resp, string(respBodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var response model.PackageResolutionResponse
			if err := json.Unmarshal(respBodyBytes, &response); err != nil {
				return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
			}
			return &response, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponse(respBodyBytes)

		case http.StatusBadRequest:
			var errResp model.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewConnectionError(errResp.Message)
			}

		case http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp model.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewConnectionError(ErrPackageResolution + " reason:" + errResp.Message)
			}
		}
	}

	return nil, NewConnectionError(ErrPackageResolution)
}

func (c *centralAPIClientImpl) ResolveToolDependencies(request *model.ToolResolutionCentralRequest, supportedPlatform, ballerinaVersion string) (*model.ToolResolutionCentralResponse, error) {
	urlStr := c.baseURL + ToolPathPrefix + ResolveDependencies

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrPackageResolution, err)
	}

	req, err := c.newRequest(http.MethodPost, urlStr, supportedPlatform, ballerinaVersion, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrPackageResolution, err)
	}

	req.Header.Set(ContentType, MediaTypeJSON)
	req.Header.Set(Accept, MediaTypeJSONContent)
	req.Header.Set(AcceptEncoding, Identity)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, ToolPathPrefix+ResolveDependencies)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
	}

	c.logResponseVerbose(resp, string(respBodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var response model.ToolResolutionCentralResponse
			if err := json.Unmarshal(respBodyBytes, &response); err != nil {
				return nil, NewConnectionErrorWithCause(ErrPackageResolution, err)
			}
			return &response, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponse(respBodyBytes)

		case http.StatusBadRequest:
			var errResp model.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewConnectionError(errResp.Message)
			}

		case http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp model.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewConnectionError(ErrPackageResolution + " reason:" + errResp.Message)
			}
		}
	}

	return nil, NewConnectionError(ErrPackageResolution)
}

func (c *centralAPIClientImpl) SearchPackage(query, supportedPlatform, ballerinaVersion string) (*model.PackageSearchResult, error) {
	urlStr := c.baseURL + Separator + "packages" + "/?q=" + url.QueryEscape(query)

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotSearch+"'"+query+"'", err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotSearch+"'"+query+"'", err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, Separator+"packages"+"/?q="+query)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotSearch+"'"+query+"'", err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var result model.PackageSearchResult
			if err := json.Unmarshal(bodyBytes, &result); err != nil {
				return nil, NewCentralClientErrorWithCause(ErrCannotSearch+"'"+query+"'", err)
			}
			return &result, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponse(bodyBytes)

		case http.StatusBadRequest:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(errResp.Message)
			}

		case http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(ErrCannotSearch + "'" + query + "' reason:" + errResp.Message)
			}
		}
	}

	return nil, NewCentralClientError(ErrCannotSearch + "'" + query + "'.")
}

func (c *centralAPIClientImpl) SearchTool(keyword, supportedPlatform, ballerinaVersion string) (*model.ToolSearchResult, error) {
	urlStr := c.baseURL + Separator + "tools" + SearchQueryParam + url.QueryEscape(keyword)

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotSearch+"'"+keyword+"'", err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotSearch+"'"+keyword+"'", err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, ToolPathPrefix+SearchQueryParam+Separator+keyword)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotSearch+"'"+keyword+"'", err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var result model.ToolSearchResult
			if err := json.Unmarshal(bodyBytes, &result); err != nil {
				return nil, NewCentralClientErrorWithCause(ErrCannotSearch+"'"+keyword+"'", err)
			}
			return &result, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponse(bodyBytes)

		case http.StatusBadRequest:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(errResp.Message)
			}

		case http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(ErrCannotSearch + "'" + keyword + "' reason:" + errResp.Message)
			}
		}
	}

	return nil, NewCentralClientError(ErrCannotSearch + "'" + keyword + "'.")
}

func (c *centralAPIClientImpl) DeprecatePackage(packageInfo, deprecationMsg, supportedPlatform, ballerinaVersion string, isUndo bool) error {
	parts := strings.Split(packageInfo, Separator)
	if len(parts) < 2 {
		return NewCentralClientError("invalid package info format")
	}

	nameParts := strings.Split(parts[1], ":")
	if len(nameParts) < 2 {
		return NewCentralClientError("invalid package info format")
	}

	existingPackage, err := c.GetPackage(parts[0], nameParts[0], nameParts[1], supportedPlatform, ballerinaVersion)
	if err != nil {
		return err
	}

	packageValue := packageInfo
	if strings.HasSuffix(packageInfo, ":*") {
		packageValue = packageInfo[:len(packageInfo)-2]
	}

	deprecated := false
	if existingPackage.IsDeprecated != nil {
		deprecated = *existingPackage.IsDeprecated
	}

	if isUndo && !deprecated {
		fmt.Fprintf(c.outStream, "package %s is not marked as deprecated in central\n", packageValue)
		return nil
	}

	var requestBody map[string]string
	var requestURL string

	if isUndo {
		requestURL = c.baseURL + PackagePathPrefix + Undeprecate + Separator + strings.ReplaceAll(packageInfo, ":", Separator)
		requestBody = map[string]string{}
	} else {
		requestBody = map[string]string{"message": deprecationMsg}
		requestURL = c.baseURL + PackagePathPrefix + Deprecate + Separator + strings.ReplaceAll(packageInfo, ":", Separator)
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return NewCentralClientErrorWithCause(ErrPackageDeprecate+"'"+packageValue+"'", err)
	}

	req, err := c.newRequest(http.MethodPut, requestURL, supportedPlatform, ballerinaVersion, bytes.NewReader(bodyBytes))
	if err != nil {
		errorMsg := ErrPackageDeprecate
		if isUndo {
			errorMsg = ErrPackageUndeprecate
		}
		return NewCentralClientErrorWithCause(errorMsg+"'"+packageValue+"'", err)
	}

	req.Header.Set(ContentType, MediaTypeJSON)
	req.Header.Set(Accept, MediaTypeJSONContent)
	req.Header.Set(AcceptEncoding, Identity)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		errorMsg := ErrPackageDeprecate
		if isUndo {
			errorMsg = ErrPackageUndeprecate
		}
		return NewCentralClientErrorWithCause(errorMsg+"'"+packageValue+"'", err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, strings.ReplaceAll(requestURL, c.baseURL, ""))

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMsg := ErrPackageDeprecate
		if isUndo {
			errorMsg = ErrPackageUndeprecate
		}
		return NewCentralClientErrorWithCause(errorMsg+"'"+packageValue+"'", err)
	}

	c.logResponseVerbose(resp, string(respBodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		if resp.StatusCode == http.StatusOK {
			var pkgResp model.Package
			if err := json.Unmarshal(respBodyBytes, &pkgResp); err == nil {
				pkgDeprecated := false
				if pkgResp.IsDeprecated != nil {
					pkgDeprecated = *pkgResp.IsDeprecated
				}

				if pkgDeprecated {
					if deprecated {
						fmt.Fprintf(c.outStream, "deprecation message is successfully updated for the package %s in central\n", packageValue)
					} else {
						fmt.Fprintf(c.outStream, "package %s marked as deprecated in central successfully\n", packageValue)
					}
				} else {
					fmt.Fprintf(c.outStream, "deprecation of the package %s is successfully undone in central\n", packageValue)
				}
				return nil
			}
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return c.handleUnauthorizedResponse(respBodyBytes)
		}

		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
			var errResp model.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return NewCentralClientError(errResp.Message)
			}
		}

		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusServiceUnavailable {
			var errResp model.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				errorMsg := ErrPackageDeprecate
				if isUndo {
					errorMsg = ErrPackageUndeprecate
				}
				return NewCentralClientError(errorMsg + "'" + packageValue + "' reason:" + errResp.Message)
			}
		}
	}

	errorMsg := ErrPackageDeprecate
	if isUndo {
		errorMsg = ErrPackageUndeprecate
	}
	return NewCentralClientError(errorMsg + "'" + packageValue + "'.")
}

func (c *centralAPIClientImpl) GetConnectors(params map[string]string, supportedPlatform, ballerinaVersion string) (interface{}, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector, err)
	}

	baseURL.Path = filepath.Join(baseURL.Path, "connectors")
	query := baseURL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	baseURL.RawQuery = query.Encode()

	req, err := c.newRequest(http.MethodGet, baseURL.String(), supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector, err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector, err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, Separator+"connectors")

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector, err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var result interface{}
		if err := json.Unmarshal(bodyBytes, &result); err == nil {
			return result, nil
		}
	}

	return nil, c.handleResponseErrors(resp, ErrCannotGetConnector, bodyBytes)
}

func (c *centralAPIClientImpl) GetConnector(id, supportedPlatform, ballerinaVersion string) (map[string]interface{}, error) {
	urlStr := c.baseURL + ConnectorPathPrefix + id

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector+"'"+id+"'", err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector+"'"+id+"'", err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, ConnectorPathPrefix+id)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector+"'"+id+"'", err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &result); err == nil {
			return result, nil
		}
	}

	return nil, c.handleResponseErrors(resp, ErrCannotGetConnector+" id:"+id, bodyBytes)
}

func (c *centralAPIClientImpl) GetConnectorByInfo(connector *model.ConnectorInfo, supportedPlatform, ballerinaVersion string) (map[string]interface{}, error) {
	urlStr := c.baseURL + ConnectorPathPrefix + connector.OrgName + Separator +
		connector.PackageName + Separator + connector.Version + Separator +
		connector.ModuleName + Separator + connector.Name

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector+"'"+connector.PackageName+"'", err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector+"'"+connector.PackageName+"'", err)
	}
	defer resp.Body.Close()

	resourcePath := ConnectorPathPrefix + connector.OrgName + Separator +
		connector.PackageName + Separator + connector.Version + Separator +
		connector.ModuleName + Separator + connector.Name
	c.logRequestConnectVerbose(req, resourcePath)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetConnector+"'"+connector.PackageName+"'", err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &result); err == nil {
			return result, nil
		}
	}

	return nil, c.handleResponseErrors(resp, ErrCannotGetConnector+" "+connector.PackageName, bodyBytes)
}

func (c *centralAPIClientImpl) GetTriggers(params map[string]string, supportedPlatform, ballerinaVersion string) (interface{}, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetTriggers, err)
	}

	baseURL.Path = filepath.Join(baseURL.Path, "triggers")
	query := baseURL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	baseURL.RawQuery = query.Encode()

	req, err := c.newRequest(http.MethodGet, baseURL.String(), supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetTriggers, err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetTriggers, err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, Separator+"triggers")

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetTriggers, err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var result interface{}
		if err := json.Unmarshal(bodyBytes, &result); err == nil {
			return result, nil
		}
	}

	return nil, c.handleResponseErrors(resp, ErrCannotGetTriggers, bodyBytes)
}

func (c *centralAPIClientImpl) GetTrigger(id, supportedPlatform, ballerinaVersion string) (map[string]interface{}, error) {
	urlStr := c.baseURL + TriggerPathPrefix + id

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetTrigger+"'"+id+"'", err)
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetTrigger+"'"+id+"'", err)
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, TriggerPathPrefix+id)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientErrorWithCause(ErrCannotGetTrigger+"'"+id+"'", err)
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &result); err == nil {
			return result, nil
		}
	}

	return nil, c.handleResponseErrors(resp, ErrCannotGetTrigger+" id:"+id, bodyBytes)
}

func (c *centralAPIClientImpl) AccessToken() string {
	return c.accessToken
}

func (c *centralAPIClientImpl) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *centralAPIClientImpl) handleUnauthorizedResponse(bodyBytes []byte) error {
	var errResp model.Error
	if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
		return NewCentralClientError("unauthorized access token. " +
			"check access token set in 'Settings.toml' file. reason: " + errResp.Message)
	}
	return NewCentralClientError("unauthorized access token. " +
		"check access token set in 'Settings.toml' file.")
}

func (c *centralAPIClientImpl) handleResponseErrors(resp *http.Response, msg string, bodyBytes []byte) error {
	contentType := resp.Header.Get(ContentType)

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		if isApplicationJSONContentType(contentType) {
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return NewCentralClientError(errResp.Message)
			}
		}
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if isApplicationJSONContentType(contentType) {
			return c.handleUnauthorizedResponse(bodyBytes)
		}
		return NewCentralClientError("unauthorized access token. check access token set in 'Settings.toml' file.")
	}

	if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusGatewayTimeout {
		if isApplicationJSONContentType(contentType) {
			var errResp model.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return NewCentralClientError(msg + " reason:" + errResp.Message)
			}
		}
	}

	return NewCentralClientError(msg)
}

func (c *centralAPIClientImpl) handleUnauthorizedResponseWithOrg(org string, bodyBytes []byte) error {
	var errResp model.Error
	if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
		return NewCentralClientError("unauthorized access token for organization: '" + org +
			"'. check access token set in 'Settings.toml' file. reason: " + errResp.Message)
	}
	return NewCentralClientError("unauthorized access token for organization: '" + org +
		"'. check access token set in 'Settings.toml' file.")
}

func (c *centralAPIClientImpl) logRequestInitVerbose(req *http.Request) {
	if c.verboseEnabled {
		fmt.Fprintf(c.outStream, "* Trying %s\n", req.URL.String())
	}
}

func (c *centralAPIClientImpl) logRequestConnectVerbose(req *http.Request, resourceURL string) {
	if c.verboseEnabled {
		fmt.Fprintf(c.outStream, "* Connected to %s\n", c.baseURL)
		fmt.Fprintf(c.outStream, "> %s %s HTTP\n", req.Method, resourceURL)
		fmt.Fprintf(c.outStream, "> Host: %s\n", c.baseURL)
		for name, values := range req.Header {
			for _, value := range values {
				if name == "Authorization" {
					fmt.Fprintf(c.outStream, "> %s: Bearer ************************************\n", name)
				} else {
					fmt.Fprintf(c.outStream, "> %s: %s\n", name, value)
				}
			}
		}
		fmt.Fprintln(c.outStream, ">")
	}
}

func (c *centralAPIClientImpl) logResponseVerbose(resp *http.Response, bodyContent string) {
	if c.verboseEnabled {
		fmt.Fprintf(c.outStream, "< HTTP %s\n", resp.Status)

		for name, values := range resp.Header {
			for _, value := range values {
				fmt.Fprintf(c.outStream, "> %s: %s\n", name, value)
			}
		}
		fmt.Fprintln(c.outStream, "< ")
		if bodyContent != "" {
			fmt.Fprintln(c.outStream, bodyContent)
		}
		fmt.Fprintf(c.outStream, "* Connection to host %s left intact \n\n", c.baseURL)
	}
}

func getBearerToken(accessToken string) string {
	return "Bearer " + accessToken
}

func isApplicationJSONContentType(contentType string) bool {
	return strings.HasPrefix(contentType, MediaTypeJSONContent)
}

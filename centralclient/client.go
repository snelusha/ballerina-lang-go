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

	"ballerina-lang-go/centralclient/models"
)

type CentralAPIClient interface {
	GetPackage(orgNamePath, packageNamePath, version, supportedPlatform, ballerinaVersion string) (*models.Package, error)
	GetPackageVersions(orgNamePath, packageNamePath, supportedPlatform, ballerinaVersion string) ([]string, error)
	PullPackage(org, name, versiom, packagePathInBalaCache, supportedPlatform, ballerinaVersion string, isBuild bool) error
	ResolvePackageNames(request *models.PackageNameResolutionRequest, supportedPlatform, ballerinaVersion string) (*models.PackageNameResolutionResponse, error)
	ResolveDependencies(request *models.PackageResolutionRequest, supportedPlatform, ballerinaVersion string) (*models.PackageResolutionResponse, error)
	GetConnectors(params map[string]string, supportedPlatform, ballerinaVersion string) (any, error)
	GetConnector(id, supportedPlatform, ballerinaVersion string) (map[string]any, error)
	GetConnectorByInfo(connector *models.ConnectorInfo, supportedPlatform, ballerinaVersion string) (map[string]any, error)
	GetTriggers(params map[string]string, supportedPlatform, vallerinaVersion string) (any, error)
	GetTrigger(id, supportedPlatform, ballerinaVersion string) (map[string]any, error)
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

func (c *centralAPIClientImpl) GetPackage(orgNamePath, packageNamePath, version, supportedPlatform, ballerinaVersion string) (*models.Package, error) {
	packageSignature := fmt.Sprintf("%s%s%s:%s", orgNamePath, Separator, packageNamePath, version)
	resourceURL := fmt.Sprintf("%s%s%s%s", PackagePathPrefix, orgNamePath, Separator, packageNamePath)

	urlStr := fmt.Sprintf("%s%s", c.baseURL, resourceURL)
	if version != "" {
		urlStr = fmt.Sprintf("%s/%s", urlStr, version)
	}

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindPackage, packageSignature, err.Error()))
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindPackage, packageSignature, err.Error()))
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, resourceURL)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindPackage, packageSignature, err.Error()))
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var pkg models.Package
			if err := json.Unmarshal(bodyBytes, &pkg); err != nil {
				return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindPackage, packageSignature, err.Error()))
			}
			return &pkg, nil

		case http.StatusNotFound:
			var errResp models.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
				if strings.Contains(errResp.Message, "package not found for:") {
					return nil, NewNoPackageError(errResp.Message)
				}
				return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindPackage, packageSignature, err.Error()))
			}

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponseWithOrg(orgNamePath, bodyBytes)

		case http.StatusBadRequest, http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp models.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(errResp.Message)
			}
		}
	}

	return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotFindPackage, packageSignature))
}

func (c *centralAPIClientImpl) GetPackageVersions(orgNamePath, packageNamePath, supportedPlatform, ballerinaVersion string) ([]string, error) {
	packageSignature := fmt.Sprintf("%s%s%s", orgNamePath, Separator, packageNamePath)
	resourceURL := fmt.Sprintf("%s%s%s%s", PackagePathPrefix, orgNamePath, Separator, packageNamePath)

	urlStr := fmt.Sprintf("%s%s", c.baseURL, resourceURL)

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindVersions, packageSignature, err.Error()))
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindVersions, packageSignature, err.Error()))
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, resourceURL)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindVersions, packageSignature, err.Error()))
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var versions []string
			if err := json.Unmarshal(bodyBytes, &versions); err != nil {
				return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindVersions, packageSignature, err.Error()))
			}
			return versions, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponseWithOrg(orgNamePath, bodyBytes)

		case http.StatusNotFound:
			var errResp models.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
				if strings.Contains(errResp.Message, "package not found for:") {
					return []string{}, nil
				}
				return nil, NewCentralClientError(fmt.Sprintf("%s%s. reason: %s", ErrCannotFindVersions, packageSignature, err.Error()))
			}

		case http.StatusBadRequest, http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp models.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(errResp.Message)
			}
		}
	}

	return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotFindVersions, packageSignature))
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
				retryCount = retryCount + 1
				continue
			}
			return err
		}
		break
	}
	return nil
}

func (c *centralAPIClientImpl) ResolvePackageNames(request *models.PackageNameResolutionRequest, supportedPlatform, ballerinaVersion string) (*models.PackageNameResolutionResponse, error) {
	urlStr := fmt.Sprintf("%s%s%s", c.baseURL, PackagePathPrefix, ResolveModules)

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
	}

	req, err := c.newRequest(http.MethodPost, urlStr, supportedPlatform, ballerinaVersion, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
	}

	req.Header.Set(ContentType, MediaTypeJSON)
	req.Header.Set(Accept, MediaTypeJSONContent)
	req.Header.Set(AcceptEncoding, Identity)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, fmt.Sprintf("%s%s", PackagePathPrefix, ResolveModules))

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
	}

	c.logResponseVerbose(resp, string(respBodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var resolutionResponse models.PackageNameResolutionResponse
			if err := json.Unmarshal(respBodyBytes, &resolutionResponse); err != nil {
				return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
			}
			return &resolutionResponse, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponse(bodyBytes)

		case http.StatusBadRequest:
			var errResp models.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(errResp.Message)
			}

		case http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp models.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(fmt.Sprintf("%s reason: %s", ErrPackageResolution, errResp.Message))
			}
		}
	}

	return nil, NewConnectionError(ErrPackageResolution)
}

func (c *centralAPIClientImpl) ResolveDependencies(request *models.PackageResolutionRequest, supportedPlatform, ballerinaVersion string) (*models.PackageResolutionResponse, error) {
	urlStr := fmt.Sprintf("%s%s%s", c.baseURL, PackagePathPrefix, ResolveDependencies)

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
	}

	req, err := c.newRequest(http.MethodPost, urlStr, supportedPlatform, ballerinaVersion, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
	}

	req.Header.Set(AcceptEncoding, Identity)
	req.Header.Set(Accept, ApplicationJSON)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, fmt.Sprintf("%s%s", PackagePathPrefix, ResolveDependencies))

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
	}

	c.logResponseVerbose(resp, string(respBodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		switch resp.StatusCode {
		case http.StatusOK:
			var resolutionResponse models.PackageResolutionResponse
			if err := json.Unmarshal(respBodyBytes, &resolutionResponse); err != nil {
				return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrPackageResolution, err.Error()))
			}
			return &resolutionResponse, nil

		case http.StatusUnauthorized:
			return nil, c.handleUnauthorizedResponse(bodyBytes)

		case http.StatusBadRequest:
			var errResp models.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(errResp.Message)
			}

		case http.StatusInternalServerError, http.StatusServiceUnavailable:
			var errResp models.Error
			if err := json.Unmarshal(respBodyBytes, &errResp); err == nil && errResp.Message != "" {
				return nil, NewCentralClientError(fmt.Sprintf("%s reason: %s", ErrPackageResolution, errResp.Message))
			}
		}
	}

	return nil, NewConnectionError(ErrPackageResolution)
}

func (c *centralAPIClientImpl) GetConnectors(params map[string]string, supportedPlatform, ballerinaVersion string) (any, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	baseURL.Path = filepath.Join(baseURL.Path, ConnectorsPath)
	query := baseURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	baseURL.RawQuery = query.Encode()

	req, err := c.newRequest(http.MethodGet, baseURL.String(), supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, fmt.Sprintf("%s%s", Separator, ConnectorsPath))

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var connectors any
		if err := json.Unmarshal(bodyBytes, &connectors); err != nil {
			return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
		}
		return connectors, nil
	}

	return nil, c.handleResponseErrors(resp, ErrCannotGetConnector, bodyBytes)
}

func (c *centralAPIClientImpl) GetConnector(id, supportedPlatform, ballerinaVersion string) (map[string]any, error) {
	urlStr := fmt.Sprintf("%s%s%s", c.baseURL, ConnectorPathPrefix, id)

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, fmt.Sprintf("%s%s", ConnectorPathPrefix, id))

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)

	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var connector map[string]any
		if err := json.Unmarshal(bodyBytes, &connector); err != nil {
			return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
		}
		return connector, nil
	}

	return nil, c.handleResponseErrors(resp, fmt.Sprintf("%sid: %s", ErrCannotGetConnector, id), bodyBytes)
}

func (c *centralAPIClientImpl) GetConnectorByInfo(connector *models.ConnectorInfo, supportedPlatform, ballerinaVersion string) (map[string]any, error) {
	resourcePath := ConnectorPathPrefix + connector.OrgName + Separator + connector.PackageName + Separator + connector.Version + Separator + connector.ModuleName + Separator + connector.Name
	urlStr := fmt.Sprintf("%s%s", c.baseURL, resourcePath)

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, resourcePath)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var connector map[string]any
		if err := json.Unmarshal(bodyBytes, &connector); err != nil {
			return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
		}
		return connector, nil
	}

	return nil, c.handleResponseErrors(resp, fmt.Sprintf("%s'%s'", ErrCannotGetConnector, connector.PackageName), bodyBytes)
}

func (c *centralAPIClientImpl) GetTriggers(params map[string]string, supportedPlatform, ballerinaVersion string) (any, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetConnector, err.Error()))
	}

	baseURL.Path = filepath.Join(baseURL.Path, ConnectorsPath)
	query := baseURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	baseURL.RawQuery = query.Encode()

	req, err := c.newRequest(http.MethodGet, baseURL.String(), supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetTriggers, err.Error()))
	}

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetTriggers, err.Error()))
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, fmt.Sprintf("%s%s", Separator, ConnectorsPath))

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetTriggers, err.Error()))
	}

	c.logResponseVerbose(resp, string(bodyBytes))

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var connectors any
		if err := json.Unmarshal(bodyBytes, &connectors); err != nil {
			return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetTriggers, err.Error()))
		}
		return connectors, nil
	}

	return nil, c.handleResponseErrors(resp, ErrCannotGetTriggers, bodyBytes)
}

func (c *centralAPIClientImpl) GetTrigger(id, supportedPlatform, ballerinaVersion string) (map[string]any, error) {
	urlStr := fmt.Sprintf("%s%s%s", c.baseURL, TriggerPathPrefix, id)
	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetTrigger, err.Error()))
	}
	c.logRequestInitVerbose(req)
	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetTrigger, err.Error()))
	}
	defer resp.Body.Close()
	c.logRequestConnectVerbose(req, fmt.Sprintf("%s%s", TriggerPathPrefix, id))
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetTrigger, err.Error()))
	}
	c.logResponseVerbose(resp, string(bodyBytes))
	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) && resp.StatusCode == http.StatusOK {
		var trigger map[string]any
		if err := json.Unmarshal(bodyBytes, &trigger); err != nil {
			return nil, NewCentralClientError(fmt.Sprintf("%s%s", ErrCannotGetTrigger, err.Error()))
		}
		return trigger, nil
	}
	return nil, c.handleResponseErrors(resp, fmt.Sprintf("%s id: %s", ErrCannotGetTrigger, id), bodyBytes)
}

func (c *centralAPIClientImpl) AccessToken() string {
	return c.accessToken
}

func (c *centralAPIClientImpl) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *centralAPIClientImpl) pullPackageInternal(org, name, version, packagePathInBalaCache, supportedPlatform, ballerinaVersion string, isBuild bool) error {
	resourceURL := fmt.Sprintf("%s%s%s%s", PackagePathPrefix, org, Separator, name)
	enableOutputStream := os.Getenv(EnableOutputStream) == "true"
	packageSignature := fmt.Sprintf("%s%s%s", org, Separator, name)

	urlStr := fmt.Sprintf("%s%s", c.baseURL, resourceURL)

	if version != "" {
		urlStr = fmt.Sprintf("%s/%s", urlStr, version)
		packageSignature = fmt.Sprintf("%s:%s", packageSignature, version)
	} else {
		urlStr = fmt.Sprintf("%s/%s", urlStr, "*")
		packageSignature = fmt.Sprintf("%s:%s", packageSignature, "*")
	}

	logFormatter := NewLogFormatter(isBuild)

	req, err := c.newRequest(http.MethodGet, urlStr, supportedPlatform, ballerinaVersion, nil)
	if err != nil {
		return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("%s'%s'", ErrCannotPullPackage, packageSignature)))
	}

	req.Header.Set(AcceptEncoding, Identity)
	req.Header.Set(Accept, ApplicationOctetStream)

	c.logRequestInitVerbose(req)

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("%s'%s'", ErrCannotPullPackage, packageSignature)))
	}
	defer resp.Body.Close()

	c.logRequestConnectVerbose(req, resourceURL)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("%s'%s'", ErrCannotPullPackage, packageSignature)))
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
				return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("%s'%s'", ErrCannotPullPackage, packageSignature)))
			}

			downloadReq.Header.Set(AcceptEncoding, Identity)
			downloadReq.Header.Set(ContentDisposition, balaFileName)

			c.logRequestInitVerbose(downloadReq)

			downloadResp, err := client.Do(downloadReq)
			if err != nil {
				return NewCentralClientError(logFormatter.formatLog(fmt.Sprintf("%s'%s'", ErrCannotPullPackage, packageSignature)))
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

			errorMsg := logFormatter.formatLog(fmt.Sprintf("%s'%s'. BALA content download from '%s' failed.", ErrCannotPullPackage, packageSignature, balaURL))
			return c.handleResponseErrors(downloadResp, errorMsg, bodyBytes)
		}

		errorMsg := logFormatter.formatLog(fmt.Sprintf("%s'%s'. BALA content download from '%s' failed.", ErrCannotPullPackage, packageSignature, balaURL))
		return NewCentralClientError(errorMsg)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return c.handleUnauthorizedResponseWithOrg(org, bodyBytes)
	}

	contentType := resp.Header.Get(ContentType)
	if isApplicationJSONContentType(contentType) {
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
			var errResp models.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return NewCentralClientError(fmt.Sprintf("error: %s", errResp.Message))
			}
		}

		if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusServiceUnavailable {
			var errResp models.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				errorMsg := logFormatter.formatLog(fmt.Sprintf("%s'%s' from the remote repository '%s'. reason: %s",
					ErrCannotPullPackage, packageSignature, urlStr, errResp.Message))
				return NewCentralClientError(errorMsg)
			}
		}
	}

	errorMsg := logFormatter.formatLog(fmt.Sprintf("%s'%s' from the remote repository '%s'.",
		ErrCannotPullPackage, packageSignature, urlStr))
	return NewCentralClientError(errorMsg)
}

func (c *centralAPIClientImpl) getHTTPClient() *http.Client {
	if c.httpClient != nil {
		return c.httpClient
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	if c.proxyURL != nil {
		proxyURL := c.proxyURL

		if c.proxyUsername != "" && c.proxyPassword != "" {
			proxyURL = &url.URL{
				Scheme: c.proxyURL.Scheme,
				Host:   c.proxyURL.Host,
				Path:   c.proxyURL.Path,
				User:   url.UserPassword(c.proxyUsername, c.proxyPassword),
			}
		}

		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Transport: &customRetryTransport{
			transport:      transport,
			maxRetries:     c.maxRetries,
			verboseEnabled: c.verboseEnabled,
			outStream:      c.outStream,
			baseURL:        c.baseURL,
		},
		Timeout: c.callTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	c.httpClient = client
	return client
}

type customRetryTransport struct {
	transport      http.RoundTripper
	maxRetries     int
	verboseEnabled bool
	outStream      io.Writer
	baseURL        string
}

func (r *customRetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	retryCount := 0

	for retryCount <= r.maxRetries {
		// Clone request body for potential retries
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, err = io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		resp, err = r.transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode < 500 || retryCount == r.maxRetries {
			return resp, nil
		}

		var bodyContent string
		if resp.Body != nil {
			bodyBytes, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr == nil {
				bodyContent = string(bodyBytes)
			}
		}

		r.logRetryVerbose(resp, bodyContent, req, retryCount+1)

		retryCount = retryCount + 1

		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	return resp, err
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

func (c *centralAPIClientImpl) handleResponseErrors(resp *http.Response, msg string, bodyBytes []byte) error {
	contentType := resp.Header.Get(ContentType)

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		if isApplicationJSONContentType(contentType) {
			var errResp models.Error
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
			var errResp models.Error
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Message != "" {
				return NewCentralClientError(fmt.Sprintf("%s reason: %s", msg, errResp.Message))
			}
		}
	}

	return NewCentralClientError(msg)
}

func (c *centralAPIClientImpl) handleUnauthorizedResponse(bodyBytes []byte) error {
	var errResp models.Error
	if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
		return NewCentralClientError(fmt.Sprintf("unauthorized access token. check access token set in 'Settings.toml' file. reason: %s", errResp.Message))
	}
	return NewCentralClientError("unauthorized access token. check access token set in 'Settings.toml' file.")
}

func (c *centralAPIClientImpl) handleUnauthorizedResponseWithOrg(org string, bodyBytes []byte) error {
	var errResp models.Error
	if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
		return NewCentralClientError(fmt.Sprintf("unauthorized access token for organization: '%s'. check access token set in 'Settings.toml' file. reason: %s", org, errResp.Message))
	}
	return NewCentralClientError(fmt.Sprintf("unauthorized access token for organization: '%s'. check access token set in 'Settings.toml' file.", org))
}

func (r *customRetryTransport) logRetryVerbose(resp *http.Response, bodyContent string, req *http.Request, retryCount int) {
	if !r.verboseEnabled || r.outStream == nil {
		return
	}

	fmt.Fprintf(r.outStream, "< HTTP %d %s\n", resp.StatusCode, resp.Status)

	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Fprintf(r.outStream, "> %s: %s\n", name, value)
		}
	}

	fmt.Fprintln(r.outStream, "< ")

	if bodyContent != "" {
		fmt.Fprintln(r.outStream, bodyContent)
	}

	fmt.Fprintf(r.outStream, "* Connection to host %s left intact \n\n", r.baseURL)
	fmt.Fprintf(r.outStream, "* Retrying request to %s due to %d response code. Retry attempt: %d\n",
		req.URL.String(), resp.StatusCode, retryCount)
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
	return fmt.Sprintf("Bearer %s", accessToken)
}

func isApplicationJSONContentType(contentType string) bool {
	return strings.HasPrefix(contentType, MediaTypeJSONContent)
}

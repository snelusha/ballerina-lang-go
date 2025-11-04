// Package centralclient holds constants for the CentralAPIClient and related logic.
package centralclient

// HTTP Headers
const (
	BallerinaPlatform                 = "Ballerina-Platform"
	Identity                          = "identity"
	ResolvedRequestedURI              = "RESOLVED_REQUESTED_URI"
	SSL                               = "SSL"
	Authorization                     = "Authorization"
	ContentType                       = "Content-Type"
	AcceptEncoding                    = "Accept-Encoding"
	UserAgent                         = "User-Agent"
	Location                          = "Location"
	Accept                            = "Accept"
	ContentDisposition                = "Content-Disposition"
	ApplicationOctetStream            = "application/octet-stream"
	ApplicationJSON                   = "application/json"
	BallerinaCentralTelemetryDisabled = "Ballerina-Central-Telemetry-Disabled"
)

// JSON field names
const (
	Organization     = "organization"
	Version          = "version"
	BalaURL          = "balaURL"
	Platform         = "platform"
	AnyPlatform      = "any"
	PkgName          = "name"
	IsDeprecated     = "isdeprecated"
	Digest           = "digest"
	DeprecateMessage = "deprecatemessage"
)

// Environment variables
const (
	EnableOutputStream           = "ENABLE_OUTPUT_STREAM"
	SysPropCentralVerboseEnabled = "CENTRAL_VERBOSE_ENABLED"
	TestModeActive               = "TEST_MODE_ACTIVE"
	BallerinaStageCentral        = "BALLERINA_STAGE_CENTRAL"
	BallerinaDevCentral          = "BALLERINA_DEV_CENTRAL"
)

// Repository URLs
const (
	ProductionRepo = "central.ballerina.io"
	StagingRepo    = "staging-central.ballerina.io"
	DevRepo        = "dev-central.ballerina.io"
)

// API paths and query parameters
const (
	PackagesPath        = "/packages/"
	ToolsPath           = "/tools/"
	ConnectorsPath      = "/connectors/"
	TriggersPath        = "/triggers/"
	SearchQueryParam    = "?q="
	Separator           = "/"
	ResolveDependencies = "resolve-dependencies"
	ResolveModules      = "resolve-modules"
	Deprecate           = "deprecate"
	Undeprecate         = "undeprecate"
	PackagePathPrefix   = Separator + "packages" + Separator
	ToolPathPrefix      = Separator + "tools" + Separator
	ConnectorPathPrefix = Separator + "connectors" + Separator
	TriggerPathPrefix   = Separator + "triggers" + Separator
)

// Hash and crypto
const (
	SHA256          = "sha-256="
	SHA256Algorithm = "SHA-256"
)

// Progress and download constants
const (
	BytesForKB               = 1024
	ProgressBarByteThreshold = 5
	UpdateIntervalMillis     = 1000
)

// Error messages
const (
	ErrCannotFindPackage  = "error: could not connect to remote repository to find package: "
	ErrCannotFindVersions = "error: could not connect to remote repository to find versions for: "
	ErrCannotPush         = "error: failed to push the package: "
	ErrCannotPullPackage  = "error: failed to pull the package: "
	ErrCannotSearch       = "error: failed to search packages: "
	ErrCannotGetConnector = "error: failed to find connector: "
	ErrCannotGetTriggers  = "error: failed to find triggers: "
	ErrCannotGetTrigger   = "error: failed to find the trigger: "
	ErrPackageDeprecate   = "error: failed to deprecate the package: "
	ErrPackageUndeprecate = "error: failed to undo deprecation of the package: "
	ErrPackageResolution  = "error: while connecting to central: "
)

// Default timeouts and retry
const (
	DefaultConnectTimeout = 60
	DefaultReadTimeout    = 60
	DefaultWriteTimeout   = 60
	DefaultCallTimeout    = 0
	MaxRetry              = 1
	ConnectionReset       = "Connection reset"
)

// Media types
const (
	MediaTypeJSON        = "application/json; charset=utf-8"
	MediaTypeJSONContent = "application/json"
)

// Package centralclient holds constants for the CentralAPIClient and related logic.
package centralclient

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

// System property names
const (
	SysPropCentralVerboseEnabled = "CENTRAL_VERBOSE_ENABLED"
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

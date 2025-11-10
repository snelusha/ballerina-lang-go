package centralclient

type CentralClientError struct {
	message string
}

func (e *CentralClientError) Error() string {
	return e.message
}

func NewCentralClientError(message string) *CentralClientError {
	return &CentralClientError{message: message}
}

type ConnectionError struct {
	CentralClientError
}

func NewConnectionError(message string) *ConnectionError {
	return &ConnectionError{
		CentralClientError: CentralClientError{message: message},
	}
}

type NoPackageError struct {
	CentralClientError
}

func NewNoPackageError(message string) *NoPackageError {
	return &NoPackageError{
		CentralClientError: CentralClientError{message: message},
	}
}

type PackageAlreadyExistsError struct {
	CentralClientError
	version string
}

func (e *PackageAlreadyExistsError) Version() string {
	return e.version
}

func NewPackageAlreadyExistsError(message string, version string) *PackageAlreadyExistsError {
	return &PackageAlreadyExistsError{
		CentralClientError: CentralClientError{message: message},
		version:            version,
	}
}

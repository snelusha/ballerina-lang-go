package centralclient

import "fmt"

type CentralClientError struct {
	message string
	cause   error
}

func (e *CentralClientError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *CentralClientError) Unwrap() error {
	return e.cause
}

func NewCentralClientError(message string) *CentralClientError {
	return &CentralClientError{message: message}
}

func NewCentralClientErrorWithCause(message string, cause error) *CentralClientError {
	return &CentralClientError{message: message, cause: cause}
}

type ConnectionError struct {
	CentralClientError
}

func NewConnectionError(message string) *ConnectionError {
	return &ConnectionError{
		CentralClientError: CentralClientError{message: message},
	}
}

func NewConnectionErrorWithCause(message string, cause error) *ConnectionError {
	return &ConnectionError{
		CentralClientError: CentralClientError{message: message, cause: cause},
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

package errors

type IndexOutOfBoundsError struct {
	message string
}

func (e IndexOutOfBoundsError) Error() string {
	return e.message
}

func NewIndexOutOfBoundsError(message string) *IndexOutOfBoundsError {
	return &IndexOutOfBoundsError{
		message: message,
	}
}

type IllegalArgumentError struct {
	message string
}

func (e IllegalArgumentError) Error() string {
	return e.message
}

func NewIllegalArgumentError(message string) *IllegalArgumentError {
	return &IllegalArgumentError{
		message: message,
	}
}

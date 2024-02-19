package apperror

type ResourceStateError struct {
	message string
}

func (e *ResourceStateError) Error() string {
	return e.message
}

func NewResourceStateError(message string) error {
	return NewClientError(&ResourceStateError{
		message: message,
	})
}

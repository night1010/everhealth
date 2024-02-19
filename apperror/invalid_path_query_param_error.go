package apperror

type InvalidPathQueryParamError struct {
	err error
}

func (e *InvalidPathQueryParamError) Error() string {
	return e.err.Error()
}

func NewInvalidPathQueryParamError(err error) error {
	return NewClientError(&InvalidPathQueryParamError{err: err})
}

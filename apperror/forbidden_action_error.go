package apperror

type ForbiddenActionError struct {
	message string
}

func (e *ForbiddenActionError) Error() string {
	return e.message
}

func NewForbiddenActionError(message string) error {
	return NewClientError(&ForbiddenActionError{message: message}).Forbidden()
}

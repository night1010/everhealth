package apperror

import "net/http"

type ClientError struct {
	err      error
	httpCode int
}

func (e *ClientError) Error() string {
	return e.err.Error()
}

func (e *ClientError) UnWrap() error {
	return e.err
}

func (e *ClientError) HttpStatusCode() int {
	return e.httpCode
}

func (e *ClientError) BadRequest() error {
	e.httpCode = http.StatusBadRequest
	return e
}
func (e *ClientError) Unauthorized() error {
	e.httpCode = http.StatusUnauthorized
	return e
}
func (e *ClientError) Forbidden() error {
	e.httpCode = http.StatusForbidden
	return e
}
func (e *ClientError) NotFound() error {
	e.httpCode = http.StatusNotFound
	return e
}
func (e *ClientError) Conflict() error {
	e.httpCode = http.StatusConflict
	return e
}
func NewClientError(err error) *ClientError {
	httpCode := http.StatusBadRequest

	return &ClientError{
		err:      err,
		httpCode: httpCode,
	}
}

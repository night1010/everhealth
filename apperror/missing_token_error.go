package apperror

import "errors"

var ErrMissingToken = errors.New("missing authorization token")

func NewMissingTokenError() error {
	return NewClientError(ErrMissingToken).Unauthorized()
}

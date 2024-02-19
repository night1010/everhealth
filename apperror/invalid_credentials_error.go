package apperror

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewInvalidCredentialsError() error {
	return NewClientError(ErrInvalidCredentials).Unauthorized()
}

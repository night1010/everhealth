package apperror

import "errors"

var ErrInvalidToken = errors.New("invalid token")

func NewInvalidTokenError() error {
	return NewClientError(ErrInvalidToken).Unauthorized()
}

package apperror

import "fmt"

type ResourceAlreadyExistError struct {
	resource string
	field    string
	value    any
}

func (e *ResourceAlreadyExistError) Error() string {
	return fmt.Sprintf("%s with %s '%s' already exist", e.resource, e.field, e.value)
}

func NewResourceAlreadyExistError(resource string, field string, value any) error {
	return NewClientError(&ResourceAlreadyExistError{
		resource: resource,
		field:    field,
		value:    value,
	})
}

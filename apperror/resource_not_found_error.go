package apperror

import "fmt"

type ResourceNotFoundError struct {
	resource string
	field    string
	value    any
}

func (e *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s with %s: %v not found", e.resource, e.field, e.value)
}

func NewResourceNotFoundError(resource string, field string, value any) error {
	return NewClientError(&ResourceNotFoundError{
		resource: resource,
		field:    field,
		value:    value,
	}).NotFound()
}

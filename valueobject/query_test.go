package valueobject_test

import (
	"testing"

	"github.com/night1010/everhealth/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestQuery_Condition(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		operation valueobject.Operator
		value     any
		expected  any
	}{
		{
			name:      "success",
			field:     "x",
			operation: valueobject.Equal,
			value:     "y",
			expected:  "y",
		},
		{
			name:      "error",
			field:     "x",
			operation: valueobject.Equal,
			value:     "",
			expected:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := valueobject.NewQuery().Condition(tt.field, tt.operation, tt.value)

			value := q.GetConditionValue(tt.field)

			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestQuery_GetConditionValue(t *testing.T) {
	tests := []struct {
		name  string
		field string
		value any
	}{
		{"success", "title", "A"},
		{"not found", "author", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := valueobject.NewQuery().Condition("title", valueobject.Equal, "A")

			value := q.GetConditionValue(tt.field)

			assert.Equal(t, tt.value, value)
		})
	}
}

func TestQuery_GetConditions(t *testing.T) {
	var expected = []*valueobject.Condition{
		{
			Field:     "a",
			Operation: valueobject.LessThan,
			Value:     "c",
		},
	}

	q := valueobject.NewQuery().Condition("a", valueobject.LessThan, "c")

	assert.Equal(t, expected, q.GetConditions())
}

func TestQuery_WithPage(t *testing.T) {
	page := 2
	q := valueobject.NewQuery().WithPage(page)

	assert.Equal(t, page, q.GetPage())
}

func TestQuery_WithLimit(t *testing.T) {
	limit := 2
	q := valueobject.NewQuery().WithLimit(limit)

	assert.Equal(t, limit, *q.GetLimit())
}

func TestQuery_Order(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tests := []struct {
			name     string
			sortBy   string
			order    valueobject.Order
			expected string
		}{
			{name: "a asc", sortBy: "a", order: valueobject.OrderAsc, expected: "a asc"},
			{name: "b desc", sortBy: "b", order: valueobject.OrderDesc, expected: "b desc"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				q := valueobject.NewQuery().WithSortBy(tt.sortBy).WithOrder(tt.order)

				assert.Equal(t, tt.expected, q.GetOrder())
			})
		}
	})

	t.Run("empty", func(t *testing.T) {
		q := valueobject.NewQuery()
		expected := ""

		assert.Equal(t, expected, q.GetOrder())
	})
}

func TestQuery_WithJoin(t *testing.T) {
	entity := "x"
	var expected = []*valueobject.Association{
		{
			Type:   valueobject.AssociationTypeJoin,
			Entity: "x",
		},
	}

	q := valueobject.NewQuery().WithJoin(entity)

	assert.Equal(t, expected, q.GetAssociations())
}

func TestQuery_WithPreload(t *testing.T) {
	entity := "x"
	var expected = []*valueobject.Association{
		{
			Type:   valueobject.AssociationTypePreload,
			Entity: "x",
		},
	}

	q := valueobject.NewQuery().WithPreload(entity)

	assert.Equal(t, expected, q.GetAssociations())
}

func TestQuery_Lock(t *testing.T) {
	q := valueobject.NewQuery().Lock()

	assert.True(t, q.IsLocked())
}

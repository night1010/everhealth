package valueobject_test

import (
	"testing"

	"github.com/night1010/everhealth/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewCondition(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		operation valueobject.Operator
		value     any
		want      *valueobject.Condition
	}{
		{
			name:      "ILIKE",
			field:     "x",
			operation: valueobject.ILike,
			value:     "y",
			want:      &valueobject.Condition{Field: "x", Operation: valueobject.ILike, Value: "%y%"},
		},
		{
			name:      "LIKE",
			field:     "x",
			operation: valueobject.Like,
			value:     "y",
			want:      &valueobject.Condition{Field: "x", Operation: valueobject.Like, Value: "%y%"},
		},
		{
			name:      "NOT ILIKE",
			field:     "x",
			operation: valueobject.NotILike,
			value:     "y",
			want:      &valueobject.Condition{Field: "x", Operation: valueobject.NotILike, Value: "%y%"},
		},
		{
			name:      "NOT LIKE",
			field:     "x",
			operation: valueobject.NotLike,
			value:     "y",
			want:      &valueobject.Condition{Field: "x", Operation: valueobject.NotLike, Value: "%y%"},
		},
		{
			name:      "empty value",
			field:     "x",
			operation: valueobject.Equal,
			value:     "",
			want:      nil,
		},
		{
			name:      "equal",
			field:     "x",
			operation: valueobject.Equal,
			value:     "y",
			want:      &valueobject.Condition{Field: "x", Operation: valueobject.Equal, Value: "y"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valueobject.NewCondition(tt.field, tt.operation, tt.value)

			assert.Equal(t, tt.want, c)
		})
	}
}

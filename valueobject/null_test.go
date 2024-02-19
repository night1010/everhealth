package valueobject_test

import (
	"testing"
	"time"

	"github.com/night1010/everhealth/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewNullString(t *testing.T) {
	input := "A"
	expected := valueobject.NullString{NullString: struct {
		String string
		Valid  bool
	}{String: "A", Valid: true}}
	
	nullable := valueobject.NewNullString(input)

	assert.Equal(t, expected, nullable)
}

func TestNullString_MarshalJSON(t *testing.T) {
	tests := []struct {
		text     string
		valid    bool
		expected string
	}{
		{text: "A", valid: true, expected: "\"A\""},
		{text: "", valid: false, expected: "null"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			nullable := valueobject.NullString{NullString: struct {
				String string
				Valid  bool
			}{String: tt.text, Valid: tt.valid}}

			j, _ := nullable.MarshalJSON()

			assert.Equal(t, tt.expected, string(j))
		})
	}
}

func TestNewNullInt32(t *testing.T) {
	input := 10
	expected := valueobject.NullInt32{NullInt32: struct {
		Int32 int32
		Valid bool
	}{Int32: 10, Valid: true}}

	nullable := valueobject.NewNullInt32(input)

	assert.Equal(t, expected, nullable)
}

func TestNullInt32_MarshalJSON(t *testing.T) {
	tests := []struct {
		text     int32
		valid    bool
		expected any
	}{
		{text: 1, valid: true, expected: "1"},
		{text: 0, valid: false, expected: "null"},
	}

	for _, tt := range tests {
		t.Run(string(tt.text), func(t *testing.T) {
			nullable := valueobject.NullInt32{NullInt32: struct {
				Int32 int32
				Valid bool
			}{Int32: tt.text, Valid: tt.valid}}

			j, _ := nullable.MarshalJSON()

			assert.Equal(t, tt.expected, string(j))
		})
	}
}

func TestNewNullTime(t *testing.T) {
	input := time.Now()
	expected := valueobject.NullTime{NullTime: struct {
		Time  time.Time
		Valid bool
	}{Time: input, Valid: true}}

	nullable := valueobject.NewNullTime(input)

	assert.Equal(t, expected, nullable)
}

func TestNullTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		time     time.Time
		valid    bool
		expected string
	}{
		{
			time:     time.Date(2000, 01, 01, 0, 0, 0, 0, time.Local),
			valid:    true,
			expected: "\"2000-01-01\"",
		},
		{
			time:     time.Time{},
			valid:    false,
			expected: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.time.String(), func(t *testing.T) {
			nullable := valueobject.NullTime{NullTime: struct {
				Time  time.Time
				Valid bool
			}{Time: tt.time, Valid: tt.valid}}

			j, _ := nullable.MarshalJSON()

			assert.Equal(t, tt.expected, string(j))
		})
	}
}

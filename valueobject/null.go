package valueobject

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullString struct {
	sql.NullString
}

func NewNullString(s string) NullString {
	nullableString := sql.NullString{
		String: s,
		Valid:  true,
	}
	return NullString{NullString: nullableString}
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return []byte("null"), nil
}

type NullInt32 struct {
	sql.NullInt32
}

func NewNullInt32(i int) NullInt32 {
	nullableInt := sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}
	return NullInt32{NullInt32: nullableInt}
}

func (n NullInt32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int32)
	}
	return []byte("null"), nil
}

type NullTime struct {
	sql.NullTime
}

func NewNullTime(t time.Time) NullTime {
	nullableTime := sql.NullTime{
		Time:  t,
		Valid: true,
	}
	return NullTime{NullTime: nullableTime}
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time.Format("2006-01-02"))
	}
	return []byte("null"), nil
}

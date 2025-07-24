package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	return time.Time(*t).MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = Time{}
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*t = Time(parsed)
	return nil
}

func (t *Time) ToDateString() string {
	if t == nil {
		return ""
	}
	return time.Time(*t).Format("2006-01-02")
}

func (t *Time) ToString() string {
	if t == nil {
		return ""
	}
	return time.Time(*t).Format("2006-01-02 15:04:05")
}

func (t Time) Value() (driver.Value, error) {
	tlt := time.Time(t)
	if tlt.IsZero() {
		return nil, nil
	}
	return tlt, nil
}

func (t *Time) Scan(v interface{}) error {
	if v == nil {
		*t = Time{} // 设置为零值
		return nil
	}
	if value, ok := v.(time.Time); ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

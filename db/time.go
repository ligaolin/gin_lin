package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return fmt.Appendf(nil, "\"%v\"", tTime.Format("2006-01-02 15:04:05")), nil
}

func (t *Time) ToString() string {
	return time.Time(*t).Format("2006-01-02 15:04:05")
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}
func (t *Time) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

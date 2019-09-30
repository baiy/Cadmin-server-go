package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time time.Time

// MarshalJSON 序列化为JSON
func (t Time) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t Time) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}

// 数据库插入
func (t Time) Value() (driver.Value, error) {
	return t.String(), nil
}

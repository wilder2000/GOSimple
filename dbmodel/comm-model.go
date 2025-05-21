package dbmodel

import (
	"database/sql/driver"
	"fmt"
	"time"
		"github.com/wilder2000/GOSimple/glog"
)

type JSONTime struct {
	time.Time
}

// MarshalJSON 实现它的json序列化方法
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	//var stamp = fmt.Sprintf("\"%s\"", time.Time(jt).Format("2006-01-02 15:04:05"))
	tt := jt.Time
	var stamp = fmt.Sprintf("\"%s\"", tt.Format("01-02 15:04"))
	return []byte(stamp), nil
}
func (j *JSONTime) Scan(value interface{}) error {
	tt, ok := value.(time.Time)
	if !ok {
		glog.Logger.InfoF("转换为time 失败")
	}
	j.Time = tt
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j JSONTime) Value() (driver.Value, error) {
	return j.Time, nil
}

package comm

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

const (
	DateFormat          = "20060102"            //必须是这个数字，不能乱写
	TimeFormat          = "2006-01-02 15:04:05" //必须是这个数字，不能乱写
	DateExp             = `^\d{4}\d{1,2}\d{1,2}`
	ZoneShangHai        = 8 * 3600
	ZoneCST             = "CST"
	YYYY_MM_DD          = "2006-01-02"
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
)

func NowTime() string {
	return time.Now().UTC().In(LZ()).Format(TimeFormat)
}
func NowSHTime() time.Time {
	return time.Now().UTC().In(LZ())
}
func PareDate(str string) (date time.Time, err error) {
	//fmt.Printf("try parse %s", str)
	date, err = time.Parse(DateFormat, str)
	date = date.In(LZ())
	return
}
func LZ() *time.Location {
	return time.FixedZone(ZoneCST, ZoneShangHai)
}
func PareTime(str string) (date time.Time, err error) {
	date, err = time.Parse(TimeFormat, str)
	//date = date.In(LZ())
	return
}
func LocalTime() time.Time {
	t := time.Now().UTC()
	date, err := time.Parse(TimeFormat, t.Format(TimeFormat))
	if err != nil {
		panic(err)
	}
	date = date.In(LZ())
	return date
}

type JSONTime sql.NullTime

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		t.Valid = false
		return nil
	}

	var now time.Time
	if len(string(data)) == len(YYYY_MM_DD)+2 {
		now, err = time.ParseInLocation(`"`+YYYY_MM_DD+`"`, string(data), time.Local)
		t.Valid = true
		t.Time = now
	} else {
		now, err = time.ParseInLocation(`"`+YYYY_MM_DD_HH_MM_SS+`"`, string(data), time.Local)
		t.Valid = true
		t.Time = now
	}
	return
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(YYYY_MM_DD_HH_MM_SS)+2)
	b = append(b, '"')
	b = t.Time.AppendFormat(b, YYYY_MM_DD_HH_MM_SS)
	b = append(b, '"')
	return b, nil
}
func (t *JSONTime) String() string {
	if !t.Valid {
		return "null"
	}
	return t.Time.Format(YYYY_MM_DD_HH_MM_SS)
}

func (t *JSONTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Scan value time.Time
func (t *JSONTime) Scan(v interface{}) error {
	return (*sql.NullTime)(t).Scan(v)
}

func NewJSONTime(t time.Time) JSONTime {
	if t.IsZero() {
		return JSONTime{Valid: false}
	}
	return JSONTime{Valid: true, Time: t}
}

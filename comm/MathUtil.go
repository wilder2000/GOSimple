package comm

import (
	"crypto/md5"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/wilder2000/GOSimple/glog"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenExcelColumn(num int) string {
	result := ""
	for num > 0 {
		num-- // 转换为从0开始的索引，因为A=1但在计算中我们从0开始
		remainder := num % 26
		result = string('A'+remainder) + result
		num /= 26
	}
	return result
}
func IsMobile(mobile string) bool {
	result, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, mobile)
	return result
}
func DivRoundUpInt32(n, a uintptr) int32 {
	return int32((n + a - 1) / a)
}
func DivRoundUpInt(n, a uintptr) int {
	return int((n + a - 1) / a)
}
func IToString(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

const (
	DateFormat   = "20060102"            //必须是这个数字，不能乱写
	TimeFormat   = "2006-01-02 15:04:05" //必须是这个数字，不能乱写
	DateExp      = `^\d{4}\d{1,2}\d{1,2}`
	ZoneShangHai = 8 * 3600
	ZoneCST      = "CST"
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
func UUID() string {
	uniqueID, err := uuid.NewRandom()
	if err != nil {
		glog.Logger.ErrorF("could not generate UUIDv4,for:$s", err.Error())
		return ""
	}
	return strings.ToUpper(uniqueID.String())
}
func LowerUUID() string {
	uniqueID, err := uuid.NewRandom()
	if err != nil {
		glog.Logger.ErrorF("could not generate UUIDv4,for:$s", err.Error())
		return ""
	}
	return strings.ToLower(uniqueID.String())
}
func MinNum(arr ...int) int {
	min := arr[0] //假设数组的第一位为最小值
	for _, val := range arr {
		if min > val {
			min = val
		}
	}
	return min
}
func MaxNum(arr ...int) int {
	max := arr[0] //假设数组的第一位为最大值
	for i := 0; i < len(arr); i++ {
		if max < arr[i] {
			max = arr[i]
		}
	}
	return max
}

//e2c569be17396eca2a2e3c11578123ed

func MD5(messages ...string) string {
	h := md5.New()
	for _, msg := range messages {
		io.WriteString(h, msg)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

const (
	YYYY_MM_DD          = "2006-01-02"
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
)

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

type RandomObject struct {
}

func (r *RandomObject) Init() {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子
}
func (r *RandomObject) Next() string {
	randomNumber := rand.Intn(999999) // 生成一个0到999999之间的随机数
	return fmt.Sprintf("%06d", randomNumber)
}

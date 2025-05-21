package http

import (
	"encoding/json"
	"errors"
	"fmt"
	gh "net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/glog"
	"github.com/wilder2000/GOSimple/http"
	"github.com/xuri/excelize/v2"
)

func DeleteObject[T any](para *DRequest, c *gin.Context) {
	glog.Logger.InfoF("Received delete:%s", c.Request.RequestURI)
	glog.Logger.InfoF("where map:%v", para.Where)
	var mt = new(T)
	db := database.DBHander.Debug()
	db, err := http.Where(db, para.Where)
	if err != nil {
		c.JSON(gh.StatusOK, http.FailedResponseCode(1, "", err.Error()))
		return
	}
	db.Delete(mt)
	if db.RowsAffected == 1 {
		c.JSON(gh.StatusOK, http.SuccessResponse("删除成功"))
	} else {
		c.JSON(gh.StatusOK, http.FailedResponseCode(1, "", db.Error))
	}

}
func UpdateObject[T any](para *URequest, c *gin.Context) {
	glog.Logger.InfoF("Received update:%s", c.Request.RequestURI)
	glog.Logger.InfoF("field map:%v", para.Fields)

	var mt = new(T)
	db, err := http.Where(database.DBHander, para.Where)
	if err != nil {
		c.JSON(gh.StatusOK, http.FailedResponseCode(1, "", db.Error))
		return
	}
	db = db.Model(mt).Updates(para.Fields)
	if db.RowsAffected == 1 {
		c.JSON(gh.StatusOK, http.SuccessResponse("更新成功"))
	} else {
		c.JSON(gh.StatusOK, http.FailedResponseCode(1, "", db.Error))
	}

}

func CreateObject[T any](para *ARequest, c *gin.Context) {
	glog.Logger.InfoF("Received create:%s", c.Request.RequestURI)
	glog.Logger.InfoF("field map:%v", para.Fields)

	var mt = new(T)
	str := para.ObjectString
	glog.Logger.InfoF("Object string=%s", str)
	err := json.Unmarshal([]byte(str), &mt)
	if err != nil {
		c.JSON(gh.StatusOK, http.FailedResponse(err.Error(), ""))
		return
	}
	db := database.DBHander.Model(mt).Create(mt)
	if db.RowsAffected == 1 {
		c.JSON(gh.StatusOK, http.SuccessResponse("创建成功"))
	} else {
		errMsg := db.Error.Error()
		glog.Logger.ErrorF("Add role failed.%s", errMsg)
		if strings.Index(errMsg, "1062") > 0 {
			c.JSON(gh.StatusOK, http.FailedResponseCode(http.DataExistFound, "重复约束触发", db.Error.Error()))
		} else {
			c.JSON(gh.StatusOK, http.FailedResponse("增加失败", db.Error.Error()))
		}
	}

}
func ConstructData(obj any, fields map[string]string) (any, error) {
	vv := reflect.ValueOf(obj).Elem()
	tmp := reflect.New(vv.Elem().Type()).Elem()
	tmp.Set(vv.Elem())
	for fname, fvalue := range fields {
		glog.Logger.InfoF("try fo cast field:%s which value is%s", fname, fvalue)
		ff := tmp.FieldByName(fname)
		switch ff.Kind() {
		case reflect.Int:
		case reflect.Int32:
		case reflect.Int64:
		case reflect.Int8:
		case reflect.Int16:
			value, err := FieldtoInt(ff.Kind(), fvalue)
			if err != nil {
				return nil, err
			}
			ff.SetInt(value)
		case reflect.Bool:
			value, err := strconv.ParseBool(fvalue)
			if err != nil {
				return nil, err
			}
			ff.SetBool(value)
		case reflect.Float32:
			value, err := strconv.ParseFloat(fvalue, 32)
			if err != nil {
				return nil, err
			}
			ff.SetFloat(value)
		case reflect.Float64:
			value, err := strconv.ParseFloat(fvalue, 64)
			if err != nil {
				return nil, err
			}
			ff.SetFloat(value)
		case reflect.Struct:
			if ff.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
				t, er := comm.PareTime(fvalue)
				if er != nil {
					glog.Logger.ErrorF("field %s value cast to time failed. value=%v", fname, fvalue)
					return nil, er
				}
				ff.Set(reflect.ValueOf(t))
				continue
			}
		case reflect.String:
			ff.SetString(fvalue)
		default:
			return nil, errors.New("not supported field type" + ff.Kind().String())
		}

	}
	return tmp, nil
}
func FieldtoInt(k reflect.Kind, fv string) (int64, error) {
	var bitSize int
	// that the result must fit into. Bit sizes 0, 8, 16, 32, and 64
	// correspond to int, int8, int16, int32, and int64.
	switch k {
	case reflect.Int:
		bitSize = 0
	case reflect.Int8:
		bitSize = 8
	case reflect.Int16:
		bitSize = 16
	case reflect.Int32:
		bitSize = 32
	case reflect.Int64:
		bitSize = 64
	default:
		return 0, errors.New("not supported kind " + k.String())
	}
	value, err := strconv.ParseInt(fv, 10, bitSize)
	return value, err
}
func Int64to32(vv int64) int32 {
	idPointer := (*int32)(unsafe.Pointer(&vv))
	return *idPointer
}
func Int64toInt(vv int64) int {
	idPointer := (*int)(unsafe.Pointer(&vv))
	return *idPointer
}

// QueryObject 通用查询实现
func QueryObject[T any](para *QRequest, c *gin.Context) {
	glog.Logger.InfoF("Received query:%s", c.Request.RequestURI)
	glog.Logger.InfoF("query map:%v", *para)

	res, err := SelectPage[T](para)

	if err != nil {
		c.JSON(gh.StatusOK, http.FailedResponse("query failed", err))
	} else {
		qres := &QResponse{}
		qres.PageSize = res.PageSize
		qres.PageIndex = res.PageIndex
		qres.TotalPages = res.TotalPages
		qres.Data = res.Rows
		if para.Attach && len(res.Rows) > 1 {

			//有导出需求,填写上文件链接
			fn, url := AttachMgr.RequestFile()
			qres.Attach = url
			glog.Logger.InfoF("attach file: %s", fn)
			createAttachFile(res.Rows, fn)

		}
		c.JSON(gh.StatusOK, http.SuccessResponse(qres))
	}
}
func createAttachFile[T any](data []T, file string) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			glog.Logger.ErrorF("%s", err.Error())
		}
	}()
	sheetName := "exported"
	sheet1, err := f.NewSheet(sheetName)
	if err != nil {
		glog.Logger.ErrorF("Create output excel failed:%s", err.Error())
		return
	}

	first := data[0]
	val0 := reflect.ValueOf(first)
	typ0 := reflect.TypeOf(first)

	for i := 0; i < val0.NumField(); i++ {
		field := typ0.Field(i)
		cKey := comm.GenExcelColumn(i+1) + "0"
		f.SetCellStr(sheetName, cKey, field.Name)
	}
	row := 1
	for _, ca := range data {
		row++
		rowStr := strconv.Itoa(row)
		val := reflect.ValueOf(ca)
		for i := 0; i < val.NumField(); i++ {
			value := val.Field(i)
			cKey := comm.GenExcelColumn(i+1) + rowStr
			f.SetCellStr(sheetName, cKey, fmt.Sprintf("%v", value))
		}

		//f.SetCellStr(sheetName, "E"+rowStr, fmt.Sprintf("%f", ca.Geo.Latitute))
		//f.SetCellStr(sheetName, "F"+rowStr, fmt.Sprintf("%f", ca.Geo.Longitude))

		f.SetActiveSheet(sheet1)
	}
	if err := f.SaveAs(file); err != nil {
		glog.Logger.ErrorF("save attach file error: %s", err.Error())
	}

}

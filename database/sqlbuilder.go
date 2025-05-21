package database

import (
	"reflect"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/glog"
)

func sql() goqu.DialectWrapper {
	return goqu.Dialect("mysql")
}

// LookWhere 截取where开始的sql
func LookWhere(sql string) string {
	glog.Logger.InfoF("try to look where %s", sql)
	i := strings.LastIndex(strings.ToUpper(sql), "WHERE")
	iL := strings.LastIndex(strings.ToUpper(sql), "LIMIT")
	iO := strings.LastIndex(strings.ToUpper(sql), "OFFSET")
	glog.Logger.InfoF("where:%d limit:%d offset:%d", i, iL, iO)
	if i <= 0 && iL <= 0 && iO <= 0 {
		return ""
	} else if i >= 0 {
		return sql[i:]
	} else if iL < 0 || iO < 0 {
		return sql[comm.MaxNum(iL, iO):]
	} else if iL > 0 && iO > 0 {
		return sql[comm.MinNum(iL, iO):]
	} else {
		glog.Logger.Error("Not implement")
		return ""
	}
}
func map2Ep(where map[string]interface{}) []goqu.Expression {
	var condition []goqu.Expression
	for k, v := range where {
		ki := reflect.TypeOf(v).Kind()
		if ki == reflect.Slice {
			va := v.([]interface{})
			if len(va) == 0 {
				continue
			}
		} else if ki == reflect.String {
			va := v.(string)
			if len(va) == 0 {
				continue
			}
			if va == "%%" {
				continue
			}
		}
		p := strings.LastIndex(k, " ")
		if p <= 0 {
			return nil
		}
		kname := k[0:p]
		if strings.HasSuffix(k, ">=") {
			condition = append(condition, goqu.I(kname).Gte(v))
		} else if strings.HasSuffix(k, ">") {
			condition = append(condition, goqu.I(kname).Gt(v))
		} else if strings.HasSuffix(k, "<=") {
			condition = append(condition, goqu.I(kname).Lte(v))
		} else if strings.HasSuffix(k, "<") {
			condition = append(condition, goqu.I(kname).Lt(v))
		} else if strings.HasSuffix(k, "=") {
			condition = append(condition, goqu.I(kname).Eq(v))
		} else if strings.HasSuffix(k, "in") {
			condition = append(condition, goqu.I(kname).In(v))
		} else if strings.HasSuffix(k, "like") {
			condition = append(condition, goqu.I(kname).ILike(v))
		}
	}
	return condition
}

func SQLExpression(where map[string]interface{}) (string, *goqu.SelectDataset) {
	str := map2Ep(where)
	sqlBuilder := sql().From("xx").Where(str...)
	sqlStr, _, err := sqlBuilder.ToSQL()
	if err != nil {
		glog.Logger.Error(err.Error())
	}
	return LookWhere(sqlStr), sqlBuilder
}
func ToCount(sql string, countSQL string) string {
	i := strings.Index(strings.ToUpper(sql), "FROM")
	if i <= 0 {
		return ""
	} else {
		if len(countSQL) > 0 {
			return "SELECT " + countSQL + " " + sql[i:]
		} else {
			return "SELECT COUNT(*) " + sql[i:]
		}

	}
}
func ToDefCount(sql string) string {
	return ToCount(sql, "")
}

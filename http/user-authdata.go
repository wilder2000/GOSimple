package http

import (
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/glog"
)

const (
	OPERS_SQL = "select distinct d.operatorid from s_users a left join s_groupuser b on a.id=b.userid  left join s_rolegroup c on b.groupid=c.groupid left join s_roleoperator d on d.roleid=c.roleid    where a.email=?;"
	URL_SQL   = "select s.operatorid,s.url from s_urlmappings s where s.operatorid in (select distinct d.operatorid from s_users a left join s_groupuser b on a.id=b.userid  left join s_rolegroup c on b.groupid=c.groupid left join s_roleoperator d on d.roleid=c.roleid    where a.id=?)"
)

//根据用户查询可以操作的Operatorid

func UserOperators(user string) map[int32]interface{} {
	operMap := make(map[int32]interface{})
	var operList []UserOperator
	err := database.DBHander.Raw(OPERS_SQL, user).Scan(&operList).Error
	if err != nil {
		glog.Logger.ErrorF("Load user operator failed. %s", err.Error())
	} else {
		for _, oper := range operList {
			operMap[oper.OperatorID] = oper
		}
	}
	return operMap
}

//根据用户查询可以访问的URL

func UserAllUrlList(user string) map[string]interface{} {
	operMap := make(map[string]interface{})
	var urlList []UserAllowAccess
	err := database.DBHander.Raw(URL_SQL, user).Scan(&urlList).Error
	if err != nil {
		glog.Logger.ErrorF("Load user operator failed. %s", err.Error())
	} else {
		glog.Logger.InfoF("found url count %d", len(urlList))
		for _, urlObj := range urlList {
			operMap[urlObj.Url] = urlObj
			glog.Logger.InfoF("add user url %s", urlObj.Url)
		}
	}
	return operMap
}

type UserOperator struct {
	//Owner     string    `gorm:"column:owner" json:"owner"`
	OperatorID int32 `gorm:"column:operatorid" json:"operator"`
}
type UserAllowAccess struct {
	Url string `gorm:"column:url" json:"url"`
}

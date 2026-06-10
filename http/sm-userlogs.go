package http

import (
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
)

func init() {
	if database.DBHander != nil {
		database.DBHander.AutoMigrate(&dbmodel.SLog{})
	}
}

func dblog(account string, ip string, status int32) {
	log := dbmodel.SLog{
		Account: account,
		IP:      ip,
		Status:  status,
	}
	if err := database.DBHander.Create(&log).Error; err != nil {
		glog.Logger.ErrorF("login log write failed: %s", err.Error())
	}
}

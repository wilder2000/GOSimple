package http

import (
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
)

func dblog(us dbmodel.SUser, ip string) {
	log := dbmodel.SLog{
		IP:      ip,
		Account: us.ID,
		//Logintime: comm.LocalTime(),
	}
	res := database.DBHander.Create(&log)
	if res.RowsAffected == 1 {
		glog.Logger.InfoF("login log write success.")
	} else {
		glog.Logger.InfoF("login log write failed.")
	}
}

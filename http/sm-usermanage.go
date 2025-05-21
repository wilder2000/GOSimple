package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
	"github.com/wilder2000/GOSimple/http"
	hp "net/http"
	"strconv"
)

type UserMgrController struct {
	AbstractController[UpdateRequest]
}

func (s UserMgrController) UrlPath() string {
	return "/um"
}
func (s UserMgrController) Execute(req *UpdateRequest, c *gin.Context) {
	glog.Logger.InfoF("Received update:%s", c.Request.RequestURI)
	switch req.Code {
	case CmdUpdateUserName:
		uname, okName := req.Fields["name"]
		uid, okID := req.Where["id"]
		if okName && okID {
			db := database.DBHander.Model(dbmodel.SUser{}).Where("id=?", uid).Update("name", uname)
			if db.RowsAffected == 1 {
				c.JSON(hp.StatusOK, SuccessResponse(""))
			} else {
				c.JSON(hp.StatusOK, FailedResponse("Failed to update", db.Error))
			}
		} else {
			c.JSON(hp.StatusOK, FailedResponse("Not found the name or id field in request map", ""))
		}
	case CmdUpdateUserSex:
		usex, okSex := req.Fields["sex"]
		uid, okUid := req.Where["id"]
		if okSex && okUid {
			db := database.DBHander.Model(dbmodel.SUser{}).Where("id=?", uid).Update("sex", usex)
			if db.RowsAffected == 1 {
				c.JSON(hp.StatusOK, SuccessResponse(""))
			} else {
				c.JSON(hp.StatusOK, FailedResponse("Failed to update", db.Error))
			}
		} else {
			c.JSON(hp.StatusOK, FailedResponse("Not found the name or id field in request map", ""))
		}
	default:
		c.JSON(hp.StatusOK, FailedResponse("Not defined cmd:"+strconv.Itoa(req.Code), ""))
	}
}

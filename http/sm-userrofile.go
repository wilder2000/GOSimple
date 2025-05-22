package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/dbmodel"
	hp "net/http"
	"strconv"
)

type UserProfileController struct {
	AbstractController[GetRequest]
}

func (s UserProfileController) UrlPath() string {
	return "/upro"
}

func (s UserProfileController) Execute(para *GetRequest, c *gin.Context) {
	cmd := para.Code
	switch cmd {
	case CmdGetUserBasic:
		fn := para.FilterName
		fv := para.FilterValue
		var usr dbmodel.SUser
		db := database.DBHander.Raw("select * from s_users where id=@uid", sql.Named(fn, fv)).Find(&usr)
		if db.Error != nil {
			c.JSON(hp.StatusOK, FailedResponse(db.Error.Error(), ""))
		} else {
			deps, err := FindUserDeps(fv)
			if err != nil {
				c.JSON(hp.StatusOK, FailedResponse("Load user department failed."+err.Error(), ""))
			} else {
				formatter := FormatUser(usr)
				formatter.Department = deps
				c.JSON(hp.StatusOK, SuccessResponse(formatter))
			}

		}

	default:
		c.JSON(hp.StatusOK, FailedResponse("Not implement cmd:"+strconv.Itoa(cmd), ""))
	}
}

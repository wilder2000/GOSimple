package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
	hp "net/http"
	"strings"
)

type RoleAddController struct {
	AbstractController[AddRoleRequest]
}

func (s RoleAddController) UrlPath() string {
	return "/radd"
}

func (s RoleAddController) Execute(para *AddRoleRequest, c *gin.Context) {
	sr := dbmodel.SRole{
		Name:       para.Name,
		Createtime: comm.LocalTime(),
	}
	db := database.DBHander.Create(&sr)
	if db.RowsAffected == 1 {
		c.JSON(hp.StatusOK, SuccessResponse("Role添加成功"))
	} else {
		errMsg := db.Error.Error()
		glog.Logger.ErrorF("Add role failed.%s", errMsg)
		if strings.Index(errMsg, "1062") > 0 {
			c.JSON(hp.StatusOK, FailedResponseCode(DataExistFound, "重复的角色", db.Error.Error()))
		} else {
			c.JSON(hp.StatusOK, FailedResponse("Role增加失败", db.Error.Error()))
		}
	}
}

type AddRoleRequest struct {
	Name string `json:"name" binding:"required,min=2,max=12"`
}

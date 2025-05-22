package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/dbmodel"
)

var (
	TargetQueryFunc  = make(map[string]HandleQueryTarget)
	TargetCreateFunc = make(map[string]HandleCreateTarget)
	TargetDeleteFunc = make(map[string]HandleDeleteTarget)
	TargetUpdateFunc = make(map[string]HandleUpdateTarget)

	TargetPreUpdateFunc = make(map[string]func())
	AttachMgr           = AttachManager{}
)

func init() {
	RegObject[dbmodel.SUser]("user")
	RegObject[dbmodel.SRole]("role")
	RegObject[dbmodel.SGroup]("group")
	RegObject[dbmodel.SGroupuser]("groupuser")
	RegObject[dbmodel.SRolegroup]("rolegroup")
	RegObject[dbmodel.SOperator]("operator")
	RegObject[dbmodel.SRoleoperator]("roleoper")
	RegObject[dbmodel.SDepartment]("depart")
	RegObject[dbmodel.SDepuser]("depuser")

	AttachMgr.InitHome()
}
func RegObject[T any](t string) {
	TargetQueryFunc[t] = func(para *QRequest, c *gin.Context) { QueryObject[T](para, c) }
	TargetCreateFunc[t] = func(para *ARequest, c *gin.Context) { CreateObject[T](para, c) }
	TargetDeleteFunc[t] = func(para *DRequest, c *gin.Context) { DeleteObject[T](para, c) }
	TargetUpdateFunc[t] = func(para *URequest, c *gin.Context) { UpdateObject[T](para, c) }
	obj := comm.CreateGenericInstance[T]()
	comm.RegisterStruct(obj)
}

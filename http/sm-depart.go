package http

import (
	"database/sql"
	"math"
	hp "net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
)

const (
	DepUserSql = "select u.id,u.email,du.departmentid,departmentid is not null as selected from s_users u left join s_depusers du  on u.id=du.userid and du.departmentid=@did where u.email like '%@uemail%' order by selected desc"
	Did        = "did"
)

type DepartmentController struct {
	AbstractController[QRequest]
}

func (s DepartmentController) UrlPath() string {
	return "/dm"
}
func (s DepartmentController) Execute(req *QRequest, c *gin.Context) {
	glog.Logger.InfoF("Received query:%s", c.Request.RequestURI)
	switch req.Code {

	case CmdQueryUnionDepartments:
		//查询部门拥有的用户，并上没有加入的用户
		rsql := strings.ReplaceAll(DepUserSql, "@uemail", comm.IToString(req.Where["name"]))
		findDepUsers[dbmodel.SUserWithDep](req, c, rsql, Did, Did)

	default:
		c.JSON(hp.StatusOK, FailedResponse("Not defined cmd:"+strconv.Itoa(req.Code), ""))
	}
}

func findDepUsers[T any](req *QRequest, c *gin.Context, rawSql string, key string, sqlName string) {
	name, ok := req.Where[key]
	if !ok {
		c.JSON(hp.StatusOK, FailedResponse("Not found '"+key+"' in where map", ""))
		return
	}
	page := RequestToPage[T](*req)
	var rowData []*T
	did := comm.IToString(name)
	glog.Logger.InfoF("did%s", did)
	db := database.DBHander.Raw(rawSql, sql.Named(sqlName, did))
	var totalRows int64
	db.Debug().Count(&totalRows)
	glog.Logger.InfoF("count totalRows=%d", totalRows)
	glog.Logger.InfoF("count error=", db.Error)

	page.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(page.PageSize)))
	page.TotalPages = totalPages

	glog.Logger.InfoF("limit:%d", page.Limit())
	glog.Logger.InfoF("Offset:%d", page.Offset())
	rawSql += " limit @limit offset @offset"
	db2 := database.DBHander.Raw(rawSql, sql.Named(sqlName, did), sql.Named("limit", page.Limit()), sql.Named("offset", page.Offset()))

	db2.Debug().Find(&rowData)
	page.Rows = rowData
	c.JSON(hp.StatusOK, SuccessResponse(PageToResponse(page)))
}

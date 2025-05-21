package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
	"github.com/wilder2000/GOSimple/http"
	hp "net/http"
)

type UserDelController struct {
	AbstractController[QRequest]
}

func (s UserDelController) UrlPath() string {
	return "/uquery"
}

func (s UserDelController) Execute(para *QRequest, c *gin.Context) {
	glog.Logger.InfoF("Received query:%s", c.Request.RequestURI)
	// con := para.Where
	glog.Logger.InfoF("query map:%v", *para)

	// pg := NewPage[dbmodel.SUser]()
	// pg.PageSize = para.PageSize
	// pg.PageIndex = para.PageIndex
	res, err := SelectPage[dbmodel.SUser](para)
	if err != nil {
		c.JSON(hp.StatusOK, FailedResponse("query failed", err))
	} else {
		qres := &QueryResponse[dbmodel.SUser]{}
		qres.PageSize = res.PageSize
		qres.PageIndex = res.PageIndex
		qres.TotalPages = res.TotalPages
		qres.Data = res.Rows
		c.JSON(hp.StatusOK, SuccessResponse(qres))
	}
}

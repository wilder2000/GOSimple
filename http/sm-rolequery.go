package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/http"
)

type RoleQueryController struct {
	AbstractController[QRequest]
}

func (s RoleQueryController) UrlPath() string {
	return "/rquery"
}

func (s RoleQueryController) Execute(para *QRequest, c *gin.Context) {
	QueryPage[dbmodel.SRole](para, c)
}

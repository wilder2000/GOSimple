package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/http"
	hp "net/http"
)

type PwdController struct {
	AbstractController[ChangePWD]
}

func (s PwdController) UrlPath() string {
	return "/pwd"
}

func (s PwdController) Execute(para *ChangePWD, c *gin.Context) {
	err := UserProxy.ChangePwd(para.Password, para.Email)
	if err != nil {
		c.JSON(hp.StatusOK, FailedResponse("update user password failed", err))
	} else {
		c.JSON(hp.StatusOK, SuccessResponse("update user password success"))
	}

}

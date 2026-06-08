package http

import (
	"github.com/gin-gonic/gin"
	hp "net/http"
)

type PwdController struct {
	AbstractController[ChangePWD]
}

func (s PwdController) OperatorId() int32 { return OPER_ID_VIEWER }
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

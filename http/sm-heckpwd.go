package http

import (
	"github.com/gin-gonic/gin"
	hp "net/http"
)

type CheckPwdController struct {
	AbstractController[CheckPWD]
}

func (s CheckPwdController) UrlPath() string {
	return "/cpwd"
}

func (s CheckPwdController) Execute(para *CheckPWD, c *gin.Context) {
	lUser := &LoginInput{
		Email:    para.Email,
		Password: para.Password,
	}
	_, err2 := UserProxy.Login(*lUser)
	if err2 != nil {
		c.JSON(hp.StatusOK, FailedResponseCode(err2.Code(), "pwd is invalid", err2.Error()))
		return
	}
	c.JSON(hp.StatusOK, SuccessResponse(para))
}

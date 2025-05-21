package http

import (
	gh "net/http"

	"github.com/gin-gonic/gin"
)

const REQCreate = "/mif/c"

func HandleCreate(c *gin.Context) {
	if ValidHttpMethods(c) {
		var paraModel ARequest
		if err := c.ShouldBind(&paraModel); err == nil {
			hf, exist := TargetCreateFunc[paraModel.Target]
			if !exist {
				c.JSON(gh.StatusOK, FailedResponse("Not found mif data map target %"+paraModel.Target, ""))
				return
			}
			hf(&paraModel, c)
			return
		} else {
			HandlErr(c, err)
		}
	}
}

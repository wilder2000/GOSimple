package http

import (
	"github.com/gin-gonic/gin"
	gh "net/http"
)

const REQQuery = "/mif/q"

func HandleQuery(c *gin.Context) {
	if ValidHttpMethods(c) {
		var paraModel QRequest
		if err := c.ShouldBind(&paraModel); err == nil {
			//next
			hf, exist := TargetQueryFunc[paraModel.Target]
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

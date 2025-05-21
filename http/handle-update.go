package http

import (
	"github.com/gin-gonic/gin"
	gh "net/http"
)

const REQUpdate = "/mif/u"

func HandleUpdate(c *gin.Context) {
	if ValidHttpMethods(c) {
		var paraModel URequest
		if err := c.ShouldBind(&paraModel); err == nil {
			hf, exist := TargetUpdateFunc[paraModel.Target]
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

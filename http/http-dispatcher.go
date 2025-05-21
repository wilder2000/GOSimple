package http

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/glog"
)

// PreProcess 前置检查和处理
func PreProcess(c *gin.Context) {
	glog.Logger.InfoF("Received Http:%s", c.Request.RequestURI)

	userID, err := ParseHttpRequest(*c)
	if err != nil {
		glog.Logger.InfoF("Can't parse token for %s", err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnAuth(err.Error()))
		return
	}
	glog.Logger.InfoF("%s access.", userID)
	//user, err := UserProxy.GetUserByID(userID)
	////检查用户有没有过期，如果过期，需要踢出用户，todo
	//if err != nil {
	//	glog.Logger.InfoF("Can't find userid=%s", err.Error())
	//	c.AbortWithStatusJSON(http.StatusForbidden, UnAuth(nil))
	//	return
	//}
	//SaveUser(user, c)

	//检查权限 /v1
	requri := c.Request.RequestURI

	acc := lookupAccessInfo(userID)

	glog.Logger.InfoF("access uri:%s", requri)
	//url as: /api/v1/pulldataversion
	requriN := requri[7:len(requri)]
	glog.Logger.InfoF("access valid uri:%s", requriN)
	if _, ok := accessMappings[userID]; !ok {
		glog.Logger.InfoF("try load user access url list. %s", userID)
		accs := UserAllUrlList(userID)
		accessMappings[userID] = accs
	}
	accessList := accessMappings[userID]
	if _, exist := accessList[requriN]; exist {
		//存在，则可以访问
		glog.Logger.InfoF("success valid url.")
		if acc.SuccessTry() {
			c.Next()
		} else {
			glog.Logger.InfoF("success valid url. but user lock please wait a moment.")
			c.AbortWithStatusJSON(http.StatusForbidden, UnAuth("please wait a moment."))
		}

	} else {
		acc.FailedTry()
		glog.Logger.InfoF("failed valid url.")
		c.AbortWithStatusJSON(http.StatusForbidden, UnAuth("not allowed access the url."))
	}

}

var (
	accessMappings = make(map[string]map[string]interface{})
	accessTimesMap sync.Map
)

func UnAuth(data interface{}) Response {
	if data == nil {
		return FailedResponseCode(RNeedLogin, "Unauthorized", "no data")
	} else {
		return FailedResponseCode(RNeedLogin, "Unauthorized", data)
	}

}

type AccessTimes struct {
	Userid      string
	LastAccess  time.Time
	LastReset   time.Time
	FobbidTime  time.Time
	FailedTimes int
	Forbid      bool
}

// 重置状态

func (r AccessTimes) reset() {
	r.LastReset = comm.LocalTime()
	r.FobbidTime = r.LastReset
	r.LastAccess = r.LastReset
	r.FailedTimes = 0
	r.Forbid = false
}
func (r AccessTimes) SuccessTry() bool {
	if !r.Forbid {
		return true
	}
	continueMin := time.Now().Sub(r.FobbidTime).Minutes()
	if r.Forbid && continueMin > config.AConfig.Security.ForbidAccessTime {
		//达到了解禁时间
		r.reset()
		return true
	}

	return false
}
func (r AccessTimes) FailedTry() {

	continueMin := time.Now().Sub(r.FobbidTime).Minutes()
	if r.Forbid && continueMin > config.AConfig.Security.ForbidAccessTime {
		//达到了解禁时间,进入另一个检测周期
		r.reset()
		r.nextFailed()
	}
	r.nextFailed()
}

// 失败次数增加

func (r AccessTimes) nextFailed() bool {
	r.LastAccess = time.Now()
	r.FailedTimes = r.FailedTimes + 1
	if r.FailedTimes == config.AConfig.Security.MaxTryTimes {
		r.FobbidTime = time.Now()
		r.Forbid = true
		return true
	}
	return false
}

func lookupAccessInfo(useid string) AccessTimes {
	acc, exist := accessTimesMap.Load(useid)
	if !exist {
		acc = AccessTimes{
			LastAccess:  comm.LocalTime(),
			FailedTimes: 0,
			Userid:      useid,
		}
		accessTimesMap.Store(useid, acc)
	}
	accObj, ok := acc.(AccessTimes)
	if ok {
		return accObj
	} else {
		glog.Logger.ErrorF("lookup Access info failed.")
		accNew := AccessTimes{
			LastAccess:  comm.LocalTime(),
			FailedTimes: 0,
			Userid:      useid,
		}
		accessTimesMap.Store(useid, accNew)
		return accNew
	}
}

package http

import (
	"database/sql"
	"fmt"
	"math/rand"
	hp "net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
	"gorm.io/gorm"
)

const (
	UserDepsSQL = "select * from s_depusers du inner  join s_departments dt on du.departmentid=dt.id and du.userid=@uid"
	mobileCode  = "vcode"
)

type userHandler struct {
	userService Service
}

func NewUserHandler(userService Service) *userHandler {
	return &userHandler{userService}
}

// 对静态资源进行检测，一般配合nginx的ngx_http_auth_request_module插件进行

func (h *userHandler) FileAccessValid(c *gin.Context) {
	glog.Logger.InfoF("Received File Access Request:%s", c.Request.RequestURI)

	userID, err := ValidTokenHttpRequest(*c)
	if err != nil {
		glog.Logger.InfoF("Can't parse token for %s", err.Error())
		c.AbortWithStatusJSON(hp.StatusUnauthorized, UnAuth(err.Error()))
		return
	}
	glog.Logger.InfoF("%s access.", userID)
	c.JSON(hp.StatusOK, SuccessResponse(""))
}

// RequestMobileCode 获取验证码
func (h *userHandler) RequestMobileCode(c *gin.Context) {
	mobile := c.PostForm("mobile")

	if comm.IsMobile(mobile) {
		//通过
		code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
		//调用发送接口
		c.Set(mobileCode, code)
		c.JSON(hp.StatusOK, SuccessResponse(code))
		return
	} else {
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "mobile code  failed", "mobile invalid"))
		return
	}
}
func (h *userHandler) RequestUser(c *gin.Context) {
	var input RequestUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		cerr, _ := Format(err)
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "Request user failed.", cerr))
		return
	}

	var user dbmodel.SUser
	user.ID = input.Uid

	result := database.DBHander.Take(&user)
	if result.RowsAffected == 1 {

		formatter := FormatUser(user)
		c.JSON(hp.StatusOK, SuccessResponse(formatter))
	} else {
		c.JSON(hp.StatusOK, FailedResponseCode(UserNotFound, "Request failed", ""))
	}
}

//没有初始数据，按token 进行验证
//如果已经注册就登录，没有就注册
//如果服务器已经有email,name,请求上没有，就返回，
//如果请求上有，服务器也有，用app更新服务器上的

func (h *userHandler) UIDLoginRegist(c *gin.Context) {
	var input RegistUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		cerr, _ := Format(err)
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "Regist Login failed.", cerr))
		return
	}
	key := config.AConfig.AppSecret
	if key.AccessKey == input.AccessKey && key.SecretKey == input.SecretKey {
		//通过
		var user dbmodel.SUser
		user.ID = input.Uid

		result := database.DBHander.Take(&user)
		if result.RowsAffected == 1 {
			//在其它终端或者之前安装并注册过
			glog.Logger.InfoF("try regist new user but found old id=%s", input.Uid)
			sname := strings.TrimSpace(input.Name)
			if sname == input.Uid || strings.HasPrefix(input.Email, input.Uid) {
				//这种情况为name=000918.f8b575fe6e3e4214866a9d2cde4b96ba.1424,email=000918.f8b575fe6e3e4214866a9d2cde4b96ba.1424@iqnc.cn
				//检查服务器有没有正确的数据，有取出返回
				glog.Logger.InfoF("Client input wrong data. Client should update by server's name and email. they are: id=%s,name=%s,email=%s", input.Uid, sname, input.Email)
			} else if len(sname) > 0 {

				if upErr := database.DBHander.Model(&user).Updates(dbmodel.SUser{Name: sname, Email: input.Email, Modtime: comm.LocalTime()}).Error; upErr != nil {
					//update error.
					glog.Logger.ErrorF("Login failed. for ", upErr.Error())
					c.JSON(hp.StatusOK, FailedResponseCode(RServerError, "Login failed", upErr.Error()))
					return
				}

			}

		} else {
			//首次注册
			user.Aliasname = randName()
			user.Name = input.Name
			user.Icon = config.IconHome + "/" + config.AConfig.DefaultAvatar
			user.Createtime = comm.LocalTime()
			user.Modtime = user.Createtime
			user.Email = input.Email
			user.State = dbmodel.UserStateAutoRegister
			txerr := database.DBHander.Transaction(func(tx *gorm.DB) error {
				if err2 := tx.Create(&user).Error; err2 != nil {
					return err2
				}
				//自动加入终端用户组
				sgu := dbmodel.SGroupuser{
					Groupid: 1,
					Userid:  user.ID,
				}
				if err3 := tx.Create(&sgu).Error; err3 != nil {
					return err3
				}
				return nil
			})
			if txerr != nil {
				c.JSON(hp.StatusOK, FailedResponseCode(RFailed, "Login failed", txerr.Error()))
				return
			}
		}

		dblog(user, c.Request.RemoteAddr)
		glog.Logger.InfoF("GenerateToken user%v", user.ID)
		_, err3 := RefreshToken(*c, user.ID)
		if err3 != nil {

			c.JSON(hp.StatusOK, FailedResponseCode(RFailed, "Login failed", nil))
			return
		}
		SaveUser(&user)
		formatter := FormatUser(user)
		c.JSON(hp.StatusOK, SuccessResponse(formatter))

	} else {
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "Account creation failed", nil))
		return
	}
}
func (h *userHandler) UIDLoginWithExist(c *gin.Context) {
	var input LoginExistInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		cerr, _ := Format(err)
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "Login failed.", cerr))
		return
	}
	key := config.AConfig.AppSecret
	if key.AccessKey == input.AccessKey && key.SecretKey == input.SecretKey {
		//通过
		var user dbmodel.SUser
		user.ID = input.Uid

		tx := database.DBHander.Take(&user)
		if tx.RowsAffected <= 0 {
			msg := fmt.Sprintf("login failed for user not found. %s", user.ID)
			glog.Logger.Error(msg)
			c.JSON(hp.StatusOK, FailedResponseCode(RFailed, msg, nil))
			return
		}

		dblog(user, c.Request.RemoteAddr)
		glog.Logger.InfoF("GenerateToken user%v", user.ID)
		_, err3 := RefreshToken(*c, user.ID)
		if err3 != nil {

			c.JSON(hp.StatusOK, FailedResponseCode(RFailed, "Login failed", nil))
			return
		}
		SaveUser(&user)
		formatter := FormatUser(user)
		c.JSON(hp.StatusOK, SuccessResponse(formatter))

	} else {
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "Account creation failed", nil))
		return
	}
}

// MobileLogin 手机号登录
func (h *userHandler) MobileLogin(c *gin.Context) {
	mobile := c.PostForm("mobile")
	returnCode := c.PostForm("vcode")
	sendCode, exist := c.Get(mobileCode)

	glog.Logger.InfoF("sendcode %s", sendCode)
	if exist && returnCode == sendCode {
		//通过
		var user dbmodel.SUser
		err := database.DBHander.Where("mobile=?", mobile).Find(&user).Error
		if err != nil {
			c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "mobile code  failed", "db error."))
		} else {
			SaveUser(&user)
			_, err3 := RefreshToken(*c, user.ID)
			if err3 != nil {
				c.JSON(hp.StatusOK, FailedResponseCode(RFailed, "Login failed", nil))
				return
			}
			c.JSON(hp.StatusOK, SuccessResponse("login success."))
		}

	} else {
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "mobile code  failed", "code invalid"))
		return
	}
}
func (h *userHandler) UpdateAliasName(c *gin.Context) {
	var input UpdateAliasInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errMsg, _ := Format(err)
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "Input parameters wrong..", errMsg))
		return
	}
	glog.Logger.InfoF("input=%v", input)
	glog.Logger.InfoF("input uid=%s,alias=%s", input.Uid, input.AliasName)
	database.DBHander.Transaction(func(tx *gorm.DB) error {

		user := dbmodel.SUser{
			ID: input.Uid,
		}
		if err4 := tx.Model(&user).Updates(dbmodel.SUser{Aliasname: input.AliasName, Modtime: comm.LocalTime()}).Error; err4 != nil {
			return err4
		}
		return nil
	})
	c.JSON(hp.StatusOK, SuccessResponse("Update success."))
}
func (h *userHandler) DeleteAccount(c *gin.Context) {

	uid, err := ParseHttpRequest(*c)
	if err != nil {
		c.JSON(hp.StatusForbidden, FailedResponse("No auth.", ""))
		return
	}
	var input DeleteAccountInput
	err2 := c.ShouldBindJSON(&input)
	if err != nil {
		err3, _ := Format(err2)
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "Uid not found.", err3))
		return
	}
	if input.Uid != uid {
		c.JSON(hp.StatusForbidden, FailedResponse("uid is wrong.", ""))
		return
	}
	database.DBHander.Transaction(func(tx *gorm.DB) error {

		if err4 := tx.Where("userid=?", input.Uid).Delete(&dbmodel.SGroupuser{}).Error; err4 != nil {
			return err4
		}
		if err5 := tx.Where("userid=?", input.Uid).Delete(&dbmodel.SDepuser{}).Error; err5 != nil {
			return err5
		}
		if err6 := tx.Where("userid=?", input.Uid).Delete(&dbmodel.SResource{}).Error; err6 != nil {
			return err6
		}
		if err7 := tx.Where("id=?", input.Uid).Delete(&dbmodel.SUser{}).Error; err7 != nil {
			return err7
		}
		return nil
	})
	c.JSON(hp.StatusOK, SuccessResponse(uid+" all data had been delete success."))
}

// UpdateMobile 更新手机号
func (h *userHandler) UpdateMobile(c *gin.Context) {
	mobile := c.PostForm("mobile")
	returnCode := c.PostForm("vcode")
	uid := c.PostForm("uid")
	sendCode, exist := c.Get(mobileCode)
	glog.Logger.InfoF("sendcode %s", sendCode)
	if exist && returnCode == sendCode {
		//通过

		dbR := database.DBHander.Model(&dbmodel.SUser{}).Where("id=?", uid).Update("mobile", mobile)
		if dbR.RowsAffected == 1 {

			var user dbmodel.SUser
			err := database.DBHander.Where("mobile=?", mobile).Find(&user).Error
			if err != nil {
				c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "mobile code  failed", "db error."))
			} else {
				SaveUser(&user)
				_, err3 := RefreshToken(*c, user.ID)
				if err3 != nil {
					c.JSON(hp.StatusOK, FailedResponseCode(RFailed, "Login failed", nil))
				}
			}

			c.JSON(hp.StatusOK, SuccessResponse("update mobile success."))
		} else {
			c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "mobile code  failed", "db error."))
		}
	} else {
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "mobile code  failed", "code invalid"))
		return
	}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		cerr, _ := Format(err)
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "Account creation failed", cerr))
		return
	}

	newUser, merr := h.userService.RegisterUser(input, randName())

	if merr != nil {
		c.JSON(hp.StatusOK, FailedResponseCode(merr.Code(), merr.Error(), nil))
		return
	}

	_, err2 := GenerateToken(newUser.ID)
	if err2 != nil {
		c.JSON(hp.StatusOK, FailedResponseCode(RFailed, "Account creation failed", nil))
		return
	}

	formatter := FormatUser(newUser)
	c.JSON(hp.StatusOK, SuccessResponse(formatter))
}

// EmailLogin 后台管理登录，使用userid,email名称暂时不改了
func (h *userHandler) EmailLogin(c *gin.Context) {
	glog.Logger.Info("login called:" + c.Request.RequestURI)
	var input LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := FormatValidation(err)
		errorMessage := gin.H{"error": errors}
		c.JSON(hp.StatusOK, FailedResponseCode(LoginParaFormat, "Login failed", errorMessage))
		return
	}
	loggeninUser, err2 := h.userService.Login(input)
	if err2 != nil {
		glog.Logger.InfoF("Login but not found user in db,user=%s,for error:%s", input.Email, err2.Error())
		errorMessage := gin.H{"errors": "认证失败"}

		c.JSON(hp.StatusOK, FailedResponseCode(err2.Code(), "Login failed", errorMessage))
		return
	}

	_, err3 := RefreshToken(*c, loggeninUser.ID)
	if err3 != nil {
		c.JSON(hp.StatusOK, FailedResponseCode(RFailed, "Login failed", nil))
		return
	}
	formatter := FormatUser(loggeninUser)
	deps, err := FindUserDeps(loggeninUser.ID)
	if err != nil {
		c.JSON(hp.StatusOK, FailedResponseCode(RFailed, "Login failed", nil))
		return
	}
	formatter.Department = deps
	SaveUser(&loggeninUser)
	c.JSON(hp.StatusOK, SuccessResponse(formatter))
}
func FindUserDeps(uid string) ([]dbmodel.SDepartment, error) {
	var deps []dbmodel.SDepartment
	db := database.DBHander.Raw(UserDepsSQL, sql.Named("uid", uid)).Find(&deps)
	if db.Error != nil {
		return nil, db.Error
	}
	return deps, nil
}
func (h *userHandler) CHPwd(c *gin.Context) {
	glog.Logger.Info("check pwd called:" + c.Request.RequestURI)
	var input LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := FormatValidation(err)
		errorMessage := gin.H{"error": errors}
		c.JSON(hp.StatusOK, FailedResponseCode(LoginParaFormat, "Login failed", errorMessage))
		return
	}
	loggeninUser, err2 := h.userService.Login(input)
	if err2 != nil {
		glog.Logger.InfoF("Login but not found user in db,user=%s,for error:%s", input.Email, err2.Error())
		errorMessage := gin.H{"errors": "认证失败"}

		c.JSON(hp.StatusOK, FailedResponseCode(err2.Code(), "Login failed", errorMessage))
		return
	}

	//token, err := auth.RefreshToken(*c, loggeninUser.ID)
	//if err != nil {
	//	c.JSON(http.StatusOK, FailedResponseCode(comm.RFailed, "Login failed", nil))
	//	return
	//}
	formatter := FormatUser(loggeninUser)

	c.JSON(hp.StatusOK, SuccessResponse(formatter))
}
func (h *userHandler) ReportErrors(c *gin.Context) {
	var errorInput ErrorsInput
	err := c.ShouldBindJSON(&errorInput)
	if err != nil {
		cerr, _ := Format(err)
		c.JSON(hp.StatusUnprocessableEntity, FailedResponseCode(CommParaFormat, "Input Parameters failed", cerr))
		return
	}
	var debugLog dbmodel.SDebug
	debugLog.Userid = errorInput.Uid
	debugLog.Envinfo = errorInput.Envinfo
	debugLog.Detail = errorInput.Detail
	if err2 := database.DBHander.Create(&debugLog).Error; err2 != nil {
		glog.Logger.InfoF("App report error. but write db failed:%s", err2.Error())
	} else {
		glog.Logger.Info("App report error write db success.")
	}
	c.JSON(hp.StatusOK, SuccessResponse(""))
}
func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		cerr, _ := Format(err)
		c.JSON(hp.StatusUnprocessableEntity, FailedResponseCode(CommParaFormat, "Email checking failed", cerr))
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}

		c.JSON(hp.StatusUnprocessableEntity, FailedResponse("Email checking failed", errorMessage))
		return
	}
	data := gin.H{
		"is_available": isEmailAvailable,
	}
	var metaMessage string

	metaMessage = "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}
	c.JSON(hp.StatusOK, FailedResponse(metaMessage, data))

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		c.JSON(hp.StatusBadRequest, FailedResponse("Failed to upload avatar image", data))
		return
	}
	userid := c.GetHeader("userid")
	if userid == "" {
		c.JSON(hp.StatusBadRequest, FailedResponse("Failed to upload avatar image", "userid not found in header."))
		return
	}
	currentUser, exist := LookupUser(userid)
	if !exist {
		c.JSON(hp.StatusBadRequest, FailedResponse("Failed to upload avatar image", "userid not login."))
		return
	}

	userID := currentUser.ID
	path := fmt.Sprintf("/%s", file.Filename)
	glog.Logger.InfoF("path=%s", path)
	imgHome := config.AvatorHome()
	oldIcon := imgHome + currentUser.Icon
	rerr := os.Remove(oldIcon)
	if rerr != nil {
		glog.Logger.InfoF("remove old icon failed. %s", oldIcon)
	}
	err = c.SaveUploadedFile(file, imgHome+path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		c.JSON(hp.StatusBadRequest, FailedResponse("Failed to upload avatar image", data))
		return
	}
	var userDb = dbmodel.SUser{}
	userDb.ID = userID
	userDb.Icon = config.IconHome + path
	if err3 := database.DBHander.Model(&userDb).Update("icon", userDb.Icon).Error; err3 != nil {
		c.JSON(hp.StatusBadRequest, FailedResponse("Failed to upload avatar image", ""))
		return
	}
	c.JSON(hp.StatusOK, SuccessResponse(FormatUser(userDb)))
}

var (
	sampleNames []string
)

func randName() string {

	if sampleNames == nil {
		var items []dbmodel.SItem
		res := database.DBHander.Where("type=1").Find(&items)
		if res.Error != nil {
			glog.Logger.Warn("load items from db failed. so use Application.yml names")
			for _, name := range config.AConfig.UserNames {
				sampleNames = append(sampleNames, name)
			}
		} else {
			for _, item := range items {
				sampleNames = append(sampleNames, item.Name)
			}
		}

	}
	n := len(sampleNames)
	i := rand.Intn(n)
	return sampleNames[i]
}

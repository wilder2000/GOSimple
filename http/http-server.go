package http

import (
	"context"
	"fmt"
	gh "net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/dbscript"
	"github.com/wilder2000/GOSimple/glog"
)

func init() {
	fmt.Printf("%s\n", config.Logo)
	fmt.Printf("%s\n", config.LogoTitle)
	RegMapping[ChangePWD](new(PwdController))
	RegMapping[CheckPWD](new(CheckPwdController))
	RegMapping[QRequest](new(UserQueryController))
	RegMapping[QRequest](new(RoleQueryController))
	RegMapping[AddRoleRequest](new(RoleAddController))
	RegMapping[QRequest](new(UserGroupsController))
	RegMapping[QRequest](new(OperatorController))
	RegMapping[QRequest](new(DepartmentController))
	RegMapping[UpdateRequest](new(UserMgrController))
	RegMapping[GetRequest](new(UserProfileController))

}

var (
	mappings       = make(map[string]gin.HandlerFunc)
	noAuthMappings = make(map[string]gin.HandlerFunc)
)

type HttpController struct {
	Path   string
	Action gin.HandlerFunc
}

func newController[M any](c HTTPController[M]) *AbstractController[M] {
	ty := reflect.ValueOf(c)
	fi := ty.Elem().FieldByName("AbstractController")
	if fi.Type().ConvertibleTo(reflect.TypeOf(AbstractController[M]{})) {
		cc := fi.Interface().(AbstractController[M])
		cc.HTTPController = c
		return &cc
	}
	return nil
}
func RegMapping[M any](c HTTPController[M]) {
	ctrl := newController(c)
	glog.Logger.InfoF("Try to regist HTTPController %s ", c.UrlPath())
	mappings[c.UrlPath()] = ctrl.Prepare
}

func initController(e *gin.Engine) (*gin.RouterGroup, []string) {
	uh := NewUserHandler(UserProxy)
	e.POST("/api/emllogin", uh.EmailLogin)
	if config.AConfig.Security.Registration {
		e.POST("/api/reguser", uh.RegisterUser)
	}

	e.POST("/api/reqmcode", uh.RequestMobileCode)
	e.POST("/api/token_valid", uh.FileAccessValid)
	//e.POST("/autoreguser", uh.AutoRegUser)
	e.POST("/api/updmobile", uh.UpdateMobile)

	e.POST("/api/moblogin", uh.MobileLogin)
	e.POST("/api/newreglogin", uh.UIDLoginRegist)
	e.POST("/api/loginexist", uh.UIDLoginWithExist)

	for path, hl := range noAuthMappings {
		e.POST(path, hl)
		glog.Logger.InfoF("NO Auth Mapping:%s Function:%s", path, reflect.TypeOf(hl).Name())
	}

	proUrlGrp := e.Group("/api/v1", PreProcess)
	urlMapping := make([]string, 0)

	registMapping := func(path string, hf gin.HandlerFunc) {
		proUrlGrp.POST(path, hf)
		urlMapping = append(urlMapping, path)
		glog.Logger.InfoF("Mapping=========:%s Function:%s", path, reflect.TypeOf(hf).Name())
	}

	for path, hl := range mappings {
		registMapping(path, hl)
	}
	registMapping("/avatorup", uh.UploadAvatar)
	registMapping("/requestuser", uh.RequestUser)
	registMapping("/delaccount", uh.DeleteAccount)
	registMapping("/modalias", uh.UpdateAliasName)
	registMapping("/reperror", uh.ReportErrors)
	registMapping(REQCreate, HandleCreate)
	registMapping(REQQuery, HandleQuery)
	registMapping(REQDelete, HandleDelete)
	registMapping(REQUpdate, HandleUpdate)
	return proUrlGrp, urlMapping
}
func registGinFunc(rg *gin.RouterGroup, path string, hf gin.HandlerFunc) {
	rg.POST(path, hf)
}

// Start http server startWebServer func

func startWebServer(address string, readout time.Duration, wout time.Duration, actions []HttpController, noauthActions []HttpController, install bool) {

	router := gin.Default()
	staticDir := config.AConfig.StaticDir
	if web := config.AConfig.Web; web != nil {
		for k, v := range web {
			fmt.Printf("web static :%s -> %s\n", k, v)
			router.Static(k, v)
		}
	}
	if len(strings.TrimSpace(staticDir.AbsoluteFileDir)) > 0 {
		router.Static(staticDir.RelativePath, staticDir.AbsoluteFileDir)
		glog.Logger.InfoF("mapping www url:%s file dir:%s", staticDir.RelativePath, staticDir.AbsoluteFileDir)
	} else {
		glog.Logger.InfoF("no mapping www path config.")
	}

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	if config.AConfig.AccessControlAllowOrigin {
		glog.Logger.InfoF("AccessControlAllowOrigin true")
		glog.Logger.InfoF("AllowHost: ", config.AConfig.AccessControlAllowHost)
		glog.Logger.InfoF("AllowMethods: ", config.AConfig.AccessControlAllowMethods)
		glog.Logger.InfoF("AllowHeaders: ", config.AConfig.AccessControlAllowHeaders)
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", config.AConfig.AccessControlAllowHost) // 允许任何源
			c.Writer.Header().Set("Access-Control-Allow-Methods", config.AConfig.AccessControlAllowMethods)
			c.Writer.Header().Set("Access-Control-Allow-Headers", config.AConfig.AccessControlAllowHeaders)
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return // 直接返回204状态码
			}
			c.Next() // 继续执行其他的中间件链
		})
	} else {
		glog.Logger.InfoF("AccessControlAllowOrigin false")
	}

	for _, mapping := range noauthActions {
		noAuthMappings[mapping.Path] = mapping.Action
	}
	rr, urlList := initController(router)
	for _, mapping := range actions {
		glog.Logger.InfoF(" POST Mapping:%s", mapping.Path)
		rr.POST(mapping.Path, mapping.Action)
	}

	if !install {

		srv := &gh.Server{
			Addr:           address,
			Handler:        router,
			ReadTimeout:    readout * time.Second,
			WriteTimeout:   wout * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && errors.Is(err, gh.ErrServerClosed) {
				glog.Logger.InfoF("http server error: %s\n", err)
			} else {
				glog.Logger.InfoF("Http Server started success. Binding :%s", address)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		glog.Logger.Info("Shutting down server...")
		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			glog.Logger.ErrorF("Server forced to shutdown:", err)
		}
	} else {
		//install
		dbscript.Install(urlList)
	}

	glog.Logger.Info("Server exiting")

}

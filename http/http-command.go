package http

import (
	"fmt"
	hp "net/http"
	"reflect"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/glog"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}
type DefaultResponse[T interface{}] struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    T      `json:"data"`
}

//	type QueryRequest struct {
//		PageIndex int            `json:"PageIndex"`
//		PageSize  int            `json:"PageSize"`
//		Code      int            `json:"Code"` //操作码
//		Where     map[string]any `json:"Where"`
//		SqlItems  []interface{}  `json:"SqlItems"`
//	}
type DeleteRequest struct {
	Code  int            `json:"Code"` //操作码
	Where map[string]any `json:"Where"`
}
type GetRequest struct {
	Code        int    `json:"Code"` //操作码
	FilterName  string `json:"FilterName" binding:"required"`
	FilterValue string `json:"FilterValue" binding:"required"`
}
type UpdateRequest struct {
	Code   int            `json:"Code"` //操作码
	Fields map[string]any `json:"Fields"`
	Where  map[string]any `json:"Where"`
}

func (r UpdateRequest) LookField(key string) (string, bool) {
	v, ok := r.Fields[key]
	if ok {
		return fmt.Sprintf("%v", v), true
	} else {
		return "", false
	}
}

func (r UpdateRequest) TryField(key string) (string, error) {
	v, ok := r.Fields[key]
	if ok {
		return fmt.Sprintf("%v", v), nil
	} else {
		return "", NewMVCError(ErrorParaNotExist, key+" not found in request.")
	}
}

type QueryResponse[T any] struct {
	PageIndex  int    `json:"PageIndex"`
	PageSize   int    `json:"PageSize"`
	TotalPages int    `json:"TotalPages"`
	TotalRows  int64  `json:"TotalRows"`
	Message    string `json:"message"`
	Code       int    `json:"code"`
	Data       []*T   `json:"Data"`
}

type WError interface {
	error
	Code() int
}

func RequestToPage[T any](req QRequest) *Page[T] {
	return &Page[T]{
		PageSize:  req.PageSize,
		PageIndex: req.PageIndex,
	}
}
func PageToResponse[T any](res *Page[T]) *QueryResponse[T] {
	return &QueryResponse[T]{
		PageSize:   res.PageSize,
		PageIndex:  res.PageIndex,
		TotalPages: res.TotalPages,
		TotalRows:  res.TotalRows,
		Data:       res.Rows,
	}
}

type AuthService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

func apiResponse(message string, code int, data interface{}) Response {

	jsonResponse := Response{
		Message: message,
		Code:    code,
		Data:    data,
	}
	return jsonResponse
}
func SuccessResponse(data interface{}) Response {
	rep := apiResponse("success", RSuccess, data)
	return rep
}
func FailedResponse(err string, data interface{}) Response {
	rep := apiResponse(err, RFailed, data)
	return rep
}
func FailedResponseCode(ec int, err string, data interface{}) Response {
	rep := apiResponse(err, ec, data)
	return rep
}

func HandlErr(c *gin.Context, err error) {
	errs, ok := Format(err)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		c.JSON(hp.StatusOK, FailedResponse("Server internal error.", err.Error()))
		return
	}
	// validator.ValidationErrors类型错误则进行翻译
	c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "valid failed", errs))
	return
}

type MVCError struct {
	errcode int
	message string
}

func NewMVCError(ec int, msg string) *MVCError {
	err := &MVCError{
		errcode: ec,
		message: msg,
	}
	return err
}
func (M MVCError) Error() string {
	return M.message
}

func (M MVCError) Code() int {
	return M.errcode
}

type HTTPController[T any] interface {
	Execute(para *T, c *gin.Context)
	Prepare(c *gin.Context)
	UrlPath() string
}
type AbstractController[T any] struct {
	HTTPController[T]
}

func (b *AbstractController[T]) Execute(para *T, c *gin.Context) {
	glog.Logger.InfoF("HTTPController %+v", b.HTTPController)
	glog.Logger.InfoF("HTTPController Execute %+v", b.HTTPController.Execute)
	b.HTTPController.Execute(para, c)
}

func (b *AbstractController[T]) Prepare(c *gin.Context) {
	var paraModel T

	if err := c.ShouldBind(&paraModel); err == nil {
		//自己实现一个统一的controler,再分发，再简化controler的实现
		if c.Request.Method == "POST" || c.Request.Method == "GET" {
			b.Execute(&paraModel, c)
			return
		} else {
			emsg := "Not implement method:" + c.Request.Method
			glog.Logger.Error(emsg)
			c.JSON(hp.StatusOK, emsg)
		}

	} else {
		errs, ok := Format(err)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(hp.StatusOK, FailedResponse("Server internal error.", err.Error()))
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(hp.StatusOK, FailedResponseCode(CommParaFormat, "valid failed", errs))
		return
	}

}

type HttpController struct {
	Path   string
	Action gin.HandlerFunc
}

func RegMapping[M any](c HTTPController[M]) {
	ctrl := newController(c)
	glog.Logger.InfoF("Try to regist HTTPController %s ", c.UrlPath())
	mappings[c.UrlPath()] = ctrl.Prepare
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

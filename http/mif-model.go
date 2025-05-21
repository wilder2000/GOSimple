package http

import (
	"math"

	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/glog"
)

type QRequest struct {
	PageIndex int      `json:"PageIndex"`
	PageSize  int      `json:"PageSize"`
	Code      int      `json:"Code"`   //操作码
	Target    string   `json:"Target"` //目标对象
	Order     string   `json:"order"`  //目标对象
	Fields    []string `json:"fields"` //指定返回字段
	Attach    bool     `json:"attach"` //是否导出文件
	// SqlItems  []interface{}  `json:"SqlItems"`
	Where map[string]any `json:"Where"`
}

type ARequest struct {
	Target       string            `json:"Target"` //目标对象名称
	Fields       map[string]string `json:"Fields"`
	ObjectString string            `json:"ObjectString"`
}

type URequest struct {
	Target string         `json:"Target"` //目标对象名称
	Fields map[string]any `json:"Fields"`
	Where  map[string]any `json:"Where"`
}

type DRequest struct {
	Target string         `json:"Target"` //目标对象名称
	Where  map[string]any `json:"Where"`
}
type QResponse struct {
	PageIndex  int    `json:"PageIndex"`
	PageSize   int    `json:"PageSize"`
	TotalPages int    `json:"TotalPages"`
	Message    string `json:"message"`
	Code       int    `json:"code"`
	Attach     string `json:"attach"`
	Data       any    `json:"Data"`
}
type HandleQueryTarget func(para *QRequest, c *gin.Context)
type HandleCreateTarget func(para *ARequest, c *gin.Context)
type HandleDeleteTarget func(para *DRequest, c *gin.Context)
type HandleUpdateTarget func(para *URequest, c *gin.Context)

type QueryTanslate[T any] interface {
	Translate(raw []T) []T
}

const (
	DefPageSize = 20
)

func NewPage[T any]() *Page[T] {
	p := &Page[T]{
		PageSize:  DefPageSize,
		PageIndex: 1,
	}
	return p
}

type Page[T any] struct {
	PageSize   int    `json:"pageSize,omitempty" form:"pageSize"`
	PageIndex  int    `json:"pageIndex,omitempty" form:"pageIndex"`
	Sort       string `json:"sort,omitempty" form:"sort"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Rows       []*T   `json:"rows"`
}

func (p *Page[T]) InitPage(totalRows int64) {
	p.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(p.PageSize)))
	p.TotalPages = totalPages

	glog.Logger.InfoF("total rows%d", p.TotalRows)
	glog.Logger.InfoF("total page%d", p.TotalPages)
}

func (p *Page[T]) Offset() int {
	return (p.CurrentPage() - 1) * p.Limit()
}

func (p *Page[T]) Limit() int {
	if p.PageSize == 0 {
		p.PageSize = DefPageSize
	}
	return p.PageSize
}

func (p *Page[T]) CurrentPage() int {
	if p.PageIndex == 0 {
		p.PageIndex = 1
	}
	return p.PageIndex
}

func (p *Page[T]) GetSort() string {
	// if p.Sort == "" {
	// 	p.Sort = "Id desc"
	// }
	return p.Sort
}

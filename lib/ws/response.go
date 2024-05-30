package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/v2/lib/bizerr"
	"github.com/hdget/hdutils/convert"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type PageResponse struct {
	Response
	Total int64 `json:"total"`
}

//const (
//	_                     = bizerr.ErrCodeModuleRoot + iota
//	ErrCodeServerInternal // 内部错误
//	ErrCodeUnauthorized   // 未授权
//	ErrCodeInvalidRequest // 非法请求
//	ErrCodeForbidden      // 拒绝访问
//
//)

// Success respond with data
// empty args respond with 'ok' message
// args[0] is the response data
func Success(c *gin.Context, args ...interface{}) {
	var ret Response
	switch len(args) {
	case 0:
		ret.Data = "ok"
	case 1:
		ret.Data = args[0]
	}

	c.PureJSON(http.StatusOK, ret)
}

func Error(c *gin.Context, code int, msg string) {
	ret := &Response{
		Code: code,
		Msg:  msg,
	}
	c.PureJSON(http.StatusOK, ret)
}

// SuccessRaw respond with raw data
func SuccessRaw(c *gin.Context, result interface{}) {
	var content string
	switch t := result.(type) {
	case string:
		content = t
	case []byte:
		content = convert.BytesToString(t)
	default:
		v, _ := json.Marshal(result)
		content = convert.BytesToString(v)
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "application/json; charset=utf-8")
	_, _ = c.Writer.WriteString(content)
}

// SuccessPages respond with pagination data
func SuccessPages(c *gin.Context, total int64, pages interface{}) {
	c.PureJSON(http.StatusOK, PageResponse{
		Response: Response{
			Data: pages,
		},
		Total: total,
	})
}

// Failure grpc http错误
func Failure(c *gin.Context, err error) {
	c.PureJSON(http.StatusOK, bizerr.Convert(err))
}

func InvalidRequest(c *gin.Context, err error) {
	c.PureJSON(http.StatusBadRequest, bizerr.Convert(err))
}

func Forbidden(c *gin.Context, err error) {
	c.PureJSON(http.StatusForbidden, bizerr.Convert(err))
}

func Unauthorized(c *gin.Context, err error) {
	c.PureJSON(http.StatusUnauthorized, bizerr.Convert(err))
}

func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
}

func PermanentRedirect(c *gin.Context, location string) {
	c.Redirect(http.StatusMovedPermanently, location)
}

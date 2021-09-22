package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/lib/err"
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

const (
	_                     = err.ErrCodeModuleRoot + iota
	ErrCodeInvalidRequest // 非法请求
	ErrCodeInternalError  // 内部错误
	ErrCodeUnauthorized   // 未授权
)

var (
	ErrInvalidRequest = err.New("invalid request", ErrCodeInvalidRequest)
	ErrForbidden      = err.New("unauthorized request", ErrCodeUnauthorized)
)

// RespondSuccess respond with data
// empty args respond with 'ok' message
// args[0] is the response data
func RespondSuccess(c *gin.Context, args ...interface{}) {
	var ret Response
	switch len(args) {
	case 0:
		ret.Data = "ok"
	case 1:
		ret.Data = args[0]
	}
	c.PureJSON(http.StatusOK, ret)
}

// RespondPages respond with pagination data
func RespondPages(c *gin.Context, total int64, pages interface{}) {
	c.PureJSON(http.StatusOK, PageResponse{
		Response: Response{
			Data: pages,
		},
		Total: total,
	})
}

func RespondError(c *gin.Context, err error) {
	c.PureJSON(http.StatusOK, err2response(err))
}

func RespondInvalidRequest(c *gin.Context) {
	c.PureJSON(http.StatusOK, err2response(ErrInvalidRequest))
}

func Forbidden(c *gin.Context) {
	c.PureJSON(http.StatusForbidden, err2response(ErrForbidden))
}

func err2response(e error) *Response {
	code := ErrCodeInternalError
	codeErr, ok := e.(*err.CodeError)
	if ok {
		code = codeErr.Code()
	}
	return &Response{
		Msg:  e.Error(),
		Code: code,
	}
}

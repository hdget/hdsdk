package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"github.com/hdget/hdsdk/lib/bizerr"
	"github.com/hdget/hdsdk/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	_                     = bizerr.ErrCodeModuleRoot + iota
	ErrCodeServerInternal // 内部错误
	ErrCodeUnauthorized   // 未授权
	ErrCodeInvalidRequest // 非法请求
	ErrCodeForbidden      // 拒绝访问

)

var (
	errInvalidRequest = bizerr.New(ErrCodeInvalidRequest, "invalid request")
	errForbidden      = bizerr.New(ErrCodeForbidden, "forbidden")
	errUnauthorized   = bizerr.New(ErrCodeUnauthorized, "unauthorized")
)

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
		content = utils.BytesToString(t)
	default:
		v, _ := json.Marshal(result)
		content = utils.BytesToString(v)
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
	c.PureJSON(http.StatusOK, fromStatusError(err))
}

func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
}

func PermanentRedirect(c *gin.Context, location string) {
	c.Redirect(http.StatusMovedPermanently, location)
}

func InvalidRequest(c *gin.Context) {
	c.PureJSON(http.StatusBadRequest, err2response(errInvalidRequest))
}

func Forbidden(c *gin.Context) {
	c.PureJSON(http.StatusForbidden, err2response(errForbidden))
}

func Unauthorized(c *gin.Context) {
	c.PureJSON(http.StatusUnauthorized, err2response(errUnauthorized))
}

// fromStatusError 从grpc status error获取额外的错误信息
func fromStatusError(err error) interface{} {
	if err == nil {
		return nil
	}

	cause := errors.Cause(err)
	st, ok := status.FromError(cause)
	if ok {
		details := st.Details()
		if len(details) > 0 {
			var pbErr bizerr.Error
			e := proto.Unmarshal(st.Proto().Details[0].GetValue(), &pbErr)
			if e == nil {
				return pbErr
			}
		}
	}

	return bizerr.Error{
		Code:    int32(codes.Internal),
		Message: err.Error(),
	}
}

func err2response(e error) *Response {
	code := ErrCodeServerInternal
	be, ok := e.(*bizerr.BizError)
	if ok {
		code = be.Code
	}
	return &Response{
		Msg:  e.Error(),
		Code: code,
	}
}

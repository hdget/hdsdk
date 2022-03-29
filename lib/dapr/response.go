package dapr

import (
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/lib/err"
)

// ApiError 默认的错误结构
type ApiError struct {
	Code int
	Msg  string
}

const DefaultErrorCode = 1

// Reply with response
func Reply(event *common.InvocationEvent, resp interface{}) *common.Content {
	// 判断响应是否是错误
	e, ok := resp.(error)
	if ok {
		var apiError interface{}
		// 如果响应是错误，检查err是否是CodeErr
		ce, ok := e.(err.CodeError)
		if ok {
			apiError = &ApiError{
				Code: ce.Code(),
				Msg:  ce.Error(),
			}
		} else {
			apiError = &ApiError{
				Code: DefaultErrorCode,
				Msg:  e.Error(),
			}
		}

		return getDaprContent(event, apiError)
	}

	// 返回正常数据
	return getDaprContent(event, resp)
}

func getDaprContent(event *common.InvocationEvent, resp interface{}) *common.Content {
	data, err := json.Marshal(resp)
	if err != nil {
		return nil
	}

	return &common.Content{
		ContentType: ContentTypeJson,
		Data:        data,
		DataTypeURL: event.DataTypeURL,
	}
}

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

var DefaultApiError = ApiError{
	Code: 1,
	Msg:  "api error",
}

// Success reply with success response
func Success(event *common.InvocationEvent, resp interface{}) (*common.Content, error) {
	return getDaprContent(event, resp), nil
}

// Error Reply with response
func Error(event *common.InvocationEvent, e error) (*common.Content, error) {
	apiError := &ApiError{
		Code: DefaultApiError.Code,
		Msg:  e.Error(),
	}

	// 如果响应是错误，检查err是否是CodeErr
	ce, ok := e.(err.CodeError)
	if ok {
		apiError.Code = ce.Code()
	}

	return getDaprContent(event, apiError), e
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

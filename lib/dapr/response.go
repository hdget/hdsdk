package dapr

import (
	"encoding/json"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/lib/bizerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Error(event *common.InvocationEvent, err error) (*common.Content, error) {
	be, ok := err.(*bizerr.BizError)
	if ok {
		st, _ := status.New(codes.Internal, "").WithDetails(&bizerr.Error{
			Code:    int32(be.Code),
			Message: be.Message,
		})
		return nil, st.Err()
	}
	return nil, err
}

func Success(event *common.InvocationEvent, resp interface{}) (*common.Content, error) {
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return &common.Content{
		ContentType: ContentTypeJson,
		Data:        data,
		DataTypeURL: event.DataTypeURL,
	}, nil
}

package dapr

import (
	"encoding/json"
	"errors"
	"github.com/dapr/go-sdk/service/common"
	"github.com/hdget/hdsdk/lib/bizerr"
	"github.com/hdget/hdsdk/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Error(err error) (*common.Content, error) {
	var be *bizerr.BizError
	ok := errors.As(err, &be)
	if ok {
		st, _ := status.New(codes.Internal, "").WithDetails(&bizerr.Error{
			Code: be.Code,
			Msg:  be.Msg,
		})
		return nil, st.Err()
	}
	return nil, err
}

func Success(event *common.InvocationEvent, resp any) (*common.Content, error) {
	var err error
	var data []byte
	switch t := resp.(type) {
	case string:
		data = utils.StringToBytes(t)
	case []byte:
		data = t
	default:
		data, err = json.Marshal(resp)
		if err != nil {
			return nil, err
		}
	}

	return &common.Content{
		ContentType: ContentTypeJson,
		Data:        data,
		DataTypeURL: event.DataTypeURL,
	}, nil
}

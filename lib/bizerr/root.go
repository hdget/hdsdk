package bizerr

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BizError struct {
	Code int
	Msg  string
}

// New an error support error code
func New(code int, message string) *BizError {
	return &BizError{
		Code: code,
		Msg:  message,
	}
}

func (be BizError) Error() string {
	return be.Msg
}

// FromStatusError 从grpc status error获取额外的错误信息
func FromStatusError(err error) interface{} {
	if err == nil {
		return nil
	}

	cause := errors.Cause(err)
	st, ok := status.FromError(cause)
	if ok {
		details := st.Details()
		if len(details) > 0 {
			var pbErr Error
			e := proto.Unmarshal(st.Proto().Details[0].GetValue(), &pbErr)
			if e == nil {
				return pbErr
			}
		}
	}

	return Error{
		Code: int32(codes.Internal),
		Msg:  err.Error(),
	}
}

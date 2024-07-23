package bizerr

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type BizError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

// New an error support error code
func New[T Integer](code T, message string) *BizError {
	return &BizError{
		Code: int32(code),
		Msg:  message,
	}
}

func InternalError(message string) *BizError {
	return &BizError{
		Code: 10001, // 10001为自定义的内部错误编码
		Msg:  message,
	}
}

func (be BizError) Error() string {
	return be.Msg
}

// Convert 从grpc status error获取额外的错误信息
func Convert(err error) *BizError {
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
				return &BizError{
					Code: pbErr.Code,
					Msg:  pbErr.Msg,
				}
			}
		}
	}

	return &BizError{
		Code: int32(codes.Internal),
		Msg:  err.Error(),
	}
}

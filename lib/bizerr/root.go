package bizerr

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

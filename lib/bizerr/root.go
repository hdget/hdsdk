package bizerr

type BizError struct {
	Code    int
	Message string
}

// New an error support error code
func New(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

func (be BizError) Error() string {
	return be.Message
}

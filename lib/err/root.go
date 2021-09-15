package err

type CodeError struct {
	code int
	msg  string
}


// ErrCodeStart 错误代码开始值
const ErrCodeStart = 10000

// ErrCodeModuleRoot define error code module, e,g: 10000, 20000, 30000...
const (
	ErrCodeModuleRoot = ErrCodeStart * (1 + iota)
)

// define common error code
const (
	ErrCodeUnknown = ErrCodeModuleRoot + iota // unknown error code
	ErrCodeInternal // internal error
)

// New an error support error code
func New(msg string, args ...int) *CodeError {
	code := ErrCodeUnknown
	if len(args) > 0 {
		code = args[0]
	}

	return &CodeError{
		code: code,
		msg:  msg,
	}
}

func (ce *CodeError) Code() int {
	return ce.code
}

func (ce *CodeError) Error() string {
	return ce.msg
}


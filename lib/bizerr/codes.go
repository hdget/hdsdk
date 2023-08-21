package bizerr

// ErrCodeStart 业务逻辑错误代码开始值
const ErrCodeStart = 10000

// ErrCodeModuleRoot define error code module, e,g: 10000, 20000, 30000...
const (
	ErrCodeModuleRoot = ErrCodeStart * (1 + iota)
)

// define common error code
const (
	_               = ErrCodeModuleRoot + iota // unknown error code
	ErrCodeInternal                            // internal error
)

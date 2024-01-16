package errdef

import "github.com/pkg/errors"

var (
	ErrInvalidLogConfig = errors.New("invalid logger configer")
)

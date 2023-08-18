package typlog

import "github.com/pkg/errors"

var (
	ErrInvalidLogConfig = errors.New("invalid log config")
	ErrNoProvider       = errors.New("log provider not found")
)

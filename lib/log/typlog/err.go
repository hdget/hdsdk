package typlog

import "github.com/pkg/errors"

var (
	ErrEmptyLogConfig   = errors.New("empty log config")
	ErrInvalidLogConfig = errors.New("invalid log config")
	ErrNoProvider       = errors.New("log provider not found")
)

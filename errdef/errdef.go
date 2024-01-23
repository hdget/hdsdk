package errdef

import "github.com/pkg/errors"

var (
	ErrInvalidConfig = errors.New("invalid config")
	ErrEmptyConfig   = errors.New("empty config")
)

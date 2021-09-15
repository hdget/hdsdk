package types

import "errors"

var (
	ErrEmptyConfig   = errors.New("empty config")
	ErrInvalidConfig = errors.New("invalid config")
)

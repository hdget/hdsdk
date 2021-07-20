package types

import "errors"

var (
	ErrMissingValue  = errors.New("(MISSING)")
	ErrEmptyConfig   = errors.New("empty config")
	ErrInvalidConfig = errors.New("invalid config")
	ErrNoProvider    = errors.New("provider not found")
	ErrNoCapability  = errors.New("capability not defined")
)

package errdef

import "github.com/pkg/errors"

var (
	ErrConfigProviderNotReady = errors.New("config provider not ready")
	ErrInvalidCapability      = errors.New("invalid capability")
	ErrInvalidConfig          = errors.New("invalid config")
	ErrEmptyConfig            = errors.New("empty config")
)

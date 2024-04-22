package errdef

import "github.com/pkg/errors"

var (
	ErrConfigProviderNotFound = errors.New("config provider not found")
	ErrConfigProviderNotReady = errors.New("config provider not ready")
	ErrEmptyProvider          = errors.New("empty provider")
	ErrInvalidCapability      = errors.New("invalid capability")
	ErrInvalidConfig          = errors.New("invalid config")
	ErrEmptyConfig            = errors.New("empty config")
	ErrEmptyDb                = errors.New("empty db")
	ErrEmptyDbBuilder         = errors.New("no db builder specified")
	ErrValueNotSettable       = errors.New("value is not settable")
)

package errdef

import "github.com/pkg/errors"

var (
	ErrInvalidConfig    = errors.New("invalid config")
	ErrEmptyConfig      = errors.New("empty config")
	ErrEmptyDb          = errors.New("empty db")
	ErrEmptyDbBuilder   = errors.New("no db builder specified")
	ErrValueNotSettable = errors.New("value is not settable")
)

package intf

import "errors"

var (
	ErrEmptyConfig   = errors.New("empty configer")
	ErrInvalidConfig = errors.New("invalid configer")
)

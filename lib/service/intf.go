package service

import "context"

type Handler func(ctx context.Context, args ...any) (any, error)
type Func func(ctx context.Context) error

type Module interface {
	GetApp() string
	GetName() string
}

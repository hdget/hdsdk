package gokit

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/hdget/hdsdk/types"
)

type InjectFunction func(logger types.LogProvider, endpointName string) (clientOptions []interface{}, serverOptions []interface{})

// MsMiddleware 因为有的中间件是处理的时候需要做一些动作，
// 有的中间件是需要在ClientOption或者serverOption前做一些动作
// 我们将不同类型的中间件抽象出来成未MsMiddleware
type MsMiddleware struct {
	Middleware      endpoint.Middleware
	InjectFunctions map[string]InjectFunction
}

const (
	MDW_TRACER       = "trace"
	MDW_CIRCUITBREAK = "circuitbreak"
	MDW_RATELIMIT    = "ratelimit"
)

type NewMiddlewareFunc func(config *MicroServiceConfig) *MsMiddleware

var (
	NewMdwFunctions = map[string]NewMiddlewareFunc{
		MDW_TRACER:       NewMdwTrace,
		MDW_CIRCUITBREAK: NewMdwCircuitBreak,
		MDW_RATELIMIT:    NewMdwRateLimit,
	}
)

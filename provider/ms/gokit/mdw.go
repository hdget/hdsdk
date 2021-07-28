package gokit

import "github.com/go-kit/kit/endpoint"

const (
	MDW_CIRCUITBREAK = "circuitbreak"
	MDW_RATELIMIT    = "ratelimit"
)

type NewMiddlewareFunc func(config *MicroServiceConfig) endpoint.Middleware

var (
	NewMdwFunctions = map[string]NewMiddlewareFunc{
		MDW_CIRCUITBREAK: NewMdwCircuitBreak,
		MDW_RATELIMIT:    NewMdwRateLimit,
	}
)

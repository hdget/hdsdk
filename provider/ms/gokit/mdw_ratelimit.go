package gokit

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
)

// RateLimitConfig 服务限流选项
type RateLimitConfig struct {
	Limit float64 // 每秒允许多少个事件
	Burst int     // 突发最多允许多少个
}

var (
	defaultRateLimitConfig = &RateLimitConfig{
		Limit:  0,
		Burst: 0,
	}
)

// NewMdwRateLimit 服务限流, limited to 1 request per second with burst of 100 requests.
// Note, rate is defined as a number of requests per second.
func NewMdwRateLimit(config *MicroServiceConfig) endpoint.Middleware {
	if config.RateLimit == nil {
		return nil
	}

	return ratelimit.NewErroringLimiter(
		rate.NewLimiter(
			rate.Limit(config.RateLimit.Limit),
			config.RateLimit.Burst,
		),
	)
}




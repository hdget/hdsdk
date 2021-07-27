package gokit

import (
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/sony/gobreaker"
	"time"
)

// CircuitBreakConfig 服务熔断配置
type CircuitBreakConfig struct {
	MaxRequests  uint32  // 半开后最大允许通过的请求，如果最大请求为0，熔断器值允许一个请求通过
	Interval     int     // 熔断器在关闭状态周期清除内部计数器的间隔，如果为0，则在关闭状态不清除计数器
	Timeout      int     // 在熔断器半开后进入开放状态的时间，如果为0，则默认设置为60秒
	Requests     uint32  // 连续请求数量
	FailureRatio float64 // 请求失败率
}

var (
	defaultCircuitBreakConfig = &CircuitBreakConfig{
		MaxRequests:  0,
		Interval:     0,
		Timeout:      0,
		Requests:     0,
		FailureRatio: 0,
	}
)

// NewMdwCircuitBreak 服务熔断
func NewMdwCircuitBreak(config *MicroServiceConfig) endpoint.Middleware {
	if config.CircuitBreak == nil {
		return nil
	}

	settings := gobreaker.Settings{
		MaxRequests: config.CircuitBreak.MaxRequests,
		Interval:    time.Second * time.Duration(config.CircuitBreak.Interval),
		Timeout:     time.Second * time.Duration(config.CircuitBreak.Timeout),
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= config.CircuitBreak.Requests && failureRatio >= config.CircuitBreak.FailureRatio
		},
	}

	return circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(settings))
}


package ws

import (
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/v2/intf"
	"strings"
	"time"
)

var SkipHttpMethods = []string{
	"PRI",
	"HEAD",
}

func newLoggerMiddleware(logger intf.LoggerProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求方式
		reqMethod := c.Request.Method
		if pie.Contains(SkipHttpMethods, reqMethod) {
			c.Next()
			return
		}

		//请求路由
		reqUrl := c.Request.RequestURI
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)

		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := getRealIP(c)

		// 日志格式
		logger.Debug("http debug", "ip", clientIP, "method", reqMethod, "code", statusCode, "latency", latencyTime, "uri", reqUrl)
	}
}

// GetRealIP 获取真实IP
func getRealIP(c *gin.Context) string {
	xForwardInfo := c.GetHeader("X-Forwarded-For")
	if xForwardInfo != "" {
		ips := strings.Split(xForwardInfo, ",")
		// 返回第一个真实IP
		if len(ips) >= 1 {
			return ips[0]
		}
	}
	return c.ClientIP()
}

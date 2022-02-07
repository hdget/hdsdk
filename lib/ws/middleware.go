package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils"
	"time"
)

var SkipHttpMethods = []string{
	"PRI",
	"HEAD",
}

func mdwLogger(logger types.LogProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求方式
		reqMethod := c.Request.Method
		if utils.StringSliceContains(SkipHttpMethods, reqMethod) {
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
		clientIP := GetRealIP(c)

		// 日志格式
		logger.Debug("http debug", "ip", clientIP, "method", reqMethod, "code", statusCode, "latency", latencyTime, "uri", reqUrl)
	}
}

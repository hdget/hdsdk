package hdgin

import (
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

// // GetRealIP 获取真实IP
func GetRealIP(c *gin.Context) string {
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

// IsPrivateIp 检查是否时内网ip
func IsPrivateIp(ipStr string) bool {
	address := net.ParseIP(ipStr)
	if address == nil {
		return false
	}

	if address.IsLoopback() || address.IsLinkLocalUnicast() || address.IsLinkLocalMulticast() {
		return true
	}

	return address.IsPrivate()
}

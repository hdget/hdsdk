package ws

import (
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	bearerPrefix        = "bearer "
	bearerPrefixLength  = len(bearerPrefix)
	headerAuthorization = "Authorization"
)

// GetBearerToken get bearer token from http header
func GetBearerToken(c *gin.Context) string {
	authHeader := c.Request.Header.Get(headerAuthorization)
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(strings.ToLower(authHeader), bearerPrefix) {
		return ""
	}

	return authHeader[bearerPrefixLength:]
}

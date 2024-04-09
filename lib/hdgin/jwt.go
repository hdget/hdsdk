package hdgin

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

var (
	bearerPrefix        = "bearer"
	headerAuthorization = "Authorization"

	ErrTokenNotFound                  = errors.New("token not found")
	ErrInvalidHttpAuthorizationHeader = errors.New("invalid http authorization header")
)

// GetAuthorizationToken get bearer token or other token from http header
func GetAuthorizationToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get(headerAuthorization)
	if authHeader == "" {
		return "", ErrTokenNotFound
	}

	authHeaderParts := strings.Split(authHeader, " ")
	switch len(authHeaderParts) {
	case 1: // Authorization: Token
		return authHeaderParts[0], nil
	case 2: // Authorization: Bearer Token
		if strings.HasSuffix(strings.ToLower(authHeaderParts[0]), bearerPrefix) {
			return authHeaderParts[1], nil
		}
	}

	return "", ErrInvalidHttpAuthorizationHeader
}

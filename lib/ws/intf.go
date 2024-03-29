package ws

import "github.com/gin-gonic/gin"

type Route struct {
	Method  HttpMethod
	Path    string
	Handler gin.HandlerFunc
}

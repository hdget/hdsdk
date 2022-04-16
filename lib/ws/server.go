package ws

import (
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type HttpMethod int

const (
	Get HttpMethod = iota
	Post
	Delete
)
const waitTime = 3 * time.Second

var (
	ErrDuplicateRouterGroup = errors.New("duplicate router group")
	ErrRouterGroupNotFound  = errors.New("router group not found")
)

// SetReleaseMode set gin to release mode
func SetReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
}

// SetDebugMode set gin to debug mode
func SetDebugMode() {
	gin.SetMode(gin.DebugMode)
}

// SetTestMode set gin to test mode
func SetTestMode() {
	gin.SetMode(gin.TestMode)
}

func addToRouterGroup(routerGroup *gin.RouterGroup, method HttpMethod, path string, handler gin.HandlerFunc) {
	switch method {
	case Get:
		routerGroup.GET(path, handler)
	case Post:
		routerGroup.POST(path, handler)
	case Delete:
		routerGroup.DELETE(path, handler)
	}
}

func addToRouter(router *gin.Engine, method HttpMethod, path string, handler gin.HandlerFunc) {
	switch method {
	case Get:
		router.GET(path, handler)
	case Post:
		router.POST(path, handler)
	case Delete:
		router.DELETE(path, handler)
	}
}

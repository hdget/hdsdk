package ws

import (
	"github.com/gin-gonic/gin"
	"hdsdk/types"
	"io"
)

type Route struct {
	Method  HttpMethod
	Path    string
	Handler gin.HandlerFunc
}

// NewRouter create a gin router
func NewRouter(logger types.LogProvider) *gin.Engine {
	// new router
	router := gin.New()

	// set gin to logout to stdout and file
	gin.DefaultWriter = io.MultiWriter(logger.GetStdLogger().Writer())

	// add basic middleware
	router.Use(
		gin.Recovery(),
		mdwLogger(logger), // logger middleware
	)

	return router
}

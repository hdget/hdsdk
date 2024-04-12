package ws

import (
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
	"net/http"
)

type httpServerImpl struct {
	*baseServer
}

func NewHttpServer(logger intf.LoggerProvider, address string, options ...ServerOption) (WebServer, error) {
	return &httpServerImpl{
		baseServer: newBaseServer(logger, address, options...),
	}, nil
}

func (w httpServerImpl) Start() error {
	if err := w.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

package ws

import (
	"github.com/hdget/hdsdk/types"
	"github.com/pkg/errors"
	"net/http"
)

type httpServerImpl struct {
	*baseServer
}

func NewHttpServer(logger types.LogProvider, address string, options ...ServerOption) (WebServer, error) {
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

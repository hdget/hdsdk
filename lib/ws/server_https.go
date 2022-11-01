package ws

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kabukky/httpscerts"
	"github.com/pkg/errors"
	"hdsdk"
	"hdsdk/types"
	"hdsdk/utils/parallel"
	"net/http"
	"syscall"
)

type HttpsServer struct {
	*HttpServer
	CertPath string
	KeyPath  string
}

func NewHttpsServer(logger types.LogProvider, certPath, keyPath, address string) (WebServer, error) {
	srv := &HttpsServer{
		HttpServer: createHttpServer(logger, address),
		CertPath:   certPath,
		KeyPath:    keyPath,
	}

	// Check if the cert files are available.
	if err := httpscerts.Check(srv.CertPath, srv.KeyPath); err != nil {
		// If they are not available, generate new ones.
		if err = httpscerts.Generate(srv.CertPath, srv.KeyPath, address); err != nil {
			return nil, errors.Wrap(err, "generate secure credential")
		}
	}

	return srv, nil
}

func (srv *HttpsServer) Run() {
	listenFunc := func() error {
		return srv.ListenAndServeTLS(srv.CertPath, srv.KeyPath)
	}

	var group parallel.Group
	{
		group.Add(listenFunc, srv.shutdown)
	}
	{
		group.Add(
			parallel.SignalActor(
				context.Background(),
				syscall.SIGINT,
				syscall.SIGQUIT,
				syscall.SIGTERM,
				syscall.SIGKILL,
			),
		)
	}

	if err := group.Run(); err != nil && err != http.ErrServerClosed {
		hdsdk.Logger.Error("https server quit", "error", err)
	}
}

func createHttpServer(logger types.LogProvider, address string) *HttpServer {
	router := NewRouter(logger)
	return &HttpServer{
		Server: &http.Server{
			Addr:    address,
			Handler: router,
		},
		router:       router,
		routerGroups: make(map[string]*gin.RouterGroup),
	}
}

package ws

import (
	"context"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdsdk/utils/parallel"
	"github.com/kabukky/httpscerts"
	"github.com/pkg/errors"
	"net/http"
	"syscall"
)

type MyHttpsServer struct {
	*MyHttpServer
	CertPath string
	KeyPath  string
}

func NewHttpsServer(logger types.LogProvider, certPath, keyPath, address string) (*MyHttpsServer, error) {
	srv := &MyHttpsServer{
		MyHttpServer: NewHttpServer(logger, address),
		CertPath:     certPath,
		KeyPath:      keyPath,
	}

	err := srv.setupCerts(address)
	if err != nil {
		return nil, err
	}

	return &MyHttpsServer{
		MyHttpServer: NewHttpServer(logger, address),
		CertPath:     certPath,
		KeyPath:      keyPath,
	}, nil
}

func (srv *MyHttpsServer) Run() {
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

func (srv *MyHttpsServer) setupCerts(address string) error {
	// Check if the cert files are available.
	if err := httpscerts.Check(srv.CertPath, srv.KeyPath); err != nil {
		// If they are not available, generate new ones.
		if err = httpscerts.Generate(srv.CertPath, srv.KeyPath, address); err != nil {
			return errors.Wrap(err, "generate secure credential")
		}
	}
	return nil
}

package ws

import (
	"context"
	"github.com/hdget/hdsdk"
	"github.com/hdget/hdsdk/types"
	"github.com/hdget/hdutils/parallel"
	"github.com/kabukky/httpscerts"
	"github.com/pkg/errors"
	"net/http"
	"syscall"
)

type httpsServerImpl struct {
	*baseServer
	CertPath string
	KeyPath  string
}

func NewHttpsServer(logger types.LogProvider, address, certPath, keyPath string, options ...ServerOption) (WebServer, error) {
	// Check if the cert files are available.
	if err := httpscerts.Check(certPath, keyPath); err != nil {
		// If they are not available, generate new ones.
		if err = httpscerts.Generate(certPath, keyPath, address); err != nil {
			return nil, errors.Wrap(err, "generate secure credential")
		}
	}

	return &httpsServerImpl{
		baseServer: newBaseServer(logger, address),
		CertPath:   certPath,
		KeyPath:    keyPath,
	}, nil
}

func (w httpsServerImpl) Run() {
	listenFunc := func() error {
		return w.ListenAndServeTLS(w.CertPath, w.KeyPath)
	}

	var group parallel.Group
	{
		group.Add(listenFunc, w.shutdown)
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

	if err := group.Run(); err != nil && errors.Is(err, http.ErrServerClosed) {
		hdsdk.Logger.Error("https server quit", "error", err)
	}
}

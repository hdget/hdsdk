package hdgin

import (
	"github.com/hdget/hdsdk/v1/intf"
	"github.com/kabukky/httpscerts"
	"github.com/pkg/errors"
	"net/http"
)

type httpsServerImpl struct {
	*baseServer
	CertPath string
	KeyPath  string
}

func NewHttpsServer(logger intf.LoggerProvider, address, certPath, keyPath string, options ...ServerOption) (WebServer, error) {
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

func (w httpsServerImpl) Start() error {
	if err := w.ListenAndServeTLS(w.CertPath, w.KeyPath); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

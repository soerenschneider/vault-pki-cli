package backends

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
)

type MultiBackend struct {
	backends []pki.CertBackend
}

func NewMultiBackend(backends ...pki.CertBackend) (*MultiBackend, error) {
	if nil == backends || len(backends) == 0 {
		return nil, errors.New("no backends supplied")
	}

	return &MultiBackend{backends: backends}, nil
}

func (b *MultiBackend) Write(certData *pki.CertData) error {
	for _, backend := range b.backends {
		err := backend.Write(certData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MultiBackend) Read() (*x509.Certificate, error) {
	var err error

	for _, backend := range b.backends {
		var cert *x509.Certificate
		cert, err = backend.Read()
		if err == nil {
			return cert, err
		}
	}

	return nil, err
}

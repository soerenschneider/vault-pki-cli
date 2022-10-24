package sink

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
)

type MultiSink struct {
	sinks []pki.CertSink
}

func NewMultiSink(sinks ...pki.CertSink) (*MultiSink, error) {
	if nil == sinks || len(sinks) == 0 {
		return nil, errors.New("no sink supplied")
	}

	return &MultiSink{sinks: sinks}, nil
}

func (b *MultiSink) Write(certData *pki.CertData) error {
	for _, backend := range b.sinks {
		err := backend.Write(certData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *MultiSink) Read() (*x509.Certificate, error) {
	var err error

	for _, backend := range b.sinks {
		var cert *x509.Certificate
		cert, err = backend.Read()
		if err == nil {
			return cert, err
		}
	}

	return nil, err
}

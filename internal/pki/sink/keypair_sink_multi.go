package sink

import (
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
)

type MultiKeyPairSink struct {
	sinks []*KeyPairSink
}

func MultiKeyPairSinkFromConfig(config *conf.Config) (*MultiKeyPairSink, error) {
	sinks, err := KeyPairSinkFromConfig(config)
	if err != nil {
		return nil, err
	}

	return NewMultiKeyPairSink(sinks...)
}

func NewMultiKeyPairSink(sinks ...*KeyPairSink) (*MultiKeyPairSink, error) {
	if nil == sinks {
		return nil, errors.New("no sinks provided")
	}

	return &MultiKeyPairSink{sinks: sinks}, nil
}

func (f *MultiKeyPairSink) WriteCert(certData *pki.CertData) error {
	var errs []error
	for _, sink := range f.sinks {
		err := sink.WriteCert(certData)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors writing: %v", errs)
	}
	return nil
}

func (f *MultiKeyPairSink) ReadCert() (*x509.Certificate, error) {
	for _, sink := range f.sinks {
		cert, err := sink.ReadCert()
		if err == nil {
			return cert, err
		}
	}

	return nil, errors.New("could not read any cert")
}

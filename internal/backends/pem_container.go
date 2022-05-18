package backends

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg"
)

type PemContainerBackend struct {
	container pki.KeyPod
}

func NewPemContainerBackend(pod pki.KeyPod) (*PemContainerBackend, error) {
	if nil == pod {
		return nil, errors.New("empty pod provided")
	}

	return &PemContainerBackend{container: pod}, nil
}

func (f *PemContainerBackend) Write(certData *pki.CertData) error {
	if certData == nil {
		return errors.New("empty certData provided")
	}

	concatenatedPemData := certData.AsContainer()
	return f.container.Write([]byte(concatenatedPemData))
}

func (f *PemContainerBackend) Read() (*x509.Certificate, error) {
	data, err := f.container.Read()
	if err != nil {
		return nil, err
	}
	return pkg.ParseCertPem(data)
}

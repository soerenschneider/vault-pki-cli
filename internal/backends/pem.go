package backends

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"log"
)

type PemBackend struct {
	ca         pki.KeyPod
	cert       pki.KeyPod
	privateKey pki.KeyPod
}

func NewPemBackend(cert, privateKey, chain pki.KeyPod) (*PemBackend, error) {
	if nil == cert {
		return nil, errors.New("empty cert pod provided")
	}

	if nil == privateKey {
		return nil, errors.New("empty private key pod provided")
	}

	return &PemBackend{cert: cert, privateKey: privateKey, ca: chain}, nil
}

func (f *PemBackend) Write(certData *pki.CertData) error {
	if certData.HasCaChain() && f.ca != nil {
		log.Println("--------------------------")
		if err := f.ca.Write(append(certData.CaChain, "\n"...)); err != nil {
			return err
		}
	}

	if err := f.cert.Write(append(certData.Certificate, "\n"...)); err != nil {
		return err
	}

	if certData.HasPrivateKey() {
		return f.privateKey.Write(append(certData.PrivateKey, "\n"...))
	}

	return nil
}

func (f *PemBackend) Read() (*x509.Certificate, error) {
	data, err := f.cert.Read()
	if err != nil {
		return nil, err
	}

	return pkg.ParseCertPem(data)
}

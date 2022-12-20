package sink

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"log"
)

type PemSink struct {
	ca         pki.StorageImplementation
	cert       pki.StorageImplementation
	privateKey pki.StorageImplementation
}

func NewPemSink(cert, privateKey, chain pki.StorageImplementation) (*PemSink, error) {
	if nil == cert {
		return nil, errors.New("empty cert storage provided")
	}

	if nil == privateKey {
		return nil, errors.New("empty private key storage provided")
	}

	return &PemSink{cert: cert, privateKey: privateKey, ca: chain}, nil
}

func (f *PemSink) WriteCert(certData *pki.CertData) error {
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

func (f *PemSink) ReadCert() (*x509.Certificate, error) {
	data, err := f.cert.Read()
	if err != nil {
		return nil, err
	}

	return pkg.ParseCertPem(data)
}

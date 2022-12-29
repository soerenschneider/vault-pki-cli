package sink

import (
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg"
)

// CsrSink offers an interface to read/write keypair data (certificate and private key) and optional ca data.
type CsrSink struct {
	ca   pki.StorageImplementation
	cert pki.StorageImplementation
	csr  pki.StorageImplementation
}

const (
	csrId = "csr"
)

func CsrSinkFromConfig(storageConfig []map[string]string) (*CsrSink, error) {
	var certVal string
	var csrVal string
	var caVal string
	for _, conf := range storageConfig {
		val, ok := conf[certId]
		if !ok {
			return nil, fmt.Errorf("can not build storage, missing '%s' in storage configuration", certId)
		}
		certVal = val

		val, ok = conf[csrId]
		if !ok {
			return nil, fmt.Errorf("can not build storage, missing '%s' in storage configuration", csrId)
		}
		csrVal = val

		val, ok = conf[caId]
		if ok {
			caVal = val
		}
	}

	builder, err := storage.GetBuilder()
	if err != nil {
		return nil, err
	}
	certStorageImpl, err := builder.BuildFromUri(certVal)
	if err != nil {
		return nil, err
	}

	csrStorageImpl, err := builder.BuildFromUri(csrVal)
	if err != nil {
		return nil, err
	}

	var caStorageImpl pki.StorageImplementation
	if len(caVal) > 0 {
		caStorageImpl, err = builder.BuildFromUri(caVal)
		if err != nil {
			return nil, err
		}
	}

	return NewCsrSink(certStorageImpl, csrStorageImpl, caStorageImpl)
}

func NewCsrSink(cert, csr, chain pki.StorageImplementation) (*CsrSink, error) {
	if nil == cert {
		return nil, errors.New("empty cert storage provided")
	}

	if nil == csr {
		return nil, errors.New("empty private key storage provided")
	}

	return &CsrSink{cert: cert, csr: csr, ca: chain}, nil
}

func (f *CsrSink) WriteCert(certData *pki.CertData) error {
	if certData.HasCaData() && f.ca != nil {
		if err := f.ca.Write(append(certData.CaData, "\n"...)); err != nil {
			return err
		}
	}

	if err := f.cert.Write(append(certData.Certificate, "\n"...)); err != nil {
		return err
	}

	return nil
}

func (f *CsrSink) ReadCert() (*x509.Certificate, error) {
	data, err := f.cert.Read()
	if err != nil {
		return nil, err
	}

	return pkg.ParseCertPem(data)
}

func (f *CsrSink) ReadCsr() ([]byte, error) {
	return f.csr.Read()
}

func (f *CsrSink) WriteSignature(cert *pki.Signature) error {
	if cert.HasCaData() && f.ca != nil {
		if err := f.ca.Write(append(cert.CaData, "\n"...)); err != nil {
			return err
		}
	}

	if err := f.cert.Write(append(cert.Certificate, "\n"...)); err != nil {
		return err
	}

	return nil
}

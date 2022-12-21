package sink

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
)

type YubikeyBackend struct {
	yubikey pki.StorageImplementation
}

func NewYubikeySink(pod *storage.YubikeyPod) (*YubikeyBackend, error) {
	if nil == pod {
		return nil, errors.New("empty yubikey pod provided")
	}

	return &YubikeyBackend{yubikey: pod}, nil
}

func (f *YubikeyBackend) WriteCert(certData *pki.CertData) error {
	var dataPortion []byte

	if certData.HasCaData() {
		dataPortion = append(dataPortion, certData.CaData...)
		dataPortion = append(dataPortion, []byte("\n")...)
	}

	dataPortion = append(dataPortion, certData.Certificate...)
	dataPortion = append(dataPortion, []byte("\n")...)

	if certData.HasPrivateKey() {
		dataPortion = append(dataPortion, certData.PrivateKey...)
	}

	return f.yubikey.Write(dataPortion)
}

func (f *YubikeyBackend) ReadCert() (*x509.Certificate, error) {
	cert, err := f.yubikey.Read()
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(cert)
}

package sink

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/pods"
)

type YubikeyBackend struct {
	yubikey pki.KeyPod
}

func NewYubikeySink(pod *pods.YubikeyPod) (*YubikeyBackend, error) {
	if nil == pod {
		return nil, errors.New("empty yubikey pod provided")
	}

	return &YubikeyBackend{yubikey: pod}, nil
}

func (f *YubikeyBackend) Write(certData *pki.CertData) error {
	var dataPortion []byte

	if certData.HasCaChain() {
		dataPortion = append(dataPortion, certData.CaChain...)
		dataPortion = append(dataPortion, []byte("\n")...)
	}

	dataPortion = append(dataPortion, certData.Certificate...)
	dataPortion = append(dataPortion, []byte("\n")...)

	if certData.HasPrivateKey() {
		dataPortion = append(dataPortion, certData.PrivateKey...)
	}

	return f.yubikey.Write(dataPortion)
}

func (f *YubikeyBackend) Read() (*x509.Certificate, error) {
	cert, err := f.yubikey.Read()
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(cert)
}

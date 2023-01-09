package sink

import (
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg"
)

// KeyPairSink offers an interface to read/write keypair data (certificate and private key) and optional ca data.
type KeyPairSink struct {
	ca         pki.StorageImplementation
	cert       pki.StorageImplementation
	privateKey pki.StorageImplementation
}

const (
	certId = "cert"
	keyId  = "key"
	caId   = "ca"
)

func KeyPairSinkFromConfig(config *conf.Config) ([]*KeyPairSink, error) {
	var sinks []*KeyPairSink
	for _, conf := range config.StorageConfig {
		var certVal string
		var keyVal string
		var caVal string

		val, ok := conf[certId]
		if ok {
			certVal = val
		}

		val, ok = conf[keyId]
		if !ok {
			return nil, fmt.Errorf("can not build storage, missing '%s' in storage configuration", keyId)
		}
		keyVal = val

		val, ok = conf[caId]
		if ok {
			caVal = val
		}

		builder, err := storage.GetBuilder()
		if err != nil {
			return nil, err
		}

		var certSink pki.StorageImplementation
		if len(certVal) > 0 {
			certSink, err = builder.BuildFromUri(certVal)
			if err != nil {
				return nil, err
			}
		}

		keySink, err := builder.BuildFromUri(keyVal)
		if err != nil {
			return nil, err
		}

		var caSink pki.StorageImplementation
		if len(caVal) > 0 {
			caSink, err = builder.BuildFromUri(caVal)
			if err != nil {
				return nil, err
			}
		}

		sink, err := NewKeyPairSink(certSink, keySink, caSink)
		if err != nil {
			return nil, fmt.Errorf("can't build individual sink: %v", err)
		}

		sinks = append(sinks, sink)
	}

	return sinks, nil
}

func NewKeyPairSink(cert, privateKey, chain pki.StorageImplementation) (*KeyPairSink, error) {
	if nil == privateKey {
		return nil, errors.New("empty private key storage provided")
	}

	if nil == cert && nil != chain {
		return nil, errors.New("please specify either a storage for the certificate or remove the storage for the chain")
	}

	return &KeyPairSink{cert: cert, privateKey: privateKey, ca: chain}, nil
}

func (f *KeyPairSink) WriteCert(certData *pki.CertData) error {
	if nil == certData {
		return errors.New("got nil as certData")
	}

	if f.cert == nil {
		if !certData.HasPrivateKey() {
			return errors.New("WriteCert(): can not write data, cert data contains no private key and no cert storage provided")
		}

		var data []byte
		if certData.HasCaData() {
			data = append(certData.CaData, "\n"...)
		}
		data = append(data, append(certData.Certificate, "\n"...)...)
		data = append(data, append(certData.PrivateKey, "\n"...)...)

		return f.privateKey.Write(data)
	}

	if certData.HasCaData() && f.ca != nil {
		if err := f.ca.Write(append(certData.CaData, "\n"...)); err != nil {
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

func (f *KeyPairSink) ReadCert() (*x509.Certificate, error) {
	data, err := f.cert.Read()
	if err != nil {
		return nil, err
	}

	return pkg.ParseCertPem(data)
}

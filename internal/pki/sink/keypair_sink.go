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
		sink, err := buildSink(conf)
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, sink)
	}

	return sinks, nil
}

func NewKeyPairSink(cert, privateKey, chain pki.StorageImplementation) (*KeyPairSink, error) {
	if nil == privateKey {
		return nil, errors.New("empty private key storage provided")
	}

	return &KeyPairSink{cert: cert, privateKey: privateKey, ca: chain}, nil
}

func (f *KeyPairSink) ReadCert() (*x509.Certificate, error) {
	var source pki.StorageImplementation
	if f.cert != nil {
		source = f.cert
	} else {
		source = f.privateKey
	}

	data, err := source.Read()
	if err != nil {
		return nil, err
	}

	return pkg.ParseCertPem(data)
}

func (f *KeyPairSink) WriteCert(certData *pki.CertData) error {
	if nil == certData {
		return errors.New("got nil as certData")
	}

	// case 1: write cert, ca and private key to same storage
	if f.cert == nil && f.ca == nil {
		return f.writeToPrivateSlot(certData)
	}

	// case 2: write cert and private to a same storage, write ca (if existent) to dedicated storage
	if f.cert == nil && f.ca != nil {
		return f.writeToCertAndCaSlot(certData)
	}

	// case 3: write to individual storage
	return f.writeToIndividualSlots(certData)
}

func (f *KeyPairSink) writeToPrivateSlot(certData *pki.CertData) error {
	var data []byte
	data = append(certData.Certificate, "\n"...)
	if certData.HasCaData() {
		data = append(data, append(certData.CaData, "\n"...)...)
	}
	data = append(data, append(certData.PrivateKey, "\n"...)...)

	return f.privateKey.Write(data)
}

func (f *KeyPairSink) writeToCertAndCaSlot(certData *pki.CertData) error {
	var data []byte
	data = append(certData.Certificate, "\n"...)
	data = append(data, append(certData.PrivateKey, "\n"...)...)

	if err := f.privateKey.Write(data); err != nil {
		return err
	}

	if certData.HasCaData() {
		return f.ca.Write(append(certData.CaData, "\n"...))
	}

	return nil
}

func (f *KeyPairSink) writeToIndividualSlots(certData *pki.CertData) error {
	var certRaw []byte
	certRaw = append(certData.Certificate, "\n"...)
	if certData.HasCaData() && f.ca == nil {
		certRaw = append(certRaw, append(certData.CaData, "\n"...)...)
	}

	if err := f.cert.Write(certRaw); err != nil {
		return err
	}

	if certData.HasCaData() && f.ca != nil {
		if err := f.ca.Write(append(certData.CaData, "\n"...)); err != nil {
			return err
		}
	}

	if certData.HasPrivateKey() {
		return f.privateKey.Write(append(certData.PrivateKey, "\n"...))
	}

	return nil
}

func buildSink(conf map[string]string) (*KeyPairSink, error) {
	builder, err := storage.GetBuilder()
	if err != nil {
		return nil, err
	}

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

	return NewKeyPairSink(certSink, keySink, caSink)
}

package pki

import (
	"crypto/x509"
)

// IssueSink defines pluggable sink to write certificate data to.
type IssueSink interface {
	WriteCert(cert *CertData) error
	ReadCert() (*x509.Certificate, error)
}

type CrlSink interface {
	WriteCrl(crlData []byte) error
}

type CsrSink interface {
	ReadCsr() ([]byte, error)
	WriteCert(cert []byte) error
}

type CaSink interface {
	WriteCa(certData []byte) error
}

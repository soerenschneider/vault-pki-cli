package pki

import (
	"crypto/x509"
)

// Sinks are providing storage interaction for different subcommands (issue, sign, ...). The sinks in turn re-use
// concrete storage implementations of `StorageImplementation`.

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
	WriteSignature(cert *Signature) error
}

type CaSink interface {
	WriteCa(certData []byte) error
}

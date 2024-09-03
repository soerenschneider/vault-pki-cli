package pki

import (
	"crypto/x509"

	"github.com/soerenschneider/vault-pki-cli/pkg"
)

// Sinks are providing storage interaction for different subcommands (issue, sign, ...). The sinks in turn re-use
// concrete storage implementations of `StorageImplementation`.

// IssueStorage defines pluggable sink to write certificate data to.
type IssueStorage interface {
	WriteCert(cert *pkg.CertData) error
	ReadCert() (*x509.Certificate, error)
}

type CrlStorage interface {
	WriteCrl(crlData []byte) error
}

type CsrStorage interface {
	ReadCsr() ([]byte, error)
	WriteSignature(cert *pkg.Signature) error
}

type CaStorage interface {
	WriteCa(certData []byte) error
}

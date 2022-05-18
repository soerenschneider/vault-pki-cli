package pki

import (
	"crypto/x509"
)

// CertBackend defines pluggable backends to write certificate data to.
type CertBackend interface {
	Write(cert *CertData) error
	Read() (*x509.Certificate, error)
}

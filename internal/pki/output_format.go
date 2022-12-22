package pki

import (
	"crypto/x509"
)

// CertSink defines pluggable sink to write certificate data to.
type CertSink interface {
	Write(cert *CertData) error
	Read() (*x509.Certificate, error)
}

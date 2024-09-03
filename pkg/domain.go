package pkg

import (
	"crypto/x509"
	"errors"
	"strings"
)

var (
	ErrNoCertFound     = errors.New("no existing cert found")
	ErrRunHook         = errors.New("error running hook")
	ErrCertInvalidData = errors.New("could not parse cert data")
	ErrWriteCert       = errors.New("could not write certificate data")
	ErrIssueCert       = errors.New("error while issuing cert")
	ErrRevokeCert      = errors.New("error while revoking cert")
	ErrSignCert        = errors.New("error while signing cert")
	ErrTidyCert        = errors.New("error while tidying up cert storage")
)

type IssueStatus int

const (
	Issued  IssueStatus = iota
	Noop    IssueStatus = iota
	Unknown IssueStatus = iota
)

type IssueResult struct {
	ExistingCert *x509.Certificate
	IssuedCert   *x509.Certificate
	Status       IssueStatus
}

type SignatureArgs struct {
	CommonName string
	Ttl        string
	IpSans     []string
	AltNames   []string
}

type IssueArgs struct {
	CommonName string
	Ttl        string
	IpSans     []string
	AltNames   []string
}

type CertData struct {
	PrivateKey  []byte
	Certificate []byte
	CaData      []byte
	Csr         []byte
}

func (certData *CertData) AsContainer() string {
	var buffer strings.Builder

	if certData.HasCaData() {
		buffer.Write(certData.CaData)
		buffer.Write([]byte("\n"))
	}

	buffer.Write(certData.Certificate)
	buffer.Write([]byte("\n"))

	if certData.HasPrivateKey() {
		buffer.Write(certData.PrivateKey)
		buffer.Write([]byte("\n"))
	}

	return buffer.String()
}

func (cert *CertData) HasPrivateKey() bool {
	return len(cert.PrivateKey) > 0
}

func (cert *CertData) HasCertificate() bool {
	return len(cert.Certificate) > 0
}

func (cert *CertData) HasCaData() bool {
	return len(cert.CaData) > 0
}

type Signature struct {
	Certificate []byte
	CaData      []byte
	Serial      string
}

func (cert *Signature) HasCaData() bool {
	return len(cert.CaData) > 0
}

package pki

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/cenkalti/backoff/v3"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/renew_strategy"
	"golang.org/x/net/context"
)

// StorageImplementation is a simple wrapper around a key artifact (cert, key, ca, crl, csr). This enables decoupling
// from the actual resource (file-based, kubernetes, network, ..) and make it interchangeable.
type StorageImplementation interface {
	Read() ([]byte, error)
	CanRead() error
	Write([]byte) error
	CanWrite() error
}

type PkiClient interface {
	// Issue issues a new certificate from the PKI
	Issue(ctx context.Context, args pkg.IssueArgs) (*pkg.CertData, error)

	// Sign signs a CSR
	Sign(ctx context.Context, csr string, args pkg.SignatureArgs) (*pkg.Signature, error)

	// Revoke revokes a certificate by its serial number
	Revoke(ctx context.Context, serial string) error

	// ReadAcme reads a previously acquired letsencrypt certificate from Vault
	ReadAcme(ctx context.Context, commonName string) (*pkg.CertData, error)

	// Tidy cleans up the PKI blob storage of dangling certificates
	Tidy(ctx context.Context) error

	// FetchCa returns the CA for the configured mount
	FetchCa(binary bool) ([]byte, error)

	// FetchCaChain returns the whole CA chain for the configured mount
	FetchCaChain() ([]byte, error)

	// FetchCrl returns the CRL of the configured mount
	FetchCrl(binary bool) ([]byte, error)
}

type RenewStrategy interface {
	Renew(cert *x509.Certificate) (bool, error)
}

type PkiService struct {
	pkiImpl  PkiClient
	strategy RenewStrategy
}

func NewPkiService(pki PkiClient, strategy RenewStrategy) (*PkiService, error) {
	if pki == nil {
		return nil, errors.New("empty pki impl provided")
	}

	if strategy == nil {
		strategy = &renew_strategy.StaticRenewal{Decision: true}
	}

	return &PkiService{
		pkiImpl:  pki,
		strategy: strategy,
	}, nil
}

func (p *PkiService) Revoke(ctx context.Context, serial string) error {
	if len(serial) == 0 {
		return errors.New("can't revoke, empty cert serial provided")
	}

	op := func() error {
		return p.pkiImpl.Revoke(ctx, serial)
	}

	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrRevokeCert, err)
	}

	return nil
}

func (p *PkiService) Tidy(ctx context.Context) error {
	op := func() error {
		return p.pkiImpl.Tidy(ctx)
	}

	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrTidyCert, err)
	}

	return nil
}

func (p *PkiService) ReadAcme(ctx context.Context, format IssueStorage, commonName string) (pkg.IssueResult, error) {
	ret := pkg.IssueResult{
		Status: pkg.Unknown,
	}

	var err error
	ret.ExistingCert, err = format.ReadCert()
	if err == nil && ret.ExistingCert == nil {
		ret.Status = pkg.Noop
	}

	var cert *pkg.CertData
	op := func() error {
		var err error
		cert, err = p.pkiImpl.ReadAcme(ctx, commonName)
		return err
	}

	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return ret, fmt.Errorf("%w: %v", pkg.ErrIssueCert, err)
	}

	ret.IssuedCert, err = pkg.ParseCertPem(cert.Certificate)
	if err != nil {
		return ret, fmt.Errorf("received cert data invalid: %w: %v", pkg.ErrCertInvalidData, err)
	}

	if ret.ExistingCert != nil {
		if !bytes.Equal(ret.ExistingCert.Raw, ret.IssuedCert.Raw) {
			ret.Status = pkg.Issued
		}
	}

	if err := format.WriteCert(cert); err != nil {
		return ret, fmt.Errorf("%w: %v", pkg.ErrWriteCert, err)
	}

	return ret, nil
}

func (p *PkiService) shouldIssue(cert *x509.Certificate) (bool, error) {
	if cert == nil {
		return true, errors.New("nil pointer supplied")
	}

	if !pkg.IsCertExpired(*cert) {
		if err := p.Verify(cert); err != nil {
			return true, fmt.Errorf("cert exists but can not be verified against ca: %w", err)
		}
	}

	return p.strategy.Renew(cert)
}

func (p *PkiService) Issue(ctx context.Context, format IssueStorage, args pkg.IssueArgs) (pkg.IssueResult, error) {
	ret := pkg.IssueResult{
		Status: pkg.Unknown,
	}

	var err error
	ret.ExistingCert, err = format.ReadCert()
	if (err != nil || ret.ExistingCert == nil) && !errors.Is(err, pkg.ErrNoCertFound) {
		log.Warn().Err(err).Msg("Could not read certificate")
	}

	issueNewCert, err := p.shouldIssue(ret.ExistingCert)
	if ret.ExistingCert != nil && err == nil && !issueNewCert {
		ret.Status = pkg.Noop
		return ret, nil
	}

	var issuedCertData *pkg.CertData
	op := func() error {
		var err error
		issuedCertData, err = p.pkiImpl.Issue(ctx, args)
		return err
	}
	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return ret, fmt.Errorf("%w: %v", pkg.ErrIssueCert, err)
	}

	ret.IssuedCert, err = pkg.ParseCertPem(issuedCertData.Certificate)
	if err != nil {
		return ret, fmt.Errorf("received cert data invalid: %w: %v", pkg.ErrCertInvalidData, err)
	}

	if err := format.WriteCert(issuedCertData); err != nil {
		return ret, fmt.Errorf("%w: %v", pkg.ErrWriteCert, err)
	}

	ret.Status = pkg.Issued
	return ret, nil
}

func (p *PkiService) Verify(cert *x509.Certificate) error {
	var caData []byte
	op := func() error {
		var err error
		caData, err = p.pkiImpl.FetchCaChain()
		return err
	}

	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return err
	}

	caBlock, _ := pem.Decode(caData)
	ca, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return err
	}

	return verifyCertAgainstCa(cert, ca)
}

func verifyCertAgainstCa(cert, ca *x509.Certificate) error {
	if cert == nil || ca == nil {
		return errors.New("empty cert(s) supplied")
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(ca)

	verifyOptions := x509.VerifyOptions{
		Roots: certPool,
	}
	_, err := cert.Verify(verifyOptions)
	return err
}

func (p *PkiService) Sign(ctx context.Context, sink CsrStorage, args pkg.SignatureArgs) error {
	csr, err := sink.ReadCsr()
	if err != nil {
		return err
	}

	var resp *pkg.Signature
	op := func() error {
		var err error
		resp, err = p.pkiImpl.Sign(ctx, string(csr), args)
		return err
	}

	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrSignCert, err)
	}

	_, err = pkg.ParseCertPem(resp.Certificate)
	if err != nil {
		return fmt.Errorf("received cert data invalid: %w: %v", pkg.ErrCertInvalidData, err)
	}

	err = sink.WriteSignature(resp)
	if err != nil {
		return fmt.Errorf("%w: %v", pkg.ErrWriteCert, err)
	}

	return nil
}

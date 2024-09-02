package pki

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"
	"golang.org/x/net/context"
)

type IssueOutcome int

const (
	Issued    = 0
	NotNeeded = 1
	Error     = 2
)

var ErrNoCertFound = errors.New("data not found")

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

type PkiService struct {
	pkiImpl  PkiClient
	strategy issue_strategies.IssueStrategy
}

func NewPkiService(pki PkiClient, strategy issue_strategies.IssueStrategy) (*PkiService, error) {
	if pki == nil {
		return nil, errors.New("empty pki impl provided")
	}

	if strategy == nil {
		strategy = &issue_strategies.StaticRenewal{Decision: true}
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

	log.Info().Msgf("Attempting to revoke certificate %s", serial)
	op := func() error {
		return p.pkiImpl.Revoke(ctx, serial)
	}

	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return fmt.Errorf("could not revoke certificate: %w", err)
	}

	log.Info().Msgf("Revoking certificate successful")
	return nil
}

func (p *PkiService) Tidy(ctx context.Context) error {
	op := func() error {
		return p.pkiImpl.Tidy(ctx)
	}

	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return fmt.Errorf("could not tidy certificate storage: %w", err)
	}

	log.Info().Msgf("Tidy blob storage scheduled")
	return nil
}

func (p *PkiService) ReadAcme(ctx context.Context, format IssueStorage, commonName string) (bool, error) {
	var changed bool
	certData, err := format.ReadCert()
	if err != nil || certData == nil {
		log.Info().Msg("No existing local certdata available")
		changed = true
	}

	log.Info().Msgf("Trying to read certificate for domain '%s'", commonName)
	var cert *pkg.CertData
	op := func() error {
		var err error
		cert, err = p.pkiImpl.ReadAcme(ctx, commonName)
		return err
	}

	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return changed, fmt.Errorf("error reading certificate %w", err)
	}

	log.Info().Msg("Certificate successfully read from Vault kv2")
	// Update metrics for the just received cert
	x509Cert, err := pkg.ParseCertPem(cert.Certificate)
	if err != nil {
		log.Error().Msgf("Could not parse certificate data: %v", err)
	} else {
		log.Info().Msgf("Read certificate valid until %v (%s)", x509Cert.NotAfter, time.Until(x509Cert.NotAfter).Round(time.Second))
	}

	if !changed {
		changed = !bytes.Equal(certData.Raw, x509Cert.Raw)
	}

	err = format.WriteCert(cert)
	if err != nil {
		return changed, fmt.Errorf("could not write bundle to backend: %v", err)
	}

	return changed, nil
}

func (p *PkiService) shouldIssue(format IssueStorage) (bool, error) {
	cert, err := format.ReadCert()
	if err != nil || cert == nil {
		if errors.Is(err, ErrNoCertFound) {
			log.Info().Msg("No existing certificate found")
			return true, nil
		} else {
			log.Warn().Msgf("Could not read certificate: %v", err)
			return true, err
		}
	}

	if !pkg.IsCertExpired(*cert) {
		if err := p.Verify(cert); err != nil {
			return true, fmt.Errorf("cert exists but can not be verified against ca: %w", err)
		}
	}

	log.Info().Msgf("Certificate %s successfully parsed", pkg.FormatSerial(cert.SerialNumber))
	return p.strategy.Renew(cert)
}

func (p *PkiService) Issue(ctx context.Context, format IssueStorage, args pkg.IssueArgs) (IssueOutcome, error) {
	shouldIssue, err := p.shouldIssue(format)
	if err == nil && !shouldIssue {
		log.Info().Msg("Cert exists and does not need a renewal")
		return NotNeeded, nil
	} else {
		log.Info().Err(err).Msg("Going to renew certificate")
	}

	var cert *pkg.CertData
	op := func() error {
		var err error
		cert, err = p.pkiImpl.Issue(ctx, args)
		return err
	}
	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	log.Info().Msg("Issuing new certificate")
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return Error, fmt.Errorf("error issuing certificate %v", err)
	}
	log.Info().Msg("New certificate successfully issued")

	// Update metrics for the just received blob
	x509Cert, err := pkg.ParseCertPem(cert.Certificate)
	if err != nil {
		log.Error().Msgf("Could not parse certificate data: %v", err)
	} else {
		log.Info().Msgf("New certificate valid until %v (%s)", x509Cert.NotAfter, time.Until(x509Cert.NotAfter).Round(time.Second))
	}

	err = format.WriteCert(cert)
	if err != nil {
		return Error, fmt.Errorf("could not write bundle to backend: %v", err)
	}

	return Issued, nil
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

	log.Info().Msgf("Received CA with serial %s", pkg.FormatSerial(ca.SerialNumber))
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

	log.Info().Msg("Trying to sign certificate")
	var resp *pkg.Signature
	op := func() error {
		var err error
		resp, err = p.pkiImpl.Sign(ctx, string(csr), args)
		return err
	}
	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
	if err := backoff.Retry(op, backoffImpl); err != nil {
		return fmt.Errorf("error signing CSR: %v", err)
	}
	log.Info().Msgf("CSR has been successfully signed using serial %s", resp.Serial)

	// Update metrics for the just received blob
	x509Cert, err := pkg.ParseCertPem(resp.Certificate)
	if err != nil {
		log.Error().Msgf("Could not parse certificate data: %v", err)
	} else {
		log.Info().Msgf("New certificate valid until %v (%s)", x509Cert.NotAfter, time.Until(x509Cert.NotAfter).Round(time.Second))
	}

	err = sink.WriteSignature(resp)
	if err != nil {
		return fmt.Errorf("could not write certificate file to backend: %v", err)
	}
	return nil
}

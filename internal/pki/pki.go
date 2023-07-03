package pki

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"
)

type IssueOutcome int

const (
	Issued    = 0
	NotNeeded = 1
	Error     = 2
)

type Pki interface {
	// Issue issues a new certificate from the PKI
	Issue(opts *conf.Config) (*CertData, error)

	// Sign signs a CSR
	Sign(csr string, opts *conf.Config) (*Signature, error)

	// Revoke revokes a certificate by its serial number
	Revoke(serial string) error

	// ReadAcme reads a previously acquired letsencrypt certificate from Vault
	ReadAcme(commonName string, config *conf.Config) (*CertData, error)

	// Tidy cleans up the PKI blob storage of dangling certificates
	Tidy() error

	// Cleanup cleans up the used resources of the client is not related to PKI operations
	Cleanup() error

	// FetchCa returns the CA for the configured mount
	FetchCa(binary bool) ([]byte, error)

	// FetchCaChain returns the whole CA chain for the configured mount
	FetchCaChain() ([]byte, error)

	// FetchCrl returns the CRL of the configured mount
	FetchCrl(binary bool) ([]byte, error)
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

type PkiCli struct {
	pkiImpl  Pki
	strategy issue_strategies.IssueStrategy
}

func NewPki(pki Pki, strategy issue_strategies.IssueStrategy) (*PkiCli, error) {
	if pki == nil {
		return nil, errors.New("empty pki impl provided")
	}
	if strategy == nil {
		strategy = &issue_strategies.StaticRenewal{Decision: true}
	}

	return &PkiCli{
		pkiImpl:  pki,
		strategy: strategy,
	}, nil
}

func updateCertificateMetrics(cert *x509.Certificate) {
	if cert == nil {
		return
	}

	secondsTotal := cert.NotAfter.Sub(cert.NotBefore).Seconds()
	internal.MetricCertLifetimeTotal.WithLabelValues(cert.Subject.CommonName).Set(secondsTotal)
	secondsUntilExpiration := time.Until(cert.NotAfter).Seconds()

	percentage := math.Max(0, secondsUntilExpiration*100./secondsTotal)

	internal.MetricCertExpiry.WithLabelValues(cert.Subject.CommonName).Set(float64(cert.NotAfter.UnixMilli()))
	internal.MetricCertLifetimePercent.WithLabelValues(cert.Subject.CommonName).Set(percentage)
}

func (p *PkiCli) Revoke(serial string) error {
	log.Info().Msgf("Attempting to revoke certificate %s", serial)
	err := p.pkiImpl.Revoke(serial)
	if err != nil {
		return fmt.Errorf("could not revoke certificate: %v", err)
	}

	log.Info().Msgf("Revoking certificate successful")
	return nil
}

func (p *PkiCli) Tidy() error {
	err := p.pkiImpl.Tidy()
	if err != nil {
		return fmt.Errorf("could not tidy certificate storage: %v", err)
	}

	log.Info().Msgf("Tidy blob storage scheduled")
	return nil
}

func (p *PkiCli) cleanup() {
	log.Info().Msg("Cleaning up the backend...")
	err := p.pkiImpl.Cleanup()
	if err != nil {
		log.Error().Msgf("Cleanup of the backend failed: %v", err)
	}
}

func (p *PkiCli) ReadAcme(format IssueSink, opts *conf.Config) (bool, error) {
	var changed bool
	certData, err := format.ReadCert()
	if err != nil || certData == nil {
		log.Info().Msg("No existing local certdata available")
		changed = true
	}

	log.Info().Msgf("Trying to read certificate for domain '%s'", opts.CommonName)
	cert, err := p.pkiImpl.ReadAcme(opts.CommonName, opts)
	if err != nil {
		return changed, fmt.Errorf("error issuing certificate %w", err)
	}
	log.Info().Msg("Certificate successfully read from Vault kv2")

	// Update metrics for the just received cert
	x509Cert, err := pkg.ParseCertPem(cert.Certificate)
	if err != nil {
		internal.MetricCertParseErrors.WithLabelValues(opts.CommonName).Set(1)
		log.Error().Msgf("Could not parse certificate data: %v", err)
	} else {
		log.Info().Msgf("Read certificate valid until %v (%s)", x509Cert.NotAfter, time.Until(x509Cert.NotAfter).Round(time.Second))
		updateCertificateMetrics(x509Cert)
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

func (p *PkiCli) shouldIssue(format IssueSink) (bool, error) {
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

	if err = p.Verify(cert); err != nil {
		return true, fmt.Errorf("cert exists but can not be verified against ca: %w", err)
	}

	log.Info().Msgf("Certificate %s successfully parsed", pkg.FormatSerial(cert.SerialNumber))
	updateCertificateMetrics(cert)
	return p.strategy.Renew(cert)
}

func (p *PkiCli) Issue(format IssueSink, opts *conf.Config) (IssueOutcome, error) {
	defer p.cleanup()

	shouldIssue, err := p.shouldIssue(format)
	if err == nil && !shouldIssue {
		log.Info().Msg("Cert exists and does not need a renewal")
		return NotNeeded, nil
	} else {
		log.Warn().Err(err).Msg("Going to renew certificate")
	}

	log.Info().Msg("Issuing new certificate")
	cert, err := p.pkiImpl.Issue(opts)
	if err != nil {
		return Error, fmt.Errorf("error issuing certificate %v", err)
	}
	log.Info().Msg("New certificate successfully issued")

	// Update metrics for the just received blob
	x509Cert, err := pkg.ParseCertPem(cert.Certificate)
	if err != nil {
		internal.MetricCertParseErrors.WithLabelValues(opts.CommonName).Set(1)
		log.Error().Msgf("Could not parse certificate data: %v", err)
	} else {
		log.Info().Msgf("New certificate valid until %v (%s)", x509Cert.NotAfter, time.Until(x509Cert.NotAfter).Round(time.Second))
		updateCertificateMetrics(x509Cert)
	}

	err = format.WriteCert(cert)
	if err != nil {
		return Error, fmt.Errorf("could not write bundle to backend: %v", err)
	}

	return Issued, nil
}

func (p *PkiCli) Verify(cert *x509.Certificate) error {
	caData, err := p.pkiImpl.FetchCaChain()
	if err != nil {
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

func (p *PkiCli) Sign(sink CsrSink, opts *conf.Config) error {
	defer p.cleanup()

	csr, err := sink.ReadCsr()
	if err != nil {
		return err
	}

	log.Info().Msg("Trying to sign certificate")
	resp, err := p.pkiImpl.Sign(string(csr), opts)
	if err != nil {
		return fmt.Errorf("error signing CSR: %v", err)
	}
	log.Info().Msgf("CSR has been successfully signed using serial %s", resp.Serial)

	// Update metrics for the just received blob
	x509Cert, err := pkg.ParseCertPem(resp.Certificate)
	if err != nil {
		internal.MetricCertParseErrors.WithLabelValues(opts.CommonName).Set(1)
		log.Error().Msgf("Could not parse certificate data: %v", err)
	} else {
		log.Info().Msgf("New certificate valid until %v (%s)", x509Cert.NotAfter, time.Until(x509Cert.NotAfter).Round(time.Second))
		updateCertificateMetrics(x509Cert)
	}

	err = sink.WriteSignature(resp)
	if err != nil {
		return fmt.Errorf("could not write certificate file to backend: %v", err)
	}
	return nil
}

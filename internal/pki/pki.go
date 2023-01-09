package pki

import (
	"bytes"
	"crypto/x509"
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

	ReadAcme(commonName string, config *conf.Config) (*CertData, error)

	// Tidy cleans up the PKI blob storage of dangling certificates
	Tidy() error

	// Cleanup cleans up the used resources of the client is not related to PKI operations
	Cleanup() error
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

func shouldIssueNewCertificate(x509Cert *x509.Certificate, strategy issue_strategies.IssueStrategy) (bool, error) {
	if x509Cert == nil {
		return true, errors.New("empty cert provided")
	}

	log.Info().Msgf("Certificate %s successfully parsed", pkg.FormatSerial(x509Cert.SerialNumber))
	updateCertificateMetrics(x509Cert)
	return strategy.Renew(x509Cert)
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
		changed = bytes.Compare(certData.Raw, x509Cert.Raw) != 0
	}

	err = format.WriteCert(cert)
	if err != nil {
		return changed, fmt.Errorf("could not write bundle to backend: %v", err)
	}

	return changed, nil
}

func (p *PkiCli) Issue(format IssueSink, opts *conf.Config) (IssueOutcome, error) {
	defer p.cleanup()
	certData, err := format.ReadCert()
	if err == nil && certData != nil {
		renew, err := shouldIssueNewCertificate(certData, p.strategy)
		if err == nil && !renew {
			log.Info().Msg("Not renewing certificate: certificate does not need renewal, yet")
			return NotNeeded, nil
		}
		if err != nil {
			log.Error().Msgf("Got error while deciding whether to renew certificate, proceeding to renew: %v", err)
		}
	}
	log.Info().Msgf("Could not read certificate: %v", err)
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

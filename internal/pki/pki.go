package pki

import (
	"crypto/x509"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"
)

// KeyPod is a simple wrapper around a key (which is just a byte stream itself). This way, we decouple
// the implementation (file-based, memory, network, ..) and make it easily swap- and testable.
type KeyPod interface {
	Read() ([]byte, error)
	CanRead() error
	Write(string) error
	CanWrite() error
}

type IssueOutcome int

const (
	Issued    = 0
	NotNeeded = 1
	Error     = 2
)

type Pki interface {
	// Issue issues a new certificate from the PKI
	Issue(opts conf.IssueArguments) (*IssuedCert, error)

	// Sign signs a CSR
	Sign(csr KeyPod, opts conf.SignArguments) (*Signature, error)

	// Revoke revokes a certificate by its serial number
	Revoke(serial string) error

	// Tidy cleans up the PKI cert storage of dangling certificates
	Tidy() error

	// Cleanup cleans up the used resources of the client is not related to PKI operations
	Cleanup() error
}

type IssuedCert struct {
	PrivateKey  []byte
	Certificate []byte
	CaChain     []byte
}

type Signature struct {
	Certificate []byte
	CaChain     []byte
	Serial      string
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
	internal.MetricCertLifetimeTotal.Set(secondsTotal)
	secondsUntilExpiration := time.Until(cert.NotAfter).Seconds()

	percentage := math.Max(0, secondsUntilExpiration*100./secondsTotal)

	internal.MetricCertExpiry.Set(float64(cert.NotAfter.UnixMilli()))
	internal.MetricCertLifetimePercent.Set(percentage)
}

func shouldIssueNewCertificate(certFile KeyPod, strategy issue_strategies.IssueStrategy) (bool, error) {
	log.Info().Msg("A certificate already exists, trying to parse it")
	cert, err := parseCert(certFile)
	if err != nil {
		internal.MetricCertParseErrors.Set(1)
		return true, fmt.Errorf("could not parse existing certificate data: %v", err)
	}

	log.Info().Msgf("Certificate %s successfully parsed", pkg.FormatSerial(cert.SerialNumber))
	updateCertificateMetrics(cert)
	return strategy.Renew(cert)
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

	log.Info().Msgf("Tidy cert storage scheduled")
	return nil
}

func (p *PkiCli) cleanup() {
	log.Info().Msg("Cleaning up the backend...")
	err := p.pkiImpl.Cleanup()
	if err != nil {
		log.Error().Msgf("Cleanup of the backend failed: %v", err)
	}
}

func (p *PkiCli) Issue(certFile, privateKeyFile KeyPod, opts conf.IssueArguments) (IssueOutcome, error) {
	defer p.cleanup()
	if certFile.CanRead() == nil {
		renew, err := shouldIssueNewCertificate(certFile, p.strategy)
		if err == nil && !renew {
			log.Info().Msg("Not renewing certifcate: certificate does not need renewal, yet")
			return NotNeeded, nil
		}
		if err != nil {
			log.Error().Msgf("Got error while deciding whether to renew certifcate, proceeding to renew: %v", err)
		}
	}

	log.Info().Msg("Issuing new certificate")
	cert, err := p.pkiImpl.Issue(opts)
	if err != nil {
		return Error, fmt.Errorf("error issuing certificate %v", err)
	}
	log.Info().Msg("New certificate successfully issued")

	// Update metrics for the just received cert
	x509Cert, err := pkg.ParseCertPem(cert.Certificate)
	if err != nil {
		internal.MetricCertParseErrors.Set(1)
		log.Error().Msgf("Could not parse certificate data: %v", err)
	} else {
		log.Info().Msgf("New certificate valid until %v (%s)", x509Cert.NotAfter, time.Until(x509Cert.NotAfter).Round(time.Second))
		updateCertificateMetrics(x509Cert)
	}

	err = certFile.Write(string(cert.Certificate))
	if err != nil {
		return Error, fmt.Errorf("could not write certificate file to backend: %v", err)
	}

	err = privateKeyFile.Write(string(cert.PrivateKey))
	if err != nil {
		return Error, fmt.Errorf("could not write private key to backend: %v", err)
	}

	return Issued, nil
}

func (p *PkiCli) Sign(certFile, csrFile KeyPod, opts conf.SignArguments) error {
	defer p.cleanup()

	log.Info().Msg("Issuing new certificate")
	resp, err := p.pkiImpl.Sign(csrFile, opts)
	if err != nil {
		return fmt.Errorf("error signing CSR: %v", err)
	}
	log.Info().Msgf("CSR has been successfully signed using serial %s", resp.Serial)

	// Update metrics for the just received cert
	x509Cert, err := pkg.ParseCertPem(resp.Certificate)
	if err != nil {
		internal.MetricCertParseErrors.Set(1)
		log.Error().Msgf("Could not parse certificate data: %v", err)
	} else {
		log.Info().Msgf("New certificate valid until %v (%s)", x509Cert.NotAfter, time.Until(x509Cert.NotAfter).Round(time.Second))
		updateCertificateMetrics(x509Cert)
	}

	err = certFile.Write(string(resp.Certificate))
	if err != nil {
		return fmt.Errorf("could not write certificate file to backend: %v", err)
	}

	return nil
}

func parseCert(certFile KeyPod) (*x509.Certificate, error) {
	content, err := certFile.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read certificate data: %v", err)
	}

	return pkg.ParseCertPem(content)
}

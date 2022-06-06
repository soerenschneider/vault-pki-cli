package conf

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/pods"
)

type Backend struct {
	CertificateFile string
	PrivateKeyFile  string
	CaFile          string
	FileOwner       string
	FileGroup       string
}

type IssueArguments struct {
	CommonName          string
	Ttl                 string
	IpSans              []string
	AltNames            []string
	ForceNewCertificate bool

	Backends []Backend

	PostIssueHooks []string

	CertificateLifetimeThresholdPercentage float64

	YubikeyPin  string
	YubikeySlot uint32

	MetricsFile string
}

func (c *IssueArguments) UsesYubikey() bool {
	return c.Backends == nil || len(c.Backends) == 0 || len(c.Backends[0].CertificateFile) == 0
}

func (c *Backend) Validate() (errs []error) {
	ownerDefined := len(c.FileOwner) > 0
	groupDefined := len(c.FileGroup) > 0
	if !ownerDefined && groupDefined {
		errs = append(errs, fmt.Errorf("only '%s' defined but not '%s'", FLAG_FILE_GROUP, FLAG_FILE_OWNER))
	}
	if ownerDefined && !groupDefined {
		errs = append(errs, fmt.Errorf("only '%s' defined but not '%s'", FLAG_FILE_OWNER, FLAG_FILE_GROUP))
	}

	emptyPrivateKeyFile := len(c.PrivateKeyFile) == 0
	if emptyPrivateKeyFile {
		errs = append(errs, fmt.Errorf("must provide private key file '%s'", FLAG_ISSUE_PRIVATE_KEY_FILE))
	}

	return
}

func (c *IssueArguments) Validate() []error {
	errs := make([]error, 0)

	if len(c.CommonName) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_COMMON_NAME))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		errs = append(errs, fmt.Errorf("'%s' must be [5, 90]", FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE))
	}

	for _, backend := range c.Backends {
		validationErrs := backend.Validate()
		errs = append(errs, validationErrs...)
	}

	emptyYubikeySlot := c.YubikeySlot == FLAG_ISSUE_YUBIKEY_SLOT_DEFAULT
	if len(c.Backends) == 0 && emptyYubikeySlot {
		errs = append(errs, fmt.Errorf("must either provide '%s' or both '%s' and '%s'", FLAG_ISSUE_YUBIKEY_SLOT, FLAG_CERTIFICATE_FILE, FLAG_ISSUE_PRIVATE_KEY_FILE))
	}

	if len(c.Backends) > 0 && !emptyYubikeySlot {
		errs = append(errs, errors.New("can't provide yubi key slot AND file-based backends"))
	}

	if !emptyYubikeySlot {
		err := pods.ValidateSlot(c.YubikeySlot)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid yubikey slot '%d': %v", c.YubikeySlot, err))
		}
	}

	return errs
}

func (c *IssueArguments) PrintConfig() {
	log.Info().Msg("------------- Printing issue cmd values -------------")
	if c.YubikeySlot != FLAG_ISSUE_YUBIKEY_SLOT_DEFAULT {
		log.Info().Msgf("%s=%x", FLAG_ISSUE_YUBIKEY_SLOT, c.YubikeySlot)
	}

	if len(c.YubikeyPin) > 0 {
		log.Info().Msgf("%s=%s", FLAG_ISSUE_YUBIKEY_PIN, "*** (Redacted)")
	}

	for n, backend := range c.Backends {
		if len(backend.CaFile) > 0 {
			log.Info().Msgf("%s[%d]=%s", FLAG_CA_FILE, n, backend.CaFile)
		}

		if len(backend.CertificateFile) > 0 {
			log.Info().Msgf("%s[%d]=%s", FLAG_CERTIFICATE_FILE, n, backend.CertificateFile)
		}

		if len(backend.PrivateKeyFile) > 0 {
			log.Info().Msgf("%s[%d]=%s", FLAG_ISSUE_PRIVATE_KEY_FILE, n, backend.PrivateKeyFile)
		}

		if len(backend.FileOwner) > 0 {
			log.Info().Msgf("%s[%d]=%s", FLAG_FILE_OWNER, n, backend.FileOwner)
		}

		if len(backend.FileGroup) > 0 {
			log.Info().Msgf("%s[%d]=%s", FLAG_FILE_GROUP, n, backend.FileGroup)
		}
	}

	log.Info().Msgf("%s=%s", FLAG_ISSUE_TTL, c.Ttl)
	log.Info().Msgf("%s=%s", FLAG_ISSUE_COMMON_NAME, c.CommonName)
	log.Info().Msgf("%s=%s", FLAG_ISSUE_METRICS_FILE, c.MetricsFile)
	log.Info().Msgf("%s=%t", FLAG_ISSUE_FORCE_NEW_CERTIFICATE, c.ForceNewCertificate)
	log.Info().Msgf("%s=%f", FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, c.CertificateLifetimeThresholdPercentage)
	if len(c.IpSans) > 0 {
		log.Info().Msgf("%s=%v", FLAG_ISSUE_IP_SANS, c.IpSans)
	}
	if len(c.AltNames) > 0 {
		log.Info().Msgf("%s=%v", FLAG_ISSUE_ALT_NAMES, c.AltNames)
	}
	for n, hook := range c.PostIssueHooks {
		log.Info().Msgf("%s[%d]='%s'", FLAG_ISSUE_HOOKS, n, hook)
	}
	log.Info().Msgf("------------- Finished printing config values -------------")
}

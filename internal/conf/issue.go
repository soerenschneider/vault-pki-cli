package conf

import (
	"fmt"
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/pods"
)

type IssueArguments struct {
	CommonName          string
	Ttl                 string
	IpSans              []string
	AltNames            []string
	ForceNewCertificate bool

	CertificateFile string
	PrivateKeyFile  string
	ChainFile       string
	FileOwner       string
	FileGroup       string

	CertificateLifetimeThresholdPercentage float64

	YubikeyPin  string
	YubikeySlot uint32

	MetricsFile string
}

func (c *IssueArguments) UsesYubikey() bool {
	if len(c.CertificateFile) == 0 {
		return true
	}

	return false
}

func (c *IssueArguments) Validate() []error {
	errs := make([]error, 0)

	if len(c.CommonName) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_COMMON_NAME))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		errs = append(errs, fmt.Errorf("'%s' must be [5, 90]", FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE))
	}

	ownerDefined := len(c.FileOwner) > 0
	groupDefined := len(c.FileGroup) > 0
	if !ownerDefined && groupDefined {
		errs = append(errs, fmt.Errorf("only '%s' defined but not '%s'", FLAG_FILE_GROUP, FLAG_FILE_OWNER))
	}
	if ownerDefined && !groupDefined {
		errs = append(errs, fmt.Errorf("only '%s' defined but not '%s'", FLAG_FILE_OWNER, FLAG_FILE_GROUP))
	}

	emptyYubikeySlot := c.YubikeySlot == FLAG_ISSUE_YUBIKEY_SLOT_DEFAULT
	emptyCertificateFile := len(c.CertificateFile) == 0
	emptyPrivateKeyFile := len(c.PrivateKeyFile) == 0
	if (emptyPrivateKeyFile || emptyCertificateFile) && emptyYubikeySlot {
		errs = append(errs, fmt.Errorf("must either provide '%s' or both '%s' and '%s'", FLAG_ISSUE_YUBIKEY_SLOT, FLAG_CERTIFICATE_FILE, FLAG_ISSUE_PRIVATE_KEY_FILE))
	}

	if !emptyCertificateFile && !emptyPrivateKeyFile && !emptyYubikeySlot {
		errs = append(errs, fmt.Errorf("can't provide yubi key slot and both '%s' and '%s'", FLAG_CERTIFICATE_FILE, FLAG_ISSUE_PRIVATE_KEY_FILE))
	}

	if !emptyYubikeySlot {
		_, err := pods.TranslateSlot(c.YubikeySlot)
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

	if len(c.CertificateFile) > 0 {
		log.Info().Msgf("%s=%s", FLAG_CERTIFICATE_FILE, c.CertificateFile)
	}

	if len(c.PrivateKeyFile) > 0 {
		log.Info().Msgf("%s=%s", FLAG_ISSUE_PRIVATE_KEY_FILE, c.PrivateKeyFile)
	}

	if len(c.FileOwner) > 0 {
		log.Info().Msgf("%s=%s", FLAG_FILE_OWNER, c.FileOwner)
	}

	if len(c.FileGroup) > 0 {
		log.Info().Msgf("%s=%s", FLAG_FILE_GROUP, c.FileGroup)
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
	log.Info().Msgf("------------- Finished printing config values -------------")
}

package conf

import (
	"fmt"
	log "github.com/rs/zerolog/log"
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

	MetricsFile string
}

func (c *IssueArguments) Validate() []error {
	errs := make([]error, 0)
	if len(c.CertificateFile) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_CERTIFICATE_FILE))
	}

	if len(c.PrivateKeyFile) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_PRIVATE_KEY_FILE))
	}

	if len(c.CommonName) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_COMMON_NAME))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		errs = append(errs, fmt.Errorf("'%s' must be [5, 90]", FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE))
	}

	return errs
}

func (c *IssueArguments) PrintConfig() {
	log.Info().Msg("------------- Printing issue cmd values -------------")
	log.Info().Msgf("%s=%s", FLAG_CERTIFICATE_FILE, c.CertificateFile)
	log.Info().Msgf("%s=%s", FLAG_ISSUE_PRIVATE_KEY_FILE, c.PrivateKeyFile)
	log.Info().Msgf("%s=%s", FLAG_FILE_OWNER, c.FileOwner)
	log.Info().Msgf("%s=%s", FLAG_FILE_GROUP, c.FileGroup)
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

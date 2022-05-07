package conf

import (
	"fmt"
	log "github.com/rs/zerolog/log"
)

type SignArguments struct {
	CommonName string
	Ttl        string
	IpSans     []string
	AltNames   []string

	CsrFile         string
	CertificateFile string
	ChainFile       string
	FileOwner       string
	FileGroup       string

	MetricsFile string
}

func (c *SignArguments) Validate() []error {
	errs := make([]error, 0)
	if len(c.CertificateFile) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_CERTIFICATE_FILE))
	}

	if len(c.CsrFile) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_CSR_FILE))
	}

	if len(c.CommonName) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_COMMON_NAME))
	}

	return errs
}

func (c *SignArguments) PrintConfig() {
	log.Info().Msg("------------- Printing sign cmd values -------------")
	log.Info().Msgf("%s=%s", FLAG_CSR_FILE, c.CsrFile)
	log.Info().Msgf("%s=%s", FLAG_CERTIFICATE_FILE, c.CertificateFile)
	log.Info().Msgf("%s=%s", FLAG_FILE_OWNER, c.FileOwner)
	log.Info().Msgf("%s=%s", FLAG_FILE_GROUP, c.FileGroup)
	log.Info().Msgf("%s=%s", FLAG_ISSUE_TTL, c.Ttl)
	log.Info().Msgf("%s=%s", FLAG_ISSUE_COMMON_NAME, c.CommonName)
	log.Info().Msgf("%s=%s", FLAG_ISSUE_METRICS_FILE, c.MetricsFile)
	if len(c.IpSans) > 0 {
		log.Info().Msgf("%s=%v", FLAG_ISSUE_IP_SANS, c.IpSans)
	}
	if len(c.AltNames) > 0 {
		log.Info().Msgf("%s=%v", FLAG_ISSUE_ALT_NAMES, c.AltNames)
	}
	log.Info().Msgf("------------- Finished printing config values -------------")
}

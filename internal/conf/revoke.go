package conf

import (
	"fmt"
	log "github.com/rs/zerolog/log"
)

type RevokeArguments struct {
	CertificateFile string
}

func (c *RevokeArguments) Validate() []error {
	errs := make([]error, 0)
	if len(c.CertificateFile) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_CERTIFICATE_FILE))
	}

	return errs
}

func (c *RevokeArguments) PrintConfig() {
	log.Info().Msgf("------------- Printing revoke cmd values --------------")
	log.Info().Msgf("%s=%s", FLAG_CERTIFICATE_FILE, c.CertificateFile)
	log.Info().Msgf("------------- Finished printing config values -------------")
}

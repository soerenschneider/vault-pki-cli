package conf

import (
	"github.com/rs/zerolog/log"
)

type FetchArguments struct {
	CertificateFile string
	DerEncoded      bool
}

func (c *FetchArguments) Validate() []error {
	return nil
}

func (c *FetchArguments) PrintConfig() {
	log.Info().Msg("------------- Printing issue cmd values -------------")
	log.Info().Msgf("%s=%s", FLAG_CERTIFICATE_FILE, c.CertificateFile)
	log.Info().Msgf("%s=%s", FLAG_DER_ENCODED, c.DerEncoded)
	log.Info().Msgf("------------- Finished printing config values -------------")
}
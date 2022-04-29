package conf

import (
	"github.com/rs/zerolog/log"
)

type FetchArguments struct {
	OutputFile string
	DerEncoded bool
}

func (c *FetchArguments) Validate() []error {
	return nil
}

func (c *FetchArguments) PrintConfig() {
	log.Info().Msg("------------- Printing fetch cmd values -------------")
	if len(c.OutputFile) > 0 {
		log.Info().Msgf("%s=%s", FLAG_OUTPUT_FILE, c.OutputFile)
	}
	log.Info().Msgf("%s=%t", FLAG_DER_ENCODED, c.DerEncoded)
	log.Info().Msgf("------------- Finished printing config values -------------")
}

package main

import (
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func readCaCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca",
		Short: "Read pki ca cert from vault",
		Run:   readCaEntryPoint,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// See https://github.com/spf13/viper/issues/233#issuecomment-386791444
			viper.BindPFlag(conf.FLAG_OUTPUT_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_OUTPUT_FILE))
			viper.BindPFlag(conf.FLAG_DER_ENCODED, cmd.PersistentFlags().Lookup(conf.FLAG_DER_ENCODED))

			return nil
		},
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "Write ca certificate to this output file")
	getCaCmd.PersistentFlags().BoolP(conf.FLAG_DER_ENCODED, "d", false, "Use DER encoding")
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return getCaCmd
}

func readCaEntryPoint(ccmd *cobra.Command, args []string) {
	PrintVersionInfo()
	configFile := viper.GetViper().GetString(conf.FLAG_CONFIG_FILE)
	var config *conf.Config
	if len(configFile) > 0 {
		var err error
		config, err = readConfig(configFile)
		if err != nil {
			log.Fatal().Msgf("Could not load desired config file: %s: %v", configFile, err)
		}
		log.Info().Msgf("Read config from file %s", viper.ConfigFileUsed())
	}

	if len(config.VaultAddress) == 0 {
		log.Error().Msg("Missing vault address, quitting")
		os.Exit(1)
	}

	if len(config.VaultMountPki) == 0 {
		log.Error().Msg("Missing vault pki mount, quitting")
		os.Exit(1)
	}

	config.FetchArguments.PrintConfig()
	certData, err := vault.FetchCert(config.VaultAddress, config.VaultMountPki, config.DerEncoded)
	if err != nil {
		os.Exit(1)
	}

	handleFetchedData(certData, *config)
}

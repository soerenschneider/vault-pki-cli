package main

import (
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func readCaChainCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca-chain",
		Short: "Read pki ca cert chain from vault",
		Run:   fetchCaChainEntryPoint,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// See https://github.com/spf13/viper/issues/233#issuecomment-386791444
			viper.BindPFlag(conf.FLAG_OUTPUT_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_OUTPUT_FILE))
			return initializeConfig(cmd)
		},
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "Write ca certificate chain to this file")
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return getCaCmd
}

func fetchCaChainEntryPoint(ccmd *cobra.Command, args []string) {
	initializeConfig(ccmd)

	PrintVersionInfo()
	configFile := viper.GetViper().GetString(conf.FLAG_CONFIG_FILE)
	if len(configFile) > 0 {
		err := readConfig(configFile)
		if err != nil {
			log.Fatal().Msgf("Could not load desired config file: %s: %v", configFile, err)
		}
		log.Info().Msgf("Read config from file %s", viper.ConfigFileUsed())
	}

	config := NewConfigFromViper()
	if len(config.VaultAddress) == 0 {
		log.Error().Msg("Missing vault address, quitting")
		os.Exit(1)
	}

	if len(config.VaultMountPki) == 0 {
		log.Error().Msg("Missing vault pki mount, quitting")
		os.Exit(1)
	}

	config.FetchArguments.PrintConfig()
	certData, err := vault.FetchCertChain(config.VaultAddress, config.VaultMountPki)
	if err != nil {
		os.Exit(1)
	}

	handleFetchedData(certData, config)
}

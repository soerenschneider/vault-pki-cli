package main

import (
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
	"os"
)

func readCaChainCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca-chain",
		Short: "ReadCert pki ca cert chain from vault",
		Run:   fetchCaChainEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteCert ca certificate chain to this file")
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return getCaCmd
}

func fetchCaChainEntryPoint(ccmd *cobra.Command, args []string) {

	PrintVersionInfo()
	config, err := config()
	if err != nil {
		log.Fatal().Err(err)
	}

	if len(config.VaultAddress) == 0 {
		log.Error().Msg("Missing vault address, quitting")
		os.Exit(1)
	}

	if len(config.VaultMountPki) == 0 {
		log.Error().Msg("Missing vault pki mount, quitting")
		os.Exit(1)
	}

	certData, err := vault.FetchCertChain(config.VaultAddress, config.VaultMountPki)
	if err != nil {
		os.Exit(1)
	}

	handleFetchedData(certData, *config)
}

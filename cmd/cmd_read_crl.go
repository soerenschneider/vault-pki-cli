package main

import (
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
	"os"
)

func readCrlCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-crl",
		Short: "ReadCert pki crl from vault",
		Run:   readCrlEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteCert CRL to this file")
	getCaCmd.PersistentFlags().BoolP(conf.FLAG_DER_ENCODED, "d", false, "Use DER encoding")
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return getCaCmd
}

func readCrlEntryPoint(ccmd *cobra.Command, args []string) {
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

	crlData, err := vault.FetchCrl(config.VaultAddress, config.VaultMountPki, config.DerEncoded)
	if err != nil {
		os.Exit(1)
	}

	handleFetchedData(crlData, *config)
}

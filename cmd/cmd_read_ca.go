package main

import (
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
	"os"
)

func readCaCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca",
		Short: "ReadCert pki ca cert from vault",
		Run:   readCaEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteCert ca certificate to this output file")
	getCaCmd.PersistentFlags().BoolP(conf.FLAG_DER_ENCODED, "d", false, "Use DER encoding")
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return getCaCmd
}

func readCaEntryPoint(ccmd *cobra.Command, args []string) {
	PrintVersionInfo()
	config, err := config()
	if err != nil {
		log.Fatal().Err(err)
	}

	certData, err := vault.FetchCert(config.VaultAddress, config.VaultMountPki, config.DerEncoded)
	if err != nil {
		os.Exit(1)
	}

	handleFetchedData(certData, *config)
}

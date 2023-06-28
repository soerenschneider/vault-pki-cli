package main

import (
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
)

func readCaCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca",
		Short: "ReadCert pki ca cert from vault",
		Run:   readCaEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteSignature ca certificate to this output file")
	getCaCmd.PersistentFlags().BoolP(conf.FLAG_DER_ENCODED, "d", false, "Use DER encoding")
	if err := getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE); err != nil {
		log.Fatal().Err(err).Msg("could not mark flag required")
	}

	return getCaCmd
}

func readCaEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()
	config, err := config()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get config")
	}

	if len(config.VaultAddress) == 0 {
		log.Fatal().Msg("missing vault address, quitting")
	}

	if len(config.VaultMountPki) == 0 {
		log.Fatal().Msg("missing vault pki mount, quitting")
	}

	storage.InitBuilder(config)
	certData, err := vault.FetchCert(config.VaultAddress, config.VaultMountPki, config.DerEncoded)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not read cert data from vault: %v", err)
	}

	sink, err := sink.CaSinkFromConfig(config.StorageConfig)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not build ca sink from config: %v", err)
	}

	if err = sink.WriteCa(certData); err != nil {
		log.Fatal().Err(err).Msgf("could not write ca: %v", err)
	}
}

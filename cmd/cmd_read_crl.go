package main

import (
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
)

func readCrlCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-crl",
		Short: "ReadCert pki crl from vault",
		Run:   readCrlEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteSignature CRL to this file")
	getCaCmd.PersistentFlags().BoolP(conf.FLAG_DER_ENCODED, "d", false, "Use DER encoding")
	if err := getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE); err != nil {
		log.Fatal().Err(err).Msg("could not mark flag required")
	}

	return getCaCmd
}

func readCrlEntryPoint(_ *cobra.Command, _ []string) {
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

	vaultClient, err := api.NewClient(getVaultConfig(config))
	if err != nil {
		log.Fatal().Err(err).Msgf("could not build vault client")
	}

	pkiImpl, err := vault.NewVaultPki(vaultClient, &vault.NoAuth{}, config)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not build rotation client")
	}

	storage.InitBuilder(config)
	crlData, err := pkiImpl.FetchCrl(config.DerEncoded)
	if err != nil {
		log.Fatal().Err(err).Msg("could not fetch crl")
	}

	sink, err := sink.CrlSinkFromConfig(config.StorageConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build crl sink from config")
	}

	if err = sink.WriteCrl(crlData); err != nil {
		log.Fatal().Err(err).Msg("could not write crl")
	}
}

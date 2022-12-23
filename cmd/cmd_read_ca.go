package main

import (
	"errors"
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
		RunE:  readCaEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteSignature ca certificate to this output file")
	getCaCmd.PersistentFlags().BoolP(conf.FLAG_DER_ENCODED, "d", false, "Use DER encoding")
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return getCaCmd
}

func readCaEntryPoint(ccmd *cobra.Command, args []string) error {
	PrintVersionInfo()
	config, err := config()
	if err != nil {
		log.Fatal().Err(err)
	}

	if len(config.VaultAddress) == 0 {
		return errors.New("missing vault address, quitting")
	}

	if len(config.VaultMountPki) == 0 {
		return errors.New("missing vault pki mount, quitting")
	}

	storage.InitBuilder(config)
	certData, err := vault.FetchCert(config.VaultAddress, config.VaultMountPki, config.DerEncoded)
	if err != nil {
		log.Error().Msgf("Could not read cert data from vault: %v", err)
		return err
	}

	sink, err := sink.CaSinkFromConfig(config.StorageConfig)
	if err != nil {
		log.Error().Msgf("could not build ca sink from config: %v", err)
		return err
	}

	return sink.WriteCa(certData)
}

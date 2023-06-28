package main

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
	"os"
)

func readCaChainCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca-chain",
		Short: "ReadCert pki ca cert chain from vault",
		RunE:  fetchCaChainEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteSignature ca certificate chain to this file")
	if err := getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE); err != nil {
		log.Fatal().Err(err).Msg("could not mark flag required")
	}

	return getCaCmd
}

func fetchCaChainEntryPoint(ccmd *cobra.Command, args []string) error {

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
	certData, err := vault.FetchCertChain(config.VaultAddress, config.VaultMountPki)
	if err != nil {
		os.Exit(1)
	}

	sink, err := sink.CaSinkFromConfig(config.StorageConfig)
	if err != nil {
		log.Error().Msgf("could not build ca sink from config: %v", err)
		return err
	}

	return sink.WriteCa(certData)
}

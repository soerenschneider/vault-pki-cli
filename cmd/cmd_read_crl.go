package main

import (
	"errors"
	"fmt"
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
		RunE:  readCrlEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteSignature CRL to this file")
	getCaCmd.PersistentFlags().BoolP(conf.FLAG_DER_ENCODED, "d", false, "Use DER encoding")
	if err := getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE); err != nil {
		log.Fatal().Err(err).Msg("could not mark flag required")
	}

	return getCaCmd
}

func readCrlEntryPoint(ccmd *cobra.Command, args []string) error {
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
	crlData, err := vault.FetchCrl(config.VaultAddress, config.VaultMountPki, config.DerEncoded)
	if err != nil {
		return fmt.Errorf("could not fetch crl from vault: %v", err)
	}

	sink, err := sink.CrlSinkFromConfig(config.StorageConfig)
	if err != nil {
		return fmt.Errorf("could not build crl sink from config: %v", err)
	}

	return sink.WriteCrl(crlData)
}

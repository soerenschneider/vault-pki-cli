package main

import (
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg/vault"
	"github.com/spf13/cobra"
)

func readCrlCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-crl",
		Short: "Read pki crl from vault",
		Run:   readCrlEntryPoint,
	}

	getCaCmd.Flags().Uint64(conf.FLAG_RETRIES, conf.FLAG_RETRIES_DEFAULT, "How many retries to perform for non-permanent errors")
	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteSignature CRL to this file")
	getCaCmd.PersistentFlags().BoolP(conf.FLAG_DER_ENCODED, "d", false, "Use DER encoding")
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE) // nolint:errcheck

	return getCaCmd
}

func readCrlEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()
	config, err := config()
	DieOnErr(err, "could not get config")

	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "could not build vault client")

	opts := []vault.VaultOpts{
		vault.WithPkiMount(config.VaultMountPki),
		vault.WithKv2Mount(config.VaultMountKv2),
		vault.WithAcmePrefix(config.AcmePrefix),
	}

	pkiImpl, err := vault.NewVaultPki(vaultClient.Logical(), config.VaultPkiRole, opts...)
	DieOnErr(err, "could not build crl client")

	storage.InitBuilder(config)
	crlData, err := pkiImpl.FetchCrl(config.DerEncoded)
	DieOnErr(err, "could not fetch crl")

	sink, err := storage.CrlStorageFromConfig(config.StorageConfig)
	DieOnErr(err, "could not build crl sink from config")

	err = sink.WriteCrl(crlData)
	DieOnErr(err, "could not write crl")
}

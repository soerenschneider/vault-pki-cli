package main

import (
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
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return getCaCmd
}

func readCrlEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()
	config, err := config()
	DieOnErr(err, "could not get config")

	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "could not build vault client")

	pkiImpl, err := vault.NewVaultPki(vaultClient, &vault.NoAuth{}, config)
	DieOnErr(err, "could not build rotation client")

	storage.InitBuilder(config)
	crlData, err := pkiImpl.FetchCrl(config.DerEncoded)
	DieOnErr(err, "could not fetch crl")

	sink, err := sink.CrlSinkFromConfig(config.StorageConfig)
	DieOnErr(err, "could not build crl sink from config")

	err = sink.WriteCrl(crlData)
	DieOnErr(err, "could not write crl")
}

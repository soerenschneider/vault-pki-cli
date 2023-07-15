package main

import (
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
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE) // nolint:errcheck

	return getCaCmd
}

func readCaEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()
	config, err := config()
	DieOnErr(err, "could not get config")

	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "could not build vault client")

	pkiImpl, err := vault.NewVaultPki(vaultClient, &vault.NoAuth{}, config)
	DieOnErr(err, "could not build rotation client")

	storage.InitBuilder(config)
	certData, err := pkiImpl.FetchCa(config.DerEncoded)
	DieOnErr(err, "Could not read cert data from vault")

	sink, err := sink.CaSinkFromConfig(config.StorageConfig)
	DieOnErr(err, "could not build ca sink from config")

	err = sink.WriteCa(certData)
	DieOnErr(err, "could not write ca")
}

package main

import (
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
)

func readCaChainCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca-chain",
		Short: "Read pki ca cert chain from vault",
		Run:   fetchCaChainEntryPoint,
	}

	getCaCmd.PersistentFlags().StringP(conf.FLAG_OUTPUT_FILE, "o", "", "WriteSignature ca certificate chain to this file")
	getCaCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE) // nolint:errcheck

	return getCaCmd
}

func fetchCaChainEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()
	config, err := config()
	DieOnErr(err, "can't get config")

	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "could not build vault client")

	pkiImpl, err := vault.NewVaultPki(vaultClient, &vault.NoAuth{}, config)
	DieOnErr(err, "could not build rotation client")

	storage.InitBuilder(config)
	certData, err := pkiImpl.FetchCaChain()
	DieOnErr(err, "can't fetch ca chain")

	sink, err := sink.CaSinkFromConfig(config.StorageConfig)
	DieOnErr(err, "could not build ca sink from config")

	err = sink.WriteCa(certData)
	DieOnErr(err, "could not write data")
}

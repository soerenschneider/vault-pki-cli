package main

import (
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg/vault"
	"github.com/spf13/cobra"
)

func readCaChainCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca-chain",
		Short: "Read pki ca cert chain from vault",
		Run:   fetchCaChainEntryPoint,
	}

	getCaCmd.Flags().Uint64(conf.FLAG_RETRIES, conf.FLAG_RETRIES_DEFAULT, "How many retries to perform for non-permanent errors")
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

	opts := []vault.VaultOpts{
		vault.WithPkiMount(config.VaultMountPki),
		vault.WithKv2Mount(config.VaultMountKv2),
		vault.WithAcmePrefix(config.AcmePrefix),
	}

	pkiImpl, err := vault.NewVaultPki(vaultClient.Logical(), config.VaultPkiRole, opts...)
	DieOnErr(err, "could not build rotation client")

	storage.InitBuilder(config)
	certData, err := pkiImpl.FetchCaChain()
	DieOnErr(err, "can't fetch ca chain")

	sink, err := storage.CaStorageFromConfig(config.StorageConfig)
	DieOnErr(err, "could not build ca sink from config")

	err = sink.WriteCa(certData)
	DieOnErr(err, "could not write data")
}

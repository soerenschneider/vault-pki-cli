package main

import (
	"github.com/cenkalti/backoff/v3"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg/vault"
	"github.com/spf13/cobra"
)

func readCaCmd() *cobra.Command {
	var getCaCmd = &cobra.Command{
		Use:   "read-ca",
		Short: "Read pki ca cert from vault",
		Run:   readCaEntryPoint,
	}

	getCaCmd.Flags().Uint64(conf.FLAG_RETRIES, conf.FLAG_RETRIES_DEFAULT, "How many retries to perform for non-permanent errors")
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

	opts := []vault.VaultOpts{
		vault.WithPkiMount(config.VaultMountPki),
		vault.WithKv2Mount(config.VaultMountKv2),
		vault.WithAcmePrefix(config.AcmePrefix),
	}

	pkiImpl, err := vault.NewVaultPki(vaultClient.Logical(), config.VaultPkiRole, opts...)
	DieOnErr(err, "could not build rotation client")

	storage.InitBuilder(config)

	var certData []byte
	op := func() error {
		var err error
		certData, err = pkiImpl.FetchCa(config.DerEncoded)
		return err
	}
	backoffImpl := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 15)

	err = backoff.Retry(op, backoffImpl)
	DieOnErr(err, "Could not read cert data from vault")

	sink, err := storage.CaStorageFromConfig(config.StorageConfig)
	DieOnErr(err, "could not build ca sink from config")

	err = sink.WriteCa(certData)
	DieOnErr(err, "could not write ca")
}

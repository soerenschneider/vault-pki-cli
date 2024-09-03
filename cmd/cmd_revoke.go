package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg/vault"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func getRevokeCmd() *cobra.Command {
	var revokeCmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a x509 cert",
		Run:   revokeCertEntryPoint,
	}

	revokeCmd.Flags().Uint64(conf.FLAG_RETRIES, conf.FLAG_RETRIES_DEFAULT, "How many retries to perform for non-permanent errors")
	revokeCmd.Flags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "Certificate to read serial from")

	return revokeCmd
}

func revokeCertEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()
	config, err := config()
	DieOnErr(err, "could not get config")

	err = config.Validate()
	DieOnErr(err, "invalid config")

	storage.InitBuilder(config)
	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "could not build vault client")

	authStrategy, err := buildAuthImpl(config)
	DieOnErr(err, "could not build auth strategy")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = authStrategy.Login(ctx, vaultClient)
	DieOnErr(err, "can't login to vault")

	opts := []vault.VaultOpts{
		vault.WithPkiMount(config.VaultMountPki),
		vault.WithKv2Mount(config.VaultMountKv2),
		vault.WithAcmePrefix(config.AcmePrefix),
	}

	vaultBackend, err := vault.NewVaultPki(vaultClient.Logical(), config.VaultPkiRole, opts...)
	DieOnErr(err, "could not build rotation client")

	pkiImpl, err := pki.NewPkiService(vaultBackend, nil)
	DieOnErr(err, "could not build pki impl")

	sink, err := storage.MultiKeyPairStorageFromConfig(config)
	DieOnErr(err, "could not build keypair")

	cert, err := sink.ReadCert()
	DieOnErr(err, "can not read certificate")

	if pkg.IsCertExpired(*cert) {
		log.Info().Msg("Certificate is expired, no revocation needed")
		return
	}

	serial := pkg.FormatSerial(cert.SerialNumber)
	err = pkiImpl.Revoke(ctx, serial)
	DieOnErr(err, "could not revoke cert")
}

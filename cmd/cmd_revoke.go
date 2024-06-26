package main

import (
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/spf13/cobra"
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

	vaultBackend, err := vault.NewVaultPki(vaultClient, authStrategy, config)
	DieOnErr(err, "could not build rotation client")

	pkiImpl, err := pki.NewPki(vaultBackend, nil, config)
	DieOnErr(err, "could not build pki impl")

	sink, err := sink.MultiKeyPairSinkFromConfig(config)
	DieOnErr(err, "could not build keypair")

	cert, err := sink.ReadCert()
	DieOnErr(err, "can not read certificate")

	if pkg.IsCertExpired(*cert) {
		log.Info().Msg("Certificate is expired, no revocation needed")
		return
	}

	serial := pkg.FormatSerial(cert.SerialNumber)
	err = pkiImpl.Revoke(serial)
	DieOnErr(err, "could not revoke cert")
}

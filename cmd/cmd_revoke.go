package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/spf13/cobra"
	"strings"
)

func getRevokeCmd() *cobra.Command {
	var revokeCmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a x509 cert",
		Run:   revokeCertEntryPoint,
	}

	revokeCmd.Flags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "Certificate to read serial from")

	return revokeCmd
}

func revokeCertEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()
	config, err := config()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get config")
	}

	errors := append(config.Validate(), config.Validate()...)
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		log.Fatal().Msgf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
	}

	storage.InitBuilder(config)
	vaultClient, err := api.NewClient(getVaultConfig(config))
	if err != nil {
		log.Fatal().Err(err).Msgf("could not build vault client")
	}

	authStrategy, err := buildAuthImpl(vaultClient, config)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not build auth strategy")
	}

	vaultBackend, err := vault.NewVaultPki(vaultClient, authStrategy, config)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not build rotation client")
	}

	pkiImpl, err := pki.NewPki(vaultBackend, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not build pki impl")
	}

	sink, err := sink.MultiKeyPairSinkFromConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not build keypair")
	}

	content, err := sink.ReadCert()
	if err != nil {
		log.Fatal().Err(err).Msgf("can not read certificate")
	}

	serial := pkg.FormatSerial(content.SerialNumber)
	if err = pkiImpl.Revoke(serial); err != nil {
		log.Fatal().Err(err).Msg("could not revoke cert")
	}

}

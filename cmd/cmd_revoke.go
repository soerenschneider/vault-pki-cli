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
		RunE:  revokeCertEntryPoint,
	}

	revokeCmd.Flags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "Certificate to read serial from")

	return revokeCmd
}

func revokeCertEntryPoint(ccmd *cobra.Command, args []string) error {
	PrintVersionInfo()
	config, err := config()
	if err != nil {
		log.Fatal().Err(err)
	}

	errors := append(config.Validate(), config.Validate()...)
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		return fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
	}

	storage.InitBuilder(config)
	vaultClient, err := api.NewClient(getVaultConfig(config))
	if err != nil {
		return fmt.Errorf("could not build vault client: %v", err)
	}

	authStrategy, err := buildAuthImpl(vaultClient, config)
	if err != nil {
		return fmt.Errorf("could not build auth strategy: %v", err)
	}

	vaultBackend, err := vault.NewVaultPki(vaultClient, authStrategy, config)
	if err != nil {
		return fmt.Errorf("could not build rotation client: %v", err)
	}

	pkiImpl, err := pki.NewPki(vaultBackend, nil)
	if err != nil {
		return fmt.Errorf("could not build pki impl: %v", err)
	}

	sink, err := sink.MultiKeyPairSinkFromConfig(config)

	content, err := sink.ReadCert()
	if err != nil {
		return fmt.Errorf("can not read certificate: %v", err)
	}

	serial := pkg.FormatSerial(content.SerialNumber)
	return pkiImpl.Revoke(serial)

}

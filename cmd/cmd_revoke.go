package main

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/spf13/cobra"
	"os"
)

func getRevokeCmd() *cobra.Command {
	var revokeCmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a x509 cert",
		Run:   revokeCertEntryPoint,
	}

	revokeCmd.PersistentFlags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "Certificate to read serial from")

	revokeCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return revokeCmd
}

func revokeCertEntryPoint(ccmd *cobra.Command, args []string) {

	PrintVersionInfo()
	config, err := config()
	if err != nil {
		log.Fatal().Err(err)
	}

	err = revokeCert(*config)
	if err == nil {
		os.Exit(0)
	}
	log.Error().Msgf("Error revoking cert: %v", err)
	os.Exit(1)
}

func revokeCert(config conf.Config) error {
	/*
		errors := append(config.Validate(), config.Validate()...)
		if len(errors) > 0 {
			fmtErrors := make([]string, len(errors))
			for i, er := range errors {
				fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
			}
			return fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
		}

		vaultClient, err := api.NewClient(getVaultConfig(&config))
		if err != nil {
			return fmt.Errorf("could not build vault client: %v", err)
		}

		authStrategy, err := buildAuthImpl(vaultClient, &config)
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

		sink, err := buildOutput(config)

		content, err := sink.ReadCert()
		if err != nil {
			return fmt.Errorf("can not read certificate: %v", err)
		}

		serial, err := pkg.GetFormattedSerial(content.Raw)
		if err != nil {
			return fmt.Errorf("could not read certificate serial number: %v", err)
		}
		return pkiImpl.Revoke(serial)

	*/
	return nil
}

func buildRevokeSink(config conf.Config) (pki.CrlSink, error) {
	// TODO:
	return nil, errors.New("no sinks")
}

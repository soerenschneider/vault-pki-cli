package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/pods"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getRevokeCmd() *cobra.Command {
	var revokeCmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a x509 cert",
		Run:   revokeCertEntryPoint,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// See https://github.com/spf13/viper/issues/233#issuecomment-386791444
			viper.BindPFlag(conf.FLAG_CERTIFICATE_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_CERTIFICATE_FILE))

			return nil
		},
	}

	revokeCmd.PersistentFlags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "Certificate to read serial from")

	revokeCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return revokeCmd
}

func revokeCertEntryPoint(ccmd *cobra.Command, args []string) {

	PrintVersionInfo()
	configFile := viper.GetViper().GetString(conf.FLAG_CONFIG_FILE)
	var config *conf.Config
	if len(configFile) > 0 {
		var err error
		config, err = readConfig(configFile)
		if err != nil {
			log.Fatal().Msgf("Could not load desired config file: %s: %v", configFile, err)
		}
		log.Info().Msgf("Read config from file %s", viper.ConfigFileUsed())
	}

	config.PrintConfig()
	config.RevokeArguments.PrintConfig()

	err := revokeCert(*config)
	if err == nil {
		os.Exit(0)
	}
	log.Error().Msgf("Error revoking cert: %v", err)
	os.Exit(1)
}

func revokeCert(config conf.Config) error {
	errors := append(config.Validate(), config.RevokeArguments.Validate()...)
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
	certPod := &pods.FsPod{FilePath: config.RevokeArguments.CertificateFile}

	content, err := certPod.Read()
	if err != nil {
		return fmt.Errorf("can not read certificate: %v", err)
	}

	serial, err := pkg.GetFormattedSerial(content)
	if err != nil {
		return fmt.Errorf("could not read certificate serial number: %v", err)
	}
	return pkiImpl.Revoke(serial)
}

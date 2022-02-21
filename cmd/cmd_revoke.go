package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/vault-pki-cli/internal"
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
	var signCmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a x509 cert",
		Run:   revokeCertEntryPoint,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// See https://github.com/spf13/viper/issues/233#issuecomment-386791444
			viper.BindPFlag(conf.FLAG_CERTIFICATE_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_CERTIFICATE_FILE))

			return initializeConfig(cmd)
		},
	}

	signCmd.PersistentFlags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "Certificate to read serial from")

	signCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)

	return signCmd
}

func revokeCertEntryPoint(ccmd *cobra.Command, args []string) {
	initializeConfig(ccmd)

	log.Info().Msgf("Version %s (%s)", internal.BuildVersion, internal.CommitHash)
	configFile := viper.GetViper().GetString(conf.FLAG_CONFIG_FILE)
	if len(configFile) > 0 {
		err := readConfig(configFile)
		if err != nil {
			log.Fatal().Msgf("Could not load desired config file: %s: %v", configFile, err)
		}
		log.Info().Msgf("Read config from file %s", viper.ConfigFileUsed())
	}

	config := NewConfigFromViper()
	config.PrintConfig()
	config.RevokeArguments.PrintConfig()

	err := revokeCert(config)
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

	signingImpl, err := vault.NewVaultSigner(vaultClient, authStrategy, config)
	if err != nil {
		return fmt.Errorf("could not build rotation client: %v", err)
	}

	vaultPki, err := pki.NewPki(signingImpl, nil)
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
	return vaultPki.Revoke(serial)
}

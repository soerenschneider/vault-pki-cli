package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getReadAcmeCmd() *cobra.Command {
	var issueCmd = &cobra.Command{
		Use:   "read-acme",
		Short: "Reads an Acmevault x509 cert",
		Run:   readAcmeEntryPoint,
	}

	issueCmd.Flags().StringP(conf.FLAG_ISSUE_COMMON_NAME, "", "", "Specifies the requested CN for the certificate. If the CN is allowed by role policy, it will be issued.")
	issueCmd.Flags().StringP(conf.FLAG_METRICS_FILE, "", conf.FLAG_ACME_METRICS_FILE_DEFAULT, "File to write metrics to")
	issueCmd.Flags().StringP(conf.FLAG_READACME_ACME_PREFIX, "", conf.FLAG_READACME_ACME_PREFIX_DEFAULT, "Prefix for Acmevault kv2 secret paths")
	issueCmd.Flags().StringP(conf.FLAG_VAULT_MOUNT_KV2, "", conf.FLAG_VAULT_MOUNT_KV2_DEFAULT, "Mount path for kv2 secret")

	viper.SetDefault(conf.FLAG_METRICS_FILE, conf.FLAG_ACME_METRICS_FILE_DEFAULT)
	viper.SetDefault(conf.FLAG_READACME_ACME_PREFIX, conf.FLAG_READACME_ACME_PREFIX_DEFAULT)
	viper.SetDefault(conf.FLAG_VAULT_MOUNT_KV2, conf.FLAG_VAULT_MOUNT_KV2_DEFAULT)

	return issueCmd
}

func readAcmeEntryPoint(ccmd *cobra.Command, args []string) {
	PrintVersionInfo()

	config, err := config()
	if err != nil {
		log.Fatal().Err(err)
	}

	if len(config.CommonName) == 0 {
		log.Fatal().Msgf("No '%s' specified", conf.FLAG_ISSUE_COMMON_NAME)
	}
	config.Print()

	storage.InitBuilder(config)
	errs := readAcmeCert(config)
	if len(errs) > 0 {
		log.Error().Msgf("reading cert not successful, %v", errs)
		internal.MetricSuccess.WithLabelValues(config.CommonName).Set(0)
	} else {
		internal.MetricSuccess.WithLabelValues(config.CommonName).Set(1)
	}
	internal.MetricRunTimestamp.WithLabelValues(config.CommonName).SetToCurrentTime()
	if !config.Daemonize && len(config.MetricsFile) > 0 {
		internal.WriteMetrics(config.MetricsFile)
	}
}

func readAcmeCert(config *conf.Config) (errors []error) {
	vaultClient, err := api.NewClient(getVaultConfig(config))
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build vault client: %w", err))
		return
	}

	authStrategy, err := buildAuthImpl(vaultClient, config)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build auth strategy: %w", err))
		return
	}

	vaultBackend, err := vault.NewVaultPki(vaultClient, authStrategy, config)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build rotation client: %w", err))
		return
	}

	pkiImpl, err := pki.NewPki(vaultBackend, nil)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build pki impl: %w", err))
		return
	}

	sink, err := sink.MultiKeyPairSinkFromConfig(config)
	if err != nil {
		errors = append(errors, fmt.Errorf("can't build certificate output: %w", err))
		return
	}

	changed, err := pkiImpl.ReadAcme(sink, config)
	if err != nil {
		log.Error().Msgf("could not read certificate: %v", err)
		errors = append(errors, err)
	}

	if changed {
		log.Info().Msg("Detected update between local cert on disk and the read certificate")
		runPostIssueHooks(config)
	} else {
		log.Info().Msg("No update detected, local certificate and remote cert identical")
	}

	return errors
}

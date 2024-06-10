package main

import (
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
	issueCmd.Flags().StringP(conf.FLAG_METRICS_FILE, "", "", "File to write metrics to")
	issueCmd.Flags().StringP(conf.FLAG_READACME_ACME_PREFIX, "", conf.FLAG_READACME_ACME_PREFIX_DEFAULT, "Prefix for Acmevault kv2 secret paths")
	issueCmd.Flags().StringP(conf.FLAG_VAULT_MOUNT_KV2, "", conf.FLAG_VAULT_MOUNT_KV2_DEFAULT, "Mount path for kv2 secret")
	issueCmd.Flags().Uint64(conf.FLAG_RETRIES, conf.FLAG_RETRIES_DEFAULT, "How many retries to perform for non-permanent errors")

	viper.SetDefault(conf.FLAG_METRICS_FILE, "")
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
	internal.MetricSuccess.WithLabelValues(config.CommonName).Set(0)
	config.Print()

	storage.InitBuilder(config)
	err = readAcmeCert(config)
	if err != nil {
		log.Error().Err(err).Msg("reading cert not successful")
	} else {
		internal.MetricSuccess.WithLabelValues(config.CommonName).Set(1)
	}
	internal.MetricRunTimestamp.WithLabelValues(config.CommonName).SetToCurrentTime()
	if !config.Daemonize && len(config.MetricsFile) > 0 {
		err := internal.WriteMetrics(config.MetricsFile)
		if err != nil {
			log.Error().Err(err).Msg("could not write metrics")
		}
	}
}

func readAcmeCert(config *conf.Config) error {
	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "can't build client")

	authStrategy, err := buildAuthImpl(config)
	DieOnErr(err, "can't build auth")

	vaultBackend, err := vault.NewVaultPki(vaultClient, authStrategy, config)
	DieOnErr(err, "can't build vault pki")

	pkiImpl, err := pki.NewPki(vaultBackend, nil, config)
	DieOnErr(err, "can't build pki impl")

	sink, err := sink.MultiKeyPairSinkFromConfig(config)
	DieOnErr(err, "can't build sink")

	changed, err := pkiImpl.ReadAcme(sink, config)
	DieOnErr(err, "can't read acme cert")

	if !changed {
		log.Info().Msg("No update detected, local certificate and remote cert identical")
		return nil
	}

	log.Info().Msg("Detected update between local cert on disk and the read certificate")
	return runPostIssueHooks(config)
}

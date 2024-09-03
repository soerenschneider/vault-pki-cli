package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = authStrategy.Login(ctx, vaultClient)
	DieOnErr(err, "can't login to vault")

	opts := []vault.VaultOpts{
		vault.WithPkiMount(config.VaultMountPki),
		vault.WithKv2Mount(config.VaultMountKv2),
		vault.WithAcmePrefix(config.AcmePrefix),
	}

	vaultBackend, err := vault.NewVaultPki(vaultClient.Logical(), config.VaultPkiRole, opts...)
	DieOnErr(err, "can't build vault pki")

	pkiImpl, err := pki.NewPkiService(vaultBackend, nil)
	DieOnErr(err, "can't build pki impl")

	sink, err := storage.MultiKeyPairStorageFromConfig(config)
	DieOnErr(err, "can't build sink")

	result, err := pkiImpl.ReadAcme(ctx, sink, config.CommonName)
	if err != nil {
		labels := prometheus.Labels{
			"cn":  "config.CommonName",
			"err": internal.TranslateErrToPromLabel(err),
		}
		internal.MetricCertErrors.With(labels).Inc()
		DieOnErr(err, "can't read acme cert")
	}

	switch result.Status {
	case pkg.Issued:
		internal.UpdateCertificateMetrics(result.IssuedCert)
		log.Info().Msg("Detected update between local cert on disk and the read certificate")
		commandCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
		return runPostIssueHooks(commandCtx, config)
	case pkg.Noop:
		log.Info().Msg("No update detected for certificate")
		internal.UpdateCertificateMetrics(result.ExistingCert)
		return nil
	}

	return nil
}

package main

import (
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/spf13/viper"
	"go.uber.org/multierr"

	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const tickDuration = 1 * time.Hour

func getIssueCmd() *cobra.Command {
	var issueCmd = &cobra.Command{
		Use:   "issue",
		Short: "Issue a x509 cert",
		Run:   issueCertEntryPoint,
	}

	issueCmd.Flags().BoolP(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE, "", false, "Issue a new certificate regardless of the current certificate's lifetime")
	issueCmd.Flags().Float64P(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, "", conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT, "Create new certificate when a given threshold of its overall lifetime has been reached")
	issueCmd.Flags().StringP(conf.FLAG_ISSUE_COMMON_NAME, "", "", "Specifies the requested CN for the certificate. If the CN is allowed by role policy, it will be issued.")
	issueCmd.Flags().StringP(conf.FLAG_ISSUE_TTL, "", conf.FLAG_ISSUE_TTL_DEFAULT, "Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used. Note that the role values default to system values if not explicitly set.")
	issueCmd.Flags().StringP(conf.FLAG_METRICS_FILE, "", "", "File to write metrics to")
	issueCmd.Flags().StringP(conf.FLAG_ISSUE_METRICS_ADDR, "", conf.FLAG_ISSUE_METRICS_ADDR_DEFAULT, "File to write metrics to")
	issueCmd.Flags().BoolP(conf.FLAG_ISSUE_DAEMONIZE, "", conf.FLAG_ISSUE_DAEMONIZE_DEFAULT, "Run as daemon")
	issueCmd.Flags().StringArrayP(conf.FLAG_ISSUE_IP_SANS, "", []string{}, "Specifies requested IP Subject Alternative Names, in a comma-delimited list. Only valid if the role allows IP SANs (which is the default).")
	issueCmd.Flags().StringArrayP(conf.FLAG_ISSUE_ALT_NAMES, "", []string{}, "Specifies requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses; they will be parsed into their respective fields. If any requested names do not match role policy, the entire request will be denied.")
	issueCmd.Flags().StringSlice(conf.FLAG_ISSUE_HOOKS, []string{}, "Run commands after issuing a new certificate.")
	issueCmd.Flags().StringSlice(conf.FLAG_ISSUE_BACKEND_CONFIG, []string{}, "Backend config.")

	viper.SetDefault(conf.FLAG_ISSUE_TTL, conf.FLAG_ISSUE_TTL_DEFAULT)
	viper.SetDefault(conf.FLAG_ISSUE_DAEMONIZE, conf.FLAG_ISSUE_DAEMONIZE_DEFAULT)
	viper.SetDefault(conf.FLAG_ISSUE_METRICS_ADDR, conf.FLAG_ISSUE_METRICS_ADDR_DEFAULT)
	viper.SetDefault(conf.FLAG_METRICS_FILE, "")
	viper.SetDefault(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT)

	//issueCmd.MarkFlagRequired(conf.FLAG_ISSUE_COMMON_NAME)

	return issueCmd
}

func issueCertEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()
	config, err := config()
	DieOnErr(err, "could not get config")
	config.Print()

	err = config.ValidateIssue()
	DieOnErr(err, "invalid config")

	if config.Daemonize && len(config.MetricsAddr) > 0 {
		log.Info().Msgf("Starting metrics server at '%s'", config.MetricsAddr)
		go func() {
			err := internal.StartMetricsServer(config.MetricsAddr)
			DieOnErr(err, "could not start metrics server")
		}()
	}

	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	pkiImpl, sink := buildDependencies(config)
	for {
		err = issueCert(config, pkiImpl, sink)
		if err != nil {
			log.Error().Err(err).Msg("issuing cert not successful")
			internal.MetricSuccess.WithLabelValues(config.CommonName).Set(0)
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

		if !config.Daemonize {
			done <- true
		}

		select {
		case <-interrupt:
			log.Info().Msg("Received signal")
			return
		case <-done:
			if err != nil {
				log.Fatal().Err(err).Msg("encountered errors")
			}
			return
		case <-ticker.C:
			continue
		}
	}
}

func buildRenewalStrategy(config *conf.Config) (issue_strategies.IssueStrategy, error) {
	if config.ForceNewCertificate {
		return &issue_strategies.StaticRenewal{Decision: true}, nil
	}

	return issue_strategies.NewPercentage(config.CertificateLifetimeThresholdPercentage)
}

func buildDependencies(config *conf.Config) (*pki.PkiCli, pki.IssueSink) {
	storage.InitBuilder(config)

	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "can't build client")

	authStrategy, err := buildAuthImpl(vaultClient, config)
	DieOnErr(err, "can't build auth")

	vaultBackend, err := vault.NewVaultPki(vaultClient, authStrategy, config)
	DieOnErr(err, "can't build vault pki")

	strat, err := buildRenewalStrategy(config)
	DieOnErr(err, "can't build renewal strategy")

	pkiImpl, err := pki.NewPki(vaultBackend, strat)
	DieOnErr(err, "can't build pki impl")

	sink, err := sink.MultiKeyPairSinkFromConfig(config)
	DieOnErr(err, "can't build sink")

	return pkiImpl, sink
}

func issueCert(config *conf.Config, pkiImpl *pki.PkiCli, sink pki.IssueSink) error {
	var serial string
	x509cert, err := sink.ReadCert()
	if err == nil {
		serial = pkg.FormatSerial(x509cert.SerialNumber)
	}

	outcome, err := pkiImpl.Issue(sink, config)

	if outcome == pki.Issued && err == nil && len(serial) > 0 {
		if err := runPostIssueHooks(config); err != nil {
			log.Error().Err(err).Msg("Encountered errors while running post-issue hooks")
		}

		if err := pkiImpl.Revoke(serial); err != nil {
			log.Warn().Err(err).Msg("Revoking serial %s failed")
		}
	}

	tidyStorage(pkiImpl)

	return err
}

func tidyStorage(pkiImpl *pki.PkiCli) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
	if r.Intn(100) >= 90 {
		log.Info().Msgf("Tidying up certificate storage")
		err := pkiImpl.Tidy()
		if err != nil {
			log.Error().Msgf("Tidying up certificate storage failed: %v", err)
		}
	}
}

func runPostIssueHooks(config *conf.Config) error {
	var err error
	for _, hook := range config.PostHooks {
		log.Info().Msgf("Running command '%s'", hook)
		parsed := strings.Split(hook, " ")
		cmd := exec.Command(parsed[0], parsed[1:]...) // #nosec G204
		cmdErr := cmd.Run()
		if cmdErr != nil {
			err = multierr.Append(err, errors.Errorf("error running command '%s': %v", parsed[0], cmdErr))
		}
	}

	return err
}

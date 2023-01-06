package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const tickDuration = 1 * time.Hour

func getIssueCmd() *cobra.Command {
	var issueCmd = &cobra.Command{
		Use:   "issue",
		Short: "Issue a x509 cert",
		RunE:  issueCertEntryPoint,
	}

	issueCmd.Flags().BoolP(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE, "", false, "Issue a new certificate regardless of the current certificate's lifetime")
	issueCmd.Flags().Float64P(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, "", conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT, "Create new certificate when a given threshold of its overall lifetime has been reached")
	issueCmd.Flags().StringP(conf.FLAG_ISSUE_COMMON_NAME, "", "", "Specifies the requested CN for the certificate. If the CN is allowed by role policy, it will be issued.")
	issueCmd.Flags().StringP(conf.FLAG_ISSUE_TTL, "", conf.FLAG_ISSUE_TTL_DEFAULT, "Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used. Note that the role values default to system values if not explicitly set.")
	issueCmd.Flags().StringP(conf.FLAG_ISSUE_METRICS_FILE, "", conf.FLAG_ISSUE_METRICS_FILE_DEFAULT, "File to write metrics to")
	issueCmd.Flags().StringP(conf.FLAG_ISSUE_METRICS_ADDR, "", conf.FLAG_ISSUE_METRICS_ADDR_DEFAULT, "File to write metrics to")
	issueCmd.Flags().BoolP(conf.FLAG_ISSUE_DAEMONIZE, "", conf.FLAG_ISSUE_DAEMONIZE_DEFAULT, "Run as daemon")
	issueCmd.Flags().StringArrayP(conf.FLAG_ISSUE_IP_SANS, "", []string{}, "Specifies requested IP Subject Alternative Names, in a comma-delimited list. Only valid if the role allows IP SANs (which is the default).")
	issueCmd.Flags().StringArrayP(conf.FLAG_ISSUE_ALT_NAMES, "", []string{}, "Specifies requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses; they will be parsed into their respective fields. If any requested names do not match role policy, the entire request will be denied.")
	issueCmd.Flags().StringSlice(conf.FLAG_ISSUE_HOOKS, []string{}, "Run commands after issuing a new certificate.")
	issueCmd.Flags().StringSlice(conf.FLAG_ISSUE_BACKEND_CONFIG, []string{}, "Backend config.")

	viper.SetDefault(conf.FLAG_ISSUE_TTL, conf.FLAG_ISSUE_TTL_DEFAULT)
	viper.SetDefault(conf.FLAG_ISSUE_DAEMONIZE, conf.FLAG_ISSUE_DAEMONIZE_DEFAULT)
	viper.SetDefault(conf.FLAG_ISSUE_METRICS_ADDR, conf.FLAG_ISSUE_METRICS_ADDR_DEFAULT)
	viper.SetDefault(conf.FLAG_ISSUE_METRICS_FILE, conf.FLAG_ISSUE_METRICS_FILE_DEFAULT)
	viper.SetDefault(conf.FLAG_ISSUE_YUBIKEY_SLOT, conf.FLAG_ISSUE_YUBIKEY_SLOT_DEFAULT)
	viper.SetDefault(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT)

	//issueCmd.MarkFlagRequired(conf.FLAG_ISSUE_COMMON_NAME)

	return issueCmd
}

func issueCertEntryPoint(ccmd *cobra.Command, args []string) error {
	PrintVersionInfo()

	config, err := config()
	if err != nil {
		log.Fatal().Err(err)
	}

	config.Print()

	if config.Daemonize && len(config.MetricsAddr) > 0 {
		log.Info().Msgf("Starting metrics server at '%s'", config.MetricsAddr)
		go internal.StartMetricsServer(config.MetricsAddr)
	}

	errors := config.ValidateIssue()
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		log.Fatal().Msgf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", "))
	}

	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	done := make(chan bool, 1)

	var errs []error
	for {
		storage.InitBuilder(config)
		errs = issueCert(config)
		if len(errs) > 0 {
			log.Error().Msgf("issuing cert not successful, %v", errs)
			internal.MetricSuccess.WithLabelValues(config.CommonName).Set(0)
		} else {
			internal.MetricSuccess.WithLabelValues(config.CommonName).Set(1)
		}
		internal.MetricRunTimestamp.WithLabelValues(config.CommonName).SetToCurrentTime()
		if !config.Daemonize && len(config.MetricsFile) > 0 {
			internal.WriteMetrics(config.MetricsFile)
		}

		if !config.Daemonize {
			done <- true
		}

		select {
		case <-interrupt:
			log.Info().Msg("Received signal")
			if len(errs) > 0 {
				return fmt.Errorf("encountered errors: %v", errs)
			}
			return nil
		case <-done:
			if len(errs) > 0 {
				return fmt.Errorf("encountered errors: %v", errs)
			}
			return nil
		case <-ticker.C:
			continue
		}
	}
}

func issueCert(config *conf.Config) (errors []error) {
	vaultClient, err := api.NewClient(getVaultConfig(config))
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build vault client: %v", err))
		return
	}

	authStrategy, err := buildAuthImpl(vaultClient, config)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build auth strategy: %v", err))
		return
	}

	vaultBackend, err := vault.NewVaultPki(vaultClient, authStrategy, config)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build rotation client: %v", err))
		return
	}

	var strat issue_strategies.IssueStrategy
	if config.ForceNewCertificate {
		strat = &issue_strategies.StaticRenewal{Decision: true}
	} else {
		strat, err = issue_strategies.NewPercentage(config.CertificateLifetimeThresholdPercentage)
		if err != nil {
			errors = append(errors, fmt.Errorf("could not build strategy: %v", err))
			return
		}
	}

	pkiImpl, err := pki.NewPki(vaultBackend, strat)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build pki impl: %v", err))
		return
	}

	sink, err := sink.MultiKeyPairSinkFromConfig(config)
	if err != nil {
		errors = append(errors, fmt.Errorf("can't build certificate output: %v", err))
		return
	}

	var serial string
	x509cert, err := sink.ReadCert()
	if err == nil {
		serial = pkg.FormatSerial(x509cert.SerialNumber)
	}

	outcome, err := pkiImpl.Issue(sink, config)
	if err != nil {
		log.Error().Msgf("could not issue new certificate: %v", err)
		errors = append(errors, err)
	}

	if outcome == pki.Issued && err == nil && len(serial) > 0 {
		errs := runPostIssueHooks(config)
		if len(errs) > 0 {
			log.Error().Msgf("Encountered errors while running post-issue hooks: %v", errs)
		}

		err := pkiImpl.Revoke(serial)
		if err != nil {
			log.Warn().Msgf("Revoking serial %s failed: %v", serial, err)
		}
	}

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100) >= 90 {
		log.Info().Msgf("Tidying up certificate storage")
		err := pkiImpl.Tidy()
		if err != nil {
			log.Error().Msgf("Tidying up certificate storage failed: %v", err)
		}
	}

	return
}

func runPostIssueHooks(config *conf.Config) (errs []error) {
	for _, hook := range config.PostHooks {
		log.Info().Msgf("Running command '%s'", hook)
		parsed := strings.Split(hook, " ")
		cmd := exec.Command(parsed[0], parsed[1:]...)
		err := cmd.Run()
		if err != nil {
			errs = append(errs, errors.Errorf("error running command '%s': %v", parsed[0], err))
		}
	}

	return
}

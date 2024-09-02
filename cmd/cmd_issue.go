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
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg/vault"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
	"golang.org/x/net/context"

	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const daemonRunInterval = 1 * time.Hour

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
	issueCmd.Flags().Uint64(conf.FLAG_RETRIES, conf.FLAG_RETRIES_DEFAULT, "How many retries to perform for non-permanent errors")

	viper.SetDefault(conf.FLAG_ISSUE_TTL, conf.FLAG_ISSUE_TTL_DEFAULT)
	viper.SetDefault(conf.FLAG_RETRIES, conf.FLAG_RETRIES_DEFAULT)
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
	DieOnErr(err, "invalid config", config)

	internal.MetricSuccess.WithLabelValues(config.CommonName).Set(0)
	internal.MetricRunTimestamp.WithLabelValues(config.CommonName).SetToCurrentTime()

	pkiImpl, sink := buildDependencies(config)
	ctx, cancel := context.WithCancel(context.Background())
	err = issueCert(ctx, config, pkiImpl, sink)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	if config.Daemonize {
		go runAsDaemon(ctx, config, pkiImpl, sink)
	} else {
		done <- true
	}

	select {
	case <-interrupt:
		log.Info().Msgf("got interrupt")
		cancel()
	case <-done:
		cancel()
	}

	DieOnErr(err, "issuing cert not successful", config)
	if len(config.MetricsFile) > 0 {
		if err := internal.WriteMetrics(config.MetricsFile); err != nil {
			log.Warn().Err(err).Msg("could not write metrics")
		}
	}
}

func runAsDaemon(ctx context.Context, config *conf.Config, pkiImpl *pki.PkiService, sink pki.IssueStorage) {
	if config.Daemonize && len(config.MetricsAddr) > 0 {
		log.Info().Msgf("Starting metrics server at '%s'", config.MetricsAddr)
		go func() {
			err := internal.StartMetricsServer(config.MetricsAddr)
			DieOnErr(err, "could not start metrics server", config)
		}()
	}

	ticker := time.NewTicker(daemonRunInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := issueCert(ctx, config, pkiImpl, sink)
			if err != nil {
				log.Error().Err(err).Msg("issuing cert not successful")
			}
		case <-ctx.Done():
			return
		}
	}
}

func issueCert(ctx context.Context, config *conf.Config, pkiImpl *pki.PkiService, sink pki.IssueStorage) error {
	var serial string
	cert, err := sink.ReadCert()
	if err == nil {
		serial = pkg.FormatSerial(cert.SerialNumber)
	}

	args := pkg.IssueArgs{
		CommonName: config.CommonName,
		Ttl:        config.Ttl,
		IpSans:     config.IpSans,
		AltNames:   config.AltNames,
	}

	outcome, err := pkiImpl.Issue(ctx, sink, args)
	if err != nil {
		return err
	}
	internal.MetricSuccess.WithLabelValues(config.CommonName).Set(1)

	if outcome == pki.Issued {
		// overwrite outer 'err'
		err = runPostIssueHooks(config)

		if !pkg.IsCertExpired(*cert) {
			err := pkiImpl.Revoke(ctx, serial)
			if err != nil {
				log.Warn().Err(err).Str("serial", serial).Msg("Revoking cert failed")
			}
		}
	}

	tidyStorage(ctx, pkiImpl)
	return err
}

func buildRenewalStrategy(config *conf.Config) (issue_strategies.IssueStrategy, error) {
	if config.ForceNewCertificate {
		return &issue_strategies.StaticRenewal{Decision: true}, nil
	}

	return issue_strategies.NewPercentage(config.CertificateLifetimeThresholdPercentage)
}

func buildDependencies(config *conf.Config) (*pki.PkiService, pki.IssueStorage) {
	storage.InitBuilder(config)

	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "can't build client", config)

	authStrategy, err := buildAuthImpl(config)
	DieOnErr(err, "can't build auth", config)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = authStrategy.Login(ctx, vaultClient)
	DieOnErr(err, "can't login to vault")

	opts := []vault.VaultOpts{
		vault.WithPkiMount(config.VaultMountPki),
		vault.WithKv2Mount(config.VaultMountKv2),
		vault.WithAcmePrefix(config.AcmePrefix),
	}

	vaultBackend, err := vault.NewVaultPki(vaultClient.Logical(), config.VaultPkiRole, opts...)
	DieOnErr(err, "can't build vault pki", config)

	strat, err := buildRenewalStrategy(config)
	DieOnErr(err, "can't build renewal strategy", config)

	pkiImpl, err := pki.NewPkiService(vaultBackend, strat)
	DieOnErr(err, "can't build pki impl", config)

	sink, err := storage.MultiKeyPairStorageFromConfig(config)
	DieOnErr(err, "can't build sink", config)

	return pkiImpl, sink
}

func tidyStorage(ctx context.Context, pkiImpl *pki.PkiService) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
	if r.Intn(100) >= 90 {
		log.Info().Msgf("Tidying up certificate storage")
		err := pkiImpl.Tidy(ctx)
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

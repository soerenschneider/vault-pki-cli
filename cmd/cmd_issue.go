package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/conf/issue_sinks"
	"github.com/soerenschneider/vault-pki-cli/internal/pods"
	"github.com/soerenschneider/vault-pki-cli/internal/sink"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getIssueCmd() *cobra.Command {
	var issueCmd = &cobra.Command{
		Use:   "issue",
		Short: "Issue a x509 cert",
		Run:   issueCertEntryPoint,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// See https://github.com/spf13/viper/issues/233#issuecomment-386791444
			viper.BindPFlag(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE))
			viper.BindPFlag(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE))
			viper.BindPFlag(conf.FLAG_FILE_OWNER, cmd.PersistentFlags().Lookup(conf.FLAG_FILE_OWNER))
			viper.BindPFlag(conf.FLAG_FILE_GROUP, cmd.PersistentFlags().Lookup(conf.FLAG_FILE_GROUP))

			viper.BindPFlag(conf.FLAG_ISSUE_COMMON_NAME, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_COMMON_NAME))
			viper.BindPFlag(conf.FLAG_ISSUE_TTL, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_TTL))
			viper.BindPFlag(conf.FLAG_ISSUE_METRICS_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_METRICS_FILE))
			viper.BindPFlag(conf.FLAG_ISSUE_IP_SANS, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_IP_SANS))
			viper.BindPFlag(conf.FLAG_ISSUE_ALT_NAMES, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_ALT_NAMES))
			viper.BindPFlag(conf.FLAG_ISSUE_HOOKS, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_HOOKS))
			viper.BindPFlag(conf.FLAG_CONFIG_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_CONFIG_FILE))

			return nil
		},
	}

	issueCmd.PersistentFlags().BoolP(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE, "", false, "Issue a new certificate regardless of the current certificate's lifetime")
	issueCmd.PersistentFlags().Float64P(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, "", conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT, "Create new certificate when a given threshold of its overall lifetime has been reached")
	issueCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_COMMON_NAME, "", "", "Specifies the requested CN for the certificate. If the CN is allowed by role policy, it will be issued.")
	issueCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_TTL, "", conf.FLAG_ISSUE_TTL_DEFAULT, "Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used. Note that the role values default to system values if not explicitly set.")
	issueCmd.PersistentFlags().StringP(conf.FLAG_CONFIG_FILE, "", "", "Config.")
	issueCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_METRICS_FILE, "", conf.FLAG_ISSUE_METRICS_FILE_DEFAULT, "File to write metrics to")
	issueCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_IP_SANS, "", []string{}, "Specifies requested IP Subject Alternative Names, in a comma-delimited list. Only valid if the role allows IP SANs (which is the default).")
	issueCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_ALT_NAMES, "", []string{}, "Specifies requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses; they will be parsed into their respective fields. If any requested names do not match role policy, the entire request will be denied.")
	issueCmd.PersistentFlags().StringSlice(conf.FLAG_ISSUE_HOOKS, []string{}, "Run commands after issuing a new certificate.")

	issueCmd.MarkFlagRequired(conf.FLAG_ISSUE_COMMON_NAME)

	return issueCmd
}

func issueCertEntryPoint(ccmd *cobra.Command, args []string) {
	PrintVersionInfo()

	viper.SetConfigType("yaml")
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
	conf.ParseFlags(config)

	err := config.BuildBackends()
	if err != nil {
		log.Fatal().Err(err).Msg("could not build all sink")
	}

	config.Validate()

	config.PrintConfig()
	config.IssueArguments.PrintConfig()

	errs := issueCert(*config)
	if len(errs) > 0 {
		log.Error().Msgf("issuing cert not successful, %v", errs)
		internal.MetricSuccess.Set(0)
	} else {
		internal.MetricSuccess.Set(1)
	}
	internal.MetricRunTimestamp.SetToCurrentTime()
	if len(config.IssueArguments.MetricsFile) > 0 {
		internal.WriteMetrics(config.IssueArguments.MetricsFile)
	}

	if len(errs) == 0 {
		os.Exit(0)
	}
	os.Exit(1)
}

func issueCert(config conf.Config) (errors []error) {
	errors = append(config.Validate(), config.IssueArguments.Validate()...)
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		errors = append(errors, fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", ")))
		return
	}

	vaultClient, err := api.NewClient(getVaultConfig(&config))
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build vault client: %v", err))
		return
	}

	authStrategy, err := buildAuthImpl(vaultClient, &config)
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
	if config.IssueArguments.ForceNewCertificate {
		strat = &issue_strategies.StaticRenewal{Decision: true}
	} else {
		strat, err = issue_strategies.NewPercentage(config.IssueArguments.CertificateLifetimeThresholdPercentage)
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

	format, err := buildOutput(config)
	if err != nil {
		errors = append(errors, fmt.Errorf("can't build certificate output: %v", err))
		return
	}

	var serial string
	x509cert, err := format.Read()
	if err == nil {
		serial = pkg.FormatSerial(x509cert.SerialNumber)
	}

	outcome, err := pkiImpl.Issue(format, config.IssueArguments)
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

func buildOutput(config conf.Config) (pki.CertSink, error) {
	var builtSinks []pki.CertSink

	for _, backend := range config.Backends {
		switch backend.GetType() {
		case issue_sinks.FsType:
			sink, err := buildPemBackend(backend.(*issue_sinks.FilesystemSink))
			if err != nil {
				return nil, fmt.Errorf("could not build pem sink: %v", err)
			}
			builtSinks = append(builtSinks, sink)
		case issue_sinks.K8sType:
			sink, err := buildK8sBackend(backend.(*issue_sinks.K8sSink))
			if err != nil {
				return nil, fmt.Errorf("could not build k8s sink: %v", err)
			}
			builtSinks = append(builtSinks, sink)
		case issue_sinks.YubiType:
			sink, err := buildYubikeyBackend(backend.(*issue_sinks.YubikeySink))
			if err != nil {
				return nil, fmt.Errorf("could not build yubi sink: %v", err)
			}
			builtSinks = append(builtSinks, sink)
		default:
			return nil, fmt.Errorf("unknown sink requested: %s", backend.GetType())
		}
	}

	return sink.NewMultiSink(builtSinks...)
}

func buildK8sBackend(sinkConfig *issue_sinks.K8sSink) (pki.CertSink, error) {
	return sink.NewK8sBackend(
		sink.WithNamespace(sinkConfig.Namespace),
		sink.WithSecretName(sinkConfig.SecretName),
	)
}

func buildPemBackend(sinkConfig *issue_sinks.FilesystemSink) (pki.CertSink, error) {
	privateKeyPod, err := pods.NewFsPod(sinkConfig.PrivateKeyFile, sinkConfig.FileOwner, sinkConfig.FileGroup)
	if err != nil {
		return nil, fmt.Errorf("could not init private-key-file: %v", err)
	}

	certPod, err := pods.NewFsPod(sinkConfig.CertificateFile, sinkConfig.FileOwner, sinkConfig.FileGroup)
	if err != nil {
		return nil, fmt.Errorf("could not init cert-file: %v", err)
	}

	var caPod pki.KeyPod
	if len(sinkConfig.CaFile) > 0 {
		var err error
		caPod, err = pods.NewFsPod(sinkConfig.CaFile, sinkConfig.FileOwner, sinkConfig.FileGroup)
		if err != nil {
			return nil, fmt.Errorf("could not init ca-file: %v", err)
		}
	}
	return sink.NewPemSink(certPod, privateKeyPod, caPod)
}

func buildYubikeyBackend(sinkConfig *issue_sinks.YubikeySink) (pki.CertSink, error) {
	pin := sinkConfig.YubikeyPin
	if len(pin) == 0 {
		var err error
		pin, err = QueryYubikeyPin()
		if err != nil {
			return nil, err
		}
	}

	yubikey, err := pods.NewYubikeyPod(sinkConfig.YubikeySlot, pin)
	if err != nil {
		return nil, fmt.Errorf("can't init yubikey: %v", err)
	}

	yubikeyBackend, err := sink.NewYubikeySink(yubikey)
	if err != nil {
		return nil, fmt.Errorf("can't build yubikey backend: %v", err)
	}
	return yubikeyBackend, nil

	return nil, errors.New("todo")
}

func runPostIssueHooks(config conf.Config) (errs []error) {
	for _, hook := range config.PostIssueHooks {
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

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/soerenschneider/vault-pki-cli/internal/backends"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/pods"
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

			viper.BindPFlag(conf.FLAG_CERTIFICATE_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_CERTIFICATE_FILE))
			viper.BindPFlag(conf.FLAG_CA_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_CA_FILE))
			viper.BindPFlag(conf.FLAG_ISSUE_PRIVATE_KEY_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_PRIVATE_KEY_FILE))
			viper.BindPFlag(conf.FLAG_ISSUE_YUBIKEY_SLOT, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_YUBIKEY_SLOT))
			viper.BindPFlag(conf.FLAG_ISSUE_YUBIKEY_PIN, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_YUBIKEY_PIN))
			viper.BindPFlag(conf.FLAG_ISSUE_COMMON_NAME, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_COMMON_NAME))
			viper.BindPFlag(conf.FLAG_ISSUE_TTL, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_TTL))
			viper.BindPFlag(conf.FLAG_ISSUE_METRICS_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_METRICS_FILE))
			viper.BindPFlag(conf.FLAG_ISSUE_IP_SANS, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_IP_SANS))
			viper.BindPFlag(conf.FLAG_ISSUE_ALT_NAMES, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_ALT_NAMES))
			viper.BindPFlag(conf.FLAG_ISSUE_HOOKS, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_HOOKS))

			return initializeConfig(cmd)
		},
	}

	issueCmd.PersistentFlags().BoolP(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE, "", false, "Issue a new certificate regardless of the current certificate's lifetime")
	issueCmd.PersistentFlags().Float64P(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, "", conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT, "Create new certificate when a given threshold of its overall lifetime has been reached")
	issueCmd.PersistentFlags().Uint32(conf.FLAG_ISSUE_YUBIKEY_SLOT, math.MaxUint32, "Yubikey slot to write x509 data to")
	issueCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_YUBIKEY_PIN, "", "", "PIN to access to Yubikey PIV")
	issueCmd.PersistentFlags().StringSliceP(conf.FLAG_CERTIFICATE_FILE, "c", []string{}, "File to write the certificate to")
	issueCmd.PersistentFlags().StringSliceP(conf.FLAG_ISSUE_PRIVATE_KEY_FILE, "p", []string{}, "File to write the private key to")
	issueCmd.PersistentFlags().StringSlice(conf.FLAG_CA_FILE, []string{}, "File to write the CA certificate to")
	issueCmd.PersistentFlags().StringSlice(conf.FLAG_FILE_OWNER, []string{}, "Owner of the written files")
	issueCmd.PersistentFlags().StringSlice(conf.FLAG_FILE_GROUP, []string{}, "Group of the written files")
	issueCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_COMMON_NAME, "", "", "Specifies the requested CN for the certificate. If the CN is allowed by role policy, it will be issued.")
	issueCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_TTL, "", "48h", "Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used. Note that the role values default to system values if not explicitly set.")
	issueCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_METRICS_FILE, "", conf.FLAG_ISSUE_METRICS_FILE_DEFAULT, "File to write metrics to")
	issueCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_IP_SANS, "", []string{}, "Specifies requested IP Subject Alternative Names, in a comma-delimited list. Only valid if the role allows IP SANs (which is the default).")
	issueCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_ALT_NAMES, "", []string{}, "Specifies requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses; they will be parsed into their respective fields. If any requested names do not match role policy, the entire request will be denied.")
	issueCmd.PersistentFlags().StringSlice(conf.FLAG_ISSUE_HOOKS, []string{}, "Run commands after issuing a new certificate.")

	issueCmd.MarkFlagRequired(conf.FLAG_ISSUE_COMMON_NAME)

	return issueCmd
}

func issueCertEntryPoint(ccmd *cobra.Command, args []string) {
	PrintVersionInfo()

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
	config.IssueArguments.PrintConfig()

	err := issueCert(config)
	if len(err) > 0 {
		log.Error().Msgf("issuing cert not successful, %v", err)
		internal.MetricSuccess.Set(0)
	} else {
		internal.MetricSuccess.Set(1)
	}
	internal.MetricRunTimestamp.SetToCurrentTime()
	if len(config.IssueArguments.MetricsFile) > 0 {
		internal.WriteMetrics(config.IssueArguments.MetricsFile)
	}

	if len(err) == 0 {
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

func buildOutput(config conf.Config) (pki.CertBackend, error) {
	if config.UsesYubikey() {
		log.Info().Msg("Building yubikey backend to write cert data to")
		return buildYubikeyBackend(config)
	} else {
		log.Info().Msg("Building pem backend to write cert data to")
		return buildPemBackend(config)
	}

	return nil, errors.New("can't decide which backend to build")
}

func buildPemBackend(config conf.Config) (pki.CertBackend, error) {
	var pemBackends []pki.CertBackend
	for _, backend := range config.Backends {

		privateKeyPod, err := pods.NewFsPod(backend.PrivateKeyFile, backend.FileOwner, backend.FileGroup)
		if err != nil {
			return nil, fmt.Errorf("could not init private-key-file: %v", err)
		}

		certPod, err := pods.NewFsPod(backend.CertificateFile, backend.FileOwner, backend.FileGroup)
		if err != nil {
			return nil, fmt.Errorf("could not init cert-file: %v", err)
		}

		var caPod pki.KeyPod
		if len(backend.CaFile) > 0 {
			var err error
			caPod, err = pods.NewFsPod(backend.CaFile, backend.FileOwner, backend.FileGroup)
			if err != nil {
				return nil, fmt.Errorf("could not init ca-file: %v", err)
			}
		}
		pamBackend, err := backends.NewPemBackend(certPod, privateKeyPod, caPod)
		if err != nil {
			return nil, err
		}
		pemBackends = append(pemBackends, pamBackend)
	}

	return backends.NewMultiBackend(pemBackends...)
}

func buildYubikeyBackend(config conf.Config) (pki.CertBackend, error) {
	pin := config.YubikeyPin
	if len(pin) == 0 {
		var err error
		pin, err = QueryYubikeyPin()
		if err != nil {
			return nil, err
		}
	}

	yubikey, err := pods.NewYubikeyPod(config.YubikeySlot, pin)
	if err != nil {
		return nil, fmt.Errorf("can't init yubikey: %v", err)
	}

	yubikeyBackend, err := backends.NewYubikeyBackend(yubikey)
	if err != nil {
		return nil, fmt.Errorf("can't build yubikey backend: %v", err)
	}
	return yubikeyBackend, nil
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

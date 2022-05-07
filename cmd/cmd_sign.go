package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/pods"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getSignCmd() *cobra.Command {
	var signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Sign a CSR",
		Run:   signCertEntryPoint,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// See https://github.com/spf13/viper/issues/233#issuecomment-386791444
			viper.BindPFlag(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE))
			viper.BindPFlag(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE))
			viper.BindPFlag(conf.FLAG_FILE_OWNER, cmd.PersistentFlags().Lookup(conf.FLAG_FILE_OWNER))
			viper.BindPFlag(conf.FLAG_FILE_GROUP, cmd.PersistentFlags().Lookup(conf.FLAG_FILE_GROUP))

			viper.BindPFlag(conf.FLAG_CERTIFICATE_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_CERTIFICATE_FILE))
			viper.BindPFlag(conf.FLAG_CSR_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_CSR_FILE))
			viper.BindPFlag(conf.FLAG_ISSUE_COMMON_NAME, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_COMMON_NAME))
			viper.BindPFlag(conf.FLAG_ISSUE_TTL, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_TTL))
			viper.BindPFlag(conf.FLAG_ISSUE_METRICS_FILE, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_METRICS_FILE))
			viper.BindPFlag(conf.FLAG_ISSUE_IP_SANS, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_IP_SANS))
			viper.BindPFlag(conf.FLAG_ISSUE_ALT_NAMES, cmd.PersistentFlags().Lookup(conf.FLAG_ISSUE_ALT_NAMES))

			return initializeConfig(cmd)
		},
	}

	signCmd.PersistentFlags().BoolP(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE, "", false, "Issue a new certificate regardless of the current certificate's lifetime")
	signCmd.PersistentFlags().Float64P(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, "", conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT, "Create new certificate when a given threshold of its overall lifetime has been reached")
	signCmd.PersistentFlags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "File to write the certificate to")
	signCmd.PersistentFlags().StringP(conf.FLAG_CSR_FILE, "", "", "The CSR file to sign")
	signCmd.PersistentFlags().StringP(conf.FLAG_FILE_OWNER, "", conf.FLAG_FILE_OWNER_DEFAULT, "Owner of the written files")
	signCmd.PersistentFlags().StringP(conf.FLAG_FILE_GROUP, "", conf.FLAG_FILE_GROUP_DEFAULT, "Group of the written files")
	signCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_COMMON_NAME, "", "", "Specifies the requested CN for the certificate. If the CN is allowed by role policy, it will be issued.")
	signCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_TTL, "", "48h", "Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used. Note that the role values default to system values if not explicitly set.")
	signCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_METRICS_FILE, "", conf.FLAG_ISSUE_METRICS_FILE_DEFAULT, "File to write metrics to")
	signCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_IP_SANS, "", []string{}, "Specifies requested IP Subject Alternative Names, in a comma-delimited list. Only valid if the role allows IP SANs (which is the default).")
	signCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_ALT_NAMES, "", []string{}, "Specifies requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses; they will be parsed into their respective fields. If any requested names do not match role policy, the entire request will be denied.")

	signCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)
	signCmd.MarkFlagRequired(conf.FLAG_CSR_FILE)
	signCmd.MarkFlagRequired(conf.FLAG_ISSUE_COMMON_NAME)

	return signCmd
}

func signCertEntryPoint(ccmd *cobra.Command, args []string) {
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
	config.SignArguments.PrintConfig()

	err := signCert(config)
	if len(err) > 0 {
		log.Error().Msgf("signing cert not successful, %v", err)
		internal.MetricSuccess.Set(0)
	} else {
		internal.MetricSuccess.Set(1)
	}
	internal.MetricRunTimestamp.SetToCurrentTime()
	if len(config.SignArguments.MetricsFile) > 0 {
		internal.WriteMetrics(config.SignArguments.MetricsFile)
	}

	if len(err) == 0 {
		os.Exit(0)
	}
	os.Exit(1)
}

func signCert(config conf.Config) (errors []error) {
	errors = append(config.Validate(), config.SignArguments.Validate()...)
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

	pkiImpl, err := pki.NewPki(vaultBackend, &issue_strategies.StaticRenewal{Decision: false})
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build pki impl: %v", err))
		return
	}

	csrPod, err := pods.NewFsPod(config.PrivateKeyFile, config.SignArguments.FileOwner, config.SignArguments.FileGroup)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not init private-key-file: %v", err))
		return
	}
	certPod, err := pods.NewFsPod(config.SignArguments.CertificateFile, config.SignArguments.FileOwner, config.SignArguments.FileGroup)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not init cert-file: %v", err))
		return
	}

	err = pkiImpl.Sign(certPod, csrPod, config.SignArguments)
	if err != nil {
		log.Error().Msgf("could not sign CSR: %v", err)
		errors = append(errors, err)
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

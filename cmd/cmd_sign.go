package main

import (
	"math/rand"
	"time"

	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"

	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"

	"github.com/hashicorp/vault/api"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getSignCmd() *cobra.Command {
	var signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Sign a CSR",
		Run:   signCertEntryPoint,
	}

	signCmd.PersistentFlags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "File to write the certificate to")
	signCmd.PersistentFlags().StringP(conf.FLAG_CSR_FILE, "", "", "The CSR file to sign")
	signCmd.PersistentFlags().StringP(conf.FLAG_FILE_OWNER, "", "", "Owner of the written files. Defaults to the current user.")
	signCmd.PersistentFlags().StringP(conf.FLAG_FILE_GROUP, "", "", "Group of the written files. Defaults to the current user's primary group.")
	signCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_COMMON_NAME, "", "", "Specifies the requested CN for the certificate. If the CN is allowed by role policy, it will be issued.")
	signCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_TTL, "", conf.FLAG_ISSUE_TTL_DEFAULT, "Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used. Note that the role values default to system values if not explicitly set.")
	signCmd.PersistentFlags().StringP(conf.FLAG_METRICS_FILE, "", "", "File to write metrics to")
	signCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_IP_SANS, "", []string{}, "Specifies requested IP Subject Alternative Names, in a comma-delimited list. Only valid if the role allows IP SANs (which is the default).")
	signCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_ALT_NAMES, "", []string{}, "Specifies requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses; they will be parsed into their respective fields. If any requested names do not match role policy, the entire request will be denied.")

	signCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)
	signCmd.MarkFlagRequired(conf.FLAG_CSR_FILE)
	signCmd.MarkFlagRequired(conf.FLAG_ISSUE_COMMON_NAME)

	return signCmd
}

func signCertEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()

	config, err := config()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get config")
	}

	storage.InitBuilder(config)
	err = signCert(config)
	if err != nil {
		log.Error().Err(err).Msgf("signing CSR not successful")
		internal.MetricSuccess.WithLabelValues(config.CommonName).Set(0)
	} else {
		internal.MetricSuccess.WithLabelValues(config.CommonName).Set(1)
	}
	internal.MetricRunTimestamp.WithLabelValues(config.CommonName).SetToCurrentTime()
	if len(config.MetricsFile) > 0 {
		if err := internal.WriteMetrics(config.MetricsFile); err != nil {
			log.Error().Err(err).Msg("could not write metrics")
		}
	}

	if err != nil {
		log.Fatal().Err(err).Msg("encountered errors")
	}
}

func signCert(config *conf.Config) error {
	if err := config.Validate(); err != nil {
		return err
	}

	vaultClient, err := api.NewClient(getVaultConfig(config))
	if err != nil {
		return err
	}

	authStrategy, err := buildAuthImpl(vaultClient, config)
	if err != nil {
		return err
	}

	vaultBackend, err := vault.NewVaultPki(vaultClient, authStrategy, config)
	if err != nil {
		return err
	}

	pkiImpl, err := pki.NewPki(vaultBackend, &issue_strategies.StaticRenewal{Decision: false})
	if err != nil {
		return err
	}

	sink, err := sink.CsrSinkFromConfig(config.StorageConfig)
	if err != nil {
		return err
	}

	err = pkiImpl.Sign(sink, config)

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
	if r.Intn(100) >= 90 {
		log.Info().Msgf("Tidying up certificate storage")
		err := pkiImpl.Tidy()
		if err != nil {
			log.Error().Msgf("Tidying up certificate storage failed: %v", err)
		}
	}

	return err
}

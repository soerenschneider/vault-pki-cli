package main

import (
	"fmt"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/pki/sink"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"math/rand"
	"strings"
	"time"

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
		RunE:  signCertEntryPoint,
	}

	signCmd.PersistentFlags().StringP(conf.FLAG_CERTIFICATE_FILE, "c", "", "File to write the certificate to")
	signCmd.PersistentFlags().StringP(conf.FLAG_CSR_FILE, "", "", "The CSR file to sign")
	signCmd.PersistentFlags().StringP(conf.FLAG_FILE_OWNER, "", "", "Owner of the written files. Defaults to the current user.")
	signCmd.PersistentFlags().StringP(conf.FLAG_FILE_GROUP, "", "", "Group of the written files. Defaults to the current user's primary group.")
	signCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_COMMON_NAME, "", "", "Specifies the requested CN for the certificate. If the CN is allowed by role policy, it will be issued.")
	signCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_TTL, "", conf.FLAG_ISSUE_TTL_DEFAULT, "Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used. Note that the role values default to system values if not explicitly set.")
	signCmd.PersistentFlags().StringP(conf.FLAG_ISSUE_METRICS_FILE, "", conf.FLAG_ISSUE_METRICS_FILE_DEFAULT, "File to write metrics to")
	signCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_IP_SANS, "", []string{}, "Specifies requested IP Subject Alternative Names, in a comma-delimited list. Only valid if the role allows IP SANs (which is the default).")
	signCmd.PersistentFlags().StringArrayP(conf.FLAG_ISSUE_ALT_NAMES, "", []string{}, "Specifies requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses; they will be parsed into their respective fields. If any requested names do not match role policy, the entire request will be denied.")

	signCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)
	signCmd.MarkFlagRequired(conf.FLAG_CSR_FILE)
	signCmd.MarkFlagRequired(conf.FLAG_ISSUE_COMMON_NAME)

	return signCmd
}

func signCertEntryPoint(ccmd *cobra.Command, args []string) error {
	PrintVersionInfo()

	config, err := config()
	if err != nil {
		log.Fatal().Err(err)
	}

	storage.InitBuilder(config)
	errs := signCert(config)
	if len(errs) > 0 {
		log.Error().Msgf("signing CSR not successful, %v", errs)
		internal.MetricSuccess.WithLabelValues(config.CommonName).Set(0)
	} else {
		internal.MetricSuccess.WithLabelValues(config.CommonName).Set(1)
	}
	internal.MetricRunTimestamp.WithLabelValues(config.CommonName).SetToCurrentTime()
	if len(config.MetricsFile) > 0 {
		internal.WriteMetrics(config.MetricsFile)
	}

	if len(errs) == 0 {
		return fmt.Errorf("encountered errors: %v", errs)
	}
	return nil
}

func signCert(config *conf.Config) (errors []error) {
	errors = append(config.Validate(), config.Validate()...)
	if len(errors) > 0 {
		fmtErrors := make([]string, len(errors))
		for i, er := range errors {
			fmtErrors[i] = fmt.Sprintf("\"%s\"", er)
		}
		errors = append(errors, fmt.Errorf("invalid config, %d errors: %s", len(errors), strings.Join(fmtErrors, ", ")))
		return
	}

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

	pkiImpl, err := pki.NewPki(vaultBackend, &issue_strategies.StaticRenewal{Decision: false})
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build pki impl: %v", err))
		return
	}

	sink, err := sink.CsrSinkFromConfig(config.StorageConfig)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build sink: %v", err))
		return
	}

	err = pkiImpl.Sign(sink, config)
	if err != nil {
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

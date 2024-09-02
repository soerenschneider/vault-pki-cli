package main

import (
	"math/rand"
	"time"

	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"github.com/soerenschneider/vault-pki-cli/pkg/pki"
	"github.com/soerenschneider/vault-pki-cli/pkg/vault"
	"golang.org/x/net/context"

	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/pkg/issue_strategies"

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
	signCmd.PersistentFlags().Uint64(conf.FLAG_RETRIES, conf.FLAG_RETRIES_DEFAULT, "How many retries to perform for non-permanent errors")

	signCmd.MarkFlagRequired(conf.FLAG_CERTIFICATE_FILE)  // nolint:errcheck
	signCmd.MarkFlagRequired(conf.FLAG_CSR_FILE)          // nolint:errcheck
	signCmd.MarkFlagRequired(conf.FLAG_ISSUE_COMMON_NAME) // nolint:errcheck

	return signCmd
}

func signCertEntryPoint(_ *cobra.Command, _ []string) {
	PrintVersionInfo()

	config, err := config()
	DieOnErr(err, "could not get config")

	internal.MetricSuccess.WithLabelValues(config.CommonName).Set(0)
	internal.MetricRunTimestamp.WithLabelValues(config.CommonName).SetToCurrentTime()

	storage.InitBuilder(config)
	err = signCert(config)
	DieOnErr(err, "encountered errors", config)
	internal.MetricSuccess.WithLabelValues(config.CommonName).Set(1)

	if len(config.MetricsFile) > 0 {
		if err := internal.WriteMetrics(config.MetricsFile); err != nil {
			log.Error().Err(err).Msg("could not write metrics")
		}
	}
}

func signCert(config *conf.Config) error {
	if err := config.Validate(); err != nil {
		return err
	}

	vaultClient, err := buildVaultClient(config)
	DieOnErr(err, "can't build vault client")

	authStrategy, err := buildAuthImpl(config)
	DieOnErr(err, "can't build auth impl")

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
	DieOnErr(err, "can't build vault pki")

	pkiImpl, err := pki.NewPkiService(vaultBackend, &issue_strategies.StaticRenewal{Decision: false})
	DieOnErr(err, "can't build pki impl")

	sink, err := storage.CsrStorageFromConfig(config.StorageConfig)
	DieOnErr(err, "can't build sink")

	args := pkg.SignatureArgs{
		CommonName: config.CommonName,
		Ttl:        config.Ttl,
		IpSans:     config.IpSans,
		AltNames:   config.AltNames,
	}

	err = pkiImpl.Sign(ctx, sink, args)

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
	if r.Intn(100) >= 90 {
		log.Info().Msgf("Tidying up certificate storage")
		err := pkiImpl.Tidy(ctx)
		if err != nil {
			log.Error().Msgf("Tidying up certificate storage failed: %v", err)
		}
	}

	return err
}

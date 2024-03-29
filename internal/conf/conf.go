package conf

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

var sensitiveVars = map[string]struct{}{
	FLAG_VAULT_AUTH_APPROLE_ID:        {},
	FLAG_VAULT_AUTH_APPROLE_SECRET_ID: {},
	FLAG_VAULT_AUTH_TOKEN:             {},
}

type Config struct {
	VaultAddress      string `mapstructure:"vault-address" validate:"required"`
	VaultAuthMethod   string `mapstructure:"vault-auth-method" validate:"required"`
	VaultToken        string `mapstructure:"vault-auth-token" validate:"required_if=VaultAuthMethod token"`
	VaultAuthK8sRole  string `mapstructure:"vault-auth-k8s-role" validate:"required_if=VaultAuthMethod k8s"`
	VaultRoleId       string `mapstructure:"vault-auth-role-id" validate:"required_if=VaultAuthMethod approle"`
	VaultSecretId     string `mapstructure:"vault-auth-secret-id" validate:"required_if=VaultSecretIdFile '' VaultAuthMethod approle,excluded_unless=VaultSecretIdFile ''"`
	VaultSecretIdFile string `mapstructure:"vault-auth-secret-id-file" validate:"required_if=VaultSecretId '' VaultAuthMethod approle,excluded_unless=VaultSecretId ''"`
	VaultMountApprole string `mapstructure:"vault-approle-mount" validate:"required_if=VaultAuthMethod approle"`
	VaultMountPki     string `mapstructure:"vault-pki-mount" validate:"required"`
	VaultMountKv2     string `mapstructure:"vault-kv2-mount"`
	VaultPkiRole      string `mapstructure:"vault-pki-role-name" validate:"required"`

	Daemonize bool `mapstructure:"daemonize"`

	CommonName string   `mapstructure:"common-name"`
	Ttl        string   `mapstructure:"ttl"`
	IpSans     []string `mapstructure:"ip-sans"`
	AltNames   []string `mapstructure:"alt-names"`

	AcmePrefix string `mapstructure:"acme-prefix"`

	MetricsFile string `mapstructure:"metrics-file"`
	MetricsAddr string `mapstructure:"metrics-addr"`

	ForceNewCertificate bool                `mapstructure:"force-new-certificate"`
	StorageConfig       []map[string]string `mapstructure:"storage"`

	PostHooks                              []string `mapstructure:"post-hooks"`
	CertificateLifetimeThresholdPercentage float64  `mapstructure:"lifetime-threshold-percent"`

	DerEncoded bool
}

func (c *Config) Print() {
	log.Info().Msg("---")
	log.Info().Msg("Active config values:")
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsZero() {
			fieldName := val.Type().Field(i).Tag.Get("mapstructure")
			_, isSensitive := sensitiveVars[fieldName]
			if isSensitive {
				log.Info().Msgf("%s=*** (redacted)", fieldName)
			} else {
				log.Info().Msgf("%s=%v", fieldName, val.Field(i))
			}
		}
	}
	log.Info().Msg("---")
}

var (
	validate *validator.Validate
	once sync.Once
)

func (c *Config) Validate() error {
	once.Do(func() {
		validate = validator.New()
	})

	return validate.Struct(c)
}

func (c *Config) ValidateIssue() error {
	err := c.Validate()

	if len(c.CommonName) == 0 {
		err = multierr.Append(err, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_COMMON_NAME))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		err = multierr.Append(err, fmt.Errorf("'%s' must be [5, 90]", FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE))
	}

	return err
}

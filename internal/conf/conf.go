package conf

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

type Config struct {
	VaultAddress      string `mapstructure:"vault-address"`
	VaultToken        string `mapstructure:"vault-token"`
	VaultAuthK8sRole  string `mapstructure:"vault-k8s-role"`
	VaultRoleId       string `mapstructure:"vault-role-id"`
	VaultSecretId     string `mapstructure:"vault-secret-id"`
	VaultSecretIdFile string `mapstructure:"vault-secret-id-file"`
	VaultMountPki     string `mapstructure:"vault-mount-pki"`
	VaultMountApprole string `mapstructure:"vault-mount-approle"`
	VaultPkiRole      string `mapstructure:"vault-pki-role-name"`

	Daemonize bool `mapstructure:"daemonize"`

	CommonName string   `mapstructure:"common-name"`
	Ttl        string   `mapstructure:"ttl"`
	IpSans     []string `mapstructure:"ip-sans"`
	AltNames   []string `mapstructure:"alt-names"`

	MetricsFile string `mapstructure:"metrics-file"`
	MetricsAddr string `mapstructure:"metrics-addr"`

	ForceNewCertificate bool
	StorageConfig       []map[string]string `mapstructure:"storage"`

	PostIssueHooks                         []string `mapstructure:"post-hooks""`
	CertificateLifetimeThresholdPercentage float64  `mapstructure:"lifetime-threshold-percent"`

	DerEncoded bool
}

func (c *Config) Print() {
	log.Info().Msg("Active config values")
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsZero() {
			log.Info().Msgf("%s=%v", val.Type().Field(i).Name, val.Field(i))
		}
	}
	log.Info().Msg("---")
}

func (c *Config) Validate() []error {
	errs := make([]error, 0)

	emptyVaultToken := len(c.VaultToken) == 0
	emptyVaultAuthK8sRole := len(c.VaultAuthK8sRole) == 0
	emptyRoleId := len(c.VaultRoleId) == 0
	emptySecretId := len(c.VaultSecretId) == 0 && len(c.VaultSecretIdFile) == 0
	emptyAppRoleAuth := emptySecretId || emptyRoleId

	numAuthMethodsProvided := 0
	if !emptyVaultToken {
		numAuthMethodsProvided += 1
	}
	if !emptyAppRoleAuth {
		numAuthMethodsProvided += 1
	}
	if !emptyVaultAuthK8sRole {
		numAuthMethodsProvided += 1
	}

	if numAuthMethodsProvided == 0 {
		errs = append(errs, errors.New("no vault auth info provided. supply either token, AppRole or k8s auth info"))
	} else if numAuthMethodsProvided > 1 {
		errs = append(errs, fmt.Errorf("must provide only a single vault auth method, %d were provided", numAuthMethodsProvided))
	}

	if len(c.VaultSecretId) > 0 && len(c.VaultSecretIdFile) > 0 {
		errs = append(errs, fmt.Errorf("both '%s' and '%s' auth info provided, don't know what to pick", FLAG_VAULT_AUTH_APPROLE_SECRET_ID, FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE))
	}

	if len(c.VaultAddress) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_ADDRESS))
	}

	if len(c.VaultMountApprole) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_APPROLE_MOUNT))
	}

	if len(c.VaultMountPki) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_PKI_MOUNT))
	}

	if len(c.VaultPkiRole) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_PKI_BACKEND_ROLE))
	}

	return errs
}

func (c *Config) ValidateIssue() []error {
	errs := c.Validate()

	if len(c.CommonName) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_COMMON_NAME))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		errs = append(errs, fmt.Errorf("'%s' must be [5, 90]", FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE))
	}

	return errs
}

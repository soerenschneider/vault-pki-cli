package conf

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

type Config struct {
	VaultAddress      string `mapstructure:"vault_address"`
	VaultToken        string `mapstructure:"vault_token"`
	VaultRoleId       string `mapstructure:"vault_role_id"`
	VaultSecretId     string `mapstructure:"vault_secret_id"`
	VaultSecretIdFile string `mapstructure:"vault_secret_id_file"`
	VaultMountPki     string `mapstructure:"vault_mount_pki"`
	VaultMountApprole string `mapstructure:"vault_mount_approle"`
	VaultPkiRole      string `mapstructure:"vault_pki_role"`

	CommonName string   `mapstructure:"common_name"`
	Ttl        string   `mapstructure:"ttl"`
	IpSans     []string `mapstructure:"ip_sans"`
	AltNames   []string `mapstructure:"alt_names"`

	MetricsFile string

	ForceNewCertificate bool
	StorageConfig       []map[string]string `mapstructure:"storage"`

	PostIssueHooks                         []string `mapstructure:"post_hooks""`
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
	emptyRoleId := len(c.VaultRoleId) == 0
	emptySecretId := len(c.VaultSecretId) == 0 && len(c.VaultSecretIdFile) == 0
	emptyAppRoleAuth := emptySecretId || emptyRoleId
	if emptyAppRoleAuth && emptyVaultToken {
		errs = append(errs, fmt.Errorf("neither '%s' nor AppRole auth info provided", FLAG_VAULT_TOKEN))
	}

	if !emptyAppRoleAuth && !emptyVaultToken {
		errs = append(errs, fmt.Errorf("both '%s' and AppRole auth info provided, don't know what to pick", FLAG_VAULT_TOKEN))
	}

	if len(c.VaultSecretId) > 0 && len(c.VaultSecretIdFile) > 0 {
		errs = append(errs, fmt.Errorf("both '%s' and '%s' auth info provided, don't know what to pick", FLAG_VAULT_SECRET_ID, FLAG_VAULT_SECRET_ID_FILE))
	}

	if len(c.VaultAddress) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_ADDRESS))
	}

	if len(c.VaultMountApprole) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_MOUNT_APPROLE))
	}

	if len(c.VaultMountPki) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_MOUNT_PKI))
	}

	if len(c.VaultPkiRole) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_VAULT_PKI_BACKEND_ROLE))
	}

	return errs
}

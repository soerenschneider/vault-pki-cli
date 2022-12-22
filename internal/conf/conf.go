package conf

import (
	"fmt"
	log "github.com/rs/zerolog/log"
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

	SignArguments   `mapstructure:"sign"`
	IssueArguments  `mapstructure:"issue"`
	RevokeArguments `mapstructure:"revoke"`
	FetchArguments  `mapstructure:"fetch"`
}

func NewDefaultConfig() Config {
	return Config{
		VaultMountPki:     FLAG_VAULT_MOUNT_PKI_DEFAULT,
		VaultMountApprole: FLAG_VAULT_MOUNT_APPROLE_DEFAULT,
		VaultPkiRole:      FLAG_VAULT_PKI_BACKEND_ROLE_DEFAULT,

		SignArguments: SignArguments{
			Ttl:         FLAG_ISSUE_TTL_DEFAULT,
			FileOwner:   FLAG_FILE_OWNER_DEFAULT,
			MetricsFile: FLAG_ISSUE_METRICS_FILE_DEFAULT,
		},

		IssueArguments: IssueArguments{
			Ttl:                                    FLAG_ISSUE_TTL_DEFAULT,
			CertificateLifetimeThresholdPercentage: FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT,
			MetricsFile:                            FLAG_ISSUE_METRICS_FILE_DEFAULT,
		},

		RevokeArguments: RevokeArguments{},

		FetchArguments: FetchArguments{
			DerEncoded: false,
		},
	}
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

func (c *Config) PrintConfig() {
	log.Info().Msg("------------- Printing common config values -------------")
	log.Info().Msgf("%s=%s", FLAG_VAULT_ADDRESS, c.VaultAddress)
	if len(c.VaultToken) > 0 {
		log.Info().Msgf("%s=*** (sensitive output)", FLAG_VAULT_TOKEN)
	}
	if len(c.VaultRoleId) > 0 {
		log.Info().Msgf("%s=*** (sensitive output)", FLAG_VAULT_ROLE_ID)
	}
	if len(c.VaultSecretId) > 0 {
		log.Info().Msgf("%s=*** (sensitive output)", FLAG_VAULT_SECRET_ID)
	}
	if len(c.VaultSecretIdFile) > 0 {
		log.Info().Msgf("%s=%s", FLAG_VAULT_SECRET_ID_FILE, c.VaultSecretIdFile)
	}

	log.Info().Msgf("%s=%s", FLAG_VAULT_MOUNT_PKI, c.VaultMountPki)
	log.Info().Msgf("%s=%s", FLAG_VAULT_MOUNT_APPROLE, c.VaultMountApprole)
	log.Info().Msgf("%s=%s", FLAG_VAULT_PKI_BACKEND_ROLE, c.VaultPkiRole)
}

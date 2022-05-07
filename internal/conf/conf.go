package conf

import (
	"fmt"

	log "github.com/rs/zerolog/log"
)

type Config struct {
	VaultAddress      string
	VaultToken        string
	VaultRoleId       string
	VaultSecretId     string
	VaultSecretIdFile string
	VaultMountPki     string
	VaultMountApprole string
	VaultPkiRole      string

	SignArguments
	IssueArguments
	RevokeArguments
	FetchArguments
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

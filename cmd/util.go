package main

import (
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
)

func getVaultConfig(conf *conf.Config) *api.Config {
	vaultConfig := api.DefaultConfig()
	vaultConfig.MaxRetries = 13
	vaultConfig.Address = conf.VaultAddress
	return vaultConfig
}

func buildAuthImpl(client *api.Client, conf *conf.Config) (vault.AuthMethod, error) {
	token := conf.VaultToken
	if len(token) > 0 {
		log.Info().Msg("Building 'token' vault auth...")
		return vault.NewTokenAuth(token)
	}

	if len(conf.VaultAuthK8sRole) > 0 {
		log.Info().Msg("Building 'kubernetes' vault auth...")
		return vault.NewVaultKubernetesAuth(client, conf.VaultAuthK8sRole)
	}

	if conf.VaultAuthImplicit {
		log.Info().Msg("Building 'implicit' vault auth...")
		return vault.NewTokenImplicitAuth(), nil
	}

	approleData := make(map[string]string)
	approleData[vault.KeyRoleId] = conf.VaultRoleId
	approleData[vault.KeySecretId] = conf.VaultSecretId
	approleData[vault.KeySecretIdFile] = conf.VaultSecretIdFile

	log.Info().Msg("Building 'approle' vault auth...")
	return vault.NewAppRoleAuth(client, approleData, conf.VaultMountApprole)
}

func PrintVersionInfo() {
	log.Info().Msgf("Version %s (%s)", internal.BuildVersion, internal.CommitHash)
}

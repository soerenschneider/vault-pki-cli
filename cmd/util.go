package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

func viperOrEnv(viperKey, envKey string) string {
	val := viper.GetViper().GetString(viperKey)
	if len(val) == 0 {
		return os.Getenv(envKey)
	}
	return val
}

func NewConfigFromViper() conf.Config {
	config := conf.Config{}

	config.VaultAddress = viperOrEnv(conf.FLAG_VAULT_ADDRESS, "VAULT_ADDR")
	config.VaultToken = viper.GetViper().GetString(conf.FLAG_VAULT_TOKEN)
	config.VaultRoleId = viper.GetViper().GetString(conf.FLAG_VAULT_ROLE_ID)
	config.VaultSecretId = viper.GetViper().GetString(conf.FLAG_VAULT_SECRET_ID)
	config.VaultSecretIdFile = getExpandedFile(viper.GetViper().GetString(conf.FLAG_VAULT_SECRET_ID_FILE))
	config.VaultMountApprole = viper.GetViper().GetString(conf.FLAG_VAULT_MOUNT_APPROLE)
	config.VaultMountPki = viper.GetViper().GetString(conf.FLAG_VAULT_MOUNT_PKI)
	config.VaultPkiRole = viper.GetViper().GetString(conf.FLAG_VAULT_PKI_BACKEND_ROLE)
	config.IssueArguments.CertificateFile = getExpandedFile(viper.GetViper().GetString(conf.FLAG_CERTIFICATE_FILE))
	config.RevokeArguments.CertificateFile = getExpandedFile(viper.GetViper().GetString(conf.FLAG_CERTIFICATE_FILE))

	// fetch cmd args
	config.FetchArguments.OutputFile = getExpandedFile(viper.GetViper().GetString(conf.FLAG_OUTPUT_FILE))
	config.FetchArguments.DerEncoded = viper.GetViper().GetBool(conf.FLAG_DER_ENCODED)

	config.FileOwner = viper.GetViper().GetString(conf.FLAG_FILE_OWNER)
	config.FileGroup = viper.GetViper().GetString(conf.FLAG_FILE_GROUP)

	config.PrivateKeyFile = getExpandedFile(viper.GetViper().GetString(conf.FLAG_ISSUE_PRIVATE_KEY_FILE))

	config.ForceNewCertificate = viper.GetViper().GetBool(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE)
	config.CertificateLifetimeThresholdPercentage = viper.GetViper().GetFloat64(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE)
	config.MetricsFile = viper.GetViper().GetString(conf.FLAG_ISSUE_METRICS_FILE)

	// issue cmd args
	config.CommonName = viper.GetViper().GetString(conf.FLAG_ISSUE_COMMON_NAME)
	config.Ttl = viper.GetViper().GetString(conf.FLAG_ISSUE_TTL)
	config.IpSans = viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_IP_SANS)
	config.AltNames = viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_ALT_NAMES)

	return config
}

func getVaultConfig(conf *conf.Config) *api.Config {
	vaultConfig := api.DefaultConfig()
	vaultConfig.MaxRetries = 13
	vaultConfig.Address = conf.VaultAddress
	return vaultConfig
}

func buildAuthImpl(client *api.Client, conf *conf.Config) (vault.AuthMethod, error) {
	token := conf.VaultToken
	if len(token) > 0 {
		return vault.NewTokenAuth(token)
	}

	approleData := make(map[string]string)
	approleData[vault.KeyRoleId] = conf.VaultRoleId
	approleData[vault.KeySecretId] = conf.VaultSecretId
	approleData[vault.KeySecretIdFile] = conf.VaultSecretIdFile

	return vault.NewAppRoleAuth(client, approleData, conf.VaultMountApprole)
}

func getExpandedFile(filename string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if strings.HasPrefix(filename, "~/") {
		return filepath.Join(dir, filename[2:])
	}

	if strings.HasPrefix(filename, "$HOME/") {
		return filepath.Join(dir, filename[6:])
	}

	return filename
}

func handleFetchedData(certData []byte, config conf.Config) {
	if len(config.FetchArguments.OutputFile) == 0 {
		fmt.Println(string(certData))
		os.Exit(0)
	}

	err := ioutil.WriteFile(config.FetchArguments.OutputFile, certData, 0644)
	if err != nil {
		log.Error().Msgf("Error writing cert: %v", err)
		os.Exit(1)
	}
}

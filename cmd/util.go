package main

import (
	"fmt"
	"github.com/ilius/go-askpass"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
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

// getNumBackends checks how many arguments per slice have been supplied to determine
// the size of the backends slice.
func getNumBackends() int {
	max := len(viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_PRIVATE_KEY_FILE))

	if len(viper.GetViper().GetStringSlice(conf.FLAG_CERTIFICATE_FILE)) > max {
		max = len(viper.GetViper().GetStringSlice(conf.FLAG_CERTIFICATE_FILE))
	}
	if len(viper.GetViper().GetStringSlice(conf.FLAG_CA_FILE)) > max {
		max = len(viper.GetViper().GetStringSlice(conf.FLAG_CA_FILE))
	}
	if len(viper.GetViper().GetStringSlice(conf.FLAG_FILE_OWNER)) > max {
		max = len(viper.GetViper().GetStringSlice(conf.FLAG_FILE_OWNER))
	}
	if len(viper.GetViper().GetStringSlice(conf.FLAG_FILE_GROUP)) > max {
		max = len(viper.GetViper().GetStringSlice(conf.FLAG_FILE_GROUP))
	}

	return max
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

	// Revoke subcmd
	config.RevokeArguments.CertificateFile = getExpandedFile(viper.GetViper().GetString(conf.FLAG_CERTIFICATE_FILE))

	// Fetch subcmd
	config.FetchArguments.OutputFile = getExpandedFile(viper.GetViper().GetString(conf.FLAG_OUTPUT_FILE))
	config.FetchArguments.DerEncoded = viper.GetViper().GetBool(conf.FLAG_DER_ENCODED)

	// Issue subcmd
	config.IssueArguments.Backends = make([]conf.Backend, getNumBackends())
	for n, val := range viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_PRIVATE_KEY_FILE) {
		config.IssueArguments.Backends[n].PrivateKeyFile = getExpandedFile(val)
	}
	for n, val := range viper.GetViper().GetStringSlice(conf.FLAG_CERTIFICATE_FILE) {
		config.IssueArguments.Backends[n].CertificateFile = getExpandedFile(val)
	}
	for n, val := range viper.GetViper().GetStringSlice(conf.FLAG_CA_FILE) {
		config.IssueArguments.Backends[n].CaFile = getExpandedFile(val)
	}
	for n, val := range viper.GetViper().GetStringSlice(conf.FLAG_FILE_OWNER) {
		config.IssueArguments.Backends[n].FileOwner = getExpandedFile(val)
	}
	for n, val := range viper.GetViper().GetStringSlice(conf.FLAG_FILE_GROUP) {
		config.IssueArguments.Backends[n].FileGroup = getExpandedFile(val)
	}

	config.IssueArguments.PostIssueHooks = viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_HOOKS)

	config.IssueArguments.MetricsFile = viper.GetViper().GetString(conf.FLAG_ISSUE_METRICS_FILE)

	config.IssueArguments.ForceNewCertificate = viper.GetViper().GetBool(conf.FLAG_ISSUE_FORCE_NEW_CERTIFICATE)
	config.IssueArguments.CertificateLifetimeThresholdPercentage = viper.GetViper().GetFloat64(conf.FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE)

	config.IssueArguments.CommonName = viper.GetViper().GetString(conf.FLAG_ISSUE_COMMON_NAME)
	config.IssueArguments.Ttl = viper.GetViper().GetString(conf.FLAG_ISSUE_TTL)
	config.IssueArguments.IpSans = viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_IP_SANS)
	config.IssueArguments.AltNames = viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_ALT_NAMES)

	config.IssueArguments.YubikeyPin = viper.GetViper().GetString(conf.FLAG_ISSUE_YUBIKEY_PIN)
	config.IssueArguments.YubikeySlot = viper.GetViper().GetUint32(conf.FLAG_ISSUE_YUBIKEY_SLOT)

	// Sign subcmd
	config.SignArguments.CertificateFile = getExpandedFile(viper.GetViper().GetString(conf.FLAG_CERTIFICATE_FILE))
	config.SignArguments.FileOwner = viper.GetViper().GetString(conf.FLAG_FILE_OWNER)
	config.SignArguments.FileGroup = viper.GetViper().GetString(conf.FLAG_FILE_GROUP)
	config.SignArguments.MetricsFile = viper.GetViper().GetString(conf.FLAG_ISSUE_METRICS_FILE)
	config.SignArguments.CommonName = viper.GetViper().GetString(conf.FLAG_ISSUE_COMMON_NAME)
	config.SignArguments.Ttl = viper.GetViper().GetString(conf.FLAG_ISSUE_TTL)
	config.SignArguments.IpSans = viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_IP_SANS)
	config.SignArguments.AltNames = viper.GetViper().GetStringSlice(conf.FLAG_ISSUE_ALT_NAMES)
	config.SignArguments.CsrFile = viper.GetViper().GetString(conf.FLAG_CSR_FILE)

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

func QueryYubikeyPin() (string, error) {
	pin, err := askpass.Askpass("Please enter Yubikey PIN (won't echo)", false, "")
	if err != nil {
		return "", fmt.Errorf("can not read pin for yubikey: %v", err)
	}

	return pin, nil
}

func PrintVersionInfo() {
	log.Info().Msgf("Version %s (%s), YubikeySupport=%s", internal.BuildVersion, internal.CommitHash, internal.YubiKeySupport)
}

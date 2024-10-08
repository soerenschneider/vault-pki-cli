package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/spf13/pflag"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix             = "VAULT_PKI_CLI"
	defaultConfigFilename = "config"
)

func main() {
	initLogging()

	root := &cobra.Command{
		Use:   "vault-pki-cli",
		Short: fmt.Sprintf("Interact with Vault PKI - %s", internal.BuildVersion),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			var errs []error
			cmd.Flags().Visit(func(flag *pflag.Flag) {

				err := viper.BindPFlag(flag.Name, cmd.Flags().Lookup(flag.Name))
				if err != nil {
					errs = append(errs, err)
				}

			})
			if len(errs) > 0 {
				return fmt.Errorf("can't bind flags: %v", errs)
			}
			return nil
		},
	}

	root.PersistentFlags().BoolP(conf.FLAG_DEBUG, "v", false, "Enable verbose logging")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_ADDRESS, "a", "", "Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_AUTH_TOKEN, "t", "", "Vault token to use for authentication. Can not be used in conjunction with AppRole login data.")
	root.PersistentFlags().BoolP(conf.FLAG_VAULT_AUTH_IMPLICIT, "i", false, "Try to implicitly authenticate to vault using VAULT_TOKEN env var or ~/.vault-token file.")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_AUTH_K8S_ROLE, "k", "", "Kubernetes role to authenticate against vault")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_AUTH_APPROLE_ID, "r", "", "Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_AUTH_APPROLE_SECRET_ID, "s", "", "Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE, "", "", "Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_APPROLE_MOUNT, "", conf.FLAG_VAULT_MOUNT_APPROLE_DEFAULT, "Path where the AppRole auth method is mounted.")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_PKI_MOUNT, "", conf.FLAG_VAULT_MOUNT_PKI_DEFAULT, "Path where the PKI secret engine is mounted.")
	root.PersistentFlags().StringP(conf.FLAG_VAULT_PKI_BACKEND_ROLE, "", conf.FLAG_VAULT_PKI_BACKEND_ROLE_DEFAULT, "The name of the PKI role backend.")
	root.PersistentFlags().StringP(conf.FLAG_CONFIG_FILE, "", "", "File to read the config from")

	root.AddCommand(getRevokeCmd())
	root.AddCommand(getIssueCmd())
	root.AddCommand(getSignCmd())
	root.AddCommand(readCaCmd())
	root.AddCommand(readCaChainCmd())
	root.AddCommand(readCrlCmd())
	root.AddCommand(getReadAcmeCmd())
	root.AddCommand(versionCmd)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func expandPath(path string) string {
	if !strings.Contains(path, "~") {
		return path
	}

	usr, err := user.Current()
	if err != nil {
		return path
	}

	dir := usr.HomeDir

	if path == "~" {
		return dir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:])
	}

	return path
}

func config() (*conf.Config, error) {
	viper.SetDefault(conf.FLAG_VAULT_PKI_MOUNT, conf.FLAG_VAULT_MOUNT_PKI_DEFAULT)
	viper.SetDefault(conf.FLAG_VAULT_APPROLE_MOUNT, conf.FLAG_VAULT_MOUNT_APPROLE_DEFAULT)
	viper.SetDefault(conf.FLAG_VAULT_PKI_BACKEND_ROLE, conf.FLAG_VAULT_PKI_BACKEND_ROLE_DEFAULT)

	viper.SetConfigName(defaultConfigFilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/vault-pki-cli")
	viper.AddConfigPath("/etc/vault-pki-cli/")

	if viper.IsSet(conf.FLAG_CONFIG_FILE) {
		configFile := expandPath(viper.GetString(conf.FLAG_CONFIG_FILE))
		viper.SetConfigFile(configFile)
	}

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil && viper.IsSet(conf.FLAG_CONFIG_FILE) {
		log.Fatal().Msgf("Can't read config: %v", err)
	}

	var config *conf.Config

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal().Msgf("unable to decode into struct, %v", err)
	}

	setupLogLevel(config.Debug)
	return config, nil
}

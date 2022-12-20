package main

import (
	"fmt"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/spf13/pflag"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix             = "VAULT_PKI_CLI"
	defaultConfigFilename = "config"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})

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
				log.Info().Msgf("%s=%v", flag.Name, flag.Value)

			})
			if len(errs) > 0 {
				return fmt.Errorf("can't bind flags: %v", errs)
			}
			return nil
		},
	}

	root.AddCommand(getRevokeCmd())
	root.AddCommand(getIssueCmd())
	root.AddCommand(getSignCmd())
	root.AddCommand(readCaCmd())
	root.AddCommand(readCaChainCmd())
	root.AddCommand(readCrlCmd())
	root.AddCommand(versionCmd)

	root.Flags().BoolP("debug", "v", false, "Enable debug logging")
	root.Flags().StringP(conf.FLAG_VAULT_ADDRESS, "a", "", "Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.")
	root.Flags().StringP(conf.FLAG_VAULT_TOKEN, "t", "", "Vault token to use for authentication. Can not be used in conjunction with AppRole login data.")
	root.Flags().StringP(conf.FLAG_VAULT_ROLE_ID, "r", "", "Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	root.Flags().StringP(conf.FLAG_VAULT_SECRET_ID, "s", "", "Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	root.Flags().StringP(conf.FLAG_VAULT_SECRET_ID_FILE, "", "", "Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.")
	root.Flags().StringP(conf.FLAG_VAULT_MOUNT_PKI, "", conf.FLAG_VAULT_MOUNT_PKI_DEFAULT, "Path where the PKI secret engine is mounted.")
	root.Flags().StringP(conf.FLAG_VAULT_MOUNT_APPROLE, "", conf.FLAG_VAULT_MOUNT_APPROLE_DEFAULT, "Path where the AppRole auth method is mounted.")
	root.Flags().StringP(conf.FLAG_VAULT_PKI_BACKEND_ROLE, "", conf.FLAG_VAULT_PKI_BACKEND_ROLE_DEFAULT, "The name of the PKI role backend.")

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func config() (*conf.Config, error) {
	viper.SetDefault(strings.Replace(conf.FLAG_VAULT_MOUNT_PKI, "-", "_", -1), conf.FLAG_VAULT_MOUNT_PKI_DEFAULT)
	viper.SetDefault(strings.Replace(conf.FLAG_VAULT_MOUNT_APPROLE, "-", "_", -1), conf.FLAG_VAULT_MOUNT_APPROLE_DEFAULT)
	viper.SetDefault(strings.Replace(conf.FLAG_VAULT_PKI_BACKEND_ROLE, "-", "_", -1), conf.FLAG_VAULT_PKI_BACKEND_ROLE_DEFAULT)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/vault-pki-cli")
	viper.AddConfigPath("/etc/vault-pki-cli/")

	if viper.IsSet(conf.FLAG_CONFIG_FILE) {
		configFile := viper.GetString(conf.FLAG_CONFIG_FILE)
		log.Info().Msgf("Trying to read config from '%s'", configFile)
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

	return config, nil
}

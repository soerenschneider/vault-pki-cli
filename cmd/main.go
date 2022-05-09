package main

import (
	"fmt"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
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

	root := &cobra.Command{Use: "vault-pki-cli", Short: fmt.Sprintf("Interact with Vault PKI - %s", internal.BuildVersion)}

	root.AddCommand(getRevokeCmd())
	root.AddCommand(getIssueCmd())
	root.AddCommand(getSignCmd())
	root.AddCommand(readCaCmd())
	root.AddCommand(readCaChainCmd())
	root.AddCommand(readCrlCmd())
	root.AddCommand(versionCmd)

	root.PersistentFlags().BoolP("debug", "v", false, "Enable debug logging")

	root.PersistentFlags().StringP(conf.FLAG_VAULT_ADDRESS, "a", "", "Vault instance to connect to. If not specified, falls back to env var VAULT_ADDR.")
	viper.BindPFlag(conf.FLAG_VAULT_ADDRESS, root.PersistentFlags().Lookup(conf.FLAG_VAULT_ADDRESS))

	root.PersistentFlags().StringP(conf.FLAG_VAULT_TOKEN, "t", "", "Vault token to use for authentication. Can not be used in conjunction with AppRole login data.")
	viper.BindPFlag(conf.FLAG_VAULT_TOKEN, root.PersistentFlags().Lookup(conf.FLAG_VAULT_TOKEN))

	root.PersistentFlags().StringP(conf.FLAG_VAULT_ROLE_ID, "r", "", "Vault role_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(conf.FLAG_VAULT_ROLE_ID, root.PersistentFlags().Lookup(conf.FLAG_VAULT_ROLE_ID))

	root.PersistentFlags().StringP(conf.FLAG_VAULT_SECRET_ID, "s", "", "Vault secret_id to use for AppRole login. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(conf.FLAG_VAULT_SECRET_ID, root.PersistentFlags().Lookup(conf.FLAG_VAULT_SECRET_ID))

	root.PersistentFlags().StringP(conf.FLAG_VAULT_SECRET_ID_FILE, "", "", "Flat file to read Vault secret_id from. Can not be used in conjuction with Vault token flag.")
	viper.BindPFlag(conf.FLAG_VAULT_SECRET_ID_FILE, root.PersistentFlags().Lookup(conf.FLAG_VAULT_SECRET_ID_FILE))

	root.PersistentFlags().StringP(conf.FLAG_VAULT_MOUNT_PKI, "", conf.FLAG_VAULT_MOUNT_PKI_DEFAULT, "Path where the PKI secret engine is mounted.")
	viper.BindPFlag(conf.FLAG_VAULT_MOUNT_PKI, root.PersistentFlags().Lookup(conf.FLAG_VAULT_MOUNT_PKI))

	root.PersistentFlags().StringP(conf.FLAG_VAULT_MOUNT_APPROLE, "", conf.FLAG_VAULT_MOUNT_APPROLE_DEFAULT, "Path where the AppRole auth method is mounted.")
	viper.BindPFlag(conf.FLAG_VAULT_MOUNT_APPROLE, root.PersistentFlags().Lookup(conf.FLAG_VAULT_MOUNT_APPROLE))

	root.PersistentFlags().StringP(conf.FLAG_VAULT_PKI_BACKEND_ROLE, "", conf.FLAG_VAULT_PKI_BACKEND_ROLE_DEFAULT, "The name of the PKI role backend.")
	viper.BindPFlag(conf.FLAG_VAULT_PKI_BACKEND_ROLE, root.PersistentFlags().Lookup(conf.FLAG_VAULT_PKI_BACKEND_ROLE))

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.GetViper()

	v.SetConfigName(defaultConfigFilename)

	v.AddConfigPath("$HOME/.config/vault-pki-cli")
	v.AddConfigPath("/etc/vault-pki-cli/")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return nil
}

func readConfig(filepath string) error {
	viper.SetConfigFile(filepath)
	return viper.ReadInConfig()
}

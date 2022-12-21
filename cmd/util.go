package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/ilius/go-askpass"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"k8s.io/client-go/kubernetes"
)

func getVaultConfig(conf *conf.Config) *api.Config {
	vaultConfig := api.DefaultConfig()
	vaultConfig.MaxRetries = 13
	vaultConfig.Address = conf.VaultAddress
	return vaultConfig
}

func getKubernetesClient(conf *conf.Config) *kubernetes.Clientset {
	return nil
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

func handleFetchedData(certData []byte, config conf.Config) {
	// TODO: Fix!!

	/*
		if len(config.Backends) == 0 {
			fmt.Println(string(certData))
			os.Exit(0)
		}


		//sink, err := buildOutput(config)
		//sink.WriteSignature(certData)
		//if err != nil {
		//	log.Error().Msgf("Error writing cert: %v", err)
		//	os.Exit(1)
		//}

	*/
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

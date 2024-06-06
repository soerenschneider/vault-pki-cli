package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"golang.org/x/term"
)

type ZeroLogAdapter struct {
	logger *zerolog.Logger
}

func (z *ZeroLogAdapter) Error(format string, args ...interface{}) {
	z.logger.Error().Msgf(format, args...)
}

func (z *ZeroLogAdapter) Info(format string, args ...interface{}) {
	z.logger.Info().Msgf(format, args...)
}

func (z *ZeroLogAdapter) Debug(format string, args ...interface{}) {
	z.logger.Debug().Msgf(format, args...)
}

func (z *ZeroLogAdapter) Warn(format string, args ...interface{}) {
	z.logger.Debug().Msgf(format, args...)
}

func getVaultConfig(conf *conf.Config) *api.Config {
	vaultConfig := api.DefaultConfig()
	vaultConfig.MaxRetries = 5
	vaultConfig.Address = conf.VaultAddress
	return vaultConfig
}

func DieOnErr(err error, msg string, config ...*conf.Config) {
	if err == nil {
		return
	}

	if len(config) > 0 && len(config[0].MetricsFile) > 0 {
		internal.MetricSuccess.WithLabelValues(config[0].CommonName).Set(0)
		if err := internal.WriteMetrics(config[0].MetricsFile); err != nil {
			log.Error().Err(err).Msg("could not write metrics")
		}
	}

	log.Fatal().Err(err).Msg(msg)
}

func buildVaultClient(config *conf.Config) (*api.Client, error) {
	vaultConfig := getVaultConfig(config)
	vaultClient, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}
	vaultClient.SetLogger(&ZeroLogAdapter{logger: &log.Logger})

	return vaultClient, nil
}

func buildAuthImpl(conf *conf.Config) (vault.AuthMethod, error) {
	switch conf.VaultAuthMethod {
	case "token":
		log.Info().Msg("Building 'token' vault auth...")
		return vault.NewTokenAuth(conf.VaultToken)
	case "kubernetes":
		log.Info().Msg("Building 'kubernetes' vault auth...")
		return vault.NewVaultKubernetesAuth(conf.VaultAuthK8sRole)
	case "approle":
		approleData := make(map[string]string)
		approleData[vault.KeyRoleId] = conf.VaultRoleId
		approleData[vault.KeySecretId] = conf.VaultSecretId
		approleData[vault.KeySecretIdFile] = conf.VaultSecretIdFile

		log.Info().Msg("Building 'approle' vault auth...")
		return vault.NewAppRoleAuth(approleData, conf.VaultMountApprole)
	case "implicit":
		log.Info().Msg("Building 'implicit' vault auth...")
		return vault.NewTokenImplicitAuth(), nil
	}

	return nil, fmt.Errorf("unknown auth strategy '%s'", conf.VaultAuthMethod)
}

func PrintVersionInfo() {
	log.Info().Msgf("Version %s (%s)", internal.BuildVersion, internal.CommitHash)
}

func setupLogLevel(debug bool) {
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
}

func initLogging() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
		})
	}
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/vault"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"go.uber.org/multierr"
	"golang.org/x/net/context"
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

func buildAuthImpl(conf *conf.Config) (api.AuthMethod, error) {
	switch conf.VaultAuthMethod {
	case "kubernetes":
		log.Debug().Msg("Building 'kubernetes' vault auth...")
		return kubernetes.NewKubernetesAuth(conf.VaultAuthK8sRole)
	case "approle":
		log.Debug().Msg("Building 'approle' vault auth...")
		secretId := &approle.SecretID{}
		if len(conf.VaultSecretIdFile) > 0 {
			secretId.FromFile = conf.VaultSecretIdFile
		} else {
			secretId.FromString = conf.VaultSecretId
		}
		return approle.NewAppRoleAuth(conf.VaultRoleId, secretId)
	case "implicit":
		log.Debug().Msg("Building 'implicit' vault auth...")
		return vault.NewNoAuth(), nil
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
	//#nosec:G115
	if term.IsTerminal(int(os.Stdout.Fd())) {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
		})
	}
}

func runPostIssueHooks(ctx context.Context, config *conf.Config) error {
	if len(config.PostHooks) > 0 {
		log.Info().Msg("Running post-issue hooks")
	}

	var err error
	for _, hook := range config.PostHooks {
		log.Info().Msgf("Running command '%s'", hook)
		parsed := strings.Split(hook, " ")
		cmd := exec.CommandContext(ctx, parsed[0], parsed[1:]...) //#nosec G204
		cmdErr := cmd.Run()
		if cmdErr != nil {
			err = multierr.Append(err, fmt.Errorf("%w: cmd '%s' failed: %v", pkg.ErrRunHook, parsed[0], cmdErr))
		}
	}

	return err
}

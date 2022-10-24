package conf

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf/issue_sinks"
	"github.com/spf13/viper"
)

type Backend interface {
	PrintConfig()
	Validate() []error
	GetType() string
}

type IssueArguments struct {
	CommonName          string   `mapstructure:"common_name"`
	Ttl                 string   `mapstructure:"ttl"`
	IpSans              []string `mapstructure:"ip_sans"`
	AltNames            []string `mapstructure:"alt_names"`
	ForceNewCertificate bool
	BackendConfig       []issue_sinks.SinkConfig `mapstructure:"backend_config"`

	Backends []Backend

	PostIssueHooks                         []string `mapstructure:"post_hooks""`
	CertificateLifetimeThresholdPercentage float64  `mapstructure:"threshold"`
	MetricsFile                            string
}

func ParseFlags(config *Config) {
	if viper.IsSet(FLAG_ISSUE_COMMON_NAME) {
		config.VaultSecretIdFile = viper.GetString(FLAG_ISSUE_COMMON_NAME)
	}
	if viper.IsSet(FLAG_ISSUE_TTL) {
		config.VaultSecretIdFile = viper.GetString(FLAG_ISSUE_TTL)
	}
	if viper.IsSet(FLAG_ISSUE_IP_SANS) {
		config.VaultSecretIdFile = viper.GetString(FLAG_ISSUE_IP_SANS)
	}
	if viper.IsSet(FLAG_ISSUE_ALT_NAMES) {
		config.VaultSecretIdFile = viper.GetString(FLAG_ISSUE_ALT_NAMES)
	}
	config.ForceNewCertificate = viper.GetBool(FLAG_ISSUE_FORCE_NEW_CERTIFICATE)

	if viper.IsSet(FLAG_ISSUE_HOOKS) {
		config.PostIssueHooks = viper.GetStringSlice(FLAG_ISSUE_HOOKS)
	}

	if viper.IsSet(FLAG_ISSUE_METRICS_FILE) {
		config.IssueArguments.MetricsFile = viper.GetString(FLAG_ISSUE_METRICS_FILE)
	}
}

func (c *IssueArguments) BuildBackends() error {
	for _, conf := range c.BackendConfig {
		backend, err := buildBackend(conf)
		if err != nil {
			return fmt.Errorf("could not build backend: %v", err)
		}
		c.Backends = append(c.Backends, backend)
	}
	c.BackendConfig = []issue_sinks.SinkConfig{}
	return nil
}

func buildBackend(backend issue_sinks.SinkConfig) (Backend, error) {
	backendType, ok := backend["type"]
	if !ok {
		return nil, errors.New("configured backend doesn't have a 'type' field")
	}

	switch backendType {
	case issue_sinks.K8sType:
		return issue_sinks.K8sBackendFromMap(backend)
	case issue_sinks.FsType:
		return issue_sinks.FsBackendFromMap(backend)
	case issue_sinks.YubiType:
		return issue_sinks.YubiSinkFromMap(backend)
	default:
		return nil, fmt.Errorf("unknown type: %s", backendType)
	}
}

func (c *IssueArguments) UsesYubikey() bool {
	if internal.YubiKeySupport == "false" {
		return false
	}

	return false
}

func (c *IssueArguments) Validate() []error {
	errs := make([]error, 0)

	/*
		for _, backend := range c.Backends {
			errs = append(errs, backend.Validate()...)
		}

	*/

	if len(c.CommonName) == 0 {
		errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_COMMON_NAME))
	}

	if c.CertificateLifetimeThresholdPercentage < 5 || c.CertificateLifetimeThresholdPercentage > 90 {
		errs = append(errs, fmt.Errorf("'%s' must be [5, 90]", FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE))
	}

	/*
		if internal.YubiKeySupport == "true" {
			yubikeyConfigSupplied := c.YubikeySink.YubikeySlot != FLAG_ISSUE_YUBIKEY_SLOT_DEFAULT
			if len(c.FsBackends) == 0 && len(c.K8sBackends) == 0 && !yubikeyConfigSupplied {
				errs = append(errs, fmt.Errorf("must either provide '%s' or both '%s' and '%s'", FLAG_ISSUE_YUBIKEY_SLOT, FLAG_CERTIFICATE_FILE, FLAG_ISSUE_PRIVATE_KEY_FILE))
			}

			if (len(c.FsBackends) > 0 || len(c.K8sBackends) > 0) && yubikeyConfigSupplied {
				errs = append(errs, errors.New("can't provide yubi key slot AND file-based sink"))
			}

			if yubikeyConfigSupplied {
				err := pods.ValidateSlot(c.YubikeySink.YubikeySlot)
				if err != nil {
					errs = append(errs, fmt.Errorf("invalid yubikey slot '%d': %v", c.YubikeySlot, err))
				}
			}
		} else if len(c.FsBackends) == 0 && len(c.K8sBackends) == 0 {
			errs = append(errs, errors.New("no backend to store certificate provided"))
		}

	*/
	return errs
}

func (c *IssueArguments) PrintConfig() {
	log.Info().Msg("------------- Printing issue cmd values -------------")
	log.Info().Msgf("%s=%s", FLAG_ISSUE_TTL, c.Ttl)
	log.Info().Msgf("%s=%s", FLAG_ISSUE_COMMON_NAME, c.CommonName)
	log.Info().Msgf("%s=%s", FLAG_ISSUE_METRICS_FILE, c.MetricsFile)
	log.Info().Msgf("%s=%t", FLAG_ISSUE_FORCE_NEW_CERTIFICATE, c.ForceNewCertificate)
	log.Info().Msgf("%s=%f", FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE, c.CertificateLifetimeThresholdPercentage)
	if len(c.IpSans) > 0 {
		log.Info().Msgf("%s=%v", FLAG_ISSUE_IP_SANS, c.IpSans)
	}
	if len(c.AltNames) > 0 {
		log.Info().Msgf("%s=%v", FLAG_ISSUE_ALT_NAMES, c.AltNames)
	}
	for n, hook := range c.PostIssueHooks {
		log.Info().Msgf("%s[%d]='%s'", FLAG_ISSUE_HOOKS, n, hook)
	}
	log.Info().Msgf("------------- Finished printing config values -------------")
}

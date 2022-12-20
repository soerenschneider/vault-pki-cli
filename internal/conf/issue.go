package conf

import (
	"github.com/soerenschneider/vault-pki-cli/internal"
)

type Backend interface {
	PrintConfig()
	Validate() []error
	GetType() string
}

type IssueArguments struct {
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
					err := storage.ValidateSlot(c.YubikeySink.YubikeySlot)
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

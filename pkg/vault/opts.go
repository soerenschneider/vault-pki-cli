package vault

import "errors"

func WithPkiMount(mountPath string) VaultOpts {
	return func(c *VaultPki) error {
		if len(mountPath) == 0 {
			return errors.New("empty pki mount path")
		}
		c.pkiMountPath = mountPath
		return nil
	}
}

func WithKv2Mount(mountPath string) VaultOpts {
	return func(c *VaultPki) error {
		if len(mountPath) == 0 {
			return errors.New("empty kv2 mount path")
		}
		c.kv2MountPath = mountPath
		return nil
	}
}

func WithAcmePrefix(prefix string) VaultOpts {
	return func(c *VaultPki) error {
		if len(prefix) == 0 {
			return errors.New("empty acme prefix")
		}
		c.acmePrefix = prefix
		return nil
	}
}

package vault

import (
	"errors"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"strings"
)

type AuthMethod interface {
	Authenticate() (string, error)
	Cleanup() error
}

type VaultClient struct {
	client    *api.Client
	auth      AuthMethod
	roleName  string
	mountPath string
}

func NewVaultPki(client *api.Client, auth AuthMethod, config conf.Config) (*VaultClient, error) {
	if client == nil {
		return nil, errors.New("nil client passed")
	}

	if auth == nil {
		return nil, errors.New("nil auth passed")
	}

	return &VaultClient{
		client:    client,
		auth:      auth,
		mountPath: config.VaultMountPki,
		roleName:  config.VaultPkiRole,
	}, nil
}

func (c *VaultClient) Revoke(serial string) error {
	return nil
}

func (c *VaultClient) issue(opts conf.IssueArguments) (*api.Secret, error) {
	token, err := c.auth.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("could not authenticate: %v", err)
	}

	c.client.SetToken(token)

	path := fmt.Sprintf("%s/issue/%s", c.mountPath, c.roleName)
	data := buildIssueArgs(opts)

	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return nil, fmt.Errorf("could not issue certificate: %v", err)
	}

	return secret, nil
}

func buildIssueArgs(opts conf.IssueArguments) map[string]interface{} {
	data := map[string]interface{}{
		"common_name": opts.CommonName,
		"ttl":         opts.Ttl,
		"format":      "pem",
		"ip_sans":     strings.Join(opts.IpSans, ","),
		"alt_names":   strings.Join(opts.AltNames, ","),
	}

	return data
}

func (c *VaultClient) sign(csrFile pki.KeyPod, opts conf.SignArguments) (*api.Secret, error) {
	token, err := c.auth.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("could not authenticate: %v", err)
	}

	c.client.SetToken(token)

	path := fmt.Sprintf("%s/sign/%s", c.mountPath, c.roleName)
	data, err := buildSignArgs(csrFile, opts)
	if err != nil {
		return nil, fmt.Errorf("could not build request, reading csr file failed: %v", err)
	}

	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return nil, fmt.Errorf("could not issue certificate: %v", err)
	}

	return secret, nil
}

func buildSignArgs(csrFile pki.KeyPod, opts conf.SignArguments) (map[string]interface{}, error) {
	csr, err := csrFile.Read()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"csr":         string(csr),
		"common_name": opts.CommonName,
		"ttl":         opts.Ttl,
		"format":      "pem",
		"ip_sans":     strings.Join(opts.IpSans, ","),
		"alt_names":   strings.Join(opts.AltNames, ","),
	}

	return data, nil
}

func (c *VaultClient) Tidy() error {
	token, err := c.auth.Authenticate()
	if err != nil {
		return fmt.Errorf("could not authenticate: %v", err)
	}

	c.client.SetToken(token)

	path := fmt.Sprintf("%s/tidy", c.mountPath)

	data := map[string]interface{}{
		"tidy_cert_store":    true,
		"tidy_revoked_certs": true,
		"safety_buffer":      "90m",
	}
	_, err = c.client.Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("could not issue certificate: %v", err)
	}

	return nil
}

func (c *VaultClient) Sign(csrFile pki.KeyPod, opts conf.SignArguments) (*pki.Signature, error) {
	secret, err := c.sign(csrFile, opts)
	if err != nil {
		return nil, err
	}

	cert := fmt.Sprintf("%s", secret.Data["certificate"])
	chain := fmt.Sprintf("%s", secret.Data["ca_chain"])
	serial := fmt.Sprintf("%s", secret.Data["serial"])

	return &pki.Signature{
		Certificate: []byte(cert),
		CaChain:     []byte(chain),
		Serial:      serial,
	}, nil
}

func (c *VaultClient) Issue(opts conf.IssueArguments) (*pki.IssuedCert, error) {
	secret, err := c.issue(opts)
	if err != nil {
		return nil, err
	}

	privateKey := fmt.Sprintf("%s", secret.Data["private_key"])
	cert := fmt.Sprintf("%s", secret.Data["certificate"])
	chain := fmt.Sprintf("%s", secret.Data["ca_chain"])

	return &pki.IssuedCert{
		PrivateKey:  []byte(privateKey),
		Certificate: []byte(cert),
		CaChain:     []byte(chain),
	}, nil
}

func (c *VaultClient) Cleanup() error {
	return c.auth.Cleanup()
}

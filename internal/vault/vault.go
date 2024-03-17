package vault

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/internal"
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"golang.org/x/net/context"
)

const (
	// keys of the kv2 secret's map for the respective data
	acmevaultKeyPrivateKey  = "private_key"
	acmevaultKeyCertificate = "cert"
	acmevaultKeyIssuer      = "issuer"
	acmevaultVersion        = "version"

	// the secret name (without the path) of the certificate saved by acmevault
	acmevaultKv2SecretNameCertificate = "certificate"
	// the secret name (without the path) of the private key saved by acmevault
	acmevaultKv2SecretNamePrivatekey = "privatekey"
)

type AuthMethod interface {
	Login(ctx context.Context, client *api.Client) (*api.Secret, error)
	Cleanup(ctx context.Context, client *api.Client) error
}

type VaultClient struct {
	client    *api.Client
	auth      AuthMethod
	roleName  string
	mountPath string
	config    *conf.Config
}

func NewVaultPki(client *api.Client, auth AuthMethod, config *conf.Config) (*VaultClient, error) {
	if client == nil {
		return nil, errors.New("nil client passed")
	}

	if auth == nil {
		return nil, errors.New("nil auth passed")
	}

	if config == nil {
		return nil, errors.New("nil config passed")
	}

	return &VaultClient{
		client:    client,
		auth:      auth,
		mountPath: config.VaultMountPki,
		roleName:  config.VaultPkiRole,
		config:    config,
	}, nil
}

func (c *VaultClient) Revoke(serial string) error {
	_, err := c.client.Auth().Login(context.Background(), c.auth)
	if err != nil {
		return fmt.Errorf("could not authenticate: %v", err)
	}

	path := fmt.Sprintf("%s/revoke", c.mountPath)
	data := map[string]interface{}{
		"serial_number": serial,
	}

	resp, err := c.client.Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("could not revoke certificate: %v", err)
	}

	if resp != nil && len(resp.Warnings) > 0 {
		log.Warn().Msgf("revoking cert produced warning: %v", resp.Warnings)
	}

	return nil
}

func (c *VaultClient) issue(opts *conf.Config) (*api.Secret, error) {
	_, err := c.client.Auth().Login(context.Background(), c.auth)
	if err != nil {
		return nil, fmt.Errorf("could not authenticate: %v", err)
	}

	path := fmt.Sprintf("%s/issue/%s", c.mountPath, c.roleName)
	data := buildIssueArgs(opts)

	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return nil, fmt.Errorf("could not issue certificate: %v", err)
	}

	return secret, nil
}

func buildIssueArgs(opts *conf.Config) map[string]interface{} {
	data := map[string]interface{}{
		"common_name": opts.CommonName,
		"ttl":         opts.Ttl,
		"format":      "pem",
		"ip_sans":     strings.Join(opts.IpSans, ","),
		"alt_names":   strings.Join(opts.AltNames, ","),
	}

	return data
}

func (c *VaultClient) sign(csr string, opts *conf.Config) (*api.Secret, error) {
	_, err := c.client.Auth().Login(context.Background(), c.auth)
	if err != nil {
		return nil, fmt.Errorf("could not authenticate: %v", err)
	}

	path := fmt.Sprintf("%s/sign/%s", c.mountPath, c.roleName)
	data, err := buildSignArgs(csr, opts)
	if err != nil {
		return nil, fmt.Errorf("could not build request, reading csr file failed: %v", err)
	}

	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return nil, fmt.Errorf("could not issue certificate: %v", err)
	}

	return secret, nil
}

func buildSignArgs(csr string, opts *conf.Config) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"csr":         csr,
		"common_name": opts.CommonName,
		"ttl":         opts.Ttl,
		"format":      "pem",
		"ip_sans":     strings.Join(opts.IpSans, ","),
		"alt_names":   strings.Join(opts.AltNames, ","),
	}

	return data, nil
}

func (c *VaultClient) getAcmevaultDataPath(domain string, leaf string) string {
	prefix := fmt.Sprintf("%s/data/%s", c.config.VaultMountKv2, c.config.AcmePrefix)
	return fmt.Sprintf("%s/client/%s/%s", prefix, domain, leaf)
}

func (c *VaultClient) readKv2Secret(path string) (map[string]interface{}, error) {
	secret, err := c.client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("could not read kv2 data '%s': %w", path, err)
	}
	if secret == nil {
		return nil, errors.New("read kv2 data is nil")
	}

	var data map[string]interface{}
	_, ok := secret.Data["data"]
	if !ok {
		internal.MetricCertParseErrors.WithLabelValues(c.config.CommonName).Inc()
		return nil, errors.New("read kv2 secret contains no data")
	}
	data, ok = secret.Data["data"].(map[string]interface{})
	if !ok {
		internal.MetricCertParseErrors.WithLabelValues(c.config.CommonName).Inc()
		return nil, errors.New("read kv2 data is malformed")
	}

	return data, nil
}

func (c *VaultClient) readAcmeCert(commonName string) (*pki.CertData, error) {
	path := c.getAcmevaultDataPath(commonName, acmevaultKv2SecretNameCertificate)
	data, err := c.readKv2Secret(path)
	if err != nil {
		return nil, err
	}

	rawCert, ok := data[acmevaultKeyCertificate]
	if !ok {
		internal.MetricCertParseErrors.WithLabelValues(commonName).Inc()
		return nil, errors.New("read kv2 secret does not contain certificate data")
	}
	cert, err := base64.StdEncoding.DecodeString(rawCert.(string))
	if err != nil {
		internal.MetricCertParseErrors.WithLabelValues(commonName).Inc()
		return nil, errors.New("could not base64 decode cert")
	}
	cert = bytes.TrimRight(cert, "\n")

	var version string
	versionRaw, ok := data[acmevaultVersion]
	if ok {
		version = versionRaw.(string)
	}

	var issuer []byte
	if version == "v1" {
		rawIssuer, ok := data[acmevaultKeyIssuer]
		if ok {
			ca, err := base64.StdEncoding.DecodeString(rawIssuer.(string))
			if err == nil {
				issuer = bytes.TrimRight(ca, "\n")
				// TODO: remove support in future, this is apparently a bug in acmevault
				issuer = bytes.TrimLeft(issuer, "\n")
				// TODO end
			}
		}
	} else {
		// TODO: remove support in the future
		rawIssuer, ok := data["dummyIssuer"]
		if ok {
			ca, err := base64.StdEncoding.DecodeString(rawIssuer.(string))
			if err == nil {
				issuer = bytes.TrimRight(ca, "\n")
			}
		}
	}

	return &pki.CertData{Certificate: cert, CaData: issuer}, nil
}

func (c *VaultClient) readAcmeSecret(commonName string) (*pki.CertData, error) {
	path := c.getAcmevaultDataPath(commonName, acmevaultKv2SecretNamePrivatekey)
	data, err := c.readKv2Secret(path)
	if err != nil {
		return nil, err
	}

	rawKey, ok := data[acmevaultKeyPrivateKey]
	if !ok {
		internal.MetricCertParseErrors.WithLabelValues(commonName).Inc()
		return nil, errors.New("read kv2 secret does not contain private key data")
	}

	privateKey, err := base64.StdEncoding.DecodeString(rawKey.(string))
	if err != nil {
		internal.MetricCertParseErrors.WithLabelValues(commonName).Inc()
		return nil, errors.New("could not base64 decode key")
	}

	privateKey = bytes.TrimRight(privateKey, "\n")
	return &pki.CertData{PrivateKey: privateKey}, nil
}

func (c *VaultClient) ReadAcme(commonName string, conf *conf.Config) (*pki.CertData, error) {
	if conf == nil {
		return nil, errors.New("nil config provided")
	}

	_, err := c.client.Auth().Login(context.Background(), c.auth)
	if err != nil {
		return nil, fmt.Errorf("could not authenticate: %v", err)
	}

	certData, err := c.readAcmeCert(commonName)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate data: %w", err)
	}

	secretData, err := c.readAcmeSecret(commonName)
	if err != nil {
		return nil, fmt.Errorf("could not read secret data: %w", err)
	}

	return &pki.CertData{
		PrivateKey:  secretData.PrivateKey,
		Certificate: certData.Certificate,
		CaData:      certData.CaData,
	}, nil
}

func (c *VaultClient) Tidy() error {
	_, err := c.client.Auth().Login(context.Background(), c.auth)
	if err != nil {
		return fmt.Errorf("could not authenticate: %v", err)
	}

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

func (c *VaultClient) Sign(csr string, opts *conf.Config) (*pki.Signature, error) {
	if opts == nil {
		return nil, errors.New("empty config provided")
	}

	secret, err := c.sign(csr, opts)
	if err != nil {
		return nil, err
	}

	cert := fmt.Sprintf("%s", secret.Data["certificate"])
	chain := fmt.Sprintf("%s", secret.Data["issuing_ca"])
	serial := fmt.Sprintf("%s", secret.Data["serial_number"])

	return &pki.Signature{
		Certificate: []byte(cert),
		CaData:      []byte(chain),
		Serial:      serial,
	}, nil
}

func (c *VaultClient) Issue(opts *conf.Config) (*pki.CertData, error) {
	if opts == nil {
		return nil, errors.New("empty config provided")
	}

	secret, err := c.issue(opts)
	if err != nil {
		return nil, err
	}

	privateKey := fmt.Sprintf("%s", secret.Data["private_key"])
	cert := fmt.Sprintf("%s", secret.Data["certificate"])
	chain := fmt.Sprintf("%s", secret.Data["issuing_ca"])

	return &pki.CertData{
		PrivateKey:  []byte(privateKey),
		Certificate: []byte(cert),
		CaData:      []byte(chain),
	}, nil
}

func (c *VaultClient) Cleanup() error {
	return c.auth.Cleanup(context.Background(), c.client)
}

func (c *VaultClient) FetchCa(binary bool) ([]byte, error) {
	path := fmt.Sprintf("%s/ca", c.mountPath)
	if !binary {
		path = path + "/pem"
	}

	return c.readRaw(path)
}

func (c *VaultClient) FetchCaChain() ([]byte, error) {
	path := fmt.Sprintf("/%s/ca_chain", c.mountPath)
	return c.readRaw(path)
}

func (c *VaultClient) FetchCrl(binary bool) ([]byte, error) {
	path := fmt.Sprintf("%s/crl", c.mountPath)
	if !binary {
		path += "/pem"
	}

	return c.readRaw(path)
}

func (c *VaultClient) readRaw(path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	secret, err := c.client.Logical().ReadRawWithContext(ctx, path)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(secret.Body)
}

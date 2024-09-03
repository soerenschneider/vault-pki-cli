package vault

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/pkg"
	"go.uber.org/multierr"
	"golang.org/x/net/context"
)

const (
	defaultMountPki   = "pki"
	defaultMountKv2   = "kv2"
	defaultAcmePrefix = "acme"

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

type VaultClient interface {
	ReadWithContext(ctx context.Context, path string) (*api.Secret, error)
	WriteWithContext(ctx context.Context, path string, data map[string]any) (*api.Secret, error)
	ReadRawWithContext(ctx context.Context, path string) (*api.Response, error)
}

type VaultPki struct {
	client       VaultClient
	roleName     string
	pkiMountPath string
	kv2MountPath string
	acmePrefix   string
}

type VaultOpts func(client *VaultPki) error

func NewVaultPki(client VaultClient, roleName string, opts ...VaultOpts) (*VaultPki, error) {
	if client == nil {
		return nil, errors.New("nil client passed")
	}

	ret := &VaultPki{
		client:       client,
		roleName:     roleName,
		pkiMountPath: defaultMountPki,
		kv2MountPath: defaultMountKv2,
		acmePrefix:   defaultAcmePrefix,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(ret); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return ret, errs
}

func (c *VaultPki) Revoke(ctx context.Context, serial string) error {
	path := fmt.Sprintf("%s/revoke", c.pkiMountPath)
	data := map[string]interface{}{
		"serial_number": serial,
	}

	resp, err := c.client.WriteWithContext(ctx, path, data)
	if err != nil {
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
			return backoff.Permanent(err)
		}
		return fmt.Errorf("could not revoke certificate: %v", err)
	}

	if resp != nil && len(resp.Warnings) > 0 {
		log.Warn().Msgf("revoking cert produced warning: %v", resp.Warnings)
	}

	return nil
}

func (c *VaultPki) issue(ctx context.Context, args pkg.IssueArgs) (*api.Secret, error) {
	path := fmt.Sprintf("%s/issue/%s", c.pkiMountPath, c.roleName)
	data := buildIssueRequestArgs(args)

	secret, err := c.client.WriteWithContext(ctx, path, data)
	if err != nil {
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
			return nil, backoff.Permanent(err)
		}

		return nil, fmt.Errorf("could not issue certificate: %w", err)
	}

	return secret, nil
}

func buildIssueRequestArgs(args pkg.IssueArgs) map[string]any {
	data := map[string]any{
		"common_name": args.CommonName,
		"ttl":         args.Ttl,
		"format":      "pem",
		"ip_sans":     strings.Join(args.IpSans, ","),
		"alt_names":   strings.Join(args.AltNames, ","),
	}

	return data
}

func (c *VaultPki) sign(ctx context.Context, csr string, args pkg.SignatureArgs) (*api.Secret, error) {
	path := fmt.Sprintf("%s/sign/%s", c.pkiMountPath, c.roleName)
	data := buildSignArgs(csr, args)

	secret, err := c.client.WriteWithContext(ctx, path, data)
	if err != nil {
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
			return nil, backoff.Permanent(err)
		}
		return nil, fmt.Errorf("could not issue certificate: %v", err)
	}

	return secret, nil
}

func buildSignArgs(csr string, args pkg.SignatureArgs) map[string]interface{} {
	data := map[string]interface{}{
		"csr":         csr,
		"common_name": args.CommonName,
		"ttl":         args.Ttl,
		"format":      "pem",
		"ip_sans":     strings.Join(args.IpSans, ","),
		"alt_names":   strings.Join(args.AltNames, ","),
	}

	return data
}

func (c *VaultPki) getAcmevaultDataPath(domain string, leaf string) string {
	prefix := fmt.Sprintf("%s/data/%s", c.kv2MountPath, c.acmePrefix)
	return fmt.Sprintf("%s/client/%s/%s", prefix, domain, leaf)
}

func (c *VaultPki) readKv2Secret(ctx context.Context, path string) (map[string]interface{}, error) {
	secret, err := c.client.ReadWithContext(ctx, path)
	if err != nil {
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
			return nil, backoff.Permanent(err)
		}
		return nil, fmt.Errorf("could not read kv2 data '%s': %w", path, err)
	}
	if secret == nil {
		return nil, backoff.Permanent(errors.New("read kv2 data is nil"))
	}

	var data map[string]interface{}
	_, ok := secret.Data["data"]
	if !ok {
		return nil, backoff.Permanent(errors.New("read kv2 secret contains no data"))
	}
	data, ok = secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, backoff.Permanent(errors.New("read kv2 data is malformed"))
	}

	return data, nil
}

func (c *VaultPki) readAcmeCert(ctx context.Context, commonName string) (*pkg.CertData, error) {
	path := c.getAcmevaultDataPath(commonName, acmevaultKv2SecretNameCertificate)
	data, err := c.readKv2Secret(ctx, path)
	if err != nil {
		return nil, err
	}

	rawCert, ok := data[acmevaultKeyCertificate]
	if !ok {
		return nil, backoff.Permanent(errors.New("read kv2 secret does not contain certificate data"))
	}
	cert, err := base64.StdEncoding.DecodeString(rawCert.(string))
	if err != nil {
		return nil, backoff.Permanent(errors.New("could not base64 decode cert"))
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

	return &pkg.CertData{Certificate: cert, CaData: issuer}, nil
}

func (c *VaultPki) readAcmeSecret(ctx context.Context, commonName string) (*pkg.CertData, error) {
	path := c.getAcmevaultDataPath(commonName, acmevaultKv2SecretNamePrivatekey)
	data, err := c.readKv2Secret(ctx, path)
	if err != nil {
		return nil, err
	}

	rawKey, ok := data[acmevaultKeyPrivateKey]
	if !ok {
		return nil, backoff.Permanent(errors.New("read kv2 secret does not contain private key data"))
	}

	privateKey, err := base64.StdEncoding.DecodeString(rawKey.(string))
	if err != nil {
		return nil, backoff.Permanent(errors.New("could not base64 decode key"))
	}

	privateKey = bytes.TrimRight(privateKey, "\n")
	return &pkg.CertData{PrivateKey: privateKey}, nil
}

func (c *VaultPki) ReadAcme(ctx context.Context, commonName string) (*pkg.CertData, error) {
	certData, err := c.readAcmeCert(ctx, commonName)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate data: %w", err)
	}

	secretData, err := c.readAcmeSecret(ctx, commonName)
	if err != nil {
		return nil, fmt.Errorf("could not read secret data: %w", err)
	}

	return &pkg.CertData{
		PrivateKey:  secretData.PrivateKey,
		Certificate: certData.Certificate,
		CaData:      certData.CaData,
	}, nil
}

func (c *VaultPki) Tidy(ctx context.Context) error {
	path := fmt.Sprintf("%s/tidy", c.pkiMountPath)
	data := map[string]interface{}{
		"tidy_cert_store":    true,
		"tidy_revoked_certs": true,
		"safety_buffer":      "90m",
	}
	_, err := c.client.WriteWithContext(ctx, path, data)
	if err != nil {
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
			return backoff.Permanent(err)
		}

		return err
	}

	return nil
}

func (c *VaultPki) Sign(ctx context.Context, csr string, args pkg.SignatureArgs) (*pkg.Signature, error) {
	secret, err := c.sign(ctx, csr, args)
	if err != nil {
		return nil, err
	}

	cert := fmt.Sprintf("%s", secret.Data["certificate"])
	chain := fmt.Sprintf("%s", secret.Data["issuing_ca"])
	serial := fmt.Sprintf("%s", secret.Data["serial_number"])

	return &pkg.Signature{
		Certificate: []byte(cert),
		CaData:      []byte(chain),
		Serial:      serial,
	}, nil
}

func (c *VaultPki) Issue(ctx context.Context, args pkg.IssueArgs) (*pkg.CertData, error) {
	secret, err := c.issue(ctx, args)
	if err != nil {
		return nil, err
	}

	privateKey := fmt.Sprintf("%s", secret.Data["private_key"])
	cert := fmt.Sprintf("%s", secret.Data["certificate"])
	chain := fmt.Sprintf("%s", secret.Data["issuing_ca"])

	return &pkg.CertData{
		PrivateKey:  []byte(privateKey),
		Certificate: []byte(cert),
		CaData:      []byte(chain),
	}, nil
}

func (c *VaultPki) FetchCa(binary bool) ([]byte, error) {
	path := fmt.Sprintf("%s/ca", c.pkiMountPath)
	if !binary {
		path = path + "/pem"
	}

	return c.readRaw(path)
}

func (c *VaultPki) FetchCaChain() ([]byte, error) {
	path := fmt.Sprintf("/%s/ca_chain", c.pkiMountPath)
	return c.readRaw(path)
}

func (c *VaultPki) FetchCrl(binary bool) ([]byte, error) {
	path := fmt.Sprintf("%s/crl", c.pkiMountPath)
	if !binary {
		path += "/pem"
	}

	return c.readRaw(path)
}

func (c *VaultPki) readRaw(path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	secret, err := c.client.ReadRawWithContext(ctx, path)
	if err != nil {
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && !shouldRetry(respErr.StatusCode) {
			return nil, backoff.Permanent(err)
		}

		return nil, err
	}

	return io.ReadAll(secret.Body)
}

func shouldRetry(statusCode int) bool {
	switch statusCode {
	case 400, // Bad Request
		401, // Unauthorized
		403, // Forbidden
		404, // Not Found
		405, // Method Not Allowed
		406, // Not Acceptable
		407, // Proxy Authentication Required
		409, // Conflict
		410, // Gone
		411, // Length Required
		412, // Precondition Failed
		413, // Payload Too Large
		414, // URI Too Long
		415, // Unsupported Media Type
		416, // Range Not Satisfiable
		417, // Expectation Failed
		418, // I'm a Teapot
		421, // Misdirected Request
		422, // Unprocessable Entity
		423, // Locked (WebDAV)
		424, // Failed Dependency (WebDAV)
		425, // Too Early
		426, // Upgrade Required
		428, // Precondition Required
		429, // Too Many Requests
		431, // Request Header Fields Too Large
		451: // Unavailable For Legal Reasons
		return false
	default:
		return true
	}
}

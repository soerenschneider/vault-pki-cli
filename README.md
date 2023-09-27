# vault-pki-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/vault-pki-cli)](https://goreportcard.com/report/github.com/soerenschneider/vault-pki-cli)
![test-workflow](https://github.com/soerenschneider/vault-pki-cli/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/vault-pki-cli/actions/workflows/release-container.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/vault-pki-cli/actions/workflows/golangci-lint.yaml/badge.svg)

## Features

🔐 Issues, signs and revokes x509 certificates

🔑 Reads ACME certs written by [acmevault](https://github.com/soerenschneider/acmevault) (e.g. issued by LetsEncrypt)

⛓  Reads the CA / CA chain of a PKI

📖 Reads the CRL of a PKI

📝 Supports DER and PEM formats, 

⏰ Automatically renews certificates based on its lifetime

🛂 Authenticate against Vault using Kubernetes, AppRole, (explicit) token or _implicit_ auth

🗂 Supports multiple _sinks_: Kubernetes, plain files, in-memory

💻 Runs effortlessly both on your workstation's CLI via command line flags or automated via systemd and config files on your server

🔭 Provides metrics to increase observability for robust automation

## Why would I need this?

mTLS is a strong and proven authentication mechanism and vault-pki-cli deals with some of its challenges

| mTLS challenges            | How vault-pki-cli can help                                                                                             |
|----------------------------|------------------------------------------------------------------------------------------------------------------------|
| Certificate Management     | Dramatically removes complexity for issuing, renewing, and revoking certificates and downloading CRLs                  |
| Key Distribution           | Safely distributes certificates using Vault's API                                                                      |
| Revocation Challenges      | Revocation is easy and can be performed automatically                                                                  |
| Key Storage                | Observability and automation allows for short-lived certificates to limit the blast-radius of compromised certificates |
| Certificate Expiration     | Unless Vault is down, certificates are automatically renewed after a user-defined threshold passes                     |


## Installation

### Docker / Podman
````shell
$ docker run ghcr.io/soerenschneider/vault-pki-cli:main
````

### Binaries
Head over to the [prebuilt binaries](https://github.com/soerenschneider/vault-pki-cli/releases) and download the correct binary for your system.

### From Source
As a prerequesite, you need to have [Golang SDK](https://go.dev/dl/) installed. After that, you can install vault-pki-cli from source by invoking:
```text
$ go install github.com/soerenschneider/vault-pki-cli@latest
```
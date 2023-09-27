# vault-pki-cli
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/vault-pki-cli)](https://goreportcard.com/report/github.com/soerenschneider/vault-pki-cli)
![test-workflow](https://github.com/soerenschneider/vault-pki-cli/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/vault-pki-cli/actions/workflows/release-container.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/vault-pki-cli/actions/workflows/golangci-lint.yaml/badge.svg)

## Features

🔐 Issues x509 certificates

🔏 Signs x509 certificates

⛔️ Revokes x509 certificates

🔑 Reads ACME certs written by [acmevault](https://github.com/soerenschneider/acmevault)

⛓ Reads the CA / CA chain of a PKI

📖 Reads the CRL of a PKI

🛂 Authenticate against Vault using AppRole, (explicit) token or implicit__ auth

⏰ Automatically renews certificates based on its lifetime

💻 Runs effortlessly both on your workstation's CLI via command line flags or automated via systemd and config files on your server

🔭 Provides metrics to increase observability for robust automation

## Installation

### Docker / Podman
````shell
$ docker pull ghcr.io/soerenschneider/vault-pki-cli:main
````

### Binaries
Head over to the [prebuilt binaries](https://github.com/soerenschneider/vault-pki-cli/releases) and download the correct binary for your system.

### From Source
As a prerequesite, you need to have [Golang SDK](https://go.dev/dl/) installed. After that, you can install vault-pki-cli from source by invoking:
```text
$ go install github.com/soerenschneider/vault-pki-cli@latest
```
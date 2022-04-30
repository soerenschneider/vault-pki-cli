# Changelog

### [1.3.1](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.3.0...v1.3.1) (2022-04-30)


### Bug Fixes

* metric 'cert_lifetime_seconds_total' was never updated ([253bee9](https://www.github.com/soerenschneider/vault-pki-cli/commit/253bee974cf16ed21e641a6aad97b9431c1d7920))

## [1.3.0](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.2.2...v1.3.0) (2022-04-30)


### Features

* make file owner configurable ([4520a7e](https://www.github.com/soerenschneider/vault-pki-cli/commit/4520a7ed744485dd6787ac35220c345132eac74d))

### [1.2.2](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.2.1...v1.2.2) (2022-04-30)


### Bug Fixes

* fix wrong usage of pem.Decode() ([6e7c41d](https://www.github.com/soerenschneider/vault-pki-cli/commit/6e7c41dbeddadca0e7f628fc50776c95ae043032))

### [1.2.1](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.2.0...v1.2.1) (2022-04-30)


### Bug Fixes

* fix default location of metrics file ([c6f110d](https://www.github.com/soerenschneider/vault-pki-cli/commit/c6f110d72f7097e95560511c172a4c154a9c635b))

## [1.2.0](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.1.0...v1.2.0) (2022-04-29)


### Features

* add 'fetch crl' subcommand ([8ce4eac](https://www.github.com/soerenschneider/vault-pki-cli/commit/8ce4eac09cb6ff355edc84483f38a6c0e8ec3bb3))
* new 'fetch crl' subcommand ([dcac789](https://www.github.com/soerenschneider/vault-pki-cli/commit/dcac78972a7a189596deebf9c9a5b8f0988b40ac))


### Bug Fixes

* add missing 'read ca chain' subcommand ([41d8b73](https://www.github.com/soerenschneider/vault-pki-cli/commit/41d8b73fdf3e9dd27b367e645a217a898ff5b0d6))
* error on invalid cert data ([1b0fba6](https://www.github.com/soerenschneider/vault-pki-cli/commit/1b0fba63a5735fb91a97cc88a336f662ca74a8c5))

## [1.1.0](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.0.0...v1.1.0) (2022-03-17)


### Features

* add subcmd to read ca chain from vault ([2152e97](https://www.github.com/soerenschneider/vault-pki-cli/commit/2152e97c83a9a4f9df9680178e26903809b9fbdb))
* add subcommand to read ca from vault ([8393e98](https://www.github.com/soerenschneider/vault-pki-cli/commit/8393e98004ec6ec5f69a2c64bf802cb2b2e3a91a))

## 1.0.0 (2022-02-24)


### Features

* revoke token on exit ([35a4458](https://www.github.com/soerenschneider/vault-pki-cli/commit/35a445868c50e726a6fed96cb54bb507f9bc4b0a))

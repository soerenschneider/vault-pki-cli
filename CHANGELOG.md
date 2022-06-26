# Changelog

## [1.6.2](https://github.com/soerenschneider/vault-pki-cli/compare/v1.6.1...v1.6.2) (2022-06-20)


### Bug Fixes

* fix unhandled error ([07107bb](https://github.com/soerenschneider/vault-pki-cli/commit/07107bbdbba0e0cfcc1434958835fd18ad310643))

## [1.6.1](https://github.com/soerenschneider/vault-pki-cli/compare/v1.6.0...v1.6.1) (2022-06-07)


### Bug Fixes

* Signal yubikey support via ldflags ([#33](https://github.com/soerenschneider/vault-pki-cli/issues/33)) ([19fa2f5](https://github.com/soerenschneider/vault-pki-cli/commit/19fa2f5dc41bee559fecc4e8037c001dd73bd2ba))

## [1.6.0](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.5.1...v1.6.0) (2022-06-06)


### Features

* Add support for post-issue hooks ([#28](https://www.github.com/soerenschneider/vault-pki-cli/issues/28)) ([9be4615](https://www.github.com/soerenschneider/vault-pki-cli/commit/9be4615d25ca2f921c2b58f7511b143943c72f9b))
* write to multiple backends ([#26](https://www.github.com/soerenschneider/vault-pki-cli/issues/26)) ([ba317fd](https://www.github.com/soerenschneider/vault-pki-cli/commit/ba317fd7639f4379d964924062389a32264a13b1))


### Bug Fixes

* Fix failing test ([d8aa65b](https://www.github.com/soerenschneider/vault-pki-cli/commit/d8aa65bde5df505dfb8e0feaeb973068abb3f9ad))

### [1.5.1](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.5.0...v1.5.1) (2022-05-24)


### Bug Fixes

* respect 'ca-file' option ([3d4aea2](https://www.github.com/soerenschneider/vault-pki-cli/commit/3d4aea27c777622a6a1dafb22b2128d017fce52b))

## [1.5.0](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.4.0...v1.5.0) (2022-05-18)


### Features

* add yubikey support ([#20](https://www.github.com/soerenschneider/vault-pki-cli/issues/20)) ([d85df82](https://www.github.com/soerenschneider/vault-pki-cli/commit/d85df823987dddd425ab06753331c1c088d4258a))

## [1.4.0](https://www.github.com/soerenschneider/vault-pki-cli/compare/v1.3.1...v1.4.0) (2022-05-07)


### Miscellaneous Chores

* release 1.4.0 ([9c3a391](https://www.github.com/soerenschneider/vault-pki-cli/commit/9c3a3919943c4cb71e991f6736c738806b74a7d3))

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

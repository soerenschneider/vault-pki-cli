# Changelog

## [1.13.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.12.0...v1.13.0) (2023-07-03)


### Features

* drop clunky support for yubikeys (use CSRs instead) ([96cf92c](https://github.com/soerenschneider/vault-pki-cli/commit/96cf92c0c993e46dbebc490ed6d28f5a556d8447))
* drop support for openbsd ([850c547](https://github.com/soerenschneider/vault-pki-cli/commit/850c547f46937e50197b1f10572d27ed4b04b10c))


### Bug Fixes

* fix error handling logic ([7338ef1](https://github.com/soerenschneider/vault-pki-cli/commit/7338ef1646cbf889e8ad31a9f8cc33a370e0fef1))

## [1.12.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.11.2...v1.12.0) (2023-07-03)


### Features

* verify that cert on disk belongs to ca ([047f804](https://github.com/soerenschneider/vault-pki-cli/commit/047f804ede7e3557e8183aa34f664ba8ee2eec8c))


### Bug Fixes

* add methods to interface ([0e27785](https://github.com/soerenschneider/vault-pki-cli/commit/0e2778532929a57094039815c311967ed264ff2b))
* check err value ([f658a47](https://github.com/soerenschneider/vault-pki-cli/commit/f658a47209cc337cfb15cd80a7a4d9a50e4ccbd8))
* check errors ([dd7f235](https://github.com/soerenschneider/vault-pki-cli/commit/dd7f235a1b05136b16c6c6cd03120d716a779211))
* fix catching the correct signal ([8ae574b](https://github.com/soerenschneider/vault-pki-cli/commit/8ae574b787bd3bb7441d3de746a98e0debb3f80c))
* fix logic ([decc702](https://github.com/soerenschneider/vault-pki-cli/commit/decc70202de4ef4ece537ec2a2c707543ac20f48))
* fix tag format ([19f9744](https://github.com/soerenschneider/vault-pki-cli/commit/19f97440e33f2189e99c474936a8c5cbb9c9e4a3))
* handle errors ([7915450](https://github.com/soerenschneider/vault-pki-cli/commit/7915450dcac44c454aaf875626d07f5c3f50c00b))
* return after work ([a57cb01](https://github.com/soerenschneider/vault-pki-cli/commit/a57cb01050f0b5a949d3ec31c0b60992cb452b38))
* set timeouts for http server ([2bd9d6f](https://github.com/soerenschneider/vault-pki-cli/commit/2bd9d6fe1fe196cd92a4084b93ff828d345acec0))

## [1.11.2](https://github.com/soerenschneider/vault-pki-cli/compare/v1.11.1...v1.11.2) (2023-04-19)


### Bug Fixes

* expand config path ([e4381fe](https://github.com/soerenschneider/vault-pki-cli/commit/e4381fee318453387fac39cd634e3892f8b3af66))
* respect '--force-new-certificate' flag ([882f83f](https://github.com/soerenschneider/vault-pki-cli/commit/882f83fb524c5aaec22250122c49929a6e9b9c6e))

## [1.11.1](https://github.com/soerenschneider/vault-pki-cli/compare/v1.11.0...v1.11.1) (2023-02-14)


### Bug Fixes

* fix implicit auth ([6126f4b](https://github.com/soerenschneider/vault-pki-cli/commit/6126f4b6ac5d307c4913825e829aa407074dd503))
* make expanding user dirs work ([da334fe](https://github.com/soerenschneider/vault-pki-cli/commit/da334fe3a47e7257290bf5d7d2b35e2ffb8f53f7))

## [1.11.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.10.4...v1.11.0) (2023-02-03)


### Features

* Add implicit auth mechanism ([e406c4d](https://github.com/soerenschneider/vault-pki-cli/commit/e406c4da6137e1b8051753cf2e8b2154f301555b))

## [1.10.4](https://github.com/soerenschneider/vault-pki-cli/compare/v1.10.3...v1.10.4) (2023-01-17)


### Miscellaneous Chores

* release 1.10.4 ([167444a](https://github.com/soerenschneider/vault-pki-cli/commit/167444a465286db708a0beea8ac1a5b9281fc6bc))

## [1.10.3](https://github.com/soerenschneider/vault-pki-cli/compare/v1.10.2...v1.10.3) (2023-01-14)


### Bug Fixes

* fix duplicated short-hand flag ([30231a5](https://github.com/soerenschneider/vault-pki-cli/commit/30231a5490badf1058fcb126b9cb5445967ea202))

## [1.10.2](https://github.com/soerenschneider/vault-pki-cli/compare/v1.10.1...v1.10.2) (2023-01-09)


### Bug Fixes

* don't read from nil cert storage ([134a0ad](https://github.com/soerenschneider/vault-pki-cli/commit/134a0adfedb5b4ad376ae6f02d67cbc8c9b7a85b))
* fix msgf directives ([f73c5bc](https://github.com/soerenschneider/vault-pki-cli/commit/f73c5bc53716cba25a7ec82e65475019353d4156))

## [1.10.1](https://github.com/soerenschneider/vault-pki-cli/compare/v1.10.0...v1.10.1) (2023-01-09)


### Bug Fixes

* re-arrange cert data when writing to single file ([05d6766](https://github.com/soerenschneider/vault-pki-cli/commit/05d676677c8e44a8aadfb3aaadc7337f005fa328))

## [1.10.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.9.1...v1.10.0) (2023-01-09)


### Features

* allow writing all cert data to single file ([e1bba82](https://github.com/soerenschneider/vault-pki-cli/commit/e1bba82e489fe125ca96af2ed7daced3e1cdceb9))

## [1.9.1](https://github.com/soerenschneider/vault-pki-cli/compare/v1.9.0...v1.9.1) (2023-01-07)


### Bug Fixes

* fix uid and gid parsing ([976d1e9](https://github.com/soerenschneider/vault-pki-cli/commit/976d1e93b5b36ac524770ae0caecc834341a7696))

## [1.9.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.8.0...v1.9.0) (2023-01-06)


### Features

* add 'chmod' param ([47124ee](https://github.com/soerenschneider/vault-pki-cli/commit/47124eeee4ec2b0833cfd66f653a72259af269d3))
* Enable reading certificates provided by Acmevault ([0750190](https://github.com/soerenschneider/vault-pki-cli/commit/075019001a26f73256fb0da9414e9c906404b0e3))


### Bug Fixes

* don't print sensitive values ([6e6eb86](https://github.com/soerenschneider/vault-pki-cli/commit/6e6eb86c3180f59d9e2fcd132c488fb7208aaca3))
* prevent writing 'go_' prefixed metrics to file ([de8a67b](https://github.com/soerenschneider/vault-pki-cli/commit/de8a67b2eab09aa5b445f880165ca7af649e0924))

## [1.8.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.7.0...v1.8.0) (2022-12-29)


### Features

* Add k8s auth for vault ([#86](https://github.com/soerenschneider/vault-pki-cli/issues/86)) ([ad93140](https://github.com/soerenschneider/vault-pki-cli/commit/ad93140ee4f37ae63c2b4a779b9e9f994ece9c31))
* add multi sink for keypairs ([#82](https://github.com/soerenschneider/vault-pki-cli/issues/82)) ([6416e0f](https://github.com/soerenschneider/vault-pki-cli/commit/6416e0fb90b5145903af945dc9440de9d5d23020))
* initial support for k8s backend ([#80](https://github.com/soerenschneider/vault-pki-cli/issues/80)) ([6aa4292](https://github.com/soerenschneider/vault-pki-cli/commit/6aa429290c6c9557327258a93a6cea602b32f307))
* support for running as daemon ([9d9a939](https://github.com/soerenschneider/vault-pki-cli/commit/9d9a939c2941b8a076164a6e30701d5df3d0b773))

## [1.7.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.6.3...v1.7.0) (2022-12-07)


### Features

* Implement revoke operation ([#71](https://github.com/soerenschneider/vault-pki-cli/issues/71)) ([15fd650](https://github.com/soerenschneider/vault-pki-cli/commit/15fd6500280ce0e35df4cce13dc6588e9ff0aebb))

## [1.6.3](https://github.com/soerenschneider/vault-pki-cli/compare/v1.6.2...v1.6.3) (2022-07-19)


### Bug Fixes

* Fix yubikey backend ([#39](https://github.com/soerenschneider/vault-pki-cli/issues/39)) ([ceb1234](https://github.com/soerenschneider/vault-pki-cli/commit/ceb12347fbe3ec60b656abe6b6eba17288618a55))

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

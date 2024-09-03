# Changelog

## [1.15.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.14.0...v1.15.0) (2024-09-03)


### Features

* add custom validation for ttl ([5c45898](https://github.com/soerenschneider/vault-pki-cli/commit/5c45898c6877ac16c006ab71d68e61f8a9a57f7c))
* make max amount of retries configurable ([e5244b7](https://github.com/soerenschneider/vault-pki-cli/commit/e5244b7e96fa35c18d8f2a54a547e2b6dad6c4ba))


### Bug Fixes

* **deps:** bump github.com/cenkalti/backoff/v3 from 3.0.0 to 3.2.2 ([7fd44f0](https://github.com/soerenschneider/vault-pki-cli/commit/7fd44f08a1a257696e0ceb0bd8774252fcde0101))
* **deps:** bump github.com/hashicorp/go-retryablehttp ([ff45dea](https://github.com/soerenschneider/vault-pki-cli/commit/ff45dea046bf830137a857e58996ef74e92dcd0e))
* **deps:** bump github.com/spf13/cobra from 1.8.0 to 1.8.1 ([d08bd03](https://github.com/soerenschneider/vault-pki-cli/commit/d08bd03299e392accee53d5cfd7a262d83dcd772))
* **deps:** bump golang from 1.22.3 to 1.23.0 ([2ab7cfd](https://github.com/soerenschneider/vault-pki-cli/commit/2ab7cfdefcef9accd58d6ac37c72238c788d61c7))
* **deps:** bump golang.org/x/net from 0.25.0 to 0.28.0 ([2e21493](https://github.com/soerenschneider/vault-pki-cli/commit/2e21493ae6b8472a41488c89025f85025622b422))
* **deps:** bump golang.org/x/sys from 0.20.0 to 0.24.0 ([402e80c](https://github.com/soerenschneider/vault-pki-cli/commit/402e80c5ba6df19a4744bf677703700d7de864cd))
* don't validate ttl if not specified ([2a3a4c3](https://github.com/soerenschneider/vault-pki-cli/commit/2a3a4c300d71de752e93eea168f2cca82161e3cb))

## [1.15.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.14.0...v1.15.0) (2024-09-03)


### Features

* add custom validation for ttl ([5c45898](https://github.com/soerenschneider/vault-pki-cli/commit/5c45898c6877ac16c006ab71d68e61f8a9a57f7c))
* make max amount of retries configurable ([e5244b7](https://github.com/soerenschneider/vault-pki-cli/commit/e5244b7e96fa35c18d8f2a54a547e2b6dad6c4ba))


### Bug Fixes

* **deps:** bump github.com/cenkalti/backoff/v3 from 3.0.0 to 3.2.2 ([7fd44f0](https://github.com/soerenschneider/vault-pki-cli/commit/7fd44f08a1a257696e0ceb0bd8774252fcde0101))
* **deps:** bump github.com/hashicorp/go-retryablehttp ([ff45dea](https://github.com/soerenschneider/vault-pki-cli/commit/ff45dea046bf830137a857e58996ef74e92dcd0e))
* **deps:** bump github.com/spf13/cobra from 1.8.0 to 1.8.1 ([d08bd03](https://github.com/soerenschneider/vault-pki-cli/commit/d08bd03299e392accee53d5cfd7a262d83dcd772))
* **deps:** bump golang from 1.22.3 to 1.23.0 ([2ab7cfd](https://github.com/soerenschneider/vault-pki-cli/commit/2ab7cfdefcef9accd58d6ac37c72238c788d61c7))
* **deps:** bump golang.org/x/net from 0.25.0 to 0.28.0 ([2e21493](https://github.com/soerenschneider/vault-pki-cli/commit/2e21493ae6b8472a41488c89025f85025622b422))
* don't validate ttl if not specified ([2a3a4c3](https://github.com/soerenschneider/vault-pki-cli/commit/2a3a4c300d71de752e93eea168f2cca82161e3cb))

## [1.14.0](https://github.com/soerenschneider/vault-pki-cli/compare/v1.13.5...v1.14.0) (2024-06-08)


### Features

* add flag for controlling printing debug logs ([2b9ba62](https://github.com/soerenschneider/vault-pki-cli/commit/2b9ba62a68559687b400af31f0fefb6d5dc6b6f2))
* increase resilience by using backoff mechanism ([1f6879e](https://github.com/soerenschneider/vault-pki-cli/commit/1f6879ee1214af818bc231d6718df6d20ee3ae81))


### Bug Fixes

* **deps:** bump github.com/go-playground/validator/v10 ([2e6d797](https://github.com/soerenschneider/vault-pki-cli/commit/2e6d79754c8cfc8e2277e969614588a73c88c339))
* **deps:** bump github.com/hashicorp/vault/api/auth/kubernetes ([696cded](https://github.com/soerenschneider/vault-pki-cli/commit/696cdeda0d3176b9d7c2a8850af14dde5ba5cf4b))
* **deps:** bump github.com/prometheus/client_golang ([5b96e58](https://github.com/soerenschneider/vault-pki-cli/commit/5b96e58092446c3d10c3a68b806ffa15311af220))
* **deps:** bump github.com/prometheus/common from 0.50.0 to 0.54.0 ([1a58b98](https://github.com/soerenschneider/vault-pki-cli/commit/1a58b985f85b8759417c6dd90b157ce2ded4400c))
* **deps:** bump github.com/rs/zerolog from 1.31.0 to 1.33.0 ([4381a7b](https://github.com/soerenschneider/vault-pki-cli/commit/4381a7b19c02578505b32ce21b57420fa1d1ab18))
* **deps:** bump github.com/spf13/viper from 1.17.0 to 1.18.2 ([bb5bbbf](https://github.com/soerenschneider/vault-pki-cli/commit/bb5bbbf6a99a1af1bf1b605bfd02d88ea34810f8))
* **deps:** bump github.com/spf13/viper from 1.18.2 to 1.19.0 ([28fc787](https://github.com/soerenschneider/vault-pki-cli/commit/28fc787de0639a4065321e49b38133264e05a842))
* **deps:** bump golang from 1.22.1 to 1.22.3 ([1796e46](https://github.com/soerenschneider/vault-pki-cli/commit/1796e4683cbfb0118ad0d0d2f4781f3f8a8f89a9))
* **deps:** bump k8s.io/api from 0.29.3 to 0.30.1 ([a2fe1d6](https://github.com/soerenschneider/vault-pki-cli/commit/a2fe1d6c95946c54d0c5da8799c60866f4e024f1))
* **deps:** bump k8s.io/apimachinery from 0.29.3 to 0.30.1 ([50882d1](https://github.com/soerenschneider/vault-pki-cli/commit/50882d1974fc1b640fb8e9cd693e1debce7e0bad))
* **deps:** bump k8s.io/client-go from 0.29.3 to 0.30.1 ([c52f827](https://github.com/soerenschneider/vault-pki-cli/commit/c52f8271eea0111038ae9af7f9f42f8506748de2))
* do not try to revoke expired certificates ([16f1a80](https://github.com/soerenschneider/vault-pki-cli/commit/16f1a8035214165e4fdd556b39282afab47cdc68))
* fix typos in cmd's help messages ([d92d10c](https://github.com/soerenschneider/vault-pki-cli/commit/d92d10ce9cee2e4ae76eaa563958c5641a967df1))
* mark appropriate errors as permanent ([9203aaa](https://github.com/soerenschneider/vault-pki-cli/commit/9203aaa1b61878635f982e1480b3fb873443d50c))
* only try to write metrics when metrics file is passed ([c821b78](https://github.com/soerenschneider/vault-pki-cli/commit/c821b78e457a335d0c5756466b9e9bddf9ca4713))

## [1.13.5](https://github.com/soerenschneider/vault-pki-cli/compare/v1.13.4...v1.13.5) (2024-03-21)


### Bug Fixes

* **deps:** bump github.com/go-jose/go-jose/v3 from 3.0.0 to 3.0.1 ([83fe417](https://github.com/soerenschneider/vault-pki-cli/commit/83fe417e7ef8a86b6af00595b15bf4e720b9a0b8))
* **deps:** bump github.com/go-jose/go-jose/v3 from 3.0.1 to 3.0.3 ([4edd3d0](https://github.com/soerenschneider/vault-pki-cli/commit/4edd3d02dd9dfb1ed4ebf7c1db7b6f258f84854e))
* **deps:** bump github.com/go-playground/validator/v10 ([6d9d3a6](https://github.com/soerenschneider/vault-pki-cli/commit/6d9d3a6a9323122d76d568c96f5dc905f7a1339e))
* **deps:** bump github.com/go-playground/validator/v10 ([b75fa6c](https://github.com/soerenschneider/vault-pki-cli/commit/b75fa6c41955ca32297d9a52d9fa2e3c4e88319b))
* **deps:** bump github.com/hashicorp/vault/api from 1.10.0 to 1.12.2 ([e00361f](https://github.com/soerenschneider/vault-pki-cli/commit/e00361f20ebe7812fc33b0f2937852222ddbdf09))
* **deps:** bump github.com/prometheus/client_golang ([8001609](https://github.com/soerenschneider/vault-pki-cli/commit/8001609a60138a753ceaa3736dd82a73bca946b7))
* **deps:** bump github.com/prometheus/client_golang ([f3f600c](https://github.com/soerenschneider/vault-pki-cli/commit/f3f600c73a05c9316af554a0012115c4528696b4))
* **deps:** bump github.com/prometheus/common from 0.44.0 to 0.46.0 ([65d54e9](https://github.com/soerenschneider/vault-pki-cli/commit/65d54e948ad57b937a27e929b0afee844c0a8f28))
* **deps:** bump github.com/rs/zerolog from 1.30.0 to 1.31.0 ([5e540e9](https://github.com/soerenschneider/vault-pki-cli/commit/5e540e91f78cbf53cb8fab04d99e8a737385cbfe))
* **deps:** bump github.com/spf13/cobra from 1.7.0 to 1.8.0 ([5f5aa95](https://github.com/soerenschneider/vault-pki-cli/commit/5f5aa95432b012276aab9bce8359c50e8e11eca1))
* **deps:** bump github.com/spf13/viper from 1.16.0 to 1.17.0 ([b5be99d](https://github.com/soerenschneider/vault-pki-cli/commit/b5be99db3ee3d4a1560330255aeeab308109ca82))
* **deps:** bump golang from 1.21.1 to 1.21.2 ([00ef64b](https://github.com/soerenschneider/vault-pki-cli/commit/00ef64bf71154711f5d8989ca78aa49f919a9564))
* **deps:** bump golang from 1.21.2 to 1.21.3 ([66f62a1](https://github.com/soerenschneider/vault-pki-cli/commit/66f62a1b7a143cff35f1f3c345c85ce6e8cff3cd))
* **deps:** bump golang from 1.21.3 to 1.21.6 ([97d25e7](https://github.com/soerenschneider/vault-pki-cli/commit/97d25e7a17eeacc8c61a844a8ebc2285229b8e30))
* **deps:** bump golang from 1.21.6 to 1.22.1 ([a044a37](https://github.com/soerenschneider/vault-pki-cli/commit/a044a37fda8d7e32669f7ad8b76ae4b325e95530))
* **deps:** bump golang.org/x/crypto from 0.14.0 to 0.17.0 ([5168d47](https://github.com/soerenschneider/vault-pki-cli/commit/5168d4703eafde03000c644f228ef3fbdc51446c))
* **deps:** bump golang.org/x/net from 0.14.0 to 0.15.0 ([ddbc820](https://github.com/soerenschneider/vault-pki-cli/commit/ddbc82091b6beccf1749ec37a4111ef931e16379))
* **deps:** bump golang.org/x/net from 0.15.0 to 0.17.0 ([a1f3469](https://github.com/soerenschneider/vault-pki-cli/commit/a1f3469e2f094f867dcb07d4b41c2cbe28bdd3d9))
* **deps:** bump golang.org/x/net from 0.20.0 to 0.22.0 ([90c79c3](https://github.com/soerenschneider/vault-pki-cli/commit/90c79c35cbbc1e9b8bccdbded3c5a4520c0a9b6c))
* **deps:** bump golang.org/x/sys from 0.12.0 to 0.13.0 ([b04203d](https://github.com/soerenschneider/vault-pki-cli/commit/b04203db3739552826222d9a830e3d2f00b30d4c))
* **deps:** bump golang.org/x/sys from 0.17.0 to 0.18.0 ([439b5ef](https://github.com/soerenschneider/vault-pki-cli/commit/439b5eff8aec9ef30818ef790c00ccdaafa6fbc2))
* **deps:** bump google.golang.org/protobuf from 1.31.0 to 1.33.0 ([989817f](https://github.com/soerenschneider/vault-pki-cli/commit/989817f62e768a13c9725fda2fa7364c5f7ea4a0))
* **deps:** bump k8s.io/api from 0.29.0 to 0.29.3 ([20b59f1](https://github.com/soerenschneider/vault-pki-cli/commit/20b59f18255e560b15912df3987a484615f4c054))
* **deps:** bump k8s.io/apimachinery from 0.29.0 to 0.29.3 ([f4c1e2f](https://github.com/soerenschneider/vault-pki-cli/commit/f4c1e2f98b01644590cb0dda1a342931fb8b8eb1))
* **deps:** bump k8s.io/client-go from 0.28.2 to 0.29.0 ([c2cb27c](https://github.com/soerenschneider/vault-pki-cli/commit/c2cb27c7e0e8a4e2a4e6d546eb6d1a965bfc9462))
* **deps:** bump k8s.io/client-go from 0.29.0 to 0.29.3 ([5bc0fde](https://github.com/soerenschneider/vault-pki-cli/commit/5bc0fde150992eb68fe73eb718bf22299f336edd))
* fix compilation error with prometheus lib updates ([0934a26](https://github.com/soerenschneider/vault-pki-cli/commit/0934a264de6b5187095ba4fabed94c77e935b5b9))
* fix linting issues ([d88ceca](https://github.com/soerenschneider/vault-pki-cli/commit/d88ceca725474d9fdfcaa24ddfc7f1a65dde90c2))
* improve handling cert data containing superfluous whitespaces ([8912651](https://github.com/soerenschneider/vault-pki-cli/commit/891265161bb34de46c1215f8128ad845c1df0f1f))
* remove newlines from decoded base64 strings ([e0c8f75](https://github.com/soerenschneider/vault-pki-cli/commit/e0c8f7567f9d7cd992f30da504a9e4ee9e399183))

## [1.13.4](https://github.com/soerenschneider/vault-pki-cli/compare/v1.13.3...v1.13.4) (2023-08-05)


### Bug Fixes

* write ca to cert if no dedicated ca file is given ([ba51523](https://github.com/soerenschneider/vault-pki-cli/commit/ba515237297e131f68735b56ea8301bf2b9d5a15))

## [1.13.3](https://github.com/soerenschneider/vault-pki-cli/compare/v1.13.2...v1.13.3) (2023-08-03)


### Bug Fixes

* check if config is supplied ([44f27b7](https://github.com/soerenschneider/vault-pki-cli/commit/44f27b7586ea4435fca89429bce5f1430ac396c4))

## [1.13.2](https://github.com/soerenschneider/vault-pki-cli/compare/v1.13.1...v1.13.2) (2023-08-02)


### Bug Fixes

* prevent non-critical errors from influencing success of run ([3f07080](https://github.com/soerenschneider/vault-pki-cli/commit/3f0708094fb708f59d7e9fb77c9f50ac0b386ae1))

## [1.13.1](https://github.com/soerenschneider/vault-pki-cli/compare/v1.13.0...v1.13.1) (2023-07-15)


### Bug Fixes

* add logging adapter to use by vault ([b4027f3](https://github.com/soerenschneider/vault-pki-cli/commit/b4027f3fc28d6b0b48cf4746ed9461128e86ec75))
* fix error logic ([93c8d5d](https://github.com/soerenschneider/vault-pki-cli/commit/93c8d5d7e83a75c3a5ab59528d29298b5a7c47e7))
* only panic on actual error ([c2b3b36](https://github.com/soerenschneider/vault-pki-cli/commit/c2b3b36d1e48e71e5ab20e875a45f798f7d1cc44))
* set default mod to 0600 for files ([580b286](https://github.com/soerenschneider/vault-pki-cli/commit/580b286164699ec424bb7e73bba77156f3a8f830))

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

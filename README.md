# vault-pki-cli
- can be considered a PKI swiss knife that interacts with Hashicorp Vault
- can be used as a building block to enable zero-trust-policy
- can be easily automated

## Features
- ✓ Issuing x509 certificates
- ✓ Signing CSRs
- ✓ Fetching the PKIs CRL
- ✓ Fetching the PKIs CA chain
- ✓ Revoking your certificate
- ✓ Support for YubiKey PIV

# Subcommands

## Issuing a x509 certificate

```shell
➜  ./vault-pki-cli -t test -a http://localhost:8200 issue --common-name my.example.com -p /tmp/my.example.com.key -c /tmp/my.example.com.crt
2022-05-09T10:19:00+02:00 INF Version v1.3.0 (b509559e872e9ff75e413dd6041e882efdf8e4c6)
2022-05-09T10:19:00+02:00 INF ------------- Printing common config values -------------
2022-05-09T10:19:00+02:00 INF vault-address=http://localhost:8200
2022-05-09T10:19:00+02:00 INF vault-token=*** (sensitive output)
2022-05-09T10:19:00+02:00 INF vault-mount-pki=pki_intermediate
2022-05-09T10:19:00+02:00 INF vault-mount-approle=approle
2022-05-09T10:19:00+02:00 INF vault-pki-role-name=my_role
2022-05-09T10:19:00+02:00 INF ------------- Printing issue cmd values -------------
2022-05-09T10:19:00+02:00 INF certificate-file=/tmp/my.example.com.crt
2022-05-09T10:19:00+02:00 INF private-key-file=/tmp/my.example.com.key
2022-05-09T10:19:00+02:00 INF ttl=48h
2022-05-09T10:19:00+02:00 INF common-name=my.example.com
2022-05-09T10:19:00+02:00 INF metrics-file=/tmp/vault-pki-cli.prom
2022-05-09T10:19:00+02:00 INF force-new-certificate=false
2022-05-09T10:19:00+02:00 INF lifetime-threshold-percent=33.000000
2022-05-09T10:19:00+02:00 INF ------------- Finished printing config values -------------
2022-05-09T10:19:00+02:00 INF A certificate already exists, trying to parse it
2022-05-09T10:19:00+02:00 INF Certificate 55:1d:fc:04:80:10:52:b6:75:09:7a:3e:57:1c:36:45:cb:4e:0c:fb successfully parsed
2022-05-09T10:19:00+02:00 INF Lifetime at 15.42%, 7h24m3s left (valid from '2022-05-07 15:42:33 +0000 UTC', until '2022-05-09 15:43:03 +0000 UTC')
2022-05-09T10:19:00+02:00 INF Issuing new certificate
2022-05-09T10:19:00+02:00 INF New certificate successfully issued
2022-05-09T10:19:00+02:00 INF New certificate valid until 2022-05-11 08:19:00 +0000 UTC (48h0m0s)
2022-05-09T10:19:00+02:00 INF Cleaning up the backend...
2022-05-09T10:19:00+02:00 INF Attempting to revoke certificate 55:1d:fc:04:80:10:52:b6:75:09:7a:3e:57:1c:36:45:cb:4e:0c:fb
2022-05-09T10:19:00+02:00 INF Revoking certificate successful
2022-05-09T10:19:00+02:00 INF Dumping metrics to /tmp/vault-pki-cli.prom 
```

## Signing a CSR
```shell
➜  openssl req -new -newkey rsa:2048 -nodes -keyout /tmp/my.example.com.key -out /tmp/my.example.csr
...
➜  ./vault-pki-cli -t test -a http://localhost:8200 sign --common-name my.example.com --csr-file /tmp/my.example.com.csr -c /tmp/my.example.com.crt
2022-05-09T10:25:10+02:00 INF Version v1.3.0 (b509559e872e9ff75e413dd6041e882efdf8e4c6)
2022-05-09T10:25:10+02:00 INF ------------- Printing common config values -------------
2022-05-09T10:25:10+02:00 INF vault-address=http://localhost:8200
2022-05-09T10:25:10+02:00 INF vault-token=*** (sensitive output)
2022-05-09T10:25:10+02:00 INF vault-mount-pki=pki_intermediate
2022-05-09T10:25:10+02:00 INF vault-mount-approle=approle
2022-05-09T10:25:10+02:00 INF vault-pki-role-name=my_role
2022-05-09T10:25:10+02:00 INF ------------- Printing sign cmd values -------------
2022-05-09T10:25:10+02:00 INF csr-file=/tmp/my.example.com.csr
2022-05-09T10:25:10+02:00 INF certificate-file=/tmp/my.example.com.crt
2022-05-09T10:25:10+02:00 INF ttl=48h
2022-05-09T10:25:10+02:00 INF common-name=my.example.com
2022-05-09T10:25:10+02:00 INF metrics-file=/tmp/vault-pki-cli.prom
2022-05-09T10:25:10+02:00 INF ------------- Finished printing config values -------------
2022-05-09T10:25:10+02:00 INF Issuing new certificate
2022-05-09T10:25:10+02:00 INF CSR has been successfully signed using serial 3e:29:3a:65:38:d5:55:ee:6f:65:e4:57:29:63:7e:dd:80:30:fa:20
2022-05-09T10:25:10+02:00 INF New certificate valid until 2022-05-11 08:25:10 +0000 UTC (48h0m0s)
2022-05-09T10:25:10+02:00 INF Cleaning up the backend...
2022-05-09T10:25:10+02:00 INF Dumping metrics to /tmp/vault-pki-cli.prom

```

# Configuration
Configuration seeks for config files named `config.$ext` in the following directories:
- $HOME/.config/vault-pki-cli
- /etc/vault-pki-cli/

## Configuration Flags

### General Flags
| Name                 | Type   | Mandatory | Default          | Example               | Description                                                                          |
|----------------------|--------|-----------|------------------|-----------------------|--------------------------------------------------------------------------------------|
| vault-address        | string | yes       |                  | http://localhost:8200 | Address of the Vault server                                                          |
| vault-token          | string | no*       |                  | test                  | Token to access vault. Can not be used in conjunction with approle login.            |
| vault-role-id        | string | no*       |                  | my-vault-role         | AppRole id to login to Vault. Can not be used in conjunction with token auth.        |
| vault-secret-id      | string | no*       |                  | very-secret-id        | AppRole secret_id to login to Vault. Can not be used in conjunction with token auth. |
| vault-secret-id-file | string | no*       |                  | ~/.vault_secret_id    | File to read AppRole secret_id from. Can not be used in conjunction with token auth. |
| vault-mount-pki      | string | no        | pki_intermediate |                       | Vault path where the pki secret backend is mounted                                   |
| vault-pki-role-name  | string | no        | my_role          |                       | Name of the PKI role configured in Vault                                             |
| vault-mount-approle  | string | no        | approle          |                       | Vault path where the AppRole auth method is mounted                                  |

### Issue Subcommand
| Name                       | Type   | Mandatory | Default                                      | Description                                                                                                                                  |
|----------------------------|--------|-----------|----------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| certificate-file           | string | no*       |                                              | The file to write the certificate to. Can not be used when also specifying Yubikey Slot.                                                     |
| private-key-file           | string | no*       |                                              | The file to write the private key to. Can not be used when also specifying Yubikey Slot.                                                     |
| common-name                | string | yes       |                                              | The common-name (CN) for the x509 cert                                                                                                       |
| yubi-slot                  | uint   | no*       |                                              | Defines which [YubiKey slot](https://docs.yubico.com/yesdk/users-manual/application-piv/slots.html) to use. Uses hex format, example: 0x9a   |
| yubi-pin                   | string | no        |                                              | Pin to unlock the YubiKey slot. If no PIN is provided, the tool asks you interactively for it.                                               |
| ip-sans                    | string | no        | []                                           | Specifies the requested IP Subject Alternative Names, in a comma-delimited list                                                              |
| alt-names                  | string | no        | []                                           | Specifies the requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses.                    |
| force-new-certificate      | bool   | no        | false                                        | Flag to force issuing a new certificate, thus ignoring the `lifetime-threshold-percent` option                                               |
| lifetime-threshold-percent | float  | no        | 33.                                          | Threshold of certificate lifetime before requesting a new one                                                                                |
| ttl                        | string | no        | 48h                                          | Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used        |
| owner                      | string | no        |                                              | The owner of the written files                                                                                                               |
| group                      | string | no        |                                              | The group owner of the written files                                                                                                         |
| metrics-file               | string | no        | /var/lib/node_exporter/vault_pki_issuer.prom | File to write the prometheus metrics to                                                                                                      |

### Sign Subcommand
| Name                       | Type   | Mandatory | Default                                      | Description                                                                                                                            |
|----------------------------|--------|-----------|----------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|
| certificate-file           | string | yes       |                                              | The file to write the certificate to                                                                                                   |
| csr-file                   | string | yes       |                                              | The file to read the CSR from                                                                                                          |
| common-name                | string | yes       |                                              | The common-name (CN) for the x509 cert                                                                                                 |
| ip-sans                    | string | no        | []                                           | Specifies the requested IP Subject Alternative Names, in a comma-delimited list                                                        |
| alt-names                  | string | no        | []                                           | Specifies the requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses.              |
| ttl                        | string | no        | 48h                                          | Specifies requested Time To Live. Cannot be greater than the role's max_ttl value. If not provided, the role's ttl value will be used  |
| owner                      | string | no        |                                              | The owner of the written files                                                                                                         |
| group                      | string | no        |                                              | The group owner of the written files                                                                                                   |
| metrics-file               | string | no        | /var/lib/node_exporter/vault_pki_issuer.prom | File to write the prometheus metrics to                                                                                                |

# YubiKey PIV Support
YubiKey PIV support is based on the excellent [piv-go](https://github.com/go-piv/piv-go) library which relies on platform-dependent libraries. As it needs to be compiled using `CGO_ENABLED=1` only binaries without YubiKey support are found in the releases section.

## Build with YubiKey Support
The Makefile target `build-yubikey` leverages the go build tag `yubikey` and builds a binary with support for YubiKeys.

# Testing with Vault
The folder `assets/terraform/` contains Terraform code that spins up a local PKI to use with vault-pki-cli.

```shell
export VAULT_TOKEN=test
export VAULT_ADDR=http://localhost:8200
docker run --cap-add=IPC_LOCK -d -p 8200:8200 -e "VAULT_DEV_ROOT_TOKEN_ID=$VAULT_TOKEN" -e "VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200" vault
terraform -chdir=assets/terraform apply -auto-approve
make build
./vault-pki-cli -a $VAULT_ADDR -t $VAULT_TOKEN issue -t test -c /tmp/test.crt -p /tmp/test.key --common-name bla.example.com
```

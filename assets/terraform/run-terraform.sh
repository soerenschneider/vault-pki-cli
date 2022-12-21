#!/bin/sh

export VAULT_TOKEN=vault-pki-cli
export VAULT_ADDR=http://localhost:8200

terraform init -upgrade
terraform apply

#!/bin/sh

TOKEN=vault-pki-cli

docker run -d -e VAULT_DEV_ROOT_TOKEN_ID="${TOKEN}" -p 8200:8200 vault

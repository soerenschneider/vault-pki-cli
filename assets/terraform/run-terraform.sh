#!/bin/sh

echo "kubernetes_ca_cert=<<EOT" > terraform.tfvars
kubectl get secrets vault-kubernetes-auth-secret -o json | jq -r '.data["ca.crt"]' | base64 -d >> terraform.tfvars
echo "EOT" >> terraform.tfvars

echo "kubernetes_jwt=\"$(kubectl get secrets -n default vault-kubernetes-auth-secret -o json | jq -r '.data.token' | base64 -d)\"" >> terraform.tfvars

export VAULT_TOKEN=vault-pki-cli
export VAULT_ADDR=http://localhost:8200

terraform init -upgrade
terraform apply -auto-approve
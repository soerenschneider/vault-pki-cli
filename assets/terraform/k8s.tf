resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
}

resource "vault_kubernetes_auth_backend_config" "example" {
  backend                = vault_auth_backend.kubernetes.path
  kubernetes_host        = "https://kubernetes.default.svc:443"
  kubernetes_ca_cert     = var.kubernetes_ca_cert
  token_reviewer_jwt     = var.kubernetes_jwt
  disable_iss_validation = true
}

resource "vault_kubernetes_auth_backend_role" "example" {
  backend                          = vault_auth_backend.kubernetes.path
  role_name                        = "vault-pki"
  bound_service_account_names      = ["*"]
  bound_service_account_namespaces = ["default"]
  token_ttl                        = 3600
  token_policies                   = ["default", "developer"]
  audience                         = "https://kubernetes.default.svc.cluster.local"
}

resource "vault_policy" "bla" {
  name = "developer"

  policy = <<EOT
path "kubernetes/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "pki_intermediate/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

EOT
}
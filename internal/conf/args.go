package conf

const (
	FLAG_CONFIG_FILE = "config-file"

	FLAG_VAULT_ADDRESS                  = "vault-address"
	FLAG_VAULT_TOKEN                    = "vault-token"
	FLAG_VAULT_ROLE_ID                  = "vault-role-id"
	FLAG_VAULT_SECRET_ID                = "vault-secret-id"
	FLAG_VAULT_SECRET_ID_FILE           = "vault-secret-id-file"
	FLAG_VAULT_MOUNT_PKI                = "vault-mount-pki"
	FLAG_VAULT_MOUNT_PKI_DEFAULT        = "pki_intermediate"
	FLAG_VAULT_PKI_BACKEND_ROLE         = "vault-pki-role-name"
	FLAG_VAULT_PKI_BACKEND_ROLE_DEFAULT = "my_role"

	FLAG_VAULT_MOUNT_APPROLE         = "vault-mount-approle"
	FLAG_VAULT_MOUNT_APPROLE_DEFAULT = "approle"

	FLAG_CERTIFICATE_FILE = "certificate-file"

	FLAG_ISSUE_FORCE_NEW_CERTIFICATE                 = "force-new-certificate"
	FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE         = "lifetime-threshold-percent"
	FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT = 33.
	FLAG_ISSUE_PRIVATE_KEY_FILE                      = "private-key-file"
	FLAG_ISSUE_CA_CHAIN_FILE                         = "ca-chain-file"
	FLAG_ISSUE_TTL                                   = "ttl"
	FLAG_ISSUE_IP_SANS                               = "ip-sans"
	FLAG_ISSUE_COMMON_NAME                           = "common-name"
	FLAG_ISSUE_ALT_NAMES                             = "alt-names"
	FLAG_ISSUE_METRICS_FILE                          = "metrics-file"
	FLAG_ISSUE_METRICS_FILE_DEFAULT                  = "/var/lib/node_exporter/ssh_key_sign.prom"
)

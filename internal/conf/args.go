package conf

const (
	FLAG_CONFIG_FILE = "config"
	FLAG_DEBUG       = "debug"

	FLAG_VAULT_ADDRESS                     = "vault-address"
	FLAG_VAULT_AUTH_TOKEN                  = "vault-auth-token"
	FLAG_VAULT_AUTH_IMPLICIT               = "vault-auth-implicit"
	FLAG_VAULT_AUTH_K8S_ROLE               = "vault-auth-k8s"
	FLAG_VAULT_AUTH_APPROLE_ID             = "vault-auth-role-id"
	FLAG_VAULT_AUTH_APPROLE_SECRET_ID      = "vault-auth-secret-id"
	FLAG_VAULT_AUTH_APPROLE_SECRET_ID_FILE = "vault-auth-secret-id-file"
	FLAG_VAULT_APPROLE_MOUNT               = "vault-approle-mount"
	FLAG_VAULT_PKI_MOUNT                   = "vault-pki-mount"
	FLAG_VAULT_PKI_BACKEND_ROLE            = "vault-pki-role-name"
	FLAG_VAULT_MOUNT_KV2                   = "vault-kv2-mount"

	FLAG_ISSUE_FORCE_NEW_CERTIFICATE         = "force-new-certificate"
	FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE = "lifetime-threshold-percent"
	FLAG_ISSUE_PRIVATE_KEY_FILE              = "private-key-file"
	FLAG_ISSUE_BACKEND_CONFIG                = "backend-config"
	FLAG_READACME_ACME_PREFIX                = "acme-prefix"

	FLAG_ISSUE_YUBIKEY_SLOT = "yubi-slot"
	FLAG_ISSUE_YUBIKEY_PIN  = "yubi-pin"

	FLAG_ISSUE_TTL          = "ttl"
	FLAG_ISSUE_DAEMONIZE    = "daemonize"
	FLAG_ISSUE_IP_SANS      = "ip-sans"
	FLAG_ISSUE_COMMON_NAME  = "common-name"
	FLAG_ISSUE_ALT_NAMES    = "alt-names"
	FLAG_METRICS_FILE       = "metrics-file"
	FLAG_ISSUE_METRICS_ADDR = "metrics-addr"
	FLAG_ISSUE_HOOKS        = "hooks"

	FLAG_OUTPUT_FILE = "output-file"
	FLAG_DER_ENCODED = "der-encoding"

	FLAG_CERTIFICATE_FILE = "certificate-file"
	FLAG_CA_FILE          = "ca-file"
	FLAG_CSR_FILE         = "csr-file"
	FLAG_FILE_OWNER       = "owner"
	FLAG_FILE_GROUP       = "group"
)

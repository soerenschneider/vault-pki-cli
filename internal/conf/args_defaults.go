package conf

import "math"

const (
	FLAG_VAULT_PKI_BACKEND_ROLE_DEFAULT              = "my_role"
	FLAG_VAULT_MOUNT_APPROLE_DEFAULT                 = "approle"
	FLAG_ISSUE_LIFETIME_THRESHOLD_PERCENTAGE_DEFAULT = 33.
	FLAG_ISSUE_TTL_DEFAULT                           = "48h"
	FLAG_FILE_OWNER_DEFAULT                          = "root"

	FLAG_ISSUE_YUBIKEY_SLOT_DEFAULT = math.MaxUint32

	FLAG_VAULT_MOUNT_PKI_DEFAULT    = "pki_intermediate"
	FLAG_ISSUE_METRICS_FILE_DEFAULT = "/var/lib/node_exporter/vault_pki_issuer.prom"
	FLAG_ISSUE_METRICS_ADDR_DEFAULT = ":9172"
)

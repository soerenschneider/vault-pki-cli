package conf

import "testing"

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		VaultAddress                           string
		VaultAuthMethod                        string
		VaultToken                             string
		VaultAuthK8sRole                       string
		VaultRoleId                            string
		VaultSecretId                          string
		VaultSecretIdFile                      string
		VaultMountApprole                      string
		VaultMountPki                          string
		VaultMountKv2                          string
		VaultPkiRole                           string
		Daemonize                              bool
		CommonName                             string
		Ttl                                    string
		IpSans                                 []string
		AltNames                               []string
		AcmePrefix                             string
		MetricsFile                            string
		MetricsAddr                            string
		ForceNewCertificate                    bool
		StorageConfig                          []map[string]string
		PostHooks                              []string
		CertificateLifetimeThresholdPercentage float32
		DerEncoded                             bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				VaultAddress:                           "https://vault:8200",
				VaultAuthMethod:                        "approle",
				VaultToken:                             "",
				VaultAuthK8sRole:                       "",
				VaultRoleId:                            "approle-role-id",
				VaultSecretId:                          "",
				VaultSecretIdFile:                      "asd",
				VaultMountApprole:                      "approle",
				VaultMountPki:                          "pki",
				VaultMountKv2:                          "secret",
				VaultPkiRole:                           "human",
				Daemonize:                              false,
				CommonName:                             "",
				Ttl:                                    "",
				IpSans:                                 nil,
				AltNames:                               nil,
				AcmePrefix:                             "",
				MetricsFile:                            "",
				MetricsAddr:                            "",
				ForceNewCertificate:                    false,
				StorageConfig:                          nil,
				PostHooks:                              nil,
				CertificateLifetimeThresholdPercentage: 0,
				DerEncoded:                             false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				VaultAddress:                           tt.fields.VaultAddress,
				VaultAuthMethod:                        tt.fields.VaultAuthMethod,
				VaultToken:                             tt.fields.VaultToken,
				VaultAuthK8sRole:                       tt.fields.VaultAuthK8sRole,
				VaultRoleId:                            tt.fields.VaultRoleId,
				VaultSecretId:                          tt.fields.VaultSecretId,
				VaultSecretIdFile:                      tt.fields.VaultSecretIdFile,
				VaultMountApprole:                      tt.fields.VaultMountApprole,
				VaultMountPki:                          tt.fields.VaultMountPki,
				VaultMountKv2:                          tt.fields.VaultMountKv2,
				VaultPkiRole:                           tt.fields.VaultPkiRole,
				Daemonize:                              tt.fields.Daemonize,
				CommonName:                             tt.fields.CommonName,
				Ttl:                                    tt.fields.Ttl,
				IpSans:                                 tt.fields.IpSans,
				AltNames:                               tt.fields.AltNames,
				AcmePrefix:                             tt.fields.AcmePrefix,
				MetricsFile:                            tt.fields.MetricsFile,
				MetricsAddr:                            tt.fields.MetricsAddr,
				ForceNewCertificate:                    tt.fields.ForceNewCertificate,
				StorageConfig:                          tt.fields.StorageConfig,
				PostHooks:                              tt.fields.PostHooks,
				CertificateLifetimeThresholdPercentage: tt.fields.CertificateLifetimeThresholdPercentage,
				DerEncoded:                             tt.fields.DerEncoded,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

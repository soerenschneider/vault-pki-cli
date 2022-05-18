package conf

import (
	"testing"
)

func TestIssueArguments_Validate(t *testing.T) {
	tests := []struct {
		name       string
		conf       IssueArguments
		wantErrors bool
	}{
		{
			name: "bare minimum",
			conf: IssueArguments{
				CommonName:                             "my.common.name",
				CertificateLifetimeThresholdPercentage: 10.,
				YubikeySlot:                            FLAG_ISSUE_YUBIKEY_SLOT_DEFAULT,
				PrivateKeyFile:                         "/tmp/cert.key",
				CertificateFile:                        "/tmp/cert.crt",
			},
			wantErrors: false,
		},
		{
			name: "empty",
			conf: IssueArguments{
				YubikeySlot: FLAG_ISSUE_YUBIKEY_SLOT_DEFAULT,
			},
			wantErrors: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.conf.Validate(); len(got) == 0 == tt.wantErrors {
				t.Errorf("Validate() = %t, want %t", len(got) > 0, tt.wantErrors)
			}
		})
	}
}

package pkg

import (
	"crypto/x509"
	"reflect"
	"testing"
)

var exampleContainer = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAu3ajfQZjZvpsZkAtFDbaS3LBigdV5lBpnBDdPtE8HhbHZQmH
EGXrIwE6vdRtPGsoaOssiJADezi/RV1QmtBpVMUvSU43U13EtPTpfqVKozX/hSYs
mlJLkxJ0LcUvWfyPbznmtXlcHMWCY/Lwj1ds4WPRHsRPXGsIVbhzfGBX4v8RaKc3
VnRKhoE8KOhp4fYTZQxXp/ucDaAML7Xt/et8XRUM0zOSMLqIy4p7XzopQfVeCoG1
BNpzceFMV0VpYiF5rF9oW29CfdgAotA3Uf9cq+Fy1U9kRn8GubVi2nwOVMdfTSAw
nJAoGkjtznx/wh6lDYU/I82veXMWMmm5vBZh8wIDAQABAoIBABBwV+fXzpGyNh1F
VW6nXL8vAf/LouG+fXRdGjmu+Xmd/8BBdKGgfl0kd3U8EpQwxWtl7BLRpiyBDmzT
wQTCb+oqHHpuLHXYDC7eJzee4Qus6YpQjaq+urfb72owF3Xpqt5TEoMpcEVpoISJ
QkUfooGlUipDhr4Q+LsjoKTwgeR62Po/qP8W3lgwDhqoXacY2nRtD02Z49FCuxdn
fN+IhhhFxsP+wkCF2dEqNgfRNvvYIDMDd8nRPNNmYTUE5iavX5ITJ/w0fULFnhR6
8Qc5O8tyiMtItQImCxE5EIGwftYxndpbeeCbxFRNIvaGSwu3IPLMsn2PvZMgfo8l
SoZnVDkCgYEA8oW1roLgYtBimqTNR+M+L9jvPv95Cmpo2mUTkClsEiGfJ9AfL8vj
nsKatv3urmXb1Cyz/XJ33hTcNK0OqYb4pEzHXKPMDySpUfpychmIrAlayOIgf7SS
LNcFuTx9l1Sh76coyUrW5FefyFywN1vbmrcdT/n1Y076K40cVor2oI0CgYEAxeGf
j5Bdinnr7PzvgmoMlr1EtXi/zLKPJ7fmSQbQ3BkqqApFxhKxqNIC6am2zDI+r+bR
BbDlsfo1wMjPwPz9hj3XTAYxlxKsTkILpr6TsB1SZt+Poat94USifULmQORrR9OF
cnOjZVgiB5HnS9Q8VD2nnuHCxrcNnyhzCurWrH8CgYEAqQxO0e/kXLyInubVOJDL
3ipGyhDl3D7EC8d81XYqIJFTETtfIb/rT9SyZ2+lmebiTolChR3vM9wyin0+xSiR
1GS4ani6WqvhYoVClQn7XH/Aylnk8V96rMrM8Iubt4qEvjo0kesa01vIwq7pHg1n
i/ar9f1z8N8yPn1EDYcb1lkCgYAGhNNL2HasbC3QhdiiFDpL8PpFfC/dX3iF13IX
r8jLp2yXUpdP2ifOJvT/m56xBWq5QsJaDKTUgyioLDVj5zG27WydTYrurifNADIA
EUEuSRkA2JaTveGMvUUZGU4ajyvVlutLhPG6Efg1BaJ4Bgriv5E5E7jl8Pva5Ws8
zdW6owKBgQDpGE1bneb9XFa9pontvVBVelNc7VFFUeqMXikfQhlLZigblMl87AZm
weOx2A3XCQ2ckLqKTxotTqzhe35CaRMgw8LjUKQOEJHlW8bp3RGmjixDv6dxFW6A
wy6YJqxiAD12fr/k40qOtvRMQ+mq1ucjWDdRanPhNnQvjue7AbAzOg==
-----END RSA PRIVATE KEY-----
-----BEGIN CERTIFICATE-----
MIIDaTCCAlGgAwIBAgIUFYAvaAVGHzV8HdK+w30oToUVbQkwDQYJKoZIhvcNAQEL
BQAwPDEYMBYGA1UEChMPTXkgb3JnYW5pemF0aW9uMQ4wDAYDVQQLEwVNeSBPVTEQ
MA4GA1UEAxMHUm9vdCBDQTAeFw0yMjA1MTYxMzQ3MzRaFw0yMjA1MjYxMzQ4MDRa
MDwxGDAWBgNVBAoTD015IG9yZ2FuaXphdGlvbjEOMAwGA1UECxMFTXkgT1UxEDAO
BgNVBAMTB1Jvb3QgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDb
dEMIHYXRU4OnmKEm3fJ9PHJK7MW1zwOHAshc9YfGXjRRECvhJ41HjfulXFPminJx
ovGESvA3/bG7f5eMAn4DnOzaoInbagSzhmcazE65GfngA5n+IrSjd4HTGejGRn8t
dep2TyLBASbZGoThgD5mAn5Zp+Zi0R7w0byJAS+VxriHMg9hn4U0cv1Pq7W5BnfJ
1di54nigosqiRCOrS1tC3KC6I0XE8GSHfzBe6vr5K2zrCB6uEIn1f2hGdREmIn5V
lMi9nDiv0mLlKwfRnnhdPBmVZaY0Ae9Xgl8luUuzkpPPsZis6GB1NVCLR7A/O7g7
SqpZpdycHX7Nt6CI+uHjAgMBAAGjYzBhMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMB
Af8EBTADAQH/MB0GA1UdDgQWBBRboboMDFISBt0nXJntRNrNIm1C5jAfBgNVHSME
GDAWgBRboboMDFISBt0nXJntRNrNIm1C5jANBgkqhkiG9w0BAQsFAAOCAQEAr0Ml
QT63ISVvLGJ/+qo7fiKrH5K6nbo4SB/1Xs0LJ5obq9fO15Li5vCOjAJCxF1+uSwU
jr7gxsniG2HmJZK4E4HankpUgo7mIwmCwhsYTe3cwt32HChDnyKuGFvItSww86FA
GWOSqdtAZukwjEZVWGlBUSRSOjtLFAG+pEQQFgujF0HthEsCbpkKFZAHsunlaaI1
fldqmdCjULArypPet6ATtqcK9V7kjz/QJIW28Atmbcn7tVWFxtmnHvKEyLaZKuWo
MUHZTO5tnRxeowi/YbAdgTuTEReCaITgCmN7Jz+fCrAAWB+uISN5tqrfYG2q2oXz
gNEPI1a9Ltgk/WIkiw==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDcTCCAlmgAwIBAgIUEILoLSZaM2LAFDe2IrBMEFOxry4wDQYJKoZIhvcNAQEL
BQAwPDEYMBYGA1UEChMPTXkgb3JnYW5pemF0aW9uMQ4wDAYDVQQLEwVNeSBPVTEQ
MA4GA1UEAxMHUm9vdCBDQTAeFw0yMjA1MTYxNDIyMjFaFw0yMjA1MTgxNDIyNTBa
MBkxFzAVBgNVBAMTDm15LmV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAu3ajfQZjZvpsZkAtFDbaS3LBigdV5lBpnBDdPtE8HhbHZQmH
EGXrIwE6vdRtPGsoaOssiJADezi/RV1QmtBpVMUvSU43U13EtPTpfqVKozX/hSYs
mlJLkxJ0LcUvWfyPbznmtXlcHMWCY/Lwj1ds4WPRHsRPXGsIVbhzfGBX4v8RaKc3
VnRKhoE8KOhp4fYTZQxXp/ucDaAML7Xt/et8XRUM0zOSMLqIy4p7XzopQfVeCoG1
BNpzceFMV0VpYiF5rF9oW29CfdgAotA3Uf9cq+Fy1U9kRn8GubVi2nwOVMdfTSAw
nJAoGkjtznx/wh6lDYU/I82veXMWMmm5vBZh8wIDAQABo4GNMIGKMA4GA1UdDwEB
/wQEAwIDqDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYE
FK2p0aKFLKFHT60krNUUNbGShPL2MB8GA1UdIwQYMBaAFFuhugwMUhIG3Sdcme1E
2s0ibULmMBkGA1UdEQQSMBCCDm15LmV4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUA
A4IBAQDVhkuMVUStI9IysE1uWyi4qFy5gyIzCCa1NRCd6I/NF3dBZEmCSFCLAc9f
uYFbKVwRru2pDDpwpaZKUGiYlqHfF8FlJkBoYUR6Znz4BSwZu6BFssHR10CLgNzU
SoZPhZzKElcOXq4FMVNOqU5l3606QOepcyByLCKiprkk8M8BD5y8GYvjrS3499fo
ffs/ybD5Nfu1CVXImDpy+kzWhjkVFoIkc8Xh2BXIJMZUqmYLodnaofDiqnV9+2CY
+8m6bBsKvtqr2Tm7gaKgx+M8fXIor3YSlcjOaiIkhMFkCEzurJSyyGpGf33L/tAT
e2DNjIh06nAhVh+zlgh5L7AJINrm
-----END CERTIFICATE-----`

func TestPemContainerBackend_Read_TableTest(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *x509.Certificate
		wantErr bool
	}{
		{
			name:    "test empty data",
			data:    []byte(""),
			want:    nil,
			wantErr: true,
		},
		{
			name: "test garbage data",
			data: []byte(`-----BEGIN CERTIFICATE-----
GARBAGE
-----END CERTIFICATE-----`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "edge case",
			data:    nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCertPem(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPemContainerBackend_Read(t *testing.T) {
	cert, err := ParseCertPem([]byte(exampleContainer))
	if err != nil {
		t.Fatalf("did not expect error: %v", err)
	}

	if cert.IsCA {
		t.Fatalf("expected no ca")
	}

	if cert.Subject.CommonName != "my.example.com" {
		t.Fatalf("expected 'my.example.com' as CN, got '%v'", cert.Subject.CommonName)
	}
}

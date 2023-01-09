package sink

import (
	"fmt"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/storage"
	"reflect"
	"testing"
)

const (
	privKey = "--- START PRIVATE KEY ---\nSECRETSECRET\n--- END PRIVATE KEY ---"
	cert    = "--- START CERT ---\nCERTCERTCERT\n--- END CERT ---"
	ca      = "--- START CA ---\nCACACACACA\n--- END CA ---"
)

func TestKeyPairSink_WriteCert(t *testing.T) {
	type fields struct {
		ca         pki.StorageImplementation
		cert       pki.StorageImplementation
		privateKey pki.StorageImplementation
	}
	tests := []struct {
		name     string
		fields   fields
		certData *pki.CertData
		wantErr  bool
		wantData string
	}{
		{
			name: "write ca, cert and key to single file",
			certData: &pki.CertData{
				PrivateKey:  []byte(privKey),
				Certificate: []byte(cert),
				CaData:      []byte(ca),
				Csr:         nil,
			},
			fields: fields{
				ca:         nil,
				cert:       nil,
				privateKey: &storage.BufferPod{},
			},
			wantErr:  false,
			wantData: fmt.Sprintf("%s\n%s\n%s\n", ca, cert, privKey),
		},
		{
			name: "write cert and key to single file",
			certData: &pki.CertData{
				PrivateKey:  []byte(privKey),
				Certificate: []byte(cert),
				CaData:      nil,
				Csr:         nil,
			},
			fields: fields{
				ca:         nil,
				cert:       nil,
				privateKey: &storage.BufferPod{},
			},
			wantErr:  false,
			wantData: fmt.Sprintf("%s\n%s\n", cert, privKey),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &KeyPairSink{
				ca:         tt.fields.ca,
				cert:       tt.fields.cert,
				privateKey: tt.fields.privateKey,
			}
			if err := f.WriteCert(tt.certData); (err != nil) != tt.wantErr {
				t.Errorf("WriteCert() error = %v, wantErr %v", err, tt.wantErr)
			}
			read, err := tt.fields.privateKey.Read()
			if err != nil {
				t.Errorf("Error reading b")
			}

			if !reflect.DeepEqual(string(read), tt.wantData) {
				t.Errorf("KeyPairSinkFromConfig() got = %v, want %v", string(read), tt.wantData)
			}
		})
	}
}

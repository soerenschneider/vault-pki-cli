package sink

import (
	"crypto/x509"
	"fmt"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"reflect"
	"testing"
)

type dummyImpl struct {
}

func (b dummyImpl) Write(cert *pki.CertData) error {
	fmt.Println("Write")
	return nil
}

func (b dummyImpl) Read() (*x509.Certificate, error) {
	return nil, nil
}

var backends []pki.CertSink = []pki.CertSink{
	dummyImpl{},
	dummyImpl{},
}

func TestNewMultiBackend(t *testing.T) {
	type args struct {
		backends []pki.CertSink
	}
	tests := []struct {
		name    string
		args    args
		want    *MultiSink
		wantErr bool
	}{
		{
			name:    "empty backend",
			args:    args{backends: []pki.CertSink{}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty backend",
			args: args{backends: backends},
			want: &MultiSink{
				sinks: backends,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMultiSink(tt.args.backends...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMultiSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMultiSink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiBackend_Write(t *testing.T) {
	type fields struct {
		backends []pki.CertSink
	}
	type args struct {
		certData *pki.CertData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "multiple",
			fields: fields{backends},
			args: args{
				&pki.CertData{
					PrivateKey:  []byte("test"),
					Certificate: nil,
					CaChain:     nil,
					Csr:         nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &MultiSink{
				sinks: tt.fields.backends,
			}
			if err := b.Write(tt.args.certData); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

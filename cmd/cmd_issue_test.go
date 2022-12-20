package main

import (
	"github.com/soerenschneider/vault-pki-cli/internal/conf"
	"github.com/soerenschneider/vault-pki-cli/internal/pki"
	"github.com/soerenschneider/vault-pki-cli/internal/sink"
	"reflect"
	"testing"
)

func Test_buildPemBackend(t *testing.T) {
	type args struct {
		config conf.Config
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				conf.Config{
					IssueArguments: conf.IssueArguments{
						Backends: []conf{
							{
								CertificateFile: "/tmp/cert.crt",
								PrivateKeyFile:  "/tmp/cert.key",
							},
							{
								CertificateFile: "/tmp/cert-1.crt",
								PrivateKeyFile:  "/tmp/cert-1.key",
							},
						},
					},
				},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildFilesystemSink(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildFilesystemSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.WriteCert(&pki.CertData{
				PrivateKey:  []byte("private"),
				Certificate: []byte("cert"),
				CaChain:     []byte("chain"),
				Csr:         nil,
			})
		})
	}
}

func Test_parseUri(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		want    pki.IssueSink
		wantErr bool
	}{
		{
			name: "happy path",
			uri:  "k8s://namespace/secret-name",
			want: &sink.K8sSink{
				Namespace:  "namespace",
				SecretName: "secret-name",
			},
			wantErr: false,
		},
		{
			name:    "additional path",
			uri:     "k8s://namespace/secret-name/theresmore",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unknown scheme",
			uri:     "bla://namespace/secret-name/theresmore",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUri(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUri() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseUri() got = %v, want %v", got, tt.want)
			}
		})
	}
}

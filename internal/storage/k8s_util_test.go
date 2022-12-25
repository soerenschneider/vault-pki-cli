package storage

import (
	"reflect"
	"testing"
)

func TestK8sConfigFromUri(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		want    *K8sConfig
		wantErr bool
	}{
		{
			name: "simple",
			uri:  "k8s-sec:///namespace/name",
			want: &K8sConfig{
				Namespace: "namespace",
				Name:      "name",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := K8sConfigFromUri(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("K8sConfigFromUri() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("K8sConfigFromUri() got = %v, want %v", got, tt.want)
			}
		})
	}
}
